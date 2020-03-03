/*
@Time : 2019/2/27 下午2:43 
@Author : zwcui
@Software: GoLand
*/
package models

import "encoding/xml"

//当普通微信用户向公众账号发消息时，微信服务器将POST消息的XML数据包到开发者填写的URL上

//加密后每隔10分钟定时收到component_verify_ticket
type ReceiveMessageEncryptXmlBody struct {
	XMLName       			xml.Name 		`xml:"xml" structs:"-"`
	ToUserName	   			string			`xml:"ToUserName" structs:"ToUserName"`
	Encrypt    				string   		`xml:"Encrypt" structs:"Encrypt"`
}

//在微信用户和公众号产生交互的过程中，用户的某些操作会使得微信服务器通过事件推送的形式通知到开发者在开发者中心处设置的服务器地址
type ReceiveMessageXmlBody struct {
	XMLName       				xml.Name 		`xml:"xml" structs:"-"`
	ToUserName    				string   		`description:"开发者微信号" xml:"ToUserName" structs:"ToUserName"`
	FromUserName     			string   		`description:"发送方帐号（一个OpenID）" xml:"FromUserName" structs:"FromUserName"`
	CreateTime        			int   			`description:"消息创建时间" xml:"CreateTime" structs:"CreateTime"`
	MsgType   					string   		`description:"text(文本消息),image(图片消息),voice(语音消息),video(视频消息),shortvideo(小视频消息),location(地理位置消息),link(链接消息),event(关注/取消关注事件)" xml:"MsgType" structs:"MsgType"`
	Content   					string   		`description:"文本消息内容" xml:"Content" structs:"Content"`
	MsgId   					int64   		`description:"消息id" xml:"MsgId" structs:"MsgId"`
	PicUrl   					string   		`description:"图片链接（由系统生成）" xml:"PicUrl" structs:"PicUrl"`
	MediaId   					string   		`description:"图片/语音消息媒体id，可以调用多媒体文件下载接口拉取数据" xml:"MediaId" structs:"MediaId"`
	Format   					string   		`description:"语音格式，如amr，speex等" xml:"Format" structs:"Format"`
	Recognition   				string   		`description:"语音识别结果，UTF8编码" xml:"Recognition" structs:"Recognition"`
	ThumbMediaId   				string   		`description:"视频消息缩略图的媒体id，可以调用多媒体文件下载接口拉取数据。" xml:"ThumbMediaId" structs:"ThumbMediaId"`
	Location_X   				string   		`description:"地理位置维度" xml:"Location_X" structs:"Location_X"`
	Location_Y   				string   		`description:"地理位置经度" xml:"Location_Y" structs:"Location_Y"`
	Scale   					string   		`description:"地图缩放大小" xml:"Scale" structs:"Scale"`
	Label   					string   		`description:"地理位置信息" xml:"Label" structs:"Label"`
	Title   					string   		`description:"消息标题" xml:"Title" structs:"Title"`
	Description   				string   		`description:"消息描述" xml:"Description" structs:"Description"`
	Url   						string   		`description:"消息链接" xml:"Url" structs:"Url"`
	Event  						string   		`description:"事件类型，subscribe(订阅)、unsubscribe(取消订阅)、SCAN(扫描带参数二维码，用户已关注时的事件推送)、LOCATION(上报地理位置事件)、CLICK(自定义菜单事件:点击菜单拉取消息时的事件推送)、VIEW(自定义菜单事件:点击菜单跳转链接时的事件推送)" xml:"Event" structs:"Event"`
	EventKey  					string   		`description:"事件KEY值，qrscene_为前缀，后面为二维码的参数值(用户未关注时，进行关注后的事件推送)   事件KEY值，是一个32位无符号整数，即创建二维码时的二维码scene_id(用户已关注时的事件推送)   事件KEY值，与自定义菜单接口中KEY值对应(点击菜单拉取消息时的事件推送)   事件KEY值，设置的跳转URL(点击菜单跳转链接时的事件推送)" xml:"EventKey" structs:"EventKey"`
	Ticket  					string   		`description:"二维码的ticket，可用来换取二维码图片" xml:"Ticket" structs:"Ticket"`
	Latitude  					string   		`description:"地理位置纬度" xml:"Latitude" structs:"Latitude"`
	Longitude  					string   		`description:"地理位置经度" xml:"Longitude" structs:"Longitude"`
	Precision  					string   		`description:"地理位置精度" xml:"Precision" structs:"Precision"`
}

