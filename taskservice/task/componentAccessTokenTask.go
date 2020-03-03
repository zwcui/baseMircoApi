/*
@Time : 2019/2/25 上午11:38 
@Author : zwcui
@Software: GoLand
*/
package task

import (
	"jingting_server/taskservice/models"
	"jingting_server/taskservice/base"
	"jingting_server/taskservice/util"
	"strconv"
	"io/ioutil"
	"net/url"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

//每分钟检查component_access_token是否过期，过期则重新获取
func RequestComponentAccessToken () string {
	//市场暂无第三方平台
	if beego.BConfig.RunMode == "prod" {
		return ""
	}

	util.Logger.Info("定时任务，每分钟检查component_access_token是否过期，过期则重新获取")
	//提前10分钟查看
	refreshTime := util.UnixOfBeijingTime() + 10 * 60

	var componentAccessToken models.SystemConfig
	hasAccessToken, _ := base.DBEngine.Table("system_config").Where("program='component_access_token'").Get(&componentAccessToken)
	if hasAccessToken && componentAccessToken.ProgramExpireTime > refreshTime {
		return	""
	}

	var ticket models.SystemConfig
	hasTicket, _ := base.DBEngine.Table("system_config").Where("program='component_verify_ticket'").Get(&ticket)
	if !hasTicket {
		util.Logger.Info("重新获取component_access_token时ticket未找到")
		return	""
	}

	data := "{\"component_appid\":\""+models.JTThirdPartyPlatformAppId+"\" ,\"component_appsecret\": \""+models.JTThirdPartyPlatformAppSecret+"\",\"component_verify_ticket\": \""+ticket.ProgramValue+"\"}"

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/component/api_component_token")
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(data)))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if err != nil {
		util.Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return	""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return	""
	}

	if resp.StatusCode != 200 {
		util.Logger.Info("requestComponentAccessToken err :resp.StatusCode != 200")
		return	""
	}

	util.Logger.Info("body:"+string(body))

	response := models.ComponentAccessTokenJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("requestComponentAccessToken json.Unmarshal(body, &response) err :" + err.Error())
		return	""
	}

	if response.ComponentAccessToken == "" {
		errorResponse := models.ErrorJsonBody{}
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			util.Logger.Info("generateQrcode json.Unmarshal(body, &response) err :" + err.Error())
			return ""
		}
		util.Logger.Info("RequestComponentAccessToken err: "+errorResponse.Errmsg)
		return ""
	}

	if hasAccessToken {
		componentAccessToken.ProgramValue = response.ComponentAccessToken
		componentAccessToken.ProgramExpireTime = util.UnixOfBeijingTime() + response.ExpiresIn
		base.DBEngine.Table("system_config").Where("r_id=?", componentAccessToken.RId).Cols("program_value", "program_expire_time").Update(&componentAccessToken)
	} else {
		componentAccessToken.Program = "component_access_token"
		componentAccessToken.ProgramValue = response.ComponentAccessToken
		componentAccessToken.ProgramExpireTime = util.UnixOfBeijingTime() + response.ExpiresIn
		base.DBEngine.Table("system_config").InsertOne(&componentAccessToken)
	}
	return response.ComponentAccessToken
}
