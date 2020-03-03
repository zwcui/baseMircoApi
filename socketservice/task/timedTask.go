package task

import (
	"jingting_server/socketservice/controllers"
	"time"
	"github.com/robfig/cron"
)

//初始化定时任务
func init() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.NewWithLocation(location)

	//每分钟检查心跳失效的socket连接
	c.AddFunc("@every 1m", func() {
		controllers.CheckSocketHeartbeat()
	})
	c.Start()
}

