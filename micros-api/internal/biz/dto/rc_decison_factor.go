package dto

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
}
