package data

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"micros-api/internal/biz"
	"micros-api/internal/biz/dto"
	"micros-api/pkg"
)

type RcDecisionFactorRepo struct {
	data *Data
	log  *log.Helper
}

func NewRcDecisionFactorRepo(data *Data, logger log.Logger) biz.RcDecisionFactorRepo {
	return &RcDecisionFactorRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *RcDecisionFactorRepo) GetByContentIdWithDataScope(ctx context.Context, contentId int64) (*dto.RcDecisionFactorClaimed, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	var tb dto.RcContentFactorClaim
	var dataRpc dto.RcDecisionFactorClaimed
	err = repo.data.Db.
		Table(tb.TableName()).
		Select("rc_decision_factor.*, rc_content_factor_claim.id as claim_id, rc_content_factor_claim.content_id as content_id").
		Joins("left join rc_decision_factor on rc_content_factor_claim.factor_id = rc_decision_factor.id").
		Where("rc_decision_factor.user_id = ? and content_id = ? ", dsi.UserId, contentId).
		Where("rc_content_factor_claim.deleted_at is null and rc_decision_factor.deleted_at is null").
		Order("created_at desc").
		First(&dataRpc).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = repo.data.Db.
				Table(tb.TableName()).
				Select("rc_decision_factor.*, rc_content_factor_claim.id as claim_id, rc_content_factor_claim.content_id as content_id").
				Joins("left join rc_decision_factor on rc_content_factor_claim.factor_id = rc_decision_factor.id").
				Where("rc_decision_factor.user_id in (?) and content_id = ? ", dsi.AccessibleIds, contentId).
				Where("rc_content_factor_claim.deleted_at is null and rc_decision_factor.deleted_at is null").
				Order("created_at").
				First(&dataRpc).
				Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, nil
				}
				return nil, err
			}
			return &dataRpc, nil
		}
		return nil, err
	}

	return &dataRpc, nil
}

// GetWithDataScope 根据查询单条（需鉴权）
func (repo *RcDecisionFactorRepo) GetWithDataScope(ctx context.Context, id int64) (*dto.RcDecisionFactor, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	var dataRpc dto.RcDecisionFactor
	err = repo.data.Db.
		Model(&dto.RcDecisionFactor{}).
		Where("id = ? and user_id in ?", id, dsi.AccessibleIds).
		First(&dataRpc).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dataRpc, nil
}

func (repo *RcDecisionFactorRepo) CheckContentIdAccessible(ctx context.Context, contentId int64) (bool, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return false, err
	}
	var tb dto.RcContentFactorClaim
	var count int64
	err = repo.data.Db.
		Table(tb.TableName()).
		Select("*").
		Joins("left join rc_decision_factor on rc_content_factor_claim.factor_id = rc_decision_factor.id").
		Where("rc_decision_factor.user_id in (?) and content_id = ? ", dsi.AccessibleIds, contentId).
		Where("rc_content_factor_claim.deleted_at is null and rc_decision_factor.deleted_at is null").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (repo *RcDecisionFactorRepo) CountByUscIdAndUserId(ctx context.Context, uscId string) (int64, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return 0, err
	}
	var count int64
	var model dto.RcDecisionFactor
	err = repo.data.Db.
		Table(model.TableName()).
		Where("usc_id = ? and user_id = ?", uscId, dsi.UserId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RcDecisionFactorRepo) Insert(ctx context.Context, data *dto.RcDecisionFactor) (int64, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return 0, err
	}
	data.UserId = dsi.UserId
	data.BaseModel.Gen()
	err = repo.data.Db.
		Model(&dto.RcDecisionFactor{}).
		Create(data).Error
	if err != nil {
		return 0, err
	}
	return data.Id, nil
}

