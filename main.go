package main

import (
	"chatgpt/bootstrap"
	"chatgpt/config"
	_ "chatgpt/config"
	"chatgpt/internal/handler"
	"fmt"
	"net/http"
)

func main() {
	r := bootstrap.New()

	// 微信消息处理
	r.POST("/chatgpt", handler.ReceiveMsg)
	// 用于公众号自动验证
	r.GET("/chatgpt", handler.WechatCheck)
	// 用于测试 curl "http://127.0.0.1:$PORT/test"
	r.GET("/test", handler.Test)

	fmt.Printf("启动服务，使用 curl 'http://127.0.0.1:%s/test?msg=HelloWorld' 测试一下吧\n", config.ServerPort)

	if err := http.ListenAndServe(":"+config.ServerPort, r); err != nil {
		panic(err)
	}
}
