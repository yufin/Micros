package data

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"io"
	"micros-worker/internal/biz/dto"
	"micros-worker/internal/conf"
	"path"
	"regexp"
	"strings"
	"time"

	"micros-worker/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ContentRepo struct {
	data *Data
	log  *log.Helper
}

// NewContentRepo .
func NewContentRepo(data *Data, logger log.Logger) biz.ContentRepo {
	return &ContentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ContentRepo) filteringPathBase(pathArray []string, patternExpr string) []string {
	pattern := regexp.MustCompile(patternExpr)
	var tempDirs []string
	for _, dir := range pathArray {
		if pattern.MatchString(path.Base(dir)) {
			tempDirs = append(tempDirs, dir)
		}
	}
	return tempDirs
}

func (r *ContentRepo) GetRcContentMetas(ctx context.Context, pageReq dto.PaginationReq) ([]dto.RcContentMeta, error) {
	var contentMetas []dto.RcContentMeta
	var err error
	if pageReq.PageNum == 0 {
		err = r.data.Db.Db.Model(&dto.RcContentMeta{}).Find(&contentMetas).Error
	} else {
		err = r.data.Db.Db.
			Model(&dto.RcContentMeta{}).
			Offset((pageReq.PageNum - 1) * pageReq.PageSize).
			Limit(pageReq.PageSize).
			Find(&contentMetas).
			Error
	}
	return contentMetas, err
}

func (r *ContentRepo) PubPreLoanFileMeta(ctx context.Context, metaCh *chan dto.SftpContentFileMeta) error {
	sftpCli, err := r.data.SftpPool.GetClient()
	if err != nil {
		return errors.WithStack(err)
	}
	defer r.data.SftpPool.ReturnClient(sftpCli)

	paths, err := sftpCli.Glob("/taxDataPreloanFile/*")
	if err != nil {
		return err
	}
	matchedPaths := r.filteringPathBase(paths, `^[a-zA-Z0-9]{18}$`)

	//metas := make([]dto.SftpContentFileMeta, 0)
	for _, p := range matchedPaths {
		uscId := path.Base(p)
		uscSubDirs, err := sftpCli.Glob(p + "/*")
		if err != nil {
			return errors.WithStack(err)
		}
		matchedMonthDirs := r.filteringPathBase(uscSubDirs, `^[0-9]{6}$`)
		for _, monthDir := range matchedMonthDirs {
			jsonFilePaths, err := sftpCli.Glob(monthDir + "/*.json")
			if err != nil {
				return errors.WithStack(err)
			}
			for _, jsonFilePath := range jsonFilePaths {
				tempT, err := time.Parse("200601", path.Base(monthDir))
				if err != nil {
					return errors.WithStack(err)
				}
				fInfo := dto.SftpContentFileMeta{
					UscId:          uscId,
					AttributeMonth: tempT.Format("2006-01"),
					SourcePath:     jsonFilePath,
				}
				//metas = append(metas, fInfo)
				*metaCh <- fInfo
			}
		}
	}
	return nil
}

func (r *ContentRepo) PubPostLoanFileInfo(ctx context.Context, metaCh *chan dto.SftpContentFileMeta) error {
	sftpCli, err := r.data.SftpPool.GetClient()
	if err != nil {
		return errors.WithStack(err)
	}
	defer r.data.SftpPool.ReturnClient(sftpCli)

	paths, err := sftpCli.Glob("/taxDataPreloanFile/postloan/*")
	if err != nil {
		return errors.WithStack(err)
	}

	bwPaths := make([]string, 0)
	for _, p := range paths {
		bwPath := p + "/bw"
		_, err := sftpCli.Stat(bwPath)
		if err != nil {
			if !strings.Contains(err.Error(), "file does not exist") {
				return errors.WithStack(err)
			}
		} else {
			bwPaths = append(bwPaths, bwPath)
		}
	}
	reAttrMonth := regexp.MustCompile(`\b\d{6}\b`) // Matches exactly 6 digits
	reUscId := regexp.MustCompile(`^[A-Z0-9]{18}`)
	//metas := make([]dto.SftpContentFileMeta, 0)
	for _, bwPath := range bwPaths {
		files, err := sftpCli.ReadDir(bwPath)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, file := range files {
			// get yearmonth from bwPath path
			if strings.HasSuffix(file.Name(), ".json") {

				filePath := bwPath + "/" + file.Name()
				var attributeMonth string
				var uscId string
				matchesAttrM := reAttrMonth.FindStringSubmatch(bwPath)
				if len(matchesAttrM) > 0 {
					attributeMonth = matchesAttrM[0]
				} else {
					return errors.WithStack(err)
				}
				matchesUscId := reUscId.FindStringSubmatch(file.Name())
				if len(matchesUscId) > 0 {
					uscId = matchesUscId[0]
				} else {
					return errors.WithStack(err)
				}
				tempT, err := time.Parse("200601", attributeMonth)
				if err != nil {
					return errors.WithStack(err)
				}

				meta := dto.SftpContentFileMeta{
					UscId:          uscId,
					AttributeMonth: tempT.Format("2006-01"),
					SourcePath:     filePath,
				}
				//metas = append(metas, meta)
				*metaCh <- meta

			}

		}

	}
	return nil
}

