/*
@Time : 2019/8/6 下午6:05
@Author : lianwu
@File : pushApi.go
@Software: GoLand
*/
package controllers

import (
	"errors"
	"jingting_server/messageservice/models"
	"jingting_server/messageservice/base"
	"jingting_server/messageservice/util"
	"strconv"
	"encoding/json"
	"net/url"
	"time"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"sort"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"github.com/astaxie/beego"
	"jingting_server/messageservice/proto/message"
	"context"
)

type PushController struct {
	apiController
}

func (this *PushController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/push-to-single", "post", []int{}},
	}
	this.authorAuth()
}


// @Title 推送消息给指定用户
// @Description 推送消息给指定用户
// @Param	authorId     formData  int64	  true		"authorId"
// @Param	mContent     formData  string	  true		"message content"
// @Param	actionUrl    formData  string	  false		"message actionUrl"
// @Param	type         formData  int	      true		"message type：1:购买成功消息 2:转发成功消息 3:评论成功消息 4:举报成功消息"
// @Success 200 {object} models.AppMessage
// @router /push-to-single [post]
func (this *PushController) PushToSingle() {
	authorId := this.MustInt64("authorId")
	content := this.MustString("mContent")
	fmt.Println(content)
	actionUrl := this.GetString("actionUrl")
	messageType := this.MustInt("type")


	message := models.AppMessage{}
	message.Content = content
	message.ReceiverId = authorId
	message.ActionUrl = actionUrl
	message.Type = messageType

	_, err := PushMessageToUser(authorId, &message, "", 0)

	if err != nil {
		util.Logger.Info("PushToSingle  err = ", err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.PushMessageError100)
		return
	}

	_, err = base.DBEngine.Table("app_message").InsertOne(&message)
	this.ReturnData = message

}



































//-----------------------------方法------------------------------------

//authorId目标用户
func PushMessageToUser(authorId int64, message *models.AppMessage, sound string, badge int) (ok bool, error error) {
	if len(sound) == 0 {
		sound = models.JTUMENGPUSH_SOUND_DEFAUL
	}
	deviceInfo, err := DeviceInfoWithAuthorId(authorId)
	if err != nil {
		return false, err
	}
	deviceToken := deviceInfo.DeviceToken


	switch deviceInfo.Manufacturers {
	/* 华为推送
	 推消息支持的android APP包名称最大为128个字节。
	 消息内容最大限制为2K。
	*/
	case models.DEVICE_INFO_MANUFACTURERS_HUAWEI:
		huaweiMessage := message.Content
		if len(huaweiMessage) > 500 {
			huaweiMessage = util.SubstrByLength(huaweiMessage,0,500) + "..."
		}
		error = HuaweiPushToSingle(deviceToken, huaweiMessage, message.ActionUrl)
		error = hwSingleSend_TOUCHUAN(deviceToken, huaweiMessage, message.ActionUrl)
		/* 魅族推送
		 推送标题, 【string 必填，字数限制1~32字符】
		 推送内容, 【string 必填，字数限制1~100字符】
		*/
	case models.DEVICE_INFO_MANUFACTURERS_MEIZU:
		meizuMessage := message.Content
		if len(meizuMessage) > 90 {
			meizuMessage = util.SubstrByLength(meizuMessage,0,90) + "..."
		}
		error = MeizuPushToDeviceTokens(deviceToken, meizuMessage, message.ActionUrl)
		/* 小米推送
		 推送标题 Push title：Enter less than 50 characters(including spaces)（最长50字（单个字母或空格长度为1））；
		 消息摘要 Message summary：Enter less than 128 characters(including spaces) （最长128字（单个字母或空格长度为1））；
		 推送计划名称（Push plan name）和推送计划描述（Push plan description）均要求小于200字。
		*/
	case models.DEVICE_INFO_MANUFACTURERS_XIAOMI:
		xiaomiMessage := message.Content
		if len(xiaomiMessage) > 120 {
			xiaomiMessage = util.SubstrByLength(xiaomiMessage,0,120) + "..."
		}
		error = XiaomiPushToDeviceTokens(deviceToken, xiaomiMessage, message.ActionUrl)
	default:
		umengMessage := message.Content
		if len(umengMessage) > 120 {
			umengMessage = util.SubstrByLength(umengMessage,0,120) + "..."
		}
		if deviceInfo.System == models.DEVICE_INFO_SYSTERM_ANDROID {
			/* 其他安卓推送（友盟推送）
			 通知的标题（title）不允许全是空白字符且长度小于50
			 通知的内容（text）不允许全是空白字符且长度小于128（通知的标题和内容必填，一个中英文字符均计算为1）
			*/
			error = UmengPushToUsersFoAndroid_Notification(deviceToken, "辣课", umengMessage, message.ActionUrl, sound)
			error = UmengPushToUsersFoAndroid_Message(deviceToken, "辣课", umengMessage, message.ActionUrl, sound)
		} else {
			/* IOS推送（友盟推送）
			 通知的标题（title）不允许全是空白字符且长度小于50
			 通知的内容（text）不允许全是空白字符且长度小于128（通知的标题和内容必填，一个中英文字符均计算为1）
			*/

			fmt.Println(umengMessage)
			error = UmengPushToUsersForiOS(deviceToken, umengMessage, message.ActionUrl, sound, badge)
		}
	}

	ok = true
	if error != nil {
		ok = false
	}
	return
}

