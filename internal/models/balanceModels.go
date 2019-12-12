package models

type BalanceResponse struct {
	Successful bool          `json:"successful"`
	Messages   []interface{} `json:"messages"`
	Result     BalResult     `json:"result"`
	Code       int64         `json:"code"`
}

type BalResult struct {
	ServiceList []ServiceList `json:"serviceList"`
	ShowBuyMore bool          `json:"showBuyMore"`
	RefreshTime string        `json:"refreshTime"`
}

type ServiceList struct {
	ShowDial      bool    `json:"showDial"`
	ShowBuyButton bool    `json:"showBuyButton"`
	Remaining     float64 `json:"remaining"`
	Total         float64 `json:"total"`
	PreText       string  `json:"preText"`
	MidText       string  `json:"midText"`
	PostText      string  `json:"postText"`
}
