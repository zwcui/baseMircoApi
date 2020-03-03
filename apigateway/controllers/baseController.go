package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/cache/redis"
	"jingting_server/apigateway/util"
	"fmt"
	"errors"
	"regexp"
	"strings"
	"time"
	"jingting_server/apigateway/base"
	"jingting_server/apigateway/models"
	"strconv"
	"encoding/json"
)

//数据返回结构体
type BaseController struct {
	beego.Controller
	Err        error
	ErrCode    int
	ReturnData interface{}
	TestData   []byte
	IsDirectReturn int	//是否直接返回ReturnData，不需要code与description
	IsXml	   int	//是否xml格式
}

//api接口统一controller
type apiController struct {
	BaseController
	NeedUserAuthList []RequestPathAndMethod
	NeedAuthorAuthList []RequestPathAndMethod
}

//需要验证的请求路径
type RequestPathAndMethod struct {
	PathRegexp 		string
	Method     		string
	RoleType     	[]int
}

const (
	REDIS_BASEAUTH = "BaseAuth_"
	REDIS_BASEAUTH_OPENID = "BaseAuth_Openid_"
	REDIS_BASEAUTH_UNIONID = "BaseAuth_Unionid_"
	REDIS_GZHCODE = "BaseAuth_GZHCode_"
	REDIS_WEBCODE = "BaseAuth_WebCode_"
)

const AUTHED_UID_KEY = "AUTHED_UID_KEY"

//接口调用返回code
const (
	ServerApiSuccuess 		= 1000		//调用成功
	ServerApiUndefinedFail 	= 999		//未知错误
	ServerApiIllegalParam 	= 900		//接口参数不合法

)

//返回数据的head
const (
	headerCodeKey 	= "code"
	headerDesKey 	= "description"
)

const (
	DEFAULT_PAGESIZE = 15	//默认分页15条
)


//默认请求之前加路径head身份验证，默认所有方法都需要验证，各个api可以重写该方法
func (this *apiController) Prepare(){
	util.Logger.Info("apiController Prepare")
	this.NeedUserAuthList = []RequestPathAndMethod{{".+", "post", []int{0}}, {".+", "patch", []int{0}}, {".+", "delete", []int{0}}, {".+", "put", []int{0}}}
	this.userAuth()
}

