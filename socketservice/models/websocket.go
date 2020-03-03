package models

import (
	"golang.org/x/net/websocket"
)

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

//socket签名key
const SOCKET_MESSAGE_SIGN_KEY string = "wenshixiong123socketmessage"

//连接存储
type SocketConnection struct {
	Conn				*websocket.Conn			`description:"socket连接" json:"conn"`
	ConnType				int					`description:"socket连接类型，1前台，2后台" json:"connType"`
	ExpireTime				int64				`description:"socket连接有效截止时间" json:"expireTime"`
	Token					string				`description:"用户token" json:"token"`
}

type AuthorSocket struct {
	Author					Author				`description:"用户" json:"author"`
	Connection				SocketConnection	`description:"链接" json:"connection"`
}

type AuthorSocketList struct {
	AuthorSocketList		[]AuthorSocket		`description:"用户列表" json:"authorList"`
}
