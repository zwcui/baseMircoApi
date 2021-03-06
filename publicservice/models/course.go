/*
@Time : 2019/9/9 下午2:38 
@Author : lianwu
@File : courseType.go
@Software: GoLand
*/
package models


type Course struct {
	Id       			    int64			    `description:"id" json:"id" xorm:"pk autoincr"`
	Title					string				`description:"标题" json:"title" xorm:"index(IDX_title)"`
	Subtitle				string				`description:"副标题" json:"subtitle" xorm:"index(IDX_title)"`
	FirstCategoryId		    int64				`description:"一级分类" json:"firstCategoryId"`
	SecondCategoryId	    int64				`description:"二级分类" json:"secondCategoryId"`
	ParentId            	int64           	`description:"上级id" json:"parentId" `
	Level               	int			    	`description:"课程级别，1为课程，2为章节" json:"level" xorm:"notnull default 0"`
	ContentType				int					`description:"视频类型，1为1分钟短视频，2为5-10分钟小视频，3为系列课程，4为精品课程" json:"contentType" xorm:"notnull default 0"`
	Cover					string				`description:"封面" json:"cover"`
	Description				string				`description:"文字简介" json:"description" xorm:"varchar(500)"`
	Gallery				    string				`description:"图片简介" json:"gallery" xorm:"varchar(500)"`
	Tag						string				`description:"标签" json:"tag" xorm:"index(IDX_title)"`
	TotalStudy				int					`description:"浏览数(有多少人学习，目前统计的是该课程下所有视频的播放数)" json:"totalStudy" xorm:"notnull default 0"`
	TotalComment			int					`description:"总评论数" json:"totalComment" xorm:"notnull default 0"`
	TotalLike				int					`description:"点赞数" json:"totalLike" xorm:"notnull default 0"`
	TotalCollection			int					`description:"收藏数" json:"totalCollection" xorm:"notnull default 0"`
	TotalBuy				int					`description:"购买数" json:"totalBuy" xorm:"notnull default 0"`
	AuthorId				int64				`description:"作者id" json:"authorId"`
	UId       			    int64				`description:"上传视频用户id" json:"uId"`
	CheckStatus				int					`description:"审核状态,0为未发布（未提交审核），1为已发布（审核通过），2为发布中（审核中），3为发布失败（审核未通过）" json:"checkStatus" xorm:"notnull default 0"`
	ShowStatus				int					`description:"上架状态,0为未上架，1为已上架" json:"showStatus"`
	EditionId    			int64				`description:"版本id，对应修改版本的id" json:"editionId"`
	EditionType    			int					`description:"版本类型，0为修改版本，1为上架版本(上架才会多生成一条记录，下架则删除)" json:"editionType" xorm:"notnull default 0"`
	RefuseReason			string				`description:"驳回原因" json:"refuseReason"`
	UnloadReason			string				`description:"下架原因" json:"unloadReason"`
	Price					int					`description:"普通购买价格(整套课程售价)，单位分" json:"price"`
	VipPrice				int					`description:"vip购买价格，单位分" json:"vipPrice"`
	SortNo					string				 `description:"排序号" json:"sortNo"`
	IsOnStore			    int					`description:"是否允许上架店铺 1是 0否" json:"isOnStore" xorm:"notnull default 0"`
	StoreMoney				int					`description:"上架分佣价格，单位分" json:"storeMoney" xorm:"notnull default 0"`
	IsDistribution			int					`description:"是否允许分销 1是 0否" json:"isDistribution" xorm:"notnull default 0"`
	EditStatus				int					`description:"修改状态,0为已上传，1为已更新（审核通过）2为修改" json:"editStatus" xorm:"notnull default 0"`
	LockFlag				int					`description:"锁定状态，1是0否" json:"lockFlag" xorm:"notnull default 0"`
	TopSortNo				int					`description:"置顶排序" json:"topSortNo" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//课程分类
type CourseCategory struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	Name 			    string			`description:"分类名称" json:"name" xorm:"notnull "`
	ParentId            int64           `description:"上级分类id，第一级时为0 " json:"parentId" `
	HomeShowStatus  	int			    `description:"首页是否展示 0不展示 1展示" json:"homeShowStatus" xorm:"notnull default 0"`
	HomeSort 			string			`description:"首页排序 升序：数字小的在前面" json:"homeSort" xorm:"notnull "`
	ChannelSort 		string			`description:"频道排序 升序：数字小的在前面" json:"channelSort" xorm:"notnull "`
	Level               int			    `description:"课程分类级别 " json:"level" xorm:"notnull"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//课程点赞记录
type CourseLike struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId       			int64				`description:"用户id" json:"authorId"`
	CourseId       			int64				`description:"课程id" json:"courseId"`
}

