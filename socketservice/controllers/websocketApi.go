package controllers

import (
	"net/http"
	"golang.org/x/net/websocket"
	"strconv"
	"encoding/json"
	"sort"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"jingting_server/socketservice/util"
	"jingting_server/socketservice/models"
)

//socket连接池
var UserSocketConnections map[int64]models.SocketConnection

type WSServer struct {
	ListenAddr string
}

func (this *WSServer) Handler (conn *websocket.Conn) {
	util.Logger.Info("-------------socket--Handler--start------------"+strconv.FormatInt(util.UnixOfBeijingTime(), 10))
	if UserSocketConnections == nil {
		UserSocketConnections = make(map[int64]models.SocketConnection)
	}

	util.Logger.Info("a new ws conn: conn.RemoteAddr()="+conn.RemoteAddr().String()+"  conn.LocalAddr()="+conn.LocalAddr().String())
	var err error
	//Try(func() {
	for {
		var reply string
		var socketMessage models.SocketMessage
		err = websocket.Message.Receive(conn, &reply)
		//连接出错则移除
		if err != nil {
			for k, v := range UserSocketConnections {
				if v.Conn == conn {
					delete(UserSocketConnections, k)
				}
			}
			util.Logger.Info("receive conn err:",err.Error())
			break
		}

		//util.Logger.Info("-----reply----  "+reply)
		err = json.Unmarshal([]byte(reply), &socketMessage)
		if err != nil {
			util.Logger.Info("----socketMessage--json.Unmarshal--err---- "+err.Error())
			continue
		}

		//验签
		if socketMessage.MessageSign != SignMessage(socketMessage) {
			util.Logger.Info("----socketMessage--签名验证失败, reply="+reply+"  SignMessage="+SignMessage(socketMessage))
			break
		}

		//建立心跳，区分前台后台，并当后台转至前台时补发离线消息
		if socketMessage.MessageType == -1 || socketMessage.MessageType == 0 {
			handleHeartbeat(&socketMessage, conn)
		}


	}
	//}, func(e interface{}) {
	//	util.Logger.Info("-------------------------造成循环退出err1-----------------------------------------")
	//	util.Logger.Info(e.(error).Error())
	//	util.Logger.Info("-------------------------造成循环退出err2-----------------------------------------")
	//})
	util.Logger.Info("-------------socket--Handler--end------------"+strconv.FormatInt(util.UnixOfBeijingTime(), 10))
}

func (this *WSServer) Start() (error) {
	http.Handle("/ws", websocket.Handler(this.Handler))
	util.Logger.Info("websocket----begin to listen")
	err := http.ListenAndServe(this.ListenAddr, nil)
	if err != nil {
		util.Logger.Info("ListenAndServe:", err)
		return err
	}
	util.Logger.Info("websocket----start end")
	return nil
}

