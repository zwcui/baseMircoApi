/*
@Time : 2019/10/29 下午3:52 
@Author : lianwu
@File : testACcesstoken.go
@Software: GoLand
*/
package task

import (
	"jingting_server/taskservice/models"
	"jingting_server/taskservice/util"
	"jingting_server/taskservice/base"
	"net/url"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"bytes"
	"errors"
	"time"
)

func testAccesstoken () {

	fmt.Println("定时任务，每分钟测试accesstoken是否可用")

	var gzhConfig models.SystemConfig
	base.DBEngine.Table("system_config").Where("program='ghz_access_token'").Get(&gzhConfig)




	data := "jingtingjiaoyu"
	fileName := "./testAccesstoken.txt"
	//content := ""

	f, err := os.OpenFile(fileName,os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	accessToken := getGZHAccessToken()
	util.Logger.Info("accessToken:" + accessToken)
	var flag bool = false
	_, err = GenerateWechatQrcodeWithDataStrByAccessToken(data, accessToken, 0)
	if err != nil {
		util.Logger.Info("createWechatQRCode GenerateWechatQrcodeWithDataStr err:" + err.Error())
		for errorCount := 0; errorCount < 10; errorCount++ {
			util.Logger.Info("重试第" + strconv.Itoa(errorCount) + "次")
			_, err = GenerateWechatQrcodeWithDataStrByAccessToken(data, accessToken, 0)
			if err == nil {
				flag = true
				break
			}
		}
	} else {
		flag = true
	}



	if !flag {
		if gzhConfig.ProgramExpireTime <= util.UnixOfBeijingTime() + 40 * 60 {
			Count = 1
			ErrorCount =0
		} else {
			Count += 1
			ErrorCount +=1
		}
		util.Logger.Info("createWechatQRCode GenerateWechatQrcodeWithDataStr err:" + err.Error())

		//content = util.BeijingTime().Format("2006-01-02 15:04:05") + "：" + "到期时间：" + time.Unix(gzhConfig.ProgramExpireTime, 0).Format("2006-01-02 15:04:05") +" 次数：" + strconv.Itoa(ErrorCount) + "/" + strconv.Itoa(Count) + " ：调用失败 ：" + accessToken + "\n"
		//
		//f.Write([]byte(content))
		var test models.TestAccessToken
		test.TestTime = util.BeijingTime().Format("2006-01-02 15:04:05")
		test.ExpeireTime = time.Unix(gzhConfig.ProgramExpireTime, 0).Format("2006-01-02 15:04:05")
		test.ErrorCount = ErrorCount
		test.Count = Count
		test.ErrorAndCount = strconv.Itoa(ErrorCount) + "/" + strconv.Itoa(Count)
		test.Status = "调用失败"
		test.ErrorDetail = err.Error()
		test.Accesstoken = accessToken
		base.DBEngine.Table("test_access_token").InsertOne(&test)



	} else {
		if gzhConfig.ProgramExpireTime <= util.UnixOfBeijingTime() + 40 * 60 {
			Count = 1
			ErrorCount =1
		} else {
			Count += 1
		}
		//content = util.BeijingTime().Format("2006-01-02 15:04:05") + "：" + "到期时间：" + time.Unix(gzhConfig.ProgramExpireTime, 0).Format("2006-01-02 15:04:05") +" 次数：" + strconv.Itoa(ErrorCount) + "/" + strconv.Itoa(Count) +  " ：调用成功 ：" + accessToken + "\n"
		//f.Write([]byte(content))
		var test models.TestAccessToken
		test.TestTime = util.BeijingTime().Format("2006-01-02 15:04:05")
		test.ExpeireTime = time.Unix(gzhConfig.ProgramExpireTime, 0).Format("2006-01-02 15:04:05")
		test.ErrorCount = ErrorCount
		test.Count = Count
		test.ErrorAndCount = strconv.Itoa(ErrorCount) + "/" + strconv.Itoa(Count)
		test.Status = "调用成功"
		test.Accesstoken = accessToken
		base.DBEngine.Table("test_access_token").InsertOne(&test)

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


type RequestGZHCreateMenuResponse struct {
	Errcode					int					`description:"errcode" json:"errcode"`
	Errmsg					string				`description:"errmsg" json:"errmsg"`
}

//微信生成带参数的临时二维码
//isLimit 是否永久二维码，1是0否，永久的给生成活动使用
func GenerateWechatQrcodeWithDataStrByAccessToken (dataStr string, accessToken string, isLimit int) (qrcode models.QrcodeJsonBody, err error) {
	scene := dataStr

	actionName := ""
	if isLimit == 1 {
		actionName = "QR_LIMIT_STR_SCENE"
	} else {
		actionName = "QR_STR_SCENE"
	}

	data := "{\"expire_seconds\": 2592000, \"action_name\": \""+actionName+"\", \"action_info\": {\"scene\": {\"scene_str\": \""+scene+"\"}}}"

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="+accessToken)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(data)))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if err != nil {
		util.Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return qrcode, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return qrcode, err
	}

	if resp.StatusCode != 200 {
		util.Logger.Info("generateQrcode err :resp.StatusCode != 200")
		return qrcode, err
	}

	response := models.QrcodeJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		util.Logger.Info("generateQrcode json.Unmarshal(body, &response) err :" + err.Error())
		return qrcode, err
	}

	if response.Url == "" {
		errorResponse := models.ErrorJsonBody{}
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			util.Logger.Info("generateQrcode json.Unmarshal(body, &response) err :" + err.Error())
			return qrcode, err
		}
		return qrcode, errors.New(errorResponse.Errmsg)
	}

	return response, nil
}


func getWechatQrcodeWithDataStrByAccessToken(data string, accessToken string, isLimit int, errorCount int) (qrcode models.QrcodeJsonBody, err error)  {
	response, err := GenerateWechatQrcodeWithDataStrByAccessToken(data, accessToken, isLimit)
	if err != nil {
		errorCount += 1
		util.Logger.Info("错误调用次数 :" + strconv.Itoa(errorCount))
		if errorCount > 10 {
			util.Logger.Info("GenerateWechatQrcodeWithDataStrByAccessToken err:" + err.Error())
			return qrcode, err
		} else {
			util.Logger.Info("重试第" + strconv.Itoa(errorCount) + "次")
			getWechatQrcodeWithDataStrByAccessToken(data, accessToken, 0, errorCount)
		}
	}

	return response, nil


}

