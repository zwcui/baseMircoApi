/*
@Time : 2019/3/19 下午1:49 
@Author : zwcui
@Software: GoLand
*/
package task

import (
	"jingting_server/taskservice/models"
	"jingting_server/taskservice/base"
	"jingting_server/taskservice/util"
	"strconv"
)

//发送模板消息
func sendTemplateMessage(){
	util.Logger.Info("定时任务，每分钟检查是否需要发送模板消息")
	authInfoMap := make(map[int64]models.AuthInfo)

	var templateMessageDataList []models.TemplateMessageData
	base.DBEngine.Table("template_message_data").Where("status != 1").And("send_time != 0").And("send_time < ?", util.UnixOfBeijingTime() - 60).And("send_time < ?", util.UnixOfBeijingTime() + 60).Asc("auth_info_id", "created").Find(&templateMessageDataList)
	util.Logger.Info("send_time < "+strconv.FormatInt(util.UnixOfBeijingTime() - 60, 10))
	for _, templateMessageData := range templateMessageDataList {
		util.Logger.Info("templateMessageData.id = "+strconv.FormatInt(templateMessageData.Id, 10))
		authInfo, hasAuthInfo := authInfoMap[templateMessageData.AuthInfoId]
		if !hasAuthInfo {
			base.DBEngine.Table("auth_info").Where("id=?", templateMessageData.AuthInfoId).Get(&authInfo)
			authInfoMap[templateMessageData.AuthInfoId] = authInfo
		}

		var subscriberList []models.Subscriber

		subscriberSql := "select * from subscriber where auth_info_id=? and deleted_at is null "
		if templateMessageData.SendSex != 0 {
			subscriberSql += " and sex='"+strconv.Itoa(templateMessageData.SendSex)+"' "
		}
		if templateMessageData.SendProvince != "" {
			subscriberSql += " and province='"+templateMessageData.SendProvince+"' "
		}
		if templateMessageData.SendCity != "" {
			subscriberSql += " and city='"+templateMessageData.SendCity+"' "
		}
		err := base.DBEngine.SQL(subscriberSql, authInfo.Id).Find(&subscriberList)
		if err != nil {
			util.Logger.Info("err:"+err.Error())
		}

		for _, subscriber := range subscriberList {
			err := util.RequestSendTemplateMessage(templateMessageData.TemplateId, subscriber.Openid, templateMessageData.Url, templateMessageData.Data, authInfo)
			if err != nil {
				util.Logger.Info("err:"+err.Error())
			}
		}

		templateMessageData.Status = 1
		base.DBEngine.Table("template_message_data").Where("id=?", templateMessageData.Id).Cols("status").Update(&templateMessageData)
	}
}

//删除已发送并超过60天的模板消息
func deleteTemolateMessage(){
	util.Logger.Info("定时任务，每天检查删除已发送并超过60天的模板消息")
	deleteSql := "update template_message_data set deleted_at=? where status=1 and send_time < ? and deleted_at is null "
	_, err := base.DBEngine.Exec(deleteSql, util.UnixOfBeijingTime(), util.UnixOfBeijingTime() - 60 * 24 * 60 * 60)
	if err != nil {
		util.Logger.Info("deleteTemolateMessage err:"+err.Error())
	}
}