//获取已添加至帐号下所有模板列表
type TemplateListJsonBody struct {
	TemplateList				[]Template		`description:"模板列表" json:"template_list"`
}

type Template struct {
	TemplateId       			string			`description:"模板ID" json:"template_id"`
	Title       				string			`description:"模板标题" json:"title"`
	PrimaryIndustry       		string			`description:"模板所属行业的一级行业" json:"primary_industry"`
	DeputyIndustry       		string			`description:"模板所属行业的二级行业" json:"deputy_industry"`
	Content       				string			`description:"模板内容" json:"content"`
	Example       				string			`description:"模板示例" json:"example"`
}

//微信接口异常信息
type MessageErrorJsonBody struct {
	Errcode						int 			`description:"错误码" json:"errcode"`
	Errmsg						string 			`description:"错误描述" json:"errmsg"`
	Msgid						int64 			`description:"msgid" json:"msgid"`
}

//------------------------表结构----------------------------------
//存储用户在公众号发的消息
type Message struct {
	Id       					int64			`description:"id" json:"id" xorm:"pk autoincr"`
	Appid       				string			`description:"appid" json:"appid"`
	MsgId       				int64			`description:"消息id" json:"msgId"`
	FromUserName       			string			`description:"发送方帐号（一个OpenID）" json:"fromUserName"`
	ToUserName       			string			`description:"开发者微信号" json:"toUserName"`
	MsgType   					string   		`description:"text(文本消息),image(图片消息),voice(语音消息),video(视频消息),shortvideo(小视频消息),location(地理位置消息),link(链接消息),event(关注/取消关注事件)" json:"msgType"`
	Content   					string   		`description:"文本消息内容" json:"content"`
	PicUrl   					string   		`description:"图片链接（由系统生成）" json:"picUrl"`
	MediaId   					string   		`description:"图片/语音消息媒体id，可以调用多媒体文件下载接口拉取数据" json:"mediaId"`
	Format   					string   		`description:"语音格式，如amr，speex等" json:"format"`
	Recognition   				string   		`description:"语音识别结果，UTF8编码" json:"recognition"`
	ThumbMediaId   				string   		`description:"视频消息缩略图的媒体id，可以调用多媒体文件下载接口拉取数据。" json:"thumbMediaId"`
	LocationX   				string   		`description:"地理位置维度" json:"location_X"`
	LocationY   				string   		`description:"地理位置经度" json:"location_Y"`
	Scale   					string   		`description:"地图缩放大小" json:"scale"`
	Label   					string   		`description:"地理位置信息" json:"label"`
	Title   					string   		`description:"消息标题" json:"title"`
	Description   				string   		`description:"消息描述" json:"description"`
	Url   						string   		`description:"消息链接" json:"url"`
	Event  						string   		`description:"事件类型，subscribe(订阅)、unsubscribe(取消订阅)、SCAN(扫描带参数二维码，用户已关注时的事件推送)、LOCATION(上报地理位置事件)、CLICK(自定义菜单事件:点击菜单拉取消息时的事件推送)、VIEW(自定义菜单事件:点击菜单跳转链接时的事件推送)" json:"event"`
	EventKey  					string   		`description:"事件KEY值，qrscene_为前缀，后面为二维码的参数值(用户未关注时，进行关注后的事件推送)   事件KEY值，是一个32位无符号整数，即创建二维码时的二维码scene_id(用户已关注时的事件推送)   事件KEY值，与自定义菜单接口中KEY值对应(点击菜单拉取消息时的事件推送)   事件KEY值，设置的跳转URL(点击菜单跳转链接时的事件推送)" json:"eventKey"`
	Ticket  					string   		`description:"二维码的ticket，可用来换取二维码图片" json:"ticket"`
	Latitude  					string   		`description:"地理位置纬度" json:"latitude"`
	Longitude  					string   		`description:"地理位置经度" json:"longitude"`
	Precision  					string   		`description:"地理位置精度" json:"precision"`
	Created           			int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           			int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         			int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//模板消息构成数据
type TemplateMessageData struct {
	Id       					int64			`description:"id" json:"id" xorm:"pk autoincr"`
	AuthInfoId					int64			`description:"公众号id" json:"authInfoId"`
	TemplateId					string			`description:"模板号" json:"templateId"`
	TemplateTitle				string			`description:"模板标题" json:"templateTitle"`
	PrimaryIndustry				string			`description:"模板所属行业的一级行业" json:"primaryIndustry"`
	DeputyIndustry				string			`description:"模板所属行业的二级行业" json:"deputyIndustry"`
	Data						string			`description:"模板数据" json:"deputyIndustry"`
	Url							string			`description:"模板跳转链接（海外帐号没有跳转能力）" json:"url"`
	SendTime					int64			`description:"发送时间" json:"sendTime"`
	SendSex						int 			`description:"发送性别，0全部 1男 2女" json:"sendSex"`
	SendProvince				string			`description:"发送省份" json:"sendProvince"`
	SendCity					string			`description:"发送城市" json:"sendCity"`
	Status						int				`description:"发送状态，0未发送，1已发送" json:"status"`
	Created           			int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           			int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         			int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//客服消息模板
type CustomerServiceTemplateMessage struct {
	Id       					int64			`description:"id" json:"id" xorm:"pk autoincr"`
	AuthInfoId					int64			`description:"公众号id" json:"authInfoId"`
	Title						string			`description:"标题" json:"title"`
	Description					string			`description:"描述，客服消息文字内容" json:"description"`
	Url							string			`description:"url" json:"url"`
	Picurl						string			`description:"picurl" json:"picurl"`
	Created           			int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           			int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         			int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//发送短信response
type SendSMSMessageResponse struct {
	XMLName   xml.Name `xml:"response"`
	ErrorCode int      `xml:"error"`
	Message   string   `xml:"message"`
}

// app推送消息
type AppMessage struct {
	Id          int64  `description:"消息id" json:"id" xorm:"pk autoincr"`
	ReceiverId int64  `description:"接收者id 为0则是推送给所有人" json:"receiverId" `
	ActionUrl   string `description:"跳转url" json:"actionUrl" xorm:"text"`
	Content     string `description:"内容" json:"content" xorm:"text"`
	Type        int    `description:"消息类型 1:购买成功消息 2:转发成功消息 3:评论成功消息 4:举报成功消息 5.下级关注通知 6：循环绑定通知 7:问卷回答成功通知" json:"type" `
	Created     int64  `description:"创建时间" json:"created" xorm:"created"`
	DeletedAt   int64  `description:"删除时间" json:"-" xorm:"deleted"`
}

//--------------------------结构体--------------------------------
type TemplateListContainer struct {
	TemplateList				[]Template		`description:"模板列表" json:"templateList"`
}

type ValueAndColor struct {
	Value		       			string			`description:"模板内容" json:"value"`
	Color       				string			`description:"模板内容字体颜色，不填默认为黑色" json:"color"`
}

type TemplateMessageDataListContainer struct {
	BaseListContainer
	TemplateMessageDataList		[]TemplateMessageData	`description:"模板内容" json:"templateMessageDataList"`
	TotalUnsend					int64					`description:"未发送数量" json:"totalUnsend"`
	TotalSend					int64					`description:"发送数量" json:"totalSend"`
}

type CustomerServiceTemplateMessageListContainer struct {
	BaseListContainer
	CustomerServiceTemplateMessageList		[]CustomerServiceTemplateMessage	`description:"模板内容" json:"customerServiceTemplateMessageList"`
}