//视频修改状态
type CourseEditHistory struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	CourseId       			int64				`description:"课程id" json:"courseId"`
	VideoStatus       		int  				`description:"视频内容修改状态0:未修改 1:已修改" json:"videoStatus" xorm:"notnull default 0"`
	TitleStatus       		int  				`description:"标题修改状态 	0:未修改 1:已修改" json:"titleStatus" xorm:"notnull default 0"`
	CoverStatus       		int  				`description:"封面修改状态    0:未修改 1:已修改" json:"coverStatus" xorm:"notnull default 0"`
	ClassifyStatus       	int  				`description:"分类状态      	0:未修改 1:已修改" json:"classifyStatus" xorm:"notnull default 0"`
	TagStatus       		int  				`description:"标签修改状态    0:未修改 1:已修改" json:"tagStatus" xorm:"notnull default 0"`
	PercentStatus       	int  				`description:"佣金比例修改状态 0:未修改 1:已修改" json:"percentStatus" xorm:"notnull default 0" `
	DescriptionStatus       int  				`description:"简介修改状态    0:未修改 1:已修改" json:"descriptionStatus" xorm:"notnull default 0"`
	SubtitleStatus			int				    `description:"副标题修改状态 	0:未修改 1:已修改" json:"subtitleStatus" xorm:"notnull default 0"`
	CategoryStatus			int				    `description:"分类修改状态 	0:未修改 1:已修改" json:"categoryStatus" xorm:"notnull default 0"`
	GalleryStatus			int				    `description:"图片简介修改状态 0:未修改 1:已修改" json:"galleryStatus" xorm:"notnull default 0"`
	ContentTypeStatus		int					`description:"视频类型修改状态    0:未修改 1:已修改" json:"contentTypeStatus" xorm:"notnull default 0"`
	PriceStatus       		int  				`description:"普通购买价格修改状态    0:未修改 1:已修改" json:"priceStatus" xorm:"notnull default 0"`
	VipPriceStatus			int					`description:"vip购买价格修改状态     0:未修改 1:已修改" json:"vipPriceStatus" xorm:"notnull default 0"`
	IsDistributionStatus	int					`description:"是否允许分销修改状态    0:未修改 1:已修改" json:"vipPriceStatus" xorm:"notnull default 0"`
	StoreMoneyStatus		int					`description:"分佣价格修改状态        0:未修改 1:已修改" json:"vipPriceStatus" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`


}

type CourseAndVideoCount struct {
	CourseAndCategory						    `description:"课程" xorm:"extends"`
	CourseVideoCount       int					`description:"课程包含的视频数" json:"courseVideoCount" `
}




//---------------------结构体-----------------------------

type CourseCategoryTree struct {
	Id       			int64			`description:"id" json:"id" xorm:"pk autoincr"`
	Name 			    string			`description:"分类名称" json:"name" xorm:"notnull "`
	ParentId            int64           `description:"上级分类id，第一级时为0 " json:"parentId" `
	HomeShowStatus  	int			    `description:"首页是否展示 0不展示 1展示" json:"homeShowStatus" xorm:"notnull default 0"`
	HomeSort 			string			`description:"首页排序 升序：数字小的在前面" json:"homeSort" xorm:"notnull "`
	ChannelSort 		string			`description:"频道排序 升序：数字小的在前面" json:"channelSort" xorm:"notnull "`
	Level               int			    `description:"课程分类级别 " json:"level" xorm:"notnull"`
	Created           	int64  			`description:"创建时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`

	ChildList			[]CourseCategoryTree    	`description:"子分类信息" json:"childList"`
}