func (r *ContentRepo) ReadContentByPath(ctx context.Context, path string) (*[]byte, error) {
	sftpCli, err := r.data.SftpPool.GetClient()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer r.data.SftpPool.ReturnClient(sftpCli)

	contentFile, err := sftpCli.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer contentFile.Close()
	b, err := io.ReadAll(contentFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &b, nil
}

func (r *ContentRepo) GetContentMetasByPath(ctx context.Context, sftpPath string) ([]dto.RcContentMeta, error) {
	meta := make([]dto.RcContentMeta, 0)
	err := r.data.Db.Db.
		Model(dto.RcContentMeta{}).
		Where("source_path = ?", sftpPath).
		Find(&meta).
		Error

	if err != nil {
		return nil, err
	}
	return meta, err
}

func (r *ContentRepo) CountContentMetaByPath(ctx context.Context, sftpPath string) (int64, error) {
	var count int64
	err := r.data.Db.Db.
		Model(dto.RcContentMeta{}).
		Where("source_path = ?", sftpPath).
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}
	return count, err
}

func (r *ContentRepo) InsertContentMeta(ctx context.Context, meta *dto.RcContentMeta) error {
	if meta.Id == 0 {
		meta.Gen()
	}
	err := r.data.Db.Db.
		Create(meta).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ContentRepo) DeleteRiskIndexByContentId(ctx context.Context, contentId int64) error {
	err := r.data.Db.Db.
		Unscoped().
		Where("content_id = ?", contentId).
		Delete(&dto.RcRiskIndex{}).
		Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil
		}
		return err
	}
	return nil
}

func (r *ContentRepo) InsertRawContent(ctx context.Context, meta *dto.RcContentMeta, content *map[string]any) (*mongo.InsertOneResult, error) {
	res, err := r.data.MongoDb.Client.
		Database("rc").
		Collection("raw_content").
		InsertOne(context.TODO(), bson.M{
			"content_id":      meta.Id,
			"version":         meta.Version,
			"usc_id":          meta.UscId,
			"enterprise_name": meta.EnterpriseName,
			"attribute_month": meta.AttributeMonth,
			"content":         content,
			"created_at":      time.Now().Local(),
			//"source_path":     meta.SourcePath,
		})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *ContentRepo) DeleteRawContentById(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.data.MongoDb.Client.
		Database("rc").
		Collection("raw_content").
		DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *ContentRepo) CountRawContentByPath(ctx context.Context, sftpPath string) (int64, error) {
	c, err := r.data.MongoDb.Client.
		Database("rc").
		Collection("raw_content").
		CountDocuments(context.TODO(), bson.M{"source_path": sftpPath})
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (r *ContentRepo) GetRawContentByContentId(ctx context.Context, contentId int64) (*map[string]any, error) {
	var res bson.M
	err := r.data.MongoDb.Client.
		Database("rc").
		Collection("raw_content").
		FindOne(ctx, bson.M{"content_id": contentId}).
		Decode(&res)
	if err != nil {
		return nil, err
	}
	content, ok := res["content"].(bson.M)
	if !ok {
		return nil, errors.New("content is not a map")
	}

	var m map[string]any
	b, err := bson.MarshalExtJSON(content, false, false)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *ContentRepo) SaveContentTransaction(ctx context.Context, meta *dto.RcContentMeta, content *map[string]any) error {
	meta.Gen()
	tx := r.data.Db.Db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	if err := tx.Create(meta).Error; err != nil {
		tx.Rollback()
		return err
	}
	docInsertRes, err := r.InsertRawContent(ctx, meta, content)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		if docInsertRes != nil && docInsertRes.InsertedID != nil {
			docId, ok := docInsertRes.InsertedID.(primitive.ObjectID)
			if ok {
				if delErr := r.DeleteRawContentById(ctx, docId); err != nil {
					log.Error("Failed to rollback MongoDB document:", delErr)
				}
			}
		}
		return err
	}
	return nil
}

func (r *ContentRepo) SaveRiskIndexesIdempotent(ctx context.Context, req []dto.RcRiskIndex) error {
	var count int64
	if len(req) == 0 {
		return nil
	}

	err := r.data.Db.Db.
		Model(&dto.RcRiskIndex{}).
		Where("content_id = ?", req[0].Id).
		Count(&count).
		Error
	if err != nil {
		return errors.WithStack(err)
	}

	if count > 0 {
		return nil
	}

	for i := range req {
		if req[i].Id == 0 {
			req[i].Gen()
		}
	}

	err = r.data.Db.Db.Create(&req).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ContentRepo) GetConf() *conf.Data {
	return r.data.DataConf
}
