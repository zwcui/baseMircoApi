/*
@Time : 2019/6/26 上午9:54
@Author : zwcui
@Software: GoLand
*/
package task

import (
	"jingting_server/taskservice/models"
	"jingting_server/taskservice/base"
	"strings"
	"jingting_server/taskservice/util"
	"github.com/astaxie/beego"
	"strconv"
)

//分佣每日结算
func settleAccountAmountDaily(){
	util.Logger.Info("定时任务，每天结算分佣金额")
	settle(2)
}

//分佣每周结算
func settleAccountAmountWeekly(){
	util.Logger.Info("定时任务，每周结算分佣金额")
	settle(3)
}

//分佣每月结算
func settleAccountAmountMonthly(){
	util.Logger.Info("定时任务，每月结算分佣金额")
	settle(4)
}

//分佣结算
func settle(settlementType int){
	var accountList []models.AuthorAccount
	base.DBEngine.Table("author_account").Where("settlement_type=?", settlementType).And("unsettled_amount > 0").Find(&accountList)
	ip := "#ip#"
	if beego.BConfig.RunMode == "prod" {
		ip = "#ip#"
	}

	for _, account := range accountList {
		session := base.DBEngine.NewSession()
		defer session.Close()
		err := session.Begin()
		if err != nil {
			util.Logger.Info("settle" + strconv.Itoa(settlementType) + " session.Begin() err:" + err.Error())
			session.Close()
			continue
		}

		//每日结算
		var cashOut models.CashOut
		cashOut.TransferType = 2
		cashOut.SettlementType = account.SettlementType
		cashOut.AuthorId = account.AuthorId
		cashOut.Amount = account.UnsettledAmount
		cashOut.Status = 0
		_, err = session.Table("cash_out").InsertOne(&cashOut)

		response, err := models.WeChatTransfers(ip, &cashOut, session)
		if err != nil {
			errCodeDes := strings.Split(err.Error(), ":")[0]
			errCode := strings.Split(err.Error(), ":")[1]
			util.Logger.Info("settle" + strconv.Itoa(settlementType) + " 微信转账失败 err = " + errCodeDes)
			cashOut.ErrorCode = errCode
			cashOut.ErrorCodeDes = errCodeDes
			models.TransfersFail(&cashOut, &account, session)
		} else {
			cashOut.PaymentNo = response.PaymentNo
			cashOut.TransferTime = response.PaymentTime
			cashOut.ErrorCode = response.ErrCode
			cashOut.ErrorCodeDes = response.ErrCodeDes
			models.TransfersSucceed(&cashOut, &account, session)
		}

		account.UnsettledAmount = 0
		_, err = session.Table("author_account").Where("author_id=?", account.AuthorId).Cols("unsettled_amount").Update(&account)
		if err != nil {
			session.Rollback()
			util.Logger.Info("settle" + strconv.Itoa(settlementType) + " err:" + err.Error())
		}

		err = session.Commit()
		if err != nil {
			session.Rollback()
			util.Logger.Info("settle" + strconv.Itoa(settlementType) + " session.Commit() err:" + err.Error())
		}
	}
}