type EditCourse struct {
	Id       			    int64			    `description:"id" json:"id" xorm:"pk autoincr"`
	Title					string				`description:"标题" json:"title"`
	Subtitle				string				`description:"副标题" json:"subtitle"`
	FirstCategoryId		    int64				`description:"一级分类" json:"firstCategoryId"`
	SecondCategoryId	    int64				`description:"二级分类" json:"secondCategoryId"`
	ParentId            	int64           	`description:"上级id" json:"parentId" `
	Level               	int			    	`description:"课程级别，1为课程，2为章节" json:"level" xorm:"notnull default 0"`
	ContentType				int					`description:"视频类型，1为1分钟短视频，2为5-10分钟小视频，3为系列课程，4为精品课程" json:"contentType" xorm:"notnull default 0"`
	Cover					string				`description:"封面" json:"cover"`
	Description				string				`description:"文字简介" json:"description" xorm:"varchar(500)"`
	Gallery				    string				`description:"图片简介" json:"gallery" xorm:"varchar(500)"`
	Tag						string				`description:"标签" json:"tag"`
	TotalStudy				int					`description:"浏览数(有多少人学习)" json:"totalStudy" xorm:"notnull default 0"`
	TotalComment			int					`description:"总评论数" json:"totalComment" xorm:"notnull default 0"`
	TotalLike				int					`description:"点赞数" json:"totalLike" xorm:"notnull default 0"`
	AuthorId				int64				`description:"作者id" json:"authorId"`
	UId       			    int64				`description:"上传视频用户id" json:"uId"`
	CheckStatus				int					`description:"审核状态,0为未发布（未提交审核），1为已发布（审核通过），2为发布中（审核中），3为发布失败（审核未通过）" json:"checkStatus" xorm:"notnull default 0"`
	ShowStatus				int					`description:"上架状态,0为未上架，1为已上架" json:"showStatus"`
	EditionId    			int64				`description:"版本id，对应修改版本的id" json:"editionId"`
	EditionType    			int					`description:"版本类型，0为修改版本，1为上架版本(上架才会多生成一条记录，下架则删除)" json:"editionType" xorm:"notnull default 0"`
	RefuseReason			string				`description:"驳回原因" json:"refuseReason"`
	UnloadReason			string				`description:"下架原因" json:"unloadReason"`
	Price					int					`description:"普通购买价格，单位分" json:"price"`
	VipPrice				int					`description:"vip购买价格，单位分" json:"vipPrice"`

	VideoStatus       		int  				`description:"视频内容修改状态0:未修改 1:已修改" json:"videoStatus" xorm:"notnull default 0"`
	TitleStatus       		int  				`description:"标题修改状态 	0:未修改 1:已修改" json:"titleStatus" xorm:"notnull default 0"`
	CoverStatus       		int  				`description:"封面修改状态    0:未修改 1:已修改" json:"coverStatus" xorm:"notnull default 0"`
	ClassifyStatus       	int  				`description:"分类状态      	0:未修改 1:已修改" json:"classifyStatus" xorm:"notnull default 0"`
	TagStatus       		int  				`description:"标签修改状态    0:未修改 1:已修改" json:"tagStatus" xorm:"notnull default 0"`
	PriceStatus       		int  				`description:"价格修改状态    0:未修改 1:已修改" json:"priceStatus" xorm:"notnull default 0"`
	PercentStatus       	int  				`description:"佣金比例修改状态 0:未修改 1:已修改" json:"percentStatus" xorm:"notnull default 0" `
	DescriptionStatus       int  				`description:"简介修改状态    0:未修改 1:已修改" json:"descriptionStatus" xorm:"notnull default 0"`
	SubtitleStatus			int				    `description:"副标题修改状态 	0:未修改 1:已修改" json:"subtitleStatus" xorm:"notnull default 0"`
	CategoryStatus			int				    `description:"分类修改状态 	0:未修改 1:已修改" json:"categoryStatus" xorm:"notnull default 0"`
	GalleryStatus			int				    `description:"图片简介修改状态 0:未修改 1:已修改" json:"galleryStatus" xorm:"notnull default 0"`

}

