package models

import "encoding/json"

//用户信息
type User struct {
	UId       			int64			`description:"uId" json:"uId" xorm:"pk autoincr"`
	PhoneNumber			string			`description:"手机号" json:"phoneNumber"`
	NickName 			string			`description:"昵称（登录账号）" json:"nickName" xorm:"notnull "`
	Email	 			string			`description:"电子邮箱" json:"email"`
	Password 			string			`description:"密码" json:"password" xorm:"notnull"`
	Salt	 			string			`description:"密码" json:"salt" xorm:"notnull"`
	RoleType 			string			`description:"角色，1为用户管理员（可创建用户、赋权限），2为帮助管理员（可创建帮助），3为内容发布者（可发布内容、新增组件），4为视频脉脉视频审核员（可审核up主视频信息）【roleType格式为[角色]，如[1][2]】" json:"roleType"`
	DefaultAuthInfoId	int64			`description:"默认公众号id" json:"defaultAuthInfoId"`
	DefaultAuthAppid	string 			`description:"默认授权方appid" json:"defaultAuthAppid"`
	Created           	int64  			`description:"注册时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//用户权限下的公众号
type UserAuth struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	UId       			int64			`description:"uId" json:"uId"`
	AuthInfoId 			int64			`description:"公众号id" json:"authInfoId"`
	Status	 			int				`description:"状态，1为正常" json:"status" xorm:"notnull default 0"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//公众号下的订阅者
type Subscriber struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	Openid				string			`description:"用户的标识，对当前公众号唯一" json:"openid"`
	AuthInfoId 			int64			`description:"公众号id" json:"authInfoId"`
	GzhAppId 			string			`description:"公众号AppId" json:"gzhAppId"`
	Subscribe       	int				`description:"用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。" json:"subscribe"`
	Nickname			string			`description:"用户的昵称" json:"nickname"`
	Sex					int				`description:"用户的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	Language			string			`description:"用户的语言，简体中文为zh_CN" json:"language"`
	City				string			`description:"用户所在城市" json:"city"`
	Province			string			`description:"用户所在省份" json:"province"`
	Country				string			`description:"用户所在国家" json:"country"`
	Headimgurl			string			`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
	SubscribeTime		int64			`description:"用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间" json:"subscribeTime"`
	Unionid				string			`description:"只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。" json:"unionid"`
	Remark				string			`description:"公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注" json:"remark"`
	Groupid				int				`description:"用户所在的分组ID（兼容旧的用户分组接口）" json:"groupid"`
	TagidList			string			`description:"用户被打上的标签ID列表" json:"tagidist"`
	SubscribeScene		string			`description:"返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENEPROFILE LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他" json:"subscribeScene"`
	QrScene				int64			`description:"二维码扫码场景（开发者自定义）" json:"qrScene"`
	QrSceneStr			string			`description:"二维码扫码场景描述（开发者自定义）" json:"qrSceneStr"`
	Channel      		string			`description:"渠道" json:"channel"`
	IsNew	      		int				`description:"1为新用户 2为老用户，新用户：活动期内：自然关注的、海报扫码的（助力）、渠道扫码关注进入的用户。老用户：活动期外自然关注、海报、渠道扫码进入的用户。" json:"isNew"`
	UnsubscribeTime		int64			`description:"取消订阅时间" json:"unsubscribeTime"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//---------------------结构体-----------------------------

type UserShort struct {
	UId       			int64			`description:"uId" json:"uId"`
	PhoneNumber			string			`description:"手机号" json:"phoneNumber"`
	NickName 			string			`description:"昵称（登录账号）" json:"nickName" xorm:"notnull "`
	Email	 			string			`description:"电子邮箱" json:"email"`
	RoleType 			string			`description:"角色，1为管理员（可创建用户）" json:"roleType"`
	DefaultAuthInfoId	int64			`description:"默认公众号id" json:"defaultAuthInfoId"`
	DefaultAuthAppid	string 			`description:"默认授权方appid" json:"defaultAuthAppid"`
	Created           	int64  			`description:"注册时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type UserShortContainer struct {
	UserShort			UserShort		`description:"用户信息" json:"user"`
	AuthInfoShort    	AuthInfoShort	`description:"用户权限下公众号信息" json:"authInfo"`
}

type UserRedis struct {
	EncryptedPassword	string			`description:"用户密码" json:"encryptedPassword"`
	RoleType			string			`description:"用户角色" json:"roleType"`
}

type AuthorRedis struct {
	AuthorId			int64			`description:"authorId" json:"authorId"`
	Openid				string			`description:"openid" json:"openid"`
	WebOpenid			string			`description:"webOpenid" json:"webOpenid"`
	AppOpenid			string			`description:"appOpenid" json:"appOpenid"`
	Unionid				string			`description:"unionid" json:"unionid"`
}

type UserShortListContainer struct {
	UserShortList		[]UserShort		`description:"用户信息" json:"userList"`
}

//-----------------------方法-----------------------------

func (u *User) UserToUserShort() (userDTO *UserShort, error error) {
	jsonByte, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	userD := UserShort{}
	err = json.Unmarshal(jsonByte, &userD)
	if err != nil {
		return nil, err
	}
	return &userD, nil
}


