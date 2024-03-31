package v3

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
	pipelineV1 "micros-api/api/pipeline/v1"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
	"micros-api/internal/data"
	"micros-api/pkg"
	"net"
	"net/http"
	"sync"
	"time"

	//"github.com/gogo/protobuf/proto/protojson"
	//"github.com/gogo/protobuf/proto/structpb"
	pb "micros-api/api/rc/v3"
)

type RcServiceServicer struct {
	pb.UnimplementedRcServiceServer
	log                *log.Helper
	rcProcessedContent *biz.RcProcessedContentUsecase
	rcOriginContent    *biz.RcOriginContentUsecase
	rcContentMeta      *biz.RcContentMetaUsecase
	rcDependencyData   *biz.RcDependencyDataUsecase
	rcReportOss        *biz.RcReportOssUsecase
	rcDecisionFactor   *biz.RcDecisionFactorUsecase
	rcDecisionFactorV3 *biz.RcDecisionFactorV3Usecase
	ossMetadata        *biz.OssMetadataUsecase
	mgoRc              *biz.MgoRcUsecase
	pipelineClient     *biz.ClientPipelineUsecase
	userAuth           *biz.UserAuthUsecase
	dataRepo           *data.Data
}

func NewRcServiceServicer(
	rpc *biz.RcProcessedContentUsecase,
	roc *biz.RcOriginContentUsecase,
	rdd *biz.RcDependencyDataUsecase,
	omd *biz.OssMetadataUsecase,
	rro *biz.RcReportOssUsecase,
	rdf *biz.RcDecisionFactorUsecase,
	rdfV3 *biz.RcDecisionFactorV3Usecase,
	mgo *biz.MgoRcUsecase,
	plc *biz.ClientPipelineUsecase,
	rcm *biz.RcContentMetaUsecase,
	uath *biz.UserAuthUsecase,
	data *data.Data,
	logger log.Logger) *RcServiceServicer {
	return &RcServiceServicer{
		rcOriginContent:    roc,
		rcProcessedContent: rpc,
		rcDependencyData:   rdd,
		rcReportOss:        rro,
		ossMetadata:        omd,
		rcDecisionFactor:   rdf,
		rcDecisionFactorV3: rdfV3,
		rcContentMeta:      rcm,
		mgoRc:              mgo,
		pipelineClient:     plc,
		userAuth:           uath,
		dataRepo:           data,
		log:                log.NewHelper(logger),
	}
}

// GetAhpResult 查询原始内容
// keep
func (s *RcServiceServicer) GetAhpResult(ctx context.Context, req *pb.GetAhpResultReq) (*pb.GetAhpResultResp, error) {

	cli := s.pipelineClient.GetClient(ctx)
	resp, err := cli.GetAhpScore(context.TODO(), &pipelineV1.GetAhpScoreReq{
		ClaimId: req.ClaimId,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New(http.StatusFailedDependency, "resp is empty", "record not found/data not accessible")
	}
	if !resp.Success {
		return &pb.GetAhpResultResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     resp.Msg,
			Data:    nil,
		}, nil
	}

	m := resp.Data.AsMap()
	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}

	return &pb.GetAhpResultResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

// GetReportDecisionFactor 查询企业风控参数
func (s *RcServiceServicer) GetReportDecisionFactor(ctx context.Context, req *pb.GetDecisionFactorReq) (*pb.GetDecisionFactorResp, error) {
	factor, err := s.rcDecisionFactor.GetByContentIdWithDataScope(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	// if no record found at RcDecisionFactor, try RcDecisionFactorV3
	if factor == nil {
		factor, err = s.rcDecisionFactorV3.GetByContentIdWithDataScope(ctx, req.ContentId)
		if err != nil {
			return nil, err
		}
	}

	if factor == nil {
		return &pb.GetDecisionFactorResp{
			Success: false,
			Code:    http.StatusNoContent,
			Msg:     "record not found/data not accessible",
			Data:    nil,
		}, nil
	}

	return &pb.GetDecisionFactorResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data: &pb.GetDecisionFactorResp_DecisionFactorData{
			UscId:     factor.UscId,
			LhQylx:    int32(factor.LhQylx),
			LhCylwz:   int32(factor.LhCylwz),
			LhGdct:    int32(factor.LhGdct),
			LhYhsx:    int32(factor.LhYhsx),
			LhSfsx:    int32(factor.LhSfsx),
			CreatedAt: factor.CreatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: factor.UserId,
			ClaimId:   factor.ClaimId,
		},
	}, nil
}

