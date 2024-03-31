package dto

import (
	"gorm.io/gorm"
	"micros-worker/pkg"
	"time"
)

type BaseModel struct {
	Id        int64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	CreateBy  *int64
	UpdateBy  *int64
}

type BaseField struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	CreateBy  *int64         `json:"-"`
	UpdateBy  *int64         `json:"-"`
}

func (e *BaseModel) Gen() {
	if e.Id == 0 {
		e.Id = pkg.Flake.NewFlakeId()
	}
}

type PaginationReq struct {
	PageNum  int
	PageSize int
}

type PaginationResp struct {
	Total    int64
	PageNum  int
	PageSize int
}
