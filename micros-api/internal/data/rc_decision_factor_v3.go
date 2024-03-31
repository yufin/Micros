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

type RcDecisionFactorV3Repo struct {
	data *Data
	log  *log.Helper
}

func NewRcDecisionFactorV3Repo(data *Data, logger log.Logger) biz.RcDecisionFactorV3Repo {
	return &RcDecisionFactorV3Repo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *RcDecisionFactorV3Repo) GetByContentIdWithDataScope(ctx context.Context, contentId int64) (*dto.RcDecisionFactorClaimed, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	var tb dto.RcContentFactorClaimV3
	var dataRpc dto.RcDecisionFactorClaimed
	err = repo.data.Db.
		Table(tb.TableName()).
		Select("rc_decision_factor.*, rc_content_factor_claim_v3.id as claim_id, rc_content_factor_claim_v3.content_id as content_id").
		Joins("left join rc_decision_factor on rc_content_factor_claim_v3.factor_id = rc_decision_factor.id").
		Where("rc_decision_factor.user_id = ? and content_id = ? ", dsi.UserId, contentId).
		Where("rc_content_factor_claim_v3.deleted_at is null and rc_decision_factor.deleted_at is null").
		Order("created_at desc").
		First(&dataRpc).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = repo.data.Db.
				Table(tb.TableName()).
				Select("rc_decision_factor.*, rc_content_factor_claim_v3.id as claim_id, rc_content_factor_claim_v3.content_id as content_id").
				Joins("left join rc_decision_factor on rc_content_factor_claim_v3.factor_id = rc_decision_factor.id").
				Where("rc_decision_factor.user_id in (?) and content_id = ? ", dsi.AccessibleIds, contentId).
				Where("rc_content_factor_claim_v3.deleted_at is null and rc_decision_factor.deleted_at is null").
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
func (repo *RcDecisionFactorV3Repo) GetWithDataScope(ctx context.Context, id int64) (*dto.RcDecisionFactor, error) {
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

func (repo *RcDecisionFactorV3Repo) CheckContentIdAccessible(ctx context.Context, contentId int64) (bool, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return false, err
	}
	var tb dto.RcContentFactorClaimV3
	var count int64
	err = repo.data.Db.
		Table(tb.TableName()).
		Select("*").
		Joins("left join rc_decision_factor on rc_content_factor_claim_v3.factor_id = rc_decision_factor.id").
		Where("rc_decision_factor.user_id in (?) and content_id = ? ", dsi.AccessibleIds, contentId).
		Where("rc_content_factor_claim_v3.deleted_at is null and rc_decision_factor.deleted_at is null").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (repo *RcDecisionFactorV3Repo) CountByUscIdAndUserId(ctx context.Context, uscId string) (int64, error) {
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

func (repo *RcDecisionFactorV3Repo) GetLatestIdByUscIdAndUserId(ctx context.Context, uscId string) (int64, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return 0, err
	}
	var model dto.RcDecisionFactor
	err = repo.data.Db.
		Table(model.TableName()).
		Where("usc_id = ? and user_id = ?", uscId, dsi.UserId).
		Order("created_at desc").
		First(&model).
		Error

	if err != nil {
		return 0, err
	}
	return model.Id, nil
}

func (repo *RcDecisionFactorV3Repo) Insert(ctx context.Context, data *dto.RcDecisionFactor) (int64, error) {
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

func (repo *RcDecisionFactorV3Repo) InsertClaimNoDupe(ctx context.Context, data *dto.RcContentFactorClaimV3) (int64, error) {
	// check if row exists
	var record dto.RcContentFactorClaimV3
	err := repo.data.Db.
		Model(&dto.RcContentFactorClaimV3{}).
		Where("content_id = ? and factor_id = ?", data.ContentId, data.FactorId).
		First(&record).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			data.BaseModel.Gen()
			err := repo.data.Db.
				Model(&dto.RcContentFactorClaimV3{}).
				Create(data).Error
			if err != nil {
				return 0, err
			}
			return data.Id, nil
		}
		return 0, err
	}
	return record.Id, nil
}

func (repo *RcDecisionFactorV3Repo) InsertClaim(ctx context.Context, data *dto.RcContentFactorClaimV3) (int64, error) {
	data.BaseModel.Gen()
	err := repo.data.Db.
		Model(&dto.RcContentFactorClaimV3{}).
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

func (repo *RcDecisionFactorV3Repo) ListReportClaimed(ctx context.Context, page *dto.PaginationReq, kwd string) (*[]dto.ListReportInfo, dto.PaginationInfo, error) {
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
					rc_content_meta.usc_id AS usc_id,
					rc_content_meta.id AS content_id,
					rc_content_meta.enterprise_name AS enterprise_name,
					rc_content_meta.attribute_month AS data_collect_month,
					rc_content_meta.version,
					rc_decision_factor.id AS factor_id,
					rc_decision_factor.lh_qylx AS lh_qylx,
					rc_content_meta.created_at AS created_at,
					ROW_NUMBER() OVER (
						PARTITION BY rc_content_meta.id
						ORDER BY rc_decision_factor.created_at
					) AS rn
				FROM rc_decision_factor
				JOIN rc_content_factor_claim_v3
				ON rc_decision_factor.id = rc_content_factor_claim_v3.factor_id
				AND rc_content_factor_claim_v3.deleted_at IS NULL
				JOIN rc_content_meta
				ON rc_content_factor_claim_v3.content_id = rc_content_meta.id
				AND rc_content_meta.deleted_at IS NULL
				WHERE rc_decision_factor.user_id in ? %s
			)
			SELECT *,
			COUNT(*) OVER () AS total
			FROM ordered_rows
			WHERE rn = 1 and version = 'V2.0'
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

func (repo *RcDecisionFactorV3Repo) ListReportClaimedByUscId(ctx context.Context, page *dto.PaginationReq, uscId string, version string) (*[]dto.ListReportInfo, int64, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page.PageNum - 1) * page.PageSize
	list := make([]dto.ListReportInfo, 0)

	err = repo.data.Db.
		Raw(
			`WITH ordered_rows AS (
				SELECT
					rc_content_meta.usc_id AS usc_id,
					rc_content_meta.id AS content_id,
					rc_content_meta.enterprise_name AS enterprise_name,
					rc_content_meta.attribute_month AS data_collect_month,
					rc_content_meta.version,
					rc_content_meta.status,
					rc_decision_factor.id AS factor_id,
					rc_decision_factor.lh_qylx AS lh_qylx,
					rc_decision_factor.user_id AS user_id,
					rc_content_meta.created_at AS created_at,
					ROW_NUMBER() OVER (
						PARTITION BY rc_content_meta.id
						ORDER BY rc_decision_factor.created_at
					) AS rn
				FROM rc_decision_factor
				JOIN rc_content_factor_claim_v3
				ON rc_decision_factor.id = rc_content_factor_claim_v3.factor_id
				AND rc_content_factor_claim_v3.deleted_at IS NULL
				JOIN rc_content_meta
				ON rc_content_factor_claim_v3.content_id = rc_content_meta.id
				AND rc_content_meta.deleted_at IS NULL
				WHERE rc_decision_factor.user_id in ?
				and rc_decision_factor.usc_id = ?
			)
			SELECT *,
			COUNT(*) OVER () AS total
			FROM ordered_rows
			WHERE rn = 1 and version = ?
			order by data_collect_month desc, created_at desc
			limit ? offset ?;`, dsi.AccessibleIds, uscId, version, page.PageSize, offset).
		Scan(&list).
		Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if len(list) > 0 {
		total = list[0].Total
	}
	return &list, total, nil
}

func (repo *RcDecisionFactorV3Repo) ListCompaniesWaiting(ctx context.Context, page *dto.PaginationReq, kwd string, version string) (*[]dto.ListCompaniesWaitingResp, int64, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page.PageNum - 1) * page.PageSize
	kwdStmt := ""
	if kwd != "" {
		kwdStmt = "and usc_id like '%" + kwd + "%'"
	}

	list := make([]dto.ListCompaniesWaitingResp, 0)

	err = repo.data.Db.
		Raw(
			fmt.Sprintf(
				`with res as (select usc_id,
										created_at,
										ROW_NUMBER() over (partition by usc_id order by created_at ) as rn
								 from rc_decision_factor
								 where not exists(select 1
												  from rc_content_meta
												  where rc_content_meta.usc_id = rc_decision_factor.usc_id
													and rc_content_meta.version = ?) 
								and user_id in ? %s)
					select usc_id,
						   created_at,
						   count(*) over () as total
					from res
					where rn = 1
					limit ? offset ?;`, kwdStmt), version, dsi.AccessibleIds, page.PageSize, offset,
		).Scan(&list).Error

	if err != nil {
		return nil, 0, err
	}
	var total int64
	if len(list) > 0 {
		total = list[0].Total
	}
	return &list, total, nil

}

