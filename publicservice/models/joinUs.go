/*
@Time : 2019/9/16 下午6:25 
@Author : zwcui
@Software: GoLand
*/
package models

//加盟开课
type JoinUs struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	UId		      			int64				`description:"管理员id" json:"uId"`
	Classify      			string				`description:"分类" json:"classify"`
	Title	      			string				`description:"课程名称" json:"title"`
	Cover	      			string				`description:"封面" json:"cover"`
	VideoUrl      			string				`description:"介绍视频地址" json:"videoUrl"`
	Url      			    string				`description:"介绍视频地址(安卓用)" json:"url"`
	ShowType      			int					`description:"展示类型，1为仅代理可看（伯乐用户 / 代理商），2为VIP会员可看，3为所有用户可看" json:"showType"`
	TotalDuration			int					`description:"总时长，单位秒" json:"totalDuration"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//------------------------------------结构体------------------------------------------
type JoinUsListContainer struct {
	BaseListContainer
	JoinUsList				[]JoinUs			`description:"加盟开课栏目列表" json:"joinUsList"`
}