func DeviceInfoWithAuthorId(authorId int64) (deviceInfo *models.UserSignInDeviceInfo, error error) {
	var dInfo models.UserSignInDeviceInfo
	has, err := base.DBEngine.Table("user_sign_in_device_info").Where("author_id = ?", authorId).Get(&dInfo)
	if err != nil {
		util.Logger.Info("DeviceInfoWithAuthorId err = ", err.Error())
		return nil, err
	}
	if has {
		return &dInfo, nil
	} else {
		err := errors.New("用户设备信息不存在")
		util.Logger.Info("DeviceInfoWithAuthorId err:" + err.Error())
		return nil, err
	}
	return
}

//-------------------------------------------华为------------------------------------------------------------------

func HuaweiPushToSingle(deviceToken, msg, actionUrl string) (err error) {

	err = HwSingleSend(deviceToken, msg, actionUrl)
	if err != nil {
		util.Logger.Info("HuaweiPushToSingle err  = ", err.Error())
		return err
	}
	return
}


func HwSingleSend(deviceToken, msg, actionUrl string) (err error) {
	nowTime := time.Now()

	accessToken := ""
	systemConfig := models.SystemConfig{}
	hasAccessToken, _ := base.DBEngine.Table("system_config").Where("program = 'huawei_token' ").Get(&systemConfig)
	if hasAccessToken {
		accessToken = systemConfig.ProgramValue
	}

	util.Logger.Info("HwSingleSend models.HwAccessToken")
	util.Logger.Info(accessToken)

	data := url.Values{}
	data.Set("access_token", accessToken)
	data.Set("nsp_svc", "openpush.message.api.send")
	data.Set("nsp_ts", strconv.FormatInt(nowTime.UTC().Unix(), 10))
	data.Set("device_token_list", "[\""+deviceToken+"\"]")
	if actionUrl == "" {
		data.Set("payload", "" +
			"{\"hps\":" +
			"	{\"msg\":" +
			"{\"type\":3," +
			"\"body\":{\"content\":\""+msg+"\",\"title\":\"辣课\"}," +
			"\"action\":{\"type\":1," +
			"\"param\":{\"intent\":\""+actionUrl+"\",\"appPkgName\":\"com.luosuo.rml\"}" +
			"}" +
			"}," +
			"\"ext\":{\"biTag\":\"Trump\"}" +
			"}" +
			"}")
	} else {
		data.Set("payload", "" +
			"{\"hps\":" +
			"	{\"msg\":" +
			"{\"type\":3," +
			"\"body\":{\"actionUrl\":\""+actionUrl+"\",\"content\":\""+msg+"\",\"title\":\"辣课\"}," +
			"\"action\":{\"type\":1," +
			"\"param\":{\"intent\":\""+actionUrl+"\",\"appPkgName\":\"com.luosuo.rml\"}" +
			"}" +
			"}," +
			"\"ext\":{\"biTag\":\"Trump\"}" +
			"}" +
			"}")
	}

	u, _ := url.ParseRequestURI("https://api.push.hicloud.com/pushsend.do?nsp_ctx=%7b%22ver%22%3a%221%22%2c+%22appId%22%3a%22"+models.JTHwClientId+"%22%7d")
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()+"&"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	nspStatus := 0
	if len(resp.Header["Nsp_status"]) > 0 {
		nspStatus, _ = strconv.Atoi(resp.Header["Nsp_status"][0])
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 || nspStatus > 0 {
		if nspStatus == 6 {
			util.Logger.Warn("HwSingleSend err ", string(body), resp.StatusCode, " ", nspStatus)
			_, err = HwPostToken()
			if err != nil {
				return
			}
			return HwSingleSend(deviceToken,msg,actionUrl)
		}
		util.Logger.Warn("HwSingleSend err ", string(body), resp.StatusCode, " ", nspStatus)
		err = errors.New("HwSingleSend  err failed")
		return
	}

	response := models.HwRequestResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	if response.Resultcode > 0 {
		util.Logger.Warn("HwSingleSend  err ", string(body),response.Resultcode)
		err = errors.New(response.Message)
		return
	}

	return
}

func hwSingleSend_TOUCHUAN(deviceToken, msg, actionUrl string) (err error) {
	nowTime := time.Now()

	accessToken := ""
	systemConfig := models.SystemConfig{}
	hasAccessToken, _ := base.DBEngine.Table("system_config").Where("program = 'huawei_token' ").Get(&systemConfig)
	if hasAccessToken {
		accessToken = systemConfig.ProgramValue
	}

	util.Logger.Info("hwSingleSend_TOUCHUAN models.HwAccessToken")
	util.Logger.Info(accessToken)

	data := url.Values{}
	data.Set("access_token", accessToken)
	data.Set("nsp_svc", "openpush.message.api.send")
	data.Set("nsp_ts", strconv.FormatInt(nowTime.UTC().Unix(), 10))
	data.Set("device_token_list", "[\""+deviceToken+"\"]")
	if actionUrl == "" {
		data.Set("payload", "" +
			"{\"hps\":" +
			"	{\"msg\":" +
			"{\"type\":1," +
			"\"body\":{\"content\":\""+msg+"\",\"title\":\"辣课\"}," +
			"\"action\":{\"type\":3," +
			"\"param\":{\"intent\":\""+actionUrl+"\",\"appPkgName\":\"com.luosuo.rml\"}" +
			"}" +
			"}," +
			"\"ext\":{\"biTag\":\"Trump\"}" +
			"}" +
			"}")
	} else {
		data.Set("payload", "" +
			"{\"hps\":" +
			"	{\"msg\":" +
			"{\"type\":1," +
			"\"body\":{\"actionUrl\":\""+actionUrl+"\",\"msg\":\""+msg+"\",\"time\":\""+strconv.FormatInt(util.UnixOfBeijingTime(), 10)+"\",\"title\":\"辣课\"}," +
			"}," +
			"}" +
			"}")
	}

	u, _ := url.ParseRequestURI("https://api.push.hicloud.com/pushsend.do?nsp_ctx=%7b%22ver%22%3a%221%22%2c+%22appId%22%3a%22"+models.JTHwClientId+"%22%7d")
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()+"&")) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	nspStatus := 0
	if len(resp.Header["Nsp_status"]) > 0 {
		nspStatus, _ = strconv.Atoi(resp.Header["Nsp_status"][0])
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 || nspStatus > 0 {
		if nspStatus == 6 {
			util.Logger.Warn("hwDWQW_TOUCHUAN hwSingleSendDWQW push err1 ", string(body), resp.StatusCode, " ", nspStatus)
			_, err = HwPostToken()
			if err != nil {
				return
			}
			return hwSingleSend_TOUCHUAN(deviceToken,msg,actionUrl)
		}
		util.Logger.Warn("hwDWQW_TOUCHUAN hwSingleSendDWQW push err2 ", string(body), resp.StatusCode, " ", nspStatus)
		err = errors.New("hwDWQW_TOUCHUAN hwSingleSendDWQW push failed")
		return
	}

	response := models.HwRequestResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	if response.Resultcode > 0 {
		util.Logger.Warn("hwDWQW_TOUCHUAN hwSingleSendDWQW push err3 ", string(body),response.Resultcode)
		err = errors.New(response.Message)
		return
	}

	return
}




