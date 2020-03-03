/*
@Time : 2019/6/26 上午10:30
@Author : zwcui
@Software: GoLand
*/
package task

import (
	"jingting_server/taskservice/util"
	"jingting_server/taskservice/base"
	"strings"
	"github.com/astaxie/beego"
	"jingting_server/taskservice/models"
)

//每天处理转账失败的提现
func handleFailedCashOut(){
	util.Logger.Info("定时任务，每天处理转账失败的提现")
	var cashOutList []models.CashOut
	base.DBEngine.Table("cash_out").Where("status=2").Find(&cashOutList)
	ip := "#ip#"
	if beego.BConfig.RunMode == "prod" {
		ip = "#ip#"
	}

	for _, cashOut := range cashOutList {
		session := base.DBEngine.NewSession()
		defer session.Close()
		err := session.Begin()
		if err != nil {
			util.Logger.Info("handleFailedCashOut session.Begin() err:" + err.Error())
			session.Close()
			continue
		}

		response, err := models.WeChatTransfers(ip, &cashOut, session)
		if err != nil {
			errCodeDes := strings.Split(err.Error(), ":")[0]
			errCode := strings.Split(err.Error(), ":")[1]
			util.Logger.Info("handleFailedCashOut secondAccount 微信转账失败 err = " + errCodeDes)
			cashOut.ErrorCode = errCode
			cashOut.ErrorCodeDes = errCodeDes

			session.Table("cash_out").Where("id=?", cashOut.Id).Cols("error_code", "error_code_des").Update(&cashOut)
			continue
		}
		cashOut.PaymentNo = response.PaymentNo
		cashOut.TransferTime = response.PaymentTime
		cashOut.ErrorCode = response.ErrCode
		cashOut.ErrorCodeDes = response.ErrCodeDes
		cashOut.Status = 1
		session.Table("cash_out").Where("id=?", cashOut.Id).AllCols().Update(&cashOut)

		var accountTransactionRecord models.AccountTransactionRecord
		hasRecord, _ := session.Table("account_transaction_record").Where("order_id=?", cashOut.Id).And("status=0").And("money_type = 1").And("record_type in (2,3,4) ").Get(&accountTransactionRecord)
		if hasRecord {
			accountTransactionRecord.Status = 1
			session.Table("account_transaction_record").Where("id=?", accountTransactionRecord.Id).Cols("status").Update(&accountTransactionRecord)
		}

		err = session.Commit()
		if err != nil {
			session.Rollback()
			util.Logger.Info("handleFailedCashOut session.Commit() err:" + err.Error())
		}
	}




}
