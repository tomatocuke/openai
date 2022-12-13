package handler

import (
	"chatgpt/config"
	"chatgpt/internal/service/gpt"
	"chatgpt/internal/service/wechat"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	cache sync.Map
)

// 请求失败，返回给用户的话术
const TimeoutMsg = "..."

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

	log.Println("微信接入校验失败")
}

// https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Passive_user_reply_message.html
func ReceiveMsg(w http.ResponseWriter, r *http.Request) {
	bs, _ := io.ReadAll(r.Body)
	msg := wechat.NewMsg(bs)

	// 非文本不回复
	if !msg.IsText() {
		log.Println("非文本不回复")
		echo(w, []byte("success"))
		return
	}

	var ch chan string
	v, ok := cache.Load(msg.MsgId)
	if ok {
		ch = v.(chan string)
	} else {
		ch := make(chan string)
		cache.Store(msg.MsgId, ch)
		go func() {
			defer cache.Delete(msg.MsgId)
			//  在第15s前也就是第三个请求结束前，必须获取结果，否则无法返回。
			ch <- gpt.Completions(msg.Content, time.Second*14)
		}()
	}

	// 5s内不返回，会再发起2次重试
	select {
	// ch最晚在第三个请求结束前收到信息
	case s := <-ch:
		close(ch)
		if s == "" {
			s = TimeoutMsg
		}
		bs := msg.GenerateEchoData(s)
		echo(w, bs)
		log.Println("收到消息:", msg.Content, "回复消息:", s)
	// 第一次、第二次失败走这里，腾讯已抛弃接收，发起重试，我们再退出select
	case <-time.After(time.Second * 5):
		return
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	s := gpt.Completions(msg, time.Second*30)
	w.Header().Set("Content-Type", "chatgptlication/text; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
	log.Println("收到消息:", msg, ", 回复消息:", s)
}

func echo(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "chatgptlication/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
