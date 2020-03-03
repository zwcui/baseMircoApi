package task

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"jingting_server/taskservice/base"
	"jingting_server/taskservice/models"
	"jingting_server/taskservice/remote"
	"jingting_server/taskservice/util"
	"strconv"
)

//检查直播预告是否过期
func checkExpireLive(){
	util.Logger.Info("定时任务，每分钟检查直播预告是否过期")

	expireTime := strconv.FormatInt(util.UnixOfBeijingTime() - 5 * 60, 10)
	var liveList []models.Live
	base.DBEngine.Table("live").Where("live.notice_start_time < ?",expireTime ).And("status = 0").Find(&liveList)
	for _, live := range liveList{
		live.Status = 3
		base.DBEngine.Table("live").Where("id = ?", live.Id).Cols("status").Update(&live)
		//推送消息
		appMessageOne := models.AppMessage{}
		appMessageOne.ReceiverId = live.AuthorId
		appMessageOne.ActionUrl = remote.JumpUrlWithKeyAndPramas(models.JTREVOKE_LIVE_JUMP_KEY, nil)
		appMessageOne.Content = "您发布的直播超过5分钟没有开始，系统已自动取消，请重新发布直播"
		appMessageOne.Type = 8
		_, pushErr := remote.PushMessageToUser( live.AuthorId, &appMessageOne, "", 0)
		if pushErr != nil {
			util.Logger.Info("checkExpireLive pushErr = ", pushErr.Error())
		}
		base.DBEngine.Table("app_message").InsertOne(&appMessageOne)
		//发送客服消息
		var author models.Author
		var dInfo models.UserSignInDeviceInfo
		base.DBEngine.Table("author").Where("id =?",  live.AuthorId).Get(&author)
		has, _ := base.DBEngine.Table("user_sign_in_device_info").Where("author_id = ?",  live.AuthorId).Get(&dInfo)
		if !has {
			content := "您发布的直播超过5分钟没有开始，系统已自动取消，请重新发布直播"
			util.RequestSendGZHTextCustomerServiceMessage(author.Openid, content, getGZHAccessToken())
		}

	}
}

//去腾讯云查询录制的视频（废弃，根据腾讯云回调获取）
func getLiveVideo(){
	util.Logger.Info("定时任务，去腾讯云查询录制的视频")

	var liveList []models.Live
	base.DBEngine.Table("live").Where("status=2").And("(video_url is null or video_url='')").Find(&liveList)
	for _, live := range liveList {
		//streamid = urlencode(sdkappid_roomid_userid_main)
		streamId := models.JTTencentCloudSDKAppID + "_" + strconv.FormatInt(live.RoomId, 10) + "_" + strconv.FormatInt(live.AuthorId, 10) + "_main"

		credential := common.NewCredential(models.JTTencentCloudSecretId, models.JTTencentCloudSecretKey)
		cpf := profile.NewClientProfile()

		client, _ := vod.NewClient(credential, regions.Chongqing, cpf)
		request := vod.NewSearchMediaRequest()
		request.StreamId = &streamId

		response, err := client.SearchMedia(request)
		if err != nil {
			util.Logger.Info("client.SearchMedia(&request) failed err:" +err.Error())
			continue
		}

		util.Logger.Info("SearchMedia streamId=" + streamId + " response.ToJsonString()")
		util.Logger.Info(response.ToJsonString())

		for _, media := range response.Response.MediaInfoSet {
			live.VideoFileId = *media.FileId
			if media.BasicInfo != nil {
				live.VideoCoverUrl = *media.BasicInfo.CoverUrl
				live.VideoUrl = *media.BasicInfo.MediaUrl
			}
			if media.MetaData != nil {
				live.VideoDuration = int(*media.MetaData.VideoDuration)
			}
		}

		base.DBEngine.Table("live").Where("id=?", live.Id).AllCols().Update(&live)
	}
}