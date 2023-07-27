package service

import (
	"brillinkmicros/internal/biz"
	"brillinkmicros/internal/biz/dto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/structpb"
	"gorm.io/gorm"
	"strconv"
	"time"

	//"github.com/gogo/protobuf/proto/protojson"
	//"github.com/gogo/protobuf/proto/structpb"
	pb "brillinkmicros/api/rc/v1"
)

type RcServiceServicer struct {
	pb.UnimplementedRcServiceServer
	log                *log.Helper
	rcProcessedContent *biz.RcProcessedContentUsecase
	rcOriginContent    *biz.RcOriginContentUsecase
	rcDependencyData   *biz.RcDependencyDataUsecase
	rcReportOss        *biz.RcReportOssUsecase
	ossMetadata        *biz.OssMetadataUsecase
}

func NewRcServiceServicer(
	rpc *biz.RcProcessedContentUsecase,
	roc *biz.RcOriginContentUsecase,
	rdd *biz.RcDependencyDataUsecase,
	omd *biz.OssMetadataUsecase,
	rro *biz.RcReportOssUsecase,
	logger log.Logger) *RcServiceServicer {
	return &RcServiceServicer{
		rcOriginContent:    roc,
		rcProcessedContent: rpc,
		rcDependencyData:   rdd,
		rcReportOss:        rro,
		ossMetadata:        omd,
		log:                log.NewHelper(logger),
	}
}

// ListReportInfos 获取报告列表
func (s *RcServiceServicer) ListReportInfos(ctx context.Context, req *pb.PaginationReq) (*pb.ReportInfosResp, error) {

	pageReq := &dto.PaginationReq{
		PageNum:  int(req.PageNum),
		PageSize: int(req.PageSize),
	}
	infosResp, err := s.rcOriginContent.GetInfos(ctx, pageReq)
	if err != nil {
		return nil, err
	}

	pbInfos := make([]*pb.ReportInfo, 0)
	for _, v := range *infosResp.Data {
		v := v
		available := false
		if v.ProcessedId != 0 {
			available = true
		}
		info := &pb.ReportInfo{
			ContentId:          v.ContentId,
			EnterpriseName:     v.EnterpriseName,
			UnifiedCreditId:    v.UscId,
			DataCollectMonth:   v.DataCollectMonth,
			Available:          available,
			ContentUpdatedTime: v.ProcessedUpdatedAt.Format("2006-01-02 15:04:05"),
			LhQylx:             int32(v.LhQylx),
			DepId:              v.DepId,
			// TODO: add i18n info
		}
		pbInfos = append(pbInfos, info)

	}

	return &pb.ReportInfosResp{
		PageNum:     uint32(infosResp.PageNum),
		PageSize:    uint32(infosResp.PageSize),
		Total:       uint32(infosResp.Total),
		TotalPage:   uint32(infosResp.TotalPage),
		ReportInfos: pbInfos,
	}, nil
}

// GetReportContent 获取报告内容
func (s *RcServiceServicer) GetReportContent(ctx context.Context, req *pb.ReportContentReq) (*pb.ReportContentResp, error) {
	rpcData, err := s.rcProcessedContent.GetByContentIdUpToDate(ctx, req.ContentId)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if rpcData == nil {
		return &pb.ReportContentResp{
			Content:   nil,
			Available: false,
			Msg:       "没有权限查看该报告/报告未生成",
		}, nil
	}
	if err = json.Unmarshal([]byte(rpcData.Content), &m); err != nil {
		return nil, err
	}
	var st *structpb.Struct
	st, err = structpb.NewStruct(m)
	return &pb.ReportContentResp{
		Content:   st,
		Available: true,
		Msg:       "",
	}, nil
}

