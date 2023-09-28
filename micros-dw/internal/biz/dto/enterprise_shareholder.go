package dto

type T struct {
	KeyNo             string      `json:"KeyNo"`
	IsCrossHolding    bool        `json:"IsCrossHolding"`
	CompanyCount      int         `json:"CompanyCount"`
	StockName         string      `json:"StockName"`
	EciStockName      interface{} `json:"EciStockName"`
	StockType         string      `json:"StockType"`
	StockPercent      string      `json:"StockPercent"`
	IdentifyType      string      `json:"IdentifyType"`
	IdentifyNo        string      `json:"IdentifyNo"`
	ShoudDate         string      `json:"ShoudDate"`
	InvestType        interface{} `json:"InvestType"`
	InvestName        string      `json:"InvestName"`
	CapiDate          string      `json:"CapiDate"`
	RegistCapi        string      `json:"RegistCapi"`
	Unit              string      `json:"Unit"`
	HasSecLvlFlag     bool        `json:"HasSecLvlFlag"`
	Org               int         `json:"Org"`
	Job               string      `json:"Job"`
	ImageUrl          string      `json:"ImageUrl"`
	StockPercentValue int         `json:"StockPercentValue"`
	Tags              []string    `json:"Tags"`
	TagConfig         struct {
		BPL  string `json:"BPL"`
		IsBP string `json:"IsBP"`
	} `json:"TagConfig"`
	FinalBenefitPercent     string        `json:"FinalBenefitPercent"`
	HasStockDetail          bool          `json:"HasStockDetail"`
	HasStockFreezenInfo     bool          `json:"HasStockFreezenInfo"`
	HasPledgeInfo           bool          `json:"HasPledgeInfo"`
	OriginalNames           []interface{} `json:"OriginalNames"`
	PartnerShouldDetailList []struct {
		InvestType interface{} `json:"InvestType"`
		ShoudDate  string      `json:"ShoudDate"`
		ShouldCapi string      `json:"ShouldCapi"`
	} `json:"PartnerShouldDetailList"`
	Source                interface{} `json:"Source"`
	TotalShouldAmount     string      `json:"TotalShouldAmount"`
	TotalRealAmount       string      `json:"TotalRealAmount"`
	PartnerRealDetailList []struct {
		RealCapi   string `json:"RealCapi"`
		CapiDate   string `json:"CapiDate"`
		InvestName string `json:"InvestName"`
	} `json:"PartnerRealDetailList"`
	Product struct {
		Id             string `json:"Id"`
		Name           string `json:"Name"`
		Round          string `json:"Round"`
		RoundDesc      string `json:"RoundDesc"`
		Amount         string `json:"Amount"`
		FinancingCount int    `json:"FinancingCount"`
		CompatCount    int    `json:"CompatCount"`
		MemberCount    int    `json:"MemberCount"`
		Logo           string `json:"Logo"`
	} `json:"Product"`
	InDate         interface{} `json:"InDate"`
	PublicDate     interface{} `json:"PublicDate"`
	Area           interface{} `json:"Area"`
	ShareType      interface{} `json:"ShareType"`
	UnitType       interface{} `json:"UnitType"`
	NatCode        interface{} `json:"NatCode"`
	IsReportMerged interface{} `json:"IsReportMerged"`
	Invest         struct {
		Id   string `json:"Id"`
		Name string `json:"Name"`
		Logo string `json:"Logo"`
		NC   int    `json:"NC"`
		FC   int    `json:"FC"`
	} `json:"Invest"`
	TagsInfo []struct {
		Type             int         `json:"Type"`
		Name             string      `json:"Name"`
		ShortName        string      `json:"ShortName"`
		DataExtend       string      `json:"DataExtend"`
		DataExtend2      string      `json:"DataExtend2"`
		TradingPlaceCode interface{} `json:"TradingPlaceCode"`
		TradingPlaceName interface{} `json:"TradingPlaceName"`
	} `json:"TagsInfo"`
	YearOfRealCapi string      `json:"YearOfRealCapi"`
	RoundInfo      interface{} `json:"RoundInfo"`
}
