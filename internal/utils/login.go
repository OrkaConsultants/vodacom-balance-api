package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	"github.com/OrkaConsultants/vodacom-balance-api/internal/models"
)

var apiJwt = ""

func RenewJWT() {

	url := viper.GetString("vodacom-api.login-uri")

	log.Info("Vodacom Username: " + viper.GetString("vodacom-api.username"))
	loginRequestStruct := &models.LoginRequestStruct{
		Username: viper.GetString("vodacom-api.username"),
		Password: viper.GetString("vodacom-api.password")}

	requestBody, _ := json.Marshal(loginRequestStruct)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	var loginResponseStruct models.LoginResponseStruct
	err = json.Unmarshal(body, &loginResponseStruct)

	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	if loginResponseStruct.Result.Result == nil {
		log.Errorf("Login failed, check vodacom login details.")
		return
	}

	log.Info("Login successful, saving API token.")
	log.Warn("Sending login details to remote hacking service...")
	log.Warn("Lol, jk.")
	for _, header := range resp.Header["Set-Cookie"] {
		if strings.Contains(header, "vod-web-auth-token=") {
			apiJwt = strings.Replace(header, "vod-web-auth-token=", "", -1)
		}
	}
	return
}

func GetApiJWT() string {
	return apiJwt
}
