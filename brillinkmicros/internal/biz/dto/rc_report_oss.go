package dto

type RcReportOss struct {
	BaseModel
	DepId   int64
	OssId   int64
	Version string
}

func (*RcReportOss) TableName() string {
	return "rc_report_oss"
}
