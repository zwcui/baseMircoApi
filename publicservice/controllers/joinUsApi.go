/*
@Time : 2019/9/16 下午6:28 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/base"
	"jingting_server/publicservice/util"
	"strconv"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"math/rand"
	"github.com/astaxie/beego"
)

type JoinUsController struct {
	apiController
}

func (this *JoinUsController) Prepare(){
	this.NeedAuthorAuthList = []RequestPathAndMethod{
	}
	this.authorAuth()
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/addJoinUsColumn", "post", []int{}},
		{"/deleteJoinUsColumn", "delete", []int{}},
	}
	this.userAuth()
}

// @Title 加盟开课设置栏目
// @Description 加盟开课设置栏目
// @Param	uId			    	formData			int64		  	true		"用户id"
// @Param	classify			formData			string	  		true		"分类"
// @Param	title				formData			string	  		true		"课程名称"
// @Param	cover				formData			string	  		true		"封面"
// @Param	videoUrl	    	formData			string 		  	true		"介绍视频地址"
// @Param	showType        	formData			int			  	true		"展示类型，1为仅代理可看（伯乐用户 / 代理商），2为VIP会员可看，3为所有用户可看"
// @Param	totalDuration	    formData		    int		  		true		"总时长，单位秒"
// @Success 200 {string} success
// @router /addJoinUsColumn [post]
func (this *JoinUsController) AddJoinUsColumn() {
	uId := this.MustInt64("uId")
	classify := this.MustString("classify")
	title := this.MustString("title")
	cover := this.MustString("cover")
	videoUrl := this.MustString("videoUrl")
	showType := this.MustInt("showType")
	totalDuration := this.MustInt("totalDuration")

	var joinUs models.JoinUs
	joinUs.UId = uId
	joinUs.Classify = classify
	joinUs.Title = title
	joinUs.Cover = cover
	joinUs.VideoUrl = videoUrl
	joinUs.Url = videoUrl
	joinUs.ShowType = showType
	joinUs.TotalDuration = totalDuration
	base.DBEngine.Table("join_us").InsertOne(&joinUs)

	this.ReturnData = "success"
}

// @Title 加盟开课栏目列表
// @Description 加盟开课栏目列表
// @Param	classify				query				string	  		false		"分类"
// @Param	showType        		query				int			  	false		"展示类型，0所有，1为仅代理可看（伯乐用户 / 代理商），2为VIP会员可看，3为所有用户可看"
// @Param	pageNum					query 	  			int				true		"page num start from 1"
// @Param	pageTime				query 	  			int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize				query 	  			int				false		"page size default is 15"
// @Success 200 {object} models.JoinUsListContainer
// @router /getJoinUsColumnList [get]
func (this *JoinUsController) GetJoinUsColumnList() {
	classify := this.GetString("classify", "")
	showType, _ := this.GetInt("showType", 0)
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}
	var joinUsList []models.JoinUs

	totalSql := " select count(*) from join_us where deleted_at is null "
	dataSql := " select * from join_us where deleted_at is null "

	if classify != "" {
		totalSql += " and classify = '" +classify + "'"
		dataSql += " and classify = '" +classify + "'"
	}
	if showType != 0 {
		totalSql += " and show_type = '" + strconv.Itoa(showType) + "'"
		dataSql += " and show_type = '" + strconv.Itoa(showType) + "'"
	}

	total, totalErr := base.DBEngine.SQL(totalSql).Count(new(models.JoinUs))
	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}
	dataSql += " order by join_us.created desc limit "+strconv.Itoa(pageSize*(pageNum-1))+", "+strconv.Itoa(pageSize)

	if total > 0 {
		base.DBEngine.SQL(dataSql).Find(&joinUsList)
	}
	var videoList []models.JoinUs
	if joinUsList == nil {
		joinUsList = make([]models.JoinUs, 0)
	}
	for _, video := range joinUsList {
		video.VideoUrl = generateUrlWithKey(video.VideoUrl, 0,2 * 60 * 60 ,0)
		video.Url = generateUrlWithKey(video.VideoUrl, 0,2 * 60 * 60 ,0)
		videoList = append(videoList, video)
	}
	if videoList == nil {
		videoList = make([]models.JoinUs, 0)
	}

	this.ReturnData = models.JoinUsListContainer{models.BaseListContainer{total, pageNum, pageTime}, videoList}
}

// @Title 删除加盟开课栏目
// @Description 删除加盟开课栏目
// @Param	joinUsId        		query				int64		  	true		"栏目id"
// @Success 200 {string} success
// @router /deleteJoinUsColumn [delete]
func (this *JoinUsController) DeleteJoinUsColumn() {
	joinUsId := this.MustInt64("joinUsId")

	base.DBEngine.Table("join_us").Where("id=?", joinUsId).Delete(new(models.JoinUs))

	this.ReturnData = "success"
}



//---------------------------方法-----------------------------
//视频音频链接加防盗链
//启用该功能后，视频原始 URL 将不再能直接播放
//ipNum : 0为不限制ip   duration：视频音频时长，默认过期10倍时长  exper：试看时长
func generateUrlWithKey(sourceUrl string, ipNum int, duration int, exper int) (url string){
	if sourceUrl == "" {
		return ""
	}

	util.Logger.Info("generateUrlWithKey")

	prefix := ""
	if beego.BConfig.RunMode == "prod" {
		prefix = "http://1300249557.vod2.myqcloud.com"
	} else {
		prefix = "http://1259438447.vod2.myqcloud.com"
	}

	dir := strings.Replace(sourceUrl, prefix, "", -1)
	dir = "/" + strings.Split(dir, "/")[1] + "/" + strings.Split(dir, "/")[2] + "/"
	util.Logger.Info(dir)

	//过期时间 4085458149 2039年
	var t string
	if duration == 0 {
		t = strconv.FormatInt(2192087014, 16)
	} else {
		//生成当前时间+指定时长的过期时间戳
		t = strconv.FormatInt(util.TimeParseDuration(duration * 10,true).Unix(), 16)
	}

	//限制ip数
	rlimit := strconv.Itoa(ipNum)

	//随机字符串
	us := strconv.FormatUint(uint64(rand.Uint32()), 10)

	//sign = md5(KEY + Dir + t + exper + rlimit + us)
	unsignStr := models.JTTencentCloudFangdaolianKey + dir + t
	if exper != 0 {
		unsignStr += strconv.Itoa(exper)
	}
	if ipNum !=0 {
		unsignStr += rlimit + us
	} else {
		unsignStr +=  us
	}


	//freeHasher := md5.New()
	//freeHasher.Write([]byte(freeUnsignStr))
	//freeMd5 := hex.EncodeToString(freeHasher.Sum(nil))
	//freeUrl = sourceUrl + "?t=" +  freeT + "&exper=" + strconv.Itoa(freeTime)
	//
	//forwardHasher := md5.New()
	//forwardHasher.Write([]byte(forwardUnsignStr))
	//forwardMd5 := hex.EncodeToString(forwardHasher.Sum(nil))
	//forwardUrl = sourceUrl + "?t=" + forwardT + "&exper=" + strconv.Itoa(forwardTime)

	fullHasher := md5.New()
	fullHasher.Write([]byte(unsignStr))
	fullMd5 := hex.EncodeToString(fullHasher.Sum(nil))
	if exper != 0 {
		url = sourceUrl + "?t=" + t + "&exper=" + strconv.Itoa(exper)
	} else {
		url = sourceUrl + "?t=" + t
	}

	if ipNum != 0 {
		url += "&rlimit=" + rlimit + "&us=" + us + "&sign=" + fullMd5
	} else {
		url +=  "&us=" + us + "&sign=" + fullMd5
	}


	return url
}