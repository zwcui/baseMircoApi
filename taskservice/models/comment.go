/*
@Time : 2019/6/3 上午10:36 
@Author : zwcui
@Software: GoLand
*/
package models

type Comment struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId				int64				`description:"回复人id" json:"authorId"`
	CommentType				int					`description:"评论类型，1为视频评论" json:"commentType"`
	CommentedId				int64				`description:"被评论对象id" json:"commentedId"`
	Content					string				`description:"评论内容" json:"content" xorm:"varchar(450)"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}


//--------------------------------结构体---------------------------------------

type CommentDetail struct {
	Comment					`description:"视频内容" xorm:"extends"`
	Nickname				string				`description:"用户的昵称" json:"nickname"`
	Sex						int					`description:"用户的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	City					string				`description:"用户所在城市" json:"city"`
	Province				string				`description:"用户所在省份" json:"province"`
	Country					string				`description:"用户所在国家" json:"country"`
	Headimgurl				string				`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
}

type CommentDetailListContainer struct {
	BaseListContainer
	CommentDetailList  		[]CommentDetail 	`description:"视频列表" json:"commentList"`
}