func (repo *RcDecisionFactorV3Repo) ListCompanies(ctx context.Context, page *dto.PaginationReq, kwd string, version string) (*[]dto.ListCompaniesLatest, int64, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page.PageNum - 1) * page.PageSize
	list := make([]dto.ListCompaniesLatest, 0)
	kwdStmt := ""
	if kwd != "" {
		kwdStmt = "and enterprise_name like '%" + kwd + "%'"
	}
	err = repo.data.Db.
		Raw(
			fmt.Sprintf(
				`with res as (SELECT t.enterprise_name,
										t.usc_id,
										t.created_at as last_created_at
								 FROM (SELECT m.enterprise_name,
											  m.usc_id,
											  m.created_at,
											  ROW_NUMBER() OVER (PARTITION BY m.usc_id ORDER BY m.attribute_month desc) AS rn
									   FROM rc_content_meta m
												inner JOIN (SELECT DISTINCT usc_id
															FROM rc_decision_factor
															WHERE user_id IN ? 
															and deleted_at is null) f ON m.usc_id = f.usc_id
									   WHERE m.version = ?
										and deleted_at is null %s) t
								 WHERE t.rn = 1
								 order by created_at desc)
					select *,
						   count(*) over () as total
					from res
					limit ? offset ?;`, kwdStmt), dsi.AccessibleIds, version, page.PageSize, offset).
		Scan(&list).
		Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if len(list) > 0 {
		total = list[0].Total
	}
	return &list, total, nil
}

