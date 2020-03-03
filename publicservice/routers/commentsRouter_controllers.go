package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BannerController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BannerController"],
		beego.ControllerComments{
			Method: "AddBanner",
			Router: `/addBanner`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BannerController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BannerController"],
		beego.ControllerComments{
			Method: "UpdateBanner",
			Router: `/updateBanner`,
			AllowHTTPMethods: []string{"patch"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BannerController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BannerController"],
		beego.ControllerComments{
			Method: "GetBannerList",
			Router: `/getBannerList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BannerController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BannerController"],
		beego.ControllerComments{
			Method: "DeleteBanner",
			Router: `/deleteBanner`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BugController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:BugController"],
		beego.ControllerComments{
			Method: "CreateBugRecord",
			Router: `/createBugRecord`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:EnrollController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:EnrollController"],
		beego.ControllerComments{
			Method: "AddEnroll",
			Router: `/addEnroll`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:EnrollController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:EnrollController"],
		beego.ControllerComments{
			Method: "GetEnrollList",
			Router: `/getEnrollList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"],
		beego.ControllerComments{
			Method: "AddHelp",
			Router: `/addHelp`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"],
		beego.ControllerComments{
			Method: "UpdateHelp",
			Router: `/updateHelp`,
			AllowHTTPMethods: []string{"patch"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"],
		beego.ControllerComments{
			Method: "GetHelpList",
			Router: `/getHelpList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"],
		beego.ControllerComments{
			Method: "GetAllHelpList",
			Router: `/getAllHelpList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"],
		beego.ControllerComments{
			Method: "GetHelpInfo",
			Router: `/getHelpInfo`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:HelpController"],
		beego.ControllerComments{
			Method: "DeleteHelp",
			Router: `/deleteHelp`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:JoinUsController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:JoinUsController"],
		beego.ControllerComments{
			Method: "AddJoinUsColumn",
			Router: `/addJoinUsColumn`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:JoinUsController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:JoinUsController"],
		beego.ControllerComments{
			Method: "GetJoinUsColumnList",
			Router: `/getJoinUsColumnList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:JoinUsController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:JoinUsController"],
		beego.ControllerComments{
			Method: "DeleteJoinUsColumn",
			Router: `/deleteJoinUsColumn`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "SignIn",
			Router: `/signIn`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "ResetPassword",
			Router: `/resetPassword`,
			AllowHTTPMethods: []string{"patch"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetSystemComponentAccessToken",
			Router: `/getSystemComponentAccessToken`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetSystemPreAuthCode",
			Router: `/getSystemPreAuthCode`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "PostSignature",
			Router: `/signature`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetUserInfo",
			Router: `/authed-user-info`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "VeriCode",
			Router: `/vericode`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetSystemConfig",
			Router: `/getSystemConfig`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "PostAppConfig",
			Router: `/app-config`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "PatchAppConfig",
			Router: `/app-config/:id`,
			AllowHTTPMethods: []string{"patch"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetAppConfig",
			Router: `/app-config`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/app-config/all`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "PostVersion",
			Router: `/version`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetVersionHistory",
			Router: `/version/history`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "GetVersion",
			Router: `/version/:deviceOs`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "PatchVersion",
			Router: `/version/:versionId`,
			AllowHTTPMethods: []string{"patch"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:PublicController"],
		beego.ControllerComments{
			Method: "DeleteVersion",
			Router: `/version/:versionId`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QRCodeController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QRCodeController"],
		beego.ControllerComments{
			Method: "CreateQRCode",
			Router: `/createQRCode`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QRCodeController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QRCodeController"],
		beego.ControllerComments{
			Method: "CreateWechatQRCode",
			Router: `/createWechatQRCode`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"],
		beego.ControllerComments{
			Method: "AnswerQuestionnaire",
			Router: `/answerQuestionnaire`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"],
		beego.ControllerComments{
			Method: "ShowQuestionnaire",
			Router: `/showQuestionnaire`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"],
		beego.ControllerComments{
			Method: "GetQuestionnaireResult",
			Router: `/getQuestionnaireResult`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"],
		beego.ControllerComments{
			Method: "GetFreeCourseAfterQuestionnaire",
			Router: `/getFreeCourseAfterQuestionnaire`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"],
		beego.ControllerComments{
			Method: "GetQuestionnaireAuthorStatisticsList",
			Router: `/getQuestionnaireAuthorStatisticsList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:QuestionnaireController"],
		beego.ControllerComments{
			Method: "GetQuestionnaireAnswerList",
			Router: `/getQuestionnaireAnswerList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"],
		beego.ControllerComments{
			Method: "AddRecentActivity",
			Router: `/addRecentActivity`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"],
		beego.ControllerComments{
			Method: "UpdateRecentActivity",
			Router: `/updateRecentActivity`,
			AllowHTTPMethods: []string{"patch"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"],
		beego.ControllerComments{
			Method: "GetRecentActivityList",
			Router: `/getRecentActivityList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"],
		beego.ControllerComments{
			Method: "DeleteRecentActivity",
			Router: `/deleteRecentActivity`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:RecentActivityController"],
		beego.ControllerComments{
			Method: "GetRecentActivity",
			Router: `/getRecentActivity`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:ReportController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:ReportController"],
		beego.ControllerComments{
			Method: "AddReport",
			Router: `/addReport`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/publicservice/controllers:ReportController"] = append(beego.GlobalControllerRouter["jingting_server/publicservice/controllers:ReportController"],
		beego.ControllerComments{
			Method: "GetReportList",
			Router: `/getReportList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
