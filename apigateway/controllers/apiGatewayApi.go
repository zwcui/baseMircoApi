/*
@Time : 2019/10/14 下午2:07 
@Author : zwcui
@Software: GoLand
*/
package controllers

import (
	"strings"
	"strconv"
	"jingting_server/apigateway/models"
	"jingting_server/apigateway/base"
	"jingting_server/apigateway/util"
)

type ApiGatewayController struct {
	apiController
}

func (this *ApiGatewayController) Prepare(){
	this.NeedUserAuthList = []RequestPathAndMethod{
		{"/refreshUrlConfig", "get", []int{1}},
	}
	this.userAuth()
}

// @Title 刷新网关配置信息
// @Description 刷新网关配置信息
// @Success 200 {object} []models.UrlConfig
// @router /refreshUrlConfig [get]
func (this *ApiGatewayController) RefreshUrlConfig() {

	LoadConfig()

	this.ReturnData = models.UrlConfigList
}


//-------------------------------------------方法-------------------------------------------------------------------

//加载配置
func LoadConfig(){
	var configList []models.ApiConfig
	base.DBEngine.Table("api_config").Where("status=1").Find(&configList)

	if configList == nil {
		util.Logger.Info("未查询到网关配置")
		panic("未查询到网关配置")
		return
	}

	var urlConfigList []models.UrlConfig
	for _, config := range configList {
		var urlConfig models.UrlConfig
		util.Logger.Info("开始加载【" + config.Description + "】配置")
		urlConfig.Description = config.Description

		for _, url := range strings.Split(config.Url, ",") {
			urlConfig.RequestURIArray = append(urlConfig.RequestURIArray, url)
		}

		for _, redirectAndWeight := range strings.Split(config.Weight, ",") {


			redirectUrl := strings.Split(redirectAndWeight, "@@@")[0]
			weightStr := strings.Split(redirectAndWeight, "@@@")[1]
			weight, _ := strconv.Atoi(weightStr)

			urlWeightMap := make(map[string]int)
			urlWeightMap[redirectUrl] = weight
			urlConfig.RequestRedirectArray = append(urlConfig.RequestRedirectArray, initRedirectArray(urlWeightMap)...)
		}
		urlConfigList = append(urlConfigList, urlConfig)
		util.Logger.Info("加载【" + config.Description + "】配置完成")
	}
	models.UrlConfigList = urlConfigList
}

//初始化权重数组
func initRedirectArray(urlWeight map[string]int) (array []string) {
	for url, weight := range urlWeight {
		for i := 0;i < weight;i++ {
			array = append(array, url)
		}
	}

	return array
}