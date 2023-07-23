package dto

type RcRdmResDetail struct {
	BaseModel
	ResId int64   `gorm:"column:res_id"`
	Field string  `gorm:"column:field"`
	Level int     `gorm:"column:level"`
	Score float64 `gorm:"column:score"`
}

func (*RcRdmResDetail) TableName() string {
	return "rc_rdm_result"
}
