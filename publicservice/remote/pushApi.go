/*
@Time : 2019/9/4 下午1:35 
@Author : zwcui
@Software: GoLand
*/
package remote

import (
	"jingting_server/publicservice/models"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"jingting_server/publicservice/base"
	"github.com/micro/go-micro"
	"context"
	"jingting_server/publicservice/util"
	"strconv"
	"jingting_server/publicservice/proto/message"
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

