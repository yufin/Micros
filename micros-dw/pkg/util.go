package pkg

import "github.com/google/uuid"

func RandUuid() string {
	randId, _ := uuid.NewRandom()
	return randId.String()
}
