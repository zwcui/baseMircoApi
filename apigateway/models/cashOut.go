/*
@Time : 2019/6/13 下午2:59 
@Author : lianwu
@File : transfersOrder.go
@Software: GoLand
*/
package models

import (
	"encoding/xml"
	"time"
	"os"
	"github.com/go-xorm/xorm"
	"strconv"
	"net/url"
	"io/ioutil"
	"errors"
	"math/rand"
	"net/http"
	"crypto/tls"
	"net"
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
)


//提现记录
type CashOut struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	TransferType			int					`description:"转账类型，1为用户提现 2为平台返现" json:"transferType"`
	SettlementType			int					`description:"结算类型，1为每单结算，2为每日结算，3为每周结算，4为每月结算" json:"settlementType"`
	AuthorId				int64				`description:"收款人id" json:"authorId"`
	TransferTime          	string  			`description:"转账时间" json:"transferTime"`
	PaymentNo 				string 				`description:"第三方支付交易号" json:"paymentNo"`
	Amount					int					`description:"付款金额，单位分" json:"amount" xorm:"notnull default 0"`
	Status					int					`description:"转账状态，0为待处理，1为转账成功，2为转账失败" json:"status" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
	ErrorCode               string				`description:"错误信息" json:"errorCode"`
	ErrorCodeDes            string				`description:"错误信息描述" json:"errorCodeDes"`
}

//转账记录
type TransfersHistory struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	CashOutId       		int64			    `description:"提现记录id" json:"cashOutId"`
	MchAppid                string				`description:"申请商户号的appid" json:"mchAppid"`
	MchId                   string				`description:"商户号" json:"mchId"`
	NonceStr                string				`description:"随机字符串" json:"nonceStr"`
	Sign                    string				`description:"签名" json:"sign"`
	PartnerTradeNo          string				`description:"商户订单号" json:"partnerTradeNo"`
	Openid                  string				`description:"收款人openid" json:"openid"`
	CheckName				string				`description:"是否校验用户姓名" json:"checkName"`
	ReUserName				string				`description:"收款用户真实姓名" json:"reUserName"`
	Amount					int					`description:"付款金额，单位分" json:"amount" xorm:"notnull default 0"`
	Desc					string				`description:"企业付款备注" json:"desc"`
	SpbillCreateIp       	string				`description:"ip地址" json:"spbillCreateIp"`
	TransferType			int					`description:"转账类型，1为用户提现 2为平台返现" json:"transferType"`
	AuthorId				int64				`description:"收款人id" json:"authorId"`
	ReturnPartnerTradeNo    string				`description:"接口返回的商户订单号" json:"returnPartnerTradeNo"`
	TransferTime          	string  			`description:"转账时间" json:"transferTime"`
	PaymentNo 				string 				`description:"第三方支付交易号" json:"paymentNo"`
	Status					int					`description:"转账状态，0为待处理，1为转账成功，2为转账失败" json:"status" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
	ReturnCode              string				`description:"返回状态码(通信标识)" json:"returnCode"`
	ReturnMsg               string				`description:"返回信息" json:"returnMsg"`
	ResultCode              string				`description:"业务结果" json:"resultCode"`
	ErrorCode               string				`description:"错误信息" json:"errorCode"`
	ErrorCodeDes            string				`description:"错误信息描述" json:"errorCodeDes"`
}




//---------------------------结构体-----------------------------------------

