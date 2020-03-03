/*
@Time : 2019/9/24 上午10:28 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"encoding/json"
	"jingting_server/publicservice/models"
	"net/http"
	"jingting_server/publicservice/base"
)

type BugController struct {
	apiController
}

func (this *BugController) Prepare(){

}

// @Title	提交app错误日志
// @Description 提交app错误日志
// @Param	bugRecord 			body			models.BugRecord		true		"错误信息"
// @Success	200 {string} success
// @router /createBugRecord [post]
func (this *BugController) CreateBugRecord() {
	bugRecord := models.BugRecord{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &bugRecord)
	if err != nil {
		this.Err = err
		this.ErrCode = http.StatusInternalServerError
		return
	}

	bugRecord.Id = 0
	base.DBEngine.Table("bug_record").InsertOne(&bugRecord)

	this.ReturnData = "success"
}