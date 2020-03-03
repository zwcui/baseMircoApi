/*
@Time : 2019/2/25 下午1:51 
@Author : zwcui
@Software: GoLand
*/
package models

import "encoding/json"

//公众号授权的权限
type AuthInfo struct {
	Id		        	  		int64  			`description:"id" json:"id" xorm:"pk autoincr"`
	AuthCode					string 			`description:"授权方auth_code" json:"authCode"`
	AuthCodeExpireTime			int64 			`description:"授权方auth_code过期时间" json:"authCodeExpireTime"`
	AuthAppid					string 			`description:"授权方appid" json:"authAppid"`
	AuthAccessToken				string 			`description:"授权方接口调用凭据（在授权的公众号或小程序具备API权限时，才有此返回值），也简称为令牌" json:"authAccessToken"`
	AuthAccessTokenExpireTime	int64 			`description:"有效期" json:"expireTime"`
	AuthRefreshToken			string 			`description:"接口调用凭据刷新令牌（在授权的公众号具备API权限时，才有此返回值），刷新令牌主要用于第三方平台获取和刷新已授权用户的access_token，只会在授权时刻提供，请妥善保存。 一旦丢失，只能让用户重新授权，才能再次拿到新的刷新令牌" json:"authRefreshToken"`
	AuthFunctionInfo			string 			`description:"授权给开发者的权限集列表，ID为1到26分别代表： 1、消息管理权限 2、用户管理权限 3、帐号服务权限 4、网页服务权限 5、微信小店权限 6、微信多客服权限 7、群发与通知权限 8、微信卡券权限 9、微信扫一扫权限 10、微信连WIFI权限 11、素材管理权限 12、微信摇周边权限 13、微信门店权限 15、自定义菜单权限 16、获取认证状态及信息 17、帐号管理权限（小程序） 18、开发管理与数据分析权限（小程序） 19、客服消息管理权限（小程序） 20、微信登录权限（小程序） 21、数据分析权限（小程序） 22、城市服务接口权限 23、广告管理权限 24、开放平台帐号管理权限 25、 开放平台帐号管理权限（小程序） 26、微信电子发票权限 41、搜索widget的权限 请注意： 1）该字段的返回不会考虑公众号是否具备该权限集的权限（因为可能部分具备），请根据公众号的帐号类型和认证情况，来判断公众号的接口权限。" json:"authFunctionInfo"`
	NickName					string 			`description:"授权方昵称" json:"nickName"`
	HeadImg						string 			`description:"授权方头像" json:"headImg"`
	ServiceType					int				`description:"授权方公众号类型，0代表订阅号，1代表由历史老帐号升级后的订阅号，2代表服务号" json:"serviceType"`
	VerifyType					int				`description:"授权方认证类型，-1代表未认证，0代表微信认证，1代表新浪微博认证，2代表腾讯微博认证，3代表已资质认证通过但还未通过名称认证，4代表已资质认证通过、还未通过名称认证，但通过了新浪微博认证，5代表已资质认证通过、还未通过名称认证，但通过了腾讯微博认证" json:"verifyType"`
	UserName					string 			`description:"授权方公众号的原始ID" json:"userName"`
	PrincipalName				string 			`description:"公众号的主体名称" json:"principalName"`
	Alias						string 			`description:"授权方公众号所设置的微信号，可能为空" json:"alias"`
	QrcodeUrl					string 			`description:"二维码图片的URL，开发者最好自行也进行保存" json:"qrcodeUrl"`
	OpenStore					int 			`description:"是否开通微信门店功能（0代表未开通，1代表已开通）" json:"openStore"`
	OpenScan					int 			`description:"是否开通微信扫商品功能（0代表未开通，1代表已开通）" json:"openScan"`
	OpenPay						int 			`description:"是否开通微信支付功能（0代表未开通，1代表已开通）" json:"openPay"`
	OpenCard					int 			`description:"是否开通微信卡券功能（0代表未开通，1代表已开通）" json:"openCard"`
	OpenShake					int 			`description:"是否开通微信摇一摇功能（0代表未开通，1代表已开通）" json:"openShake"`
	Status	 					int				`description:"状态，1为启用 2停用" json:"status" xorm:"notnull default 0"`
	Remark	 					string			`description:"备注" json:"remark"`
	Created           			int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           			int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         			int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//------------------------------结构体------------------------------------

type AuthInfoShort struct {
	Id		        	  		int64  			`description:"id" json:"id" xorm:"pk autoincr"`
	NickName					string 			`description:"授权方昵称" json:"nickName"`
	HeadImg						string 			`description:"授权方头像" json:"headImg"`
	ServiceType					int				`description:"授权方公众号类型，0代表订阅号，1代表由历史老帐号升级后的订阅号，2代表服务号" json:"serviceType"`
	VerifyType					int				`description:"授权方认证类型，-1代表未认证，0代表微信认证，1代表新浪微博认证，2代表腾讯微博认证，3代表已资质认证通过但还未通过名称认证，4代表已资质认证通过、还未通过名称认证，但通过了新浪微博认证，5代表已资质认证通过、还未通过名称认证，但通过了腾讯微博认证" json:"verifyType"`
	PrincipalName				string 			`description:"公众号的主体名称" json:"principalName"`
	QrcodeUrl					string 			`description:"二维码图片的URL，开发者最好自行也进行保存" json:"qrcodeUrl"`
	Status	 					int				`description:"状态，1为启用 2停用" json:"status" xorm:"notnull default 0"`
	Remark	 					string			`description:"备注" json:"remark"`
	Created           			int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           			int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         			int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type AuthInfoShortContainer struct {
	AuthInfoShort				AuthInfoShort 	`description:"用户权限下公众号信息" json:"authInfo"`
}

type UserAuthListContainer struct {
	UserAuthList				[]AuthInfoShort `description:"用户权限下公众号信息" json:"userAuthList"`
}

//-------------------------------方法----------------------------------

func (a *AuthInfo) AuthInfoToAuthInfoShort() (authInfoDTO *AuthInfoShort, error error) {
	josnByte, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	authInfo := AuthInfoShort{}
	err = json.Unmarshal(josnByte, &authInfo)
	if err != nil {
		return nil, err
	}
	return &authInfo, nil
}