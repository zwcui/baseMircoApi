/*
@Time : 2019/3/5 上午10:37 
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
	"fmt"
	"errors"
	"github.com/skip2/go-qrcode"
	"os"
	"github.com/go-xorm/core"
	"io"
	"mime/multipart"
	"github.com/astaxie/beego"
)


//微信生成带参数的临时二维码
//isLimit 是否永久二维码，1是0否，永久的给生成活动使用
func GenerateWechatQrcodeWithDataMap(dataMap map[string]string, authInfo models.AuthInfo, isLimit int) (qrcode models.QrcodeJsonBody, err error) {
	scene := ""
	for key, value := range dataMap {
		scene += key + "=" + value + "&"
	}

	return GenerateWechatQrcodeWithDataStr(scene, authInfo, isLimit)
}

//微信生成带参数的临时二维码
//isLimit 是否永久二维码，1是0否，永久的给生成活动使用
func GenerateWechatQrcodeWithDataStr(dataStr string, authInfo models.AuthInfo, isLimit int) (qrcode models.QrcodeJsonBody, err error) {
	scene := dataStr

	actionName := ""
	if isLimit == 1 {
		actionName = "QR_LIMIT_STR_SCENE"
	} else {
		actionName = "QR_STR_SCENE"
	}

	data := "{\"expire_seconds\": 2592000, \"action_name\": \""+actionName+"\", \"action_info\": {\"scene\": {\"scene_str\": \""+scene+"\"}}}"

	u, _ := url.ParseRequestURI("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token="+authInfo.AuthAccessToken)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(data)))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := client.Do(r)
	if err != nil {
		Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return qrcode, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return qrcode, err
	}

	if resp.StatusCode != 200 {
		Logger.Info("generateQrcode err :resp.StatusCode != 200")
		return qrcode, err
	}

	response := models.QrcodeJsonBody{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Info("generateQrcode json.Unmarshal(body, &response) err :" + err.Error())
		return qrcode, err
	}

	if response.Url == "" {
		errorResponse := models.ErrorJsonBody{}
		err = json.Unmarshal(body, &errorResponse)
		if err != nil {
			Logger.Info("generateQrcode json.Unmarshal(body, &response) err :" + err.Error())
			return qrcode, err
		}
		return qrcode, errors.New(errorResponse.Errmsg)
	}

	return response, nil
}

//生成自定义二维码
func GenerateLocalQrcode(content string) (qrcodeUrl string, err error){
	if content == "" {
		return "", errors.New("没有相应链接")
	}
	currentPath, _ := os.Getwd()
	qrcodePath := currentPath + "/qrcode/"
	if !IsExist(qrcodePath) {
		err := os.MkdirAll(qrcodePath,os.ModePerm)
		// 创建文件夹
		if err != nil {
			Logger.Info("mkdir failed![%v]\n", err)
			return "", err
		} else {
			Logger.Info("mkdir success!\n")
		}
	}

	qrcodeName := core.Uuid + ".png"
	var qrcodeFile *os.File
	if !IsExist(qrcodePath + qrcodeName) {
		Logger.Info("create qrcodeFile")
		qrcodeFile, err = os.Create(qrcodePath + qrcodeName)
		if err != nil {
			Logger.Info("os.Create err:" + err.Error())
			return "", err
		}
	}

	err = qrcode.WriteFile(content, qrcode.Medium, 256, qrcodePath + qrcodeName)
	if err != nil {
		Logger.Info("write error:"+err.Error())
		return "", err
	}


	//上传文件服务器
	fileServerUrl := ""
	fileUrl := ""
	if beego.BConfig.RunMode == "dev" || beego.BConfig.RunMode == "test" {
		fileServerUrl = "http://jingting.vipask.net/v1/file/?fileType=image"
		fileUrl = "http://jingting.vipask.net/v1/file"
	} else if beego.BConfig.RunMode == "prod" {
		fileServerUrl = "http://jingtingedu.vipask.net/v1/file/?fileType=image"
		fileUrl = "http://jingtingedu.vipask.net/v1/file"
	}

	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)


	_, err = body_writer.CreateFormFile("userfile", qrcodePath + qrcodeName)
	if err != nil {
		Logger.Info("error writing to buffer")
		return "", err
	}

	//fh, err := os.Open(filePath)
	//if err != nil {
	//	Logger.Info("error opening file")
	//	return "", err
	//}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, qrcodeFile, close_buf)
	fi, err := qrcodeFile.Stat()
	if err != nil {
		Logger.Info("Error Stating file: %s", qrcodePath + qrcodeName)
		return "", err
	}
	req, err := http.NewRequest("POST", fileServerUrl, request_reader)
	if err != nil {
		return "", err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Logger.Info("resp, err := client.Do(r) err:" + err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Info("body, err := ioutil.ReadAll(resp.Body) err:" + err.Error())
		return "", err
	}

	Logger.Info(string(body))

	response := models.FileServerResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		Logger.Info("json.Unmarshal(body, &response) err :" + err.Error())
		return "", err
	}

	os.Remove(qrcodePath + qrcodeName)

	data := response.Data
	if data != nil && len(data) > 0 {
		return fileUrl + data[0].Uri, nil
	} else {
		return "", errors.New("上传文件服务器出错")
	}

	return "", nil
}


