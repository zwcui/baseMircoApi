/*
@Time : 2019/2/26 下午2:01 
@Author : zwcui
@Software: GoLand
*/
package models


const (
	ErrorSpliter = "###"
)

const (
	CommonError100 			= "Common100###服务器异常，请联系管理员"

	UserError100			= "User100###您无权限添加用户"
	UserError200			= "User200###该账户昵称已存在"
	UserError300			= "User300###用户不存在"
	UserError400			= "User400###账户密码错误"
	UserError500			= "User400###账户不存在"

	AuthorAccountError100   = "AuthorAccount100###账户不存在"

	AuthorizeError100 		= "Authorize100###请求公众号appid、accessToken出错"
	AuthorizeError200 		= "Authorize200###获取授权方的基本信息出错"
	AuthorizeError300 		= "Authorize300###未找到id对应的公众号"
	AuthorizeError400 		= "Authorize400###无法删除公众号，原因："
	AuthorizeError500 		= "Authorize500###无法恢复公众号，原因："
	AuthorizeError600 		= "Authorize600###未传入授权id"
	AuthorizeError999 		= "Authorize999###"	//自定义

	ActivityError100 		= "Activity100###未找到活动信息"
	ActivityError200 		= "Activity200###创建活动二维码失败，原因："
	ActivityError300 		= "Activity300###当前用户未参加该活动"
	ActivityError400 		= "Activity400###奖品已领完"
	ActivityError500 		= "Activity500###活动已结束"
	ActivityError600 		= "Activity600###邀请好友数不足"
	ActivityError700 		= "Activity700###您已领奖，请等待奖品寄送"
	ActivityError800 		= "Activity800###未找到奖品记录"
	ActivityError900 		= "Activity900###当前活动的关键词已与下面的活动重复："
	ActivityError1000 		= "Activity1000###同一公众号活动启用上限5个"
	ActivityError1100 		= "Activity1100###删除活动前请先停用该活动"

	MessageError100 		= "Message100###查询模板列表失败，原因："
	MessageError200 		= "Message200###未找到id对应的模板消息数据"
	MessageError300 		= "Message300###未找到id对应的客服消息模板"

	HelpError100 			= "Help100###未找到上层信息"
	HelpError200 			= "Help100###未找到id对应的帮助信息"

	ArticleError100 		= "Article100###未找到相关内容"
	ArticleError200 		= "Article200###生成二维码失败"
	ArticleError300 		= "Article300###不支持报名"
	ArticleError400 		= "Article400###请勿重复报名"

	BusinessCardError100 	= "BusinessCard100###未找到对应的名片"

	VideoError100 			= "Video100###未找到对应的视频信息"
	VideoError200 			= "Video200###免费观看时间设置不合法"
	VideoError300 			= "Video300###视频审核中无法修改"
	VideoError400 			= "Video400###视频审核参数不合法"

	CommentError100 		= "Comment100###未找到对应的评论信息"
	CommentError200 		= "Comment200###您只可以删除自己的评论信息"

	AuthorSignInError100 	= "AuthorSignIn100###登录失败："

	OrderError100 			= "Order100###暂不支持其他支付方式"
	OrderError200 			= "Order200###未找到对应的商品信息"
	OrderError300 			= "Order300###该商品已下架"
	OrderError400 			= "Order400###创建订单失败"
	OrderError500 			= "Order500###账户余额不足"
	OrderError600 			= "Order600###您已购买，请勿重复购买"

	TransfersError100       = "Transfers100###您未实名认证，请进行认证"
	TransfersError200       = "Transfers200###您今日付款次数已达上限"
	TransfersError300       = "Transfers300###您今日付款总额已达上限"
	TransfersError400       = "Transfers400###转账申请正在处理中，请稍后"
	TransfersError500       = "Transfers500###转账申请已提交，我们会在12小时工作时间内帮您处理完成，请留意到账信息。"
	TransfersError600       = "Transfers600###您的单笔转账额度不得大于5000元"
	TransfersError700       = "Transfers700###您的单笔转账额度不得小于0.3元"
	TransfersError800       = "Transfers800###该用户不存在，请核对信息"
	TransfersError900       = "Transfers900###您账户余额不足提现金额"
	TransfersError1000      = "Transfers1000###该用户无平台账户"

	VeriCodeError100        = "VeriCodeError100###您输入的手机号有误，请重新输入"
	VeriCodeError200        = "VeriCodeError200###验证码发送失败"
	VeriCodeError300        = "VeriCodeError200###验证码保存失败"

	PutPhoneNumError100     = "PutPhoneNumError100###该手机号已绑定，请输入其他手机号"
	PutPhoneNumError200     = "PutPhoneNumError200###验证码错误"
	PutPhoneNumError300     = "PutPhoneNumError300###验证码已失效"

	SetPassWordError100     = "SetPassWordError100###您俩次的输入的密码不一致"
	SetPassWordError200     = "SetPassWordError200###密码加密失败"
	SetPassWordError300     = "SetPassWordError300###密码不能为空"

	PhoneNumLoginError100   = "PhoneNumLoginError100###您输入的账号名或密码错误，请重新输入"


	CreateQRCodeError100    = "CreateQRCode100###二维码生成失败："
	CreateQRCodeError200    = "CreateQRCode200###二维码读取失败："

	TencentCloudError100    = "TencentCloud100###视频鉴黄失败："
	TencentCloudError200    = "TencentCloud200###获取所有分类失败："
	TencentCloudError300    = "TencentCloud300###查询任务详情失败："
	TencentCloudError400    = "TencentCloud400###获取转码模板列表失败："
	TencentCloudError500    = "TencentCloud500###视频转码失败："

	VideoShareError100 		= "VideoShare100###分享id不能为作者本人"
	VideoShareError200 		= "VideoShare200###该id已参与分享"


	VideoEditError100 		= "VideoEdit100###您无权限修改视频"

	PushMessageError100 	= "PushMessage100###消息推送失败"

	GZHError100				= "GZH100###创建菜单失败"

)
