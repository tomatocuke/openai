package handler

import (
	"chatgpt/config"
	"chatgpt/internal/service/gpt"
	"chatgpt/internal/service/wechat"
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	success = []byte("success")
	users   sync.Map
)

type user struct {
	ch       chan string
	requests sync.Map
}

func GetUser(id string) *user {
	v, ok := users.Load(id)
	if ok {
		return v.(*user)
	}
	u := &user{
		ch: make(chan string, 5),
	}
	users.Store(id, u)
	return u
}

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

	// 其余正常询问请求
	u := GetUser(msg.FromUserName)

	// 5s超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, ok := u.requests.Load(msg.MsgId)
	if !ok {
		go func() {
			// 最久第3次要给微信回复，设置超时时间为14秒
			u.ch <- gpt.Query(msg.Content, time.Second*14)
			u.requests.Delete(msg.MsgId)
		}()
	}

	select {
	case result := <-u.ch:
		bs := msg.GenerateEchoData(result)
		echo(w, bs)
	// 超时不要回答，会重试的
	case <-ctx.Done():
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	s := gpt.Query(msg, time.Second*10)
	echo(w, []byte(s))
}

func SetMode(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	switch mode {
	case "1":
		gpt.CurrentMode = gpt.FastMode
	case "2":
		gpt.CurrentMode = gpt.NormalMode
	case "3":
		gpt.CurrentMode = gpt.MaxMode
	}
	echo(w, []byte("ok"))
}

func echo(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "chatgptlication/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
