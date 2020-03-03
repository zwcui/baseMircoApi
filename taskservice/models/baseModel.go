/*
@Time : 2019/3/11 上午11:18 
@Author : zwcui
@Software: GoLand
*/
package models

type BaseListContainer struct {
	TotalCount 			int64 			`description:"总数" json:"totalCount"`
	PageNum    			int   			`description:"当前页数" json:"pageNum"`
	PageTime   			int64 			`description:"查询时间" json:"pageTime"`
}

