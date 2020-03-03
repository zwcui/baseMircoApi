/*
@Time : 2020/1/9 上午9:50 
@Author : lianwu
@File : enroll.go
@Software: GoLand
*/
package models


type Enroll struct {
	Id       			    int64			    `description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId				int64				`description:"报名人id" json:"authorId"`
	SalesmanId	            int64				`description:"业务员id" json:"salesmanId"`
	PhoneNumber  		    string 			    `description:"用户手机号" json:"phoneNumber"`
	RealName			    string			    `description:"真实姓名" json:"realName"`
	Profession			    string			    `description:"职业" json:"profession"`
	City			        string			    `description:"城市" json:"city"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`

}

type RedeemCode struct {
	Id       			    int64			    `description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId				int64				`description:"领取人id" json:"authorId"`
	Code                    string				`description:"兑换码" json:"code"`
	OrderId				    int64				`description:"订单id" json:"orderId"`
	Status					int					`description:"兑换码状态，0未生效，1已生效，2已兑换 3已过期" json:"status" xorm:"notnull default 0"`
	ExpireTime			    int64			    `description:"兑换码过期时间" json:"expireTime"`
	Scene					int					`description:"兑换码使用场景，1 兑换课程" json:"scene" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`


}







//---------------------结构体-----------------------------

type EnrollListContainer struct {
	BaseListContainer
	EnrollList					[]EnrollDetail					`description:"报名列表" json:"enrollList"`

}

type EnrollDetail  struct {
	Enroll                                                  `description:"报名信息" xorm:"extends"`
	LastAuthorName                    string			    `description:"上级姓名" json:"lastAuthorName"`


}