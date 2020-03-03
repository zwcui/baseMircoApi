/*
@Time : 2019/10/11 下午5:18 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"jingting_server/messageservice/models"
	"jingting_server/messageservice/util"
	"encoding/base64"
	"encoding/xml"
	"jingting_server/messageservice/base"
	"strconv"
	"strings"
	"net/url"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"jingting_server/messageservice/remote"
	"bytes"
	"time"
	"github.com/pkg/errors"
)

//公众号开发
type GzhController struct {
	apiController
}

func (this *GzhController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
		//{"/createMenu", "post", []int{1}},
	}
	this.userAuth()
}

// @Title 配置公众号回调地址（只用一次）
// @Description 配置公众号回调地址（只用一次）
// @Param	signature			query		string	  		false		"signature"
// @Param	timestamp			query		string	  		false		"timestamp"
// @Param	nonce				query		string	  		false		"nonce"
// @Param	echostr				query		string	  		false		"echostr"
// @Success 200 {string} success
// @router /receiveMessage [get]
func(this *GzhController) ReceiveMessage(){
	signature := this.GetString("signature")
	timestamp := this.GetString("timestamp")
	nonce := this.GetString("nonce")
	echostr := this.GetString("echostr")

	//配置时校验
	if signature != "" && timestamp != "" && nonce != "" && echostr != "" {
		util.Logger.Info("signature = " + signature)
		util.Logger.Info("timestamp = " + timestamp)
		util.Logger.Info("nonce = " + nonce)
		util.Logger.Info("echostr = " + echostr)

		//暂不校验签名
		this.IsDirectReturn = 1
		this.ReturnData = echostr
		return
	}
}

// @Title 获取公众号信息（微信回调）
// @Description 获取公众号信息（微信回调）
// @Success 200 {string} success
// @router /receiveMessage [post]
func(this *GzhController) PostReceiveMessage(){

	response := models.ReceiveMessageEncryptXmlBody{}
	//util.Logger.Info("receiveMessage ReceiveMessageEncryptXmlBody body = ", string(this.Ctx.Input.RequestBody))
	err := xml.Unmarshal(this.Ctx.Input.RequestBody, &response)
	if err != nil {
		util.Logger.Info("xml.Unmarshal body = " + err.Error())
		this.ServeXML()
		this.StopRun()
	}

	AESKey, _ := base64.StdEncoding.DecodeString(models.JTGZHEncodingAESKey + "=")

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


	//回复空串让微信服务器不再重试
	this.IsDirectReturn = 1
	this.ReturnData = "success"

	job := func() {
		handleGZHMessage(receiveMessageXmlBody)
	}
	base.GoPool.Submit(job)

}

// @Title 创建公众号菜单
// @Description 创建公众号菜单
// @Param	param				query		string	  		false		"菜单json"
// @Success 200 {string} success
// @router /createMenu [post]
func(this *GzhController) CreateMenu(){
	param := this.GetString("param", "")

	if param == "" {
		param = "{\"button\":[{\"type\":\"view\",\"name\":\"学院\",\"key\":\"study\",\"url\":\"http://jingting.vipask.net/wx/?#/index/recommend\"},{\"type\":\"view\",\"name\":\"发现\",\"key\":\"find\",\"url\":\"http://jingting.vipask.net/wx/?#/emotions/all\"},{\"type\":\"view\",\"name\":\"我的\",\"key\":\"center\",\"url\":\"http://jingting.vipask.net/wx/?#/pcenter\"}]}"
	}

	util.Logger.Info(param)

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=" + getGZHAccessToken())
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(param)))
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		util.Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.GZHError100)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.GZHError100)
		return
	}

	if resp.StatusCode != 200 {
		util.Logger.Info("requestComponentAccessToken err :resp.StatusCode != 200")
		this.ReturnData = util.GenerateAlertMessage(models.GZHError100)
		return
	}

	util.Logger.Info("body:"+string(body))

	response := RequestGZHCreateMenuResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("CreateMenu json.Unmarshal(body, &response) err :" + err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.GZHError100)
		return
	}

	if response.Errcode == 0 {
		this.ReturnData = "success"
	} else {
		util.Logger.Info("create menu failed:" + response.Errmsg)
		this.ReturnData = util.GenerateAlertMessage(models.GZHError100)
		return
	}



	//body示例：
	//{
	//	"button":[
	//	{
	//	"type":"click",
	//	"name":"今日歌曲",
	//	"key":"V1001_TODAY_MUSIC"
	//	},
	//	{
	//	"name":"菜单",
	//	"sub_button":[
	//	{
	//	"type":"view",
	//	"name":"搜索",
	//	"url":"http://www.soso.com/"
	//	},
	//	{
	//	"type":"miniprogram",
	//	"name":"wxa",
	//	"url":"http://mp.weixin.qq.com",
	//	"appid":"wx286b93c14bbf93aa",
	//	"pagepath":"pages/lunar/index"
	//	},
	//	{
	//	"type":"click",
	//	"name":"赞一下我们",
	//	"key":"V1001_GOOD"
	//	}]
	//	}]
	//}
	//{
	//	"button": [
	//{
	//"name": "扫码",
	//"sub_button": [
	//{
	//"type": "scancode_waitmsg",
	//"name": "扫码带提示",
	//"key": "rselfmenu_0_0",
	//"sub_button": [ ]
	//},
	//{
	//"type": "scancode_push",
	//"name": "扫码推事件",
	//"key": "rselfmenu_0_1",
	//"sub_button": [ ]
	//}
	//]
	//},
	//{
	//"name": "发图",
	//"sub_button": [
	//{
	//"type": "pic_sysphoto",
	//"name": "系统拍照发图",
	//"key": "rselfmenu_1_0",
	//"sub_button": [ ]
	//},
	//{
	//"type": "pic_photo_or_album",
	//"name": "拍照或者相册发图",
	//"key": "rselfmenu_1_1",
	//"sub_button": [ ]
	//},
	//{
	//"type": "pic_weixin",
	//"name": "微信相册发图",
	//"key": "rselfmenu_1_2",
	//"sub_button": [ ]
	//}
	//]
	//},
	//{
	//"name": "发送位置",
	//"type": "location_select",
	//"key": "rselfmenu_2_0"
	//},
	//{
	//"type": "media_id",
	//"name": "图片",
	//"media_id": "MEDIA_ID1"
	//},
	//{
	//"type": "view_limited",
	//"name": "图文消息",
	//"media_id": "MEDIA_ID2"
	//}
	//]
	//}
}

// @Title TestRequestGZHUserInfoByOpenId
// @Description TestRequestGZHUserInfoByOpenId
// @Param	param				query		string	  		false		"openid"
// @Success 200 {string} success
// @router /testRequestGZHUserInfoByOpenId [post]
func(this *GzhController) TestRequestGZHUserInfoByOpenId(){
	param := this.GetString("param", "oVD2dwixwwmzAkHkQMie7CZaz2Pc")

	userInfo, err := requestGZHUserInfoByOpenId(param)
	if err != nil {
		util.Logger.Info("requestGZHUserInfoByOpenId(param) err:" + err.Error())

		this.ReturnData = err.Error()
		return
	}

	result, err := json.Marshal(userInfo)
	if err != nil {
		util.Logger.Info("json.Marshal(userInfo) err:" + err.Error())
		this.ReturnData = err.Error()
		return
	}

	util.Logger.Info(string(result))

	this.ReturnData = userInfo.Unionid
}





























//---------------------------------------方法---------------------------------------------------------


//处理公众号消息
func handleGZHMessage(receiveMessageXmlBody *models.ReceiveMessageXmlBody){
	util.Logger.Info("handleGZHMessage  start")
	util.Logger.Info("running goroutines:" + strconv.Itoa(base.GoPool.Running()))
	//存库
	var message models.Message
	message.Appid = models.JTGZHAppId
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
		//关注公众号
		var welcomeMessageConfig1 models.SystemConfig
		base.DBEngine.Table("system_config").Where("program='gzh_welcome_message1'").Get(&welcomeMessageConfig1)
		var welcomeMessageConfig2 models.SystemConfig
		base.DBEngine.Table("system_config").Where("program='gzh_welcome_message2'").Get(&welcomeMessageConfig2)
		var welcomeMessageConfig3 models.SystemConfig
		base.DBEngine.Table("system_config").Where("program='gzh_welcome_message3'").Get(&welcomeMessageConfig3)

		if welcomeMessageConfig1.ProgramValue != "" {
			util.RequestSendGZHTextCustomerServiceMessage(message.FromUserName, welcomeMessageConfig1.ProgramValue, getGZHAccessToken())
		}

		if welcomeMessageConfig2.ProgramValue != "" {
			util.RequestSendGZHTextCustomerServiceMessage(message.FromUserName, welcomeMessageConfig2.ProgramValue, getGZHAccessToken())
		}
		//是否在免费时间段
		var freeConfig models.SystemConfig
		base.DBEngine.Table("system_config").Where("program='free_course_start_end'").Get(&freeConfig)
		freeStartTimeStr := strings.Split(freeConfig.ProgramValue, ",")[0]
		freeEndTimeStr := strings.Split(freeConfig.ProgramValue, ",")[1]
		freeStartTime, _ := strconv.ParseInt(freeStartTimeStr, 10, 64)
		freeEndTime, _ := strconv.ParseInt(freeEndTimeStr, 10, 64)
		if util.UnixOfBeijingTime() >= freeStartTime && util.UnixOfBeijingTime() <= freeEndTime {
			util.RequestSendGZHTextCustomerServiceMessage(message.FromUserName, welcomeMessageConfig3.ProgramValue, getGZHAccessToken())
		}


		util.Logger.Info("扫码关注*****************************openid:" +  message.FromUserName)
		//加入订阅者记录
		var subscriber models.Subscriber
		hasSubscriber, _ := base.DBEngine.Table("subscriber").Where("openid=?", message.FromUserName).And("gzh_app_id=?", models.JTGZHAppId).And("auth_info_id=0").Get(&subscriber)
		if !hasSubscriber {
			userInfo, err := requestGZHUserInfoByOpenId(message.FromUserName)
			//增加重试
			var flag bool = false
			if err != nil || userInfo.Unionid == "" {
				util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
				for errorCount := 0; errorCount < 10; errorCount++ {
					util.Logger.Info("重试第" + strconv.Itoa(errorCount + 1) + "次")
					userInfo, err = requestGZHUserInfoByOpenId(message.FromUserName)
					if err == nil {
						flag = true
						break
					}
				}
			} else {
				flag = true
			}


			//userInfo, err := getGZHUserInfoByOpenId(message.FromUserName,0)
			if !flag {
				util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
			} else {
				subscriber.Openid = message.FromUserName
				subscriber.GzhAppId = models.JTGZHAppId
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
				base.DBEngine.Table("subscriber").InsertOne(&subscriber)
			}
		} else if subscriber.Subscribe == 0 {
			userInfo, err := requestGZHUserInfoByOpenId(message.FromUserName)
			//增加重试
			var flag bool = false
			if err != nil {
				util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
				for errorCount := 0; errorCount < 10; errorCount++ {
					util.Logger.Info("重试第" + strconv.Itoa(errorCount + 1) + "次")
					userInfo, err = requestGZHUserInfoByOpenId(message.FromUserName)
					if err == nil {
						flag = true
						break
					}
				}
			} else {
				flag = true
			}

			if !flag {
				util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
			} else {
				subscriber.Openid = message.FromUserName
				subscriber.GzhAppId = models.JTGZHAppId
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
				subscriber.Created = util.UnixOfBeijingTime()
				base.DBEngine.Table("subscriber").Where("id=?", subscriber.Id).AllCols().Update(&subscriber)
			}
		}

		//关注公众号预注册
		var author models.Author
		hasAuthor, _ := base.DBEngine.Table("author").Where("unionid = ?", subscriber.Unionid).Get(&author)
		author.Openid = subscriber.Openid
		author.Unionid = subscriber.Unionid
		author.Nickname = subscriber.Nickname
		author.Sex = subscriber.Sex
		author.City = subscriber.City
		author.Province = subscriber.Province
		author.Country = subscriber.Country
		author.Headimgurl = subscriber.Headimgurl
		author.SignInDays += countTodaySignIn(author)
		author.LastSignInTime = util.UnixOfBeijingTime()
		if !hasAuthor {
			author.RegisterSystem = 0	//公众号注册
			author.RegisterChannel = "微信公众号"
			base.DBEngine.Table("author").InsertOne(&author)

		} else {
			base.DBEngine.Table("author").Where("id = ?", author.Id).AllCols().Update(&author)
		}
		var authorAccount models.AuthorAccount
		hasAuthorAccount, _ := base.DBEngine.Table("author_account").Where("author_id = ?", author.Id).Get(&authorAccount)
		if !hasAuthorAccount {
			authorAccount.AuthorId = author.Id
			authorAccount.Amount = 0
			authorAccount.SettlementType = 0
			authorAccount.UnsettledAmount = 0
			base.DBEngine.Table("author_account").InsertOne(&authorAccount)
		}

		//存入redis
		authorRedis := models.AuthorRedis{author.Id, author.Openid, author.WebOpenid, author.AppOpenid, author.Unionid}
		authorRedisBytes, _ := json.Marshal(authorRedis)
		base.RedisCache.Put(REDIS_BASEAUTH_UNIONID + author.Unionid, string(authorRedisBytes), 60*60*2*time.Second)

		//关注或扫码，参数有观看记录信息，则推送客服消息告知
		pushCourseoViewHistoryMessage(message.FromUserName, paramMap)

		//分享名片后关注公众号永久绑定上下级关系
		if lastAuthorIdStr, ok := paramMap["lastAuthorId"]; ok {
			lastAuthorId, _ := strconv.ParseInt(lastAuthorIdStr, 10, 64)
			ok, err := remote.CreateShare(lastAuthorId, author.Id)
			if err != nil {
				util.Logger.Info("CreateShare err = ", err.Error())
			}
			if ok {
				//推送消息给上级
				message := models.AppMessage{}
				message.Content = "【您已成功邀请好友"+ author.Nickname + " 加入辣课，可在个人中心->我的推广中查看下级明细】"
				message.ReceiverId = lastAuthorId
				message.ActionUrl = JumpUrlWithKeyAndPramas(models.JTFOLLOW_JUMP_KEY, nil)
				message.Type = 5
				_, err = PushMessageToUser(lastAuthorId, &message, "", 0)
				if err != nil {
					util.Logger.Info("PushMessageToUser  err = ", err.Error())
				}
				base.DBEngine.Table("app_message").InsertOne(&message)
			}
		}
	}

	//取消订阅
	if message.MsgType == "event" && message.Event == "unsubscribe" {
		var subscriber models.Subscriber
		base.DBEngine.Table("subscriber").Where("openid=?", message.FromUserName).And("gzh_app_id=?", models.JTGZHAppId).Get(&subscriber)

		//取消订阅更新订阅者信息
		subscriber.Subscribe = 0
		subscriber.UnsubscribeTime = util.UnixOfBeijingTime()
		base.DBEngine.Table("subscriber").Where("id=?", subscriber.Id).Cols("subscribe").Update(&subscriber)
	}

	//已关注再扫码
	if message.MsgType == "event" && message.Event == "SCAN" {
		util.Logger.Info("已关注再扫码*****************************openid:" +  message.FromUserName)
		var subscriber models.Subscriber
		hasSubscriber, _ := base.DBEngine.Table("subscriber").Where("openid=?", message.FromUserName).And("gzh_app_id=?", models.JTGZHAppId).Get(&subscriber)
		if !hasSubscriber {
			userInfo, err := requestGZHUserInfoByOpenId(message.FromUserName)
			//增加重试
			var flag bool = false
			if err != nil || userInfo.Unionid == "" {
				util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
				for errorCount := 0; errorCount < 10; errorCount++ {
					util.Logger.Info("重试第" + strconv.Itoa(errorCount) + "次")
					userInfo, err = requestGZHUserInfoByOpenId(message.FromUserName)
					if err == nil {
						flag = true
						break
					}
				}
			} else {
				flag = true
			}


			//userInfo, err := getGZHUserInfoByOpenId(message.FromUserName,0)
			if !flag {
				util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
			} else {
				subscriber.Openid = message.FromUserName
				subscriber.GzhAppId = models.JTGZHAppId
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
				base.DBEngine.Table("subscriber").InsertOne(&subscriber)
			}
		} else if subscriber.Subscribe == 0 {
			userInfo, err := requestGZHUserInfoByOpenId(message.FromUserName)
			//增加重试
			var flag bool = false
			if err != nil || userInfo.Unionid == "" {
				util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
				for errorCount := 0; errorCount < 10; errorCount++ {
					util.Logger.Info("重试第" + strconv.Itoa(errorCount) + "次")
					userInfo, err = requestGZHUserInfoByOpenId(message.FromUserName)
					if err == nil {
						flag = true
						break
					}
				}
			} else {
				flag = true
			}

			if !flag {
				util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
			} else {
				subscriber.Openid = message.FromUserName
				subscriber.GzhAppId = models.JTGZHAppId
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
				subscriber.Created = util.UnixOfBeijingTime()
				base.DBEngine.Table("subscriber").Where("id=?", subscriber.Id).AllCols().Update(&subscriber)
			}
		}

		//关注公众号预注册
		var author models.Author
		hasAuthor, _ := base.DBEngine.Table("author").Where("unionid = ?", subscriber.Unionid).Get(&author)
		author.Openid = subscriber.Openid
		author.Unionid = subscriber.Unionid
		author.Nickname = subscriber.Nickname
		author.Sex = subscriber.Sex
		author.City = subscriber.City
		author.Province = subscriber.Province
		author.Country = subscriber.Country
		author.Headimgurl = subscriber.Headimgurl
		author.SignInDays += countTodaySignIn(author)
		author.LastSignInTime = util.UnixOfBeijingTime()
		if !hasAuthor {
			author.RegisterSystem = 0	//公众号注册
			author.RegisterChannel = "微信公众号"
			base.DBEngine.Table("author").InsertOne(&author)
		} else {
			base.DBEngine.Table("author").Where("id = ?", author.Id).AllCols().Update(&author)
		}
		var authorAccount models.AuthorAccount
		hasAuthorAccount, _ := base.DBEngine.Table("author_account").Where("author_id = ?", author.Id).Get(&authorAccount)
		if !hasAuthorAccount {
			authorAccount.AuthorId = author.Id
			authorAccount.Amount = 0
			authorAccount.SettlementType = 0
			authorAccount.UnsettledAmount = 0
			base.DBEngine.Table("author_account").InsertOne(&authorAccount)
		}

		//存入redis
		authorRedis := models.AuthorRedis{author.Id, author.Openid, author.WebOpenid, author.AppOpenid, author.Unionid}
		authorRedisBytes, _ := json.Marshal(authorRedis)
		base.RedisCache.Put(REDIS_BASEAUTH_UNIONID + author.Unionid, string(authorRedisBytes), 60*60*2*time.Second)

		//分享名片后关注公众号永久绑定上下级关系
		if lastAuthorIdStr, ok := paramMap["lastAuthorId"]; ok {
			lastAuthorId, _ := strconv.ParseInt(lastAuthorIdStr, 10, 64)
			ok, err := remote.CreateShare(lastAuthorId, author.Id)
			if err != nil {
				util.Logger.Info("CreateShare err = ", err.Error())
			}
			if ok {
				//推送消息给上级
				message := models.AppMessage{}
				message.Content = "【您已成功邀请好友"+ author.Nickname + " 加入辣课，可在个人中心->我的推广中查看下级明细】"
				message.ReceiverId = lastAuthorId
				message.ActionUrl = JumpUrlWithKeyAndPramas(models.JTFOLLOW_JUMP_KEY, nil)
				message.Type = 5
				_, err = PushMessageToUser(lastAuthorId, &message, "", 0)
				if err != nil {
					util.Logger.Info("PushMessageToUser  err = ", err.Error())
				}
				base.DBEngine.Table("app_message").InsertOne(&message)
			}
		}
	}

	//公众号菜单点击
	if message.MsgType == "event" && message.Event == "VIEW" && message.EventKey != "" {

	}

	//回复关键词
	if message.MsgType == "text" && message.Content != "" {

	}

	util.Logger.Info("handleGZHMessage  end")
}


//根据openid获取用户信息
func requestGZHUserInfoByOpenId(openid string) (userInfo models.UserInfoJsonBody, err error){
	accessToken := getGZHAccessToken()

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + accessToken + "&openid=" + openid + "&lang=zh_CN")

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
	util.Logger.Info(string(body))

	if strings.Contains(string(body), "\"errcode\":40001") {
		return userInfo, errors.New(string(body))
	}

	response := models.UserInfoJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("requestUserInfoByOpenId json.Unmarshal(body, &response) err :" + err.Error())
		return userInfo, err
	}

	return response, nil
}

//获取辣课公众号access_token
func getGZHAccessToken() string {
	var gzhConfig models.SystemConfig
	hasAccessToken, _ := base.DBEngine.Table("system_config").Where("program='ghz_access_token'").Get(&gzhConfig)
	if !hasAccessToken {
		resp, err := requestGZHAcceccToken()
		if err != nil {
			util.Logger.Info("requestGZHAcceccToken err:" + err.Error())
			return ""
		}
		gzhConfig.Description = "辣课公众号access_token"
		gzhConfig.Program = "ghz_access_token"
		gzhConfig.ProgramValue = resp.AccessToken
		gzhConfig.ProgramExpireTime = util.UnixOfBeijingTime() + resp.ExpiresIn
		base.DBEngine.Table("system_config").InsertOne(&gzhConfig)
		return resp.AccessToken
	} else {
		if gzhConfig.ProgramExpireTime <= util.UnixOfBeijingTime() + 40 * 60 {
			resp, err := requestGZHAcceccToken()
			if err != nil {
				util.Logger.Info("requestGZHAcceccToken err:" + err.Error())
				return ""
			}
			gzhConfig.ProgramValue = resp.AccessToken
			gzhConfig.ProgramExpireTime = util.UnixOfBeijingTime() + resp.ExpiresIn
			base.DBEngine.Table("system_config").Where("r_id=?", gzhConfig.RId).AllCols().Update(&gzhConfig)
			return resp.AccessToken
		} else {
			return gzhConfig.ProgramValue
		}
	}
}

//请求公众号access_token
func requestGZHAcceccToken() (response RequestGZHAcceccTokenResponse, err error) {

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid="+models.JTGZHAppId+"&secret="+models.JTGZHAppSecret)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)

	resp, err := client.Do(r)
	if err != nil {
		util.Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return
	}

	if resp.StatusCode != 200 {
		util.Logger.Info("requestGZHAcceccToken err :resp.StatusCode != 200")
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("requestGZHAcceccToken json.Unmarshal(body, &response) err :" + err.Error())
		return
	}

	return response, nil
}


type RequestGZHAcceccTokenResponse struct {
	AccessToken				string				`description:"获取到的凭证" json:"access_token"`
	ExpiresIn				int64				`description:"凭证有效时间，单位：秒" json:"expires_in"`
}

type RequestGZHCreateMenuResponse struct {
	Errcode					int					`description:"errcode" json:"errcode"`
	Errmsg					string				`description:"errmsg" json:"errmsg"`
}

//关注或扫码，参数有观看记录信息，则推送客服消息告知
func pushCourseoViewHistoryMessage(openid string, paramMap map[string]string){
	videoIdStr, hasVideoId := paramMap["videoId"]
	watchTimeStr, hasWatchTime := paramMap["watchTime"]
	var video  models.Video
	var author models.Author
	var chapter models.Course
	var course models.Course
	base.DBEngine.Table("video").Where("id = ?", videoIdStr).Get(&video)
	base.DBEngine.Table("author").Where("openid = ?", openid).Get(&author)
	base.DBEngine.Table("course").Where("id = ?", video.CourseId).Get(&chapter)
	base.DBEngine.Table("course").Where("id = ?", chapter.ParentId).Get(&course)
	if hasVideoId && hasWatchTime {
		//查询该用户是否转发过该视频
		var videoShare models.VideoShare
		var content string
		hasVideoShare, _ := base.DBEngine.Table("video_share").Where("author_id=?", author.Id).And("video_id=?", video.Id).And("status = ? ", 1).Get(&videoShare)
		if hasVideoShare {
			//发送客服消息
			//content = "感谢关注! \n您之前观看了<a href= '"+ base.ServerURL+"/video/wx/detail.html?videoId=" + videoIdStr + "&watchTime=" + watchTimeStr + "&authorId=" + strconv.FormatInt(author.Id,10) +"'>《"+ video.Title +"》</a>" + " \n \n可在【个人中心】查看佣金等。"
			content = "感谢关注! \n您之前观看了<a href= '"+ base.ServerURL+"/wx/?#/detail/" + strconv.FormatInt(course.Id,10) + "?authorId=" + strconv.FormatInt(author.Id,10) +"&videoId=" + videoIdStr + "&watchTime=" + watchTimeStr +"'>《"+ video.Title +"》</a>" + " \n \n可在【个人中心】查看佣金等。"

		} else {
			//发送客服消息
			//content = "感谢关注! \n您之前观看了<a href= '"+ base.ServerURL+"/video/wx/detail.html?videoId=" + videoIdStr + "&watchTime=" + watchTimeStr + "&authorId=" + strconv.FormatInt(author.Id,10) +"'>《"+ video.Title +"》</a>" + " \n点击可以继续观看下一段。\n \n可在【个人中心】查看佣金等。"
			content = "感谢关注! \n您之前观看了<a href= '"+ base.ServerURL+"/wx/?#/detail/" + strconv.FormatInt(course.Id ,10) + "?authorId=" + strconv.FormatInt(author.Id,10) +"&videoId=" + videoIdStr + "&watchTime=" + watchTimeStr +"'>《"+ video.Title +"》</a>" + " \n点击可以继续观看下一段。\n \n可在【个人中心】查看佣金等。"

		}
		util.RequestSendGZHTextCustomerServiceMessage(openid, content, getGZHAccessToken())
	}
}

var Count, ErrorCount int

func getGZHUserInfoByOpenId(openid string, errorCount int) (userInfo models.UserInfoJsonBody, err error)  {
	response, err := requestGZHUserInfoByOpenId(openid)
	if err != nil {
		errorCount += 1
		util.Logger.Info("错误调用次数 :" + strconv.Itoa(errorCount))
		if errorCount > 10 {
			util.Logger.Info("requestGZHUserInfoByOpenId err:" + err.Error())
			return userInfo, err
		} else {
			util.Logger.Info("重试第" + strconv.Itoa(errorCount) + "次")
			getGZHUserInfoByOpenId(openid, errorCount)
		}
	}

	return response, nil


}

//通过上次登录时间判断是否要累计登录天数
func countTodaySignIn(author models.Author) int {
	lastSignInTime := author.LastSignInTime
	if lastSignInTime == 0 {
		return 1
	}

	last := time.Unix(lastSignInTime, 0).Format("2006-01-02")
	now := time.Now().Format("2006-01-02")
	if last != now {
		return 1
	}

	return 0
}
//创建店铺
func createStore (author models.Author) error{

	//创建店铺
	var systemConfig models.SystemConfig
	var cover string
	hasSystemConfig, _ := base.DBEngine.Table("system_config").Where("program='default_store_cover'").Get(&systemConfig)
	if hasSystemConfig {
		cover = systemConfig.ProgramValue
	} else {
		cover = ""
	}
	store := models.Store{
		AuthorId: author.Id,
		Cover: cover,
	}
	str, _ := json.Marshal(store)
	_, createStoreErr := remote.CreateStore(string(str))
	if createStoreErr != nil {
		util.Logger.Info(" CreateStore err:" + createStoreErr.Error())
		return createStoreErr
	}

	return nil
}

