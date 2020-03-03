/*
@Time : 2019/6/17 下午2:22 
@Author : zwcui
@Software: GoLand
*/
package models

//-------------------------结构体--------------------------------------------
//任务
type TencentCloudTask struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	RequestId				string				`description:"requestId" json:"requestId"`
	TaskId					string				`description:"taskId" json:"taskId"`
	FileId					string				`description:"fileId" json:"fileId"`
	TaskType				int					`description:"任务类型，1为鉴黄，2为转码" json:"taskType"`
	Status					int					`description:"状态，0为等待中，1为已完成，2为处理中" json:"status" xorm:"notnull default 0"`
	Result					string				`description:"处理结果" json:"result"`
	ErrorInfo				string				`description:"错误信息" json:"errorInfo"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//任务查询
type TencentCloudTaskDetail struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	RequestId				string				`description:"requestId" json:"requestId"`
	TaskId					string				`description:"taskId" json:"taskId"`
	TaskType				string				`description:"Procedure：视频处理任务；EditMedia：视频编辑任务；WechatPublish：微信发布任务；ComposeMedia：制作媒体文件任务。2017:Transcode：视频转码任务；SnapshotByTimeOffset：视频截图任务；Concat：视频拼接任务；Clip：视频剪辑任务；ImageSprites：截取雪碧图任务。" json:"taskType"`
	Status					string				`description:"WAITING：等待中；PROCESSING：处理中；FINISH：已完成" json:"status"`
	Response	     		string 				`description:"返回结果" json:"response" xorm:"text"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//---------------------------常量-----------------------------------------------------
//视频转码模板，参考https://console.cloud.tencent.com/vod/video-process/template
const (
	//Video_Transcode_MP4_320 	uint64	= 10
	Video_Transcode_MP4_640 	uint64	= 20
	Video_Transcode_MP4_1280 	uint64	= 30
	Video_Transcode_MP4_1920 	uint64	= 40
	//Video_Transcode_FLV_320 	uint64	= 10046
	//Video_Transcode_FLV_640 	uint64	= 10047
	//Video_Transcode_FLV_1280 	uint64	= 10048
)