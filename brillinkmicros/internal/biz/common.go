package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/sony/sonyflake"
	"time"
)

var f *sonyflake.Sonyflake = nil

type Snowflake struct {
	Sf *sonyflake.Sonyflake
}

func (s *Snowflake) GetFlake() *sonyflake.Sonyflake {
	if s.Sf == nil {
		st := sonyflake.Settings{}
		s.Sf = sonyflake.NewSonyflake(st)
		if s.Sf == nil {
			log.Error("sonyflake not created \r\n")
			panic("sonyflake not created")
		}
	}
	return s.Sf
}

func (s *Snowflake) NewFlakeId() int64 {
	id, err := s.GetFlake().NextID()
	if err != nil {
		log.Errorf("sonyflake nextId error:%s \r\n", err)
		panic("sonyflake not created")
	}
	return int64(id)
}

var flake = Snowflake{Sf: f}

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
		e.Id = flake.NewFlakeId()
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
