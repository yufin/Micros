package data

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (repo *DwEnterpriseDataRepo) GetEntIdent(ctx context.Context, name string) (*dto.EnterpriseWaitList, error) {
	var data dto.EnterpriseWaitList
	err := repo.data.Dbs.Db.
		Model(&dto.EnterpriseWaitList{}).
		Where("enterprise_name = ?", name).
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
		return nil, err
	}
	if dataStr == "" {
		return nil, nil
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
		return nil, err
	}
	if dataStr == "" {
		return nil, nil
	}
	data := make([]string, 0)
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (repo *DwEnterpriseDataRepo) GetEntEquityTransparency(ctx context.Context, uscId string) (*dto.EnterpriseEquityTransparency, error) {
	data := bson.M{}
	err := repo.data.Mongo.Client.Database("spider").Collection("qcc_equity_penetration").FindOne(
		context.TODO(),
		bson.M{"usc_id": uscId},
	).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	detail := data["penitration"]

	jsonBytes, err := bson.MarshalExtJSON(bson.M{"arr": detail.(bson.A)}, false, false)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	err = json.Unmarshal(jsonBytes, &m)
	if err != nil {
		return nil, err
	}
	detailArr := make([]map[string]any, 0)
	for _, v := range m["arr"].([]any) {
		detailArr = append(detailArr, v.(map[string]any))
	}

	return &dto.EnterpriseEquityTransparency{
		Conclusion: data["conclusion"].(string),
		Data:       detailArr,
		UscId:      data["usc_id"].(string),
	}, nil
}
