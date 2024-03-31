package dto

type RcContentFactorClaim struct {
	BaseModel
	ContentId int64 `json:"content_id" gorm:"column:content_id" xlsx:"content_id"`
	FactorId  int64 `json:"factor_id" gorm:"column:factor_id" xlsx:"factor_id"`
}

func (*RcContentFactorClaim) TableName() string {
	return "rc_content_factor_claim"
}

type RcContentFactorClaimV3 struct {
	BaseModel
	ContentId int64 `json:"content_id" gorm:"column:content_id" xlsx:"content_id"`
	FactorId  int64 `json:"factor_id" gorm:"column:factor_id" xlsx:"factor_id"`
}

func (*RcContentFactorClaimV3) TableName() string {
	return "rc_content_factor_claim_v3"
}
