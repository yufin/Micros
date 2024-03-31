package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"micros-dw/internal/biz/dto"
	"time"
)

type DwDataRepo interface {
	GetUscIdByEnterpriseName(ctx context.Context, name string) (*dto.EnterpriseWaitList, error)
	//GetDoc(ctx context.Context, uscId string, tp time.Time, db string, coll string) (*bson.M, error)
	GetDocInDuration(ctx context.Context, uscId string, tp time.Time, coll string) (bson.M, error)
	GetDocInExtendDuration(ctx context.Context, uscId string, tp time.Time, collName string, extendDate int) (bson.M, error)
	GetDocWithDuration(ctx context.Context, uscId string, tp time.Time, coll string, extendDate int) (map[string]any, error)
	GetDocAggregated(ctx context.Context, coll string, filter *mongo.Pipeline) (*[]map[string]any, error)
	GetDoc(ctx context.Context, coll string, cond *bson.D) (*map[string]any, error)
	CountDoc(ctx context.Context, filter bson.D, coll string, db string) (int64, error)
	GetDocs(ctx context.Context, db string, coll string, orderBy string, pageSize int64, pageNum int64) ([]map[string]any, error)
}

type DwDataUsecase struct {
	repo DwDataRepo
	log  *log.Helper
}

func NewDwDataUsecase(repo DwDataRepo, logger log.Logger) *DwDataUsecase {
	return &DwDataUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *DwDataUsecase) GetDocsByPagination(ctx context.Context, db string, coll string, pageSize int64, pageNum int64, sortBy string) ([]map[string]any, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetDocsByPagination %s %d %d %s", coll, pageSize, pageNum, sortBy)
	resArr, err := uc.repo.GetDocs(ctx, db, coll, sortBy, pageSize, pageNum)
	if err != nil {
		return nil, 0, err
	}
	count, err := uc.repo.CountDoc(ctx, bson.D{}, coll, db)
	if err != nil {
		return nil, 0, err
	}
	return resArr, count, nil
}

func (uc *DwDataUsecase) GetUscIdByEnterpriseName(ctx context.Context, name string) (*dto.EnterpriseWaitList, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetEntIdent %s", name)
	return uc.repo.GetUscIdByEnterpriseName(ctx, name)
}

func (uc *DwDataUsecase) GetDocInDuration(ctx context.Context, uscId string, tp time.Time, coll string) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetDocInDuration %s %s %s", uscId, tp, coll)
	return uc.repo.GetDocInDuration(ctx, uscId, tp, coll)
}

func (uc *DwDataUsecase) GetDocInExtendDuration(ctx context.Context, uscId string, tp time.Time, collName string, extendDate int) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetDocInExtendDuration %s %s %s %d", uscId, tp, collName, extendDate)
	return uc.repo.GetDocInExtendDuration(ctx, uscId, tp, collName, extendDate)
}

func (uc *DwDataUsecase) GetDocWithDuration(ctx context.Context, uscId string, tp time.Time, coll string, extendDate int) (map[string]any, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetDocWithDuration %s %s %s %d", uscId, tp, coll, extendDate)
	return uc.repo.GetDocWithDuration(ctx, uscId, tp, coll, extendDate)
}

func (uc *DwDataUsecase) GetDocAggregated(ctx context.Context, coll string, filter *mongo.Pipeline) (*[]map[string]any, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetDocByFilter %s", coll)
	return uc.repo.GetDocAggregated(ctx, coll, filter)
}

func (uc *DwDataUsecase) GetDoc(ctx context.Context, coll string, cond *bson.D) (*map[string]any, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetDocFindOne %s", coll)
	return uc.repo.GetDoc(ctx, coll, cond)
}

func (uc *DwDataUsecase) CountDoc(ctx context.Context, filter bson.D, coll string, db string) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.DwEnterpriseUsecase.GetDocFindOne %s", coll)
	return uc.repo.CountDoc(ctx, filter, coll, db)
}
