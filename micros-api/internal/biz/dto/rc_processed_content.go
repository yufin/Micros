package dto

type RcProcessedContent struct {
	BaseModel
	ContentId int64
	Content   string
}

func (*RcProcessedContent) TableName() string {
	return "rc_processed_content"
}
