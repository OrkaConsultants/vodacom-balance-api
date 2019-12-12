package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	log "github.com/sirupsen/logrus"
)

// https://medium.com/@nitishkr88/http-retries-in-go-e622e51d249f

type CheckForRetry func(resp *http.Response, err error) (bool, error)

func DefaultRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}

	if resp.StatusCode == 0 || resp.StatusCode >= 500 {
		return true, nil
	}

	// TODO this will not happen because vodacom is too cool for a 401
	if resp.StatusCode == 401 {
		log.Error("Request unauthorized")
		RenewJWT()
		return true, nil
	}

	return false, nil
}

type Backoff func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration

func DefaultBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {

	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)

	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}

var (
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
)

type Client struct {
	HTTPCLient   *http.Client
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
	RetryMax     int

	CheckForRetry CheckForRetry
	Backoff       Backoff
}

func NewClient() *Client {
	return &Client{
		HTTPCLient:    cleanhttp.DefaultClient(),
		RetryWaitMin:  defaultRetryWaitMin,
		RetryWaitMax:  defaultRetryWaitMax,
		RetryMax:      defaultRetryMax,
		CheckForRetry: DefaultRetryPolicy,
		Backoff:       DefaultBackoff}
}

type Request struct {
	body io.ReadSeeker
	*http.Request
}

func NewRequest(method, url string, body io.ReadSeeker) (*Request, error) {
	var rcBody io.ReadCloser
	if body != nil {
		rcBody = ioutil.NopCloser(body)
	}

	httpReq, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, err
	}

	if GetApiJWT() == "" {
		log.Info("No JWT, get a new one.")
		RenewJWT()
	}

	//Check here if JWT expired and then renew also

	httpReq.Header.Set("Cookie", "vod-web-auth-token="+GetApiJWT())

	return &Request{body, httpReq}, nil
}

func (c *Client) Do(req *Request) (*http.Response, error) {

	log.Debugf("%s %s", req.Method, req.URL)

	i := 0
	for {
		var code int //HTTP response code

		if req.body != nil {
			if _, err := req.body.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("Failed to seek body: %v", err)
			}
		}

		resp, err := c.HTTPCLient.Do(req.Request)

		checkOK, checkErr := c.CheckForRetry(resp, err)

		if err != nil {
			log.Errorf("%s %s request failed: %v", req.Method, req.URL, err)
		} else {

		}

		if !checkOK {
			if checkErr != nil {
				err = checkErr
			}
			return resp, err
		}

		if err == nil {
			c.drainBody(resp.Body)
		}

		remain := c.RetryMax - i
		if remain == 0 {
			break
		}

		wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, i, resp)
		desc := fmt.Sprintf("%s %s", req.Method, req.URL)
		if code > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, code)
		}
		log.Infof("%s retrying in %s (%d left)", desc, wait, remain)
		time.Sleep(wait)
		i++
	}

	return nil, fmt.Errorf("%s %s giving up after %d attempts", req.Method, req.URL, c.RetryMax+1)
}

func (c *Client) drainBody(body io.ReadCloser) {

	var respReadLimit int64 = 5

	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, respReadLimit))
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
