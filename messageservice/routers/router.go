// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"jingting_server/messageservice/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v2",
		beego.NSNamespace("/message",
			beego.NSInclude(
				&controllers.MessageController{},
			),
		),
		beego.NSNamespace("/push",
			beego.NSInclude(
				&controllers.PushController{},
			),
		),
		beego.NSNamespace("/gzh",
			beego.NSInclude(
				&controllers.GzhController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
