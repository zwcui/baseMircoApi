/*
@Time : 2019/10/11 下午6:32 
@Author : zwcui
@Software: GoLand
*/
package task

import (
	"jingting_server/taskservice/models"
	"net/url"
	"fmt"
	"net/http"
	"jingting_server/taskservice/util"
	"io/ioutil"
	"encoding/json"
	"jingting_server/taskservice/base"
)

func checkGZHAccessToken(){
	util.Logger.Info("定时任务，每分钟检查辣课公众号access_token是否过期，过期则重新获取")

	//提前10分钟查看
	refreshTime := util.UnixOfBeijingTime() + 40 * 60

	var gzhConfig models.SystemConfig
	hasAccessToken, _ := base.DBEngine.Table("system_config").Where("program='ghz_access_token'").Get(&gzhConfig)
	if !hasAccessToken {
		resp, err := requestGZHAccessToken()
		util.Logger.Info("requestGZHAcceccToken")
		util.Logger.Info(resp)
		if err != nil {
			util.Logger.Info("requestGZHAcceccToken err:" + err.Error())
			return
		}
		gzhConfig.Description = "辣课公众号access_token"
		gzhConfig.Program = "ghz_access_token"
		gzhConfig.ProgramValue = resp.AccessToken
		gzhConfig.ProgramExpireTime = util.UnixOfBeijingTime() + resp.ExpiresIn

		util.Logger.Info(util.UnixOfBeijingTime())
		util.Logger.Info(resp.ExpiresIn)
		util.Logger.Info(gzhConfig.ProgramExpireTime)

		base.DBEngine.Table("system_config").InsertOne(&gzhConfig)
		util.Logger.Info("插入accesToken:"+gzhConfig.ProgramValue)
		return
	} else {
		if gzhConfig.ProgramExpireTime <= refreshTime {
			resp, err := requestGZHAccessToken()
			if err != nil {
				util.Logger.Info("requestGZHAcceccToken err:" + err.Error())
				return
			}
			util.Logger.Info("requestGZHAcceccToken")
			util.Logger.Info(resp)
			gzhConfig.ProgramValue = resp.AccessToken
			gzhConfig.ProgramExpireTime = util.UnixOfBeijingTime() + resp.ExpiresIn

			util.Logger.Info(util.UnixOfBeijingTime())
			util.Logger.Info(resp.ExpiresIn)
			util.Logger.Info(gzhConfig.ProgramExpireTime)

			base.DBEngine.Table("system_config").Where("r_id=?", gzhConfig.RId).AllCols().Update(&gzhConfig)
			util.Logger.Info("更新accesToken:"+gzhConfig.ProgramValue)
		}
		util.Logger.Info("accesToken:"+gzhConfig.ProgramValue)
		return
	}
}

//请求公众号access_token
func requestGZHAccessToken() (response RequestGZHAcceccTokenResponse, err error) {

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