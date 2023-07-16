package dto

import (
	"brillinkmicros/pkg"
	"time"
)

type BaseModel struct {
	Id        int64 `gorm:"primaryKey"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
	CreateBy  *int64
	UpdateBy  *int64
}

func (e *BaseModel) Gen() {
	if e.Id == 0 {
		e.Id = pkg.Flake.NewFlakeId()
	}
}

type PaginationReq struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type PaginationResp struct {
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
	PageNum   int   `json:"page_num"`
	PageSize  int   `json:"page_size"`
}