//处理心跳
func handleHeartbeat(socketMessage *models.SocketMessage, conn *websocket.Conn){
	if _, ok := UserSocketConnections[socketMessage.MessageSenderUid]; !ok {
		var socketConnection models.SocketConnection
		socketConnection.Conn = conn
		if socketMessage.MessageType == 0 {
			socketConnection.ConnType = 1
		} else if socketMessage.MessageType == -1 {
			socketConnection.ConnType = 2
		}
		socketConnection.ExpireTime = socketMessage.MessageExpireTime
		socketConnection.Token = socketMessage.MessageToken
		UserSocketConnections[socketMessage.MessageSenderUid] = socketConnection
		util.Logger.Info("-----socketConnection.heartbeat.ExpireTime----start-"+strconv.FormatInt(socketMessage.MessageSenderUid, 10)+"--"+strconv.FormatInt(UserSocketConnections[socketMessage.MessageSenderUid].ExpireTime, 10))
		//建立心跳发送未接收的消息
		handleUnsentSocketMessage(socketMessage, conn)
	} else {
		//以token作为唯一标示
		if UserSocketConnections[socketMessage.MessageSenderUid].Token != socketMessage.MessageToken {
			//账户被挤下线
			util.Logger.Info("-------您的账户已在其他地方登陆-------"+strconv.FormatInt(socketMessage.MessageSenderUid, 10))

			var kickOffSocketMessage models.SocketMessage
			kickOffSocketMessage.MessageType = 2
			kickOffSocketMessage.MessageSendTime = util.UnixOfBeijingTime()
			kickOffSocketMessage.MessageSenderUid = 0
			kickOffSocketMessage.MessageReceiverUid = socketMessage.MessageSenderUid
			kickOffSocketMessage.MessageExpireTime = 0
			kickOffSocketMessage.MessageContent = "您的账户已在其他地方登陆"
			kickOffSocketMessage.MessageSign = SignMessage(kickOffSocketMessage)
			kickOffSocketMessageBytes, _ := json.Marshal(kickOffSocketMessage)
			if err := websocket.Message.Send(UserSocketConnections[socketMessage.MessageSenderUid].Conn, string(kickOffSocketMessageBytes)); err != nil {
				util.Logger.Info("----userMessage--websocket.Message.Send 您的账户已在其他地方登陆 err:", err.Error())
				//移除出错的链接
				delete(UserSocketConnections, socketMessage.MessageSenderUid)
			}
			//新登陆的账户
			var socketConnection models.SocketConnection
			socketConnection.Conn = conn
			if socketMessage.MessageType == 0 {
				socketConnection.ConnType = 1
			} else if socketMessage.MessageType == -1 {
				socketConnection.ConnType = 2
			}
			socketConnection.ExpireTime = socketMessage.MessageExpireTime
			socketConnection.Token = socketMessage.MessageToken
			UserSocketConnections[socketMessage.MessageSenderUid] = socketConnection
		} else {
			//查看是否后台切换至前台，如果是则补发离线消息
			storedSocketConnection := UserSocketConnections[socketMessage.MessageSenderUid]
			if storedSocketConnection.ConnType == 2 && socketMessage.MessageType == 0 {
				//查看redis缓存消息，仅后台跳至前台补发直连消息socket，如果被取消，则只发取消的socket
				handleUnsentSocketMessage(socketMessage, conn)
			}

			util.Logger.Info("-----socketConnection.heartbeat.ExpireTime-----"+strconv.FormatInt(socketMessage.MessageSenderUid, 10)+"--"+strconv.FormatInt(UserSocketConnections[socketMessage.MessageSenderUid].ExpireTime, 10))
			//切换网络刷新conn，现为每次心跳刷新conn
			if conn.Request().Header.Get("deviceToken") != "" && storedSocketConnection.Conn.Request().Header.Get("deviceToken") == conn.Request().Header.Get("deviceToken") {
				util.Logger.Info("refresh storedSocketConnection")
				var socketConnection models.SocketConnection
				socketConnection.Conn = conn
				if socketMessage.MessageType == 0 {
					socketConnection.ConnType = 1
				} else if socketMessage.MessageType == -1 {
					socketConnection.ConnType = 2
				}
				socketConnection.ExpireTime = socketMessage.MessageExpireTime
				socketConnection.Token = socketMessage.MessageToken
				UserSocketConnections[socketMessage.MessageSenderUid] = socketConnection
			} else {
				socketConnection := UserSocketConnections[socketMessage.MessageSenderUid]
				if socketMessage.MessageType == 0 {
					socketConnection.ConnType = 1
				} else if socketMessage.MessageType == -1 {
					socketConnection.ConnType = 2
				}
				socketConnection.ExpireTime = socketMessage.MessageExpireTime
				UserSocketConnections[socketMessage.MessageSenderUid] = socketConnection
			}
		}
	}
}


//redis缓存消息转发
func handleUnsentSocketMessage(socketMessage *models.SocketMessage, conn *websocket.Conn){


}


//签名
func SignMessage(socketMessage models.SocketMessage) string {
	params := make(map[string]string)
	params["messageType"] = strconv.Itoa(socketMessage.MessageType)
	params["messageSendTime"] = strconv.FormatInt(socketMessage.MessageSendTime, 10)
	params["messageSenderUid"] = strconv.FormatInt(socketMessage.MessageSenderUid, 10)
	params["messageReceiverUid"] = strconv.FormatInt(socketMessage.MessageReceiverUid, 10)
	params["messageExpireTime"] = strconv.FormatInt(socketMessage.MessageExpireTime, 10)
	params["messageContent"] = socketMessage.MessageContent
	params["messageToken"] = socketMessage.MessageToken
	//params["messageAppNo"] = strconv.Itoa(socketMessage.MessageAppNo)

	keys := make([]string, len(params))

	i := 0
	for k := range params {
		keys[i] = k
		i++
	}

	//util.Logger.Info("keys", keys)
	sort.Strings(keys)
	//util.Logger.Info("sorted keys", keys)

	strTemp := ""
	for _, key := range keys {
		strTemp = strTemp + key + "=" + params[key] + "&"
	}
	strTemp += "key=" + models.SOCKET_MESSAGE_SIGN_KEY
	//util.Logger.Info("strTemp = ", strTemp)

	hasher := md5.New()
	hasher.Write([]byte(strTemp))
	md5Str := hex.EncodeToString(hasher.Sum(nil))

	//util.Logger.Info("md5 = ", md5Str)
	return strings.ToUpper(md5Str)
}

//每分钟检查失效的socket连接
func CheckSocketHeartbeat(){
	util.Logger.Info("定时任务，每分钟检查失效的socket连接")
	for uId, socketConnection := range UserSocketConnections {
		util.Logger.Info("-----定时任务  遍历users-----  util.UnixOfBeijingTime()="+strconv.FormatInt(util.UnixOfBeijingTime(), 10)+"   uid="+strconv.FormatInt(uId, 10)+"   ExpireTime="+strconv.FormatInt(socketConnection.ExpireTime, 10))
		if (socketConnection.ExpireTime + 15) <= util.UnixOfBeijingTime() {
			delete(UserSocketConnections, uId)
		}
	}
}

