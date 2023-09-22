package dto

type RcRdmResult struct {
	BaseModel
	AppType int
	DepId   int64
	Comment string
}

func (*RcRdmResult) TableName() string {
	return "rc_rdm_result"
}
