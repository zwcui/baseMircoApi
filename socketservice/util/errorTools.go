/*
@Time : 2019/2/26 下午2:11 
@Author : zwcui
@Software: GoLand
*/
package util

import (
	"strings"
	"jingting_server/socketservice/models"
)

//根据错误码获取
func getErrorCodeAndDescription(error string, suffix string) (errorCode string, errorDescription string) {
	errorCode = strings.Split(error, models.ErrorSpliter)[0]
	errorDescription = strings.Split(error, models.ErrorSpliter)[1] + suffix
	return errorCode, errorDescription
}

//生成错误返回结构体
func GenerateAlertMessage(errorArgs ...string) models.AlertMessage{
	var alert models.AlertMessage
	if len(errorArgs) > 1 {
		alert.AlertCode, alert.AlertMessage = getErrorCodeAndDescription(errorArgs[0], errorArgs[1])
	} else if len(errorArgs) == 1 {
		alert.AlertCode, alert.AlertMessage = getErrorCodeAndDescription(errorArgs[0], "")
	} else {
		alert.AlertCode, alert.AlertMessage = getErrorCodeAndDescription(models.CommonError100, "")
	}
	return alert
}
