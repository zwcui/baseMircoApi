/*
@Time : 2019/6/20 上午9:31 
@Author : lianwu
@File : CreateQRCode.go
@Software: GoLand
*/
package controllers

import (
	"strconv"
	"github.com/skip2/go-qrcode"
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/util"
	"os"
	"encoding/base64"
	"io/ioutil"
	"jingting_server/publicservice/base"
	"net/url"
	"fmt"
	"net/http"
	"encoding/json"
)


type QRCodeController struct {
	apiController
}


func (this *QRCodeController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
	}
	this.authorAuth()
}

// @Title 新增二维码
// @Description 新增二维码
// @Param   codeContent   formData        string  	true        "二维码存储内容"
// @Success 200 {string}
// @router /createQRCode [post]
func (this *QRCodeController) CreateQRCode() {
	codeContent := this.MustString("codeContent")
	//二维码生成图片
	filePath := strconv.FormatInt(util.UnixOfBeijingTime(), 10) + ".png"
	err := qrcode.WriteFile(codeContent, qrcode.Medium, 256, filePath)
	if err != nil {
		util.Logger.Info("qrcode.WriteFile err:" + err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CreateQRCodeError100, err.Error())
		return
	}
	image, err := ioutil.ReadFile(filePath)
	if err != nil {
		util.Logger.Info("ioutil.ReadFile(filePath) err:" + err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CreateQRCodeError200, err.Error())
		return
	}
	//删除图片
	del := os.Remove(filePath)
	if del != nil {
		util.Logger.Info("os.Remove(filePath) err:" + del.Error())
	}
	this.ReturnData = base64.StdEncoding.EncodeToString(image)
}

// @Title 生成微信临时二维码
// @Description 生成微信临时二维码
// @Param   data	   formData        string  	false        "二维码存储内容"
// @Success 200 {object} models.QrcodeJsonBody
// @router /createWechatQRCode [post]
func (this *QRCodeController) CreateWechatQRCode() {
	data := this.GetString("data", "")

	//var authInfo models.AuthInfo
	//hasAuthInfo, _ := base.DBEngine.Table("auth_info").Where("auth_appid = ?", models.JTGZHAppId).Get(&authInfo)
	//if !hasAuthInfo {
	//	this.ReturnData = util.GenerateAlertMessage(models.AuthorizeError300)
	//	return
	//}
	accessToken := getGZHAccessToken()
	util.Logger.Info("accessToken:" + accessToken)
	var flag bool = false
	response, err := util.GenerateWechatQrcodeWithDataStrByAccessToken(data, accessToken, 0)
	if err != nil {
		util.Logger.Info("createWechatQRCode GenerateWechatQrcodeWithDataStr err:" + err.Error())
		for errorCount := 0; errorCount < 10; errorCount++ {
			util.Logger.Info("重试第" + strconv.Itoa(errorCount) + "次")
			response, err = util.GenerateWechatQrcodeWithDataStrByAccessToken(data, accessToken, 0)
			if err == nil {
				flag = true
				break
			}
		}
	} else {
		flag = true
	}
	if flag {
		this.ReturnData = response
	} else {

		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}

}




//---------------------------------------方法---------------------------------------------------------

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

func getWechatQrcodeWithDataStrByAccessToken(data string, accessToken string, isLimit int, errorCount int) (qrcode models.QrcodeJsonBody, err error)  {
	response, err := util.GenerateWechatQrcodeWithDataStrByAccessToken(data, accessToken, isLimit)
	if err != nil {
		errorCount += 1
		util.Logger.Info("错误调用次数 :" + strconv.Itoa(errorCount))
		if errorCount > 10 {
			util.Logger.Info("GenerateWechatQrcodeWithDataStrByAccessToken err:" + err.Error())
			return qrcode, err
		} else {
			util.Logger.Info("重试第" + strconv.Itoa(errorCount) + "次")
			response, err = getWechatQrcodeWithDataStrByAccessToken(data, accessToken, 0, errorCount)

		}
	}
	util.Logger.Info("第" +strconv.Itoa(errorCount + 1 ) + "次返回")
	return response, nil


}