func (s *RcServiceServicer) InsertReportDecisionFactorReplica(ctx context.Context, req *pb.InsertReportDecisionFactorReq) (*pb.InsertReportDecisionFactorResp, error) {
	countRdf, err := s.rcDecisionFactor.CountByUscIdAndUserId(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if countRdf == 0 {
		insertReq := dto.RcDecisionFactor{
			UscId:   req.UscId,
			LhQylx:  int(req.LhQylx),
			LhCylwz: int(req.LhCylwz),
			LhGdct:  int(req.LhGdct),
			LhYhsx:  int(req.LhYhsx),
			LhSfsx:  int(req.LhSfsx),
		}
		_, err := s.rcDecisionFactor.Insert(ctx, &insertReq)
		if err != nil {
			return nil, err
		}
		return &pb.InsertReportDecisionFactorResp{
			Success: true,
			Code:    http.StatusAccepted,
			Msg:     "录入成功",
		}, nil
	}
	return &pb.InsertReportDecisionFactorResp{
		Success: true,
		Code:    0,
		Msg:     "该企业风控参数已录入.",
	}, nil
}

// InsertReportDecisionFactor 录入企业风控参数（认领企业报告）
// keep
func (s *RcServiceServicer) InsertReportDecisionFactor(ctx context.Context, req *pb.InsertReportDecisionFactorReq) (*pb.InsertReportDecisionFactorResp, error) {

	countRdf, err := s.rcDecisionFactor.CountByUscIdAndUserId(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if countRdf > 0 {
		rdfId, err := s.rcDecisionFactorV3.GetLatestIdByUscIdAndUserId(ctx, req.UscId)
		if err != nil {
			return nil, err
		}
		contentIdsV3, err := s.rcContentMeta.GetContentIdsByUscId(ctx, req.UscId)
		if err != nil {
			return nil, err
		}
		for _, contentIdV3 := range contentIdsV3 {
			claimReq := dto.RcContentFactorClaimV3{
				ContentId: contentIdV3,
				FactorId:  rdfId,
			}
			_, err := s.rcDecisionFactorV3.InsertClaimNoDupe(ctx, &claimReq)
			if err != nil {
				return nil, err
			}
		}
		return &pb.InsertReportDecisionFactorResp{
			Success: true,
			Code:    0,
			Msg:     "该企业风控参数已录入.",
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
	// insert for RcContentFactorClaim
	contentIds, err := s.rcOriginContent.GetContentIdsByUscId(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	for _, contentId := range contentIds {
		claimReq := dto.RcContentFactorClaim{
			ContentId: contentId,
			FactorId:  rdfId,
		}
		_, err := s.rcDecisionFactor.InsertClaimNoDupe(ctx, &claimReq)
		if err != nil {
			return nil, err
		}
	}

	// insert for RcContentFactorClaimV3 as well.
	contentIdsV3, err := s.rcContentMeta.GetContentIdsByUscId(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	for _, contentIdV3 := range contentIdsV3 {
		contentIdV3 := contentIdV3
		claimReq := dto.RcContentFactorClaimV3{
			ContentId: contentIdV3,
			FactorId:  rdfId,
		}
		_, err := s.rcDecisionFactorV3.InsertClaimNoDupe(ctx, &claimReq)
		if err != nil {
			return nil, err
		}
	}

	return &pb.InsertReportDecisionFactorResp{
		Success: true,
		Code:    http.StatusAccepted,
		Msg:     "",
	}, nil
}

// UpdateReportDecisionFactor 更新企业风控参数
// Todo: reimplement me
func (s *RcServiceServicer) UpdateReportDecisionFactor(ctx context.Context, req *pb.UpdateReportDecisionFactorReq) (*pb.InsertReportDecisionFactorResp, error) {
	// logic: insert to rdf, got newFactorId
	// get contentId by claimId, insert new row to claim table with contentId and newFactorId
	claimed, err := s.rcDecisionFactor.GetClaimRecord(ctx, req.ClaimId)
	if err != nil {

		// distinguish where the claimId comes from, rcDecisionFactorV3 or rcDecisionFactor
		// insert to v3 as well
		if errors.Is(err, gorm.ErrRecordNotFound) {
			claimedV3, err := s.rcDecisionFactorV3.GetClaimRecord(ctx, req.ClaimId)
			if err != nil {
				return nil, err
			}
			contentV3, err := s.rcContentMeta.Get(ctx, claimedV3.ContentId)
			if err != nil {
				return nil, err
			}
			if contentV3 == nil {
				return &pb.InsertReportDecisionFactorResp{
					Success: false,
					Code:    http.StatusNotFound,
					Msg:     err.Error(),
				}, nil
			}
			insertReqV3 := dto.RcDecisionFactor{
				UscId:   contentV3.UscId,
				LhQylx:  int(req.LhQylx),
				LhCylwz: int(req.LhCylwz),
				LhGdct:  int(req.LhGdct),
				LhYhsx:  int(req.LhYhsx),
			}
			newFactorIdV3, err := s.rcDecisionFactorV3.Insert(ctx, &insertReqV3)
			if err != nil {
				return nil, err
			}
			_, err = s.rcDecisionFactorV3.InsertClaimNoDupe(ctx, &dto.RcContentFactorClaimV3{
				ContentId: claimedV3.ContentId,
				FactorId:  newFactorIdV3,
			})
			return &pb.InsertReportDecisionFactorResp{
				Success: true,
				Code:    http.StatusAccepted,
				Msg:     "",
			}, nil
		}

		return nil, err
	}
	content, err := s.rcOriginContent.Get(ctx, claimed.ContentId)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return &pb.InsertReportDecisionFactorResp{
			Success: false,
			Code:    http.StatusNotFound,
			Msg:     "content not found",
		}, nil
	}

	insertReq := dto.RcDecisionFactor{
		UscId:   content.UscId,
		LhQylx:  int(req.LhQylx),
		LhCylwz: int(req.LhCylwz),
		LhGdct:  int(req.LhGdct),
		LhYhsx:  int(req.LhYhsx),
		LhSfsx:  int(req.LhSfsx),
	}
	newFactorId, err := s.rcDecisionFactor.Insert(ctx, &insertReq)
	if err != nil {
		return nil, err
	}

	_, err = s.rcDecisionFactor.InsertClaimNoDupe(ctx, &dto.RcContentFactorClaim{
		ContentId: claimed.ContentId,
		FactorId:  newFactorId,
	})
	if err != nil {
		return nil, err
	}
	// insert for v3 as well
	return &pb.InsertReportDecisionFactorResp{
		Success: true,
		Code:    http.StatusAccepted,
		Msg:     "",
	}, nil
}

func (s *RcServiceServicer) ListCompaniesWaiting(ctx context.Context, req *pb.ListCompanyWaitingReq) (*pb.ListCompaniesWaitingResp, error) {
	page := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	var list *[]dto.ListCompaniesWaitingResp
	var total int64
	var err error

	switch req.Version {
	case pb.ReportVersion_V3:
		list, total, err = s.rcDecisionFactorV3.ListCompaniesWaiting(ctx, page, req.UscidKwd, "V2.0")
	case pb.ReportVersion_V2_5:
		list, total, err = s.rcDecisionFactorV3.ListCompaniesWaiting(ctx, page, req.UscidKwd, "V1.0")
	default:
		return nil, errors.New(400, "unsupport report version", string(req.Version))
	}

	if err != nil {
		return nil, err
	}

	infos := make([]*pb.ListCompaniesWaitingResp_CompaniesWaiting, 0)
	for _, item := range *list {
		infos = append(infos, &pb.ListCompaniesWaitingResp_CompaniesWaiting{
			UscId:     item.UscId,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			Id:        pkg.RandUuid(),
		})
	}
	return &pb.ListCompaniesWaitingResp{
		Success: true,
		Msg:     "",
		Code:    http.StatusOK,
		Total:   uint32(total),
		Data:    infos,
	}, nil
}

// ListReport 获取报告列表
// drop
func (s *RcServiceServicer) ListReport(ctx context.Context, req *pb.ListReportKwdSearchReq) (*pb.ListReportResp, error) {
	page := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	var list *[]dto.ListReportInfo
	var pageInfo dto.PaginationInfo
	var err error

	if req.Version == pb.ReportVersion_V3 {
		list, pageInfo, err = s.rcDecisionFactorV3.ListReportClaimed(ctx, page, req.NameKwd)
	} else {
		list, pageInfo, err = s.rcDecisionFactor.ListReportClaimed(ctx, page, req.NameKwd)
	}
	if err != nil {
		return nil, err
	}

	infos := make([]*pb.ListReportResp_ReportInfo, 0)
	for _, item := range *list {
		infos = append(infos, &pb.ListReportResp_ReportInfo{
			UscId:              item.UscId,
			ContentId:          item.ContentId,
			EnterpriseName:     item.EnterpriseName,
			DataCollectMonth:   item.DataCollectMonth,
			ContentUpdatedTime: item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	sem := make(chan struct{}, 10)
	errCh := make(chan error, 1)
	done := make(chan bool, 1)
	var wg sync.WaitGroup
	var checkProcessedFunc func(item *pb.ListReportResp_ReportInfo)

	switch req.Version {
	case pb.ReportVersion_V3:
		checkProcessedFunc = func(item *pb.ListReportResp_ReportInfo) {
			defer wg.Done()
			defer func() { <-sem }()
			item.Available = true
		}
	case pb.ReportVersion_V2_5:
		checkProcessedFunc = func(item *pb.ListReportResp_ReportInfo) {
			defer wg.Done()
			defer func() { <-sem }()
			item.Available = true
		}
	case pb.ReportVersion_V2:
		checkProcessedFunc = func(item *pb.ListReportResp_ReportInfo) {
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
				item.ContentUpdatedTime = ""
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
		Success: true,
		Msg:     "",
		Code:    0,
		Total:   uint32(pageInfo.Total),
		Offset:  uint32(pageInfo.Offset),
		Data:    infos,
	}, nil
}

// ListReportByUscId 获取报告列表
func (s *RcServiceServicer) ListReportByUscId(ctx context.Context, req *pb.ListReportByUscIdReq) (*pb.ListReportByUscIdResp, error) {

	page := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	var list *[]dto.ListReportInfo
	var total int64
	var err error

	switch req.Version {
	case pb.ReportVersion_V3:
		if err := s.rcDecisionFactorV3.SyncClaimed(ctx, req.UscId, "V2.0"); err != nil {
			return nil, err
		}
		list, total, err = s.rcDecisionFactorV3.ListReportClaimedByUscId(ctx, page, req.UscId, "V2.0")
		if err != nil {
			return nil, err
		}
		err = s.userAuth.TagUserNameToContentInfo(ctx, list)
		if err != nil {
			return nil, err
		}

	case pb.ReportVersion_V2_5:
		if err := s.rcDecisionFactorV3.SyncClaimed(ctx, req.UscId, "V1.0"); err != nil {
			return nil, err
		}
		list, total, err = s.rcDecisionFactorV3.ListReportClaimedByUscId(ctx, page, req.UscId, "V1.0")
	default:
		return nil, errors.New(400, "unsupport version", string(req.Version))
	}
	if err != nil {
		return nil, err
	}

	infos := make([]*pb.ListReportByUscIdResp_ReportInfo, 0)
	for _, item := range *list {
		infos = append(infos, &pb.ListReportByUscIdResp_ReportInfo{
			ContentId:          item.ContentId,
			DataCollectMonth:   item.DataCollectMonth,
			ContentUpdatedTime: item.CreatedAt.Format("2006-01-02 15:04:05"),
			Available:          true,
			Id:                 pkg.RandUuid(),
			Status:             item.Status,
			CreateBy:           item.Username,
		})
	}
	return &pb.ListReportByUscIdResp{
		Success: true,
		Msg:     "",
		Code:    0,
		Total:   uint32(total),
		Data:    infos,
	}, nil
}

// ListCompanies 获取报告列表
// keep
func (s *RcServiceServicer) ListCompanies(ctx context.Context, req *pb.ListReportKwdSearchReq) (*pb.ListCompaniesResp, error) {
	page := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	var list *[]dto.ListCompaniesLatest
	var total int64
	var err error

	switch req.Version {
	case pb.ReportVersion_V3:
		list, total, err = s.rcDecisionFactorV3.ListCompanies(ctx, page, req.NameKwd, "V2.0")
	case pb.ReportVersion_V2_5:
		list, total, err = s.rcDecisionFactorV3.ListCompanies(ctx, page, req.NameKwd, "V1.0")
	default:
		return nil, errors.New(400, "unsupport report version", string(req.Version))
	}

	if err != nil {
		return nil, err
	}

	infos := make([]*pb.ListCompaniesResp_CompanyInfo, 0)
	for _, item := range *list {
		infos = append(infos, &pb.ListCompaniesResp_CompanyInfo{
			UscId:          item.UscId,
			EnterpriseName: item.EnterpriseName,
			LastUpdate:     item.LastCreatedAt.Format("2006-01-02 15:04:05"),
			Id:             pkg.RandUuid(),
		})
	}
	return &pb.ListCompaniesResp{
		Success: true,
		Msg:     "",
		Code:    0,
		Total:   uint32(total),
		Data:    infos,
	}, nil
}

// GetReportContent 获取企业报告内容
// keep
func (s *RcServiceServicer) GetReportContent(ctx context.Context, req *pb.GetReportContentReq) (*pb.GetReportContentResp, error) { // test dynamic client
	var accessible bool
	var err error
	if req.Version == pb.ReportVersion_V3 || req.Version == pb.ReportVersion_V2_5 {
		accessible, err = s.rcDecisionFactorV3.CheckContentIdAccessible(ctx, req.ContentId)
	} else {
		accessible, err = s.rcDecisionFactor.CheckContentIdAccessible(ctx, req.ContentId)
	}
	if err != nil {
		return nil, err
	}

	if !accessible {
		return &pb.GetReportContentResp{
			Success: false,
			Code:    http.StatusForbidden,
			Msg:     "data not accessible",
			Data:    nil,
		}, nil
	}

	m := make(map[string]interface{})

	// get processed content
	switch req.Version {
	case pb.ReportVersion_V2:
		data, err := s.rcProcessedContent.GetNewestByContentId(ctx, req.ContentId)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return &pb.GetReportContentResp{
				Success: false,
				Code:    http.StatusNoContent,
				Msg:     "record not found",
				Data:    nil,
			}, nil
		}
		err = json.Unmarshal([]byte(data.Content), &m)
		if err != nil {
			return nil, err
		}
	case pb.ReportVersion_V2_5:
		if req.Realtime {
			cli := s.pipelineClient.GetClient(ctx)
			resp, err := cli.GetContentProcess(context.TODO(), &pipelineV1.GetContentProcessReq{ContentId: req.ContentId, ReportVersion: pipelineV1.ReportVersion_V2_5})
			if err != nil {
				return nil, err
			}
			if !resp.Success {
				return &pb.GetReportContentResp{
					Success: false,
					Code:    http.StatusFailedDependency,
					Msg:     resp.Msg,
					Data:    nil,
				}, nil
			}
			m = resp.Data.AsMap()
		} else {
			data, err := s.mgoRc.GetNewestDocByContentId(ctx, req.ContentId)
			if err != nil {
				return nil, err
			}
			if data == nil {
				return &pb.GetReportContentResp{
					Success: false,
					Code:    http.StatusNoContent,
					Msg:     "record not found",
					Data:    nil,
				}, nil
			}
			b, err := bson.MarshalExtJSON(data["content"], false, false)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(b, &m)
			if err != nil {
				return nil, err
			}
		}
	case pb.ReportVersion_V3:
		httpReq, ok := kratosHttp.RequestFromServerContext(ctx)
		if !ok {
			return nil, errors.New(400, "invalid request context", "")
		}
		lang := httpReq.Header.Get("Accept-Language")
		if req.Realtime {
			cli := s.pipelineClient.GetClient(ctx)
			resp, err := cli.GetContentProcess(context.TODO(), &pipelineV1.GetContentProcessReq{ContentId: req.ContentId, ReportVersion: pipelineV1.ReportVersion_V3, Lang: lang})
			if err != nil {
				return nil, err
			}
			if !resp.Success {
				return &pb.GetReportContentResp{
					Success: false,
					Code:    http.StatusFailedDependency,
					Msg:     resp.Msg,
					Data:    nil,
				}, nil
			}
			m = resp.Data.AsMap()
		}
	default:
		return nil, errors.New(400, "invalid version", string(req.Version))
	}

	st, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}

	return &pb.GetReportContentResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    st,
	}, nil
}

func (s *RcServiceServicer) GetTradeDetail(ctx context.Context, req *pb.GetTradeDetailReq) (*pb.GetTradeDetailResp, error) {
	var accessible bool
	var err error
	if req.ReportVersion == pb.ReportVersion_V3 || req.ReportVersion == pb.ReportVersion_V2_5 {
		accessible, err = s.rcDecisionFactorV3.CheckContentIdAccessible(ctx, req.ContentId)
	} else {
		accessible, err = s.rcDecisionFactor.CheckContentIdAccessible(ctx, req.ContentId)
	}
	if err != nil {
		return nil, err
	}

	if !accessible {
		return &pb.GetTradeDetailResp{
			Success: false,
			Code:    http.StatusForbidden,
			Msg:     "data not accessible",
			Data:    nil,
		}, nil
	}

	cli := s.pipelineClient.GetClient(ctx)

	optionTimePeriod := make([]pipelineV1.GetTradeDetailReq_TimePeriodOption, 0)
	for _, period := range req.OptionTimePeriod {
		optionTimePeriod = append(optionTimePeriod, pipelineV1.GetTradeDetailReq_TimePeriodOption(period))
	}

	optionTradeFrequency := make([]pipelineV1.GetTradeDetailReq_TradeFrequencyOption, 0)
	for _, frequency := range req.OptionTradeFrequency {
		optionTradeFrequency = append(optionTradeFrequency, pipelineV1.GetTradeDetailReq_TradeFrequencyOption(frequency))
	}

	resp, err := cli.GetTradeDetail(context.TODO(), &pipelineV1.GetTradeDetailReq{
		ContentId:            req.ContentId,
		ReportVersion:        pipelineV1.ReportVersion_V3,
		OptionTimePeriod:     optionTimePeriod,
		OptionTopCus:         pipelineV1.GetTradeDetailReq_TopCusOption(req.OptionTopCus),
		OptionTradeFrequency: optionTradeFrequency,
		TradeType:            pipelineV1.GetTradeDetailReq_TradeType(req.TradeType),
	})

	m := resp.Data.AsMap()
	httpReq, ok := kratosHttp.RequestFromServerContext(ctx)
	if !ok {
		return nil, errors.New(400, "invalid request context", "")
	}

	lang := httpReq.Header.Get("Accept-Language")
	var st *structpb.Struct
	if lang == "en-US" {
		hash, err := pkg.GenMapCacheKey(m)
		if err != nil {
			return nil, err
		}
		k := "tradeDetailEn:" + hash
		bEn, err := s.dataRepo.Rdb.Client.Get(context.TODO(), k).Bytes()
		if err != nil {
			var netErr net.Error
			if errors.Is(err, redis.Nil) || (errors.As(err, &netErr) && netErr.Timeout()) {
				s.log.Infof("refresh tradeDetailEn key: %v", k)
				b, err := json.Marshal(resp.Data.AsMap())
				var respM map[string]any
				err = json.Unmarshal(b, &respM)
				if err != nil {
					return nil, err
				}
				respSt, err := structpb.NewStruct(respM)
				if err != nil {
					return nil, err
				}

				transResp, err := cli.GetJsonTranslate(context.TODO(), &pipelineV1.GetJsonTranslateReq{Data: respSt})
				if err != nil {
					return nil, err
				}
				if !transResp.Success {
					return nil, errors.New(400, "GetJsonTranslate Failed", "")
				}
				st, err = structpb.NewStruct(transResp.Data.AsMap())
				if err != nil {
					return nil, err
				}
				bEn, err = json.Marshal(transResp.Data.AsMap())
				if err != nil {
					return nil, err
				}
				//err := s.dataRepo.Rdb.Set(context.TODO(), k, bTrans, time.Hour*24*5).Err()
				err = s.dataRepo.Rdb.Client.Set(context.TODO(), k, bEn, time.Hour*24*3).Err()
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err

			}
		}
		var transCachedM map[string]any
		err = json.Unmarshal(bEn, &transCachedM)
		if err != nil {
			return nil, err
		}
		st, err = structpb.NewStruct(transCachedM)
		if err != nil {
			return nil, err
		}

	} else {
		st, err = structpb.NewStruct(m)
		if err != nil {
			return nil, err
		}
	}

	return &pb.GetTradeDetailResp{
		Success: resp.Success,
		Code:    resp.Code,
		Msg:     resp.Msg,
		Data:    st,
	}, nil
}

// GetReportDataValidationStats 获取报告数据验证统计数据
// keep
func (s *RcServiceServicer) GetReportDataValidationStats(ctx context.Context, req *pb.GetReportDataValidationStatsReq) (*pb.GetReportDataValidationStatsResp, error) { // test dynamic client
	accessible, err := s.rcDecisionFactor.CheckContentIdAccessible(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	if !accessible {
		return &pb.GetReportDataValidationStatsResp{
			Success: false,
			Code:    http.StatusForbidden,
			Msg:     "data not accessible",
			Data:    nil,
		}, nil
	}
	var version pipelineV1.ReportVersion
	switch req.Version {
	case pb.ReportVersion_V2:
		version = pipelineV1.ReportVersion_V2
	case pb.ReportVersion_V2_5:
		version = pipelineV1.ReportVersion_V2_5
	case pb.ReportVersion_V3:
		version = pipelineV1.ReportVersion_V3
	case pb.ReportVersion_Latest:
		version = pipelineV1.ReportVersion_LATEST
	default:
		return nil, errors.New(400, "invalid version", string(req.Version))
	}
	cli := s.pipelineClient.GetClient(ctx)
	resp, err := cli.GetContentValidate(context.TODO(), &pipelineV1.GetContentProcessReq{ContentId: req.ContentId, ReportVersion: version})
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return &pb.GetReportDataValidationStatsResp{
			Success: false,
			Code:    http.StatusFailedDependency,
			Msg:     resp.Msg,
		}, nil
	}

	return &pb.GetReportDataValidationStatsResp{
		Success: true,
		Code:    http.StatusOK,
		Msg:     "",
		Data:    resp.Data,
	}, nil
}

func (s *RcServiceServicer) GetReportPrintConfig(ctx context.Context, _ *emptypb.Empty) (*pb.GetReportPrintConfigResp, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}

	res, err := s.mgoRc.GetReportPrintConfig(ctx, dsi.UserId)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return &pb.GetReportPrintConfigResp{
			Success: true,
			Code:    0,
			Msg:     "",
			Config:  nil,
		}, nil
	}
	st, err := structpb.NewStruct(res.Config)
	if err != nil {
		return nil, err
	}

	return &pb.GetReportPrintConfigResp{
		Success: true,
		Code:    0,
		Msg:     "",
		Config:  st,
	}, nil
}

func (s *RcServiceServicer) UpdateReportPrintConfig(ctx context.Context, req *pb.SaveReportPrintConfigReq) (*pb.SaveReportPrintConfigResp, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}

	config := dto.ReportPrintConfig{
		CreateBy: dsi.UserId,
		CreateAt: time.Now().Local(),
		Config:   req.Config.AsMap(),
	}

	err = s.mgoRc.UpdateReportPrintConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &pb.SaveReportPrintConfigResp{
		Success: true,
		Code:    0,
		Msg:     "",
	}, nil
}
