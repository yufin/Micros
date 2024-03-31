package dto

import "time"

//type PubNoticeReq struct {
//	SenderType string         `json:"senderType" bson:"sender_type"`
//	NoticeType string         `json:"noticeType" bson:"notice_type"`
//	ConfId     string         `json:"confId" bson:"conf_id"`
//	Content    map[string]any `json:"content"`
//}

type SenderConf interface {
	GetSenderType() string
}

type NoticeConfig[SC SenderConf] struct {
	SenderType   string    `json:"sender_type" bson:"sender_type"`
	SenderConfig SC        `json:"sender_config" bson:"sender_config"`
	Title        string    `json:"title" bson:"title"`
	Comment      string    `json:"comment" bson:"comment"`
	CreateAt     time.Time `json:"create_at" bson:"create_at"`
	UpdateAt     time.Time `json:"update_at" bson:"update_at"`
}

func (n *NoticeConfig[SC]) Gen() {
	n.SenderType = n.SenderConfig.GetSenderType()
	n.CreateAt = time.Now().Local()
}

// WechatBotConfig 微信机器人配置 [SenderConf]
type WechatBotConfig struct {
	WebHook string `json:"web_hook" bson:"web_hook"`
}

func (w WechatBotConfig) GetSenderType() string {
	return "WECHAT_BOT"
}

type Sms struct {
	Uri string
}

func (w Sms) GetSenderType() string {
	return "SMS"
}

type PubNoticeByWechatBotMarkdownReq struct {
	SenderId string
	Content  string
}

type PubNoticeByWechatBotTextReq struct {
	SenderId    string
	Content     string
	MentionList []string
}
