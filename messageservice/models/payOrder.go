/*
@Time : 2019/6/5 下午3:10 
@Author : zwcui
@Software: GoLand
*/
package models

import (
	"encoding/xml"
	"sort"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"strconv"
)

//付费订单
type PayOrder struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	OrderType				int					`description:"订单类型，1为购买视频" json:"orderType"`
	AppType				    int					`description:"应用类型，0为公众号 1为app" json:"appType" xorm:"notnull default 0"`
	OrderTypeId				int64				`description:"订单对应id" json:"orderTypeId"`
	LastShareAuthorId		int64				`description:"上级分享人id" json:"lastShareAuthorId"`
	PayerId					int64				`description:"付款人id" json:"payerId"`
	PayeeId					int64				`description:"收款人id" json:"payeeId"`
	PayType					int					`description:"付款方式，1微信，2支付宝" json:"payType" xorm:"notnull default 0"`
	PayTime          		int64  				`description:"支付时间" json:"payTime"`
	PayTransactionId 		string 				`description:"第三方支付交易号" json:"payTransactionId"`
	Amount					int					`description:"付款金额，单位分" json:"amount" xorm:"notnull default 0"`
	Status					int					`description:"订单状态，0为待付款，1为已付款，2为付款失败" json:"status" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//--------------------------常量---------------------------------------------
const (
	ORDER_AND_PAY_STATUS_WAITING_PAY = iota
	ORDER_AND_PAY_STATUS_PAY_SUCCESS
	ORDER_AND_PAY_STATUS_PAY_FAILED
)


//---------------------------结构体-----------------------------------------
type WeChatRequestParams map[string]string


type ShareLevelAndId map[string]int64

//微信支付h5参数
type OrderParams struct {
	UId          int64  `description:"uId" json:"uId" structs:"uId"`
	TotalFee     int    `description:"支付总金额 单位分" json:"totalFee" valid:"Min(0)" structs:"totalFee"`
	TotalMinute  int    `description:"支付总时间 单位分钟" json:"totalMinute" structs:"totalMinute"`
	ExchangeRate int    `description:"豆币对人民币汇率" json:"exchangeRate" valid:"Min(0)" structs:"exchangeRate"`
	PayType      int    `description:"支付方式 1微信 2支付宝  4苹果" json:"payType" valid:"Range(1, 2)" structs:"payType"`
	Noncestr     string `description:"随机字符串" json:"noncestr" valid:"Required" structs:"noncestr"`
	Timestamp    string `description:"时间戳" json:"timestamp" valid:"Required" structs:"timestamp"`
	Sign         string `description:"签名" json:"sign" valid:"Required" structs:"sign"`
	OrderSource  string `description:"充值来源" json:"orderSource" valid:"Required" structs:"orderSource"`
	Location         string  `description:"地址" json:"location"`
	LocationCode     string  `description:"地址码" json:"locationCode"`
	AppNo	     int  	`description:"app号，1是得问，2是掌律" json:"appNo"`
	IsFxPackage	   	 int  	`description:"是否分销商购买礼包" json:"isFxPackage" structs:"isFxPackage"`
	IsH5Pay	   	 int  	`description:"是否h5支付" json:"isH5Pay" structs:"isH5Pay"`
}

type WeChatPayUnifiedOrderRequestBody struct {
	XMLName        xml.Name `xml:"xml"`
	AppId          string   `xml:"appid"`
	MchId          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIp string   `xml:"spbill_create_ip"`
	NotifyUrl      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
	Openid	       string   `xml:"openid"`
}

type WeChatPayUnifiedOrderResponseBody struct {
	XMLName    xml.Name `xml:"xml" structs:"-"`
	ReturnCode string   `xml:"return_code" structs:"return_code"`
	ReturnMsg  string   `xml:"return_msg" structs:"return_msg"`
	AppId      string   `xml:"appid" structs:"appid"`
	MchId      string   `xml:"mch_id" structs:"mch_id"`
	DeviceInfo string   `xml:"device_info" structs:"device_info"`
	NonceStr   string   `xml:"nonce_str" structs:"nonce_str"`
	Sign       string   `xml:"sign" structs:"sign"`
	ResultCode string   `xml:"result_code" structs:"result_code"`
	ErrCode    string   `xml:"err_code" structs:"err_code"`
	ErrCodeDes string   `xml:"err_code_des" structs:"err_code_des"`
	TradeType  string   `xml:"trade_type" structs:"trade_type"`
	PrepayId   string   `xml:"prepay_id" structs:"prepay_id"`
	MwebUrl   string   `xml:"mweb_url" structs:"mweb_url"`
}

//type CreateOrderResponse struct {
//	WechatPartnerid string `description:"微信商户号" json:"wechatPartnerid"`
//	WechatPrepayid  string `description:"微信预支付交易会话ID" json:"wechatPrepayid"`
//	WechatPackage   string `description:"微信扩展字段" json:"wechatPackage"`
//	WechatNoncestr  string `description:"微信随机字符串" json:"wechatNoncestr"`
//	WechatTimestamp string `description:"微信时间戳" json:"wechatTimestamp"`
//	WechatSign      string `description:"微信签名" json:"wechatSign"`
//	WechatMwebUrl   string `description:"微信h5支付url" json:"wechatMwebUrl"`
//}

type CreateOrderResponse struct {
	WechatAppId 	string `description:"应用appId" json:"appId"`
	WechatTimeStamp string `description:"微信时间戳" json:"timeStamp"`
	WechatNoncestr  string `description:"微信随机字符串" json:"nonceStr"`
	WechatPackage   string `description:"微信扩展字段" json:"package"`
	WechatSignType  string `description:"微信签名" json:"signType"`
	WechatPaySign   string `description:"微信签名" json:"paySign"`

	WechatPartnerId string `description:"商户号" json:"partnerId"`
	WechatPrepayId  string `description:"预支付交易会话ID" json:"prepayId"`
	WechatSign      string `description:"微信签名" json:"sign"`
}

type WeChatPayCallBackBody struct {
	XMLName       xml.Name `xml:"xml" structs:"-"`
	ReturnCode    string   `xml:"return_code" structs:"return_code"`
	ReturnMsg     string   `xml:"return_msg" structs:"return_msg"`
	AppId         string   `xml:"appid" structs:"appid"`
	MchId         string   `xml:"mch_id" structs:"mch_id"`
	DeviceInfo    string   `xml:"device_info" structs:"device_info"`
	NonceStr      string   `xml:"nonce_str" structs:"nonce_str"`
	Sign          string   `xml:"sign" structs:"sign"`
	ResultCode    string   `xml:"result_code" structs:"result_code"`
	ErrCode       string   `xml:"err_code" structs:"err_code"`
	ErrCodeDes    string   `xml:"err_code_des" structs:"err_code_des"`
	OpenId        string   `xml:"openid" structs:"openid"`
	IsSubscribe   string   `xml:"is_subscribe" structs:"is_subscribe"`
	TradeType     string   `xml:"trade_type" structs:"trade_type"`
	BankType      string   `xml:"bank_type" structs:"bank_type"`
	TotalFee      int      `xml:"total_fee" structs:"total_fee"`
	FeeType       string   `xml:"fee_type" structs:"fee_type"`
	CashFee       int      `xml:"cash_fee" structs:"cash_fee"`
	CashFeeType   string   `xml:"cash_fee_type" structs:"cash_fee_type"`
	TransactionId string   `xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string   `xml:"out_trade_no" structs:"out_trade_no"`
	Attach        string   `xml:"attach" structs:"attach"`
	TimeEnd       string   `xml:"time_end" structs:"time_end"`
}

