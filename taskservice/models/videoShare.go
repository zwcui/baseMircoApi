/*
@Time : 2019/6/12 下午3:56
@Author : zwcui
@Software: GoLand
*/
package models

type VideoShare struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	VideoId       		int64			`description:"videoId" json:"videoId"`
	LastAuthorId       	int64			`description:"分销上级id" json:"lastAuthorId"`
	AuthorId	       	int64			`description:"分销当前id" json:"authorId"`
	Level		       	int				`description:"当前层级" json:"level"`
	Status		       	int				`description:"状态，1为正常" json:"status"`
	DirectNextShareNum	int				`description:"直接下级分享人数" json:"directNextShareNum" xorm:"notnull default 0"`
	DirectNextPayNum	int				`description:"直接下级购买人数" json:"directNextPayNum" xorm:"notnull default 0"`
	AllNextShareNum		int				`description:"所有下级分享人数" json:"allNextShareNum" xorm:"notnull default 0"`
	LatestForwardTime  	int64  			`description:"最新转发时间" json:"latestForwardTime"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}




//-----------------------------结构体-------------------------------------

type VideoShareContainer struct {
	BaseListContainer
	VideoShareInfoList				[]VideoShareInfo			`description:"视频分享记录列表" json:"video"`
}

//分享的视频记录
type VideoShareInfo struct {
	VideoShare                          `description:"分享视频记录" xorm:"extends"`
	CommissionNum		int				`description:"赚取的佣金" json:"commissionNum"`
	Title				string			`description:"标题" json:"title"`
	Nickname			string			`description:"视频作者的昵称" json:"nickname"`

}

type CommissionRecordContainer struct {
	BaseListContainer
	CommissionRecordList    []CommissionRecord			`description:"佣金明细列表" json:"video"`
}

type CommissionRecord struct {
	AccountCommissionRecord                     `description:"账户分佣记录" xorm:"extends"`
	PayerName               string              `description:"购买者姓名" json:"payerName"`
	Created             	int64  			    `description:"购买时间(订单创建时间)" json:"created" xorm:"created"`
}