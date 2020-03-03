/*
@Time : 2019/2/25 下午1:33 
@Author : zwcui
@Software: GoLand
*/
package task

import (
	"net/url"
	"net/http"
	"bytes"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"jingting_server/taskservice/util"
	"jingting_server/taskservice/models"
	"jingting_server/taskservice/base"
	"fmt"
	"github.com/astaxie/beego"
)

//每分钟检查pre_auth_code是否过期，过期则重新获取
func RequestPreAuthCode() string {
	//市场暂无第三方平台
	if beego.BConfig.RunMode == "prod" {
		return ""
	}

	util.Logger.Info("定时任务，每分钟检查pre_auth_code是否过期，过期则重新获取")
	//提前3分钟查看
	//refreshTime := util.UnixOfBeijingTime() + 3 * 60

	var preAuthCode models.SystemConfig
	hasPreAuthCode, _ := base.DBEngine.Table("system_config").Where("program='pre_auth_code'").Get(&preAuthCode)
	//if hasPreAuthCode && preAuthCode.ProgramExpireTime > refreshTime {
	//	return	""
	//}

	var componentAccessToken models.SystemConfig
	hasAccessToken, _ := base.DBEngine.Table("system_config").Where("program='component_access_token'").Get(&componentAccessToken)
	if !hasAccessToken {
		util.Logger.Info("重新获取pre_auth_code时componentAccessToken未找到")
		return	""
	}

	data := "{\"component_appid\":\""+models.JTThirdPartyPlatformAppId+"\"}"

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token="+componentAccessToken.ProgramValue)
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
		util.Logger.Info("requestPreAuthCode err :resp.StatusCode != 200")
		return	""
	}

	util.Logger.Info("body:"+string(body))

	response := models.PreAuthCodeJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("requestPreAuthCode json.Unmarshal(body, &response) err :" + err.Error())
		return	""
	}

	if response.PreAuthCode == "" {
		errorResponse := models.ErrorJsonBody{}
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			util.Logger.Info("generateQrcode json.Unmarshal(body, &response) err :" + err.Error())
			return ""
		}
		util.Logger.Info("RequestComponentAccessToken err: "+errorResponse.Errmsg)
		return ""
	}

	if hasPreAuthCode {
		preAuthCode.ProgramValue = response.PreAuthCode
		preAuthCode.ProgramExpireTime = util.UnixOfBeijingTime() + response.ExpiresIn
		base.DBEngine.Table("system_config").Where("r_id=?", preAuthCode.RId).Cols("program_value", "program_expire_time").Update(&preAuthCode)
	} else {
		preAuthCode.Program = "pre_auth_code"
		preAuthCode.ProgramValue = response.PreAuthCode
		preAuthCode.ProgramExpireTime = util.UnixOfBeijingTime() + response.ExpiresIn
		base.DBEngine.Table("system_config").InsertOne(&preAuthCode)
	}
	return response.PreAuthCode
}
