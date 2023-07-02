package main

import (
	"fmt"
	"log"
	"net/http"
	"openai/bootstrap"
	"openai/internal/config"
	"openai/internal/handler"
	"os"
)

func init() {

}

func main() {
	r := bootstrap.New()

	// 微信消息处理
	r.POST("/wx", handler.ReceiveMsg)
	// 用于公众号自动验证
	r.GET("/wx", handler.WechatCheck)
	// 用于测试 curl "http://127.0.0.1:$PORT/test"
	r.GET("/test", handler.Test)
	r.GET("/", handler.Test)

	// 设置日志
	if !config.Debug {
		SetLog()
	}

	fmt.Printf("启动服务，使用 curl 'http://127.0.0.1:%s%stest?msg=你好' 测试一下吧\n", config.Http.Port, config.Http.Prefix)
	if err := http.ListenAndServe(":"+config.Http.Port, r); err != nil {
		panic(err)
	}
}

func SetLog() {
	dir := "./log"
	file := dir + "/data.log"
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	fmt.Println("查看日志请使用 tail -f " + file)
}