func (repo *RcDecisionFactorV3Repo) GetClaimRecord(ctx context.Context, claimId int64) (*dto.RcContentFactorClaimV3, error) {
	var dataRpc dto.RcContentFactorClaimV3
	err := repo.data.Db.
		Model(&dto.RcContentFactorClaimV3{}).
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

func (repo *RcDecisionFactorV3Repo) SyncClaimed(ctx context.Context, uscId string, version string) error {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return err
	}
	unsyncedContentIds := make([]int64, 0)
	err = repo.data.Db.Raw(
		`WITH
		content_ids AS (
			SELECT id
			FROM rc_content_meta
			WHERE usc_id = ?
			AND version = ?
			AND deleted_at is null
		),
		factor_ids AS (
			SELECT id
			FROM rc_decision_factor
			WHERE usc_id = ?
			AND user_id IN ?
			AND deleted_at is null
			ORDER BY created_at DESC
		),
		claimed_content_ids AS (
			SELECT cfc.content_id
			FROM rc_content_factor_claim_v3 cfc
			JOIN factor_ids fi ON cfc.factor_id = fi.id
			where deleted_at is null
		)
		SELECT ci.*
		FROM content_ids ci
		LEFT JOIN claimed_content_ids cci ON ci.id = cci.content_id
		WHERE cci.content_id IS NULL;`, uscId, version, uscId, dsi.AccessibleIds).
		Pluck("id", &unsyncedContentIds).Error
	if err != nil {
		return err
	}

	if len(unsyncedContentIds) == 0 {
		return nil
	}

	factorIds := make([]int64, 0)
	err = repo.data.Db.Raw(
		`SELECT id
			FROM rc_decision_factor
			WHERE usc_id = ?
			  AND user_id IN ?
			  AND deleted_at is null
			ORDER BY created_at DESC`, uscId, dsi.AccessibleIds).
		Pluck("id", &factorIds).
		Error
	if err != nil {
		return err
	}

	if len(factorIds) == 0 {
		return nil
	}

	insertReqs := make([]dto.RcContentFactorClaimV3, 0)
	for _, v := range unsyncedContentIds {
		req := dto.RcContentFactorClaimV3{
			ContentId: v,
			FactorId:  factorIds[0],
		}
		req.BaseModel.Gen()
		insertReqs = append(insertReqs, req)
	}

	if err := repo.data.Db.Create(&insertReqs).Error; err != nil {
		return err
	}

	return nil
}
