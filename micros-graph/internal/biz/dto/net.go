package dto

type Net struct {
	Nodes []*Node `json:"nodes"`
	Edges []*Edge `json:"edges"`
}
