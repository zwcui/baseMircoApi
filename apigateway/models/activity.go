/*
@Time : 2019/2/27 下午4:35 
@Author : zwcui
@Software: GoLand
*/
package models


//裂变活动
type Activity struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthInfoId				int64				`description:"公众号id" json:"authInfoId"`
	Step					int					`description:"新增裂变活动的步骤" json:"step"`
	Name       				string				`description:"活动名称" json:"name"`
	KeyWords       			string				`description:"关键词" json:"keyWords"`
	Status       			int					`description:"状态，1启用 2停止 3结束" json:"status" xorm:"notnull default 0"`
	Content       			string				`description:"触发关键词推送文案" json:"content"`
	Banner       			string				`description:"裂变海报" json:"banner"`
	BannerAvatarX			int					`description:"裂变海报用户头像x坐标" json:"bannerAvatarX"`
	BannerAvatarY			int					`description:"裂变海报用户头像y坐标" json:"bannerAvatarY"`
	BannerAvatarRadius 		int					`description:"裂变海报用户头像半径" json:"bannerAvatarRadius"`
	BannerNickNameX			int					`description:"裂变海报用户昵称x坐标" json:"bannerNickNameX"`
	BannerNickNameY			int					`description:"裂变海报用户昵称y坐标" json:"bannerNickNameY"`
	BannerNickNameFontSize  int					`description:"裂变海报用户昵称尺寸" json:"bannerNickNameFontSize"`
	BannerNickNameColor		string				`description:"裂变海报用户昵称颜色，16进制" json:"bannerNickNameColor"`
	BannerQrcodeX		  	int					`description:"裂变海报二维码x坐标" json:"bannerQrcodeX"`
	BannerQrcodeY		  	int					`description:"裂变海报二维码y坐标" json:"bannerQrcodeY"`
	BannerQrcodeSideLength 	int					`description:"裂变海报二维码边长" json:"bannerQrcodeSideLength"`
	FissionCount			int					`description:"裂变人数" json:"fissionCount"`
	EndType       			string				`description:"结束类型，1时间，2奖品 多个,分隔" json:"endType"`
	EndTime       			int64				`description:"结束时间" json:"endTime"`
	ActualEndTime       	int64				`description:"实际结束时间" json:"actualEndTime"`
	EndPrizeName       		string				`description:"结束奖品名称" json:"endPrizeName"`
	EndPrizeCount       	int					`description:"结束奖品个数" json:"endPrizeCount"`
	LeftPrizeCount       	int					`description:"剩余奖品个数" json:"leftPrizeCount"`
	EndContent     			string				`description:"结束文案" json:"endContent"`
	EndPicture     			string				`description:"结束图片" json:"endPicture"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//裂变详情（根据不同规则设定不同文案）
type FissionInfo struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	ActivityId				int64				`description:"活动id" json:"activityId"`
	Level					int					`description:"裂变级别" json:"level"`
	Type					int					`description:"类型" json:"type"`
	TypeName				string				`description:"类型名称" json:"typeName"`
	Text					string				`description:"文本" json:"text"`
	Picture					string				`description:"图片" json:"picture"`
	SortNo					string				`description:"排序号" json:"sortNo"`
	Status					int					`description:"状态，1启用，2未启用" json:"status" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//活动渠道二维码
type ActivityQrcode struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	ActivityId				int64				`description:"活动id" json:"activityId"`
	Channel   				string				`description:"渠道" json:"channel"`
	Ticket       			string				`description:"获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。" json:"ticket"`
	ExpireTime				int64				`description:"微信二维码失效时间" json:"expireTime"`
	Url						string				`description:"二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片" json:"url"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//裂变人员，关注后则成为裂变人员
//扫了包含活动信息的码，或者自然进来（加入所有进行中活动），才有该数据
//已关注的用户(新老用户)扫码没有数据
type FissionMember struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Openid					string				`description:"裂变人" json:"openid"`
	LastOpenid				string				`description:"上级裂变人" json:"lastOpenid"`
	ActivityId				int64				`description:"活动id" json:"activityId"`
	Channel					string				`description:"渠道" json:"channel"`
	Status					int					`description:"状态，1为正常 2为取消关注" json:"status"`
	UnsubscribeTime			int64				`description:"取消订阅时间" json:"unsubscribeTime"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//活动加入表
