package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"micros-worker/internal/biz/dto"
	"micros-worker/pkg"
)

type CommonNoticeRepo interface {
	SaveConfig(ctx context.Context, b []byte) (*mongo.InsertOneResult, error)
	GetConfigById(ctx context.Context, idHex string) ([]byte, error)
}

type CommonNoticeUsecase struct {
	repo CommonNoticeRepo
	log  *log.Helper
}

// NewCommonNoticeUsecase new a infraNotice usecase.
func NewCommonNoticeUsecase(repo CommonNoticeRepo, logger log.Logger) *CommonNoticeUsecase {
	return &CommonNoticeUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *CommonNoticeUsecase) SaveConfig(ctx context.Context, b []byte) (*mongo.InsertOneResult, error) {
	return uc.repo.SaveConfig(ctx, b)
}

func (uc *CommonNoticeUsecase) SendWechatNoticeInMarkdown(ctx context.Context, req dto.PubNoticeByWechatBotMarkdownReq) error {
	b, err := uc.repo.GetConfigById(ctx, req.SenderId)
	if err != nil {
		return err
	}
	var conf dto.NoticeConfig[dto.WechatBotConfig]
	if err := json.Unmarshal(b, &conf); err != nil {
		return err
	}

	assertNoticeType := dto.WechatBotConfig{}.GetSenderType()
	if conf.SenderType != assertNoticeType {
		return errors.New("notice_type not match")
	}
	wb := pkg.WechatBot{Hook: conf.SenderConfig.WebHook}
	if err := wb.NoticeInMarkdown(req.Content); err != nil {
		return err
	}
	return nil
}

func (uc *CommonNoticeUsecase) SendSms(ctx context.Context, req dto.PubNoticeByWechatBotMarkdownReq) error {
	b, err := uc.repo.GetConfigById(ctx, req.SenderId)
	if err != nil {
		return err
	}
	var conf dto.NoticeConfig[dto.Sms]
	if err := json.Unmarshal(b, &conf); err != nil {
		return err
	}

	assertNoticeType := dto.WechatBotConfig{}.GetSenderType()
	if conf.SenderType != assertNoticeType {
		return errors.New("notice_type not match")
	}
	wb := pkg.WechatBot{Hook: conf.SenderConfig.Uri}
	if err := wb.NoticeInMarkdown(req.Content); err != nil {
		return err
	}
	return nil
}

func (uc *CommonNoticeUsecase) SendWechatNoticeInText(ctx context.Context, req dto.PubNoticeByWechatBotTextReq) error {
	return nil
}
