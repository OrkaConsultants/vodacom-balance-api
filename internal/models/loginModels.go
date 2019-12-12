package models

type LoginRequestStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseStruct struct {
	Messages   []Message `json:"messages"`
	Result     Result    `json:"result"`
	Successful bool      `json:"successful"`
	Code       int64     `json:"code"`
}

type Message struct {
	Message   string `json:"message"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Severity  string `json:"severity"`
}

type Result struct {
	ExceptionMessage      string      `json:"exceptionMessage"`
	Type                  string      `json:"type"`
	Result                *bool       `json:"result"`
	ValidationfieldErrors interface{} `json:"validationfieldErrors"`
}
