package biz

import (
	"time"
)

type RcOriginContentGetPageResp struct {
	PaginationResp
	Data *[]RcOriginContent
}

type RcOriginContentInfosResp struct {
	PaginationResp
	Data *[]RcOriginContentInfo
}

type RcOriginContentInfo struct {
	ContentId        int64
	UscId            string
	DataCollectMonth string
	Content          string
	StatusCode       int
	//EnterpriseName     string
	ProcessedId        int64
	ProcessedUpdatedAt time.Time
	UpdatedAt          time.Time
	CreatedAt          time.Time
}
