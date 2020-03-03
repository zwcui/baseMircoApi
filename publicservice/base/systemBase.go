package base

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/go-xorm/xorm"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/cache"
	//"github.com/garyburd/redigo/redis"
	"strconv"
	"jingting_server/publicservice/models"
	"jingting_server/publicservice/util"
	"github.com/panjf2000/ants"
)

//运行标识
const (
	RUN_MODE_DEV  = "dev"
	RUN_MODE_TEST = "test"
	RUN_MODE_PROD = "prod"
)

//数据库参数、redis参数、服务地址
var redisConn, ServerURL, apiURL, ConsulAddress string
var dbConfig databaseConfig
var rdConfig redisConfig
//协程池
var GoPool *ants.Pool
type databaseConfig struct {
	DbType   		string
	DbUser     		string
	DbPassword 		string
	DbName 			string
	DbCharset  		string
	DbHost     		string
	DbPort     		string
}

type redisConfig struct {
	RedisConn		string
	Auth			string
	Key				string
	DBNum			int
	//MaxIdle			int
	//MaxActive		int
	//IdleTimeout		time.Duration
}

//数据库引擎
var DBEngine *xorm.Engine

//Redis
var RedisCache cache.Cache
//var redisPool *redis.Pool


//系统初始化
func init(){

	if beego.BConfig.RunMode == RUN_MODE_DEV {
		ServerURL = "http://192.168.100.58:8082"
		apiURL = "http://192.168.100.58:8082"
		dbConfig.DbType = "mysql"
		dbConfig.DbHost = "192.168.100.58"
		dbConfig.DbPort = ":3306"
		dbConfig.DbUser = "cnmindstack"
		dbConfig.DbPassword = "cnmindstack"
		dbConfig.DbName = "zhangmai"
		dbConfig.DbCharset = "utf8mb4"
		rdConfig.RedisConn = "192.168.100.58:6379"
		rdConfig.Auth = ""
		rdConfig.DBNum = 0
		ConsulAddress = "127.0.0.1"
	} else if beego.BConfig.RunMode == RUN_MODE_TEST {
		ServerURL = "http://jingting.vipask.net"
		apiURL = "http://jingting.vipask.net"
		dbConfig.DbType = "mysql"
		dbConfig.DbHost = "#ip#"
		dbConfig.DbPort = ":3307"
		dbConfig.DbUser = "cnmindstack"
		dbConfig.DbPassword = "cnmindstack"
		dbConfig.DbName = "zhangmai"
		dbConfig.DbCharset = "utf8mb4"
		rdConfig.RedisConn = "#ip#:6380"
		rdConfig.Auth = ""
		rdConfig.DBNum = 0
		ConsulAddress = "#ip#"
	} else if beego.BConfig.RunMode == RUN_MODE_PROD {
		ServerURL = "http://www.jingtingedu.com"
		apiURL = "http://www.jingtingedu.com"
		dbConfig.DbType = "mysql"
		dbConfig.DbHost = "#ip#"
		dbConfig.DbPort = ":3306"
		dbConfig.DbUser = "cnmindstack"
		dbConfig.DbPassword = "cnmindstack"
		dbConfig.DbName = "jingting"
		dbConfig.DbCharset = "utf8mb4"
		rdConfig.RedisConn = "#ip#:6379"
		rdConfig.Auth = ""
		rdConfig.DBNum = 0
		ConsulAddress = "#ip#"
	} else {
		panic("运行标识错误")
	}

	initDB(dbConfig)
	initRedis(rdConfig)
}



//数据库初始化
func initDB(dbConfig databaseConfig){
	var err error
	//"root:123@/test?charset=utf8"
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=%s",
		dbConfig.DbUser, dbConfig.DbPassword, dbConfig.DbHost, dbConfig.DbPort, dbConfig.DbName, dbConfig.DbCharset)
	fmt.Println(dbUrl)
	DBEngine, err = xorm.NewEngine(dbConfig.DbType, dbUrl)
	if err != nil {
		panic("创建数据库连接Engine失败! err:"+err.Error())
	}
	DBEngine.ShowSQL(false)			//在控制台打印出生成的SQL
	DBEngine.SetMaxIdleConns(10)	//设置闲置的连接数
	DBEngine.SetMaxOpenConns(400)	//设置最大打开的连接数，默认值为0表示不限制
	//cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)	//启用一个全局的内存缓存，存放到内存中，缓存struct的记录数为1000条
	//DBEngine.SetDefaultCacher(cacher)

	err = DBEngine.Ping()
	if err != nil {
		panic("数据库连接ping失败! err:"+err.Error())
	}

	//将sql写入到文件中
	f, err := os.Create("sql.log")
	if err != nil {
		panic("创建sql.log文件失败! err:"+err.Error())
	}
	 defer f.Close()
	DBEngine.SetLogger(xorm.NewSimpleLogger(f))

	//同步表结构
	err = DBEngine.Sync2(new(models.SystemConfig), new(models.SystemAuthCode), new(models.SystemAuthAccessToken), new(models.OperationRecord), new(models.Help), new(models.UserSignInDeviceInfo), new(models.Banner), new(models.Report), new(models.JoinUs), new(models.RecentActivity), new(models.BugRecord), new(models.AppConfig), new(models.Version), new(models.Questionnaire), new(models.QuestionnaireQuestion), new(models.QuestionnaireAnswer), new(models.QuestionnaireAuthorStatistics), new(models.Enroll))
	if err != nil {
		panic("同步表结构失败! err:"+err.Error())
	}
}

//初始化redis
func initRedis(rdConfig redisConfig){
	var err error
	RedisCache, err = cache.NewCache("redis", `{"conn":"`+rdConfig.RedisConn+`", "key":"`+rdConfig.Key+`", "dbNum":"`+strconv.Itoa(rdConfig.DBNum)+`"}`)
	if err != nil {
		panic("redis初始化失败！err:"+err.Error())
	}
	RedisCache.Put("lastStartTime", strconv.FormatInt(util.UnixOfBeijingTime(), 10), 0)
}

