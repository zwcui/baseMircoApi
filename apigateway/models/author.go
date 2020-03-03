/*
@Time : 2019/5/31 下午5:26 
@Author : zwcui
@Software: GoLand
*/
package models

import "encoding/json"

type Author struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	Openid				string			`description:"用户的标识，对当前公众号唯一(关注、授权)" json:"openid"`
	WebOpenid			string			`description:"用户的标识，对当前网站应用唯一" json:"webOpenid"`
	AppOpenid			string			`description:"用户的标识，对当前app唯一" json:"appOpenid"`
	Unionid				string			`description:"只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。" json:"unionid" xorm:"notnull"`
	AccessToken  		string 			`description:"接口调用凭证" json:"access_token"`
	ExpiresIn    		int    			`description:"access_token接口调用凭证超时时间，单位（秒）" json:"expires_in"`
	RefreshToken 		string 			`description:"用户刷新access_token" json:"refresh_token"`
	Scope        		string 			`description:"用户授权的作用域，使用逗号（,）分隔" json:"scope"`
	Nickname			string			`description:"用户的昵称" json:"nickname"`
	Sex					int				`description:"用户的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	City				string			`description:"用户所在城市" json:"city"`
	Province			string			`description:"用户所在省份" json:"province"`
	Country				string			`description:"用户所在国家" json:"country"`
	Headimgurl			string			`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
	Privilege  			string 			`description:"用户特权信息，json数组，如微信沃卡用户为（chinaunicom）" json:"privilege"`
	PhoneNumber  		string 			`description:"用户手机号" json:"phoneNumber"`
	Password            string          `description:"密码,加密存储"`
	Salt                string          `description:"salt"`
	IsTest              int             `description:"是否测试账号 0:否 1:是" xorm:"notnull default 0"`
	LoginUuid			string			`description:"登录uuid" json:"loginUuid"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//------------------------------结构体-----------------------------------
type AuthorShort struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	Openid				string			`description:"用户的标识，对当前公众号唯一" json:"openid"`
	WebOpenid			string			`description:"用户的标识，对当前网站应用唯一" json:"webOpenid"`
	Unionid				string			`description:"只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。" json:"unionid"`
	Nickname			string			`description:"用户的昵称" json:"nickname"`
	Sex					int				`description:"用户的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	City				string			`description:"用户所在城市" json:"city"`
	Province			string			`description:"用户所在省份" json:"province"`
	Country				string			`description:"用户所在国家" json:"country"`
	Headimgurl			string			`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
	PhoneNumber  		string 			`description:"用户手机号" json:"phoneNumber"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type AuthorDetail struct {
	Author				AuthorShort		`description:"up主" json:"author"`
}

/*
	如果开发者拥有多个移动应用、网站应用和公众帐号，
	可通过获取用户基本信息中的unionid来区分用户的唯一性，
	因为只要是同一个微信开放平台帐号下的移动应用、网站应用和公众帐号，用户的unionid是唯一的。
	换句话说，同一用户，对同一个微信开放平台下的不同应用，unionid是相同的
 */

//微信登录获取accessToken结构体
type WechatAccessTokenResponse struct {
	AccessToken  		string 			`description:"接口调用凭证" json:"access_token"`
	ExpiresIn    		int    			`description:"access_token接口调用凭证超时时间，单位（秒）" json:"expires_in"`
	RefreshToken 		string 			`description:"用户刷新access_token" json:"refresh_token"`
	Openid       		string 			`description:"授权用户唯一标识" json:"openid"`
	Scope        		string 			`description:"用户授权的作用域，使用逗号（,）分隔" json:"scope"`
	Unionid      		string 			`description:"当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。" json:"unionid"`
}

//微信登录刷新accessToken结构体
type WechatRefreshTokenResponse struct {
	AccessToken  		string 			`description:"接口调用凭证" json:"access_token"`
	ExpiresIn    		int    			`description:"access_token接口调用凭证超时时间，单位（秒）" json:"expires_in"`
	RefreshToken 		string 			`description:"用户刷新access_token" json:"refresh_token"`
	Openid       		string 			`description:"授权用户唯一标识" json:"openid"`
	Scope        		string 			`description:"用户授权的作用域，使用逗号（,）分隔" json:"scope"`
}

//微信登录通过accessToken获取用户信息结构体
type WechatUserInfoResponse struct {
	Openid     			string   		`description:"普通用户的标识，对当前开发者帐号唯一" json:"openid"`
	Nickname  		 	string   		`description:"普通用户昵称" json:"nickname"`
	Sex        			int      		`description:"普通用户性别，1为男性，2为女性" json:"sex"`
	Province   			string   		`description:"普通用户个人资料填写的省份" json:"province"`
	City       			string   		`description:"普通用户个人资料填写的城市" json:"city"`
	Country    			string   		`description:"国家，如中国为CN" json:"country"`
	Headimgurl 			string   		`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空" json:"headimgurl"`
	Privilege  			[]string 		`description:"用户特权信息，json数组，如微信沃卡用户为（chinaunicom）" json:"privilege"`
	Unionid    			string   		`description:"用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。" json:"unionid"`
}

//根据用户id返回用户基本信息结构体
type UserInfo struct {
	Openid     			string   		`description:"普通用户的标识，对当前公众号唯一" json:"openid"`
	WebOpenid			string			`description:"普通用户的标识，对当前网站应用唯一" json:"webOpenid"`
	Nickname  		 	string   		`description:"普通用户昵称" json:"nickname"`
	Sex        			int      		`description:"普通用户性别，1为男性，2为女性" json:"sex"`
	Province   			string   		`description:"普通用户个人资料填写的省份" json:"province"`
	City       			string   		`description:"普通用户个人资料填写的城市" json:"city"`
	Country    			string   		`description:"国家，如中国为CN" json:"country"`
	Headimgurl 			string   		`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空" json:"headimgurl"`
	Privilege  			[]string 		`description:"用户特权信息，json数组，如微信沃卡用户为（chinaunicom）" json:"privilege"`
	Unionid    			string   		`description:"用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。" json:"unionid"`
	PhoneNumber  		string 			`description:"用户手机号" json:"phoneNumber"`
}

//------------------------方法------------------------------------------
//author转authorshort
func (author *Author) AuthorToAuthorShort() (authorShort *AuthorShort, error error) {
	josnByte, err := json.Marshal(author)
	if err != nil {
		return nil, err
	}

	authorShortD := AuthorShort{}
	err = json.Unmarshal(josnByte, &authorShortD)
	if err != nil {
		return nil, err
	}
	return &authorShortD, nil
}
