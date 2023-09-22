package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
)

func NewRcRdmResultRepo(data *Data, logger log.Logger) biz.RcRdmResultRepo {
	return &RcRdmResultRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

type RcRdmResultRepo struct {
	data *Data
	log  *log.Helper
}

func (r RcRdmResultRepo) GetUpToDate(ctx context.Context, depId int64) (*dto.RcRdmResult, error) {
	var dataRdr dto.RcRdmResult
	err := r.data.Db.
		Model(&dataRdr).
		Where("dep_id = ?", depId).
		Order("created_at desc").
		First(&dataRdr).
		Error
	if err != nil {
		return nil, err
	}
	return &dataRdr, nil
}
