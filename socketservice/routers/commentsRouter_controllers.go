package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["jingting_server/socketservice/controllers:SocketController"] = append(beego.GlobalControllerRouter["jingting_server/socketservice/controllers:SocketController"],
		beego.ControllerComments{
			Method: "GetAuthorSocketOnlineList",
			Router: `/getAuthorSocketOnlineList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
