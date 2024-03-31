package dto

type RcRiskIndex struct {
	BaseModel
	ContentId  int64  `json:"content_id"`
	RiskDec    string `json:"risk_dec"`
	IndexDec   string `json:"index_dec"`
	IndexValue string `json:"index_value"`
	IndexFlag  string `json:"index_flag"`
}

func (RcRiskIndex) TableName() string {
	return "rc_risk_index"
}
