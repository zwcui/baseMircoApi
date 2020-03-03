/*
@Time : 2020/2/7 上午11:42 
@Author : lianwu
@File : live.go
@Software: GoLand
*/
package models

type Live struct {
	Id                    int64   `description:"直播Id" json:"id" xorm:"pk autoincr"`
	Title                 string  `description:"标题" json:"title" `
	CoverUrl              string  `description:"封面" json:"coverUrl" `
	CoverWidth            int     `description:"封面宽度" json:"coverWidth" `
	CoverHeight           int     `description:"封面高度" json:"coverHeight" `
	AuthorId              int64   `description:"主播id" json:"authorId" `
	Status                int     `description:"直播状态 0:直播未开始 1:直播中 2:直播结束 3：撤销直播" json:"status"  xorm:"notnull default 0"`
	CloseTime             int64   `description:"直播关闭时间" json:"closeTime"`
	StartTime             int64   `description:"直播开始时间" json:"startTime"`
	ViewNum               int     `description:"当前观看人数" json:"viewNum" xorm:"notnull default 0"`
	RoomId                int64   `description:"直播所属房间ID" json:"roomId" xorm:"notnull default 0"`
	FlvDownstreamAddress  string  `description:"flv播放地址" json:"flvDownstreamAddress"`
	VideoUrl              string  `description:"录制视频播放地址" json:"videoUrl"`
	VideoDuration		  int     `description:"录制视频总时长" json:"videoDuration"`
	VideoCoverUrl         string  `description:"录制视频封面地址" json:"videoCoverUrl"`
	VideoFileId           string  `description:"录制视频fileId" json:"videoFileId"`
	NoticeStartTime       int64   `description:"直播预告开始时间" json:"noticeStartTime"`
	ConnectMicrophoneCount int    `description:"连麦次数" json:"connectMicrophoneCount" xorm:"notnull default 0"`
	Created           	  int64   `description:"注册时间" json:"created" xorm:"created"`
	Updated           	  int64   `description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	  int64   `description:"删除时间" json:"deleted" xorm:"deleted"`
	ScreenState			  int     `description:"屏幕状态，0：竖屏，1：横屏" xorm:"notnull default 0" json:"screenState"`
}



//课程列表
type LiveListContainer struct {
	BaseListContainer
	LiveList					[]Live					`description:"课程列表" json:"liveList"`
}