func (repo *RcDecisionFactorRepo) InsertClaimNoDupe(ctx context.Context, data *dto.RcContentFactorClaim) (int64, error) {
	// check if row exists
	var record dto.RcContentFactorClaim
	err := repo.data.Db.
		Model(&dto.RcContentFactorClaim{}).
		Where("content_id = ? and factor_id = ?", data.ContentId, data.FactorId).
		First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			data.BaseModel.Gen()
			err := repo.data.Db.
				Model(&dto.RcContentFactorClaim{}).
				Create(data).Error
			if err != nil {
				return 0, err
			}

			msg := make([]byte, 8)
			binary.BigEndian.PutUint64(msg, uint64(data.Id))
			_, err = repo.data.Nw.Js.Publish("task.rc.report.claimed.newId", msg)
			if err != nil {
				return 0, err
			}

			return data.Id, nil
		}
		return 0, err
	}
	return record.Id, nil
}

func (repo *RcDecisionFactorRepo) InsertClaim(ctx context.Context, data *dto.RcContentFactorClaim) (int64, error) {
	data.BaseModel.Gen()
	err := repo.data.Db.
		Model(&dto.RcContentFactorClaim{}).
		Create(data).Error
	if err != nil {
		return 0, err
	}

	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(data.Id))
	_, err = repo.data.Nw.Js.Publish("task.rc.report.claimed.newId", msg)
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (repo *RcDecisionFactorRepo) ListReportClaimed(ctx context.Context, page *dto.PaginationReq, kwd string) (*[]dto.ListReportInfo, dto.PaginationInfo, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, dto.PaginationInfo{}, err
	}

	offset := (page.PageNum - 1) * page.PageSize
	list := make([]dto.ListReportInfo, 0)
	kwdStmt := ""
	if kwd != "" {
		kwdStmt = "and enterprise_name like '%" + kwd + "%'"
	}

	err = repo.data.Db.
		Raw(
			fmt.Sprintf(
				`WITH ordered_rows AS (
				SELECT
					rc_origin_content.usc_id AS usc_id,
					rc_origin_content.id AS content_id,
					rc_origin_content.enterprise_name AS enterprise_name,
					rc_origin_content.year_month AS data_collect_month,
					rc_decision_factor.id AS factor_id,
					rc_decision_factor.lh_qylx AS lh_qylx,
					rc_origin_content.created_at AS created_at,
					ROW_NUMBER() OVER (
						PARTITION BY rc_origin_content.id
						ORDER BY rc_decision_factor.created_at
					) AS rn
				FROM rc_decision_factor
				JOIN rc_content_factor_claim
				ON rc_decision_factor.id = rc_content_factor_claim.factor_id
				AND rc_content_factor_claim.deleted_at IS NULL
				JOIN rc_origin_content
				ON rc_content_factor_claim.content_id = rc_origin_content.id
				AND rc_origin_content.deleted_at IS NULL
				WHERE rc_decision_factor.user_id in ? %s
			)
			SELECT *,
			COUNT(*) OVER () AS total
			FROM ordered_rows
			WHERE rn = 1
			order by data_collect_month desc, created_at desc
			limit ? offset ?;`, kwdStmt), dsi.AccessibleIds, page.PageSize, offset).
		Scan(&list).
		Error
	if err != nil {
		return nil, dto.PaginationInfo{}, err
	}
	var total int64
	if len(list) > 0 {
		total = list[0].Total
	}

	return &list, dto.PaginationInfo{
		Total:  total,
		Offset: int64(offset),
	}, nil
}

func (repo *RcDecisionFactorRepo) GetClaimRecord(ctx context.Context, claimId int64) (*dto.RcContentFactorClaim, error) {
	var dataRpc dto.RcContentFactorClaim
	err := repo.data.Db.
		Model(&dto.RcContentFactorClaim{}).
		Where("id = ?", claimId).
		First(&dataRpc).
		Error
	if err != nil {
		//if errors.Is(err, gorm.ErrRecordNotFound) {
		//	return nil, nil
		//}
		return nil, err
	}
	return &dataRpc, nil
}
