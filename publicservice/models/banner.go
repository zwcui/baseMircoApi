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
	Position   int    `description:"广告位置，1为首页Banner，2为首页广告，3为频道banner，4为个人中心广告，5为vip个人中心广告（app跳转至公众号开店）6分销市场的首页banner 7分销市场首页广告" json:"position"  valid:"Required" xorm:"notnull default 0"`
	Status     int      `description:"状态，0为无效，1为有效" json:"status" xorm:"notnull default 0"`
	SortNo     int      `description:"排序号，越小越前" json:"sortNo" `
	NeedSignIn int      `description:"点击banner是否需要登录，1是0否" json:"needSignIn" xorm:"notnull default 0"`
	Created    int64  `description:"创建时间" xorm:"created" json:"created"`
	Updated    int64   `description:"更新时间" xorm:"updated" json:"updated"`
	DeletedAt  int64  `description:"删除时间" json:"-" xorm:"deleted"`
}





//---------------------结构体-----------------------------


type BannerListContainer struct {
	BannerList		[]Banner		`description:"banner列表" json:"bannerList"`
}
