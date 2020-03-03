/*
@Time : 2019/9/17 下午6:16 
@Author : zwcui
@Software: GoLand
*/
package models

//近期活动
type RecentActivity struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	UId		      			int64				`description:"管理员id" json:"uId"`
	Title	      			string				`description:"活动名称" json:"title"`
	Banner	      			string				`description:"活动海报" json:"banner"`
	Content      			string				`description:"活动介绍" json:"content" xorm:"text"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//------------------------------------结构体------------------------------------------
type RecentActivityContainer struct {
	BaseListContainer
	RecentActivityList		[]RecentActivity			`description:"近期活动列表" json:"recentActivityList"`
}
type RecentActivityDetailContainer struct {
	RecentActivity          RecentActivity   `description:"活动详情" json:"recentActivity"`

}