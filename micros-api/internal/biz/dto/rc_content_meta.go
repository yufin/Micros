package dto

//create table rc_content_meta
//(
//id              bigint       not null
//primary key,
//usc_id          char(18)     not null comment '统一社会信用代码',
//enterprise_name varchar(255) null comment '企业名称',
//attribute_month varchar(10)  null comment '数据所属年月(eg:2020-12)',
//source_path     varchar(512) null,
//version         varchar(10)  null,
//created_at      datetime(3)  null,
//updated_at      datetime(3)  null,
//deleted_at      datetime(3)  null,
//create_by       bigint       null,
//update_by       bigint       null
//);

type RcContentMeta struct {
	BaseModel
	UscId          string `json:"usc_id" gorm:"column:usc_id" xlsx:"usc_id"`
	EnterpriseName string `json:"enterprise_name" gorm:"column:enterprise_name" xlsx:"enterprise_name"`
	AttributeMonth string `json:"attribute_month" gorm:"column:attribute_month" xlsx:"attribute_month"`
	SourcePath     string `json:"source_path" gorm:"column:source_path" xlsx:"source_path"`
	Version        string `json:"version" gorm:"column:version" xlsx:"version"`
}

func (*RcContentMeta) TableName() string {
	return "rc_content_meta"
}