//刷新token
func HwPostToken() (token string, err error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Add("client_id", models.JTHwClientId)
	data.Add("client_secret", models.JTHwClientSecret)

	u, _ := url.ParseRequestURI(models.JTHw_ACCESS_TOKEN_URL)
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		util.Logger.Info("HwPostToken client  Do err = ", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("HwPostToken ioutil err = ", err.Error())
		return "", err
	}

	response := models.HwTokenResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("HwPostToken Unmarshal err = ", err.Error())
		return "", err
	}
	token = response.AccessToken
	models.HwAccessToken = token
	systemConfig := models.SystemConfig{}
	has, err1 := base.DBEngine.Table("system_config").Where("program = 'huawei_token' ").Get(&systemConfig)
	if err1 != nil {
		util.Logger.Info("HwPostToken get systemConfig err = ", err.Error())
		return "", err
	}
	if has {
		systemConfig.ProgramValue = token
		systemConfig.ProgramExpireTime = util.UnixOfBeijingTime() + response.ExpiresIn
		base.DBEngine.Table("system_config").Where("r_id = ? ", systemConfig.RId).Cols("program_value", "program_expire_time").Update(&systemConfig)
	} else {
		systemConfig.Program = "huawei_token"
		systemConfig.ProgramValue = token
		systemConfig.ProgramExpireTime = util.UnixOfBeijingTime() + response.ExpiresIn
		base.DBEngine.Table("system_config").InsertOne(&systemConfig)
	}
	return
}




