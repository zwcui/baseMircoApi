package models

//app弹出提示
type AlertMessage struct {
	AlertCode			string		`description:"提示信息码，forward开头表示跳转actionurl" json:"alertCode"`
	AlertMessage		string		`description:"提示信息" json:"alertMessage"`
}

type SystemConfig struct {
	RId          		int64  		`description:"配置编号" json:"rId" xorm:"pk autoincr"`
	Description  		string 		`description:"描述" json:"description"`
	Program      		string 		`description:"配置参数" json:"program"`
	ProgramValue 		string 		`description:"参数的值" json:"programValue"`
	ProgramExpireTime 	int64 		`description:"参数过期时间" json:"programExpireTime"`
	Created      		int64  		`description:"创建时间" json:"created" xorm:"created"`
	DeletedAt    		int64  		`description:"删除时间" json:"-" xorm:"deleted"`
}

//预授权码
type SystemAuthCode struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	PreAuthCode			string			`description:"预授权码" json:"preAuthCode"`
	ExpireTime			int64			`description:"预授权码过期时间" json:"expireTime"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//接口调用凭据
type SystemAuthAccessToken struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	AuthCodeId       	int64			`description:"authCodeId" json:"authCodeId"`
	AuthorizerAccessToken string		`description:"接口调用凭据" json:"authorizerAccessToken"`
	ExpireTime			int64			`description:"接口调用凭据过期时间" json:"expireTime"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//后台操作记录
type OperationRecord struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	AuthInfoId			int64			`description:"公众号id" json:"authInfoId"`
	Operator			string			`description:"操作人" json:"operator"`
	RequestUrl			string			`description:"请求地址" json:"requestUrl"`
	RequestMethod		string			`description:"请求方法" json:"requestMethod"`
	RequestRemoteAddr	string			`description:"请求来源" json:"requestRemoteAddr"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//--------------------结构体------------------------
type FileServerResponse struct {
	Header 				FileServerResponseHeader 			`json:"header"`
	Data   				[]FileServerResponseData 			`json:"data"`
}

type FileServerResponseHeader struct {
	Code        		int 			`json:"code"`
	Description 		string 			`json:"description"`
}

type FileServerResponseData struct {
	Uri					string			`json:"uri"`
	Size				int64			`json:"size"`
	FileType			string			`json:"fileType"`
}

type WcpCommonReturnData struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type WcpAuthTokenData struct {
	WcpCommonReturnData
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

type WcpRefreshTokenData struct {
	WcpCommonReturnData
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WcpRefreshTicketData struct {
	WcpCommonReturnData
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

//-----------------------常量---------------------------
const (
	wCP_HOST       = "api.weixin.qq.com"
	wCP_HOST_DWQW       = "api.weixin.qq.com"
	wCP_APP_ID     = "wxd777a0c26a6427ee"
	wCP_APP_ID_DWQW     = "wx2227f529deb9c09d"
	wCP_APP_SECRET = "074e65135a3b66d2bad791bb85f0993a"
	wCP_APP_SECRET_DWQW = "de073002b556d729d5926db36ffcd591"
)