/*
@Time : 2019/11/1 上午10:03 
@Author : lianwu
@File : storeApi.go
@Software: GoLand
*/
package remote

import (
	"context"
	"jingting_server/messageservice/proto/store"
	"jingting_server/messageservice/util"
)

//调用storeservice的创建店铺服务
func CreateStore(store string) (ok bool, error error) {

	service := getRemoteService()

	createClient := go_micro_storeservice.NewCreateStoreService("go.micro.storeservice", (*service).Client())

	rsp, err := createClient.CreateStore(context.Background(), &go_micro_storeservice.CreateStoreParam{Store:store})
	util.Logger.Info("remote CreateShare rsp:")
	util.Logger.Info(rsp)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

