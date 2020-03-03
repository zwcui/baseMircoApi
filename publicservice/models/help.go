/*
@Time : 2019/3/21 上午11:50 
@Author : zwcui
@Software: GoLand
*/
package models

type Help struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	Name				string			`description:"名称" json:"name"`
	Level				int				`description:"级别 1为分类，2为章节，3为标题" json:"level"`
	LastId				int64			`description:"上级id" json:"lastId"`
	Content	     		string 			`description:"内容" json:"content" xorm:"text"`
	Created           	int64  			`description:"注册时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}


//----------------------------结构体--------------------------------------------
type HelpDetailListContainer struct {
	BaseListContainer
	HelpDetailList 		[]HelpDetail	`description:"帮助列表" json:"helpList"`
}

type HelpDetail struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	Classify			string			`description:"分类" json:"classify"`
	Chapter				string			`description:"章节" json:"chapter"`
	Title				string			`description:"标题" json:"title"`
	Level				int				`description:"级别 1为分类，2为章节，3为标题" json:"level"`
	LastId				int64			`description:"上级id" json:"lastId"`
	Content	     		string 			`description:"内容" json:"content" xorm:"text"`
	Created           	int64  			`description:"注册时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type AllHelp struct {
	ClassifyList		[]HelpClassify	`description:"分类" json:"helpList"`
}

type ClassifyListContainer struct {
	ClassifyList		[]HelpClassify	`description:"分类" json:"helpList"`
}

type ChapterListContainer struct {
	ChapterList			[]HelpChapter	`description:"章节" json:"chapterList"`
}

type TitleContainer struct {
	Title				HelpTitle		`description:"标题" json:"title"`
}

type HelpClassify struct {
	Classify			Help			`description:"分类" json:"classify"`
	ChapterList			[]HelpChapter	`description:"章节" json:"chapterList"`
}

type HelpChapter struct {
	Chapter				Help			`description:"章节" json:"chapter"`
	TitleList			[]HelpTitle		`description:"标题" json:"titleList"`
}

type HelpTitle struct {
	Title				Help			`description:"标题" json:"title"`
}


