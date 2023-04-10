package wechat

import (
	"encoding/xml"
	"time"
)

type Msg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	Content      string   `xml:"Content"`
	Recognition  string   `xml:"Recognition"`

	MsgId int64 `xml:"MsgId,omitempty"`
}

func NewMsg(data []byte) *Msg {
	var msg Msg
	if err := xml.Unmarshal(data, &msg); err != nil {
		return nil
	}
	return &msg
}

func (msg *Msg) GenerateEchoData(s string) []byte {
	data := Msg{
		ToUserName:   msg.FromUserName,
		FromUserName: msg.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      s,
	}
	bs, _ := xml.Marshal(&data)
	return bs
}
