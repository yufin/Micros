package data

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RcRdmResDetailRepo struct {
	data *Data
	log  *log.Helper
}

func NewRcRdmResDetailRepo(data *Data, logger log.Logger) biz.RcRdmResDetailRepo {
	return &RcRdmResDetailRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r RcRdmResDetailRepo) GetByResId(ctx context.Context, resId int64) (*[]dto.RcRdmResDetail, error) {
	var dataRdr dto.RcRdmResDetail
	rdrList := make([]dto.RcRdmResDetail, 0)
	err := r.data.Db.
		Model(&dataRdr).
		Where("res_id = ?", resId).
		Find(&rdrList).
		Error
	if err != nil {
		return nil, err
	}
	return &rdrList, nil
}

func (r RcRdmResDetailRepo) GetByResIdAndLevel(ctx context.Context, resId int64, level int) (*[]dto.RcRdmResDetail, error) {
	var dataRdr dto.RcRdmResDetail
	rdrList := make([]dto.RcRdmResDetail, 0)
	err := r.data.Db.
		Model(&dataRdr).
		Where("res_id = ? and level = ?", resId, level).
		Find(&rdrList).
		Error
	if err != nil {
		return nil, err
	}
	return &rdrList, nil
}
