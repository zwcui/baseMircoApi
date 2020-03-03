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
)

type RecentActivityController struct {
	apiController
}

func (this *RecentActivityController) Prepare(){
	this.NeedAuthorAuthList = []RequestPathAndMethod{
	}
	this.authorAuth()
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/addRecentActivity", "post", []int{}},
		{"/updateRecentActivity", "patch", []int{}},
		{"/deleteRecentActivity", "delete", []int{}},
	}
	this.userAuth()
}

// @Title 新增近期活动
// @Description 新增近期活动
// @Param	uId			    	formData			int64		  	true		"用户id"
// @Param	title				formData			string	  		true		"活动名称"
// @Param	banner				formData			string	  		true		"活动海报"
// @Param	content	    		formData			string 		  	true		"活动介绍"
// @Success 200 {string} success
// @router /addRecentActivity [post]
func (this *RecentActivityController) AddRecentActivity() {
	uId := this.MustInt64("uId")
	title := this.MustString("title")
	banner := this.MustString("banner")
	content := this.MustString("content")

	var recentActivity models.RecentActivity
	recentActivity.UId = uId
	recentActivity.Banner = banner
	recentActivity.Title = title
	recentActivity.Content = content
	base.DBEngine.Table("recent_activity").InsertOne(&recentActivity)

	this.ReturnData = "success"
}

// @Title 修改近期活动
// @Description 修改近期活动
// @Param	recentActivityId        	query				int64		  	true		"近期活动id"
// @Param	title						query				string	  		true		"活动名称"
// @Param	banner						query				string	  		true		"活动海报"
// @Param	content	    				query				string 		  	true		"活动介绍"
// @Success 200 {string} success
// @router /updateRecentActivity [patch]
func (this *RecentActivityController) UpdateRecentActivity() {
	recentActivityId := this.MustInt64("recentActivityId")
	title := this.MustString("title")
	banner := this.MustString("banner")
	content := this.MustString("content")

	var recentActivity models.RecentActivity
	hasRecentActivity, _ := base.DBEngine.Table("recent_activity").Where("id=?", recentActivityId).Get(&recentActivity)
	if !hasRecentActivity {
		this.ReturnData = util.GenerateAlertMessage(models.RecentActivityError100)
		return
	}

	recentActivity.Title = title
	recentActivity.Banner = banner
	recentActivity.Content = content
	base.DBEngine.Table("recent_activity").Where("id=?", recentActivityId).AllCols().Update(&recentActivity)

	this.ReturnData = "success"
}

// @Title 近期活动列表
// @Description 近期活动列表
// @Param	pageNum					query 	  			int				true		"page num start from 1"
// @Param	pageTime				query 	  			int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize				query 	  			int				false		"page size default is 15"
// @Success 200 {object} models.RecentActivityContainer
// @router /getRecentActivityList [get]
func (this *RecentActivityController) GetRecentActivityList() {
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	total, totalErr := base.DBEngine.Table("recent_activity").Count(new(models.RecentActivity))

	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	var recentActivityList []models.RecentActivity
	if total > 0 {
		base.DBEngine.Table("recent_activity").Desc("created").Limit(pageSize, pageSize*(pageNum-1)).Find(&recentActivityList)
	}

	if recentActivityList == nil {
		recentActivityList = make([]models.RecentActivity, 0)
	}

	this.ReturnData = models.RecentActivityContainer{models.BaseListContainer{total, pageNum, pageTime}, recentActivityList}
}

// @Title 删除近期活动
// @Description 删除近期活动
// @Param	recentActivityId        	query				int64		  	true		"近期活动id"
// @Success 200 {string} success
// @router /deleteRecentActivity [delete]
func (this *RecentActivityController) DeleteRecentActivity() {
	recentActivityId := this.MustInt64("recentActivityId")

	base.DBEngine.Table("recent_activity").Where("id=?", recentActivityId).Delete(new(models.RecentActivity))

	this.ReturnData = "success"
}

// @Title 活动详情
// @Description 活动详情
// @Param	id				    query		int64	  		true		"活动id"
// @Success 200 {object} models.RecentActivityDetailContainer
// @router /getRecentActivity [get]
func (this *RecentActivityController) GetRecentActivity() {

	id := this.MustInt64("id")
	var activity models.RecentActivity

	has, _ := base.DBEngine.Table("recent_activity").Where("id = ?", id).Get(&activity)
	if has {
		this.ReturnData = models.RecentActivityDetailContainer{activity}
	} else {
		this.ReturnData = util.GenerateAlertMessage(models.RecentActivityError100)
	}


}