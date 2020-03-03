/*
@Time : 2019/11/12 下午1:24 
@Author : zwcui
@Software: GoLand
*/
package util

import "math"

//四舍五入两位小数
func Rounding(value float64) float64 {
	return math.Trunc(value * 1e2 + 0.5) * 1e-2
}

//舍去第三位向后，保留两位小数
func NoRounding(value float64) float64 {
	return math.Trunc(value * 1e2) * 1e-2
}