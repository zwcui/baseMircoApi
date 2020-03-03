/*
@Time : 2019/9/20 下午1:59 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	go_micro_socketMessageservice "jingting_server/socketservice/proto/socketMessage"
	"strconv"
	"jingting_server/socketservice/util"
	"jingting_server/socketservice/models"
	"jingting_server/socketservice/base"
)

type SocketController struct {
	apiController
}

func (this *SocketController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
		//{"/getAuthorSocketOnlineList", "get", []int{}},
	}
	this.userAuth()
}

//@Title 获取socket在线用户列表
//@Description 获取socket在线用户列表
//@Success 200 {object} []models.AuthorSocketList
//@router /getAuthorSocketOnlineList [get]
func (this *SocketController) GetAuthorSocketOnlineList() {

	uIdList := make([]int64, 0)
	util.Logger.Info("len(UserSocketConnections)="+strconv.Itoa(len(UserSocketConnections)))
	for uId, _ := range UserSocketConnections {
		uIdList = append(uIdList, uId)
	}

	var authorList []models.Author
	if uIdList != nil && len(uIdList) > 0 {
		base.DBEngine.Table("author").In("id", uIdList).Find(&authorList)
	}

	if authorList == nil {
		authorList = make([]models.Author, 0)
	}

	var authorSocketList []models.AuthorSocket
	for _, author := range authorList {
		var authorSocket models.AuthorSocket
		authorSocket.Author = author
		authorSocket.Connection = UserSocketConnections[author.Id]
		authorSocketList = append(authorSocketList, authorSocket)
	}

	if authorSocketList == nil {
		authorSocketList = make([]models.AuthorSocket, 0)
	}

	this.ReturnData = models.AuthorSocketList{authorSocketList}
}




//-----------------------------方法------------------------------------
func (this *SocketController) PushMessageSocketToUser(ctx context.Context, req *go_micro_socketMessageservice.PushSocketMessageToUserParam, rsp *go_micro_socketMessageservice.PushSocketMessageToUserResponse) error{
	util.Logger.Info("service PushMessageSocketToUser")
	util.Logger.Info("req")
	util.Logger.Info(req)

	receiverId, _ := strconv.ParseInt(req.ReceiverId, 10, 64)
	socketConnection, ok := UserSocketConnections[receiverId]
	content := req.Content
	if ok {
		err := websocket.Message.Send(socketConnection.Conn, content)
		if err != nil {
			util.Logger.Info("PushMessageSocketToUser  err = ", err.Error())
			this.ReturnData = util.GenerateAlertMessage(models.PushMessageError100)
			return err
		}
	}else {
		util.Logger.Info("PushMessageSocketToUser  err = socketConnectionErr")
		this.ReturnData = util.GenerateAlertMessage(models.PushMessageError100)
		return errors.New("push socket message failed: user is offline")
	}
	rsp.Result = "success"
	return nil
}




