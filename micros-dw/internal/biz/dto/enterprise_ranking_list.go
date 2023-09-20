package dto

import "time"

type EnterpriseRankingList struct {
	UscId                 string    `gorm:"column:usc_id" json:"uscId"`
	RankingPosition       int       `gorm:"column:ranking_position" json:"rankingPosition"`
	ListTitle             string    `gorm:"column:list_title" json:"listTitle"`
	ListType              string    `gorm:"column:list_type" json:"listType"`
	ListSource            string    `gorm:"column:list_source" json:"listSource"`
	ListParticipantsTotal int       `gorm:"column:list_participants_total" json:"listParticipantsTotal"`
	ListPublishedDate     time.Time `gorm:"column:list_published_date" json:"listPublishDate"`
	ListUrlQcc            string    `gorm:"column:list_url_qcc" json:"listUrlQcc"`
	ListUrlOrigin         string    `gorm:"column:list_url_origin" json:"listUrlOrigin"`
}
