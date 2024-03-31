package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"micros-api/internal/biz"
)

type UserAuthRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserAuthRepo(data *Data, logger log.Logger) biz.UserAuthRepo {
	return &UserAuthRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *UserAuthRepo) GetUsernameByUserId(ctx context.Context, userId int64) (string, error) {
	var userName string
	err := repo.data.DbBl.
		Table("system_users").
		Select("username").
		Where("id = ?", userId).
		Scan(&userName).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return userName, nil
}
