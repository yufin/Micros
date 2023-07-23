package service

import (
	"brillinkmicros/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	pb "brillinkmicros/api/rc/v1"
)

type RcRdmServiceServicer struct {
	pb.UnimplementedRcRdmServiceServer
	log              *log.Helper
	rcRdmResult      *biz.RcRdmResultUsecase
	rcRdmResDetail   *biz.RcRdmResDetailUsecase
	rcDependencyData *biz.RcDependencyDataUsecase
}

func NewRcRdmServiceServicer(
	rrr *biz.RcRdmResultUsecase,
	rrrd *biz.RcRdmResDetailUsecase,
	rdd *biz.RcDependencyDataUsecase,
	logger log.Logger,
) *RcRdmServiceServicer {
	return &RcRdmServiceServicer{
		rcRdmResDetail:   rrrd,
		rcRdmResult:      rrr,
		rcDependencyData: rdd,
		log:              log.NewHelper(logger),
	}
}

func (s *RcRdmServiceServicer) GetAhpResult(ctx context.Context, req *pb.GetAhpResultDetailReq) (*pb.AhpResultResp, error) {
	_, err := s.rcDependencyData.Get(ctx, req.DepId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return &pb.AhpResultResp{
				Available: false,
				Msg:       "没有权限下载该报告",
			}, nil
		}
		return nil, err
	}

	rrr, err := s.rcRdmResult.GetUpToDate(ctx, req.DepId)
	if err != nil {
		return nil, err
	}

	rrrdList, err := s.rcRdmResDetail.GetByResIdAndLevel(ctx, rrr.Id, int(req.Level))
	if err != nil {
		return nil, err
	}

	data := make([]*pb.AhpResDetail, 0)
	for _, rrrd := range *rrrdList {
		data = append(data, &pb.AhpResDetail{
			Field: rrrd.Field,
			Score: float32(rrrd.Score),
			Level: int32(rrrd.Level),
		})
	}
	return &pb.AhpResultResp{
		Available: true,
		Msg:       "success",
		Data:      data,
	}, nil
}