// GetReportContentByDepIdNoDs 根据部门id获取报告内容,没有权限验证
func (s *RcServiceServicer) GetReportContentByDepIdNoDs(ctx context.Context, req *pb.ReportContentByDepIdReq) (*pb.ReportContentResp, error) {
	rpcData, err := s.rcProcessedContent.GetContentUpToDateByDepId(ctx, req.DepId, 1)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if rpcData == nil {
		return &pb.ReportContentResp{
			Content:   nil,
			Available: false,
			Msg:       "没有权限查看该报告/报告未生成",
		}, nil
	}
	if err = json.Unmarshal([]byte(rpcData.Content), &m); err != nil {
		return nil, err
	}
	var st *structpb.Struct
	st, err = structpb.NewStruct(m)
	return &pb.ReportContentResp{
		Content:   st,
		Available: true,
		Msg:       "",
	}, nil
}

// RefreshReportContent 刷新报告内容
func (s *RcServiceServicer) RefreshReportContent(ctx context.Context, req *pb.RefreshReportContentReq) (*pb.RefreshReportContentResp, error) {
	if req.ContentId == 0 {
		return &pb.RefreshReportContentResp{
			Success:    false,
			MsgPubTime: time.Now().Format("2006-01-02 15:04:05"),
			Msg:        "contentId为空值",
		}, nil
	}

	success, err := s.rcProcessedContent.RefreshReportContent(ctx, req.ContentId)
	if err != nil {
		return &pb.RefreshReportContentResp{
			Success:    false,
			MsgPubTime: time.Now().Format("2006-01-02 15:04:05"),
			Msg:        "",
		}, err
	}
	return &pb.RefreshReportContentResp{
		Success:    success,
		MsgPubTime: time.Now().Format("2006-01-02 15:04:05"),
		Msg:        "",
	}, nil
}

