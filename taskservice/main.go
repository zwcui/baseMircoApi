package main

import (
	_ "github.com/astaxie/beego/cache/redis"
	_ "jingting_server/taskservice/task"
	_ "github.com/mkevac/debugcharts"

	"github.com/astaxie/beego"
	"time"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
	"jingting_server/taskservice/base"
	"github.com/micro/go-micro"
	"jingting_server/taskservice/util"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"runtime"
)

func main() {

	beego.BConfig.WebConfig.DirectoryIndex = true

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
		micro.Name("go.micro.taskservice"),
		micro.Version("1.0.0"),
		//健康检查
		micro.RegisterTTL(30 * time.Second),	//注册服务的过期时间
		micro.RegisterInterval(20 * time.Second),//间隔多久再次注册服务
		micro.Registry(reg),
	)

	service.Init()

	// run server
	go func() {
		util.Logger.Info("start to service.Run()")
		if err := service.Run(); err != nil {
			panic(err)
		}
	}()
}

//性能监控，地址：http://localhost:9017/debug/charts/
func performanceMonitoring()  {
	go dummyAllocations()
	go dummyCPUUsage()
	go func() {
		log.Fatal(http.ListenAndServe(":9017", handlers.CompressHandler(http.DefaultServeMux)))
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