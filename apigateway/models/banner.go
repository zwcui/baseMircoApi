/*
@Time : 2019/9/9 上午9:32 
@Author : lianwu
@File : banner.go
@Software: GoLand
*/
package models

type Banner struct {
	Id        int64  `description:"id" json:"id" xorm:"pk autoincr"`
	CoverUrl   string `description:"图片url" json:"coverUrl" valid:"Required"`
	ActionUrl  string `description:"跳转url" json:"actionUrl"`
	Position   int    `description:"广告位置，1为首页Banner，2为首页广告，3为频道banner，4为个人中心广告" json:"position"  valid:"Required" xorm:"notnull default 0"`
	Status     int      `description:"状态，0为无效，1为有效" json:"status" xorm:"notnull default 0"`
	SortNo     int      `description:"排序号，越小越前" json:"sortNo" `
	Created    int64  `description:"创建时间" xorm:"created" json:"created"`
	Updated    int64   `description:"更新时间" xorm:"updated" json:"updated"`
	DeletedAt  int64  `description:"删除时间" json:"-" xorm:"deleted"`
}





//---------------------结构体-----------------------------


type BannerListContainer struct {
	BannerList		[]Banner		`description:"banner列表" json:"bannerList"`
}
