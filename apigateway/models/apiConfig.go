/*
@Time : 2019/10/14 上午9:20 
@Author : zwcui
@Software: GoLand
*/
package models

type ApiConfig struct {
	Id       				int64				`description:"id" json:"id" xorm:"pk autoincr"`
	Description				string				`description:"描述" json:"description"`
	Url						string				`description:"接口跳转uri，多个,分隔" json:"url" xorm:"text"`
	Weight					string				`description:"uri对应的地址与权重，地址与权重用@@@分隔，多个,分隔" json:"weight" xorm:"text"`
	Status       			int					`description:"状态，1启用 2停止" json:"status" xorm:"notnull default 0"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}


//-----------------------------------------变量-------------------------------------------------
var UrlConfigList []UrlConfig

//-----------------------------------------结构体-------------------------------------------------

type UrlConfig struct {
	Description				string				`description:"描述" json:"description"`
	RequestURIArray 		[]string			`description:"请求地址uri" json:"requestURIArray"`
	//RequestRedirectArray 	[]RedirectWeight	`description:"请求跳转地址" json:"requestRedirectArray"`
	RequestRedirectArray 	[]string			`description:"请求跳转地址" json:"requestRedirectArray"`
}

//type RedirectWeight struct {
//	RedirectUrl				string				`description:"服务器地址" json:"redirectUrl"`
//	WeightStart				int					`description:"权重起始（包含）" json:"weightStart"`
//	WeightEnd				int					`description:"权重结束（不包含）" json:"weightEnd"`
//}