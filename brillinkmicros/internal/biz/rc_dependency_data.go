package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type RcDependencyData struct {
	BaseModel
	ContentId       *int64
	AttributedMonth *string
	UscId           string
	LhQylx          int
	LhCylwz         int
	LhGdct          int
	LhQybq          int
	LhYhsx          int
	LhSfsx          int
	AdditionData    string
}

func (*RcDependencyData) TableName() string {
	return "rc_dependency_data"
}

type RcDependencyDataRepo interface {
	GetByContentId(ctx context.Context, contentId int64) (*RcDependencyData, error)
	Get(ctx context.Context, id int64) (*RcDependencyData, error)
	Insert(ctx context.Context, insertReq *RcDependencyData) (int64, error)
	Update(ctx context.Context, updateReq *RcDependencyData) (int64, error)
	Delete(ctx context.Context, id int64) (bool, error)
	GetDefaultContentIdForInsertDependencyData(ctx context.Context, uscId string) ([]int64, error)
	CheckIsInsertDepdDataDuplicate(ctx context.Context, uscId string) (bool, error)
}

type RcDependencyDataUsecase struct {
	repo RcDependencyDataRepo
	log  *log.Helper
}

func NewRcDependencyDataUsecase(repo RcDependencyDataRepo, logger log.Logger) *RcDependencyDataUsecase {
	return &RcDependencyDataUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RcDependencyDataUsecase) GetByContentId(ctx context.Context, contentId int64) (*RcDependencyData, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.GetByContentId %d", contentId)
	return uc.repo.GetByContentId(ctx, contentId)
}

func (uc *RcDependencyDataUsecase) Get(ctx context.Context, id int64) (*RcDependencyData, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.Get %d", id)
	return uc.repo.Get(ctx, id)
}

func (uc *RcDependencyDataUsecase) Insert(ctx context.Context, insertReq *RcDependencyData) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.Insert %v", insertReq)
	return uc.repo.Insert(ctx, insertReq)
}

func (uc *RcDependencyDataUsecase) Update(ctx context.Context, updateReq *RcDependencyData) (int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.Update %v", updateReq)
	return uc.repo.Update(ctx, updateReq)
}

func (uc *RcDependencyDataUsecase) Delete(ctx context.Context, id int64) (bool, error) {
	uc.log.WithContext(ctx).Infof("biz.RcDependencyDataUsecase.Delete %d", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *RcDependencyDataUsecase) GetDefaultContentIdForInsertDependencyData(ctx context.Context, uscId string) ([]int64, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.GetDefaultContentIdForInsertDependencyData %s", uscId)
	return uc.repo.GetDefaultContentIdForInsertDependencyData(ctx, uscId)
}

func (uc *RcDependencyDataUsecase) CheckIsInsertDepdDataDuplicate(ctx context.Context, uscId string) (bool, error) {
	uc.log.WithContext(ctx).Infof("biz.RcOriginContentUsecase.CheckIsInsertDepdDataDuplicate %s", uscId)
	return uc.repo.CheckIsInsertDepdDataDuplicate(ctx, uscId)
}