//-------------------------------------------魅族------------------------------------------------------------------

func MeizuPushToDeviceTokens(deviceTokens, msg, actionUrl string) (err error) {
	if len(actionUrl) == 0 {
		actionUrl = models.JTDefault_ActionUrl
	}
	msgData := meizuMsgWithContentAncActiongUrl(msg, actionUrl)
	jsonBytes, err := json.Marshal(msgData)
	if err != nil {
		util.Logger.Info("MeizuPushToDeviceTokens Marshal err = ", err.Error())
		return err
	}
	messageJson := string(jsonBytes)

	data := make(map[string]string)
	data["messageJson"] = messageJson
	data["appId"] = models.JTMeizuAppId
	data["pushIds"] = deviceTokens

	err = meizuPush("garcia/api/server/push/varnished/pushByPushId", data)


	return
}



func meizuMsgWithContentAncActiongUrl(content, actionUrl string) (meizuMsg map[string]interface{}) {
	meizuMsg = make(map[string]interface{})

	//advanceInfo := map[string]interface{}{"suspend":1,"clearNoticeBar":1}
	pushTimeInfo := map[string]int{"offLine": 1, "validTime": 24}
	meizuMsg["pushTimeInfo"] = pushTimeInfo

	clickTypeInfo := map[string]interface{}{"clickType": 2, "url": actionUrl}
	meizuMsg["clickTypeInfo"] = clickTypeInfo

	noticeExpandInfo := map[string]interface{}{"noticeExpandType": 0}
	meizuMsg["noticeExpandInfo"] = noticeExpandInfo

	noticeBarInfo := map[string]interface{}{"noticeBarType": 0, "title": "辣课", "content": content}
	meizuMsg["noticeBarInfo"] = noticeBarInfo

	return meizuMsg
}


func meizuSignParamers(data map[string]string) (sign string) {
	keys := make([]string, len(data))

	i := 0
	for k := range data {
		keys[i] = k
		i++
	}

	//fmt.Println("keys", keys)
	sort.Strings(keys)
	//fmt.Println("sorted keys", keys)

	strTemp := ""
	for _, key := range keys {
		strTemp = strTemp + key + "=" + data[key]
	}

	strTemp = strTemp + models.JTMeizuAppSecret

	hasher := md5.New()
	hasher.Write([]byte(strTemp))
	sign = hex.EncodeToString(hasher.Sum(nil))

	return
}


