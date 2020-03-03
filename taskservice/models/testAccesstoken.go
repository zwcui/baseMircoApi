/*
@Time : 2019/10/31 上午9:45 
@Author : lianwu
@File : testAccesstoken.go
@Software: GoLand
*/
package models


//记录测试数据
type TestAccessToken struct {

	Id       			int64			`description:"Id" json:"d" xorm:"pk autoincr"`
	TestTime			string			`description:"调用时间" json:"testTime"`
	ExpeireTime			string			`description:"过期时间" json:"expeireTime"`
	Status              string			`description:"状态" json:"status"`
	ErrorAndCount		string			`description:"错误占比" json:"errorAndCount"`
	Count			    int			    `description:"调用总次数" json:"count"`
	ErrorCount			int			    `description:"错误总次数" json:"errorCount"`
	Accesstoken 		string			`description:"accesstoken" json:"accesstoken"`
	ErrorDetail 		string			`description:"错误信息" json:"errorDetail"`

	Created           	int64  			`description:"注册时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}
