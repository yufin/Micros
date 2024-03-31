package dto

import "time"

type RcDecisionFactor struct {
	BaseModel
	UscId   string
	LhQylx  int
	LhCylwz int
	LhGdct  int
	LhQybq  int
	LhYhsx  int
	LhSfsx  int
	UserId  int64
}

func (*RcDecisionFactor) TableName() string {
	return "rc_decision_factor"
}

type ListReportInfo struct {
	UscId            string
	ContentId        int64
	EnterpriseName   string
	DataCollectMonth string
	FactorId         int64
	LhQylx           int
	Total            int64
	Version          string
	CreatedAt        time.Time
	Status           int32
	UserId           int64
	Username         string
}

type ListCompaniesInfo struct {
	UscId          string
	EnterpriseName string
	AttributeMonth string
	Total          int64
}

type ListCompaniesLatest struct {
	UscId          string
	EnterpriseName string
	LastCreatedAt  time.Time
	Total          int64
}

type RcDecisionFactorClaimed struct {
	BaseModel
	UscId     string
	LhQylx    int
	LhCylwz   int
	LhGdct    int
	LhQybq    int
	LhYhsx    int
	LhSfsx    int
	UserId    int64
	ClaimId   int64
	ContentId int64
}

type ListCompaniesWaitingResp struct {
	UscId     string    `json:"usc_id"`
	CreatedAt time.Time `json:"created_at"`
	Total     int64     `json:"total"`
}
