package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type WechatBot struct {
	Hook string `json:"hook"`
}

func (w WechatBot) sendMsg(body map[string]any) error {
	cli := resty.New()
	resp, err := cli.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(w.Hook)
	if err != nil {
		return errors.WithStack(err)
	}
	if !resp.IsSuccess() {
		return errors.New(
			fmt.Sprintf(
				"wechatBot send message failed, resp: %v",
				resp.String()),
		)
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return errors.WithStack(err)
	}
	errCode, ok := result["errcode"].(float64)
	if !ok {
		return errors.WithStack(errors.New("errcode parse type assert failed"))
	}
	if errCode != 0 {
		return errors.New(
			fmt.Sprintf(
				"wechatBot send message failed, resp: %v",
				resp.String()),
		)
	}
	return nil
}

func (w WechatBot) NoticeInMarkdown(content string) error {
	body := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]any{
			"content": content,
		},
	}
	return w.sendMsg(body)
}

func (w WechatBot) NoticeInText(content string, mentionedList []string) error {
	body := map[string]any{
		"msgtype": "text",
		"text": map[string]any{
			"content":        content,
			"mentioned_list": mentionedList,
		},
	}
	return w.sendMsg(body)
}
