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

// implement biz.RcProcessedContentRepo

func NewRcProcessedContentRepo(data *Data, logger log.Logger) biz.RcProcessedContentRepo {
	return &RcProcessedContentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// RcProcessedContentRepo is a data access object for table rc_processed_content.
// impl of biz.RcProcessedContentRepo
type RcProcessedContentRepo struct {
	data *Data
	log  *log.Helper
}

// Get .
// impl of biz.RcProcessedContentRepo Get
// db operation
func (repo *RcProcessedContentRepo) Get(ctx context.Context, id int64) (*dto.RcProcessedContent, error) {
	var dataRpc *dto.RcProcessedContent
	err := repo.data.Db.
		Table(dataRpc.TableName()).
		Where("id = ?", id).
		First(&dataRpc).
		Error
	repo.log.WithContext(ctx).Infof("RcProcessedContentRepo biz.Get %d", id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dataRpc, nil
}

// GetByContentIdUpToDate .
// impl of biz.RcProcessedContentRepo GetList
// db operation
func (repo *RcProcessedContentRepo) GetByContentIdUpToDate(ctx context.Context, contentId int64) (*dto.RcProcessedContent, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	var dataRpc *dto.RcProcessedContent
	err = repo.data.Db.
		Table(fmt.Sprintf("%s as rpc", dataRpc.TableName())).
		Select("rpc.*").
		Joins("INNER JOIN rc_dependency_data AS rdd ON rpc.content_id = rdd.content_id").
		Where("rpc.content_id = ?", contentId).
		Where("rdd.create_by IN (?)", dsi.AccessibleIds).
		Order("rpc.created_at desc").
		Limit(1).
		First(&dataRpc).
		Error

	repo.log.WithContext(ctx).Infof("RcProcessedContentRepo biz.GetList %v", contentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dataRpc, nil
}

// GetByContentIdUpToDateByUser .
// impl of biz.RcProcessedContentRepo GetList
// db operation
func (repo *RcProcessedContentRepo) GetContentUpToDateByDepId(ctx context.Context, depId int64, allowedUserId int64) (*dto.RcProcessedContent, error) {
	dsi, err := pkg.ParseBlDataScope(ctx)
	if err != nil {
		return nil, err
	}
	if dsi.UserId != allowedUserId {
		return nil, errors.New(400, "not allowed", "")
	}
	var dataRpc *dto.RcProcessedContent
	err = repo.data.Db.
		Table(fmt.Sprintf("%s as rpc", dataRpc.TableName())).
		Select("rpc.*").
		Joins("INNER JOIN rc_dependency_data AS rdd ON rpc.content_id = rdd.content_id").
		Where("rdd.id = ?", depId).
		Order("rpc.created_at desc").
		Limit(1).
		First(&dataRpc).
		Error

	repo.log.WithContext(ctx).Infof("RcProcessedContentRepo biz.GetContentUpToDateByDepId %v", depId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return dataRpc, nil
}

func (repo *RcProcessedContentRepo) RefreshReportContent(ctx context.Context, contentId int64) (bool, error) {
	err := func() error {
		msg := make([]byte, 8)
		binary.BigEndian.PutUint64(msg, uint64(contentId))
		_, err := repo.data.Nw.Js.Publish("task.rc.content.newId", msg)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return false, err
	}
	return true, nil
}
