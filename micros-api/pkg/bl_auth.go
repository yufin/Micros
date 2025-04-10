package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/transport/http"
)

const BlDataScopeHeaderKey = "BL-DATA-SCOPES"
const DataScopeFullAccess = 1

type DataScopeInfo struct {
	UserId        int64   `json:"userId"`
	AccessType    int     `json:"accessType"`
	AccessibleIds []int64 `json:"accessibleIds"`
}

func ParseBlDataScope(ctx context.Context) (*DataScopeInfo, error) {
	req, ok := http.RequestFromServerContext(ctx)
	if !ok {
		return nil, errors.New("error: request from context not found")
	}
	ds := req.Header.Get(BlDataScopeHeaderKey)
	var dsi DataScopeInfo
	if err := json.Unmarshal([]byte(ds), &dsi); err != nil {
		return nil, err
	}

	dsi.AccessibleIds = func(sli []int64, n int64) []int64 {
		for _, value := range sli {
			if value == n {
				return sli
			}
		}
		return append(sli, n)
	}(dsi.AccessibleIds, dsi.UserId)

	return &dsi, nil
}

func ApplyBlDataScope(dsi *DataScopeInfo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if dsi.AccessType == DataScopeFullAccess {
			return db.Where("create_by IN (?)", dsi.AccessibleIds)
		}
		return db
	}
}
