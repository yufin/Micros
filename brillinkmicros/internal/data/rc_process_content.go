package data

import (
	"brillinkmicros/internal/biz"
	"context"
	"encoding/binary"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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
func (repo *RcProcessedContentRepo) Get(ctx context.Context, id int64) (*biz.RcProcessedContent, error) {
	var dataRpc *biz.RcProcessedContent
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
func (repo *RcProcessedContentRepo) GetByContentIdUpToDate(ctx context.Context, contentId int64) (*biz.RcProcessedContent, error) {
	var dataRpc *biz.RcProcessedContent
	err := repo.data.Db.
		Table(dataRpc.TableName()).
		Where("content_id = ?", contentId).
		Order("updated_at desc").
		First(&dataRpc).
		Error
	repo.log.WithContext(ctx).Infof("RcProcessedContentRepo biz.GetList %v", contentId)
	if err != nil {
		return nil, err
	}
	return dataRpc, nil
}

func (repo *RcProcessedContentRepo) RefreshReportContent(ctx context.Context, contentId int64) (bool, error) {
	err := func() error {
		msg := make([]byte, 8)
		binary.BigEndian.PutUint64(msg, uint64(contentId))
		_, err := repo.data.Nw.js.Publish("task.rskc.content.process.newId", msg)
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
