/*
@Time : 2019/5/6 上午11:18 
@Author : zwcui
@Software: GoLand
*/
package models

//----------------------------------表结构------------------------------------------
//名片
type BusinessCard struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Title      				string				`description:"名片标题" json:"title"`
	Name      				string				`description:"姓名" json:"name"`
	Phone      				string				`description:"电话" json:"phone"`
	Wechat      			string				`description:"微信" json:"wechat"`
	Email      				string				`description:"邮箱" json:"email"`
	Company      			string				`description:"公司" json:"company"`
	Position      			string				`description:"职位" json:"position"`
	Avatar      			string				`description:"头像" json:"avatar"`
	Introduction      		string				`description:"简介" json:"introduction"`
	Other	      			string				`description:"其他" json:"other"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}


//------------------------------------结构体--------------------------------------------
type BusinessCardListContainer struct {
	BaseListContainer
	BusinessCardList		[]BusinessCard		`description:"名片列表" json:"businessCardList"`
}

type BusinessCardDetailContainer struct {
	BusinessCard			BusinessCard		`description:"名片" json:"businessCardDetail"`
}
