package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["jingting_server/apigateway/controllers:ApiGatewayController"] = append(beego.GlobalControllerRouter["jingting_server/apigateway/controllers:ApiGatewayController"],
		beego.ControllerComments{
			Method: "RefreshUrlConfig",
			Router: `/refreshUrlConfig`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
