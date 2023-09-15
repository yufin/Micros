package data

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"micros-dw/internal/biz"
	"micros-dw/internal/biz/dto"
)

type DwEnterpriseDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewDwEnterpriseRepo(data *Data, logger log.Logger) biz.DwEnterpriseRepo {
	return &DwEnterpriseDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *DwEnterpriseDataRepo) GetEntIdent(ctx context.Context, name string) (string, error) {
	var uscId *string
	err := repo.data.Dbs.Db.
		Table("enterprise_wait_list").
		Select("usc_id").
		Where("enterprise_name = ?", name).
		Order("created_at desc").
		Limit(1).
		Scan(&uscId).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	if uscId == nil {
		return "", nil
	}
	return *uscId, nil
}

func (repo *DwEnterpriseDataRepo) GetEntInfo(ctx context.Context, uscId string) (*dto.EnterpriseInfo, error) {
	var data dto.EnterpriseInfo
	err := repo.data.Dbs.Db.
		Model(&dto.EnterpriseInfo{}).
		Where("usc_id = ?", uscId).
		Order("created_at desc").
		First(&data).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (repo *DwEnterpriseDataRepo) GetEntCredential(ctx context.Context, uscId string) (*[]dto.EnterpriseCertification, error) {
	data := make([]dto.EnterpriseCertification, 0)
	err := repo.data.Dbs.Db.
		Model(&dto.EnterpriseCertification{}).
		Where("usc_id = ?", uscId).
		Order("certification_date desc").
		Scan(&data).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (repo *DwEnterpriseDataRepo) GetEntRankingList(ctx context.Context, uscId string) (*[]dto.EnterpriseRankingList, error) {
	data := make([]dto.EnterpriseRankingList, 0)
	err := repo.data.Dbs.Db.
		Raw(
			`select usc_id,
			ranking_position,
			list_title,
			list_type,
			list_source,
			list_participants_total,
			list_published_date,
			list_url_qcc,
			list_url_origin
			from enterprise_ranking
					 left join ranking_list on enterprise_ranking.list_id = ranking_list.id
			where usc_id = ?
			order by list_published_date desc`, uscId).
		Scan(&data).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (repo *DwEnterpriseDataRepo) GetEntIndustry(ctx context.Context, uscId string) (*[]string, error) {
	var dataStr string
	err := repo.data.Dbs.Db.
		Table("enterprise_industry").
		Select("industry_data").
		Where("usc_id = ?", uscId).
		Order("created_at desc").
		Limit(1).
		Scan(&dataStr).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	data := make([]string, 0)
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *DwEnterpriseDataRepo) GetEntProduct(ctx context.Context, uscId string) (*[]string, error) {
	var dataStr string
	err := repo.data.Dbs.Db.
		Table("enterprise_product").
		Select("product_data").
		Where("usc_id = ?", uscId).
		Order("created_at desc").
		Limit(1).
		Scan(&dataStr).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	data := make([]string, 0)
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}
