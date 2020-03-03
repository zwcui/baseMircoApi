/*
@Time : 2019/6/5 下午3:26 
@Author : zwcui
@Software: GoLand
*/
package models

//账户信息
type AuthorAccount struct {
	AuthorId				int64				`description:"authorId" json:"authorId" xorm:"pk"`
	Amount					int					`description:"账户金额，单位分" json:"amount"`
	SettlementType			int					`description:"结算类型，0手动结算, 1为每单结算，2为每日结算，3为每周结算，4为每月结算" json:"settlementType" xorm:"notnull default 0"`
	UnsettledAmount			int					`description:"待结算金额" json:"unsettledAmount"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

// 收支明细
type AccountTransactionRecord struct {
	Id       				int64  				`description:"记录id" json:"id" xorm:"pk autoincr"`
	AuthorId		        int64  				`description:"用户id" json:"authorId"`
	Money      				int    				`description:"交易金额，单位分" json:"money"`
	MoneyType  				int    				`description:"交易类型，0:收入(消费),1:支出(提现)" json:"moneyType"`
	RecordType 				int    				`description:"记录类型，0:充值, 1:购买视频, 2:用户提现, 3:平台返现, 4:平台分佣" json:"recordType"`
	RecordName 				string 				`description:"记录类型，0:充值, 1:购买视频, 2:用户提现, 3:平台返现, 4:平台分佣" json:"recordName"`
	OrderId   	 			int64  				`description:"记录类型为0、1、4: 订单编号，payOrder表Id, 记录类型为2、3: 提现单编号 cashOut表Id" json:"orderId"`
	Status     				int    				`description:"交易状态 0:交易未执行 1:交易已执行 2:交易取消" json:"status" xorm:"notnull default 1"`
	ShareVideoId			int64 				`description:"分佣视频id" json:"shareVideoId"`
	ShareAmount				int    				`description:"分佣时购买价，单位分" json:"shareAmount"`
	ShareLevel				int    				`description:"分佣时层级" json:"shareLevel"`
	ShareLevelPercent		int    				`description:"分佣时层级比例" json:"shareLevelPercent"`
	Created           		int64  				`description:"创建时间" json:"created" xorm:"created"`
	Updated           		int64  				`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         		int64  				`description:"删除时间" json:"deleted" xorm:"deleted"`
}

//---------------------------------------------------------------------------------------------

type AccountTransactionRecordContainer struct {
	BaseListContainer
	AccountTransactionRecordList	[]AccountTransactionRecord	`description:"收支明细列表" json:"accountTransactionRecordList"`
}

type AuthorAccountContainer struct{
	AuthorAccount    `description:"账户信息" xorm:"extends"`
	UserInfo		 `description:"用户信息" xorm:"extends"`
}

//佣金明细返回结构体
type AccountCommissionRecord struct {
	OrderId   	 			int64  				`description:"交易类型为0:订单编号，payOrder表Id 交易类型为1: 提现单编号 cashOut表Id" json:"orderId"`
	Money      				int    				`description:"交易金额，单位分" json:"money"`
	ShareVideoId			int64 				`description:"分佣视频id" json:"shareVideoId"`
	ShareAmount				int    				`description:"分佣时购买价，单位分" json:"shareAmount"`
	ShareLevel				int    				`description:"分佣时层级" json:"shareLevel"`
	ShareLevelPercent		int    				`description:"分佣时层级比例" json:"shareLevelPercent"`

}










