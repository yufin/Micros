package dto

type EnterpriseShareholder struct {
	KeyNo         string      `json:"KeyNo"`
	Name          string      `json:"Name"`
	No            string      `json:"No"`
	CreditCode    string      `json:"CreditCode"`
	EconKind      string      `json:"EconKind"`
	Status        string      `json:"Status"`
	RegistCapi    string      `json:"RegistCapi"`
	FundedRatio   string      `json:"FundedRatio"`
	InvestInDate  int         `json:"InvestInDate"`
	InvestOutDate interface{} `json:"InvestOutDate"`
	StartDate     string      `json:"StartDate"`
	Oper          struct {
		HasImage     bool   `json:"HasImage"`
		CompanyCount int    `json:"CompanyCount"`
		Org          int    `json:"Org"`
		KeyNo        string `json:"KeyNo"`
		Name         string `json:"Name"`
		OperType     int    `json:"OperType"`
	} `json:"Oper"`
	MultipleOper struct {
		OperType int `json:"OperType"`
		OperList []struct {
			Assignor     interface{} `json:"Assignor"`
			HasImage     bool        `json:"HasImage"`
			CompanyCount int         `json:"CompanyCount"`
			Org          int         `json:"Org"`
			KeyNo        string      `json:"KeyNo"`
			Name         string      `json:"Name"`
		} `json:"OperList"`
		OperTypeDesc string `json:"OperTypeDesc"`
	} `json:"MultipleOper"`
	ImageUrl     string `json:"ImageUrl"`
	IndustryItem struct {
		IndustryCode       string `json:"IndustryCode"`
		Industry           string `json:"Industry"`
		SubIndustryCode    string `json:"SubIndustryCode"`
		SubIndustry        string `json:"SubIndustry"`
		MiddleCategoryCode string `json:"MiddleCategoryCode"`
		MiddleCategory     string `json:"MiddleCategory"`
		SmallCategoryCode  string `json:"SmallCategoryCode"`
		SmallCategory      string `json:"SmallCategory"`
	} `json:"IndustryItem"`
	Product  interface{} `json:"Product"`
	Invest   interface{} `json:"Invest"`
	TagsInfo []struct {
		Type             int         `json:"Type"`
		Name             string      `json:"Name"`
		ShortName        interface{} `json:"ShortName"`
		DataExtend       string      `json:"DataExtend"`
		DataExtend2      interface{} `json:"DataExtend2"`
		TradingPlaceCode interface{} `json:"TradingPlaceCode"`
		TradingPlaceName interface{} `json:"TradingPlaceName"`
	} `json:"TagsInfo"`
	SecLvlBP struct {
		KeyNo        string `json:"KeyNo"`
		Name         string `json:"Name"`
		Org          int    `json:"Org"`
		StockPercent string `json:"StockPercent"`
		ShouldCapi   string `json:"ShouldCapi"`
	} `json:"SecLvlBP"`
	InvesterCount int    `json:"InvesterCount"`
	RealDate      int    `json:"RealDate"`
	RealCapi      string `json:"RealCapi"`
	ShouldDate    int    `json:"ShouldDate"`
	ShouldCapi    string `json:"ShouldCapi"`
	StockType     string `json:"StockType"`
	TotalPercent  string `json:"TotalPercent"`
	Area          struct {
		ProvinceCode string `json:"ProvinceCode"`
		ProvinceName string `json:"ProvinceName"`
		CityCode     int    `json:"CityCode"`
		CityName     string `json:"CityName"`
		CountyCode   int    `json:"CountyCode"`
		CountyName   string `json:"CountyName"`
	} `json:"Area"`
	OperName         string `json:"OperName"`
	Province         string `json:"Province"`
	ProvinceName     string `json:"ProvinceName"`
	MultiShouldDate  bool   `json:"MultiShouldDate"`
	MultiRealDate    bool   `json:"MultiRealDate"`
	IsMultiRealCapi  bool   `json:"isMultiRealCapi"`
	IsMultiShoulCapi bool   `json:"isMultiShoulCapi"`
}
