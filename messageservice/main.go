package main

import (
	_ "github.com/astaxie/beego/cache/redis"
	_ "jingting_server/messageservice/routers"
	_ "github.com/mkevac/debugcharts"
	"github.com/astaxie/beego"
	"github.com/panjf2000/ants"
	"github.com/astaxie/beego/plugins/cors"
	"jingting_server/messageservice/base"
	"jingting_server/messageservice/util"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"jingting_server/messageservice/controllers"
	"time"
	"log"
	"net/http"
	"runtime"
	"github.com/gorilla/handlers"
	"jingting_server/messageservice/proto/message"
)

func main() {

	beego.BConfig.WebConfig.DirectoryIndex = true

	if beego.BConfig.RunMode != "prod" {
		beego.BConfig.WebConfig.StaticDir["/v2/message/swagger"] = "swagger"
	}

	// 跨域
	beego.InsertFilter("/v2/*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Api-Version", "Source", "AuthInfoId"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	//协程池
	base.GoPool, _ = ants.NewPool(1000)

	defer util.Logger.Flush()
	defer base.GoPool.Release()

	//服务注册
	registServiceToConsul()

	//性能监控
	//performanceMonitoring()

	beego.Run()
}


func registServiceToConsul(){
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			base.ConsulAddress,
		}
	})

	// 初始化服务
	service := micro.NewService(
		micro.Name("go.micro.messageservice"),
		micro.Version("1.0.0"),
		//健康检查
		micro.RegisterTTL(30 * time.Second),	//注册服务的过期时间
		micro.RegisterInterval(20 * time.Second),//间隔多久再次注册服务
		micro.Registry(reg),
	)

	service.Init()
	// 注册 Handler
	go_micro_messageservice.RegisterPushMessageToUserServiceHandler(service.Server(), new(controllers.PushController))
	go_micro_messageservice.RegisterJumpUrlWithKeyAndParamsServiceHandler(service.Server(), new(controllers.PushController))

	// run server
	go func() {
		if err := service.Run(); err != nil {
			panic(err)
		}
	}()
}

//性能监控，地址：http://localhost:9011/debug/charts/
func performanceMonitoring()  {
	go dummyAllocations()
	go dummyCPUUsage()
	go func() {
		log.Fatal(http.ListenAndServe(":9011", handlers.CompressHandler(http.DefaultServeMux)))
	}()
}

func dummyCPUUsage() {
	var a uint64
	var t = time.Now()
	for {
		t = time.Now()
		a += uint64(t.Unix())
	}
}

func dummyAllocations() {
	var d []uint64

	for {
		for i := 0; i < 2*1024*1024; i++ {
			d = append(d, 42)
		}
		time.Sleep(time.Second * 10)
		d = make([]uint64, 0)
		runtime.GC()
		time.Sleep(time.Second * 10)
	}
}