package dto

type RcDependencyData struct {
	BaseModel
	ContentId       *int64
	AttributedMonth *string
	UscId           string
	LhQylx          int
	LhCylwz         int
	LhGdct          int
	LhQybq          int
	LhYhsx          int
	LhSfsx          int
	AdditionData    string
}

func (*RcDependencyData) TableName() string {
	return "rc_dependency_data"
}
