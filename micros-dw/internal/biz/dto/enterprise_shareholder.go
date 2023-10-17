package dto

type EnterpriseShareholder struct {
	ShareholderName string `bson:"shareholder_name"`
	ShareholderType string `bson:"shareholder_type"`
	CapitalAmount   string `bson:"capital_amount"`
	RealAmount      string `bson:"real_amount"`
	CapitalType     string `bson:"capital_type"`
	Percent         string `bson:"percent"`
}
