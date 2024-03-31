package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"micros-api/internal/biz/dto"
	"time"
)

type ArtifactDataRepo interface {
	SaveMany(ctx context.Context, db string, coll string, req []interface{}) error
	Get(ctx context.Context, db string, coll string, pageSize int64, pageNum int64, orderBy string, filter bson.D, result any) error
	CountDoc(ctx context.Context, filter bson.D, coll string, db string) (int64, error)
	DeleteOne(ctx context.Context, db string, coll string, id string) error
}

type ArtifactDataUsecase struct {
	repo ArtifactDataRepo
	log  *log.Helper
}

func NewArtifactDataUsecase(repo ArtifactDataRepo, logger log.Logger) *ArtifactDataUsecase {
	return &ArtifactDataUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ArtifactDataUsecase) InsertEnterpriseComment(ctx context.Context, req []dto.EnterpriseCommentInsertReq) error {
	uc.log.WithContext(ctx).Infof("biz.ArtifactDataUsecase.InsertEnterpriseComment len:%v", len(req))
	reqTrans := make([]interface{}, 0)
	for _, v := range req {
		v.CreateAt = time.Now().Local()
		reqTrans = append(reqTrans, v)
	}
	return uc.repo.SaveMany(ctx, "artifacts", "enterprise_comment", reqTrans)
}

func (uc *ArtifactDataUsecase) InsertProductEvalRule(ctx context.Context, req []dto.ProductEvalRuleInsertResp) error {
	uc.log.WithContext(ctx).Infof("biz.ArtifactDataUsecase.InsertProductEvalRule len:%v", len(req))
	reqTrans := make([]interface{}, 0)
	for _, v := range req {
		v.CreateAt = time.Now().Local()
		reqTrans = append(reqTrans, v)
	}
	return uc.repo.SaveMany(ctx, "artifacts", "product_eval_rule", reqTrans)
}

func (uc *ArtifactDataUsecase) GetProductEvalRule(ctx context.Context, pageSize int64, pageNum int64, filter bson.D) ([]map[string]any, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.ArtifactDataUsecase.GetProductEvalRule pageSize:%v pageNum:%v filter:%v", pageSize, pageNum, filter)
	res := make([]dto.ProductEvalRule, 0)
	err := uc.repo.Get(ctx, "artifacts", "product_eval_rule", pageSize, pageNum, "create_at", filter, &res)
	if err != nil {
		return nil, 0, err
	}
	count, err := uc.repo.CountDoc(ctx, filter, "product_eval_rule", "artifacts")
	if err != nil {
		return nil, 0, err
	}
	resTrans := make([]map[string]any, 0)
	b, err := json.Marshal(res)
	if err != nil {
		return nil, 0, err
	}
	err = json.Unmarshal(b, &resTrans)
	if err != nil {
		return nil, 0, err
	}
	return resTrans, count, nil
}

func (uc *ArtifactDataUsecase) GetEnterpriseComment(ctx context.Context, pageSize int64, pageNum int64, filter bson.D) ([]map[string]any, int64, error) {
	uc.log.WithContext(ctx).Infof("biz.ArtifactDataUsecase.GetEnterpriseComment pageSize:%v pageNum:%v filter:%v", pageSize, pageNum, filter)
	res := make([]dto.EnterpriseComment, 0)
	err := uc.repo.Get(ctx, "artifacts", "enterprise_comment", pageSize, pageNum, "create_at", filter, &res)
	if err != nil {
		return nil, 0, err
	}
	count, err := uc.repo.CountDoc(ctx, filter, "enterprise_comment", "artifacts")
	if err != nil {
		return nil, 0, err
	}
	resTrans := make([]map[string]any, 0)
	b, err := json.Marshal(res)
	if err != nil {
		return nil, 0, err
	}
	err = json.Unmarshal(b, &resTrans)
	if err != nil {
		return nil, 0, err
	}
	return resTrans, count, nil
}

func (uc *ArtifactDataUsecase) DeleteProductEvalRuleById(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("biz.ArtifactDataUsecase.DeleteProductEvalRuleById id:%v", id)
	return uc.repo.DeleteOne(ctx, "artifacts", "product_eval_rule", id)
}

func (uc *ArtifactDataUsecase) DeleteEnterpriseComment(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("biz.ArtifactDataUsecase.DeleteProductEvalRuleById id:%v", id)
	return uc.repo.DeleteOne(ctx, "artifacts", "enterprise_comment", id)
}
