package config

import (
	"flag"
	"fmt"
	"os"
)

var (
	ServerPort string
	ApiKey     string
	WxToken    string
)

func init() {
	flag.StringVar(&ServerPort, "PORT", os.Getenv("PORT"), "服务端口号")
	flag.StringVar(&ApiKey, "API_KEY", os.Getenv("API_KEY"), "OpenAI的API_KEY")
	flag.StringVar(&WxToken, "WX_TOKEN", os.Getenv("WX_TOKEN"), "微信公众号令牌")
	flag.Parse()
	if ApiKey == "" {
		fmt.Println("API_KEY 不能为空")
		os.Exit(0)
	}
	if WxToken == "" {
		fmt.Println("WX_TOKEN 未设置，不能用于公众号服务")
	}
	if ServerPort == "" {
		ServerPort = "8080"
	}
}
