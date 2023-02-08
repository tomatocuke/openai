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
	r.POST("/", handler.ReceiveMsg)
	// 用于公众号自动验证
	r.GET("/", handler.WechatCheck)
	// 用于测试 curl "http://127.0.0.1:$PORT/test"
	r.GET("/test", handler.Test)
	// 更改模式
	r.GET("/mode", handler.SetMode)

	fmt.Printf("启动服务，使用 curl 'http://127.0.0.1:%s/test?msg=你好哇' 测试一下吧\n", config.ServerPort)

	if err := http.ListenAndServe(":"+config.ServerPort, r); err != nil {
		panic(err)
	}
}
