/*
@Time : 2019/2/25 下午2:20 
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

//定时任务，每分钟检查公众号access_token是否过期，过期则重新获取
func requestAuthorizerAccessToken(){
	//市场暂无第三方平台
	if beego.BConfig.RunMode == "prod" {
		return
	}

	util.Logger.Info("定时任务，每分钟检查公众号access_token是否过期，过期则重新获取")
	//提前10分钟查看
	refreshTime := util.UnixOfBeijingTime() + 10 * 60

	var componentAccessToken models.SystemConfig
	hasAccessToken, _ := base.DBEngine.Table("system_config").Where("program='component_access_token'").Get(&componentAccessToken)
	if !hasAccessToken {
		util.Logger.Info("requestAuthorizerAccessToken component_access_token 未找到")
		return
	}

	var authList []models.AuthInfo
	base.DBEngine.Table("auth_info").Where("auth_refresh_token is not null and auth_refresh_token != ''").And("auth_appid != 'wx570bc396a51b8ff8'").Find(&authList)
	for _, auth := range authList {

		if auth.AuthAccessTokenExpireTime > refreshTime {
			continue
		}

		data := "{\"component_appid\":\""+models.JTThirdPartyPlatformAppId+"\",\"authorizer_appid\":\""+auth.AuthAppid+"\",\"authorizer_refresh_token\":\""+auth.AuthRefreshToken+"\"}"

		u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=" + componentAccessToken.ProgramValue)
		urlStr := fmt.Sprintf("%v", u)

		client := &http.Client{}
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(data)))
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add("Content-Length", strconv.Itoa(len(data)))

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
			util.Logger.Info("requestAuthorizerAccessToken err :resp.StatusCode != 200")
			return
		}

		response := models.RefreshAuthorizerAccessTokenJsonBody{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			util.Logger.Info("requestAuthorizerAccessToken json.Unmarshal(body, &response) err :" + err.Error())
			return
		}

		if response.AuthorizerAccessToken != "" {
			auth.AuthAccessToken = response.AuthorizerAccessToken
			auth.AuthAccessTokenExpireTime = util.UnixOfBeijingTime() + response.ExpiresIn
			auth.AuthRefreshToken = response.AuthorizerRefreshToken
			base.DBEngine.Table("auth_info").Where("id=?", auth.Id).Cols("auth_access_token", "auth_access_token_expire_time", "auth_refresh_token").Update(&auth)
		} else {
			util.Logger.Info("requestAuthorizerAccessToken response.AuthorizerAccessToken == nil ")
			util.Logger.Info(string(body))
		}


		util.Logger.Info("刷新了" + auth.NickName + "(appid:" + auth.AuthAppid + ")的access_token")
	}
}
