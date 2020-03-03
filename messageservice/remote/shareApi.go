/*
@Time : 2019/9/29 下午3:30 
@Author : lianwu
@File : share.go
@Software: GoLand
*/
package remote

import (
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro"
	"jingting_server/messageservice/base"
	"jingting_server/messageservice/util"
	"context"
	"jingting_server/messageservice/proto/share"
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

//调用videoservice的生成绑定关系服务
func CreateShare(lastAuthorId int64, authorId int64) (ok bool, error error) {

	service := getRemoteService()

	pushClient := go_micro_videoservice.NewCreateShareService("go.micro.videoservice", (*service).Client())

	rsp, err := pushClient.CreateShare(context.Background(), &go_micro_videoservice.CreateShareParam{LastAuthorId:lastAuthorId, AuthorId:authorId})
	util.Logger.Info("remote CreateShare rsp:")
	util.Logger.Info(rsp)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