type WeChatTransfersRequestBody struct {
	XMLName        xml.Name `xml:"xml"`
	MchAppid       string   `xml:"mch_appid"`
	MchId          string   `xml:"mchid"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	PartnerTradeNo string   `xml:"partner_trade_no"`
	Openid         string   `xml:"openid"`
	CheckName      string   `xml:"check_name"`
	ReUserName	   string   `xml:"re_user_name"`
	Amount         int      `xml:"amount"`
	Desc           string   `xml:"desc"`
	SpbillCreateIp string   `xml:"spbill_create_ip"`


}


type WeChatTransfersResponseBody struct {
	XMLName    		xml.Name `xml:"xml" structs:"-"`
	ReturnCode      string   `xml:"return_code" structs:"return_code"`
	ReturnMsg       string   `xml:"return_msg" structs:"return_msg"`
	AppId           string   `xml:"appid" structs:"appid"`
	MchId           string   `xml:"mchid" structs:"mchid"`
	DeviceInfo      string   `xml:"device_info" structs:"device_info"`
	NonceStr        string   `xml:"nonce_str" structs:"nonce_str"`
	ResultCode      string   `xml:"result_code" structs:"result_code"`
	ErrCode         string   `xml:"err_code" structs:"err_code"`
	ErrCodeDes      string   `xml:"err_code_des" structs:"err_code_des"`
	PartnerTradeNo  string   `xml:"partner_trade_no" structs:"partner_trade_no"`
	PaymentNo   	string   `xml:"payment_no" structs:"payment_no"`
	PaymentTime     string   `xml:"payment_time" structs:"payment_time"`
}

type CreateTransfersHistoryResponse struct {
	PartnerTradeNo 	string `description:"商户订单号" json:"partnerTradeNo"`
	PaymentNo 		string `description:"微信付款单号" json:"PaymentNo"`
	PaymentTime  	string `description:"付款成功时间" json:"PaymentTime"`
}

type CashOutContainer struct {
	BaseListContainer
	CashOutList	[]CashOut	`description:"提现记录列表" json:"cashOutList"`
}

//-------------------------------方法-------------------------------------------------------------------

//企业向微信用户个人付款至微信零钱
func WeChatTransfers(spbillCreateIp string, cashOut *CashOut, session *xorm.Session) (responsebody *WeChatTransfersResponseBody, err error) {

	var author Author
	var openid string
	has, err := session.Table("author").Where("id = ?", cashOut.AuthorId).Get(&author)
	if err != nil {
		fmt.Println("weChatTransfers get author err = ", err.Error())
		return nil, err
	}
	if has {
		openid = author.Openid
	} else {
		err = errors.New("该用户不存在，请核对信息")
		return nil, err
	}
	requestBody := WeChatTransfersRequestBody{}
	requestBody.SpbillCreateIp = spbillCreateIp
	requestBody.Openid = openid
	requestBody.Amount = cashOut.Amount
	requestBody.PartnerTradeNo = "CashOut" + beego.BConfig.RunMode +strconv.FormatInt(cashOut.Id,10)
	requestBody.CheckName = "NO_CHECK"
	if cashOut.TransferType == 1 {
		requestBody.Desc = "用户提现转账"
	} else {
		requestBody.Desc = "平台返现转账"
	}
	requestBody.MchAppid = JTGZHAppId
	requestBody.MchId = JTWechatPayBusinessNum
	requestBody.NonceStr = strconv.FormatUint(uint64(rand.Uint32()), 10)


	params := WeChatRequestParams{}
	params["mch_appid"] = requestBody.MchAppid
	params["mchid"] = requestBody.MchId
	params["openid"] = requestBody.Openid
	params["nonce_str"] = requestBody.NonceStr
	params["partner_trade_no"] = requestBody.PartnerTradeNo
	params["spbill_create_ip"] = requestBody.SpbillCreateIp
	params["desc"] = requestBody.Desc
	params["amount"] = strconv.Itoa(requestBody.Amount)
	params["check_name"] = requestBody.CheckName
	//	log.Println("wechatpay params ", params)
	fmt.Println("wechatpay params ", params)
	sign := params.WechatPaySign()
	requestBody.Sign = sign


	//调微信接口
	var Url *url.URL
	Url, err = url.Parse("https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers")
	if err != nil {
		fmt.Println("wechatTransfer prseurl err = ", err.Error())
		return nil, err
	}

	//	fmt.Println(Url.String())

	xmlBodyBytes, err := xml.MarshalIndent(requestBody, "  ", "    ")
	if err != nil {
		fmt.Println("wechatTransfer request err = ", err.Error())
		return nil, err
	}

	fmt.Println("xmlStr = ", string(xmlBodyBytes))

	certPath := getAPPRootPath() + JTCertPath
	keyPath := getAPPRootPath() + JTKeyPath
	cli, err := newTLSClient(certPath, keyPath)
	if err != nil {
		fmt.Println("wechatTransfer 证书 err = ", err.Error())
		return nil, err
	}

	res, err := cli.Post(Url.String(), "application/xml; charset=utf-8", bytes.NewBuffer(xmlBodyBytes))
	if err != nil {
		fmt.Println("wechatTransfer request err = ", err.Error())
		return nil, err
	}
	defer res.Body.Close()


	//	fmt.Println("response Status:", resp.Status)
	//	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("wechatTransfer readbody err = ", err.Error())
		return nil, err
	}
	//	fmt.Println("WeChatPay Body:", string(body))
	response := WeChatTransfersResponseBody{}
	err = xml.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("wechatTransfer Unmarshal response err = ", err.Error())
		return nil, err
	}

	//创建转账记录
	var transfersHistory TransfersHistory
	transfersHistory.CashOutId = cashOut.Id
	transfersHistory.MchAppid = requestBody.MchAppid
	transfersHistory.MchId = requestBody.MchId
	transfersHistory.NonceStr = requestBody.NonceStr
	transfersHistory.PartnerTradeNo = requestBody.PartnerTradeNo
	transfersHistory.Openid = requestBody.Openid
	transfersHistory.CheckName = requestBody.CheckName
	transfersHistory.ReUserName = requestBody.ReUserName
	transfersHistory.Amount = requestBody.Amount
	transfersHistory.Desc = requestBody.Desc
	transfersHistory.SpbillCreateIp = requestBody.SpbillCreateIp
	transfersHistory.Sign = requestBody.Sign
	transfersHistory.TransferType = cashOut.TransferType
	transfersHistory.AuthorId = cashOut.AuthorId

	if response.ResultCode == "SUCCESS" && response.ReturnCode == "SUCCESS" {
		transfersHistory.Status = 1
	} else {
		transfersHistory.Status = 2
	}
	transfersHistory.PaymentNo = response.PaymentNo
	transfersHistory.TransferTime = response.PaymentTime
	transfersHistory.ReturnCode = response.ReturnCode
	transfersHistory.ReturnMsg = response.ReturnMsg
	transfersHistory.ResultCode = response.ResultCode
	transfersHistory.ErrorCode = response.ErrCode
	transfersHistory.ErrorCodeDes = response.ErrCodeDes
	_, err = session.Table("transfers_history").InsertOne(&transfersHistory)
	if err != nil {
		fmt.Println("wechatTransfer err = ", err.Error())
		return nil,err
	}

	//	fmt.Println("reponse struct :", response)

	if response.ReturnCode == "FAIL" {
		err = errors.New("返回内容结果异常失败")
		fmt.Println("response.ReturnCode == FAIL  err:" + response.ReturnMsg)
		return nil, err
	}


	if response.ReturnCode != "SUCCESS" {
		err = errors.New(response.ReturnMsg)
		fmt.Println("wechatTransfer err = ", err.Error())
		return nil, err
	}

	if response.ResultCode != "SUCCESS" {
		err = errors.New("请求微信转账接口失败")
		if len(response.ErrCodeDes) > 0 {
			err = errors.New(response.ErrCode+":"+response.ErrCodeDes)
		}
		fmt.Println("wechatTransfer err = ", err.Error())
		return nil, err
	}


	return &response, nil
}
//转账成功
func TransfersSucceed(cashOut *CashOut, authorAccount *AuthorAccount, session *xorm.Session)(err error){

	//根据微信返回结果回填
	//1.修改提现记录状态，回填第三方交易单号以及转账成功时间
	cashOut.Status = 1
	var recordType int
	if cashOut.TransferType == 1 {
		recordType = 2
	} else {
		recordType = 3
	}
	_, err = session.Table("cash_out").Where("id = ?", cashOut.Id).AllCols().Update(cashOut)
	if err != nil {
		fmt.Println("transfersSucceed update(cash_out)  err = ", err.Error())
		return err
	}
	err = UpdateAccount(authorAccount, cashOut.Amount, 1, recordType, cashOut.Id, 1, nil, session)
	if err != nil {
		fmt.Println("transfersSucceed UpdateAccount err = ", err.Error())
		return err
	}
	return

}




//转账失败
func TransfersFail(cashOut *CashOut, authorAccount *AuthorAccount, session *xorm.Session)(err error){

	//根据微信返回结果回填
	//1.修改提现记录状态，回填第三方交易单号以及转账成功时间
	cashOut.Status = 2
	var recordType int
	if cashOut.TransferType == 1 {
		recordType = 2
	} else {
		recordType = 3
	}
	_, err = session.Table("cash_out").Where("id = ?", cashOut.Id).AllCols().Update(cashOut)
	if err != nil {
		fmt.Println("transfersFail Update(cashOut) err = ", err.Error())
		return err
	}
	err = UpdateAccount(authorAccount, cashOut.Amount, 1, recordType, cashOut.Id, 0, nil, session)
	if err != nil {
		fmt.Println("transfersFail UpdateAccount err = ", err.Error())
		return err
	}

	return

}


//amount 大于0
func UpdateAccount (authorAccount *AuthorAccount, amount int, moneyType int, recordType int, orderId int64, status int, param map[string]string, session *xorm.Session) (err error){
	if authorAccount == nil {
		fmt.Println("account not exist")
		return errors.New("account not exist")
	}
	if amount == 0 {
		fmt.Println("amount = 0 ")
		return errors.New("amount = 0 ")
	}
	if moneyType == 0 {
		authorAccount.Amount += amount
	} else {
		authorAccount.Amount -= amount
	}
	if authorAccount.Amount < 0 {
		authorAccount.Amount = 0
	}
	_, err = session.Table("author_account").Where("author_id = ?", authorAccount.AuthorId).Cols("amount").Update(authorAccount)
	if err != nil {
		fmt.Println("update author_account err = ", err.Error())
		return errors.New("update author_account err = " + err.Error())
	}
	var record AccountTransactionRecord
	var recordName string
	record.AuthorId = authorAccount.AuthorId
	record.Money = amount
	record.MoneyType = moneyType
	record.OrderId = orderId
	record.Status = status
	if recordType == 0 {
		recordName = "充值"
	} else if recordType == 1 {
		recordName = "购买视频"
	} else if recordType == 2 {
		recordName = "用户提现"
	} else if recordType == 3 {
		recordName = "平台返现"
	} else if recordType == 4 {
		recordName = "平台分佣"
	}
	if moneyType == 0 {
		recordName += "收入"
	} else {
		recordName += "支出"
	}
	record.RecordType = recordType
	record.RecordName = recordName

	//平台分佣时记录基础金额、层级和比例
	if recordType == 4 && param != nil {
		if shareVideoId, ok := param["shareVideoId"]; ok {
			record.ShareVideoId, _ = strconv.ParseInt(shareVideoId, 10, 64)
		}
		if shareAmount, ok := param["shareAmount"]; ok {
			record.ShareAmount, _ = strconv.Atoi(shareAmount)
		}
		if shareLevel, ok := param["shareLevel"]; ok {
			record.ShareLevel, _ = strconv.Atoi(shareLevel)
		}
		if shareLevelPercent, ok := param["shareLevelPercent"]; ok {
			record.ShareLevelPercent, _ = strconv.Atoi(shareLevelPercent)
		}
	}

	_, err = session.Table("account_transaction_record").InsertOne(&record)
	if err != nil {
		fmt.Println("update author_account err = ", err.Error())
		return errors.New("insert author_account err = " + err.Error())
	}
	return
}






//创建支持双向证书认证的 http.Client.
func newTLSClient(certPath, keyPath string) (httpClient *http.Client, err error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return tlsClient(tlsConfig)
}

func tlsClient(tlsConfig *tls.Config) (*http.Client, error) {

	dialTLS := func(network, addr string) (net.Conn, error) {
		return tls.DialWithDialer(&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}, network, addr, tlsConfig)
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			DialTLS:               dialTLS,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}, nil
}

func getAPPRootPath()string{
	currentPath, _ := os.Getwd()
	fmt.Println("current path:", currentPath)
	return currentPath
}