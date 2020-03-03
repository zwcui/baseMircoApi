/*
@Time : 2019/3/6 下午6:12 
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

//定时任务，检查活动是否结束
func checkActivityIsEnd(){
	util.Logger.Info("定时任务，每小时检查活动是否结束")
	var activityList []models.Activity
	base.DBEngine.Table("activity").Where("status=1").And("end_type is not null").And("end_type != 0").Find(&activityList)
	if activityList == nil {
		activityList = make([]models.Activity, 0)
	}

	for _, activity := range activityList {
		if (activity.EndType == "1" && activity.EndTime < util.UnixOfBeijingTime()) ||
				(activity.EndType == "2" && activity.EndPrizeCount == 0) ||
				((activity.EndType == "1,2" || activity.EndType == "2,1") && (activity.EndTime < util.UnixOfBeijingTime() || activity.EndPrizeCount == 0)) {
			activity.Status = 3
			activity.ActualEndTime = util.UnixOfBeijingTime()
			base.DBEngine.Table("activity").Where("id=?", activity.Id).Cols("status", "actual_end_time").Update(&activity)

			//活动结束 推送结束文案与图片
			var authInfo models.AuthInfo
			base.DBEngine.Table("auth_info").Where("id=?", activity.AuthInfoId).Get(&authInfo)

			var joinActivityList []models.JoinActivity
			base.DBEngine.Table("join_activity").Where("activity_id=?", activity.Id).And("status=1").Find(&joinActivityList)

			mediaId := ""
			if joinActivityList != nil && len(joinActivityList) > 0 && activity.EndPicture != "" {
				//上传至微信
				mediaId, _ = util.UploadOnlineMedia(activity.EndPicture, authInfo, 1)
			}

			if joinActivityList != nil && len(joinActivityList) > 0 {
				util.Logger.Info("活动【"+activity.Name+"】结束，推送给"+strconv.Itoa(len(joinActivityList))+"个参与者，开始："+strconv.FormatInt(util.UnixOfBeijingTime(), 10))
			}

			for _, joinActivity := range joinActivityList {
				util.RequestSendTextCustomerServiceMessage(joinActivity.Openid, activity.EndContent, authInfo)
				util.RequestSendMediaCustomerServiceMessage(joinActivity.Openid, mediaId, authInfo)
			}

			if joinActivityList != nil && len(joinActivityList) > 0 {
				util.Logger.Info("活动【"+activity.Name+"】结束，推送给"+strconv.Itoa(len(joinActivityList))+"个参与者，结束："+strconv.FormatInt(util.UnixOfBeijingTime(), 10))
			}
		}
	}
}
