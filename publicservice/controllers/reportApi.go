/*
@Time : 2019/9/16 上午11:19 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"strconv"
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/base"
	"jingting_server/publicservice/util"
	"jingting_server/publicservice/remote"
)

type ReportController struct {
	apiController
}

func (this *ReportController) Prepare(){
	this.NeedAuthorAuthList = []RequestPathAndMethod{
		{"/addReport", "post", []int{}},
	}
	this.authorAuth()

	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/getReportList", "get", []int{}},
	}
	this.userAuth()
}

// @Title 举报
// @Description 举报
// @Param	authorId	    formData		int64		  	true		"用户id"
// @Param	type	        formData		int  		  	true		"举报类型 1:视频 2:评论 3:课程"
// @Param	valueId			formData		int64	  		true		"type为1时:视频id, type为2时:评论id, type为3时:课程id"
// @Param	classify	    formData		int  		  	true		"举报内容类别 1:违法违禁, 2:色情, 3:低俗, 4:赌博诈骗, 5:血腥暴力, 6:人生攻击, 7:与其他视频相同, 8:不良封面标题, 9;青少年不良信息, 10:其他 "
// @Param	content	        formData		string		  	true		"举报内容"
// @Success 200 {string} success
// @router /addReport [post]
func (this *ReportController) AddReport() {
	valueId := this.MustInt64("valueId")
	authorId := this.MustInt64("authorId")
	classify := this.MustInt("classify")
	reportType := this.MustInt("type")
	content := this.MustString("content")

	this.NeedSameAuthor(authorId)
	author := models.Author{}

	//验证用户
	has, err := base.DBEngine.Table("author").Where("id = ?", authorId).Get(&author)
	if err != nil {
		util.Logger.Info("AddReport get Author err = ", err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
		return
	}
	if !has {
		this.ReturnData = util.GenerateAlertMessage(models.UserError300)
		return
	}

	var report models.Report
	report.ValueId = valueId
	report.AuthorId = authorId
	report.Type = reportType
	report.Classify = classify
	report.Content = content

	base.DBEngine.Table("report").InsertOne(&report)

	//举报成功后给app用户返回通知
	params := make(map[string]string)
	if reportType == 1 {
		params = map[string]string{
			"videoId" : strconv.FormatInt(valueId,10),
			"receiverId" : strconv.FormatInt(authorId,10),
		}
	} else if reportType == 2 {
		params = map[string]string{
			"commentId" : strconv.FormatInt(valueId,10),
			"receiverId" : strconv.FormatInt(authorId,10),
		}
	} else if reportType == 3 {
		params = map[string]string{
			"courseId" : strconv.FormatInt(valueId,10),
			"receiverId" : strconv.FormatInt(authorId,10),
		}
	}

	appMessage := models.AppMessage{}
	appMessage.ReceiverId = authorId
	appMessage.ActionUrl = remote.JumpUrlWithKeyAndPramas(models.JTREPORT_JUMP_KEY, params)
	appMessage.Content = "您的举报信息已收到，我们会尽快核实处理。"
	appMessage.Type = 4
	_, pushErr := remote.PushMessageToUser(authorId, &appMessage, "", 0)
	if pushErr != nil {
		util.Logger.Info("addReport pushErr = ", pushErr.Error())
	}
	base.DBEngine.Table("app_message").InsertOne(&appMessage)

	this.ReturnData = "success"
}

// @Title 举报列表
// @Description 举报列表
// @Param	pageNum						query 	  			int				true		"page num start from 1"
// @Param	pageTime					query 	  			int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize					query 	  			int				false		"page size default is 15"
// @Success 200 {object} models.ReportListContainer
// @router /getReportList [get]
func (this *ReportController) GetReportList() {
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	total, totalErr := base.DBEngine.Table("report").Count(new(models.Report))

	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	var reportList []models.Report
	if total > 0 {
		err := base.DBEngine.Table("report").Desc("created").Limit(pageSize, pageSize*(pageNum-1)).Find(&reportList)
		if err != nil {
			this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+err.Error())
			return
		}
	}

	if reportList == nil {
		reportList = make([]models.Report, 0)
	}

	this.ReturnData = models.ReportListContainer{models.BaseListContainer{total, pageNum, pageTime}, reportList}
}
