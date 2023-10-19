package v3

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
	"sync"

	//"github.com/gogo/protobuf/proto/protojson"
	//"github.com/gogo/protobuf/proto/structpb"
	pb "micros-api/api/rc/v3"
)

type RcServiceServicer struct {
	pb.UnimplementedRcServiceServer
	log                *log.Helper
	rcProcessedContent *biz.RcProcessedContentUsecase
	rcOriginContent    *biz.RcOriginContentUsecase
	rcDependencyData   *biz.RcDependencyDataUsecase
	rcReportOss        *biz.RcReportOssUsecase
	rcDecisionFactor   *biz.RcDecisionFactorUsecase
	ossMetadata        *biz.OssMetadataUsecase
	mgoRc              *biz.MgoRcUsecase
}

func NewRcServiceServicer(
	rpc *biz.RcProcessedContentUsecase,
	roc *biz.RcOriginContentUsecase,
	rdd *biz.RcDependencyDataUsecase,
	omd *biz.OssMetadataUsecase,
	rro *biz.RcReportOssUsecase,
	rdf *biz.RcDecisionFactorUsecase,
	mgo *biz.MgoRcUsecase,
	logger log.Logger) *RcServiceServicer {
	return &RcServiceServicer{
		rcOriginContent:    roc,
		rcProcessedContent: rpc,
		rcDependencyData:   rdd,
		rcReportOss:        rro,
		ossMetadata:        omd,
		rcDecisionFactor:   rdf,
		mgoRc:              mgo,
		log:                log.NewHelper(logger),
	}
}

func (s *RcServiceServicer) InsertReportDecisionFactor(ctx context.Context, req *pb.InsertDependencyDataReq) (*pb.SetDependencyDataResp, error) {

	countRdf, err := s.rcDecisionFactor.CountByUscIdAndUserId(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if countRdf > 0 {
		return &pb.SetDependencyDataResp{
			Success:     false,
			IsGenerated: false,
			Code:        200,
			Msg:         "该企业风控参数已录入.",
		}, nil
	}

	insertReq := dto.RcDecisionFactor{
		UscId:   req.UscId,
		LhQylx:  int(req.LhQylx),
		LhCylwz: int(req.LhCylwz),
		LhGdct:  int(req.LhGdct),
		LhYhsx:  int(req.LhYhsx),
		LhSfsx:  int(req.LhSfsx),
	}
	rdfId, err := s.rcDecisionFactor.Insert(ctx, &insertReq)
	if err != nil {
		return nil, err
	}
	contentIds, err := s.rcOriginContent.GetContentIdsByUscId(ctx, req.UscId)
	if err != nil {
		return nil, err
	}

	for _, contentId := range contentIds {
		contentId := contentId
		claimReq := dto.RcContentFactorClaim{
			ContentId: contentId,
			FactorId:  rdfId,
		}
		_, err := s.rcDecisionFactor.InsertClaimNoDupe(ctx, &claimReq)
		if err != nil {
			return nil, err
		}
	}

	return &pb.SetDependencyDataResp{
		Success:     true,
		IsGenerated: true,
		Code:        200,
		Msg:         "",
	}, nil
}

func (s *RcServiceServicer) ListReport(ctx context.Context, req *pb.ListReportKwdSearchReq) (*pb.ListReportResp, error) {
	page := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	list, pageInfo, err := s.rcDecisionFactor.ListReportClaimed(ctx, page, req.KwdName)
	if err != nil {
		return nil, err
	}
	infos := make([]*pb.ReportInfo, 0)
	for _, item := range *list {
		infos = append(infos, &pb.ReportInfo{
			UscId:            item.UscId,
			ContentId:        item.ContentId,
			EnterpriseName:   item.EnterpriseName,
			DataCollectMonth: item.DataCollectMonth,
			FactorId:         item.FactorId,
			LhQylx:           int32(item.LhQylx),
		})
	}

	sem := make(chan struct{}, 10)
	errCh := make(chan error, 1)
	done := make(chan bool, 1)
	var wg sync.WaitGroup
	var checkProcessedFunc func(item *pb.ReportInfo)

	switch req.Version {
	case pb.ReportVersion_V3:
		checkProcessedFunc = func(item *pb.ReportInfo) {
			defer wg.Done()
			defer func() { <-sem }()
			processedId, createdAt, err := s.mgoRc.GetNewestDocInfoByContentId(ctx, item.ContentId)
			if err != nil {
				errCh <- err
				return
			}
			if processedId != "" {
				item.Available = true
				item.ContentUpdatedTime = createdAt.Format("2006-01-02 15:04:05")
			} else {
				item.Available = false
			}
		}
	case pb.ReportVersion_V2:
		checkProcessedFunc = func(item *pb.ReportInfo) {
			defer wg.Done()
			defer func() { <-sem }()
			processedId, createdAt, err := s.rcProcessedContent.GetNewestRowInfoByContentId(ctx, item.ContentId)
			if err != nil {
				errCh <- err
				return
			}
			if processedId != 0 {
				item.Available = true
				item.ContentUpdatedTime = createdAt.Format("2006-01-02 15:04:05")
			} else {
				item.Available = false
			}
		}
	default:
		return nil, errors.New(400, "invalid version", string(req.Version))
	}

	for _, item := range infos {
		wg.Add(1)
		sem <- struct{}{}
		go checkProcessedFunc(item)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case err := <-errCh:
		return nil, err
	case <-done:
		break
	}

	return &pb.ListReportResp{
		Success:     true,
		Msg:         "",
		Code:        0,
		Total:       uint32(pageInfo.Total),
		Offset:      uint32(pageInfo.Offset),
		ReportInfos: infos,
	}, nil
}
