package models


const (
	DEVICE_INFO_SYSTERM_ANDROID = 1
	DEVICE_INFO_SYSTERM_iOS     = 2
)

const (
	DEVICE_INFO_MANUFACTURERS_HUAWEI = 1
	DEVICE_INFO_MANUFACTURERS_MEIZU  = 2
	DEVICE_INFO_MANUFACTURERS_XIAOMI = 3
)

//用户登录信息
type UserSignInDeviceInfo struct {
	AuthorId     int64  `description:"authorId" json:"authorId" xorm:"pk"`
	System       int    `description:"设备系统类型 0 未知 1 android 2 ios 3 h5" json:"system" valid:"Range(0, 3)" json:"system"`
	Manufacturers int    `description:"厂商 0 未知 1 华为 2 魅族 3 小米" json:"manufacturers" `
	DeviceToken   string `description:"deviceToken" json:"deviceToken" xorm:"varchar(200)"`
	DeviceModel   string `description:"设备型号" json:"deviceModel" xorm:"varchar(50)"`
	SystemVersion string `description:"设备系统版本" json:"systemVersion" xorm:"varchar(30)"`
	AppVersion    string `description:"app版本" json:"appVersion" xorm:"varchar(30) "`
	Created       int64  `description:"创建时间" json:"created" xorm:"created"`
	Updated       int64  `description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt     int64  `description:"删除时间" json:"deleted" xorm:"deleted"`
}




type DeviceTokenListContainer struct {
	DeviceToken   string
	Manufacturers int
}