type CourseShort struct {
	Id       			    int64			    `description:"id" json:"id" xorm:"pk autoincr"`
	Title					string				`description:"标题" json:"title"`
	Subtitle				string				`description:"副标题" json:"subtitle"`
	FirstCategoryId		    int64				`description:"一级分类" json:"firstCategoryId"`
	SecondCategoryId	    int64				`description:"二级分类" json:"secondCategoryId"`
	ParentId            	int64           	`description:"上级id" json:"parentId" `
	Level               	int			    	`description:"课程级别，1为课程，2为章节" json:"level" xorm:"notnull default 0"`
	ContentType				int					`description:"视频类型，1为1分钟短视频，2为5-10分钟小视频，3为系列课程，4为精品课程" json:"contentType" xorm:"notnull default 0"`
	Cover					string				`description:"封面" json:"cover"`
	Description				string				`description:"文字简介" json:"description" xorm:"varchar(500)"`
	Gallery				    string				`description:"图片简介" json:"gallery" xorm:"varchar(500)"`
	Tag						string				`description:"标签" json:"tag"`
	TotalStudy				int					`description:"浏览数(有多少人学习)" json:"totalStudy" xorm:"notnull default 0"`
	TotalComment			int					`description:"总评论数" json:"totalComment" xorm:"notnull default 0"`
	TotalLike				int					`description:"点赞数" json:"totalLike" xorm:"notnull default 0"`
	TotalCollection			int					`description:"收藏数" json:"totalCollection"`
	AuthorId				int64				`description:"作者id" json:"authorId"`
	UId       			    int64				`description:"上传视频用户id" json:"uId"`
	Price					int					`description:"普通购买价格，单位分" json:"price"`
	VipPrice				int					`description:"vip购买价格，单位分" json:"vipPrice"`
	StoreMoney				int					`description:"上架分佣价格，单位分" json:"storeMoney" xorm:"notnull default 0"`
	SortNo					string				`description:"排序号" json:"sortNo"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type CourseCategoryH5ListContainer struct {
	CourseCategoryList		    []CourseCategoryTree		`description:"课程分类信息" json:"courseCategoryList"`
}

type CourseCategoryListContainer struct {
	CourseCategoryList		    []CourseCategory		 	`description:"课程分类信息" json:"courseCategoryList"`
}

//课程列表
type CourseListContainer struct {
	BaseListContainer
	CourseList					[]CourseAndVideoCount					`description:"课程列表" json:"courseList"`
}

type CourseInfo struct {
	CourseShort										        `description:"课程" xorm:"extends"`
	LikeStatus       		    int				            `description:"点赞状态 0:否 1:是" json:"likeStatus"`
	CollectionStatus            int				            `description:"收藏状态 0:否 1:是" json:"collectionStatus"`
	BuyStatus                   int				            `description:"是否购买 0:否 1:是" json:"buyStatus"`
	StoreStatus                 int				            `description:"店铺中是否存在该课程 0:不存在 1:存在" json:"storeStatus"`

}

type CourseInfoForBackend struct {
	Course						Course						`description:"课程" json:"course"`
}

type CourseAndCategory struct {
	Course						    						`description:"课程" xorm:"extends"`
	Name                       string                       `description:"课程分类" json:"name"`

}

//章节
type CourseVideoInfoForBackend struct {
	Course						Course						`description:"课程" json:"course"`
	VideoList					[]Video						`description:"视频列表" json:"videoList"`
}

type CourseVideoInfoForBackendListContainer struct {
	CourseList					[]CourseVideoInfoForBackend `description:"课程" json:"courseList"`
}

type EditCourseContainer struct {
	EditCourse							                    `description:"课程修改内容" xorm:"extends"`
}

type CourseInfoContainer struct {
	CourseInfo				   CourseInfo			        `description:"课程修改内容" json:"course"`
}

type EditAndShowVideo map[int64]int64

type SortMap map[int]int
