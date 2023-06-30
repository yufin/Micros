package common

import (
	"context"
	"encoding/json"
	"errors"

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
	return &dsi, nil
}
