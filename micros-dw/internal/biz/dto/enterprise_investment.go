package dto

type EnterpriseInvestment struct {
	EnterpriseName    string `bson:"enterprise_name"`
	Operator          string `bson:"operator"`
	ShareholdingRatio string `bson:"shareholding_ratio"`
	InvestedAmount    string `bson:"capital_amount"`
	StartDate         string `bson:"start_date"`
	Status            string `bson:"status"`
}
