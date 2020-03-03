/*
@Time : 2019/9/4 下午1:35 
@Author : zwcui
@Software: GoLand
*/
package remote

import (
	"crypto/md5"
	"encoding/hex"
	"jingting_server/taskservice/models"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"jingting_server/taskservice/base"
	"github.com/micro/go-micro"
	"jingting_server/taskservice/proto/message"
	"context"
	go_micro_socketMessageservice "jingting_server/taskservice/proto/socketMessage"
	"jingting_server/taskservice/util"
	"sort"
	"strconv"
	"strings"
)

func getRemoteService() *micro.Service {
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			base.ConsulAddress,
		}
	})

	// 初始化服务
	service := micro.NewService(
		micro.Registry(reg),
	)
	service.Init()
	return &service
}

//调用messageservice的推送服务
func PushMessageToUser(authorId int64, message *models.AppMessage, sound string, badge int) (ok bool, error error) {

	service := getRemoteService()

	pushClient := go_micro_messageservice.NewPushMessageToUserService("go.micro.messageservice", (*service).Client())

	rsp, err := pushClient.PushMessageToUser(context.Background(), &go_micro_messageservice.PushMessageToUserParam{Content:message.Content, ReceiverId:strconv.FormatInt(authorId, 10), ActionUrl:message.ActionUrl, Type:strconv.Itoa(message.Type)})
	util.Logger.Info("remote PushMessageToUser rsp:")
	util.Logger.Info(rsp)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

//调用messageservice的生成actionurl
func JumpUrlWithKeyAndPramas(key string, params map[string]string) (urlStr string) {

	service := getRemoteService()

	pushClient := go_micro_messageservice.NewJumpUrlWithKeyAndParamsService("go.micro.messageservice", (*service).Client())

	rsp, err := pushClient.JumpUrlWithKeyAndParams(context.Background(), &go_micro_messageservice.JumpUrlWithKeyAndParamsParam{Key:key, Map:params})
	util.Logger.Info("remote JumpUrlWithKeyAndPramas rsp:")
	util.Logger.Info(rsp)
	if err != nil {
		util.Logger.Info("remote JumpUrlWithKeyAndPramas err:"+err.Error())
		return ""
	} else {
		return rsp.UrlStr
	}
}

//调用socketservice的推送服务
func PushSocketMessageToUser(authorId int64, content string ) (ok bool, error error) {

	service := getRemoteService()

	pushClient := go_micro_socketMessageservice.NewPushSocketMessageToUserService("go.micro.socketservice", (*service).Client())

	rsp, err := pushClient.PushMessageSocketToUser(context.Background(), &go_micro_socketMessageservice.PushSocketMessageToUserParam{Content:content, ReceiverId:strconv.FormatInt(authorId, 10)})
	util.Logger.Info("remote PushSocketMessageToUser rsp:")
	util.Logger.Info(rsp)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
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
