/*
@Time : 2019/3/13 上午9:49 
@Author : zwcui
@Software: GoLand
*/
package util

import (
	"net/url"
	"net/http"
	"bytes"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"jingting_server/messageservice/models"
	"errors"
	"fmt"
)

//----------------------------客服消息--------------------------------------------

//发送文本客服消息
func RequestSendGZHTextCustomerServiceMessage(openid string, content string, accessToken string)(err error) {
	if content == "" {
		return nil
	}
	dataMap := map[string]string{
		"content" : content,
	}
	err = requestSendGZHCustomerServiceMessage(openid, 1, dataMap, accessToken)
	if err != nil {
		for errorCount := 0; errorCount < 10; errorCount++ {
			Logger.Info("重试第" + strconv.Itoa(errorCount + 1) + "次")
			err = requestSendGZHCustomerServiceMessage(openid, 1, dataMap, accessToken)
			if err == nil {
				break
			}
		}
	}
	return err
}

//发送图片客服消息
func RequestSendGZHMediaCustomerServiceMessage(openid string, mediaId string, accessToken string)(err error) {
	if mediaId == "" {
		return nil
	}
	dataMap := map[string]string{
		"media_id" : mediaId,
	}
	return requestSendGZHCustomerServiceMessage(openid, 2, dataMap, accessToken)
}

//发送图文客服消息
func RequestSendGZHPicAndTextCustomerServiceMessage(openid string, title string, description string, url string, picurl string, accessToken string)(err error) {
	if picurl == "" {
		return nil
	}
	dataMap := make(map[string]string)
	if url != "" {
		dataMap = map[string]string{
			"title" : title,
			"description" : description,
			"url" : url,
			"picurl" : picurl,
		}
	} else {
		dataMap = map[string]string{
			"title" : title,
			"description" : description,
			"picurl" : picurl,
		}
	}

	return requestSendGZHCustomerServiceMessage(openid, 6, dataMap, accessToken)
}

//发送客服消息
//msgType消息类型，1文本消息，2图片消息，3语音消息，4视频消息，5音乐消息，6图文消息（点击跳转到外链），7图文消息（点击跳转到图文消息页面），8菜单消息，9卡券，10小程序卡片（要求小程序与公众号已关联）
func requestSendGZHCustomerServiceMessage(openid string, msgType int, dataMap map[string]string, accessToken string)(err error) {
	msgTypeStr := ""
	data := ""

	switch msgType {
	case 1:
		msgTypeStr = "text"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"text\",\"text\":{\"content\":\"Hello World\"}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}}"
		break
	case 2:
		msgTypeStr = "image"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"image\",\"image\":{\"media_id\":\"MEDIA_ID\"}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}}"
		break
	case 3:
		msgTypeStr = "voice"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"voice\",\"voice\":{\"media_id\":\"MEDIA_ID\"}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}}"
		break
	case 4:
		msgTypeStr = "video"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"video\",\"video\":{\"media_id\":\"MEDIA_ID\",\"thumb_media_id\":\"MEDIA_ID\",\"title\":\"TITLE\",\"description\":\"DESCRIPTION\"}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}}"
		break
	case 5:
		msgTypeStr = "music"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"music\",\"music\":{\"title\":\"MUSIC_TITLE\",\"description\":\"MUSIC_DESCRIPTION\",\"musicurl\":\"MUSIC_URL\",\"hqmusicurl\":\"HQ_MUSIC_URL\",\"thumb_media_id\":\"THUMB_MEDIA_ID\"}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}}"
		break
	case 6:
		msgTypeStr = "news"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"news\",\"news\":{\"articles\": [{\"title\":\"Happy Day\",\"description\":\"Is Really A Happy Day\",\"url\":\"URL\",\"picurl\":\"PIC_URL\"}]}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{\"articles\": [{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}]}}"
		break
	case 7:
		msgTypeStr = "mpnews"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"mpnews\",\"mpnews\":{\"media_id\":\"MEDIA_ID\"}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}}"
		break
	case 8:
		msgTypeStr = "msgmenu"
		//data := "{\"touser\": \"OPENID\",\"msgtype\": \"msgmenu\",\"msgmenu\": {\"head_content\": \"您对本次服务是否满意呢? \",\"list\": [{\"id\": \"101\",\"content\": \"满意\"},{\"id\": \"102\",\"content\": \"不满意\"}],\"tail_content\": \"欢迎再次光临\"}}"
		data = "{\"touser\": \"OPENID\",\"msgtype\": \"msgmenu\",\"msgmenu\": {\"head_content\": \"您对本次服务是否满意呢? \",\"list\": [{\"id\": \"101\",\"content\": \"满意\"},{\"id\": \"102\",\"content\": \"不满意\"}],\"tail_content\": \"欢迎再次光临\"}}"
		break
	case 9:
		msgTypeStr = "wxcard"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"wxcard\",\"wxcard\":{\"card_id\":\"123dsdajkasd231jhksad\"}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}}"
		break
	case 10:
		msgTypeStr = "miniprogrampage"
		//data := "{\"touser\":\"OPENID\",\"msgtype\":\"miniprogrampage\",\"miniprogrampage\":{\"title\":\"title\",\"appid\":\"appid\",\"pagepath\":\"pagepath\",\"thumb_media_id\":\"thumb_media_id\"}}"
		data = "{\"touser\":\"" + openid + "\",\"msgtype\":\"" + msgTypeStr + "\",\"" + msgTypeStr + "\":{"
		for key, value := range dataMap {
			data += "\""+key+"\":\""+value+"\","
		}
		data = SubstrByLength(data, 0, len(data)-1)
		data += "}}"
		break
	}

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + accessToken)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(data)))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if err != nil {
		Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return	err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return	err
	}

	if resp.StatusCode != 200 {
		Logger.Info("requestSendCustomerServiceMessage err :resp.StatusCode != 200")
		return	err
	}

	Logger.Info("body:"+string(body))

	errorResponse := models.MessageErrorJsonBody{}
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		Logger.Info("requestSendCustomerServiceMessage json.Unmarshal(body, &response) err :" + err.Error())
		return err
	}

	if errorResponse.Errcode == 0 {
		return nil
	} else {
		return errors.New(errorResponse.Errmsg)
	}
}


