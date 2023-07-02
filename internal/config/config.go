package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	Debug bool

	Http struct {
		Port   string
		Proxy  string
		Prefix string
	}

	OpenAI struct {
		Key string

		Params struct {
			Api         string
			Model       string
			Prompt      string
			Temperature float32
			MaxTokens   uint16
		}

		MaxQuestionLength int
	}

	Wechat struct {
		Token        string
		Timeout      int
		SubscribeMsg string
	}
	// User struct {
	// 	QueryTimesDaily int64
	// }
)

func init() {

	// 读取配置
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("解析配置文件config.yaml失败:", err.Error())
		os.Exit(0)
	}

	viper.UnmarshalKey("debug", &Debug)
	viper.UnmarshalKey("http", &Http)
	viper.UnmarshalKey("openai", &OpenAI)
	viper.UnmarshalKey("wechat", &Wechat)
	// viper.UnmarshalKey("user", &User)

	if Http.Prefix == "" {
		Http.Prefix = "/"
	}

	if OpenAI.Key == "" {
		fmt.Println("OpenAI的Key不能为空")
		os.Exit(0)
	}

	if Wechat.Token == "" {
		fmt.Println("未设置公众号token，公众号功能不可用")
	}

	if Wechat.Timeout < 3 || Wechat.Timeout > 13 {
		Wechat.Timeout = 8
	}

	fmt.Println(OpenAI.Params)
}
