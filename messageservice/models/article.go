/*
@Time : 2019/4/24 上午10:58 
@Author : zwcui
@Software: GoLand
*/
package models

//发布内容
type Article struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	UId       				int64				`description:"发布者uId" json:"uId"`
	Title      				string				`description:"标题" json:"title"`
	Abstract   				string				`description:"摘要" json:"abstract"`
	Cover   				string				`description:"封面" json:"cover"`
	Tag		   				string				`description:"标签，多个,隔开" json:"tag"`
	Content		   			string				`description:"内容" json:"content" xorm:"text"`
	ContentUrl	   			string				`description:"内容地址" json:"contentUrl"`
	BannerUrl	   			string				`description:"海报地址" json:"bannerUrl"`
	ReadCount				int					`description:"查看量" json:"readCount" xorm:"notnull default 0"`
	ForwardCount			int					`description:"转发人数" json:"forwardCount" xorm:"notnull default 0"`
	PluginType	   			string				`description:"插件类型，1为报名" json:"pluginType"`
	Status       			int					`description:"状态，0草稿未发布 1已发布" json:"status" xorm:"notnull default 0"`
	Qrcode		   			string				`description:"链接二维码" json:"qrcode"`
	BannerQrcode   			string				`description:"海报二维码" json:"bannerQrcode"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//内容报名
type ArticleSignUp struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	ArticleId 				int64				`description:"内容id" json:"articleId"`
	Name      				string				`description:"姓名" json:"name"`
	Phone      				string				`description:"手机号" json:"phone"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//内容分享层级
type ArticleShare struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	LastShareId				int64				`description:"上级分享id" json:"lastShareId"`
	Openid					string				`description:"分享人" json:"openid"`
	LastOpenid				string				`description:"上级分享人" json:"lastOpenid"`
	ArticleId 				int64				`description:"内容id" json:"articleId"`
	Level					int					`description:"活动层级，最上面为1级" json:"level"`
	DirectNextShareNum		int					`description:"直接下级分享人数" json:"directNextShareNum" xorm:"notnull default 0"`
	AllNextShareNum			int					`description:"所有下级分享人数" json:"allNextShareNum" xorm:"notnull default 0"`
	Status					int					`description:"状态，1为正常" json:"status"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}


//-------------------------结构体--------------------------------
type ArticleDetail struct {
	Article					`description:"内容" xorm:"extends"`
	AuthorName      		string				`description:"作者" json:"authorName"`
}

type ArticleListContainer struct {
	BaseListContainer
	ArticleList  			[]ArticleDetail 	`description:"内容列表" json:"articleList"`
}

type ArticleContainer struct {
	Article					Article				`description:"内容" json:"article"`
}

type ArticleShareListContainer struct {
	ArticleShareList  		[]ArticleShareDetail `description:"分享列表" json:"articleShareList"`
}

type ArticleShareDetail struct {
	ArticleShare							`description:"分享" xorm:"extends"`
	Nickname				string			`description:"用户的昵称" json:"nickname"`
	Sex						int				`description:"用户的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	Headimgurl				string			`description:"用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
}
