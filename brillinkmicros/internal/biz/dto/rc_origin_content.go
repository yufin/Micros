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

func (*RcOriginContent) TableName() string {
	return "rc_origin_content"
}

type RcOriginContentInfosResp struct {
	PaginationResp
	Data *[]RcOriginContentInfo
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
	DepId              int64
	ProcessedUpdatedAt time.Time
	UpdatedAt          time.Time
	CreatedAt          time.Time
}

type RcOriginContentInfoV3 struct {
	ContentId          int64
	UscId              string
	DataCollectMonth   string
	LhQylx             int
	EnterpriseName     string
	ProcessedId        string
	DepId              int64
	ProcessedUpdatedAt time.Time
	UpdatedAt          time.Time
	CreatedAt          time.Time
}

type RcOriginContentInfosRespV3 struct {
	PaginationResp
	Data *[]RcOriginContentInfoV3
}
