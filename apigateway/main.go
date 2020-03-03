package main

import (
	_ "github.com/astaxie/beego/cache/redis"
	_ "jingting_server/apigateway/routers"
	"jingting_server/apigateway/util"
	_ "github.com/mkevac/debugcharts"
	"strings"
	"net/http"
	"net/url"
	"fmt"
	"net/http/httputil"
	"log"
	"time"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
	"jingting_server/apigateway/base"
	"github.com/micro/go-micro"
	"math/rand"
	"github.com/gorilla/handlers"
	"runtime"
	"jingting_server/apigateway/models"
	"encoding/json"
	"jingting_server/apigateway/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"strconv"
)

//接口网关，支持权重
func main() {

	//加载参数
	controllers.LoadConfig()

	//协程池
	//base.GoPool, _ = ants.NewPool(1000)

	//服务注册
	registServiceToConsul()

	//性能监控
	//performanceMonitoring()

	//启动接口
	go  startBeego()

	//监听服务
	startServer()

	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	test()
	//	defer wg.Done()
	//}()
	//wg.Wait()
}

type Service struct {
	urlConfigList []models.UrlConfig
}

func (this *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//head校验防止恶意请求
	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(userAgent, "Apache") || strings.Contains(userAgent, "apache") || strings.Contains(userAgent, "Java") || strings.Contains(userAgent, "java") {
		util.Logger.Info("成功拦截恶意请求--header")
		util.Logger.Info("RequestURI "+r.RequestURI)
		util.Logger.Info(r.Header)
		util.Logger.Info("成功拦截恶意请求--header")
		fmt.Fprintf(w, "401 Not authorized")
		return
	}

	//统一身份校验
	//authorization := r.Header.Get("Authorization")

	//ants.Submit(func() {
	//go func() {
	//网关停服提示
	apigatewayOnlineKey := "apigateway_online_status"
	var config models.SystemConfig
	if base.RedisCache.IsExist(apigatewayOnlineKey) {
		apigatewayOnlineStatusRedis := base.RedisCache.Get(apigatewayOnlineKey)
		json.Unmarshal(apigatewayOnlineStatusRedis.([]byte), &config)
	} else {
		hasConfig, _ := base.DBEngine.Table("system_config").Where("program=?", apigatewayOnlineKey).Get(&config)
		if hasConfig {
			apigatewayOnlineStatusRedisBytes, _ := json.Marshal(config)
			base.RedisCache.Put(apigatewayOnlineKey, string(apigatewayOnlineStatusRedisBytes), 24 * 60 * 60 * time.Second)
		}
	}
	if config.ProgramValue != "1" {
		result := "{\"header\": {\"code\": \"1000\",\"description\": \"success\"},\"data\": {\"alertCode\": \"apigateway100\",\"alertMessage\": \""+config.Description+"\"}}"
		w.Write([]byte(result))
		return
	}

	//获得跳转地址
	var remote *url.URL
	remote = this.getRequestURL(r)
	if remote == nil {
		fmt.Fprintf(w, "404 Not Found")
		return
	}

	//转发
	util.Logger.Info("proxy from "+r.URL.String()+" to "+remote.String())
	proxy := httputil.NewSingleHostReverseProxy(remote)

	//ctx := r.Context()
	//timeout := 50 * time.Millisecond
	//ctx, cancel := context.WithTimeout(ctx, timeout)
	//r.WithContext(ctx)
	proxy.ServeHTTP(w, r)
	//defer cancel()
	//}()
}

//根据配置信息获取跳转url
func (this *Service) getRequestURL(r *http.Request) (u *url.URL) {
	for _, urlConfig := range models.UrlConfigList {
		for _, uri := range urlConfig.RequestURIArray {
			if strings.Contains(r.RequestURI, uri) {
				rand := rand.Intn(len(urlConfig.RequestRedirectArray))
				u, _ = url.Parse(urlConfig.RequestRedirectArray[rand])
				return u
			}
		}
	}
	return nil
}



//初始化权重并启动服务
func startServer() {
	// 注册被代理的服务器 (url)
	service := &Service{
		urlConfigList: models.UrlConfigList,
	}
	fmt.Println("start to ListenAndServe")
	err := http.ListenAndServe(":9999", service)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}

//启动接口
func startBeego(){
	beego.BConfig.WebConfig.DirectoryIndex = true

	if beego.BConfig.RunMode != "prod" {
		beego.BConfig.WebConfig.StaticDir["/v2/apiGateway/swagger"] = "swagger"
	}

	// 跨域
	beego.InsertFilter("/v2/*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Api-Version", "Source", "AuthInfoId"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	defer util.Logger.Flush()

	beego.Run()
}

//注册服务至consul
func registServiceToConsul(){
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			base.ConsulAddress,
		}
	})

	// 初始化服务
	service := micro.NewService(
		micro.Name("go.micro.apigateway"),
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


//测试权重
func test(){

	controllers.LoadConfig()

	//var user *handle
	//user = &handle{urlWeight: userUrlWeight, redirectArray: userDefineRequestRedirectArray}

	service := &Service{
		urlConfigList: models.UrlConfigList,
	}

	count10 := 0
	count20 := 0
	count30 := 0
	count40 := 0


	util.Logger.Info("start")

	r := http.Request{}
	r.RequestURI = "/v2/user"

	for i:=0;i<1000;i++ {

		url := service.getRequestURL(&r)
		//util.Logger.Info(url)


		//util.Logger.Info(url.RequestURI())
		if strings.HasPrefix(url.String(), "http://127.0.0.1:10") {
			count10 ++
		} else if strings.HasPrefix(url.String(), "http://127.0.0.1:20") {
			count20 ++
		} else if strings.HasPrefix(url.String(), "http://127.0.0.1:30") {
			count30 ++
		} else if strings.HasPrefix(url.String(), "http://127.0.0.1:40") {
			count40 ++
		}
	}
	util.Logger.Info("end")
	util.Logger.Info("count10="+strconv.Itoa(count10))
	util.Logger.Info("count20="+strconv.Itoa(count20))
	util.Logger.Info("count30="+strconv.Itoa(count30))
	util.Logger.Info("count40="+strconv.Itoa(count40))
}


//性能监控，地址：http://localhost:9030/debug/charts/
func performanceMonitoring()  {
	go dummyAllocations()
	go dummyCPUUsage()
	go func() {
		log.Fatal(http.ListenAndServe(":9030", handlers.CompressHandler(http.DefaultServeMux)))
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