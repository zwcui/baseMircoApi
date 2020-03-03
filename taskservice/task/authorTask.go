/*
@Time : 2019/9/16 上午11:26 
@Author : zwcui
@Software: GoLand
*/
package task

import (
	"jingting_server/taskservice/util"
	"jingting_server/taskservice/base"
)

func checkAuthorVIPEndTime(){
	util.Logger.Info("定时任务，每10分钟检查会员是否过期")

	updateSql := "update author set author.is_vip=0 where author.vip_end_time > 0 and author.vip_end_time < ? "

	_, err := base.DBEngine.Exec(updateSql, util.UnixOfBeijingTime())
	if err != nil {
		util.Logger.Info("checkAuthorVIPEndTime updateSql err:" + err.Error())
	}
}
