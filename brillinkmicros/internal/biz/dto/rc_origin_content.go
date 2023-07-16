package dto

import "time"

type RcOriginContent struct {
	BaseModel
	UscId          string
	EnterpriseName string
	YearMonth      string
	Content        string
	StatusCode     int
}

type RcOriginContentInfosResp struct {
	PaginationResp
	Data *[]RcOriginContentInfo
}

func (*RcOriginContent) TableName() string {
	return "rskc_origin_content"
}

type RcOriginContentGetPageResp struct {
	PaginationResp
	Data *[]RcOriginContent
}

type RcOriginContentInfo struct {
	ContentId          int64
	UscId              string
	DataCollectMonth   string
	LhQylx             int
	EnterpriseName     string
	ProcessedId        int64
	ProcessedUpdatedAt time.Time
	UpdatedAt          time.Time
	CreatedAt          time.Time
}
