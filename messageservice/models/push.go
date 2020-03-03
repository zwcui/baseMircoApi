package models





var (
   	UmengAuthToken              string
	HwAccessToken               string
)


type XmRequestResponse struct {
	Result      string `json:"result"`
	Description string `json:"description"`
	Code        int    `json:"code"`
	Info        string `json:"info"`
	Reason      string `json:"reason"`
}

type HwTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type HwRequestResponse struct {
	Resultcode int    `json:"resultcode"`
	Message    string `json:"message"`
}


type MeizuPushResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UmengPushResponseData struct {
	ErrorCode string `json:"error_code"`
}

type UmengPushResponse struct {
	Ret  string                `json:"ret"`
	Data UmengPushResponseData `json:"data"`
}







