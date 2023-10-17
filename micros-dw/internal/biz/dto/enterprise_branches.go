package dto

import "time"

type EnterpriseBranches struct {
	EnterpriseName string    `bson:"enterprise_name"`
	Operator       string    `bson:"operator"`
	Area           string    `bson:"area"`
	StartDate      time.Time `bson:"start_date"`
	Status         string    `bson:"status"`
}
