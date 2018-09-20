package main

import (
	"errors"
	"net/http"
	"time"

	"go_web_demo/config"
	"go_web_demo/model"
	v "go_web_demo/pkg/version"
	"go_web_demo/router"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

var (
	// 命令行指定命令参数时的配置: 这里就是命令行启动中跟参数 ./apiserver -c config.yaml 指定配置文件参数名为 config.yaml
	cfg     = pflag.StringP("config", "c", "", "apiserver config file path")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	// 提供版本信息
	pflag.Parse()
	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", " ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		return
	}

	// 初始化配置文件
	// 配置文件中没有的字段，读取时不会报错，只是返回空值，这个要注意下！
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// 初始化数据库
	model.DB.Init()
	// 程序停止时关闭数据库
	defer model.DB.Close()

	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	middlewares := []gin.HandlerFunc{}

	// 加载路由
	router.Load(
		g,
		middlewares...,
	)

	// 启动新线程检测服务启动后是否可用
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	// 开启安全服务
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}

	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// 检查服务器能否连通
func pingServer() error {
	cnt_try := viper.GetInt("max_ping_count")
	url := viper.GetString("url")
	for i := 0; i < cnt_try; i++ {
		resp, err := http.Get(url + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		log.Info("waitting fot the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the server")
}
