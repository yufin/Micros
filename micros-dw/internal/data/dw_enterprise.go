package data

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"micros-dw/internal/biz"
	"micros-dw/internal/biz/dto"
	"time"
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
	err := repo.data.Mongo.Client.Database("dw").Collection("basic_penetration").FindOne(
		context.TODO(),
		bson.M{"usc_id": uscId},
	).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	detail := data["penetration_tree"]

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
		Conclusion: data["penetration_conclusion"].(string),
		Data:       detailArr,
		UscId:      data["usc_id"].(string),
	}, nil
}

func (repo *DwEnterpriseDataRepo) GetShareholders(ctx context.Context, uscId string) (*[]dto.EnterpriseShareholder, error) {
	var res bson.M
	if err := repo.data.Mongo.Client.Database("dw").
		Collection("basic_shareholders").
		FindOne(context.TODO(), bson.M{"usc_id": uscId}).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	// bson.M to json
	var shareholder struct {
		Shareholder []dto.EnterpriseShareholder `bson:"shareholders"`
	}

	b, err := bson.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(b, &shareholder)
	if err != nil {
		return nil, err
	}
	return &shareholder.Shareholder, nil
}

func (repo *DwEnterpriseDataRepo) GetInvestments(ctx context.Context, uscId string) (*[]dto.EnterpriseInvestment, error) {
	var res bson.M
	if err := repo.data.Mongo.Client.Database("dw").
		Collection("basic_investment").
		FindOne(context.TODO(), bson.M{"usc_id": uscId}).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	var investments struct {
		Investment []dto.EnterpriseInvestment `bson:"investment"`
	}

	b, err := bson.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(b, &investments)
	if err != nil {
		return nil, err
	}
	return &investments.Investment, nil
}

func (repo *DwEnterpriseDataRepo) GetBranches(ctx context.Context, uscId string) (*[]dto.EnterpriseBranches, error) {
	var res bson.M
	if err := repo.data.Mongo.Client.Database("dw").
		Collection("basic_branchlist").
		FindOne(context.TODO(), bson.M{"usc_id": uscId},
			options.FindOne().SetProjection(bson.M{"branch_list": 1, "_id": 0})).Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	var branches struct {
		Branches []dto.EnterpriseBranches `bson:"branch_list"`
	}

	b, err := bson.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(b, &branches)
	if err != nil {
		return nil, err
	}
	return &branches.Branches, nil
}

func (repo *DwEnterpriseDataRepo) GetDocInDuration(ctx context.Context, uscId string, tp time.Time, coll string) (bson.M, error) {
	var res bson.M

	filter := bson.D{
		{"$and", []bson.D{
			{{"date", bson.D{{"$lte", tp}}}},
			{{"check_date", bson.D{{"$gte", tp}}}},
			{{"usc_id", uscId}},
		}},
	}

	err := repo.data.Mongo.Client.Database("dw").
		Collection(coll).
		FindOne(context.TODO(), filter).
		Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return res, nil
}

func (repo *DwEnterpriseDataRepo) GetDocInExtendDuration(ctx context.Context, uscId string, tp time.Time, collName string, extendDate int) (bson.M, error) {
	var res bson.M

	pipeline := mongo.Pipeline{
		{
			{"$addFields", bson.D{
				{"valid_time_end", bson.D{{"$add", bson.A{"$check_date", extendDate * 24 * 60 * 60000}}}},
			}},
		},
		{
			{"$match", bson.D{
				{"$and", bson.A{
					bson.D{{"date", bson.D{{"$lte", tp}}}},
					bson.D{{"valid_time_end", bson.D{{"$gte", tp}}}},
					bson.D{{"usc_id", uscId}},
				}},
			}},
		},
		{
			{"$project", bson.D{
				{"valid_time_end", 0},
			}},
		},
	}
	cur, err := repo.data.Mongo.Client.Database("dw").Collection(collName).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if cur.TryNext(context.TODO()) {
		if err := cur.Decode(&res); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, nil
			}
			return nil, err
		}
	}
	return res, nil
}

func (repo *DwEnterpriseDataRepo) GetDocWithDuration(ctx context.Context, uscId string, tp time.Time, coll string, extendDate int) (map[string]any, error) {
	res, err := repo.GetDocInDuration(ctx, uscId, tp, coll)
	if err != nil {
		return nil, err
	}
	if res == nil {
		res, err = repo.GetDocInExtendDuration(ctx, uscId, tp, coll, extendDate)
		if err != nil {
			return nil, err
		}
	}
	var m map[string]any
	if res != nil {
		b, err := bson.MarshalExtJSON(res, false, false)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
