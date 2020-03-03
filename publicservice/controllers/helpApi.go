/*
@Time : 2019/3/21 下午12:00 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/base"
	"jingting_server/publicservice/util"
	"strconv"
)

type HelpController struct {
	apiController
}

func (this *HelpController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/addHelp", "post", []int{2}},
		{"/deleteHelp", "delete", []int{2}},
	}
	this.userAuth()
}

// @Title 新增帮助（h5使用）
// @Description 新增帮助（h5使用）
// @Param	name					formData			string  		true		"名称"
// @Param	level					formData			int		  		true		"级别 1为分类，2为章节，3为标题"
// @Param	lastId					formData			int64	  		false		"上级id"
// @Param	content					formData			string	  		false		"文章内容"
// @Success 200 {string} success
// @router /addHelp [post]
func (this *HelpController) AddHelp() {
	name := this.MustString("name")
	level := this.MustInt("level")
	lastId, _ := this.GetInt64("lastId", 0)
	content := this.GetString("content")

	if lastId != 0 {
		hasLastHelp, _ := base.DBEngine.Table("help").Where("id=?", lastId).Get(new(models.Help))
		if !hasLastHelp {
			this.ReturnData = util.GenerateAlertMessage(models.HelpError100)
			return
		}
	}

	var help models.Help
	help.Name = name
	help.Level = level
	help.LastId = lastId
	help.Content = content
	base.DBEngine.Table("help").InsertOne(&help)

	this.ReturnData = "success"
}

// @Title 新增帮助（h5使用）
// @Description 新增帮助（h5使用）
// @Param	id						formData			int64	  		true		"id"
// @Param	name					formData			string  		true		"名称"
// @Param	level					formData			int		  		true		"级别 1为分类，2为章节，3为标题"
// @Param	lastId					formData			int64	  		false		"上级id"
// @Param	content					formData			string	  		false		"文章内容"
// @Success 200 {string} success
// @router /updateHelp [patch]
func (this *HelpController) UpdateHelp() {
	id := this.MustInt64("id")
	name := this.MustString("name")
	level := this.MustInt("level")
	lastId, _ := this.GetInt64("lastId", 0)
	content := this.GetString("content")

	var help models.Help
	hasHelp, _ := base.DBEngine.Table("help").Where("id=?", id).Get(&help)
	if !hasHelp {
		this.ReturnData = util.GenerateAlertMessage(models.HelpError200)
		return
	}

	help.Name = name
	help.Level = level
	help.LastId = lastId
	help.Content = content
	base.DBEngine.Table("help").Where("id=?", id).AllCols().Update(&help)

	this.ReturnData = "success"
}

// @Title 获取帮助列表（h5使用）
// @Description 获取帮助列表（h5使用）
// @Param	level						query 	  			int				true		"级别 1为分类，2为章节，3为标题"
// @Param	pageNum						query 	  			int				true		"page num start from 1"
// @Param	pageTime					query 	  			int64			false		"page time should be empty when pagenum == 1"
// @Param	pageSize					query 	  			int				false		"page size default is 15"
// @Success 200 {object} models.HelpDetailListContainer
// @router /getHelpList [get]
func (this *HelpController) GetHelpList() {
	level := this.MustInt("level")
	pageNum := this.MustInt("pageNum")
	pageTime, _ := this.GetInt64("pageTime", util.UnixOfBeijingTime())
	pageSize := this.GetPageSize("pageSize")
	if pageNum <= 1 {
		pageTime = util.UnixOfBeijingTime()
	}

	totalSql := "select count(1) from help where deleted_at is null and level='"+strconv.Itoa(level)+"' "
	dataSql := ""
	if level == 1 {
		dataSql = "select help.id as id, help.name as classify, help.level as level, help.created as created, help.updated as updated from help where help.deleted_at is null and help.level=1 order by help.created asc"
	} else if level == 2 {
		dataSql = "select chapter.id as id, chapter.name as chapter, classify.name as classify, chapter.level as level, chapter.last_id as last_id, chapter.created as created, chapter.updated as updated from help chapter left join help classify on classify.id=chapter.last_id where chapter.deleted_at is null and classify.deleted_at is null and chapter.level=2 order by chapter.created asc"
	} else if level == 3 {
		dataSql = "select title.id as id, title.name as title, title.content as content, chapter.name as chapter, classify.name as classify, title.level as level, title.last_id as last_id, title.created as created, title.updated as updated from help title left join help chapter on chapter.id=title.last_id left join help classify on classify.id=chapter.last_id where title.deleted_at is null and chapter.deleted_at is null and classify.deleted_at is null and title.level=3 order by classify.name asc, classify.created asc, chapter.name asc, chapter.created asc, title.created asc"
	}
	dataSql += " limit "+strconv.Itoa(pageSize*(pageNum-1))+", "+strconv.Itoa(pageSize)

	total, totalErr := base.DBEngine.SQL(totalSql).Count(new(models.Help))
	if totalErr != nil {
		this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+totalErr.Error())
		return
	}

	var helpDetailList []models.HelpDetail
	if total > 0 {
		err := base.DBEngine.SQL(dataSql).Find(&helpDetailList)
		if err != nil {
			this.ReturnData = util.GenerateAlertMessage(models.CommonError100, " err:"+err.Error())
			return
		}
	}

	if helpDetailList == nil {
		helpDetailList = make([]models.HelpDetail, 0)
	}

	this.ReturnData = models.HelpDetailListContainer{models.BaseListContainer{total, pageNum, pageTime}, helpDetailList}
}

// @Title 获取帮助列表（h5官网使用）
// @Description 获取帮助列表（h5官网使用）
// @Success 200 {object} models.AllHelp
// @router /getAllHelpList [get]
func (this *HelpController) GetAllHelpList() {
	var helpClassifyList []models.HelpClassify

	var classifyList []models.Help
	base.DBEngine.Table("help").Where("level=1").Asc("created").Find(&classifyList)
	if classifyList == nil {
		classifyList = make([]models.Help, 0)
	}

	for _, classify := range classifyList {
		var chapterList []models.Help
		base.DBEngine.Table("help").Where("level=2").And("last_id=?", classify.Id).Asc("created").Find(&chapterList)
		if chapterList == nil {
			chapterList = make([]models.Help, 0)
		}

		var helpChapterList []models.HelpChapter
		for _, chapter := range chapterList {
			var helpChapter models.HelpChapter
			var helpTitleList []models.HelpTitle
			var titleList []models.Help
			base.DBEngine.Table("help").Where("level=3").And("last_id=?", chapter.Id).Asc("created").Find(&titleList)
			if titleList == nil {
				titleList = make([]models.Help, 0)
			}

			for _, title := range titleList {
				helpTitleList = append(helpTitleList, models.HelpTitle{title})
			}

			if helpTitleList == nil {
				helpTitleList = make([]models.HelpTitle, 0)
			}

			helpChapter.Chapter = chapter
			helpChapter.TitleList = helpTitleList
			helpChapterList = append(helpChapterList, helpChapter)
		}

		var helpClassify models.HelpClassify
		helpClassify.Classify = classify
		helpClassify.ChapterList = helpChapterList
		helpClassifyList = append(helpClassifyList, helpClassify)
	}

	if helpClassifyList == nil {
		helpClassifyList = make([]models.HelpClassify, 0)
	}

	this.ReturnData = models.AllHelp{helpClassifyList}
}

// @Title 获取帮助详情信息（h5官网使用）
// @Description 获取帮助详情信息（h5官网使用）
// @Param	id						query 	  			int64				true		"id"
// @Success 200 {object} models.Help
// @router /getHelpInfo [get]
func (this *HelpController) GetHelpInfo() {
	id := this.MustInt64("id")

	var help models.Help
	hasHelp, _ := base.DBEngine.Table("help").Where("id=?", id).Get(&help)
	if !hasHelp {
		this.ReturnData = util.GenerateAlertMessage(models.HelpError200)
		return
	}

	if help.Level == 1 {
		var helpClassifyList []models.HelpClassify

		var chapterList []models.Help
		base.DBEngine.Table("help").Where("level=2").And("last_id=?", help.Id).Asc("created").Find(&chapterList)
		if chapterList == nil {
			chapterList = make([]models.Help, 0)
		}

		var helpChapterList []models.HelpChapter
		for _, chapter := range chapterList {
			var helpChapter models.HelpChapter
			var helpTitleList []models.HelpTitle
			var titleList []models.Help
			base.DBEngine.Table("help").Where("level=3").And("last_id=?", chapter.Id).Asc("created").Find(&titleList)
			if titleList == nil {
				titleList = make([]models.Help, 0)
			}

			for _, title := range titleList {
				helpTitleList = append(helpTitleList, models.HelpTitle{title})
			}

			helpChapter.Chapter = chapter
			helpChapter.TitleList = helpTitleList
			helpChapterList = append(helpChapterList, helpChapter)
		}

		var helpClassify models.HelpClassify
		helpClassify.Classify = help
		helpClassify.ChapterList = helpChapterList
		helpClassifyList = append(helpClassifyList, helpClassify)

		if helpClassifyList == nil {
			helpClassifyList = make([]models.HelpClassify, 0)
		}

		this.ReturnData = models.ClassifyListContainer{helpClassifyList}
	} else if help.Level == 2 {
		var helpChapterList []models.HelpChapter

		var helpChapter models.HelpChapter
		var helpTitleList []models.HelpTitle
		var titleList []models.Help
		base.DBEngine.Table("help").Where("level=3").And("last_id=?", help.Id).Asc("created").Find(&titleList)
		if titleList == nil {
			titleList = make([]models.Help, 0)
		}

		for _, title := range titleList {
			helpTitleList = append(helpTitleList, models.HelpTitle{title})
		}

		helpChapter.Chapter = help
		helpChapter.TitleList = helpTitleList
		helpChapterList = append(helpChapterList, helpChapter)

		if helpChapterList == nil {
			helpChapterList = make([]models.HelpChapter, 0)
		}

		this.ReturnData = models.ChapterListContainer{helpChapterList}
	} else if help.Level == 3 {

		this.ReturnData = models.TitleContainer{models.HelpTitle{help}}
	}

	return
}

// @Title 删除帮助（h5使用）
// @Description 删除帮助（h5使用）
// @Param	id						query 	  			int64				true		"id"
// @Success 200 {string} success
// @router /deleteHelp [delete]
func (this *HelpController) DeleteHelp() {
	id := this.MustInt64("id")

	var help models.Help
	hasHelp, _ := base.DBEngine.Table("help").Where("id=?", id).Get(&help)
	if !hasHelp {
		this.ReturnData = util.GenerateAlertMessage(models.HelpError200)
		return
	}

	if help.Level == 1 {
		deleteSql := "update help set deleted_at=? where level=3 and last_id in (select id from help where level=2 and last_id=?)"
		base.DBEngine.Exec(deleteSql, util.UnixOfBeijingTime(), id)

		deleteSql = "update help set deleted_at=? where level=2 and last_id=?"
		base.DBEngine.Exec(deleteSql, util.UnixOfBeijingTime(), id)

		base.DBEngine.Table("help").Where("id=?", id).Delete(&help)
	} else if help.Level == 2 {
		deleteSql := "update help set deleted_at=? where level=3 and last_id=?"
		base.DBEngine.Exec(deleteSql, util.UnixOfBeijingTime(), id)

		base.DBEngine.Table("help").Where("id=?", id).Delete(&help)
	} else if help.Level == 3 {
		base.DBEngine.Table("help").Where("id=?", id).Delete(&help)
	}

	this.ReturnData = "success"
}








