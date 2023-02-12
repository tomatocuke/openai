package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"openai/internal/config"
	"openai/internal/service/fiter"
	"openai/internal/service/openai"
	"openai/internal/service/wechat"
	"strings"
	"sync"
	"time"
)

var (
	success  = []byte("success")
	warn     = "警告，检测到敏感词"
	requests sync.Map // K - 消息ID ， V - chan string
)

func WechatCheck(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	signature := query.Get("signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")
	echostr := query.Get("echostr")

	// 校验
	if wechat.CheckSignature(signature, timestamp, nonce, config.WxToken) {
		w.Write([]byte(echostr))
		return
	}

	log.Println("此接口为公众号验证访问，公众号接入校验失败")
}

// https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Passive_user_reply_message.html
// 微信服务器在五秒内收不到响应会断掉连接，并且重新发起请求，总共重试三次
func ReceiveMsg(w http.ResponseWriter, r *http.Request) {
	bs, _ := io.ReadAll(r.Body)
	msg := wechat.NewMsg(bs)

	// 非文本不回复(返回success表示不回复)
	if !msg.IsText() {
		log.Println("非文本不回复")
		echo(w, success)
		return
	}

	if fiter.Check(msg.Content) != "" {
		warnWx := msg.GenerateEchoData(warn)
		echo(w, warnWx)
		return
	}

	// 5s超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var ch chan string
	v, ok := requests.Load(msg.MsgId)
	if !ok {
		ch = make(chan string)
		requests.Store(msg.MsgId, ch)
		go func(id int64, msg string) {
			// 15s不回复微信，则失效
			timeout := time.Second * 14

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			isFast := !strings.Contains(msg, "代码") && !strings.Contains(msg, "详细")
			result := openai.Query(isFast, msg, timeout)
			ch <- result

			// 定期关闭
			<-ctx.Done()
			close(ch)
			requests.Delete(id)

		}(msg.MsgId, msg.Content)
	} else {
		ch = v.(chan string)
	}

	select {
	case result := <-ch:
		if fiter.Check(result) != "" {
			result = warn
		}
		bs := msg.GenerateEchoData(result)
		echo(w, bs)
	// 超时不要回答，会重试的
	case <-ctx.Done():
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	mode := r.URL.Query().Get("mode")
	isFast := true
	if mode == "full" {
		isFast = false
	}
	if fiter.Check(msg) != "" {
		echo(w, []byte(warn))
		return
	}
	s := openai.Query(isFast, msg, time.Second*30)
	echo(w, []byte(s))
}

func echo(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