//对路径进行校验
//公众号管理系统用户体系
func (this *apiController) userAuth(){
	pathNeedAuth := false
	var roleType []int
	for _, value := range this.NeedUserAuthList {
		if ok, _ := regexp.MatchString(value.PathRegexp + "$", this.Ctx.Request.URL.Path); ok && strings.ToUpper(this.Ctx.Request.Method) == strings.ToUpper(value.Method) {
			pathNeedAuth = true
			roleType = value.RoleType
			break
		}
	}

	//要求head中放Authorization，内容格式为 "Basic nickName:123456"  密码为加密后的密文，加密方式为base64
	if pathNeedAuth {
		nickName, encryptedPassword, ok := this.Ctx.Request.BasicAuth()
		if !ok {
			w := this.Ctx.ResponseWriter
			w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+"empty auth"+`"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			this.ServeJSON()
			this.StopRun()
		}
		redisTemp := base.RedisCache.Get(REDIS_BASEAUTH + nickName)
		if redisTemp == nil {
			user, err := UserWithNickName(nickName)
			if err != nil {
				util.Logger.Info(err.Error())
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + err.Error() + `"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
			if user == nil {
				util.Logger.Info("user not exists " + nickName)
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + nickName + " user not exists "+`"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
			//校验密码
			passwordByte, _ := util.Base64Decode([]byte(encryptedPassword))
			password := string(passwordByte)
			hashedPwd, _ := util.EncryptPasswordWithSalt(password, user.Salt)
			if hashedPwd != user.Password {
				util.Logger.Info("password incorrect")
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + "password error" + `"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			} else {
				hasAuth := false
				if roleType == nil || len(roleType) == 0 {
					hasAuth = true
				} else if roleType[0] == 0 && len(roleType) == 1 {
					hasAuth = true
				} else {
					for _, role := range roleType {
						if role == 0 {
							continue
						}
						if strings.Contains(user.RoleType, "["+strconv.Itoa(role)+"]") {
							hasAuth = true
							break
						}
					}
				}
				if !hasAuth {
					util.Logger.Info("do not have right auth")
					w := this.Ctx.ResponseWriter
					w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + "do not have right auth" + `"`)
					w.WriteHeader(401)
					w.Write([]byte("401 Unauthorized\n"))
					this.ServeJSON()
					this.StopRun()
				}
			}

			//存入redis
			if user != nil {
				userRedis := models.UserRedis{encryptedPassword, user.RoleType}
				userRedisBytes, _ := json.Marshal(userRedis)
				base.RedisCache.Put(REDIS_BASEAUTH+nickName, string(userRedisBytes), 60*60*2*time.Second)
			}
		} else {
			var userRedis models.UserRedis
			err := json.Unmarshal(redisTemp.([]byte), &userRedis)
			if err != nil {
				util.Logger.Info(err.Error())
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + err.Error() + `"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
			if encryptedPassword != userRedis.EncryptedPassword {
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + "password redis error" + `"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			} else {
				hasAuth := false
				if roleType == nil || len(roleType) == 0 {
					hasAuth = true
				} else if roleType[0] == 0 && len(roleType) == 1 {
					hasAuth = true
				} else {
					for _, role := range roleType {
						if role == 0 {
							continue
						}
						if strings.Contains(userRedis.RoleType, "["+strconv.Itoa(role)+"]") {
							hasAuth = true
							break
						}
					}
				}
				if !hasAuth {
					util.Logger.Info("do not have right auth")
					w := this.Ctx.ResponseWriter
					w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + "do not have right auth" + `"`)
					w.WriteHeader(401)
					w.Write([]byte("401 Unauthorized\n"))
					this.ServeJSON()
					this.StopRun()
				}
			}
		}

	}

	//保存后台操作记录
	nickName, _, _ := this.Ctx.Request.BasicAuth()
	authInfoIdStr := this.Ctx.Request.Header.Get("authInfoId")
	if authInfoIdStr == "" {
		return
	}

	var operationRecord models.OperationRecord
	if authInfoIdStr != "" {
		authInfoId, _ := strconv.ParseInt(authInfoIdStr, 10, 64)
		operationRecord.AuthInfoId = authInfoId
	}
	operationRecord.Operator = nickName
	operationRecord.RequestUrl = this.Ctx.Request.RequestURI
	operationRecord.RequestMethod = this.Ctx.Request.Method
	operationRecord.RequestRemoteAddr = this.Ctx.Request.RemoteAddr
	base.DBEngine.Table("operation_record").InsertOne(&operationRecord)
}

