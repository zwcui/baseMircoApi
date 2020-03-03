/*
@Time : 2019/9/9 上午10:12 
@Author : lianwu
@File : bannerApi.go
@Software: GoLand
*/
package controllers

import (
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/base"
	"jingting_server/publicservice/util"
	"strings"
)

type BannerController struct {
	apiController
}

func (this *BannerController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/addBanner", "post", []int{1,2}},
		{"/updateBanner", "patch", []int{1,2}},



	}
	this.userAuth()
	this.NeedAuthorAuthList = []RequestPathAndMethod{

	}
	this.authorAuth()
}


// @Title 新增Banner（h5使用）
// @Description 新增Banner（h5使用）
// @Param	coverUrl		formData		string	  		true		"图片url"
// @Param	actionUrl  		formData		string	  		true		"跳转url"
// @Param	position		formData		int      		true		"广告位置，1为首页Banner，2为首页广告，3为频道banner，4为个人中心广告，5为vip个人中心广告"
// @Param	status  		formData		int     		true		"状态，0为无效，1为有效"
// @Param	sortNo  		formData		int     		true		"排序号，越小越前"
// @Param	needSignIn 		formData		int     		true		"点击banner是否需要登录，1是0否"
// @Success 200 {string} success
// @router /addBanner [post]
func (this *BannerController) AddBanner() {
	coverUrl := this.MustString("coverUrl")
	actionUrl := this.MustString("actionUrl")
	position := this.MustInt("position")
	status := this.MustInt("status")
	sortNo := this.MustInt("sortNo")
	needSignIn := this.MustInt("needSignIn")

	var banner models.Banner
	banner.ActionUrl = actionUrl
	banner.CoverUrl = coverUrl
	banner.Position = position
	banner.Status = status
	banner.SortNo = sortNo
	banner.NeedSignIn = needSignIn
	_, err :=base.DBEngine.Table("banner").InsertOne(&banner)
	if err != nil {
		util.Logger.Info("AddBanner err = ", err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
	}
	this.ReturnData = "success"
}

// @Title 修改Banner（h5使用）
// @Description 修改Banner（h5使用）
// @Param	id      		formData		int64     		true		"排序号，越小越前"
// @Param	coverUrl		formData		string	  		false		"图片url"
// @Param	actionUrl  		formData		string	  		false		"跳转url"
// @Param	position		formData		int      		false		"广告位置，1为首页Banner，2为首页广告，3为频道banner，4为个人中心广告，5为vip个人中心广告"
// @Param	status  		formData		int     		false		"状态，0为无效，1为有效"
// @Param	sortNo  		formData		int     		false		"排序号，越小越前"
// @Param	needSignIn 		formData		int     		false		"点击banner是否需要登录，1是0否"
// @Success 200 {string} success
// @router /updateBanner [patch]
func (this *BannerController) UpdateBanner() {

	id := this.MustInt64("id")
	coverUrl := this.GetString("coverUrl", "")
	actionUrl := this.GetString("actionUrl","")
	position, _ := this.GetInt("position",-1)
	status, _ := this.GetInt("status",-1)
	sortNo, _ := this.GetInt("sortNo",-1)
	needSignIn, _ := this.GetInt("needSignIn",-1)

	var banner models.Banner
	has, _ := base.DBEngine.Table("banner").Where("id = ?", id).Get(&banner)
	if has {

		if coverUrl != ""{
			banner.CoverUrl = coverUrl
		}
		if actionUrl != ""{
			banner.ActionUrl = actionUrl
		}
		if position != -1{
			banner.Position = position
		}
		if status != -1{
			banner.Status = status
		}
		if sortNo != -1{
			banner.SortNo = sortNo
		}
		if needSignIn != -1{
			banner.NeedSignIn = needSignIn
		}
		_, err := base.DBEngine.Table("banner").Where("id =?", id).AllCols().Update(&banner)
		if err != nil {
			util.Logger.Info("AddBanner err = ", err.Error())
			this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
			return
		}

	} else {
		this.ReturnData = util.GenerateAlertMessage(models.BannerError100)
		return
	}

	this.ReturnData = "success"
}


// @Title banner列表
// @Description banner列表
// @Param	position		query		int      		false		"广告位置，1为首页Banner，2为首页广告，3为频道banner，4为个人中心广告，5为vip个人中心广告 6分销市场的首页banner 7分销市场首页广告"
// @Success 200 {object} models.BannerListContainer
// @router /getBannerList [get]
func (this *BannerController) GetBannerList() {

	loginFlag := false
	authorization := this.Ctx.Request.Header.Get("Authorization")
	if authorization != "" {
		loginFlag = true
	}

	position, err := this.GetInt("position",0)
	var bannerList []models.Banner
	if position == 0 {
		err = base.DBEngine.Table("banner").OrderBy("sort_no").Find(&bannerList)
	} else {
		err = base.DBEngine.Table("banner").Where("position = ? ",position).OrderBy("sort_no").Find(&bannerList)
	}
	if err != nil {
		util.Logger.Info("GetBannerList err = ", err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
	}
	if bannerList == nil {
		bannerList = make([]models.Banner, 0)
	}

	for i := 0;i < len(bannerList);i++ {
		if bannerList[i].Position == 1 && !loginFlag && strings.Contains(bannerList[i].ActionUrl, "https://www.jingtingedu.com/shop/#/detail") {
			bannerList[i].ActionUrl = ""
		}
	}

	this.ReturnData = models.BannerListContainer{bannerList}
}

// @Title 删除banner(后台使用)
// @Description 删除banner(后台使用)
// @Param	id						query 	  			int64				true		"id"
// @Success 200 {string} success
// @router /deleteBanner [delete]
func (this *BannerController) DeleteBanner() {

	id := this.MustInt64("id")
	var banner models.Banner
	banner.Id = id
	_, err := base.DBEngine.Table("banner").Delete(&banner)
	if err != nil {
		util.Logger.Info("DeleteBanner err = ", err.Error())
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100)
	}
	this.ReturnData = "success"
}



