/*
@Time : 2019/2/26 下午2:26 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/task"
	"strconv"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"errors"
	"math/rand"
	"sort"
	"crypto/sha1"
	"encoding/hex"
	"jingting_server/publicservice/util"
	"jingting_server/publicservice/base"
	"github.com/astaxie/beego/validation"
	"time"
	mathRand "math/rand"
	"strings"
)

type PublicController struct {
	apiController
}

func (this *PublicController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/getSystemComponentAccessToken", "get", []int{1}},
		{"/getSystemPreAuthCode", "get", []int{1}},
		{"/app-config", "post", []int{1}},
		{"/app-config/", "patch", []int{1}},
		{"/app-config/all", "get", []int{1}},
		{"/version", "post", []int{1}},
		{"/version/history", "get", []int{1}},
		}
	this.userAuth()
}

// @Title 公众号管理者登录（h5使用）
// @Description 公众号管理者登录（h5使用）
// @Param	nickName				formData		string  		true		"登录账户"
// @Param	password				formData		string	  		true		"登录密码"
// @Success 200 {object} models.UserShortContainer
// @router /signIn [post]
func (this *PublicController) SignIn() {
	nickName := this.MustString("nickName")
	password := this.MustString("password")

	user, err := UserWithNickName(nickName)
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}
	if user == nil {
		this.ReturnData = util.GenerateAlertMessage(models.UserError300)
		return
	}

	hashedPassword, err := util.EncryptPasswordWithSalt(password, user.Salt)
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}
	if user.Password != hashedPassword {
		this.ReturnData = util.GenerateAlertMessage(models.UserError400)
		return
	}

	var authInfo models.AuthInfo
	//如果没有默认选择公众号，则自动选择一个
	if user.DefaultAuthAppid == "" {
		hasAuth, _ := base.DBEngine.Table("auth_info").Join("LEFT OUTER", "user_auth", "user_auth.auth_info_id=auth_info.id").Where("user_auth.u_id=?", user.UId).And("user_auth.status=1").Desc("auth_info.created").Limit(1, 0).Get(&authInfo)
		if hasAuth {
			user.DefaultAuthInfoId = authInfo.Id
			user.DefaultAuthAppid = authInfo.AuthAppid
			base.DBEngine.Table("user").Where("u_id=?", user.UId).Cols("default_auth_info_id").Update(user)
		}
	} else {
		base.DBEngine.Table("auth_info").Where("auth_appid = ?", user.DefaultAuthAppid).Get(&authInfo)
	}

	userShort, err := user.UserToUserShort()
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}
	if userShort == nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}
	authInfoShort, err := authInfo.AuthInfoToAuthInfoShort()
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}

	this.ReturnData = models.UserShortContainer{*userShort, *authInfoShort}
}

// @Title 重置密码（h5使用）
// @Description 登录（h5使用）
// @Param	uId						formData		int64	  		true		"uId"
// @Param	oldPassword				formData		string	  		true		"旧登录密码"
// @Param	newPassword				formData		string	  		true		"新登录密码"
// @Success 200 {string} success
// @router /resetPassword [patch]
func (this *PublicController) ResetPassword() {
	uId := this.MustInt64("uId")
	oldPassword := this.MustString("oldPassword")
	newPassword := this.MustString("newPassword")

	var user models.User
	hasUser, _ := base.DBEngine.Table("user").Where("u_id=?", uId).Get(&user)
	if !hasUser {
		this.ReturnData = util.GenerateAlertMessage(models.UserError300)
		return
	}

	hashedPassword, err := util.EncryptPasswordWithSalt(oldPassword, user.Salt)
	if err != nil {
		util.Logger.Info("err:"+err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}
	if user.Password != hashedPassword {
		this.ReturnData = util.GenerateAlertMessage(models.UserError400)
		return
	}


	newHashedPwd, newSalt, _ := util.EncryptPassword(newPassword)
	user.Password = newHashedPwd
	user.Salt = newSalt
	base.DBEngine.Table("user").Where("u_id=?", user.UId).Cols("password", "salt").Update(&user)

	base.RedisCache.Delete(REDIS_BASEAUTH+user.NickName)

	this.ReturnData = "success"
}

// @Title 获取系统参数component_access_token（h5使用）
// @Description 获取系统参数component_access_token（h5使用）
// @Success 200 {string} success
// @router /getSystemComponentAccessToken [get]
func (this *PublicController) GetSystemComponentAccessToken() {
	var systemConfig models.SystemConfig
	hasComponentAcccessToken, _ := base.DBEngine.Table("system_config").Where("program='component_access_token'").Get(&systemConfig)
	if !hasComponentAcccessToken || systemConfig.ProgramExpireTime < util.UnixOfBeijingTime() {
		this.ReturnData = task.RequestComponentAccessToken()
	} else {
		this.ReturnData = systemConfig.ProgramValue
	}
}

// @Title 获取系统参数pre_auth_code（h5使用）
// @Description 获取系统参数pre_auth_code（h5使用）
// @Success 200 {string} success
// @router /getSystemPreAuthCode [get]
func (this *PublicController) GetSystemPreAuthCode() {
	code := task.RequestPreAuthCode()
	util.Logger.Info("PreAuthCode=" + code)
	this.ReturnData = code
}


// @Title 签名
// @Description 签名
// @Param	url    formData 	string 	 true	 "url"
// @Success 200 {string}  signature,noncestr,timestamp
// @router  /signature [post]
func (this *PublicController) PostSignature() {
	inputUrl := this.MustString("url")

	ticket, err := WcpValidTicket()
	if err != nil {
		this.Err = err
		this.ErrCode = http.StatusInternalServerError
		return
	}

	params := make(map[string]string)
	params["url"] = inputUrl
	params["timestamp"] = strconv.FormatInt(util.UnixOfBeijingTime(), 10)
	params["noncestr"] = strconv.FormatUint(uint64(rand.Uint32()), 10)
	params["jsapi_ticket"] = ticket

	signature := WcpSignature(params)

	this.ReturnData = map[string]string{"signature": signature, "noncestr": params["noncestr"], "timestamp": params["timestamp"]}
}

// @Title 获取认证授权的用户信息
// @Description 获取认证授权的用户信息
// @Param	code    query 	string 	 true	 "code"
// @Success 200 {string}  参考微信文档
// @router  /authed-user-info [get]
func (this *PublicController) GetUserInfo() {
	code := this.MustString("code")
	this.ReturnData, this.Err = WcpAuthedUserInfo(code, 0)
	this.ErrCode = http.StatusInternalServerError
}

// @Title 获取验证码
// @Description  获取验证码
// @Param   areaCode      query    string  false        "区号，中国+86"
// @Param   phoneNum      query    string  true         "phoneNum"
// @Param   bind          query    int     false        "验证手机号是否绑定 0:不验证 1:验证码登录时验证 2:换绑手机号时验证"
// @Success 200 {string} success
// @router /vericode [get]
func (this *PublicController) VeriCode() {
	areaCode := this.GetString("areaCode", "")
	phoneNum := this.MustString("phoneNum")
	bind, _ := this.GetInt("bind", -1)

	//1.验证手机号
	valid := validation.Validation{}
	valid.Required(phoneNum, "phoneNum")
	valid.Numeric(phoneNum, "phoneNum")
	if valid.HasErrors() {
		util.Logger.Info("VeriCode phoneNum err = 手机号有误")
		this.ReturnData = util.GenerateAlertMessage(models.VeriCodeError100)
		return
	}
	var author models.Author
	if bind == 1 {
		//查询该手机号是否已绑定
		has, _ := base.DBEngine.Table("author").Where("phone_number = ?",phoneNum).Get(&author)
		if !has {
			this.ReturnData = util.GenerateAlertMessage(models.LoginError100)
			return
		} else if phoneNum == "18626285773" || phoneNum == "19941898991" || phoneNum == "19941898990" {
			this.ReturnData = "success"
			return
		}
	} else if bind == 2 {
		//查询该手机号是否已绑定
		has, _ := base.DBEngine.Table("author").Where("phone_number = ?",phoneNum).Get(&author)
		if has {
			this.ReturnData = util.GenerateAlertMessage(models.LoginError300)
			return
		}
	}
	//2.发送验证码
	verCodeNum := veriCodeNum()

	format := "【辣课】您本次的验证码为：%s，请勿透露此信息。"
	smsStr := fmt.Sprintf(format, verCodeNum)
	succeed := util.SendSMS(smsStr, areaCode + phoneNum, 1)
	if !succeed {
		util.Logger.Info("VeriCode verCodeNum err = 发送失败")
		this.ReturnData = util.GenerateAlertMessage(models.VeriCodeError200)
		return
	}
	//3.保存至redis
	err := base.RedisCache.Put(VeriCodeRedisKey(phoneNum), verCodeNum, 600*time.Second)
	if err != nil {
		util.Logger.Info("VeriCode redis err = ", err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.VeriCodeError300)
		return

	}
	this.ReturnData = "success"

}


// @Title 获取VIP配置参数
// @Description 获取VIP配置参数
// @Param   program      query    string   true        "参数名称,多个参数用,隔开"
// @Success 200 {object} models.SystemConfigListContainer
// @router /getSystemConfig [get]
func (this *PublicController) GetSystemConfig() {

	program := this.MustString("program")
	var programs = strings.Split(program, ",")
	var systemConfig []models.SystemConfig
	err := base.DBEngine.Table("system_config").In("program", programs).Find(&systemConfig)

	if err != nil {
		util.Logger.Info("GetSystemConfig  err = ", err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}
	if systemConfig == nil {
		systemConfig = make([]models.SystemConfig, 0)
	}
	this.ReturnData = models.SystemConfigListContainer{systemConfig}


}

// @Title	新增app配置信息（后台使用）
// @Description 新增app配置信息（后台使用）
// @Param	config 			body			models.AppConfig			true			"配置信息"
// @Success	200 {object} models.AppConfig
// @router /app-config [post]
func (this *PublicController) PostAppConfig() {
	appConfig := models.AppConfig{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &appConfig)
	if err != nil {
		this.Err = err
		this.ErrCode = http.StatusInternalServerError
		return
	}

	base.DBEngine.Table("app_config").InsertOne(&appConfig)

	this.ReturnData = appConfig
}

// @Title	修改app配置信息（后台使用）
// @Description 修改app配置信息（后台使用）
// @Param	id 	                path	        	int64			true		"配置信息id"
// @Param	version 			formData			string			false		"版本号"
// @Param	description 		formData			string			false		"描述"
// @Param	forcedUpdate     	formData			int				false		"是否强制更新"
// @Success	200 {object} models.AppConfig
// @router /app-config/:id [patch]
func (this *PublicController) PatchAppConfig() {

	id := this.MustInt64(":id")
	description := this.GetString("description", "")
	version := this.GetString("version", "")
	forcedUpdate, _ := this.GetInt("forcedUpdate", -1)

	if forcedUpdate > 1 || forcedUpdate < -1 {
		util.Logger.Info("forcedUpdate 必须为0或者1")
		this.ReturnData = util.GenerateAlertMessage(models.AppConfigError100)
		return
	}

	appConfig := models.AppConfig{}
	hasAppConfig, _ := base.DBEngine.Table("app_config").Where("id=?", id).Get(&appConfig)

	if !hasAppConfig {
		util.Logger.Info("配置信息不存在")
		this.ReturnData = util.GenerateAlertMessage(models.AppConfigError200)
		return
	}

	if len(description) > 0 {
		appConfig.Description = description
	}

	if len(version) > 0 {
		appConfig.Version = version
	}

	if forcedUpdate >= 0 {
		appConfig.ForcedUpdate = forcedUpdate
	}

	_, this.Err = base.DBEngine.Table("app_config").Where("id=?", id).AllCols().Update(&appConfig)

	this.ReturnData = "success"
}

// @Title	获取app配置信息
// @Description 获取app配置信息
// @Param	system 				query				int				true		"1 android 2 ios"
// @Param	version 			query				string			true		"版本号"
// @Success	200	{object} models.AppConfig
// @router /app-config [get]
func (this *PublicController) GetAppConfig() {
	system := this.MustInt("system")
	version := this.MustString("version")

	appConfig := models.AppConfig{}
	hasAppConfig, _ := base.DBEngine.Table("app_config").Where("system = ?", system).And("version = ?", version).Get(&appConfig)
	if !hasAppConfig {
		util.Logger.Info("配置信息不存在")
		this.ReturnData = util.GenerateAlertMessage(models.AppConfigError200)
		return
	}

	this.ReturnData = appConfig
}

// @Title 获取app配置信息列表(后台使用)
// @Description 获取app配置信息列表(后台使用)
// @Param	pageNum		query 	 int			true		"page num start from 1"
// @Param	pageTime	query 	 int64	        true		"page time should be empty when pagenum == 1"
// @Param	pageSize	query 	 int			false		"page size default is 15"
// @Success 200 {object} models.AppConfigListContainer
// @router /app-config/all [get]
func (this *PublicController) GetAll() {
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", 0)
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	var total int64
	var totalErr error
	total, totalErr = base.DBEngine.Table("app_config").Where("created < ?", pageTime).Count(new(models.AppConfig))

	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	var clist []models.AppConfig

	if total > 0 {
		base.DBEngine.Table("app_config").Where("created < ?", pageTime).Limit(pageSize, pageSize*(pageNum-1)).Desc("created").Find(&clist)
	}

	if clist == nil {
		clist = make([]models.AppConfig, 0)
	}

	this.ReturnData = models.AppConfigListContainer{models.BaseListContainer{total, pageNum, pageTime}, clist}
}

// @Title 上传版本（后台使用）
// @Description 上传最新的APK文件（IOS暂时提供位置，不做调用）（后台使用）
// @Param	deviceOs			formData		string		true		"deviceOS"
// @Param	version				formData		string		true		"版本号"
// @Param	description			formData		string		true		"描述"
// @Param	filePath			formData		string		true		"文件地址"
// @Param	channel				formData		string		true		"渠道"
// @Success	200		{string}	success
// @router /version [post]
func (this *PublicController) PostVersion() {
	deviceOS := this.MustString("deviceOs")
	versionStr := this.MustString("version")
	description := this.MustString("description")
	filePath := this.MustString("filePath")
	channel := this.MustString("channel")

	version := models.Version{}
	version.Version = versionStr
	version.Description = description
	version.FilePath = filePath
	version.DeviceOs = deviceOS
	version.Channel = channel

	base.DBEngine.Table("version").InsertOne(&version)

	this.ReturnData = "success"
}

// @Title	获取历史版本号列表（后台使用）
// @Description 获取历史版本号，以创建时间降序排序（后台使用）
// @Param	pageNum			query 			int			true		"page num start from 1"
// @Param	pageTime		query 			int64		true		"page time should be empty when pagenum == 1"
// @Param	pageSize		query 			int			false		"page size default is 15"
// @Success	200		{object}	models.VersionContainer
// @router /version/history [get]
func (this *PublicController) GetVersionHistory() {
	//分页
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", 0)
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	var total int64
	var totalErr error

	queryString := "select count(*) from version"
	queryString += " where (deleted_at is null or deleted_at=0)"

	total, totalErr = base.DBEngine.SQL(queryString).Count(new(models.Version))

	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	var versionList []models.Version
	if total > 0 {
		base.DBEngine.Table("version").Limit(pageSize, pageSize*(pageNum-1)).Desc("created").Find(&versionList)
	}

	if versionList == nil {
		versionList = make([]models.Version, 0)
	}

	this.ReturnData = models.VersionContainer{models.BaseListContainer{total, pageNum, pageTime}, versionList}
}

// @Title	获取最新的版本号
// @Description 获取当前最新的版本号
// @Param	deviceOs			path			string			true	"设备操作系统"
// @Param	channel				query			string			false	"渠道"
// @Success	200	{models.version}
// @router /version/:deviceOs [get]
func (this *PublicController) GetVersion() {
	deviceOs := this.MustString(":deviceOs")
	channel := this.GetString("channel")

	var versions []models.Version
	if channel != "" {
		base.DBEngine.Where("device_os = ?", deviceOs).And("channel = ?", channel).Desc("created").Limit(1, 0).Find(&versions)
	} else {
		base.DBEngine.Where("device_os = ?", deviceOs).Desc("created").Limit(1, 0).Find(&versions)
	}


	if versions == nil || len(versions) == 0 {
		util.Logger.Info("未找到最新版本")
		this.ReturnData = util.GenerateAlertMessage(models.VersionError100)
		return
	}

	serverUrl := base.ServerURL + "/v1/file" + versions[0].FilePath

	versions[0].FilePath = serverUrl
	this.ReturnData = versions[0]
}

// @Title	修改版本描述
// @Description 修改版本描述，
// @Param	versionId			path		string		true		"修改的版本编号"
// @Param	description			formData	string		true		"修改的版本描述内容"
// @Success	200		{string}	success
// @router /version/:versionId [patch]
func (this *PublicController) PatchVersion() {
	versionId := this.MustString(":versionId")
	description := this.MustString("description")

	version := models.Version{}
	base.DBEngine.Table("version").Where("version_id=?", versionId).Get(&version)

	version.Description = description

	base.DBEngine.Table("version").Where("version_id=?", versionId).Cols("description").Update(&version)

	this.ReturnData = "success"
}

// @Title	删除该版本
// @Description 删除该版本
// @Param	versionId	path		string	true	"要删除的版本编号"
// @Success	200		{string}	success
// @router /version/:versionId [delete]
func (this *PublicController) DeleteVersion() {
	versionId := this.MustString(":versionId")

	base.DBEngine.Table("version").Where("version_id=?", versionId).Delete(new(models.Version))


	this.ReturnData = "success"
}


//---------------------方法-----------------------------

func WcpAuthedUserInfo(code string, codeType int) (userInfo map[string]interface{}, err error) {
	redisKey := REDIS_GZHCODE + code
	accessToken := ""
	openid := ""
	var appid string
	var appSecret string
	if codeType == 0 {
		appid = models.JTGZHAppId
		appSecret = models.JTGZHAppSecret
	} else {
		appid = models.JTAPPAppId
		appSecret = models.JTAPPAppSecret
	}
	if base.RedisCache.IsExist(redisKey) {
		value := base.RedisCache.Get(redisKey)
		result := string(value.([]byte))
		accessToken = strings.Split(result, "###")[0]
		openid = strings.Split(result, "###")[1]
	} else {
		accessToken, openid, err = wcpGetAuthToken(code, appid, appSecret)
		if err != nil {
			return
		}
		//防止code被用两次，查询一次之后放入redis
		base.RedisCache.Put(redisKey, accessToken + "###" + openid, 60 * 60 * 24 * time.Second)
	}

	userInfoBytes, err := wcpRequestAuthUserInfo(accessToken, openid)

	userInfo = make(map[string]interface{})
	err = json.Unmarshal(userInfoBytes, &userInfo)
	if err != nil {
		return
	}

	userInfo["accessToken"] = accessToken

	return userInfo, nil
}

func wcpRequestAuthUserInfo(accessToken, openid string) (userInfo []byte, err error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", "api.weixin.qq.com", accessToken, openid))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response := models.WcpCommonReturnData{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	if response.Errcode != 0 {
		err = errors.New(response.Errmsg)
		return
	}

	return body, nil
}

func wcpGetAuthToken(code string, appId string, appSecret string) (accessToken, openid string, err error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", "api.weixin.qq.com", appId, appSecret, code))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response := models.WcpAuthTokenData{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}

	if response.Errcode != 0 {
		err = errors.New(response.Errmsg)
		util.Logger.Info("wcpGetAuthToken err = " + err.Error())
		return
	}

	return response.AccessToken, response.Openid, nil
}

func WcpValidTicket() (ticket string, err error) {
	systemConfig := models.SystemConfig{}
	has := false
	has, err = base.DBEngine.Table("system_config").Where("program = ?", "js_ticket").Get(&systemConfig)
	if err != nil {
		return
	}

	if has {
		if util.UnixOfBeijingTime() < (systemConfig.ProgramExpireTime - 60*30) {
			return systemConfig.ProgramValue, nil
		}
	}
	ticket, err, errCode := wcpRefreshJsapiTicket()
	if errCode == 40001 || errCode == 40014 || errCode == 42001 {
		err = nil
		_, err1 := wcpRefreshPlatformAccessToken()
		if err1 != nil {
			err = err1
			util.Logger.Info("err:"+err.Error())
			util.Logger.Info("errCode:"+strconv.Itoa(errCode))
			return
		}
		ticket, err, _ = wcpRefreshJsapiTicket()
	}
	return
}

func wcpRefreshJsapiTicket() (ticket string, err error, wcpErrCode int) {
	accessToken, err := wcpValidPlatformAccessToken()
	if err != nil {
		return
	}

	resp, err := http.Get(fmt.Sprintf("https://%s/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", "api.weixin.qq.com", accessToken))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("http请求错误，http code = %d", resp.StatusCode)), 0
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response := models.WcpRefreshTicketData{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	util.Logger.Info(response)

	if response.Errcode != 0 {
		err = errors.New(fmt.Sprintf("%s %d", response.Errmsg, response.Errcode))
		wcpErrCode = response.Errcode
		return
	}

	systemConfig := models.SystemConfig{}
	var has bool
	has, err = base.DBEngine.Table("system_config").Where("program = ?", "js_ticket").Get(&systemConfig)
	if err != nil {
		return
	}
	if has {
		systemConfig.ProgramValue = response.Ticket
		systemConfig.ProgramExpireTime = util.UnixOfBeijingTime() + int64(response.ExpiresIn)
		base.DBEngine.Table("system_config").Where("program = ?", "js_ticket").AllCols().Update(&systemConfig)
	} else {
		systemConfig.Program = "js_ticket"
		systemConfig.ProgramValue = response.Ticket
		systemConfig.ProgramExpireTime = util.UnixOfBeijingTime() + int64(response.ExpiresIn)
		base.DBEngine.Table("system_config").InsertOne(&systemConfig)
	}


	return response.Ticket, nil, 0
}

func wcpValidPlatformAccessToken() (accessToken string, err error) {
	systemConfig := models.SystemConfig{}
	has := false
	has, err = base.DBEngine.Table("system_config").Where("program = ?", "js_access_token").Get(&systemConfig)
	if err != nil {
		return
	}

	if has {
		if util.UnixOfBeijingTime() < (systemConfig.ProgramExpireTime - 60*60) {
			return systemConfig.ProgramValue, nil
		}
	}
	accessToken, err = wcpRefreshPlatformAccessToken()
	return
}

func wcpRefreshPlatformAccessToken() (accessToken string, err error) {
	var resp *http.Response
	resp, err = http.Get(fmt.Sprintf("https://%s/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", "api.weixin.qq.com", models.JTGZHAppId, models.JTGZHAppSecret))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("http请求错误，http code = %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response := models.WcpRefreshTokenData{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	util.Logger.Info(response)

	if response.Errcode != 0 {
		err = errors.New(response.Errmsg)
		return
	}

	systemConfig := models.SystemConfig{}
	var has bool
	has, err = base.DBEngine.Table("system_config").Where("program = ?", "js_access_token").Get(&systemConfig)
	if err != nil {
		return
	}
	if has {
		systemConfig.ProgramValue = response.AccessToken
		systemConfig.ProgramExpireTime = util.UnixOfBeijingTime() + int64(response.ExpiresIn)
		base.DBEngine.Table("system_config").Where("program = ?", "js_access_token").AllCols().Update(&systemConfig)
	} else {
		systemConfig.Program = "js_access_token"
		systemConfig.ProgramValue = response.AccessToken
		systemConfig.ProgramExpireTime = util.UnixOfBeijingTime() + int64(response.ExpiresIn)
		base.DBEngine.Table("system_config").InsertOne(&systemConfig)
	}

	return response.AccessToken, nil
}

func WcpSignature(intput map[string]string) string {
	keys := make([]string, len(intput))

	i := 0
	for k := range intput {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	strTemp := ""
	for _, key := range keys {
		strTemp = strTemp + key + "=" + intput[key] + "&"
	}
	strTemp = strTemp[:len(strTemp)-1]

	hasher := sha1.New()
	hasher.Write([]byte(strTemp))
	sha1Str := hex.EncodeToString(hasher.Sum(nil))

	return sha1Str
}
//生成验证码随机数
func veriCodeNum() string {
	rnd := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return vcode
}

//生成验证码存储在redis中的key
func VeriCodeRedisKey(phoneNum string) string {
	return "veriCode" + phoneNum
}
