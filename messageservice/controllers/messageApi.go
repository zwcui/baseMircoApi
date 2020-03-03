/*
@Time : 2019/2/27 下午2:36 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"jingting_server/messageservice/models"
	"jingting_server/messageservice/util"
	"encoding/xml"
	"jingting_server/messageservice/base"
	"encoding/base64"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"github.com/go-xorm/xorm"
	"runtime"
	"math/rand"
	"github.com/astaxie/beego"
	"jingting_server/messageservice/remote"
)

type MessageController struct {
	apiController
}

func (this *MessageController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/addTemplateMessageData", "post", []int{0}},
		{"/getTemplateMessageDataList", "get", []int{0}},
		{"/deleteTemplateMessageData", "delete", []int{0}},
		{"/addCustomerServiceTemplateMessage", "post", []int{0}},
		{"/getCustomerServiceTemplateMessageList", "get", []int{0}},
		{"/deleteCustomerServiceTemplateMessage", "delete", []int{0}},
	}
	this.userAuth()
}

// @Title 微信用户向公众账号发消息转发至此（微信使用）
// @Description 微信用户向公众账号发消息转发至此（微信使用）
// @Success 200 {string} success
// @router /receiveMessage/:appId [post]
func (this *MessageController) ReceiveMessage() {
	appId := this.MustString(":appId")

	response := models.ReceiveMessageEncryptXmlBody{}
	//util.Logger.Info("receiveMessage ReceiveMessageEncryptXmlBody body = ", string(this.Ctx.Input.RequestBody))
	err := xml.Unmarshal(this.Ctx.Input.RequestBody, &response)
	if err != nil {
		util.Logger.Info("xml.Unmarshal body = " + err.Error())
		this.ServeXML()
		this.StopRun()
	}

	AESKey, _ := base64.StdEncoding.DecodeString(models.JTThirdPartyPlatformEncodingAESKey + "=")

	// Decode base64
	cipherData, err := base64.StdEncoding.DecodeString(response.Encrypt)
	if err != nil {
		util.Logger.Info("Wechat Service: Decode base64 error:", err.Error())
		return
	}
	//util.Logger.Info("Decode base64:   "+string(cipherData))

	// AES Decrypt
	plainData, err := util.AesDecrypt(cipherData, AESKey)
	if err != nil {
		util.Logger.Info("util.AesDecrypt(cipherData, aesKey):", err.Error())
		return
	}
	//util.Logger.Info("util.AesDecrypt(cipherData, aesKey):   "+string(plainData))

	//Xml decoding
	receiveMessageXmlBody, _ := util.ParseEncryptTextRequestBodyToReceiveMessage(plainData)
	util.Logger.Info(receiveMessageXmlBody)

	//避免5秒内无法处理，如果微信重试消息则直接返回空
	if receiveMessageXmlBody.MsgId != 0 {
		hasStoredMessage, _ := base.DBEngine.Table("message").Where("msg_id=?", receiveMessageXmlBody.MsgId).And("created>?", util.UnixOfBeijingTime() - 60).Get(new(models.Message))
		if hasStoredMessage {
			util.Logger.Info("hasStoredMessage wx retry")
			util.Logger.Info(receiveMessageXmlBody)
			this.IsDirectReturn = 1
			this.ReturnData = "success"
			return
		}
	}

	//util.Logger.Info("runtime.NumCPU()")
	//util.Logger.Info(runtime.NumCPU())
	util.Logger.Info("runtime.NumGoroutine()")
	util.Logger.Info(runtime.NumGoroutine())

	//上线检测
	//回复关键词
	var authInfo models.AuthInfo
	base.DBEngine.Table("auth_info").Where("auth_appid=?", appId).Get(&authInfo)

	if receiveMessageXmlBody.MsgType == "text" && receiveMessageXmlBody.Content == "TESTCOMPONENT_MSG_TYPE_TEXT" {
		//response, err := util.MakeTextResponseBody(receiveMessageXmlBody.ToUserName, receiveMessageXmlBody.FromUserName, "TESTCOMPONENT_MSG_TYPE_TEXT_callback")
		response, err := util.MakeEncryptResponseBody(receiveMessageXmlBody.ToUserName, receiveMessageXmlBody.FromUserName, "TESTCOMPONENT_MSG_TYPE_TEXT_callback", strconv.Itoa(rand.Int()), strconv.FormatInt(util.UnixOfBeijingTime(), 10))
		if err != nil {
			util.Logger.Info("err:"+err.Error())
		}

		this.IsDirectReturn = 1
		this.IsXml = 1
		this.TestData = response
		//this.ReturnData = "<xml><ToUserName><![CDATA["+receiveMessageXmlBody.FromUserName+"]]></ToUserName><FromUserName><![CDATA["+receiveMessageXmlBody.ToUserName+"]]></FromUserName><CreateTime>"+strconv.FormatInt(util.UnixOfBeijingTime(), 10)+"</CreateTime><MsgType><![CDATA[text]]></MsgType><Content><![CDATA[TESTCOMPONENT_MSG_TYPE_TEXT_callback]]></Content></xml>"
		//this.ReturnData = "<xml><ToUserName>"+receiveMessageXmlBody.FromUserName+"</ToUserName><FromUserName>"+receiveMessageXmlBody.ToUserName+"</FromUserName><CreateTime>"+strconv.FormatInt(util.UnixOfBeijingTime(), 10)+"</CreateTime><MsgType>text</MsgType><Content>TESTCOMPONENT_MSG_TYPE_TEXT_callback</Content></xml>"

		return
	} else if receiveMessageXmlBody.MsgType == "text" && strings.HasPrefix(receiveMessageXmlBody.Content, "QUERY_AUTH_CODE:") {
		this.IsDirectReturn = 1
		this.ReturnData = ""

		job := func() {
			//onlineTest(receiveMessageXmlBody, &authInfo)
		}
		base.GoPool.Submit(job)

		//go onlineTest(receiveMessageXmlBody, &authInfo)

		return
	} else {
		//回复空串让微信服务器不再重试
		this.IsDirectReturn = 1
		this.ReturnData = "success"


		job := func() {
			handleMessage(appId, authInfo, receiveMessageXmlBody)
		}
		base.GoPool.Submit(job)

		//协程处理
		//go handleMessage(appId, authInfo, receiveMessageXmlBody)

		return

	}
}

//func onlineTest(receiveMessageXmlBody *models.ReceiveMessageXmlBody, authInfo *models.AuthInfo){
//	authInfo.AuthCode = strings.Split(receiveMessageXmlBody.Content, ":")[1]
//	base.DBEngine.Table("auth_info").Where("auth_appid=?", authInfo.AuthAppid).Cols("auth_code").Update(&authInfo)
//
//	err := requestAuthorizerAccessToken(authInfo)
//	if err != nil {
//		util.Logger.Info("err:"+err.Error())
//	}
//
//	err = util.RequestSendTextCustomerServiceMessage(receiveMessageXmlBody.FromUserName, authInfo.AuthCode+"_from_api", *authInfo)
//	if err != nil {
//		util.Logger.Info("err:"+err.Error())
//	} else {
//		util.Logger.Info("send online test msg: "+authInfo.AuthCode+"_from_api")
//	}
//
//}

//处理消息
func handleMessage(appId string, authInfo models.AuthInfo, receiveMessageXmlBody *models.ReceiveMessageXmlBody){
	util.Logger.Info("handleMessage  start")
	util.Logger.Info("running goroutines:" + strconv.Itoa(base.GoPool.Running()))

	//存库
	var message models.Message
	message.Appid = appId
	message.MsgId = receiveMessageXmlBody.MsgId
	message.FromUserName = receiveMessageXmlBody.FromUserName
	message.ToUserName = receiveMessageXmlBody.ToUserName
	message.MsgType = receiveMessageXmlBody.MsgType
	message.Content = receiveMessageXmlBody.Content
	message.PicUrl = receiveMessageXmlBody.PicUrl
	message.MediaId = receiveMessageXmlBody.MediaId
	message.Format = receiveMessageXmlBody.Format
	message.Recognition = receiveMessageXmlBody.Recognition
	message.ThumbMediaId = receiveMessageXmlBody.ThumbMediaId
	message.LocationX = receiveMessageXmlBody.Location_X
	message.LocationY = receiveMessageXmlBody.Location_Y
	message.Scale = receiveMessageXmlBody.Scale
	message.Label = receiveMessageXmlBody.Label
	message.Title = receiveMessageXmlBody.Title
	message.Description = receiveMessageXmlBody.Description
	message.Url = receiveMessageXmlBody.Url
	message.Event = receiveMessageXmlBody.Event
	message.EventKey = receiveMessageXmlBody.EventKey
	message.Ticket = receiveMessageXmlBody.Ticket
	message.Latitude = receiveMessageXmlBody.Latitude
	message.Longitude = receiveMessageXmlBody.Longitude
	message.Precision = receiveMessageXmlBody.Precision
	base.DBEngine.Table("message").InsertOne(&message)

	//解析参数
	paramMap := make(map[string]string)
	if message.Event != "VIEW" && message.EventKey != "" {
		if strings.HasPrefix(message.EventKey, "qrscene") {
			param := util.SubstrByLength(message.EventKey, 8, len(message.EventKey)-8)
			for _, kv := range strings.Split(param, "&") {
				if kv == "" {
					continue
				}
				key := strings.Split(kv, "=")[0]
				value := strings.Split(kv, "=")[1]
				paramMap[key] = value
			}
		} else {
			for _, kv := range strings.Split(message.EventKey, "&") {
				if kv == "" {
					continue
				}
				key := strings.Split(kv, "=")[0]
				value := strings.Split(kv, "=")[1]
				paramMap[key] = value
			}
		}
	}
	util.Logger.Info("paramMap")
	util.Logger.Info(paramMap)

	//扫码或搜索关注订阅号
	if message.MsgType == "event" && message.Event == "subscribe" {
		//加入订阅者记录
		var subscriber models.Subscriber
		hasSubscriber, _ := base.DBEngine.Table("subscriber").Where("openid=?", message.FromUserName).And("auth_info_id=?", authInfo.Id).Get(&subscriber)
		if !hasSubscriber {
			userInfo, err := requestUserInfoByOpenId(message.FromUserName, authInfo)
			if err != nil {
				util.Logger.Info("requestUserInfoByOpenId err:"+err.Error())
			} else {
				subscriber.Openid = message.FromUserName
				subscriber.AuthInfoId = authInfo.Id
				subscriber.Subscribe = userInfo.Subscribe
				subscriber.Nickname = userInfo.Nickname
				subscriber.Sex = userInfo.Sex
				subscriber.Language = userInfo.Language
				subscriber.City = userInfo.City
				subscriber.Province = userInfo.Province
				subscriber.Country = userInfo.Country
				subscriber.Headimgurl = userInfo.Headimgurl
				subscriber.SubscribeTime = userInfo.SubscribeTime
				subscriber.Unionid = userInfo.Unionid
				subscriber.Remark = userInfo.Remark
				subscriber.Groupid = userInfo.Groupid
				for _, tagid := range userInfo.TagidList {
					subscriber.TagidList += strconv.Itoa(tagid) + ","
				}
				subscriber.SubscribeScene = userInfo.SubscribeScene
				subscriber.QrScene = userInfo.QrScene
				subscriber.QrSceneStr = userInfo.QrSceneStr
				if channel, ok := paramMap["channel"]; ok {
					subscriber.Channel = channel
				}
				subscriber.IsNew = checkSubscriberIsNew(authInfo)
				base.DBEngine.Table("subscriber").InsertOne(&subscriber)
			}
		} else if subscriber.Subscribe == 0 {
			userInfo, err := requestUserInfoByOpenId(message.FromUserName, authInfo)
			if err != nil {
				util.Logger.Info("requestUserInfoByOpenId err:"+err.Error())
			} else {
				subscriber.Openid = message.FromUserName
				subscriber.AuthInfoId = authInfo.Id
				subscriber.Subscribe = userInfo.Subscribe
				subscriber.Nickname = userInfo.Nickname
				subscriber.Sex = userInfo.Sex
				subscriber.Language = userInfo.Language
				subscriber.City = userInfo.City
				subscriber.Province = userInfo.Province
				subscriber.Country = userInfo.Country
				subscriber.Headimgurl = userInfo.Headimgurl
				subscriber.SubscribeTime = userInfo.SubscribeTime
				subscriber.Unionid = userInfo.Unionid
				subscriber.Remark = userInfo.Remark
				subscriber.Groupid = userInfo.Groupid
				for _, tagid := range userInfo.TagidList {
					subscriber.TagidList += strconv.Itoa(tagid) + ","
				}
				subscriber.SubscribeScene = userInfo.SubscribeScene
				subscriber.QrScene = userInfo.QrScene
				subscriber.QrSceneStr = userInfo.QrSceneStr
				if channel, ok := paramMap["channel"]; ok {
					subscriber.Channel = channel
				} else {
					subscriber.Channel = ""
				}
				subscriber.IsNew = checkSubscriberIsNew(authInfo)
				subscriber.Created = util.UnixOfBeijingTime()
				base.DBEngine.Table("subscriber").Where("id=?", subscriber.Id).AllCols().Update(&subscriber)
			}
		}

		//如果是裂变人员则加入裂变关系
		becomeFissionMember(subscriber, paramMap, authInfo)

		//关注自动注册视频脉脉
		if loginUuid, ok := paramMap["loginUuid"]; ok {
			createAuthorByOpenid(loginUuid, &subscriber, paramMap)
		}

		//关注或扫码，参数有观看记录信息，则推送客服消息告知
		pushVideoViewHistoryMessage(message.FromUserName, authInfo, paramMap)

		//分享名片后关注公众号永久绑定上下级关系
		if lastAuthorId, ok := paramMap["lastAuthorId"]; ok {
			var author models.Author
			base.DBEngine.Table("author").Where("openid = ?", message.FromUserName).Get(&author)
			id, _ := strconv.ParseInt(lastAuthorId, 10, 64)
			_, err := remote.CreateShare(id, author.Id)
			if err != nil {
				util.Logger.Info("CreateShare err = ", err.Error())
			}
			//推送消息给上级
			message := models.AppMessage{}
			message.Content = "【您已成功邀请好友"+ author.Nickname + " 加入辣课，可在个人中心->我的推广中查看下级明细】"
			message.ReceiverId = id
			message.ActionUrl = JumpUrlWithKeyAndPramas(models.JTFOLLOW_JUMP_KEY, nil)
			message.Type = 5
			_, err = PushMessageToUser(id, &message, "", 0)
			if err != nil {
				util.Logger.Info("PushMessageToUser  err = ", err.Error())
			}
			base.DBEngine.Table("app_message").InsertOne(&message)

		}



	}

	//已关注再扫码
	if message.MsgType == "event" && message.Event == "SCAN" {
		//已关注用户通过朋友海报再次关注，不助力上级

		//反复助力上级则弹出提示
		if activityId, ok := paramMap["activity"]; ok {
			var fissionMember models.FissionMember
			hasFissionMember, _ := base.DBEngine.Table("fission_member").Where("openid=?", message.FromUserName).And("activity_id=?", activityId).And("status=1").And("last_openid is not null and last_openid != ''").Get(&fissionMember)
			if hasFissionMember {
				//用户重复助力提醒文案(为防止用户多次助力好友)
				var fissionInfo04 models.FissionInfo
				hasFissionInfo04, _ := base.DBEngine.Table("fission_info").Where("activity_id=?", activityId).And("type=4").And("status=1").Get(&fissionInfo04)
				if hasFissionInfo04 {
					content := fissionInfo04.Text
					util.RequestSendTextCustomerServiceMessage(message.FromUserName, content, authInfo)
				}
			} else {
				//新关注用户才能给好友助力（已关注公众号的用户，可以参加活动但不能为好友助力）
				var fissionInfo07 models.FissionInfo
				hasFissionInfo07, _ := base.DBEngine.Table("fission_info").Where("activity_id=?", activityId).And("type=7").And("status=1").Get(&fissionInfo07)
				if hasFissionInfo07 {
					content := fissionInfo07.Text
					util.RequestSendTextCustomerServiceMessage(message.FromUserName, content, authInfo)
				}
			}
		}

		//扫码，更新loginUuid，网站登录使用
		if loginUuid, ok := paramMap["loginUuid"]; ok {
			var subscriber models.Subscriber
			hasSubscriber, _ := base.DBEngine.Table("subscriber").Where("openid=?", message.FromUserName).And("auth_info_id=?", authInfo.Id).Get(&subscriber)
			if hasSubscriber {
				//扫码自动刷新视频脉脉账户信息
				createAuthorByOpenid(loginUuid, &subscriber, paramMap)
			} else {
				util.Logger.Info("event SCAN !hasSubscriber openid=" + message.FromUserName + " auth_info_id=" + strconv.FormatInt(authInfo.Id, 10))
			}
		}

		//关注或扫码，参数有观看记录信息，则推送客服消息告知
		pushVideoViewHistoryMessage(message.FromUserName, authInfo, paramMap)
	}

	//取消订阅
	if message.MsgType == "event" && message.Event == "unsubscribe" {
		var subscriber models.Subscriber
		base.DBEngine.Table("subscriber").Where("openid=?", message.FromUserName).And("auth_info_id=?", authInfo.Id).Get(&subscriber)

		//活动统计
		updateActivityStatisticsWhenUnsubscribe(subscriber, authInfo)

		//取消订阅不再助力
		cancelFissionMember(subscriber, authInfo)

		//取消订阅取消活动资格
		leaveActivity(subscriber, authInfo)

		//取消订阅更新订阅者信息
		subscriber.Subscribe = 0
		subscriber.UnsubscribeTime = util.UnixOfBeijingTime()
		base.DBEngine.Table("subscriber").Where("id=?", subscriber.Id).Cols("subscribe").Update(&subscriber)
	}

	//回复关键词
	if message.MsgType == "text" && message.Content != "" {
		//判断是否满足活动要求
		var activity models.Activity
		hasActivity, _ := base.DBEngine.Table("activity").Where("FIND_IN_SET('"+message.Content+"', key_words)").And("status=1").And("auth_info_id=?", authInfo.Id).Get(&activity)
		if hasActivity {
			var subscriber models.Subscriber
			base.DBEngine.Table("subscriber").Where("openid=?", message.FromUserName).And("auth_info_id=?", authInfo.Id).Get(&subscriber)

			//判断加入几次，如果到达三次则不能加入
			joinCountSql := "select count(1) from join_activity where openid='"+subscriber.Openid+"' and activity_id='"+strconv.FormatInt(activity.Id, 10)+"' and deleted_at is null "
			total, err := base.DBEngine.SQL(joinCountSql).Count(new(models.JoinActivity))
			if err != nil {
				util.Logger.Info("err:" + err.Error())
				return
			}
			if total>= 3 {
				util.RequestSendTextCustomerServiceMessage(subscriber.Openid, "您已重复报名取关3次，无法参与此活动", authInfo)
				return
			}

			//记录活动报名关键词（取关重新回复则重新记录）
			hasJoinActivity, _ := base.DBEngine.Table("join_activity").Where("openid=?", message.FromUserName).And("activity_id=?", activity.Id).And("status=1").Get(new(models.JoinActivity))
			if !hasJoinActivity {
				var activityKeywordRecord models.ActivityKeywordRecord
				activityKeywordRecord.ActivityId = activity.Id
				activityKeywordRecord.Openid = message.FromUserName
				activityKeywordRecord.Keyword = message.Content
				base.DBEngine.Table("activity_keyword_record").InsertOne(&activityKeywordRecord)
			}

			//1.发送对应文案
			if !hasJoinActivity {
				content := activity.Content
				content = strings.Replace(content, "[用户昵称]", subscriber.Nickname, -1)
				content = strings.Replace(content, "[奖品剩余数量]", strconv.Itoa(activity.LeftPrizeCount), -1)
				util.RequestSendTextCustomerServiceMessage(message.FromUserName, content, authInfo)
			}

			//2.生成海报并发送，扫码进入则渠道取自于码，否则取自于第一次关注的渠道
			mediaId, err := joinActivity(subscriber, activity, authInfo)
			if err != nil {
				util.Logger.Info("err:"+err.Error())
			} else {
				if mediaId != "" {
					util.RequestSendMediaCustomerServiceMessage(message.FromUserName, mediaId, authInfo)
				}
			}
		}
	}



	util.Logger.Info("handleMessage  end")
}







// @Title 获取已添加至帐号下所有模板列表（h5使用）
// @Description 获取已添加至帐号下所有模板列表（h5使用）
// @Param	authInfoId					query			int64	  		true		"公众号id"
// @Success 200 {string} success
// @router /getTemplateList [get]
func (this *MessageController) GetTemplateList() {
	authInfoId := this.MustInt64("authInfoId")

	var authInfo models.AuthInfo
	hasAuth, _ := base.DBEngine.Table("auth_info").Where("id=?", authInfoId).Get(&authInfo)
	if !hasAuth {
		this.ReturnData = util.GenerateAlertMessage(models.AuthorizeError300)
		return
	}

	templateList, err := util.RequestQueryTemplateList(authInfo)
	if err != nil {
		this.ReturnData = util.GenerateAlertMessage(models.MessageError100, err.Error())
		return
	}

	if templateList == nil {
		templateList = make([]models.Template, 0)
	}

	this.ReturnData = models.TemplateListContainer{templateList}
}

// @Title 新增模板数据消息（h5使用）
// @Description 新增模板数据消息（h5使用）
// @Param	authInfoId					formData			int64	  		true		"公众号id"
// @Param	templateId					formData			string	  		true		"模板号"
// @Param	templateTitle				formData			string	  		false		"模板标题"
// @Param	primaryIndustry				formData			string	  		false		"模板所属行业的一级行业"
// @Param	deputyIndustry				formData			string	  		false		"模板所属行业的二级行业"
// @Param	url							formData			string	  		false		"url"
// @Param	data						formData			string	  		true		"模板数据"
// @Param	sendTime					formData			int64	  		false		"发送时间"
// @Param	sendSex						formData			int 	  		false		"发送性别"
// @Param	sendProvince				formData			string	  		false		"发送省份"
// @Param	sendCity					formData			string	  		false		"发送城市"
// @Success 200 {string} success
// @router /addTemplateMessageData [post]
func (this *MessageController) AddTemplateMessageData() {
	authInfoId := this.MustInt64("authInfoId")
	templateId := this.MustString("templateId")
	templateTitle := this.GetString("templateTitle", "")
	primaryIndustry := this.GetString("primaryIndustry", "")
	deputyIndustry := this.GetString("deputyIndustry", "")
	url := this.GetString("url", "")
	data := this.GetString("data", "")
	sendTime, _ := this.GetInt64("sendTime", 0)
	sendSex, _ := this.GetInt("sendSex", 0)
	sendProvince := this.GetString("sendProvince", "")
	sendCity := this.GetString("sendCity", "")

	var templateMessageData models.TemplateMessageData
	templateMessageData.AuthInfoId = authInfoId
	templateMessageData.TemplateId = templateId
	templateMessageData.TemplateTitle = templateTitle
	templateMessageData.PrimaryIndustry = primaryIndustry
	templateMessageData.DeputyIndustry = deputyIndustry
	templateMessageData.Url = url
	templateMessageData.Data = data
	templateMessageData.SendTime = sendTime
	templateMessageData.SendSex = sendSex
	templateMessageData.SendProvince = sendProvince
	templateMessageData.SendCity = sendCity
	base.DBEngine.Table("template_message_data").InsertOne(&templateMessageData)

	this.ReturnData = "success"
}

// @Title 模板数据消息列表（h5使用）
// @Description 模板数据消息列表（h5使用）
// @Param	authInfoId					query				int64	  		true		"公众号id"
// @Param	templateId					query				string	  		false		"模板号"
// @Param	templateTitle				query				string	  		false		"模板标题"
// @Param	status						query				int		  		false		"发送状态，0未发送，1已发送"
// @Param	pageNum						query 	  			int				true		"page num start from 1"
// @Param	pageTime					query 	  			int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize					query 	  			int				false		"page size default is 15"
// @Success 200 {object} models.TemplateMessageDataListContainer
// @router /getTemplateMessageDataList [get]
func (this *MessageController) GetTemplateMessageDataList() {
	authInfoId := this.MustInt64("authInfoId")
	templateId := this.GetString("templateId", "")
	templateTitle := this.GetString("templateTitle", "")
	status, _ := this.GetInt64("status", -1)
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	var total int64
	var totalUnsend int64
	var totalSend int64
	var totalErr error

	total, totalErr = base.DBEngine.Table("template_message_data").Where("auth_info_id=?", authInfoId).And("template_id like '%"+templateId+"%'").And("template_title like '%"+templateTitle+"%'").Count(new(models.TemplateMessageData))
	totalUnsend, totalErr = base.DBEngine.Table("template_message_data").Where("auth_info_id=?", authInfoId).And("template_id like '%"+templateId+"%'").And("template_title like '%"+templateTitle+"%'").And("status=?", 0).Count(new(models.TemplateMessageData))
	totalSend, totalErr = base.DBEngine.Table("template_message_data").Where("auth_info_id=?", authInfoId).And("template_id like '%"+templateId+"%'").And("template_title like '%"+templateTitle+"%'").And("status=?", 1).Count(new(models.TemplateMessageData))

	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	var templateMessageDataList []models.TemplateMessageData
	if total > 0 {
		if status == -1 {
			base.DBEngine.Table("template_message_data").Where("auth_info_id=?", authInfoId).And("template_id like '%" + templateId + "%'").And("template_title like '%" + templateTitle + "%'").Limit(pageSize, pageSize*(pageNum-1)).Desc("created").Find(&templateMessageDataList)
		} else {
			base.DBEngine.Table("template_message_data").Where("auth_info_id=?", authInfoId).And("template_id like '%" + templateId + "%'").And("template_title like '%" + templateTitle + "%'").And("status=?", status).Limit(pageSize, pageSize*(pageNum-1)).Desc("created").Find(&templateMessageDataList)
		}
	}

	if templateMessageDataList == nil {
		templateMessageDataList = make([]models.TemplateMessageData, 0)
	}

	this.ReturnData = models.TemplateMessageDataListContainer{models.BaseListContainer{total, pageNum, pageTime}, templateMessageDataList, totalUnsend, totalSend}
}

// @Title 删除模板数据消息（h5使用）
// @Description 删除模板数据消息（h5使用）
// @Param	templateMessageId					query				int64	  		true		"模板消息id"
// @Success 200 {string} success
// @router /deleteTemplateMessageData [delete]
func (this *MessageController) DeleteTemplateMessageData() {
	templateMessageId := this.MustInt64("templateMessageId")

	var templateMessageData models.TemplateMessageData
	has, _ := base.DBEngine.Table("template_message_data").Where("id=?", templateMessageId).Get(&templateMessageData)
	if !has {
		this.ReturnData = util.GenerateAlertMessage(models.MessageError200)
		return
	}

	base.DBEngine.Table("template_message_data").Where("id=?", templateMessageId).Delete(&templateMessageData)

	this.ReturnData = "success"
}

// @Title 查询所有订阅者省份（h5使用）
// @Description 查询所有订阅者省份（h5使用）
// @Param	authInfoId					query				int64	  		true		"公众号id"
// @Success 200 {string}
// @router /getSubscriberProvince [get]
func (this *MessageController) GetSubscriberProvince() {
	authInfoId := this.MustInt64("authInfoId")

	sql := "select distinct province from subscriber where auth_info_id='"+strconv.FormatInt(authInfoId, 10)+"' and subscribe=1 and deleted_at is null "
	resultsSlice, err := base.DBEngine.Query(sql)
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+err.Error())
		return
	}

	var province string
	for _, value := range resultsSlice {
		province += string(value["province"]) + ","
	}

	this.ReturnData = province
}

// @Title 查询所有订阅者城市（h5使用）
// @Description 查询所有订阅者城市（h5使用）
// @Param	authInfoId					query				int64	  		true		"公众号id"
// @Param	province					query				string	  		true		"省份"
// @Success 200 {string}
// @router /getSubscriberCity [get]
func (this *MessageController) GetSubscriberCity() {
	authInfoId := this.MustInt64("authInfoId")
	province := this.MustString("province")

	sql := "select distinct city from subscriber where auth_info_id='"+strconv.FormatInt(authInfoId, 10)+"' and province='"+province+"' and subscribe=1 and deleted_at is null "
	resultsSlice, err := base.DBEngine.Query(sql)
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+err.Error())
		return
	}

	var city string
	for _, value := range resultsSlice {
		city += string(value["city"]) + ","
	}

	this.ReturnData = city
}


// @Title 新增客服消息模板（h5使用）
// @Description 新增客服消息模板（h5使用）
// @Param	authInfoId					formData			int64	  		true		"公众号id"
// @Param	title						formData			string	  		false		"标题"
// @Param	description					formData			string	  		false		"描述"
// @Param	url							formData			string	  		false		"url"
// @Param	picurl						formData			string	  		false		"picurl"
// @Success 200 {string} success
// @router /addCustomerServiceTemplateMessage [post]
func (this *MessageController) AddCustomerServiceTemplateMessage() {
	authInfoId := this.MustInt64("authInfoId")
	title := this.GetString("title", "")
	description := this.GetString("description", "")
	url := this.GetString("url", "")
	picurl := this.GetString("picurl", "")

	var customerServiceTemplateMessage models.CustomerServiceTemplateMessage
	customerServiceTemplateMessage.AuthInfoId = authInfoId
	customerServiceTemplateMessage.Title = title
	customerServiceTemplateMessage.Description = description
	customerServiceTemplateMessage.Url = url
	customerServiceTemplateMessage.Picurl = picurl
	base.DBEngine.Table("customer_service_template_message").InsertOne(&customerServiceTemplateMessage)

	this.ReturnData = "success"
}

// @Title 查看客服消息模板列表（h5使用）
// @Description 查看客服消息模板列表（h5使用）
// @Param	authInfoId					query				int64	  		true		"公众号id"
// @Param	title						query				string	  		false		"标题"
// @Param	pageNum						query 	  			int				true		"page num start from 1"
// @Param	pageTime					query 	  			int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize					query 	  			int				false		"page size default is 15"
// @Success 200 {object} models.CustomerServiceTemplateMessageListContainer
// @router /getCustomerServiceTemplateMessageList [get]
func (this *MessageController) GetCustomerServiceTemplateMessageList() {
	authInfoId := this.MustInt64("authInfoId")
	title := this.GetString("title", "")
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	total, totalErr := base.DBEngine.Table("customer_service_template_message").Where("auth_info_id=?", authInfoId).And("title like '%"+title+"%'").Count(new(models.TemplateMessageData))
	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	var customerServiceTemplateMessageList []models.CustomerServiceTemplateMessage
	if total > 0 {
		base.DBEngine.Table("customer_service_template_message").Where("auth_info_id=?", authInfoId).And("title like '%"+title+"%'").Limit(pageSize, pageSize*(pageNum-1)).Desc("created").Find(&customerServiceTemplateMessageList)
	}

	if customerServiceTemplateMessageList == nil {
		customerServiceTemplateMessageList = make([]models.CustomerServiceTemplateMessage, 0)
	}

	this.ReturnData = models.CustomerServiceTemplateMessageListContainer{models.BaseListContainer{total, pageNum, pageTime}, customerServiceTemplateMessageList}
}

// @Title 删除客服消息模板（h5使用）
// @Description 删除客服消息模板（h5使用）
// @Param	id							query				int64	  		true		"客服消息模板id"
// @Success 200 {string} success
// @router /deleteCustomerServiceTemplateMessage [delete]
func (this *MessageController) DeleteCustomerServiceTemplateMessage() {
	id := this.MustInt64("id")

	var customerServiceTemplateMessage models.CustomerServiceTemplateMessage
	hasCustomerServiceTemplateMessage, _ := base.DBEngine.Table("customer_service_template_message").Where("id=?", id).Get(&customerServiceTemplateMessage)
	if !hasCustomerServiceTemplateMessage {
		this.ReturnData = util.GenerateAlertMessage(models.MessageError300)
		return
	}

	base.DBEngine.Table("customer_service_template_message").Where("id=?", id).Delete(&customerServiceTemplateMessage)

	this.ReturnData = "success"
}

// @Title 发送客服消息模板（h5使用）
// @Description 发送客服消息模板（h5使用）
// @Param	id							formData			int64	  		true		"客服消息模板id"
// @Param	sex							formData			int		  		false		"1男 2女 0全部"
// @Param	country						formData			string	  		false		"用户所在国家"
// @Param	province					formData			string	  		false		"用户所在省份"
// @Param	city						formData			string	  		false		"用户所在城市"
// @Param	nickNameList				formData			string	  		false		"昵称列表"
// @Success 200 {string} success
// @router /sendCustomerServiceTemplateMessage [post]
func (this *MessageController) SendCustomerServiceTemplateMessage() {
	id := this.MustInt64("id")
	sex, _ := this.GetInt("sex", 0)
	country := this.GetString("country", "")
	province := this.GetString("province", "")
	city := this.GetString("city", "")
	nickNameList := this.GetString("nickNameList", "")

	var customerServiceTemplateMessage models.CustomerServiceTemplateMessage
	hasCustomerServiceTemplateMessage, _ := base.DBEngine.Table("customer_service_template_message").Where("id=?", id).Get(&customerServiceTemplateMessage)
	if !hasCustomerServiceTemplateMessage {
		this.ReturnData = util.GenerateAlertMessage(models.MessageError300)
		return
	}

	var authInfo models.AuthInfo
	base.DBEngine.Table("auth_info").Where("id=?", customerServiceTemplateMessage.AuthInfoId).Get(&authInfo)

	var subscriberList []models.Subscriber
	sql := "select * from subscriber where auth_info_id='"+strconv.FormatInt(customerServiceTemplateMessage.AuthInfoId, 10)+"' and subscribe=1 and deleted_at is null "
	if sex != 0 {
		sql += " and sex='"+strconv.Itoa(sex)+"' "
	}
	if country != "" {
		sql += " and country='"+country+"' "
	}
	if province != "" {
		sql += " and province='"+province+"' "
	}
	if city != "" {
		sql += " and city='"+city+"' "
	}
	if nickNameList != "" {
		sql += " and ( "
		for i, nickName := range strings.Split(nickNameList, ",") {
			if i == 0 {
				sql += " nickname='"+nickName+"' "
			} else {
				sql += " or nickname='"+nickName+"' "
			}
		}
		sql += " ) "
	}
	util.Logger.Info(sql)
	base.DBEngine.SQL(sql).Find(&subscriberList)

	if subscriberList == nil {
		subscriberList = make([]models.Subscriber, 0)
	}

	util.Logger.Info("发送客服消息模板 --> "+strconv.Itoa(len(subscriberList))+"个人")
	isPicAndText := false	//是否图文消息
	isPic := false	//是否纯图片消息
	if customerServiceTemplateMessage.Picurl != "" && customerServiceTemplateMessage.Description != "" {
		isPicAndText = true
	} else if customerServiceTemplateMessage.Picurl != "" && customerServiceTemplateMessage.Description == "" {
		isPic = true
	}

	mediaId := ""
	var err error
	if isPic {
		mediaId, err = util.UploadOnlineMedia(customerServiceTemplateMessage.Picurl, authInfo, 1)
		if err != nil {
			util.Logger.Info("上传客服消息图片  UploadOnlineMedia  err:"+err.Error())
		}
	}

	for _, subscribe := range subscriberList {
		if isPicAndText {
			err := util.RequestSendPicAndTextCustomerServiceMessage(subscribe.Openid, customerServiceTemplateMessage.Title, customerServiceTemplateMessage.Description, customerServiceTemplateMessage.Url, customerServiceTemplateMessage.Picurl, authInfo)
			if err != nil {
				util.Logger.Info("发送客服消息模板  RequestSendPicAndTextCustomerServiceMessage  err:"+err.Error())
			}
		} else if isPic {
			err := util.RequestSendMediaCustomerServiceMessage(subscribe.Openid, mediaId, authInfo)
			if err != nil {
				util.Logger.Info("发送客服消息模板  RequestSendMediaCustomerServiceMessage  err:"+err.Error())
			}
		} else {
			err := util.RequestSendTextCustomerServiceMessage(subscribe.Openid, customerServiceTemplateMessage.Description, authInfo)
			if err != nil {
				util.Logger.Info("发送客服消息模板  RequestSendTextCustomerServiceMessage  err:"+err.Error())
			}
		}
	}

	this.ReturnData = "success"
}



// @Title 我的消息列表
// @Description 我的消息列表
// @Param	authorId			        query				int64	  		true		"用户id"
// @Param	pageNum						query 	  			int				true		"page num start from 1"
// @Param	pageTime					query 	  			int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize					query 	  			int				false		"page size default is 15"
// @Success 200 {object} models.AppMessageListContainer
// @router /getAppMessageList [get]
func (this *MessageController) GetAppMessageList() {
	authorId := this.MustInt64("authorId")
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	total, totalErr := base.DBEngine.Table("app_message").Where("receiver_id = ?", authorId).Count(new(models.AppMessage))
	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	var appMessageList []models.AppMessage
	if total > 0 {
		base.DBEngine.Table("app_message").Where("receiver_id = ?", authorId).Limit(pageSize, pageSize*(pageNum-1)).Desc("created").Find(&appMessageList)
	}

	if appMessageList == nil {
		appMessageList = make([]models.AppMessage, 0)
	}

	this.ReturnData = models.AppMessageListContainer{models.BaseListContainer{total, pageNum, pageTime}, appMessageList}

}

// @Title TestCreateShare
// @Description TestCreateShare
// @Param	lastAuthorId				formData			int64	  		false		"上级分享id"
// @Param	authorId					formData			int64	  		true		"当前分享id"
// @Success 200 {string} success
// @router /testCreateShare [post]
func (this *MessageController) TestCreateShare() {
	lastAuthorId, _ := this.GetInt64("lastAuthorId", 0)
	authorId := this.MustInt64("authorId")

	remote.CreateShare(lastAuthorId, authorId)

	this.ReturnData = "success"
}






//--------------------------业务方法--------------------------------------------

//成为裂变人员
func becomeFissionMember(subscriber models.Subscriber, paramMap map[string]string, authInfo models.AuthInfo) (err error){
	util.Logger.Info("becomeFissionMember start")
	session := base.DBEngine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		util.Logger.Info("err:" + err.Error())
		return err
	}

	activityId, _ := paramMap["activity"]
	lastOpenid, hasOpenid := paramMap["openid"]
	var activity models.Activity
	hasActivity, _ := session.Table("activity").Where("auth_info_id=?", authInfo.Id).And("id=?", activityId).And("status=1").Get(&activity)
	if hasActivity {
		if hasOpenid {
			var lastJoinActivity models.JoinActivity
			hasLastJoinActivity, _ := session.Table("join_activity").Where("openid=?", lastOpenid).And("activity_id=?", activityId).And("status=1").Get(&lastJoinActivity)
			if hasActivity && hasLastJoinActivity {
				var fissionMember models.FissionMember
				fissionMember.Openid = subscriber.Openid
				fissionMember.LastOpenid = lastOpenid
				fissionMember.ActivityId, _ = strconv.ParseInt(activityId, 10, 64)
				if channel, hasChannel := paramMap["channel"]; hasChannel {
					fissionMember.Channel = channel
				}
				fissionMember.Status = 1
				session.Table("fission_member").InsertOne(&fissionMember)

				//推送消息
				//成功助力好友提醒文案（当B用户助力了A用户时，B用户收到的提醒文案）
				var fissionInfo01 models.FissionInfo
				hasFissionInfo01, _ := session.Table("fission_info").Where("activity_id=?", activity.Id).And("type=1").And("status=1").Get(&fissionInfo01)
				if hasFissionInfo01 {
					var lastSubscriber models.Subscriber
					session.Table("subscriber").Where("openid=?", lastOpenid).And("auth_info_id=?", authInfo.Id).Get(&lastSubscriber)
					content := fissionInfo01.Text
					content = strings.Replace(content, "[用户昵称]", subscriber.Nickname, -1)
					content = strings.Replace(content, "[邀请者昵称]", lastSubscriber.Nickname, -1)
					content = strings.Replace(content, "[奖品剩余数量]", strconv.Itoa(activity.LeftPrizeCount), -1)
					util.RequestSendTextCustomerServiceMessage(subscriber.Openid, content, authInfo)
				}
				if lastJoinActivity.PushCount < 6 {
					//收到助力好友成功助力的提醒文案（当B用户助力了A用户时，A用户收到的助力成功提醒文案）
					var fissionInfo02 models.FissionInfo
					hasFissionInfo02, _ := session.Table("fission_info").Where("activity_id=?", activity.Id).And("type=2").And("status=1").Get(&fissionInfo02)
					if hasFissionInfo02 {
						content := fissionInfo02.Text
						content = strings.Replace(content, "[达成条件需要人数]", strconv.Itoa(activity.FissionCount - (lastJoinActivity.DirectNextSubscribeNum + 1)), -1)
						content = strings.Replace(content, "[助力好友昵称]", subscriber.Nickname, -1)
						content = strings.Replace(content, "[奖品剩余数量]", strconv.Itoa(activity.LeftPrizeCount), -1)
						content = strings.Replace(content, "[已有助力好友人数]", strconv.Itoa(lastJoinActivity.DirectNextSubscribeNum + 1), -1)
						util.RequestSendTextCustomerServiceMessage(lastJoinActivity.Openid, content, authInfo)
					}

					//更新推送次数
					lastJoinActivity.PushCount += 1
					session.Table("join_activity").Where("id=?", lastJoinActivity.Id).Cols("push_count").Update(&lastJoinActivity)
				} else if lastJoinActivity.PushCount == 6 {
					//超限助力提醒文案(因微信接口调用限制，8条以后的好友助力将推送此提醒文案)
					var fissionInfo03 models.FissionInfo
					hasFissionInfo03, _ := session.Table("fission_info").Where("activity_id=?", activity.Id).And("type=3").And("status=1").Get(&fissionInfo03)
					if hasFissionInfo03 {
						content := fissionInfo03.Text
						content = strings.Replace(content, "[助力好友昵称]", subscriber.Nickname, -1)
						content = strings.Replace(content, "[奖品剩余数量]", strconv.Itoa(activity.LeftPrizeCount), -1)
						content = strings.Replace(content, "[已有助力好友人数]", strconv.Itoa(lastJoinActivity.DirectNextSubscribeNum + 1), -1)
						util.RequestSendTextCustomerServiceMessage(lastJoinActivity.Openid, content, authInfo)
					}

					//更新推送次数
					lastJoinActivity.PushCount += 1
					session.Table("join_activity").Where("id=?", lastJoinActivity.Id).Cols("push_count").Update(&lastJoinActivity)
				}

				//更新裂变人数
				err = updateLastJoinActivityAllNextSubscribeNum(&lastJoinActivity, 1, session)
				if err != nil {
					session.Rollback()
					return err
				}
				err = updateLastJoinActivityDirectNextSubscribeNum(&lastJoinActivity, 1, subscriber, activity, authInfo, session)
				if err != nil {
					session.Rollback()
					return err
				}
			}
		} else {
			//扫渠道码，加入自己的fissionMember
			var fissionMember models.FissionMember
			fissionMember.Openid = subscriber.Openid
			fissionMember.LastOpenid = ""
			fissionMember.ActivityId = activity.Id
			if channel, hasChannel := paramMap["channel"]; hasChannel {
				fissionMember.Channel = channel
			}
			fissionMember.Status = 1
			session.Table("fission_member").InsertOne(&fissionMember)
		}
	}

	//活动统计
	if hasActivity {
		channel, _ := paramMap["channel"]
		err = updateActivityChannelStatistics(activity, channel, 1, 0, -1, 0, 0, 0, session)
		if err != nil {
			session.Rollback()
			util.Logger.Info("err:" + err.Error())
			return err
		}
	} else {
		//用户自然关注判断为，所有（多个）正在进行活动的自然用户的关注用户。（关注时如果有2个活动正在进行中，则两个活动统计数据内都可以看到该自然用户），无活动不记录
		var activityList []models.Activity
		session.Table("activity").Where("auth_info_id=?", authInfo.Id).And("status=1").Find(&activityList)
		for _, activity := range activityList {
			//自然用户关注，成为当时正在进行的活动中的关注者
			var fissionMember models.FissionMember
			fissionMember.Openid = subscriber.Openid
			fissionMember.LastOpenid = ""
			fissionMember.Channel = ""
			fissionMember.ActivityId = activity.Id
			fissionMember.Status = 1
			session.Table("fission_member").InsertOne(&fissionMember)
			err = updateActivityChannelStatistics(activity, "", 1, 0, -1, 0, 0, 0, session)
			if err != nil {
				session.Rollback()
				util.Logger.Info("err:" + err.Error())
				return err
			}
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		util.Logger.Info("err:" + err.Error())
		return err
	}
	util.Logger.Info("becomeFissionMember end")
	return nil
}

//取消成为裂变人员
func cancelFissionMember(subscriber models.Subscriber, authInfo models.AuthInfo) (err error){
	util.Logger.Info("cancelFissionMember start")
	session := base.DBEngine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		util.Logger.Info("err:" + err.Error())
		return err
	}

	var fissionMemberList []models.FissionMember
	session.Table("fission_member").Where("openid=?", subscriber.Openid).And("status=1").Find(&fissionMemberList)
	for _, fissionMember := range fissionMemberList {
		var activity models.Activity
		hasActivity, _ := session.Table("activity").Where("id=?", fissionMember.ActivityId).And("status=1").Get(&activity)
		if hasActivity {
			util.Logger.Info("hasActivity  ")
			var lastJoinActivity models.JoinActivity
			session.Table("join_activity").Where("openid=?", fissionMember.LastOpenid).And("activity_id=?", activity.Id).And("status=1").Get(&lastJoinActivity)

			//减掉订阅裂变人数
			err = updateLastJoinActivityDirectNextSubscribeNum(&lastJoinActivity, -1, subscriber, activity, authInfo, session)

			//减掉订阅裂变人数
			err = updateLastJoinActivityAllNextSubscribeNum(&lastJoinActivity, -1, session)

			//增加取消订阅裂变人数
			err = updateLastJoinActivityDirectNextUnsubscribeNum(&lastJoinActivity, 1, session)

			//增加取消订阅裂变人数
			err = updateLastJoinActivityAllNextUnsubscribeNum(&lastJoinActivity, 1, session)

			if err != nil {
				session.Rollback()
				util.Logger.Info("err:" + err.Error())
				return err
			}

			//用户取消关注则助力失效（开启此功能后，当助力的好友取消关注，则推送下面的文案给上级用户）
			if fissionMember.LastOpenid != "" {
				var fissionInfo06 models.FissionInfo
				hasFissionInfo06, _ := session.Table("fission_info").Where("activity_id=?", activity.Id).And("type=6").And("status=1").Get(&fissionInfo06)
				if hasFissionInfo06 {
					content := fissionInfo06.Text
					content = strings.Replace(content, "[助力好友昵称]", subscriber.Nickname, -1)
					content = strings.Replace(content, "[已有助力好友人数]", strconv.Itoa(lastJoinActivity.DirectNextSubscribeNum), -1)
					util.RequestSendTextCustomerServiceMessage(fissionMember.LastOpenid, content, authInfo)
				}
			}
		} else {
			util.Logger.Info("!hasActivity  ")
		}

		fissionMember.Status = 2
		fissionMember.UnsubscribeTime = util.UnixOfBeijingTime()
		session.Table("fission_member").Where("id=?", fissionMember.Id).Cols("status", "unsubscribe_time").Update(&fissionMember)
	}


	//清空直接下级的FissionMember
	updateDirectNextFissionMemberSql := "update fission_member set last_openid='' where last_openid='"+subscriber.Openid+"' and status=1 "
	_, err = session.Exec(updateDirectNextFissionMemberSql)
	if err != nil {
		session.Rollback()
		util.Logger.Info("err:" + err.Error())
		return err
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		util.Logger.Info("err:" + err.Error())
		return err
	}
	util.Logger.Info("cancelFissionMember end")
	return nil
}



//加入活动，返回裂变海报mediaId
func joinActivity(subscriber models.Subscriber, activity models.Activity, authInfo models.AuthInfo)(mediaId string, err error) {
	util.Logger.Info("joinActivity start")
	session := base.DBEngine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		util.Logger.Info("err:" + err.Error())
		return "", err
	}

	var joinActivity models.JoinActivity
	hasJoined, _ := session.Table("join_activity").Where("openid=?", subscriber.Openid).And("activity_id=?", activity.Id).And("status=1").Get(&joinActivity)
	if !hasJoined {
		joinActivity.Openid = subscriber.Openid
		joinActivity.ActivityId = activity.Id
		joinActivity.Status = 1

		var fissionMember models.FissionMember
		hasFissionMember, _ := session.Table("fission_member").Where("openid=?", subscriber.Openid).And("activity_id=?", activity.Id).And("last_openid is not null and last_openid != ''").And("status=1").Get(&fissionMember)

		var lastJoinActivity models.JoinActivity
		if hasFissionMember {
			util.Logger.Info("hasFissionMember")
			//如果有上级openid，说明发展下线
			hasLastJoinActivity, _ := session.Table("join_activity").Where("openid=?", fissionMember.LastOpenid).And("activity_id=?", activity.Id).And("status=1").Get(&lastJoinActivity)
			if hasLastJoinActivity {
				util.Logger.Info("hasLastJoinActivity")
				joinActivity.LastOpenid = lastJoinActivity.Openid
				joinActivity.Level = lastJoinActivity.Level + 1
			} else {
				util.Logger.Info("!hasLastJoinActivity")
				util.Logger.Info("二维码有openid却为参加活动 openid:"+fissionMember.LastOpenid+" activityId:"+strconv.FormatInt(activity.Id, 10))
				joinActivity.LastOpenid = ""
				joinActivity.Level = 1
			}
			if fissionMember.Channel != "" {
				util.Logger.Info("fissionMember.Channel != nil")
				//扫码有渠道进来则渠道为二维码上的渠道
				joinActivity.Channel = fissionMember.Channel
			} else {
				util.Logger.Info("fissionMember.Channel == nil")
				//如果有上级openid，则渠道取自上级加入活动的渠道
				joinActivity.Channel = lastJoinActivity.Channel
			}
		} else {
			util.Logger.Info("!hasFissionMember")
			//没有则自己成为领头人
			joinActivity.LastOpenid = ""
			joinActivity.Level = 1
			//渠道取自订阅时的渠道
			joinActivity.Channel = subscriber.Channel
		}

		_, err = session.Table("join_activity").InsertOne(&joinActivity)
		if err != nil {
			session.Rollback()
			util.Logger.Info("err:" + err.Error())
			return "", err
		}

		//活动统计
		err = updateActivityChannelStatistics(activity, joinActivity.Channel, 0, 0, int64(joinActivity.Level), 1, 0, 0, session)
		if err != nil {
			session.Rollback()
			util.Logger.Info("err:" + err.Error())
			return "", err
		}
	} else {
		//重复报名：已报名用户，再次报名。未完成任务：弹出提示 “您已报名该活动，目前任务完成人数8/10，请加油”   完成任务：弹出完成提醒（后台添加）
		if joinActivity.FinishTime == 0 {
			util.RequestSendTextCustomerServiceMessage(subscriber.Openid, "您已报名该活动，目前任务完成人数"+strconv.Itoa(joinActivity.DirectNextSubscribeNum)+"/"+strconv.Itoa(activity.FissionCount)+"，请加油", authInfo)
		} else {
			//达到裂变人数提醒文案（当完成一级裂变人数助力时，推送此提醒文案或图片）
			var fissionInfo05 models.FissionInfo
			hasFissionInfo05, _ := session.Table("fission_info").Where("activity_id=?", activity.Id).And("type=5").And("status=1").Get(&fissionInfo05)
			if hasFissionInfo05 {
				var lastSubscriber models.Subscriber
				session.Table("subscriber").Where("open_id=?", joinActivity.Openid).And("auth_info_id=?", authInfo.Id).Get(&lastSubscriber)

				var nextSubscriber models.Subscriber
				_, err := session.Table("subscriber").Select("subscriber.*").Join("LEFT OUTER", "fission_member", "fission_member.openid=subscriber.openid").Where("fission_member.status=1 and fission_member.activity_id=? and fission_member.last_openid=?", activity.Id, subscriber.Openid).Desc("fission_member.created").Get(&nextSubscriber)
				if err != nil {
					util.Logger.Info("err:"+err.Error())
				}

				content := fissionInfo05.Text
				content = strings.Replace(content, "[用户昵称]", subscriber.Nickname, -1)
				content = strings.Replace(content, "[助力好友昵称]", nextSubscriber.Nickname, -1)
				content = strings.Replace(content, "[奖品剩余数量]", strconv.Itoa(activity.LeftPrizeCount), -1)
				content = strings.Replace(content, "[已有助力好友人数]", strconv.Itoa(joinActivity.DirectNextSubscribeNum), -1)
				content = strings.Replace(content, "[达成条件需要人数]", strconv.Itoa(activity.FissionCount), -1)
				if beego.BConfig.RunMode == "prod" {
					content += "\n<a href='http://zhangmai.vipask.net/userCenter/inputAddress.html?activityId="+strconv.FormatInt(activity.Id, 10)+"&openid="+subscriber.Openid+"&authInfoId="+strconv.FormatInt(authInfo.Id, 10)+"'>点击详情免费领取</a>"
				} else {
					content += "\n<a href='http://tzhangmai.vipask.net/userCenter/inputAddress.html?activityId="+strconv.FormatInt(activity.Id, 10)+"&openid="+subscriber.Openid+"&authInfoId="+strconv.FormatInt(authInfo.Id, 10)+"'>点击详情免费领取</a>"
				}
				util.RequestSendTextCustomerServiceMessage(joinActivity.Openid, content, authInfo)
				if fissionInfo05.Picture != "" {
					//上传至微信
					mediaId, err := util.UploadOnlineMedia(fissionInfo05.Picture, authInfo, 1)
					if err != nil {
						util.Logger.Info("err:"+err.Error())
					} else {
						util.RequestSendMediaCustomerServiceMessage(joinActivity.Openid, mediaId, authInfo)
					}
				}
			}
		}

	}

	//生成裂变海报
	mediaId, err = util.GenerateBanner(subscriber, joinActivity.Channel, activity, authInfo)
	if err != nil {
		session.Rollback()
		util.Logger.Info("err:" + err.Error())
		return "", err
	}
	joinActivity.BannerMediaId = mediaId
	session.Table("join_activity").Where("id=?", joinActivity.Id).Cols("banner_media_id").Update(&joinActivity)

	err = session.Commit()
	if err != nil {
		session.Rollback()
		util.Logger.Info("err:" + err.Error())
		return "", err
	}
	util.Logger.Info("joinActivity end")
	return mediaId, nil
}

//用户取消关注离开活动
func leaveActivity(subscriber models.Subscriber, authInfo models.AuthInfo)(err error) {
	util.Logger.Info("leaveActivity start")
	session := base.DBEngine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		util.Logger.Info("err:" + err.Error())
		return err
	}

	var joinActivityList []models.JoinActivity
	session.Table("join_activity").Where("openid=?", subscriber.Openid).And("status=1").Find(&joinActivityList)
	for _, joinActivity := range joinActivityList {
		var activity models.Activity
		hasActivity, _ := session.Table("activity").Where("id=?", joinActivity.ActivityId).And("status=1").Get(&activity)
		if hasActivity {
			joinActivity.Status = 2
			session.Table("join_activity").Where("id=?", joinActivity.Id).Cols("status").Update(&joinActivity)
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		util.Logger.Info("err:" + err.Error())
		return err
	}
	util.Logger.Info("leaveActivity end")
	return nil
}


//更新上级 所有下线订阅人数
//lastJoinActivity参数表示上级人员
func updateLastJoinActivityAllNextSubscribeNum(lastJoinActivity *models.JoinActivity, number int, session *xorm.Session) (err error){
	if lastJoinActivity.Id == 0 {
		return nil
	}
	util.Logger.Info("updateLastJoinActivityAllNextSubscribeNum 所有下线订阅人数 "+strconv.Itoa(number))
	lastJoinActivity.AllNextSubscribeNum += number
	if lastJoinActivity.AllNextSubscribeNum < 0 {
		lastJoinActivity.AllNextSubscribeNum = 0
	}

	//行锁避免高并发
	lockSql := "update join_activity set id=id where id=?"
	session.Exec(lockSql, lastJoinActivity.Id)

	_, err = session.Table("join_activity").Where("id=?", lastJoinActivity.Id).Cols("all_next_subscribe_num").And("status=1").Update(lastJoinActivity)

	if err != nil {
		util.Logger.Info("updateLastJoinActivityAllNextSubscribeNum err:"+err.Error())
		return err
	}

	if lastJoinActivity.LastOpenid != "" {
		var upperJoinActivity models.JoinActivity
		hasUpper, _ := session.Table("join_activity").Where("openid=?", lastJoinActivity.LastOpenid).And("activity_id=?", lastJoinActivity.ActivityId).And("status=1").Get(&upperJoinActivity)
		if hasUpper {
			updateLastJoinActivityAllNextSubscribeNum(&upperJoinActivity, number, session)
		}
	}

	return nil
}

//更新上级 所有下线取消订阅人数
//lastJoinActivity参数表示上级人员
func updateLastJoinActivityAllNextUnsubscribeNum(lastJoinActivity *models.JoinActivity, number int, session *xorm.Session) (err error){
	if lastJoinActivity.Id == 0 {
		return nil
	}
	util.Logger.Info("updateLastJoinActivityAllNextUnsubscribeNum 所有下线取消订阅人数 "+strconv.Itoa(number))
	lastJoinActivity.AllNextUnsubscribeNum += number
	if lastJoinActivity.AllNextUnsubscribeNum < 0 {
		lastJoinActivity.AllNextUnsubscribeNum = 0
	}

	//行锁避免高并发
	lockSql := "update join_activity set id=id where id=?"
	session.Exec(lockSql, lastJoinActivity.Id)

	_, err = session.Table("join_activity").Where("id=?", lastJoinActivity.Id).Cols("all_next_unsubscribe_num").And("status=1").Update(lastJoinActivity)

	if err != nil {
		util.Logger.Info("updateLastJoinActivityAllNextUnsubscribeNum err:"+err.Error())
		return err
	}

	if lastJoinActivity.LastOpenid != "" {
		var upperJoinActivity models.JoinActivity
		hasUpper, _ := session.Table("join_activity").Where("openid=?", lastJoinActivity.LastOpenid).And("activity_id=?", lastJoinActivity.ActivityId).And("status=1").Get(&upperJoinActivity)
		if hasUpper {
			updateLastJoinActivityAllNextUnsubscribeNum(&upperJoinActivity, number, session)
		}
	}

	return nil
}

//更新上级 直接下线订阅人数
//lastJoinActivity参数表示上级人员
func updateLastJoinActivityDirectNextSubscribeNum(lastJoinActivity *models.JoinActivity, number int, subscriber models.Subscriber, activity models.Activity, authInfo models.AuthInfo, session *xorm.Session) (err error){
	if lastJoinActivity.Id == 0 {
		return nil
	}
	util.Logger.Info("updateLastJoinActivityDirectNextSubscribeNum 直接下线订阅人数 "+strconv.Itoa(number))
	lastJoinActivity.DirectNextSubscribeNum += number
	if lastJoinActivity.DirectNextSubscribeNum < 0 {
		lastJoinActivity.DirectNextSubscribeNum = 0
	}

	//如果是新增，且达到裂变人数
	if number > 0 && lastJoinActivity.DirectNextSubscribeNum >= activity.FissionCount && lastJoinActivity.FinishTime == 0 {
		//达到裂变人数提醒文案（当完成一级裂变人数助力时，推送此提醒文案或图片）
		var fissionInfo05 models.FissionInfo
		hasFissionInfo05, _ := session.Table("fission_info").Where("activity_id=?", activity.Id).And("type=5").And("status=1").Get(&fissionInfo05)
		if hasFissionInfo05 {
			util.Logger.Info("hasFissionInfo05")
			util.Logger.Info(lastJoinActivity.Openid)
			util.Logger.Info(authInfo.Id)

			var lastSubscriber models.Subscriber
			session.Table("subscriber").Where("open_id=?", lastJoinActivity.Openid).And("auth_info_id=?", authInfo.Id).Get(&lastSubscriber)

			content := fissionInfo05.Text
			content = strings.Replace(content, "[用户昵称]", lastSubscriber.Nickname, -1)
			content = strings.Replace(content, "[助力好友昵称]", subscriber.Nickname, -1)
			content = strings.Replace(content, "[奖品剩余数量]", strconv.Itoa(activity.LeftPrizeCount), -1)
			content = strings.Replace(content, "[已有助力好友人数]", strconv.Itoa(lastJoinActivity.DirectNextSubscribeNum), -1)
			content = strings.Replace(content, "[达成条件需要人数]", strconv.Itoa(activity.FissionCount - lastJoinActivity.DirectNextSubscribeNum), -1)
			if beego.BConfig.RunMode == "prod" {
				content += "\n<a href='http://zhangmai.vipask.net/userCenter/inputAddress.html?activityId="+strconv.FormatInt(activity.Id, 10)+"&openid="+lastJoinActivity.Openid+"&authInfoId="+strconv.FormatInt(authInfo.Id, 10)+"'>点击详情免费领取</a>"
			} else {
				content += "\n<a href='http://tzhangmai.vipask.net/userCenter/inputAddress.html?activityId="+strconv.FormatInt(activity.Id, 10)+"&openid="+lastJoinActivity.Openid+"&authInfoId="+strconv.FormatInt(authInfo.Id, 10)+"'>点击详情免费领取</a>"
			}
			util.Logger.Info("next content")
			util.Logger.Info(content)
			util.RequestSendTextCustomerServiceMessage(lastJoinActivity.Openid, content, authInfo)

			if fissionInfo05.Picture != "" {
				//上传至微信
				mediaId, err := util.UploadOnlineMedia(fissionInfo05.Picture, authInfo, 1)
				if err != nil {
					util.Logger.Info("err:"+err.Error())
				} else {
					util.RequestSendMediaCustomerServiceMessage(lastJoinActivity.Openid, mediaId, authInfo)
				}
			}
		}

		//完成任务更新完成时间
		lastJoinActivity.FinishTime = util.UnixOfBeijingTime()

		//活动统计
		err = updateActivityChannelStatistics(activity, lastJoinActivity.Channel, 0, 0, -1, 0, 1, 0, session)
		if err != nil {
			session.Rollback()
			util.Logger.Info("err:" + err.Error())
			return err
		}
	}

	//如果是减少，且原来完成过任务，推送文案
	if number < 0 && lastJoinActivity.DirectNextSubscribeNum < activity.FissionCount && lastJoinActivity.FinishTime != 0 {
		lastJoinActivity.FinishTime = 0
	}



	//行锁避免高并发
	lockSql := "update join_activity set id=id where id=?"
	session.Exec(lockSql, lastJoinActivity.Id)

	_, err = session.Table("join_activity").Where("id=?", lastJoinActivity.Id).Cols("direct_next_subscribe_num", "finish_time").And("status=1").Update(lastJoinActivity)

	if err != nil {
		util.Logger.Info("updateLastJoinActivityDirectNextSubscribeNum err:"+err.Error())
		return err
	}

	return nil
}

//更新上级 直接下线订阅人数
//lastJoinActivity参数表示上级人员
func updateLastJoinActivityDirectNextUnsubscribeNum(lastJoinActivity *models.JoinActivity, number int, session *xorm.Session) (err error){
	if lastJoinActivity.Id == 0 {
		return nil
	}
	util.Logger.Info("updateLastJoinActivityDirectNextUnsubscribeNum 直接下线订阅人数 "+strconv.Itoa(number))
	lastJoinActivity.DirectNextUnsubscribeNum += number
	if lastJoinActivity.DirectNextUnsubscribeNum < 0 {
		lastJoinActivity.DirectNextUnsubscribeNum = 0
	}

	//行锁避免高并发
	lockSql := "update join_activity set id=id where id=?"
	session.Exec(lockSql, lastJoinActivity.Id)

	_, err = session.Table("join_activity").Where("id=?", lastJoinActivity.Id).Cols("direct_next_unsubscribe_num").And("status=1").Update(lastJoinActivity)

	if err != nil {
		util.Logger.Info("updateLastJoinActivityDirectNextUnsubscribeNum err:"+err.Error())
		return err
	}

	return nil
}

//活动统计
func updateActivityStatistics(activity models.Activity, totalSubscribeNum int64, totalUnsubscribeNum int64, maxLevel int64, totalJoinNum int64, totalFinishNum int64, totalReceivePrizeNum int64, session *xorm.Session) (err error){
	var activityStatistics models.ActivityStatistics
	hasActivityStatistics, _ := session.Table("activity_statistics").Where("activity_id=?", activity.Id).Get(&activityStatistics)
	if !hasActivityStatistics {
		activityStatistics.ActivityId = activity.Id
		session.Table("activity_statistics").InsertOne(&activityStatistics)
	}

	activityStatistics.TotalSubscribeNum += totalSubscribeNum
	activityStatistics.TotalUnsubscribeNum += totalUnsubscribeNum
	if maxLevel > 0 {
		activityStatistics.MaxLevel = maxLevel
	}
	activityStatistics.TotalJoinNum += totalJoinNum
	activityStatistics.TotalFinishNum += totalFinishNum
	activityStatistics.TotalReceivePrizeNum += totalReceivePrizeNum
	if activityStatistics.TotalSubscribeNum < 0 {
		activityStatistics.TotalSubscribeNum = 0
	}
	if activityStatistics.TotalUnsubscribeNum < 0 {
		activityStatistics.TotalUnsubscribeNum = 0
	}
	if activityStatistics.MaxLevel < 0 {
		activityStatistics.MaxLevel = 0
	}
	if activityStatistics.TotalJoinNum < 0 {
		activityStatistics.TotalJoinNum = 0
	}
	if activityStatistics.TotalFinishNum < 0 {
		activityStatistics.TotalFinishNum = 0
	}
	if activityStatistics.TotalReceivePrizeNum < 0 {
		activityStatistics.TotalReceivePrizeNum = 0
	}

	_, err = session.Table("activity_statistics").Where("id=?", activityStatistics.Id).AllCols().Update(&activityStatistics)
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		return err
	}

	return nil
}

//活动渠道统计
func updateActivityChannelStatistics(activity models.Activity, channel string, totalSubscribeNum int64, totalUnsubscribeNum int64, maxLevel int64, totalJoinNum int64, totalFinishNum int64, totalReceivePrizeNum int64, session *xorm.Session) (err error){
	//先统计总
	err = updateActivityStatistics(activity, totalSubscribeNum, totalUnsubscribeNum, maxLevel, totalJoinNum, totalFinishNum, totalReceivePrizeNum, session)
	if err != nil {
		return err
	}

	var activityChannelStatistics models.ActivityChannelStatistics
	hasChannelActivityStatistics, _ := session.Table("activity_channel_statistics").Where("activity_id=?", activity.Id).And("channel=?", channel).Get(&activityChannelStatistics)
	if !hasChannelActivityStatistics {
		activityChannelStatistics.ActivityId = activity.Id
		activityChannelStatistics.Channel = channel
		session.Table("activity_channel_statistics").InsertOne(&activityChannelStatistics)
	}

	activityChannelStatistics.TotalSubscribeNum += totalSubscribeNum
	activityChannelStatistics.TotalUnsubscribeNum += totalUnsubscribeNum
	if maxLevel > 0 {
		activityChannelStatistics.MaxLevel = maxLevel
	}
	activityChannelStatistics.TotalJoinNum += totalJoinNum
	activityChannelStatistics.TotalFinishNum += totalFinishNum
	activityChannelStatistics.TotalReceivePrizeNum += totalReceivePrizeNum
	if activityChannelStatistics.TotalSubscribeNum < 0 {
		activityChannelStatistics.TotalSubscribeNum = 0
	}
	if activityChannelStatistics.TotalUnsubscribeNum < 0 {
		activityChannelStatistics.TotalUnsubscribeNum = 0
	}
	if activityChannelStatistics.MaxLevel < 0 {
		activityChannelStatistics.MaxLevel = 0
	}
	if activityChannelStatistics.TotalJoinNum < 0 {
		activityChannelStatistics.TotalJoinNum = 0
	}
	if activityChannelStatistics.TotalFinishNum < 0 {
		activityChannelStatistics.TotalFinishNum = 0
	}
	if activityChannelStatistics.TotalReceivePrizeNum < 0 {
		activityChannelStatistics.TotalReceivePrizeNum = 0
	}

	_, err = session.Table("activity_channel_statistics").Where("id=?", activityChannelStatistics.Id).AllCols().Update(&activityChannelStatistics)
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		return err
	}

	return nil
}

//检查关注用户是新用户还是老用户
//新用户：活动期内：自然关注的、海报扫码的（助力）、渠道扫码关注进入的用户。
//老用户：活动期外自然关注、海报、渠道扫码进入的用户。
func checkSubscriberIsNew(authInfo models.AuthInfo)(isNew int){
	hasActivity, _ := base.DBEngine.Table("activity").Where("auth_info_id=?", authInfo.Id).And("status=1").Get(new(models.Activity))
	if hasActivity {
		return 1
	} else {
		return 2
	}
}

//用户取关活动统计
func updateActivityStatisticsWhenUnsubscribe(subscriber models.Subscriber, authInfo models.AuthInfo) (err error){
	util.Logger.Info("updateActivityStatisticsWhenUnsubscribe start")
	session := base.DBEngine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		util.Logger.Info("err:" + err.Error())
		return err
	}

	//判断当前进行中的活动是否关注
	sql := "select fission_member.* from fission_member left join activity on fission_member.activity_id=activity.id " +
		" where activity.deleted_at is null and activity.status=1 and activity.auth_info_id='"+strconv.FormatInt(authInfo.Id, 10)+"' " +
			" and fission_member.status=1 and fission_member.deleted_at is null and fission_member.openid='"+subscriber.Openid+"' "
	var fissionMemberList []models.FissionMember
	err = session.SQL(sql).Find(&fissionMemberList)
	if err != nil {
		session.Rollback()
		util.Logger.Info("err:" + err.Error())
		return err
	}
	if fissionMemberList != nil && len(fissionMemberList) > 0 {
		for _, fissionMember := range fissionMemberList {
			var activity models.Activity
			session.Table("activity").Where("id=?", fissionMember.ActivityId).Get(&activity)
			err = updateActivityChannelStatistics(activity, fissionMember.Channel, 0, 1, -1, 0, 0, 0, session)
			if err != nil {
				session.Rollback()
				util.Logger.Info("err:" + err.Error())
				return err
			}
		}
	} else {
		activitySql := "select * from activity where activity.deleted_at is null and activity.status=1 and activity.auth_info_id='"+strconv.FormatInt(authInfo.Id, 10)+"' " +
			" and not exists(select 1 from fission_member where fission_member.status=1 and fission_member.deleted_at is null and fission_member.openid='"+subscriber.Openid+"' and fission_member.activity_id=activity.id ) "
		var activityList []models.Activity
		err = session.SQL(activitySql).Find(&activityList)
		if err != nil {
			session.Rollback()
			util.Logger.Info("err:" + err.Error())
			return err
		}
		for _, activity := range activityList {
			err = updateActivityChannelStatistics(activity, "", 0, 1, -1, 0, 0, 0, session)
			if err != nil {
				session.Rollback()
				util.Logger.Info("err:" + err.Error())
				return err
			}
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		util.Logger.Info("err:" + err.Error())
		return err
	}
	util.Logger.Info("updateActivityStatisticsWhenUnsubscribe end")
	return nil
}

//创建author信息
func createAuthorByOpenid(loginUuid string, subscriber *models.Subscriber, paramMap map[string]string){
	session := base.DBEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		session.Rollback()
		util.Logger.Info("createAuthorByOpenid session err = ", err.Error())
		return
	}

	var author models.Author
	hasAuthor, _ := session.Table("author").Where("openid = ?", subscriber.Openid).Get(&author)
	author.Openid = subscriber.Openid
	author.Unionid = subscriber.Unionid
	author.Nickname = subscriber.Nickname
	author.Sex = subscriber.Sex
	author.City = subscriber.City
	author.Province = subscriber.Province
	author.Country = subscriber.Country
	author.Headimgurl = subscriber.Headimgurl
	author.LoginUuid = loginUuid
	if !hasAuthor {
		session.Table("author").InsertOne(&author)
	} else {
		session.Table("author").Where("id = ?", author.Id).AllCols().Update(&author)
	}

	var authorAccount models.AuthorAccount
	hasAuthorAccount, _ := session.Table("author_account").Where("author_id = ?", author.Id).Get(&authorAccount)
	if !hasAuthorAccount {
		authorAccount.AuthorId = author.Id
		authorAccount.Amount = 0
		authorAccount.SettlementType = 0
		authorAccount.UnsettledAmount = 0
		session.Table("author_account").InsertOne(&authorAccount)
	}

	//创建观看历史
	videoIdStr, hasVideoId := paramMap["videoId"]
	watchTimeStr, hasWatchTime := paramMap["watchTime"]
	if hasVideoId && hasWatchTime {
		videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)
		watchTime, _ := strconv.Atoi(watchTimeStr)
		err := addVideoViewHistory(author.Id, videoId, watchTime, session)
		if err != nil {
			util.Logger.Info("addVideoViewHistory err" + err.Error())
			session.Rollback()
			return
		}
	}

	err = session.Commit()
	if err != nil {
		session.Rollback()
		util.Logger.Info("createAuthorByOpenid commit err = ", err.Error())
		return
	}
}

//关注或扫码，参数有观看记录信息，则推送客服消息告知
func pushVideoViewHistoryMessage(openid string, authInfo models.AuthInfo, paramMap map[string]string){
	videoIdStr, hasVideoId := paramMap["videoId"]
	_, hasWatchTime := paramMap["watchTime"]
	var video  models.Video
	var author models.Author
	var chapter models.Course
	var course models.Course
	base.DBEngine.Table("video").Where("id = ?", videoIdStr).Get(&video)
	base.DBEngine.Table("author").Where("openid = ?", openid).Get(&author)
	base.DBEngine.Table("course").Where("id = ?", video.CourseId).Get(&chapter)
	base.DBEngine.Table("course").Where("id = ?", chapter.CheckStatus).Get(&course)
	if hasVideoId && hasWatchTime {
		//查询该用户是否转发过该视频
		var videoShare models.VideoShare
		var content string
		hasVideoShare, _ := base.DBEngine.Table("video_share").Where("author_id=?", author.Id).And("video_id=?", video.Id).And("status = ? ", 1).Get(&videoShare)
		if hasVideoShare {
			//发送客服消息
			//content = "感谢关注! \n您之前观看了<a href= '"+ base.ServerURL+"/video/wx/detail.html?videoId=" + videoIdStr + "&watchTime=" + watchTimeStr + "&authorId=" + strconv.FormatInt(author.Id,10) +"'>《"+ video.Title +"》</a>" + " \n \n可在【个人中心】查看佣金等。"
			content = "感谢关注! \n您之前观看了<a href= '"+ base.ServerURL+"/wx?authorId=" + strconv.FormatInt(author.Id,10) + "/#/detail/" + strconv.FormatInt(course.Id ,10) +"'>《"+ video.Title +"》</a>" + " \n \n可在【个人中心】查看佣金等。"

		} else {
			//发送客服消息
			//content = "感谢关注! \n您之前观看了<a href= '"+ base.ServerURL+"/video/wx/detail.html?videoId=" + videoIdStr + "&watchTime=" + watchTimeStr + "&authorId=" + strconv.FormatInt(author.Id,10) +"'>《"+ video.Title +"》</a>" + " \n点击可以继续观看下一段。\n \n可在【个人中心】查看佣金等。"
			content = "感谢关注! \n您之前观看了<a href= '"+ base.ServerURL+"/wx?authorId=" + strconv.FormatInt(author.Id,10) + "/#/detail/" + strconv.FormatInt(course.Id ,10) +"'>《"+ video.Title +"》</a>" + " \n点击可以继续观看下一段。\n \n可在【个人中心】查看佣金等。"

		}
		util.RequestSendTextCustomerServiceMessage(openid, content, authInfo)
	}
}






















//---------------------------微信接口请求---------------------------------------------------
//根据openid获取用户信息
func requestUserInfoByOpenId(openid string, authInfo models.AuthInfo) (userInfo models.UserInfoJsonBody, err error){
	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + authInfo.AuthAccessToken + "&openid=" + openid)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)

	resp, err := client.Do(r)
	if err != nil {
		util.Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return userInfo, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return userInfo, err
	}

	if resp.StatusCode != 200 {
		util.Logger.Info("requestUserInfoByOpenId err :resp.StatusCode != 200")
		return userInfo, err
	}

	response := models.UserInfoJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("requestUserInfoByOpenId json.Unmarshal(body, &response) err :" + err.Error())
		return userInfo, err
	}

	return response, nil
}

func addVideoViewHistory(authorId int64, videoId int64, watchTime int, session *xorm.Session) (err error){
	var videoViewHistory models.VideoViewHistory
	sql := " select * from video_view_history where video_id=" +strconv.FormatInt(videoId,10)+" and author_id =" + strconv.FormatInt(authorId,10) +" and deleted_at is null"
	hasVideoViewHistory, _ := session.SQL(sql).Get(&videoViewHistory)
	if hasVideoViewHistory {
		videoViewHistory.WatchTime = watchTime
		_, err := session.Table("video_view_history").Where("id = ?", videoViewHistory.Id).AllCols().Update(&videoViewHistory)
		if err != nil {
			util.Logger.Info("AddBrowseCount update videoViewHistory err = ", err.Error())
			return err
		}
	} else {
		videoViewHistory.VideoId = videoId
		videoViewHistory.AuthorId = authorId
		videoViewHistory.WatchTime = watchTime
		_, err = session.Table("video_view_history").InsertOne(&videoViewHistory)
		if err != nil {
			util.Logger.Info("AddBrowseCount insert videoViewHistory err = ", err.Error())
			return err
		}
	}
	return nil
}




