package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz/dto"
)

type UserAuthRepo interface {
	GetUsernameByUserId(ctx context.Context, userId int64) (string, error)
}

type UserAuthUsecase struct {
	repo UserAuthRepo
	log  *log.Helper
}

func NewUserAuthUsecase(repo UserAuthRepo, logger log.Logger) *UserAuthUsecase {
	return &UserAuthUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *UserAuthUsecase) TagUserNameToContentInfo(ctx context.Context, infos *[]dto.ListReportInfo) error {
	for i, v := range *infos {
		userName, err := uc.repo.GetUsernameByUserId(ctx, v.UserId)
		if err != nil {
			return err
		}
		(*infos)[i].Username = userName
	}
	return nil
}
