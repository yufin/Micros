package dto

type OssMetaData struct {
	BaseModel
	ObjName    string
	BucketName string
	EndPoint   string
	App        int
}

func (*OssMetaData) TableName() string {
	return "oss_metadata"
}
