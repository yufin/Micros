package dto

import (
	"time"
)

type EnterpriseCommentInsertReq struct {
	EnterpriseName string    `json:"enterprise_name" bson:"enterprise_name"`
	UscId          string    `json:"usc_id" bson:"usc_id"`
	Comment        string    `json:"comment" bson:"comment"`
	Prestige       string    `json:"prestige" bson:"prestige"`
	CreateAt       time.Time `bson:"create_at"`
	Ban            bool      `json:"ban" bson:"ban"`
}

type EnterpriseComment struct {
	Id             string    `json:"id" bson:"_id"`
	EnterpriseName string    `json:"enterprise_name" bson:"enterprise_name"`
	UscId          string    `json:"usc_id" bson:"usc_id"`
	Comment        string    `json:"comment" bson:"comment"`
	Prestige       string    `json:"prestige" bson:"prestige"`
	CreateAt       time.Time `bson:"create_at"`
	Ban            bool      `json:"ban" bson:"ban"`
}

type ProductEvalRuleInsertResp struct {
	Product         string    `json:"product" bson:"product"`
	IndustryComment string    `json:"industry_comment" bson:"industry_comment"`
	RuleSubstring   []string  `json:"rule_substring" bson:"rule_substring"`
	CreateAt        time.Time `bson:"create_at"`
	Ban             bool      `json:"ban" bson:"ban"`
}

type ProductEvalRule struct {
	Id              string    `json:"id" bson:"_id"`
	Product         string    `json:"product" bson:"product"`
	IndustryComment string    `json:"industry_comment" bson:"industry_comment"`
	RuleSubstring   []string  `json:"rule_substring" bson:"rule_substring"`
	CreateAt        time.Time `bson:"create_at"`
	Ban             bool      `json:"ban" bson:"ban"`
}
