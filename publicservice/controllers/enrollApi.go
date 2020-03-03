/*
@Time : 2020/1/14 下午5:12 
@Author : lianwu
@File : enrollApi.go
@Software: GoLand
*/
package controllers

import (
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/util"
	"strconv"
	"jingting_server/publicservice/base"
)

type EnrollController struct {
	apiController
}

func (this *EnrollController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
	}
	this.userAuth()
}

// @Title 报名领取课程
// @Description 报名领取课程
// @Param	authorId				formData		int64  		    true		"报名者id"
// @Param	salesmanId		        formData		int64  		    false		"业务员id"
// @Param	phoneNumber				formData		string  		true		"手机号"
// @Param	realName				formData		string	  		true		"真实姓名"
// @Param	profession				formData		string  		true		"职业"
// @Param	city				    formData		string	  		true		"城市"
// @Success 200 {string} success
// @router /addEnroll [post]
func (this *EnrollController) AddEnroll() {
	authorId := this.MustInt64("authorId")
	salesmanId, _ := this.GetInt64("salesmanId", 0)
	phoneNumber := this.MustString("phoneNumber")
	realName := this.MustString("realName")
	profession := this.MustString("profession")
	city := this.MustString("city")
	has, _ := base.DBEngine.Table("enroll").Where("author_id = ?", authorId).Get(new(models.Enroll))
	var enroll models.Enroll
	if has {
		this.ReturnData = util.GenerateAlertMessage(models.EnrollError100)
		return

	} else {
		enroll.AuthorId = authorId
		enroll.SalesmanId = salesmanId
		enroll.PhoneNumber = phoneNumber
		enroll.RealName = realName
		enroll.Profession = profession
		enroll.City = city
		base.DBEngine.Table("enroll").InsertOne(&enroll)
	}

	this.ReturnData = "success"

}


// @Title 报名领取课程
// @Description 报名领取课程
// @Param	authorId				query		        int64  		    false		"报名者id"
// @Param	phoneNumber				query		        string  		false		"手机号"
// @Param	lastAuthorName			query		        string	  		false		"上级姓名"
// @Param	realName				query		        string	  		false		"真实姓名"
// @Param	salesmanId		        query		        int64  		    false		"业务员id"
// @Param	pageNum					query 	  			int				true		"page num start from 1"
// @Param	pageTime				query 	  			int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize				query 	  			int				false		"page size default is 15"
// @Success 200 {object} models.EnrollListContainer
// @router /getEnrollList [get]
func (this *EnrollController) GetEnrollList() {
	authorId, _ := this.GetInt64("authorId", 0)
	phoneNumber := this.GetString("phoneNumber", "")
	realName := this.GetString("realName", "")
	lastAuthorName := this.GetString("lastAuthorName", "")
	salesmanId, _ := this.GetInt64("salesmanId", 0)
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}


	totalSql := " select count(1) from enroll left join share on enroll.author_id = share.author_id left join author on share.last_author_id = author.id where enroll.deleted_at is null"
	dataSql := " select enroll.*, author.nickname as last_author_name from enroll left join share on enroll.author_id = share.author_id left join author on share.last_author_id = author.id where enroll.deleted_at is null"

	if authorId != 0 {
		totalSql += " and enroll.author_id = " + strconv.FormatInt(authorId, 10)
		dataSql += " and enroll.author_id = " + strconv.FormatInt(authorId, 10)
	}
	if salesmanId != 0 {
		totalSql += " and enroll.salesman_id = " + strconv.FormatInt(salesmanId, 10)
		dataSql += " and enroll.salesman_id = " + strconv.FormatInt(salesmanId, 10)
	}
	if phoneNumber != "" {
		totalSql += " and enroll.phone_number = '" + phoneNumber + "'"
		dataSql += " and enroll.phone_number = '" + phoneNumber + "'"
	}
	if realName != "" {
		totalSql += " and enroll.real_name like '%" + realName + "%'"
		dataSql += " and enroll.real_name like '%" + realName + "%'"
	}
	if lastAuthorName != "" {
		totalSql += " and author.nickname like '%" + lastAuthorName + "%'"
		dataSql += " and author.nickname like '%" + lastAuthorName + "%'"
	}
	dataSql += " ORDER BY enroll.created desc limit "+strconv.Itoa(pageSize*(pageNum-1))+", "+strconv.Itoa(pageSize)
	total, totalErr := base.DBEngine.SQL(totalSql).Count(new(models.Enroll))
	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}
	var enrollList []models.EnrollDetail
	if total > 0 {
		base.DBEngine.SQL(dataSql).Find(&enrollList)
	}
	if enrollList == nil {
		enrollList = make([]models.EnrollDetail, 0)
	}

	this.ReturnData = models.EnrollListContainer{models.BaseListContainer{total, pageNum, pageTime}, enrollList}


}
