package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type conf struct {
	Http struct {
		Port  string `json:"port"`
		Proxy string `json:"proxy"`
	} `json:"http"`
	OpenAI struct {
		Key string `json:"key"`
	} `json:"openai"`
	Wechat struct {
		Token string `json:"token"`
	} `json:"wechat"`
}

var (
	C conf
)

func init() {

	// 尝试加载配置文件，否则使用参数
	if err := parseConfigFile(); err != nil {
		fmt.Println("缺少配置文件 config.json")
		os.Exit(0)
	}

	if C.OpenAI.Key == "" {
		fmt.Println("OpenAI的Key不能为空")
		os.Exit(0)
	}

	if C.Http.Port == "" {
		C.Http.Port = "9001"
	}

	if C.Wechat.Token == "" {
		fmt.Println("未设置公众号token，公众号功能不可用")
	}

}

func parseConfigFile() error {
	filename := "./config.json"
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	bs, _ := io.ReadAll(f)
	err = json.Unmarshal(bs, &C)
	return err
}
