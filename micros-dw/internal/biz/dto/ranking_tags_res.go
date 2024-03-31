package dto

import "time"

type RankingTagsRes struct {
	Data []struct {
		RankingTags struct {
			Id              string `json:"Id"`
			InstitutionName string `json:"InstitutionName"`
			Url             string `json:"Url"`
			PublishDate     string `json:"PublishDate"`
			Title           string `json:"Title"`
			Ranking         int    `json:"Ranking"`
		} `json:"ranking_tags"`
		ListDetail struct {
			Id              string `json:"Id"`
			Title           string `json:"Title"`
			InstitutionName string `json:"InstitutionName"`
			PublishDate     struct {
				Date time.Time `json:"$date"`
			} `json:"PublishDate"`
			RegionDesc string `json:"RegionDesc"`
			Url        string `json:"Url"`
			BdType     string `json:"BdType"`
			CreateTime struct {
				Date time.Time `json:"$date"`
			} `json:"create_time"`
		} `json:"list_detail"`
	} `json:"data"`
	Total int `json:"total"`
}