//取关status=2
type JoinActivity struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Openid					string				`description:"活动加入人" json:"openid"`
	LastOpenid				string				`description:"上级活动加入人" json:"lastOpenid"`
	ActivityId				int64				`description:"活动id" json:"activityId"`
	Level					int					`description:"活动层级，最上面为1级" json:"level"`
	DirectNextSubscribeNum	int					`description:"直接下线订阅人数（指裂变人数，fissionMember）" json:"directNextSubscribeNum" xorm:"notnull default 0"`
	DirectNextUnsubscribeNum int				`description:"直接下线取消订阅人数（指裂变人数，fissionMember）" json:"directNextUnsubscribeNum" xorm:"notnull default 0"`
	AllNextSubscribeNum		int					`description:"所有下线订阅人数（指裂变人数，fissionMember）" json:"allNextSubscribeNum" xorm:"notnull default 0"`
	AllNextUnsubscribeNum	int					`description:"所有下线取消订阅人数（指裂变人数，fissionMember）" json:"allNextUnsubscribeNum" xorm:"notnull default 0"`
	BannerMediaId			string				`description:"裂变海报MediaId" json:"bannerMediaId"`
	Channel					string				`description:"渠道" json:"channel"`
	Status					int					`description:"状态，1为正常" json:"status"`
	IsRewarded				int					`description:"是否得到奖励" json:"isRewarded" xorm:"notnull default 0"`
	FinishTime				int64				`description:"完成任务时间" json:"finishTime"`
	PushCount				int					`description:"助力好友提醒推送次数" json:"pushCount"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//获取奖品记录
type ReceivePrizeRecord struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthInfoId				int64				`description:"公众号id" json:"authInfoId"`
	ActivityId				int64				`description:"活动id" json:"activityId"`
	ReceiverOpenid			string				`description:"获取奖品人Openid" json:"receiverOpenid"`
	ReceiverSex				int					`description:"获取奖品人性别" json:"receiverSex"`
	ReceiverNickname		string				`description:"获取奖品人昵称" json:"receiverNickname"`
	ReceiverRealName		string				`description:"获取奖品人真实姓名" json:"receiverRealName"`
	ReceiverPhoneNumber		string				`description:"获取奖品人手机号" json:"receiverPhoneNumber"`
	ReceiverArea			string				`description:"寄送奖品地区" json:"receiverArea"`
	ReceiverAddress			string				`description:"寄送奖品地址" json:"receiverAddress"`
	TrackingNumber			string				`description:"快递单号" json:"trackingNumber"`
	Status					int					`description:"状态，0为未寄送，1为已寄送" json:"status" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//活动报名关键词记录
type ActivityKeywordRecord struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	ActivityId				int64				`description:"活动id" json:"activityId"`
	Openid					string				`description:"活动加入人" json:"openid"`
	Keyword					string				`description:"关键字" json:"keyword"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//活动统计
type ActivityStatistics struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	ActivityId				int64				`description:"活动id" json:"activityId"`
	TotalSubscribeNum		int64				`description:"总新增关注人数 活动期间内所有关注公众号人数（裂变人数）" json:"totalSubscribeNum" xorm:"notnull default 0"`
	TotalUnsubscribeNum		int64				`description:"活动期内所有取关人数" json:"totalUnsubscribeNum" xorm:"notnull default 0"`
	MaxLevel				int64				`description:"最大裂变层数" json:"maxLevel" xorm:"notnull default 0"`
	TotalJoinNum			int64				`description:"活动报名人数（含老用户）" json:"totalJoinNum" xorm:"notnull default 0"`
	TotalFinishNum			int64				`description:"完成任务人数" json:"totalFinishNum" xorm:"notnull default 0"`
	TotalReceivePrizeNum	int64				`description:"领取奖品人数" json:"totalReceivePrizeNum" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//活动渠道统计
