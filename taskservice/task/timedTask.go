package task

import (
	"time"
	"github.com/robfig/cron"
)

//初始化定时任务
func init() {
	location, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.NewWithLocation(location)

	//每分钟检查component_access_token是否过期，过期则重新获取
	c.AddFunc("@every 1m", func() {
		RequestComponentAccessToken()
	})

	//每分钟检查pre_auth_code是否过期，过期则重新获取
	c.AddFunc("@every 1m", func() {
		RequestPreAuthCode()
	})

	//每分钟检查公众号access_token是否过期，过期则重新获取
	//c.AddFunc("@every 1m", func() {
	//	requestAuthorizerAccessToken()
	//})

	//定时任务，每小时检查活动是否结束
	c.AddFunc("@every 30m", func() {
		checkActivityIsEnd()
	})

	//定时任务，每分钟检查是否需要发送模板消息
	c.AddFunc("@every 1m", func() {
		sendTemplateMessage()
	})

	//定时任务，每天检查删除已发送并超过60天的模板消息
	c.AddFunc("@every 24h", func() {
		deleteTemolateMessage()
	})

	//定时任务，分佣每日结算
	//c.AddFunc("0 0 22 * * ? ", func() {
	//	settleAccountAmountDaily()
	//})

	//定时任务，分佣每周结算
	//c.AddFunc("0 0 22 * * 7 ", func() {
	//	settleAccountAmountWeekly()
	//})

	//定时任务，分佣每月结算
	//c.AddFunc("0 0 22 25 * *", func() {
	//	settleAccountAmountMonthly()
	//})

	//定时任务，每天处理转账失败的提现
	//c.AddFunc("0 0 23 * * ? ", func() {
	//	handleFailedCashOut()
	//})

	//定时任务，每10分钟检查会员是否过期
	c.AddFunc("@every 10m", func() {
		checkAuthorVIPEndTime()
	})

	//每分钟检查辣课公众号access_token是否过期，过期则重新获取
	c.AddFunc("@every 1m", func() {
		checkGZHAccessToken()
	})

	//每分钟检查辣课公众号access_token是否有效
	c.AddFunc("@every 1m", func() {
		testAccesstoken()
	})

	//每分钟检查直播预告是否过期
	c.AddFunc("@every 1m", func() {
		checkExpireLive()
	})

	//每分钟去腾讯云查询录制的视频（废弃，根据腾讯云回调获取）
	//c.AddFunc("@every 1m", func() {
	//	getLiveVideo()
	//})

	c.Start()
}

var Count, ErrorCount int