// InsertReportDependencyData 插入企业风控参数
func (s *RcServiceServicer) InsertReportDependencyData(ctx context.Context, req *pb.InsertDependencyDataReq) (*pb.SetDependencyDataResp, error) {

	isInserted, err := s.rcDependencyData.CheckIsInsertDepdDataDuplicate(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if isInserted {
		return &pb.SetDependencyDataResp{
			Success:     false,
			IsGenerated: false,
			Code:        200,
			Msg:         "该企业风控参数已录入.",
		}, nil
	}

	contentIds, err := s.rcDependencyData.GetDefaultContentIdForInsertDependencyData(ctx, req.UscId)
	if err != nil {
		return nil, err
	}
	if len(contentIds) == 0 {
		insertReq := dto.RcDependencyData{
			UscId:   req.UscId,
			LhQylx:  int(req.LhQylx),
			LhCylwz: int(req.LhCylwz),
			LhGdct:  int(req.LhGdct),
			//LhQybq:       int(req.LhQybq),
			LhYhsx:       int(req.LhYhsx),
			LhSfsx:       int(req.LhSfsx),
			AdditionData: req.AdditionData,
			StatusCode:   0,
		}
		_, err = s.rcDependencyData.Insert(ctx, &insertReq)
		if err != nil {
			return nil, err
		}
		return &pb.SetDependencyDataResp{
			Success:     true,
			IsGenerated: false,
			Code:        200,
			Msg:         "",
		}, nil
	}

	for _, contentId := range contentIds {
		contentId := contentId
		dataRoc, err := s.rcOriginContent.Get(ctx, contentId)
		if err != nil {
			return nil, err
		}
		insertReq := dto.RcDependencyData{
			ContentId:       &contentId,
			AttributedMonth: &dataRoc.YearMonth,
			UscId:           dataRoc.UscId,
			LhQylx:          int(req.LhQylx),
			LhCylwz:         int(req.LhCylwz),
			LhGdct:          int(req.LhGdct),
			//LhQybq:          int(req.LhQybq),
			LhYhsx:       int(req.LhYhsx),
			LhSfsx:       int(req.LhSfsx),
			AdditionData: req.AdditionData,
		}
		_, err = s.rcDependencyData.Insert(ctx, &insertReq)
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

// UpdateReportDependencyData 更新企业风控参数
func (s *RcServiceServicer) UpdateReportDependencyData(ctx context.Context, req *pb.UpdateDependencyDataReq) (*pb.SetDependencyDataResp, error) {
	if req.Id == 0 {
		return nil, errors.BadRequest("Empty ContentId", "row id is required")
	}

	updateReq := dto.RcDependencyData{
		BaseModel: dto.BaseModel{
			Id: req.Id,
		},
		LhQylx:  int(req.LhQylx),
		LhCylwz: int(req.LhCylwz),
		LhGdct:  int(req.LhGdct),
		//LhQybq:       int(req.LhQybq),
		LhYhsx:       int(req.LhYhsx),
		LhSfsx:       int(req.LhSfsx),
		AdditionData: req.AdditionData,
	}
	newRddId, err := s.rcDependencyData.Update(ctx, &updateReq)
	if err != nil {
		return nil, err
	}
	fmt.Println(newRddId)
	//newRdd, err := s.rcDependencyData.Get(ctx, newRddId)
	//if err != nil {
	//	return nil, err
	//}
	return &pb.SetDependencyDataResp{
		Success: true,
		Code:    200,
		Msg:     "success",
	}, nil
}

// GetReportDependencyData 获取企业风控参数
func (s *RcServiceServicer) GetReportDependencyData(ctx context.Context, req *pb.GetDependencyDataReq) (*pb.GetDependencyDataResp, error) {
	// 优先返回本人创建数据,若无则返回dsi.AccessibleIds order by created_at asc
	contentId, err := strconv.ParseInt(req.ContentId, 10, 64)
	if err != nil {
		return nil, err
	}
	dataRdd, err := s.rcDependencyData.GetByContentId(ctx, contentId)
	if err != nil {
		return nil, err
	}
	if dataRdd == nil {
		return nil, err
	}
	return &pb.GetDependencyDataResp{
		Id:           dataRdd.Id,
		ContentId:    *dataRdd.ContentId,
		UscId:        dataRdd.UscId,
		LhQylx:       int32(dataRdd.LhQylx),
		LhCylwz:      int32(dataRdd.LhCylwz),
		LhGdct:       int32(dataRdd.LhGdct),
		LhQybq:       int32(dataRdd.LhQybq),
		LhYhsx:       int32(dataRdd.LhYhsx),
		LhSfsx:       int32(dataRdd.LhSfsx),
		AdditionData: dataRdd.AdditionData,
		CreatedAt:    dataRdd.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    dataRdd.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetReportPdfByDepId 根据depId获取报告pdf
func (s *RcServiceServicer) GetReportPdfByDepId(ctx context.Context, req *pb.ReportDownloadReq) (*pb.OssFileDownloadResp, error) {
	// :TODO 添加oss下载方法
	_, err := s.rcDependencyData.Get(ctx, req.DepId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return &pb.OssFileDownloadResp{
				Available: false,
				Msg:       "没有权限下载该报告",
			}, nil
		}
		return nil, err
	}
	ossId, err := s.rcReportOss.GetOssIdUptoDateByDepId(ctx, req.DepId)
	if err != nil {
		return nil, err
	}
	if ossId == 0 {
		return &pb.OssFileDownloadResp{
			Available: false,
			Msg:       "请等待报告pdf文件生成",
		}, nil
	}

	meta, err := s.ossMetadata.GetById(ctx, ossId)
	if err != nil {
		return nil, err
	}
	var fileName string
	if req.FileName == "" {
		fileName = "风控报告.pdf"
	} else {
		fileName = req.FileName + ".pdf"
	}

	preUrl, err := s.ossMetadata.GetDownloadUrlByObjName(ctx, fileName, meta)
	if err != nil {
		return nil, err
	}

	return &pb.OssFileDownloadResp{
		Available: true,
		Msg:       "",
		Url:       preUrl.String(),
		CreatedAt: meta.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
