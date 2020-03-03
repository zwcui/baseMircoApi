/*
@Time : 2019/2/25 上午11:01 
@Author : zwcui
@Software: GoLand
*/
package models

import "encoding/xml"

//加密后每隔10分钟定时收到component_verify_ticket
type ComponentVerifyTicketEncryptXmlBody struct {
	XMLName       			xml.Name 		`xml:"xml" structs:"-"`
	AppId    				string   		`xml:"AppId" structs:"AppId"`
	Encrypt    				string   		`xml:"Encrypt" structs:"Encrypt"`
}

//每隔10分钟定时收到component_verify_ticket，用于更新component_access_token
type ComponentVerifyTicketXmlBody struct {
	XMLName       			xml.Name 		`xml:"xml" structs:"-"`
	AppId    				string   		`xml:"AppId" structs:"AppId"`
	CreateTime     			int64   		`xml:"CreateTime" structs:"CreateTime"`
	InfoType        		string   		`xml:"InfoType" structs:"InfoType"`
	ComponentVerifyTicket   string   		`xml:"ComponentVerifyTicket" structs:"ComponentVerifyTicket"`
}

//component_access_token请求返回体
type ComponentAccessTokenJsonBody struct {
	ComponentAccessToken	string 			`description:"component_access_token" json:"component_access_token"`
	ExpiresIn	  			int64 			`description:"expires_in" json:"expires_in"`
}

//pre_auth_code请求返回体
type PreAuthCodeJsonBody struct {
	PreAuthCode				string 			`description:"pre_auth_code" json:"pre_auth_code"`
	ExpiresIn	  			int64 			`description:"expires_in" json:"expires_in"`
}


//使用授权码换取公众号或小程序的接口调用凭据和授权信息返回体
type AuthorizationInfoResponseJsonBody struct {
	AuthorizationInfoJsonBody	AuthorizationInfoResponse		`description:"授权方" json:"authorization_info"`
}

type AuthorizationInfoResponse struct {
	AuthorizerAppid				string 			`description:"授权方appid" json:"authorizer_appid"`
	AuthorizerAccessToken		string 			`description:"授权方接口调用凭据（在授权的公众号或小程序具备API权限时，才有此返回值），也简称为令牌" json:"authorizer_access_token"`
	ExpiresIn					int64 			`description:"有效期" json:"expires_in"`
	AuthorizerRefreshToken		string 			`description:"接口调用凭据刷新令牌（在授权的公众号具备API权限时，才有此返回值），刷新令牌主要用于第三方平台获取和刷新已授权用户的access_token，只会在授权时刻提供，请妥善保存。 一旦丢失，只能让用户重新授权，才能再次拿到新的刷新令牌" json:"authorizer_refresh_token"`
	FuncInfo					[]FuncscopeCategory 			`description:"" json:"func_info"`
}

type FuncscopeCategory struct {
	FuncscopeCategory			FuncscopeCategoryDetail 		`description:"" json:"funcscope_category"`
}

type FuncscopeCategoryDetail struct {
	Id							int 			`description:"id" json:"id"`
}


//获取授权方的基本信息，包括头像、昵称、帐号类型、认证类型、微信号、原始ID和二维码图片URL
type AuthorizerInfoResponseJsonBody struct {
	AuthorizerInfoJsonBody		AuthorizerInfoResponse			`description:"授权方详情" json:"authorizer_info"`
	AuthorizationInfoJsonBody	AuthorizationInfoResponseShort	`description:"授权方详情" json:"authorization_info"`
}

type AuthorizerInfoResponse struct {
	NickName					string 			`description:"授权方昵称" json:"nick_name"`
	HeadImg						string 			`description:"授权方头像" json:"head_img"`
	ServiceTypeInfo				ServiceTypeInfo	`description:"授权方公众号类型，0代表订阅号，1代表由历史老帐号升级后的订阅号，2代表服务号" json:"service_type_info"`
	VerifyTypeInfo				VerifyTypeInfo  `description:"授权方认证类型，-1代表未认证，0代表微信认证，1代表新浪微博认证，2代表腾讯微博认证，3代表已资质认证通过但还未通过名称认证，4代表已资质认证通过、还未通过名称认证，但通过了新浪微博认证，5代表已资质认证通过、还未通过名称认证，但通过了腾讯微博认证" json:"verify_type_info"`
	UserName					string 			`description:"授权方公众号的原始ID" json:"user_name"`
	PrincipalName				string 			`description:"公众号的主体名称" json:"principal_name"`
	BusinessInfo				BusinessInfo 	`description:"用以了解以下功能的开通状况（0代表未开通，1代表已开通）： open_store:是否开通微信门店功能 open_scan:是否开通微信扫商品功能 open_pay:是否开通微信支付功能 open_card:是否开通微信卡券功能 open_shake:是否开通微信摇一摇功能" json:"business_info"`
	Alias						string 			`description:"授权方公众号所设置的微信号，可能为空" json:"alias"`
	QrcodeUrl					string 			`description:"二维码图片的URL，开发者最好自行也进行保存" json:"qrcode_url"`
}

type AuthorizationInfoResponseShort struct {
	AuthorizerAppid				string 			`description:"授权方appid" json:"authorizer_appid"`
	FuncInfo					[]FuncscopeCategory 			`description:"" json:"func_info"`
}

type ServiceTypeInfo struct {
	Id							int 			`description:"id" json:"id"`
}

type VerifyTypeInfo struct {
	Id							int 			`description:"id" json:"id"`
}

type BusinessInfo struct {
	OpenStore					int 			`description:"是否开通微信门店功能" json:"open_store"`
	OpenScan					int 			`description:"是否开通微信扫商品功能" json:"open_scan"`
	OpenPay						int 			`description:"是否开通微信支付功能" json:"open_pay"`
	OpenCard					int 			`description:"是否开通微信卡券功能" json:"open_card"`
	OpenShake					int 			`description:"是否开通微信摇一摇功能" json:"open_shake"`
}


//获取（刷新）授权公众号或小程序的接口调用凭据（令牌）
type RefreshAuthorizerAccessTokenJsonBody struct {
	AuthorizerAccessToken		string 			`description:"授权方令牌" json:"authorizer_access_token"`
	ExpiresIn					int64 			`description:"有效期，为2小时" json:"expires_in"`
	AuthorizerRefreshToken		string 			`description:"刷新令牌(一旦丢失，只能让用户重新授权，才能再次拿到新的刷新令牌)" json:"authorizer_refresh_token"`
}


//微信接口异常信息
type ErrorJsonBody struct {
	Errcode						int 			`description:"错误码" json:"errcode"`
	Errmsg						string 			`description:"错误描述" json:"errmsg"`
}