package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:GzhController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:GzhController"],
		beego.ControllerComments{
			Method: "ReceiveMessage",
			Router: `/receiveMessage`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:GzhController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:GzhController"],
		beego.ControllerComments{
			Method: "PostReceiveMessage",
			Router: `/receiveMessage`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:GzhController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:GzhController"],
		beego.ControllerComments{
			Method: "CreateMenu",
			Router: `/createMenu`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:GzhController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:GzhController"],
		beego.ControllerComments{
			Method: "TestRequestGZHUserInfoByOpenId",
			Router: `/testRequestGZHUserInfoByOpenId`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "ReceiveMessage",
			Router: `/receiveMessage/:appId`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "GetTemplateList",
			Router: `/getTemplateList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "AddTemplateMessageData",
			Router: `/addTemplateMessageData`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "GetTemplateMessageDataList",
			Router: `/getTemplateMessageDataList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "DeleteTemplateMessageData",
			Router: `/deleteTemplateMessageData`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "GetSubscriberProvince",
			Router: `/getSubscriberProvince`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "GetSubscriberCity",
			Router: `/getSubscriberCity`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "AddCustomerServiceTemplateMessage",
			Router: `/addCustomerServiceTemplateMessage`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "GetCustomerServiceTemplateMessageList",
			Router: `/getCustomerServiceTemplateMessageList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "DeleteCustomerServiceTemplateMessage",
			Router: `/deleteCustomerServiceTemplateMessage`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "SendCustomerServiceTemplateMessage",
			Router: `/sendCustomerServiceTemplateMessage`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "GetAppMessageList",
			Router: `/getAppMessageList`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:MessageController"],
		beego.ControllerComments{
			Method: "TestCreateShare",
			Router: `/testCreateShare`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["jingting_server/messageservice/controllers:PushController"] = append(beego.GlobalControllerRouter["jingting_server/messageservice/controllers:PushController"],
		beego.ControllerComments{
			Method: "PushToSingle",
			Router: `/push-to-single`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

}
