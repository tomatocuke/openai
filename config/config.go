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
	flag.StringVar(&ServerPort, "port", "9002", "server port")
	flag.StringVar(&ApiKey, "api_key", "", "api key")
	flag.StringVar(&WxToken, "wx_token", "", "wx token")
	flag.Parse()
	if ApiKey == "" {
		fmt.Println(" api key cannot be empty")
		os.Exit(0)
	}
}
