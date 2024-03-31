package dto

import "time"

type ReportPrintConfig struct {
	CreateAt time.Time      `bson:"create_at" json:"create_at"`
	CreateBy int64          `bson:"create_by" json:"create_by"`
	Config   map[string]any `bson:"config" json:"config"`
}