func meizuPush(api string, data map[string]string) (err error) {
	sign := meizuSignParamers(data)
	data["sign"] = sign

	inputData := url.Values{}
	for key, v := range data {
		inputData.Set(key, v)
	}

	u, _ := url.ParseRequestURI(fmt.Sprintf("%s/%s", models.JTMEIZU_PUSH_URL, api))
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(inputData.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(inputData.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("meizuPush ioutil err = ", err.Error())
		return err
	}

	if resp.StatusCode != 200 {
		util.Logger.Warn("MEIZU push err ", string(body))
		err = errors.New("MEIZU push failed")
		return err
	}

	response := models.MeizuPushResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("meizuPush ioutil err = ", err.Error())
		return err
	}

	if response.Code != "200" {
		util.Logger.Warn("MEIZU Unmarshal err ", string(body))
		err = errors.New(response.Message)
		return err
	}

	//util.Logger.Debug(string(body))

	return
}





//-------------------------------------------小米------------------------------------------------------------------

func XiaomiPushToDeviceTokens(deviceTokens, msg, actionUrl string) (err error) {
	if len(actionUrl) == 0 {
		actionUrl = models.JTDefault_ActionUrl
	}
	err = XiaomiPush(deviceTokens, msg, actionUrl)
	return
}

func XiaomiPush(deviceTokens string, msg, actionUrl string) (err error) {
	data := url.Values{}
	data.Set("payload", msg)
	data.Set("extra.notify_effect", "2")
	data.Set("title", "辣课")
	data.Set("description", msg)
	data.Set("extra.intent_uri", actionUrl)
	data.Set("pass_through", "0")
	data.Set("restricted_package_name", models.JTXmPackageName)
	data.Set("registration_id", deviceTokens)
	data.Set("notify_type", "1")

	u, _ := url.ParseRequestURI(models.JTXM_ACCESS_URL)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	r.Header.Add("Authorization", "key="+models.JTXmAppSecret)

	resp, err := client.Do(r)
	if err != nil {
		util.Logger.Info("XiaomiPush client Do err = ", err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("XiaomiPush ioutil err = ", err.Error())
		return err
	}

	if resp.StatusCode != 200 {
		util.Logger.Warn("XiaomiPush err ", string(body))
		err = errors.New("XiaomiPush failed")
		return err
	}
	response := models.XmRequestResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("XiaomiPush Unmarshal err = ", err.Error())
		return err
	}

	if response.Code > 0 {
		util.Logger.Warn("XiaomiPush err:", string(response.Description))
		util.Logger.Warn("XiaomiPush err reason:", string(response.Reason))
		err = errors.New(response.Description)
		return err
	}

	return
}




//-------------------------------------------友盟------------------------------------------------------------------
func UmengPushToUsersFoAndroid_Notification(deviceTokens, title, msg, actionUrl, sound string) (err error)  {
	msgMap := make(map[string]string)
	msgMap["msg"] = msg
	msgMap["title"] = title
	msgMap["actionUrl"] = actionUrl
	msgMap["time"] = strconv.FormatInt(util.UnixOfBeijingTime(), 10)

	jsonBytes, err := json.Marshal(msgMap)
	if err != nil {
		util.Logger.Info("UmengPushToUsersFoAndroid err = ", err.Error())
		return err
	}

	params := make(map[string]interface{})
	params["appkey"] = models.JTUmengAndroidAppkey
	params["timestamp"] = strconv.FormatInt(util.UnixOfBeijingTime(), 10)
	if strings.Contains(deviceTokens, ",") {
		params["type"] = "listcast"
	} else {
		params["type"] = "unicast"
	}
	params["device_tokens"] = deviceTokens

	payload := make(map[string]interface{})
	body := make(map[string]string)


	body["sound"] = sound

	body["ticker"] = "辣课"
	body["title"] = title
	body["text"] = msg
	body["after_open"] = "go_custom"
	body["custom"] = string(jsonBytes)

	payload["body"] = body
	payload["display_type"] = "notification"
	payload["actionUrl"] = actionUrl
	params["payload"] = payload
	params["production_mode"] = "true"


	_, err = UmengPush(true, params)
	if err != nil {
		util.Logger.Info("UmengPushToUsersFoAndroid err = ", err.Error())
		return err
	}
	return
}

