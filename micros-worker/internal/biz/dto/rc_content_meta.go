package dto

type RcContentMeta struct {
	BaseModel
	UscId          string `json:"uscId" gorm:"comment:统一社会信用代码"`
	EnterpriseName string `json:"enterpriseName" gorm:"comment:企业名称"`
	AttributeMonth string `json:"attributeMonth" gorm:"comment:数据更新年月"`
	SourcePath     string `json:"SourcePath" gorm:"comment:sftp路径"`
	Version        string `json:"version" gorm:"comment:报文版本"`
}

func (RcContentMeta) TableName() string {
	return "rc_content_meta"
}
