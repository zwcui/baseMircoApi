package models

//socket统一结构
type SocketMessage struct {
	MessageType				int					`description:"消息类型，-1为后台建立连接，0为前台建立连接，2为被挤下线，3为直播间消息" json:"messageType" `
	MessageSendTime			int64				`description:"消息发送时间" json:"messageSendTime" `
	MessageSenderUid		int64				`description:"消息发送uid" json:"messageSenderUid" `
	MessageReceiverUid		int64				`description:"消息接受uid" json:"messageReceiverUid" `
	MessageExpireTime		int64				`description:"心跳有效时间" json:"messageExpireTime" `
	MessageContent			string				`description:"消息内容，jsonString" json:"messageContent" `
	MessageSign				string				`description:"消息签名" json:"messageSign" `
	MessageToken			string				`description:"用户token" json:"messageToken" `
}

//进出直播间消息
type LiveRoomSocketMessage struct {
	MessageType				int					`description:"消息类型，1为进入直播间推送，2为离开直播间推送，3为直播间关闭，4为连麦请求列表，5为通知老师开始连麦，6为通知观众开始连麦，7为通知老师结束连麦，8为通知观众结束连麦" json:"messageType" `
	CurrentNum				int					`description:"当前直播间人数" json:"currentNum" `
	ConnectionAuthorId		int64				`description:"当前连麦id" json:"connectionAuthorId" `
	CurrentSortNo			int					`description:"当前用户排序号" json:"currentSortNo" `
	Content					string				`description:"直播间消息" json:"content" `
}

//socket签名key
const SOCKET_MESSAGE_SIGN_KEY string = "wenshixiong123socketmessage"