func UmengPushToUsersFoAndroid_Message(deviceTokens, title, msg, actionUrl, sound string) (err error)  {
	msgMap := make(map[string]string)
	msgMap["msg"] = msg
	msgMap["title"] = title
	msgMap["actionUrl"] = actionUrl
	msgMap["time"] = strconv.FormatInt(util.UnixOfBeijingTime(), 10)

	jsonBytes, err := json.Marshal(msgMap)
	if err != nil {
		util.Logger.Info("UmengPushToUsersFoAndroid err = ", err.Error())
		return err
	}

	params := make(map[string]interface{})
	params["appkey"] = models.JTUmengAndroidAppkey
	params["timestamp"] = strconv.FormatInt(util.UnixOfBeijingTime(), 10)
	if strings.Contains(deviceTokens, ",") {
		params["type"] = "listcast"
	} else {
		params["type"] = "unicast"
	}
	params["device_tokens"] = deviceTokens

	payload := make(map[string]interface{})
	body := make(map[string]string)


	body["sound"] = sound

	body["ticker"] = "辣课"
	body["title"] = title
	body["text"] = msg
	body["custom"] = string(jsonBytes)

	payload["body"] = body
	payload["display_type"] = "message"
	payload["actionUrl"] = actionUrl
	params["payload"] = payload
	params["production_mode"] = "true"


	_, err = UmengPush(true, params)
	if err != nil {
		util.Logger.Info("UmengPushToUsersFoAndroid err = ", err.Error())
		return err
	}
	return
}

func UmengPushToAllForAndroid(title, msg, actionUrl string) (err error) {
	params := make(map[string]interface{})
	params["appkey"] = models.JTUmengAndroidAppkey
	params["timestamp"] = strconv.FormatInt(util.UnixOfBeijingTime(), 10)

	params["type"] = "broadcast"
	payload := make(map[string]interface{})
	body := make(map[string]string)

	body["ticker"] = "辣课"
	body["title"] = title
	body["text"] = msg
	payload["body"] = body
	payload["display_type"] = "notification"

	params["payload"] = payload
	if beego.BConfig.RunMode != base.RUN_MODE_PROD {
		params["production_mode"] = "false"
	} else {
		params["production_mode"] = "true"
	}

	_, err = UmengPush(true, params)
	return err
}



func UmengPushToUsersForiOS(deviceTokens, msg, actionUrl, sound string, badge int) (err error) {
	params := make(map[string]interface{})
	params["appkey"] = models.JTUmengiOSAppkey
	params["timestamp"] = strconv.FormatInt(util.UnixOfBeijingTime(), 10)
	if strings.Contains(deviceTokens, ",") {
		params["type"] = "listcast"
	} else {
		params["type"] = "unicast"
	}
	params["device_tokens"] = deviceTokens

	payload := make(map[string]interface{})
	aps := make(map[string]interface{})
	aps["alert"] = msg
	aps["sound"] = sound
	if badge == 1 {
		aps["badge"] = 0
	} else {
		aps["badge"] = 1
	}
	payload["aps"] = aps
	payload["actionUrl"] = actionUrl

	payload["time"] = strconv.FormatInt(util.UnixOfBeijingTime(), 10)

	params["payload"] = payload
	params["production_mode"] = "true"


	_, err = UmengPush(false, params)
	if err != nil {
		util.Logger.Info("UmengPushToUsersForiOS err = ", err.Error())
		return err
	}
	return
}





