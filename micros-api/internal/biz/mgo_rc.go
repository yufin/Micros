package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"micros-api/internal/biz/dto"
	"time"
)

type MgoRcRepo interface {
	GetProcessedObjIdByContentId(ctx context.Context, contentId int64) (bson.M, error)
	GetProcessedContentByObjId(ctx context.Context, objIdHex string) (bson.M, error)
	GetProcessedContentInfoByObjId(ctx context.Context, objIdHex string) (bson.M, error)
	GetContentInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosRespV3, error)
	GetContentInfosByKwd(ctx context.Context, page *dto.PaginationReq, kwd string) (*dto.RcOriginContentInfosRespV3, error)
	GetNewestDocInfoByContentId(ctx context.Context, contentId int64) (string, time.Time, error)
	GetNewestDocByContentId(ctx context.Context, contentId int64) (bson.M, error)
	GetRdmResultByClaimedId(ctx context.Context, claimId int64) (bson.M, error)
	GetCountOnDistinctUscIdForColl(ctx context.Context, collName string) (int, error)
	SaveReportPrintConfig(ctx context.Context, b []byte) (*mongo.InsertOneResult, error)
	GetReportPrintConfig(ctx context.Context, cond bson.M) ([]byte, error)
	DeleteReportPrintConfig(ctx context.Context, cond bson.M) error
}

type MgoRcUsecase struct {
	repo MgoRcRepo
	log  *log.Helper
}

func NewMgoRcUsecase(repo MgoRcRepo, logger log.Logger) *MgoRcUsecase {
	return &MgoRcUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *MgoRcUsecase) SaveReportPrintConfig(ctx context.Context, req dto.ReportPrintConfig) error {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = uc.repo.SaveReportPrintConfig(ctx, b)
	if err != nil {
		return err
	}
	return nil
}

func (uc *MgoRcUsecase) GetReportPrintConfig(ctx context.Context, createBy int64) (*dto.ReportPrintConfig, error) {
	cond := bson.M{"create_by": createBy}
	b, err := uc.repo.GetReportPrintConfig(ctx, cond)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}
	var res dto.ReportPrintConfig
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}
	return &res, err
}

func (uc *MgoRcUsecase) UpdateReportPrintConfig(ctx context.Context, config dto.ReportPrintConfig) error {
	cond := bson.M{"create_by": config.CreateBy}
	err := uc.repo.DeleteReportPrintConfig(ctx, cond)
	if err != nil {
		return err
	}
	return uc.SaveReportPrintConfig(ctx, config)
}

func (uc *MgoRcUsecase) GetCountOnDistinctUscIdForColl(ctx context.Context, collName string) (int, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetCountOnDistinctUscIdForColl %s", collName)
	return uc.repo.GetCountOnDistinctUscIdForColl(ctx, collName)
}

func (uc *MgoRcUsecase) GetRdmResultByClaimedId(ctx context.Context, claimId int64) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetRdmResultByClaimedId %d", claimId)
	return uc.repo.GetRdmResultByClaimedId(ctx, claimId)
}

func (uc *MgoRcUsecase) GetNewestDocInfoByContentId(ctx context.Context, contentId int64) (string, time.Time, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetDocInfoByContentId %d", contentId)
	return uc.repo.GetNewestDocInfoByContentId(ctx, contentId)
}

func (uc *MgoRcUsecase) GetProcessedObjIdByContentId(ctx context.Context, contentId int64) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetProcessedObjIdByContentId %d", contentId)
	return uc.repo.GetProcessedObjIdByContentId(ctx, contentId)
}

func (uc *MgoRcUsecase) GetContentInfos(ctx context.Context, page *dto.PaginationReq) (*dto.RcOriginContentInfosRespV3, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetContentInfos %d", page.PageNum)
	return uc.repo.GetContentInfos(ctx, page)
}

func (uc *MgoRcUsecase) GetProcessedContentByObjId(ctx context.Context, objIdHex string) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetProcessedContentByObjId %s", objIdHex)
	return uc.repo.GetProcessedContentByObjId(ctx, objIdHex)
}

func (uc *MgoRcUsecase) GetProcessedContentInfoByObjId(ctx context.Context, objIdHex string) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetProcessedContentInfoByObjId %s", objIdHex)
	return uc.repo.GetProcessedContentInfoByObjId(ctx, objIdHex)
}

func (uc *MgoRcUsecase) GetContentInfosByKwd(ctx context.Context, page *dto.PaginationReq, kwd string) (*dto.RcOriginContentInfosRespV3, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.SearchReportInfosByKwd %d", page.PageNum)
	return uc.repo.GetContentInfosByKwd(ctx, page, kwd)
}

func (uc *MgoRcUsecase) GetNewestDocByContentId(ctx context.Context, contentId int64) (bson.M, error) {
	uc.log.WithContext(ctx).Infof("biz.MgoRcUsecase.GetNewestDocByContentId %d", contentId)
	return uc.repo.GetNewestDocByContentId(ctx, contentId)
}
