/*
@Time : 2019/9/16 上午11:18 
@Author : zwcui
@Software: GoLand
*/
package models

//举报记录
type Report struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId       			int64				`description:"用户id" json:"authorId"`
	ValueId       			int64				`description:"type为1时:视频id, type为2时:评论id" json:"valueId"`
	Type       			    int				    `description:"举报类型 1:视频 2:评论" json:"type"`
	Content					string				`description:"举报内容" json:"content"`
	Classify				int				    `description:"举报内容类别 1:违法违禁, 2:色情, 3:低俗, 4:赌博诈骗, 5:血腥暴力, 6:人生攻击, 7:与其他视频相同, 8:不良封面标题, 9;青少年不良信息, 10:其他 " json:"classify"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`

}

type ReportListContainer struct {
	BaseListContainer
	ReportList  		    []Report 	         `description:"举报记录列表" json:"viewVideoList"`
}