type WeChatPayCallBackResponseBody struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}
//收入查询返回结构体
type Revenue struct {
	PayOrder                `description:"订单信息" xorm:"extends"`
	Sale          string    `description:"销售类型" json:"sale"`
	Money      	  int    	`description:"交易金额，单位分" json:"money"`
	PayName		  string	`description:"支付人昵称" json:"payName"`
	ShareName	  string	`description:"分享人昵称" json:"shareName"`
	Title    	  string	`description:"视频标题" json:"title"`
}

type RevenueListContainer struct {
	BaseListContainer
	RevenueList    []Revenue	`description:"收入列表" json:"revenueList"`
	RevenueMonth   int          `description:"本月收入" json:"revenueMonth"`
	AccountAmount  int          `description:"账户余额" json:"accountAmount"`
}





//---------------------------签名方法--------------------------------------------

func (this WeChatRequestParams)WechatPaySign() string {
	keys := make([]string, len(this))

	i := 0
	for k := range this {
		keys[i] = k
		i++
	}

	//	fmt.Println("keys", keys)
	sort.Strings(keys)
	//	fmt.Println("sorted keys", keys)

	strTemp := ""
	for _, key := range keys {
		if key == "scene_info" {
			continue
		}
		strTemp = strTemp + key + "=" + this[key] + "&"
	}
	strTemp += "key=" + JTWechatPayApiKey
		//fmt.Println("strTemp = ", strTemp)

	hasher := md5.New()
	hasher.Write([]byte(strTemp))
	md5Str := hex.EncodeToString(hasher.Sum(nil))

		//fmt.Println("md5 = ", md5Str)

	return strings.ToUpper(md5Str)
}

func ConvertXMLMapToSignParams(value map[string]interface{}) WeChatRequestParams {
	params := make(WeChatRequestParams)
	for k := range value {
		if k == "sign" {
			continue
		}
		switch element := value[k].(type) {
		case string:
			if len(element) > 0 {
				params[k] = element
			}
		case int:
			params[k] = strconv.Itoa(element)
		default:

		}
	}
	return params
}