//对路径进行校验
//视频上传分享用户体系
func (this *apiController) authorAuth(){
	pathNeedAuth := false
	for _, value := range this.NeedAuthorAuthList {
		if ok, _ := regexp.MatchString(value.PathRegexp + "$", this.Ctx.Request.URL.Path); ok && strings.ToUpper(this.Ctx.Request.Method) == strings.ToUpper(value.Method) {
			pathNeedAuth = true
			break
		}
	}
	if beego.BConfig.RunMode != "prod" {
		pathNeedAuth = false
	}

	//要求head中放Authorization，内容格式为 "Basic openid:unionid"  密码为加密后的密文，加密方式为base64
	if pathNeedAuth {
		unionid, encryptedOpenid, ok := this.Ctx.Request.BasicAuth()
		if !ok {
			w := this.Ctx.ResponseWriter
			w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+"empty auth"+`"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			this.ServeJSON()
			this.StopRun()
		}

		redisTemp := base.RedisCache.Get(REDIS_BASEAUTH_UNIONID + unionid)
		if redisTemp == nil {
			author, err := AuthorWithUnionid(unionid)
			if err != nil {
				util.Logger.Info(err.Error())
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + err.Error() + `"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
			if author == nil {
				util.Logger.Info("author not exists " + unionid)
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + unionid + " author not exists "+`"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
			//校验密码
			openidByte, _ := util.Base64Decode([]byte(encryptedOpenid))
			openid := string(openidByte)
			if author.Openid != openid && author.WebOpenid != openid && author.AppOpenid != openid {
				util.Logger.Info("unionid incorrect")
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+"!strings.Contains(author.Openid, openid)"+`"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}

			//存入redis
			if author != nil {
				authorRedis := models.AuthorRedis{author.Id, author.Openid, author.WebOpenid, author.AppOpenid,author.Unionid}
				authorRedisBytes, _ := json.Marshal(authorRedis)
				base.RedisCache.Put(REDIS_BASEAUTH_UNIONID + unionid, string(authorRedisBytes), 60*60*2*time.Second)
			}
		} else {
			var authorRedis models.AuthorRedis
			err := json.Unmarshal(redisTemp.([]byte), &authorRedis)
			if err != nil {
				util.Logger.Info(err.Error())
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + err.Error() + `"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
			openidByte, _ := util.Base64Decode([]byte(encryptedOpenid))
			openid := string(openidByte)
			if authorRedis.Openid != openid && authorRedis.WebOpenid != openid && authorRedis.AppOpenid != openid {
				w := this.Ctx.ResponseWriter
				w.Header().Set("WWW-Authenticate", `Base Auth failed : "` + "!strings.Contains(authorRedis.Openid, openid)" + `"`)
				w.WriteHeader(401)
				w.Write([]byte("401 Unauthorized\n"))
				this.ServeJSON()
				this.StopRun()
			}
		}
	}
}

//取参数错误返回
func (this *BaseController) Failed() {
	util.Logger.Info(this.Ctx.Request.URL.Path)
	util.Logger.Info("BaseController Failed")
	if this.ErrCode == 0 {
		this.ErrCode = ServerApiUndefinedFail
	}
	this.Data["json"] = map[string]interface{}{
		"header": map[string]string{
			headerCodeKey: fmt.Sprintf("%d", this.ErrCode),
			headerDesKey:  this.Err.Error(),
		},
	}
	this.ServeJSON()
	this.StopRun()
}

// 函数结束时,组装成json结果返回
func (this *BaseController) Finish() {
	util.Logger.Info(this.Ctx.Request.URL.Path)
	if this.Err != nil {
		this.Failed()
	}
	if this.IsDirectReturn == 1 {
		if this.IsXml == 1 {
			util.Logger.Info("xml test")
			//this.Data["xml"] = this.ReturnData
			//this.ServeXML()
			util.Logger.Info(string(this.TestData))
			this.Ctx.ResponseWriter.Write(this.TestData)
		} else {
			this.Ctx.ResponseWriter.Write([]byte(this.ReturnData.(string)))
		}
	} else {
		r := struct {
			Header interface{} `json:"header"`
			Data   interface{} `json:"data"`
		}{}

		r.Header = map[string]string{
			headerCodeKey: fmt.Sprintf("%d", ServerApiSuccuess),
			headerDesKey:  "success",
		}

		r.Data = this.ReturnData
		this.Data["json"] = r
		this.ServeJSON()
	}
}

// 如果请求的参数不存在,就直接 error返回
func (this *BaseController) MustString(key string) string {
	v := this.GetString(key)
	if v == "" {
		this.ErrCode = ServerApiIllegalParam
		this.Err = errors.New(fmt.Sprintf("require filed: %s", key))
		this.Failed()
	}
	return v
}

// 如果请求的参数不存在,就直接 error返回
func (this *BaseController) MustInt64(key string) int64 {
	v, err := this.GetInt64(key)
	if err != nil {
		this.ErrCode = ServerApiIllegalParam
		this.Err = errors.New(fmt.Sprintf("require filed: %s", key))
		this.Failed()
	}
	return v
}

// 如果请求的参数不存在,就直接 error返回
func (this *BaseController) MustFloat64(key string) float64 {
	v, err := this.GetFloat(key)
	if err != nil {
		this.ErrCode = ServerApiIllegalParam
		this.Err = errors.New(fmt.Sprintf("require filed: %s", key))
		this.Failed()
	}
	return v
}

// 如果请求的参数不存在,就直接 error返回
func (this *BaseController) MustInt(key string) int {
	v, err := this.GetInt(key)
	if err != nil {
		this.ErrCode = ServerApiIllegalParam
		this.Err = errors.New(fmt.Sprintf("require filed: %s", key))
		this.Failed()
	}
	return v
}

func (this *BaseController) GetPageSize(key string) int {
	v, _ := this.GetInt(key)
	if v == 0 {
		return DEFAULT_PAGESIZE
	}
	return v
}

//检查传入的用户和通过基本验证的用户是否为同一人
func (this *BaseController) NeedSameAuthor(authorId int64) {
	if beego.BConfig.RunMode != "prod" {
		return
	}

	unionid, _, ok := this.Ctx.Request.BasicAuth()
	if !ok {
		w := this.Ctx.ResponseWriter
		w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+"empty auth"+`"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		this.ServeJSON()
		this.StopRun()
	}

	redisTemp := base.RedisCache.Get(REDIS_BASEAUTH_UNIONID + unionid)

	if redisTemp == nil {
		author, err := AuthorWithUnionid(unionid)
		if err != nil {
			util.Logger.Info(err.Error())
			w := this.Ctx.ResponseWriter
			w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+err.Error()+`"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			this.ServeJSON()
			this.StopRun()
		}
		if author == nil {
			util.Logger.Info("author not exists ")
			w := this.Ctx.ResponseWriter
			w.Header().Set("WWW-Authenticate", `Base Auth failed :  author not exists "`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			this.ServeJSON()
			this.StopRun()
		}

		//存入redis
		//openidByte, _ := util.Base64Decode([]byte(encryptedOpenid))
		//openid := string(openidByte)
		if author != nil {
			authorRedis := models.AuthorRedis{author.Id, author.Openid, author.WebOpenid, author.AppOpenid,author.Unionid}
			authorRedisBytes, _ := json.Marshal(authorRedis)
			base.RedisCache.Put(REDIS_BASEAUTH_UNIONID + unionid, string(authorRedisBytes), 60*60*2*time.Second)
		}

		if author.Id != authorId {
			w := this.Ctx.ResponseWriter
			w.Header().Set("WWW-Authenticate", `Basic realm="access user`+strconv.FormatInt(authorId, 10)+`is forbidden for you"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized  NeedSameAuthor!!!!!(author.Id=" + strconv.FormatInt(author.Id, 10) + ", authorId=" + strconv.FormatInt(authorId, 10) + ")\n"))
			util.Logger.Debug("401 Unauthorized  NeedSameAuthor!!!!!(author.Id=" + strconv.FormatInt(author.Id, 10) + ", authorId=" + strconv.FormatInt(authorId, 10) + ")\n")
			this.ServeJSON()
			this.StopRun()
		}

	} else {
		var authorRedis models.AuthorRedis
		err := json.Unmarshal(redisTemp.([]byte), &authorRedis)
		if err != nil {
			util.Logger.Info(err.Error())
			w := this.Ctx.ResponseWriter
			w.Header().Set("WWW-Authenticate", `Base Auth failed : "`+err.Error()+`"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			this.ServeJSON()
			this.StopRun()
		}

		if authorRedis.AuthorId != authorId {
			w := this.Ctx.ResponseWriter
			w.Header().Set("WWW-Authenticate", `Basic realm="access user`+strconv.FormatInt(authorId, 10)+`is forbidden for you"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized  NeedSameAuthor!!!!!(author.Id=" + strconv.FormatInt(authorRedis.AuthorId, 10) + ", authorId=" + strconv.FormatInt(authorId, 10) + ")\n"))
			util.Logger.Debug("401 Unauthorized  NeedSameAuthor!!!!!(author.Id=" + strconv.FormatInt(authorRedis.AuthorId, 10) + ", authorId=" + strconv.FormatInt(authorId, 10) + ")\n")
			this.ServeJSON()
			this.StopRun()
		}
	}
}

//通过昵称获取用户信息
func UserWithNickName(nickName string) (user *models.User, err error) {
	var u models.User
	hasUser, err := base.DBEngine.Table("user").Where("nick_name=?", nickName).Get(&u)
	if err != nil {
		return nil, err
	}
	if !hasUser {
		return nil, nil
	}
	return &u, nil
}

//通过unionid获取用户信息
func AuthorWithUnionid(unionid string) (author *models.Author, err error) {
	var a models.Author
	hasAuthor, err := base.DBEngine.Table("author").Where("unionid=?", unionid).Get(&a)
	if err != nil {
		return nil, err
	}
	if !hasAuthor {
		return nil, nil
	}
	return &a, nil
}