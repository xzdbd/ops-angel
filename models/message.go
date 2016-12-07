package models

import (
	"encoding/xml"
	"time"
)

const (
	MsgTypeDefault          = ".*"
	MsgTypeText             = "text"
	MsgTypeImage            = "image"
	MsgTypeVoice            = "voice"
	MsgTypeVideo            = "video"
	MsgTypeLocation         = "location"
	MsgTypeLink             = "link"
	MsgTypeNews             = "news"
	MsgTypeEvent            = "event"
	MsgTypeEventSubscribe   = "subscribe"
	MsgTeypEventUnsubscribe = "unsubscribe"
)

type msgBaseReq struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
}

type msgBaseResp struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	FuncFlag     int // 位0x0001被标志时，星标刚收到的消息
}

type Request struct {
	XMLName xml.Name `xml:xml`
	msgBaseReq
	Content    string
	Location_X float32
	Location_Y float32
	Scale      int
	Label      string
	PicUrl     string
	MsgId      int64
}

// 回复文本消息
type TextResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	Content string
}

// 回复图片消息
type ImageResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	MediaId int64 //通过上传多媒体文件，得到的id。
}

// 回复语音消息
type VoiceResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	MediaId int64 //通过上传多媒体文件，得到的id。
}

// 回复视频消息
type VideoResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	MediaId     int64
	Title       string
	Description string
}

// 回复音乐消息
type MusicResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	Title        string
	Description  string
	MusicURL     string
	HQMusicUrl   string
	ThumbMediaId int64 //缩略图的媒体id，通过上传多媒体文件，得到的id
}

// 回复图文消息
type NewsResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	ArticleCount int     `xml:",omitempty"`
	Articles     []*Item `xml:"Articles>item,omitempty"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}
