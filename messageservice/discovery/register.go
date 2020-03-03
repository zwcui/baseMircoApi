/*
@Time : 2019/9/4 上午9:20 
@Author : zwcui
@Software: GoLand
*/
package discovery

import "time"

type RegisterInfo struct {
	Host           string
	Port           int
	ServiceName    string
	UpdateInterval time.Duration
}

type Register interface {
	Register(info RegisterInfo) error
	DeRegister(info RegisterInfo) error
}

