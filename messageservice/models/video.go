/*
@Time : 2019/5/23 上午9:24 
@Author : zwcui
@Software: GoLand
*/
package models

//视频信息
type Video struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Title					string				`description:"标题" json:"title"`
	Cover					string				`description:"封面" json:"cover"`
	Classify				string				`description:"分类" json:"classify"`
	Tag						string				`description:"标签" json:"tag"`
	Price					int					`description:"价格，单位分" json:"price" xorm:"notnull default 0"`
	Star					int					`description:"星级" json:"star"`
	FreeTime				int					`description:"免费观看结束时间，单位秒" json:"freeTime"`
	ForwardTime				int					`description:"转发免费观看结束时间，单位秒" json:"forwardTime"`
	TotalDuration			int					`description:"总时长，单位秒" json:"totalDuration"`
	FileId					string				`description:"未转码原始上传fileId" json:"fileId"`
	OriginalUrl				string				`description:"未加防盗链未转码原始上传完整url" json:"originalUrl"`
	Url						string				`description:"加入防盗链未转码原始上传完整url" json:"url"`
	FreeUrl					string				`description:"加入防盗链未转码原始上传免费url，如果切分未完成则前端使用防盗链试看" json:"freeUrl"`
	ForwardUrl				string				`description:"加入防盗链未转码原始上传转发url，如果切分未完成则前端使用防盗链试看" json:"forwardUrl"`

	OriginalUrl640			string				`description:"未加防盗链完整url（分辨率640）" json:"originalUrl640"`
	Url640					string				`description:"地址640（废弃）" json:"url640"`
	FreeUrl640				string				`description:"未加防盗链免费url（分辨率640）" json:"freeUrl640"`
	ForwardUrl640			string				`description:"未加防盗链转发url（分辨率640）" json:"forwardUrl640"`
	FreeFileId640			string				`description:"免费切分fileId（分辨率640）" json:"freeFileId640"`
	ForwardFileId640		string				`description:"转发切分fileId（分辨率640）" json:"forwardFileId640"`
	FullFileId640			string				`description:"完整fileId（分辨率640）" json:"fullFileId640"`

	OriginalUrl1280			string				`description:"未加防盗链完整url（分辨率1280）" json:"originalUrl1280"`
	Url1280					string				`description:"地址1280（废弃）" json:"url1280"`
	FreeUrl1280				string				`description:"未加防盗链免费url（分辨率1280）" json:"freeUrl1280"`
	ForwardUrl1280			string				`description:"未加防盗链转发url（分辨率1280）" json:"forwardUrl1280"`
	FreeFileId1280			string				`description:"免费切分fileId（分辨率1280）" json:"freeFileId1280"`
	ForwardFileId1280		string				`description:"转发切分fileId（分辨率1280）" json:"forwardFileId1280"`
	FullFileId1280			string				`description:"完整fileId（分辨率1280）" json:"fullFileId1280"`

	OriginalUrl1920			string				`description:"未加防盗链完整url（分辨率1920）" json:"originalUrl1920"`
	Url1920					string				`description:"地址1920（废弃）" json:"url1920"`
	FreeUrl1920				string				`description:"未加防盗链免费url（分辨率1920）" json:"freeUrl1920"`
	ForwardUrl1920			string				`description:"未加防盗链转发url（分辨率1920）" json:"forwardUrl1920"`
	FreeFileId1920			string				`description:"免费切分fileId（分辨率1920）" json:"freeFileId1920"`
	ForwardFileId1920		string				`description:"转发切分fileId（分辨率1920）" json:"forwardFileId1920"`
	FullFileId1920			string				`description:"完整fileId（分辨率1920）" json:"fullFileId1920"`

	AudioUrl128				string				`description:"音频完整url（清晰度128）" json:"audioUrl128"`
	AudioFreeUrl128			string				`description:"音频免费url（清晰度128）" json:"audioFreeUrl1920"`
	AudioForwardUrl128		string				`description:"音频转发免费url（清晰度128）" json:"audioForwardUrl1920"`
	AudioFreeFileId128		string				`description:"音频免费切分fileId（清晰度128）" json:"audioFreeFileId128"`
	AudioForwardFileId128	string				`description:"音频转发切分fileId（清晰度128）" json:"audioForwardFileId128"`
	AudioFullFileId128		string				`description:"音频完整fileId（清晰度128）" json:"audioFullFileId128"`

	AuthorId				int64				`description:"作者id" json:"authorId"`
	Description				string				`description:"简介" json:"description"`
	Advertisement			string				`description:"广告" json:"advertisement"`
	HasWatermark			int					`description:"是否有水印，1是0否" json:"hasWatermark" xorm:"notnull default 0"`
	Watermark				string				`description:"水印" json:"watermark"`
	AutoCheck				int					`description:"腾讯云鉴黄、恐、政，0为未鉴别，1为通过，2为未通过" json:"autoCheck" xorm:"notnull default 0"`
	AutoCheckTaskId			string				`description:"腾讯云鉴黄、恐、政" json:"autoCheckTaskId"`
	AutoCheckResult			string				`description:"腾讯云鉴黄、恐、政" json:"autoCheckResult"`
	Width					int					`description:"宽" json:"width"`
	Height					int					`description:"高" json:"height"`
	TotalPlay				int					`description:"播放数" json:"totalPlay" xorm:"notnull default 0"`
	TotalView				int					`description:"浏览数" json:"totalView" xorm:"notnull default 0"`
	TotalComment			int					`description:"评论数" json:"totalComment" xorm:"notnull default 0"`
	TotalLike				int					`description:"点赞数" json:"totalLike" xorm:"notnull default 0"`
	TotalShare				int					`description:"分享数" json:"totalShare" xorm:"notnull default 0"`
	CheckStatus				int					`description:"审核状态,0为未发布（未提交审核），1为已发布（审核通过），2为发布中（审核中），3为发布失败（审核未通过）" json:"checkStatus" xorm:"notnull default 0"`
	ShowStatus				int					`description:"上架状态,0为未上架，1为已上架" json:"showStatus" xorm:"notnull default 0"`
	FirstLevelPercent		int					`description:"1级分销商分佣比例，单位百分之" json:"firstLevelPercent" xorm:"notnull default 0"`
	SecondLevelPercent		int					`description:"2级分销商分佣比例，单位百分之" json:"secondLevelPercent" xorm:"notnull default 0"`
	EditionId    			int64				`description:"版本id，对应修改版本的id" json:"editionId"`
	EditionType    			int					`description:"版本类型，0为修改版本，1为上架版本(上架才会多生成一条记录，下架则删除)" json:"editionType" xorm:"notnull default 0"`
	RefuseReason			string				`description:"驳回原因" json:"refuseReason"`
	UnloadReason			string				`description:"下架原因" json:"unloadReason"`
	SubmitTime          	int64  				`description:"投稿时间" json:"submitTime"`

	FreeStatus              int                 `description:"免费状态 0:需要购买vip观看完整版 1:免费观看完整版" json:"freeStatus" xorm:"notnull default 0"`
	CourseId                int64      			`description:"章节id" json:"courseId"`
	EditStatus				int					`description:"修改状态,0为已上传，1为已更新（审核通过）2为修改" json:"editStatus" xorm:"notnull default 0"`

	SortNo					string				`description:"排序号" json:"sortNo"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}


//视频修改状态
type VideoEditHistory struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	VideoId       			int64				`description:"视频id" json:"videoId"`
	VideoStatus       		int  				`description:"视频内容修改状态0:未修改 1:已修改" json:"videoStatus" xorm:"notnull default 0"`
	TitleStatus       		int  				`description:"标题修改状态 	0:未修改 1:已修改" json:"titleStatus" xorm:"notnull default 0"`
	CoverStatus       		int  				`description:"封面修改状态    0:未修改 1:已修改" json:"coverStatus" xorm:"notnull default 0"`
	ClassifyStatus       	int  				`description:"分类状态      	0:未修改 1:已修改" json:"classifyStatus" xorm:"notnull default 0"`
	TagStatus       		int  				`description:"标签修改状态    0:未修改 1:已修改" json:"tagStatus" xorm:"notnull default 0"`
	PriceStatus       		int  				`description:"价格修改状态    0:未修改 1:已修改" json:"priceStatus" xorm:"notnull default 0"`
	PercentStatus       	int  				`description:"佣金比例修改状态 0:未修改 1:已修改" json:"percentStatus" xorm:"notnull default 0" `
	DescriptionStatus       int  				`description:"简介修改状态    0:未修改 1:已修改" json:"descriptionStatus" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`


}
//视频浏览记录
type VideoViewHistory struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId       			int64				`description:"用户id" json:"authorId"`
	VideoId       			int64				`description:"视频id" json:"videoId"`
	CourseId       			int64				`description:"课程id" json:"courseId"`
	WatchTime      			int					`description:"观看时间，单位秒" json:"watchTime"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`

}

//视频点赞记录
type VideoLike struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId       			int64				`description:"用户id" json:"authorId"`
	VideoId       			int64				`description:"视频id" json:"videoId"`
	LikeStatus       		int				    `description:"点赞状态 0:否 1:是" json:"likeStatus"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`

}

//举报记录
type Report struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	AuthorId       			int64				`description:"用户id" json:"authorId"`
	ValueId       			int64				`description:"type为1时:视频id, type为2时:评论id" json:"valueId"`
	Type       			    int				    `description:"举报类型 1:视频 2:评论" json:"type"`
	Content					string				`description:"举报内容" json:"content"`
	Classify				int				    `description:"举报内容类别 1:违法违禁, 2:色情, 3:低俗, 4:赌博诈骗, 5:血腥暴力, 6:人生攻击, 7:与其他视频相同, 8:不良封面标题, 9;青少年不良信息, 10:其他 " json:"classify"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`

}


//-----------------------------结构体-------------------------------------
type VideoDetail struct {
	VideoInfo			     `description:"视频基本信息" xorm:"extends"`
	Url				        string				`description:"地址(此处是根据用户是否购买转发而返回对应的url,会有三种情况)" json:"Url"`
	Url640			        string				`description:"地址640(此处是根据用户是否购买转发而返回对应的url,会有三种情况)" json:"Url640"`
	Url1280			        string				`description:"地址1280(此处是根据用户是否购买转发而返回对应的url,会有三种情况)" json:"Url1280"`
	Url1920			        string				`description:"地址1920(此处是根据用户是否购买转发而返回对应的url,会有三种情况)" json:"Url1920"`

	UrlStatus				int  				`description:"地址状态 0:已购买 1:已转发或关注公众号 2:未购买未转发" json:"urlStatus"`
	LikeStatus				int  				`description:"点赞状态 0:未点赞 1:已点赞" json:"likeStatus"`
	BuyNum					int  				`description:"购买人数" json:"buyNum"`
	BuyDate					int64  				`description:"购买日期" json:"buyDate"`
	Subscribe       	    int				    `description:"用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。" json:"subscribe"`

	Nickname				string				`description:"up主的昵称" json:"nickname"`
	Sex						int					`description:"up主的性别，值为1时是男性，值为2时是女性，值为0时是未知" json:"sex"`
	City					string				`description:"up主所在城市" json:"city"`
	Province				string				`description:"up主所在省份" json:"province"`
	Country					string				`description:"up主所在国家" json:"country"`
	Headimgurl				string				`description:"up主头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。" json:"headimgurl"`
}


//视频详情返回的基本信息结构体
type VideoInfo struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Title					string				`description:"标题" json:"title"`
	Cover					string				`description:"封面" json:"cover"`
	Classify				string				`description:"分类" json:"classify"`
	Tag						string				`description:"标签" json:"tag"`
	Price					int					`description:"价格，单位分" json:"price"`
	Star					int					`description:"星级" json:"star"`
	TotalDuration			int					`description:"总时长，单位秒" json:"totalDuration"`
	FreeTime				int					`description:"免费观看结束时间，单位秒" json:"freeTime"`
	ForwardTime				int					`description:"转发免费观看结束时间，单位秒" json:"forwardTime"`
	AuthorId				int64				`description:"作者id" json:"authorId"`
	Description				string				`description:"简介" json:"description"`
	Advertisement			string				`description:"广告" json:"advertisement"`
	Width					int					`description:"宽" json:"width"`
	Height					int					`description:"高" json:"height"`
	TotalPlay				int					`description:"播放数" json:"totalPlay"`
	TotalComment			int					`description:"评论数" json:"totalComment"`
	TotalLike				int					`description:"点赞数" json:"totalLike"`
	TotalShare				int					`description:"分享数" json:"totalShare"`
	TotalView				int					`description:"浏览数" json:"totalView"`

	RefuseReason			string				`description:"驳回原因" json:"refuseReason"`
	UnloadReason			string				`description:"下架原因" json:"unloadReason"`
	CheckStatus				int					`description:"审核状态,0为未发布（未提交审核），1为已发布（审核通过），2为发布中（审核中），3为发布失败（审核未通过）" json:"checkStatus" `
	ShowStatus				int					`description:"上架状态,0为未上架，1为已上架" json:"showStatus"`
	FirstLevelPercent		int					`description:"1级分销商分佣比例，单位百分之" json:"firstLevelPercent"`
	SecondLevelPercent		int					`description:"2级分销商分佣比例，单位百分之" json:"secondLevelPercent"`
	EditionId    			int64				`description:"版本id，对应修改版本的id" json:"editionId"`
	EditionType    			int					`description:"版本类型，0为修改版本，1为上架版本(上架才会多生成一条记录，下架则删除)" json:"editionType" `
	AutoCheck				int					`description:"腾讯云鉴黄、恐、政，0为未鉴别，1为通过，2为未通过" json:"autoCheck" `

	OriginalUrl				string				`description:"未key加密地址" json:"originalUrl"`
	OriginalUrl640			string				`description:"未key加密地址640" json:"originalUrl640"`
	OriginalUrl1280			string				`description:"未key加密地址1280" json:"originalUrl1280"`
	OriginalUrl1920			string				`description:"未key加密地址1920" json:"originalUrl1920"`
	FileId					string				`description:"fileId" json:"fileId"`
	FreeStatus              int                 `description:"免费状态 0:需要购买vip观看完整版 1:免费观看完整版" json:"freeStatus" xorm:"notnull default 0"`
	SortNo					string				`description:"排序号" json:"sortNo"`
	Created           		int64  				`description:"创建时间" json:"created"`
	SubmitTime           	int64  				`description:"投稿时间" json:"submitTime"`
}

//查看视频修改内容返回的结构体
type EditVideo struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Title					string				`description:"标题" json:"title"`
	Cover					string				`description:"封面" json:"cover"`
	Classify				string				`description:"分类" json:"classify"`
	Tag						string				`description:"标签" json:"tag"`
	Price					int					`description:"价格，单位分" json:"price"`

	FreeTime				int					`description:"免费观看结束时间，单位秒" json:"freeTime"`
	ForwardTime				int					`description:"转发免费观看结束时间，单位秒" json:"forwardTime"`
	TotalDuration			int					`description:"总时长，单位秒" json:"totalDuration"`
	Width					int					`description:"宽" json:"width"`
	Height					int					`description:"高" json:"height"`
	FileId					string				`description:"fileId" json:"fileId"`
	Description				string				`description:"简介" json:"description"`
	Advertisement			string				`description:"广告" json:"advertisement"`
	HasWatermark			int					`description:"是否有水印，1是0否" json:"hasWatermark" `
	Watermark				string				`description:"水印" json:"watermark"`

	FirstLevelPercent		int					`description:"1级分销商分佣比例，单位百分之" json:"firstLevelPercent"`
	SecondLevelPercent		int					`description:"2级分销商分佣比例，单位百分之" json:"secondLevelPercent"`
	Url						string				`description:"地址" json:"url"`
	OriginalUrl				string				`description:"未key加密地址" json:"originalUrl"`
	OriginalUrl640			string				`description:"未key加密地址640" json:"originalUrl640"`
	OriginalUrl1280			string				`description:"未key加密地址1280" json:"originalUrl1280"`
	OriginalUrl1920			string				`description:"未key加密地址1920" json:"originalUrl1920"`
	AuthorId				int64				`description:"作者id" json:"authorId"`

	VideoStatus       		int  				`description:"视频内容修改状态0:未修改 1:已修改" json:"videoStatus"`
	TitleStatus       		int  				`description:"标题修改状态 	0:未修改 1:已修改" json:"titleStatus"`
	CoverStatus       		int  				`description:"封面修改状态    0:未修改 1:已修改" json:"coverStatus"`
	ClassifyStatus       	int  				`description:"分类状态      	0:未修改 1:已修改" json:"classifyStatus"`
	TagStatus       		int  				`description:"标签修改状态    0:未修改 1:已修改" json:"tagStatus"`
	PriceStatus       		int  				`description:"价格修改状态    0:未修改 1:已修改" json:"priceStatus"`
	PercentStatus       	int  				`description:"佣金比例修改状态 0:未修改 1:已修改" json:"percentStatus"`
	DescriptionStatus       int  				`description:"简介修改状态    0:未修改 1:已修改" json:"descriptionStatus"`
	IsTest                  int  				`description:"视频作者是否为测试用户    0:否 1:是" json:"isTest"`

}
type EditVideoContainer struct {
	EditVideo							        `description:"视频修改内容" xorm:"extends"`
}

type VideoDetailContainer struct {
	VideoDetailList				[]VideoDetail			`description:"视频详情" json:"video"`
}


type VideoDetailListContainer struct {
	BaseListContainer
	VideoDetailList  		[]VideoDetail 		`description:"视频列表" json:"videoList"`
}
//视频浏览记录返回机构体
type ViewVideo struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	VideoId       			int64				`description:"视频id" json:"videoId"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	CourseId       			int64				`description:"课程id" json:"courseId"`
	WatchTime      			int					`description:"观看时间，单位秒" json:"watchTime"`
	Title					string				`description:"标题" json:"title"`
	Cover					string				`description:"封面" json:"cover"`
	TotalDuration			int					`description:"总时长，单位秒" json:"totalDuration"`
	ShowStatus				int					`description:"上架状态,0为未上架，1为已上架" json:"showStatus" xorm:"notnull default 0"`

}

