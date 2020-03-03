/*
@Time : 2019/10/31 下午8:52 
@Author : lianwu
@File : store.go
@Software: GoLand
*/
package models


//店铺
type Store struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId  				int64				`description:"用户id" json:"authorId"`
	Cover					string				`description:"封面" json:"cover"`

	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type StoreCourse struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	StoreId  				int64				`description:"店铺id" json:"storeId"`
	CourseId				int64				`description:"课程id" json:"courseId"`

	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`


}

type StoreCourseContainer struct {
	AuthorShort             AuthorShort         `description:"店铺" json:"author"`
	Store                   Store               `description:"店铺" json:"store"`
	CourseList             []Course             `description:"课程列表" json:"courseList"`

}