func UmengPush(isAndroid bool, params map[string]interface{}) (body []byte, err error) {
	signResult := ""
	var jsonBytes []byte
	if isAndroid {
		signResult, jsonBytes, err = androidPostSign(params)
	} else {
		signResult, jsonBytes, err = iOSPostSign(params)
	}
	if err != nil {
		util.Logger.Info("UmengPush PostSign err = ", err.Error())
		return nil, err
	}
	var Url *url.URL
	Url, err = url.Parse(models.JTUMENG_PUSH_URL)
	if err != nil {
		util.Logger.Info("UmengPush Parse err = ", err.Error())
		return nil, err
	}
	Url.RawQuery = "sign=" + signResult
	req, err := http.NewRequest("POST", Url.String(), bytes.NewBuffer(jsonBytes))

	fmt.Println(string(jsonBytes))
	if err != nil {
		util.Logger.Info("UmengPush NewRequest err = ", err.Error())
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.Logger.Info("UmengPush client Do err = ", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("UmengPush ioutil err = ", err.Error())
		return nil, err
	}

	if resp.StatusCode != 200 {
		response := models.UmengPushResponse{}
		err1 := json.Unmarshal(body, &response)
		if err1 == nil {
			util.Logger.Warnf("umeng push err code %s", response.Data.ErrorCode)
			util.Logger.Info("UmengPush params:", params)
			util.Logger.Info("UmengPush response Body:", string(body))
		} else {
			util.Logger.Warnf("umeng push Unmarshal err %s", err1.Error())
		}
		err = errors.New("umeng push err code " + response.Data.ErrorCode)
		return nil, err
	}
	return
}



func postSign(method, url, masterSecret string, params map[string]interface{}) (string, []byte, error) {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		util.Logger.Info("postSign err = ", err.Error())
		return "", nil, err
	}

	signSource := strings.ToUpper(method) + url + string(jsonBytes) + masterSecret

	hasher := md5.New()
	hasher.Write([]byte(signSource))

	return hex.EncodeToString(hasher.Sum(nil)), jsonBytes, nil
}

func iOSPostSign(params map[string]interface{}) (string, []byte, error) {
	return postSign("POST", models.JTUMENG_PUSH_URL, models.JTUmengiOSAppMasterSecret, params)
}

func androidPostSign(params map[string]interface{}) (string, []byte, error) {
	return postSign("POST", models.JTUMENG_PUSH_URL, models.JTUmengAndroidAppMasterSecret, params)
}

//生成通知跳转链接
//key为需要跳转的页面的key
//pramas为页面需要的参数
func JumpUrlWithKeyAndPramas(key string, pramas map[string]string) (urlStr string) {
	url := url.URL{}

	url.Scheme = "renmailian"
	url.Host = "app.renmailian.cn"

	// Path
	url.Path = key

	if pramas != nil {
		// Query Parameters
		q := url.Query()
		for key, value := range pramas {
			q.Set(key, value)
		}
		url.RawQuery = q.Encode()
	}

	return url.String()
}


//----------------------------------------micro service----------------------------------------------------------------
func (this *PushController) PushMessageToUser(ctx context.Context, req *go_micro_messageservice.PushMessageToUserParam, rsp *go_micro_messageservice.PushMessageToUserResponse) error {
	util.Logger.Info("service PushMessageToUser")
	util.Logger.Info("req")
	util.Logger.Info(req)

	message := models.AppMessage{}
	message.Content = req.Content
	message.ReceiverId, _ = strconv.ParseInt(req.ReceiverId, 0, 64)
	message.ActionUrl = req.ActionUrl
	message.Type, _ = strconv.Atoi(req.Type)

	_, err := PushMessageToUser(message.ReceiverId, &message, "", 0)

	if err != nil {
		util.Logger.Info("PushToSingle  err = ", err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.PushMessageError100)
		return err
	}

	//_, err = base.DBEngine.Table("app_message").InsertOne(&message)

	rsp.Result = "success"

	return nil
}

func (this *PushController) JumpUrlWithKeyAndParams(ctx context.Context, req *go_micro_messageservice.JumpUrlWithKeyAndParamsParam, rsp *go_micro_messageservice.JumpUrlWithKeyAndParamsResponse) error {
	util.Logger.Info("service JumpUrlWithKeyAndParams")
	util.Logger.Info("req")
	util.Logger.Info(req)

	rsp.UrlStr = JumpUrlWithKeyAndPramas(req.Key, req.Map)

	return nil
}
