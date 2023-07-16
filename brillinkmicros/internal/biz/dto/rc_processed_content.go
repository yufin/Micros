package dto

type RcProcessedContent struct {
	BaseModel
	ContentId int64
	Content   string
}

func (*RcProcessedContent) TableName() string {
	return "rskc_processed_content"
}
