package pkg

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/sony/sonyflake"
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

var Flake = Snowflake{Sf: f}
