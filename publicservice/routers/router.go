// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"jingting_server/publicservice/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v2",
		beego.NSNamespace("/qrCode",
			beego.NSInclude(
				&controllers.QRCodeController{},
			),
		),
		beego.NSNamespace("/public",
			beego.NSInclude(
				&controllers.PublicController{},
			),
		),
		beego.NSNamespace("/help",
			beego.NSInclude(
				&controllers.HelpController{},
			),
		),
		beego.NSNamespace("/banner",
			beego.NSInclude(
				&controllers.BannerController{},
			),
		),
		beego.NSNamespace("/report",
			beego.NSInclude(
				&controllers.ReportController{},
			),
		),
		beego.NSNamespace("/joinUs",
			beego.NSInclude(
				&controllers.JoinUsController{},
			),
		),
		beego.NSNamespace("/recentActivity",
			beego.NSInclude(
				&controllers.RecentActivityController{},
			),
		),
		beego.NSNamespace("/bug",
			beego.NSInclude(
				&controllers.BugController{},
			),
		),
		beego.NSNamespace("/questionnaire",
			beego.NSInclude(
				&controllers.QuestionnaireController{},
			),
		),
		beego.NSNamespace("/enroll",
			beego.NSInclude(
				&controllers.EnrollController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