type ViewVideoListContainer struct {
	BaseListContainer
	ViewVideoList  		    []ViewVideo 	    `description:"浏览视频列表" json:"viewVideoList"`
}


type ReportListContainer struct {
	BaseListContainer
	ReportList  		    []Report 	         `description:"举报记录列表" json:"viewVideoList"`
}

//预览视频
type PreviewVideoContainer struct {
	Url						string				`description:"加入防盗链地址" json:"url"`
}


type AddLikeContainer struct {
	TotalLike					int 			`description:"点赞数" json:"totalLike"`
}

//不包含视频地址的信息
type VideoShort struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Title					string				`description:"标题" json:"title"`
	Cover					string				`description:"封面" json:"cover"`
	Classify				string				`description:"分类" json:"classify"`
	Tag						string				`description:"标签" json:"tag"`
	Price					int					`description:"价格，单位分" json:"price" xorm:"notnull default 0"`
	Star					int					`description:"星级" json:"star"`
	TotalDuration			int					`description:"总时长，单位秒" json:"totalDuration"`
	AuthorId				int64				`description:"作者id" json:"authorId"`
	Width					int					`description:"宽" json:"width"`
	Height					int					`description:"高" json:"height"`
	TotalPlay				int					`description:"播放数" json:"totalPlay" xorm:"notnull default 0"`
	TotalView				int					`description:"浏览数" json:"totalView" xorm:"notnull default 0"`
	TotalComment			int					`description:"评论数" json:"totalComment" xorm:"notnull default 0"`
	TotalLike				int					`description:"点赞数" json:"totalLike" xorm:"notnull default 0"`
	TotalShare				int					`description:"分享数" json:"totalShare" xorm:"notnull default 0"`
	SubmitTime          	int64  				`description:"投稿时间" json:"submitTime"`
	FreeStatus              int                 `description:"免费状态 0:需要购买vip观看完整版 1:免费观看完整版" json:"freeStatus" xorm:"notnull default 0"`
	SortNo					string				`description:"排序号" json:"sortNo"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

type VideoListContainer struct {
	BaseListContainer
	VideoList			[]VideoDetail		`description:"视频列表" json:"videoList"`
}

//-----------------------------文件服务器响应-------------------------------

type FileServiceVideoClipResponse struct {
	Code        			int    `json:"code"`
	Description 			string `json:"description"`
	Data        			struct {
				Urls 		[]string `json:"urls"`
	} 	`json:"data"`
}