//----------------------------模板消息--------------------------------------------

//获取微信公众号下配置的模板列表
func RequestQueryGZHTemplateList(accessToken string)(templateList []models.Template, err error) {

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?access_token=" + accessToken)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return templateList, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return templateList, err
	}

	if resp.StatusCode != 200 {
		Logger.Info("requestTemplateList err :resp.StatusCode != 200")
		return templateList, errors.New("requestTemplateList err :resp.StatusCode != 200")
	}

	Logger.Info("body:" + string(body))

	response := models.TemplateListJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Info("requestTemplateList json.Unmarshal(body, &response) err :" + err.Error())
		return templateList, err
	}

	if response.TemplateList == nil {
		errorResponse := models.ErrorJsonBody{}
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			Logger.Info("requestQueryTemplateList json.Unmarshal(body, &response) err :" + err.Error())
			return templateList, err
		}
		Logger.Info("requestQueryTemplateList err: "+errorResponse.Errmsg)
		return templateList, errors.New(errorResponse.Errmsg)
	}

	return response.TemplateList, nil
}

//删除模板
func RequestDeleteGZHTemplate(templateId string, accessToken string)(err error) {
	data := "{\"template_id\" : \""+templateId+"\"}"

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/template/del_private_template?access_token=" + accessToken)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(data)))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if err != nil {
		Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return	err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return	err
	}

	if resp.StatusCode != 200 {
		Logger.Info("requestDeleteTemplate err :resp.StatusCode != 200")
		return	err
	}

	Logger.Info("body:"+string(body))

	response := models.ComponentAccessTokenJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Info("requestDeleteTemplate json.Unmarshal(body, &response) err :" + err.Error())
		return	err
	}

	errorResponse := models.ErrorJsonBody{}
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		Logger.Info("requestDeleteTemplate json.Unmarshal(body, &response) err :" + err.Error())
		return err
	}

	if errorResponse.Errmsg == "" {
		return nil
	} else {
		return errors.New(errorResponse.Errmsg)
	}
}

//发送模板消息
func RequestSendGZHTemplateMessage(templateId string, openid string, jumpUrl string, templateData string, accessToken string)(err error) {
	data := "{\"touser\":\""+openid+"\",\"template_id\":\""+templateId+"\",\"url\":\""+jumpUrl+"\",\"data\":"

	//for key, value := range dataMap {
	//	data += "\""+key+"\": {\"value\":\""+value.Value+"\",\"color\":\""+value.Color+"\"},"
	//}
	data += templateData
	data += "}"
	Logger.Info(data)

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + accessToken)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(data)))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if err != nil {
		Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return	err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return	err
	}

	if resp.StatusCode != 200 {
		Logger.Info("requestSendTemplateMessage err :resp.StatusCode != 200")
		return	err
	}

	Logger.Info("body:"+string(body))

	errorResponse := models.MessageErrorJsonBody{}
	err = json.Unmarshal(body, &errorResponse)
	if err != nil {
		Logger.Info("requestSendTemplateMessage json.Unmarshal(body, &response) err :" + err.Error())
		return err
	}

	if errorResponse.Errcode == 0 {
		return nil
	} else {
		return errors.New(errorResponse.Errmsg)
	}
}
