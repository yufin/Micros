package dto

type RcReportOss struct {
	BaseModel
	DepId   int64
	OssId   int64
	Version int
}

func (*RcReportOss) TableName() string {
	return "rc_report_oss"
}