type ActivityChannelStatistics struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	ActivityId				int64				`description:"活动id" json:"activityId"`
	Channel   				string				`description:"渠道" json:"channel"`
	TotalSubscribeNum		int64				`description:"总新增关注人数 活动期间内所有关注公众号人数（裂变人数）" json:"totalSubscribeNum" xorm:"notnull default 0"`
	TotalUnsubscribeNum		int64				`description:"活动期内所有取关人数" json:"totalUnsubscribeNum" xorm:"notnull default 0"`
	MaxLevel				int64				`description:"最大裂变层数" json:"maxLevel" xorm:"notnull default 0"`
	TotalJoinNum			int64				`description:"活动报名人数（含老用户）" json:"totalJoinNum" xorm:"notnull default 0"`
	TotalFinishNum			int64				`description:"完成任务人数" json:"totalFinishNum" xorm:"notnull default 0"`
	TotalReceivePrizeNum	int64				`description:"领取奖品人数" json:"totalReceivePrizeNum" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//----------------------结构体---------------------------
type ActicityInfoContainer struct {
	Activity				Activity			`description:"裂变活动" json:"activity"`
	FissionList				[]FissionInfo		`description:"裂变详情" json:"fissionList"`
}

type ActicityListContainer struct {
	Activity				[]Activity			`description:"裂变活动" json:"activityList"`
}

type ActivityQrcodeListContainer struct {
	ActivityQrcodeList		[]ActivityQrcode	`description:"活动渠道二维码列表" json:"activityQrcodeList"`
}

type ReceivePrizeRecordListContainer struct {
	BaseListContainer
	ReceivePrizeRecordList  []ReceivePrizeRecordDetail `description:"获取奖品列表" json:"receivePrizeRecordList"`
}

type ActivityStatisticsDetail struct {
	Activity				Activity			`description:"活动" json:"activity"`
	ActivityStatistics		ActivityStatistics	`description:"活动统计" json:"activityStatistics"`
}

type ActivityStatisticsListContainer struct {
	BaseListContainer
	ActivityStatisticsDetailList  []ActivityStatisticsDetail `description:"活动统计列表" json:"activityStatisticsDetailList"`
}

type ActivityStatisticsContainer struct {
	Activity							Activity					`description:"活动" json:"activity"`
	ActivityStatistics					ActivityStatistics			`description:"活动统计" json:"activityStatistics"`
	ActivityChannelStatisticsList		[]ActivityChannelStatistics	`description:"活动渠道统计" json:"activityChannelStatisticsList"`
}

type JoinActivityDetail struct {
	JoinActivity			`description:"活动加入人" xorm:"extends"`
	Headimgurl				string			`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
	Nickname				string			`description:"用户的昵称" json:"nickname"`
	Sex						int				`description:"用户的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	City					string			`description:"用户所在城市" json:"city"`
	Province				string			`description:"用户所在省份" json:"province"`
	Country					string			`description:"用户所在国家" json:"country"`
}

type JoinActivityDetailContainer struct {
	BaseListContainer
	JoinActivityDetailList	[]JoinActivityDetail	`description:"活动加入人" json:"joinActivityDetailList"`
}

type FissionMemberDetail struct {
	JoinActivity			`description:"活动加入人" xorm:"extends"`
	Subscribe       		int				`description:"用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。" json:"subscribe"`
	SubscriberOpenid		string			`description:"用户的标识，对当前公众号唯一" json:"subscriberOpenid"`
	UnsubscribeTime			int64			`description:"取消订阅时间" json:"unsubscribeTime"`
	Headimgurl				string			`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
	Nickname				string			`description:"用户的昵称" json:"nickname"`
	Sex						int				`description:"用户的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	City					string			`description:"用户所在城市" json:"city"`
	Province				string			`description:"用户所在省份" json:"province"`
	Country					string			`description:"用户所在国家" json:"country"`
}

type FissionMemberDetailContainer struct {
	BaseListContainer
	FissionMemberDetailList	[]FissionMemberDetail	`description:"活动裂变人" json:"fissionMemberDetailList"`
}

type ReceivePrizeRecordDetail struct {
	ReceivePrizeRecord						`description:"活动加入人" xorm:"extends"`
	FinishTime				int64			`description:"完成任务时间" json:"finishTime"`
}

//-----------------------微信接口请求返回体-----------------------------
type MediaUploadJsonBody struct {
	MediaType    			string   			`description:"type" json:"type"`
	MediaId    				string   			`description:"media_id" json:"media_id"`
	CreatedAt    			int64   			`description:"created_at" json:"created_at"`
}

type QrcodeJsonBody struct {
	Ticket       			string				`description:"获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。" json:"ticket"`
	ExpireSeconds			int64				`description:"该二维码有效时间，以秒为单位。 最大不超过2592000（即30天）。" json:"expire_seconds"`
	Url						string				`description:"二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片" json:"url"`
}

//用户个人信息
type UserInfoJsonBody struct {
	Subscribe       		int					`description:"用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。" json:"subscribe"`
	Openid					string				`description:"用户的标识，对当前公众号唯一" json:"openid"`
	Nickname				string				`description:"用户的昵称" json:"nickname"`
	Sex						int					`description:"用户的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	Language				string				`description:"用户的语言，简体中文为zh_CN" json:"language"`
	City					string				`description:"用户所在城市" json:"city"`
	Province				string				`description:"用户所在省份" json:"province"`
	Country					string				`description:"用户所在国家" json:"country"`
	Headimgurl				string				`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
	SubscribeTime			int64				`description:"用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间" json:"subscribe_time"`
	Unionid					string				`description:"只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。" json:"unionid"`
	Remark					string				`description:"公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注" json:"remark"`
	Groupid					int					`description:"用户所在的分组ID（兼容旧的用户分组接口）" json:"groupid"`
	TagidList				[]int				`description:"用户被打上的标签ID列表" json:"tagid_list"`
	SubscribeScene			string				`description:"返回用户关注的渠道来源，ADD_SCENE_SEARCH 公众号搜索，ADD_SCENE_ACCOUNT_MIGRATION 公众号迁移，ADD_SCENE_PROFILE_CARD 名片分享，ADD_SCENE_QR_CODE 扫描二维码，ADD_SCENEPROFILE LINK 图文页内名称点击，ADD_SCENE_PROFILE_ITEM 图文页右上角菜单，ADD_SCENE_PAID 支付后关注，ADD_SCENE_OTHERS 其他" json:"subscribe_scene"`
	QrScene					int64				`description:"二维码扫码场景（开发者自定义）" json:"qr_scene"`
	QrSceneStr				string				`description:"二维码扫码场景描述（开发者自定义）" json:"qr_scene_str"`
}