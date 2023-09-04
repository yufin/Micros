package dto

import (
	"fmt"
	"time"
)

type EnterpriseInfo struct {
	InfoId                        int64  `gorm:"primaryKey" json:"-"`
	UscId                         string `gorm:"column:usc_id"`
	EnterpriseTitle               string `gorm:"column:enterprise_title"`
	EnterpriseTitleEn             string `gorm:"column:enterprise_title_en"`
	BusinessRegistrationNumber    string `gorm:"column:business_registration_number"`
	EstablishedDate               *Date  `gorm:"column:established_date"`
	Region                        string `gorm:"column:region"`
	ApprovedDate                  *Date  `gorm:"column:approved_date"`
	RegisteredAddress             string `gorm:"column:registered_address"`
	RegisteredCapital             string `gorm:"column:registered_capital"`
	PaidInCapital                 string `gorm:"column:paid_in_capital"`
	EnterpriseType                string `gorm:"column:enterprise_type"`
	StuffSize                     string `gorm:"column:stuff_size"`
	StuffInsuredNumber            int    `gorm:"column:stuff_insured_number"`
	BusinessScope                 string `gorm:"column:business_scope"`
	ImportExportQualificationCode string `gorm:"column:import_export_qualification_code"`
	LegalRepresentative           string `gorm:"column:legal_representative"`
	RegistrationAuthority         string `gorm:"column:registration_authority"`
	RegistrationStatus            string `gorm:"column:registration_status"`
	TaxpayerQualification         string `gorm:"column:taxpayer_qualification"`
	OrganizationCode              string `gorm:"column:organization_code"`
	UrlQcc                        string `gorm:"column:url_qcc"`
	UrlHomepage                   string `gorm:"column:url_homepage"`
	BusinessTermStart             *Date  `gorm:"column:business_term_start"`
	BusinessTermEnd               *Date  `gorm:"column:business_term_end"`
	StatusCode                    int    `gorm:"column:status_code" json:"-"`
	BaseField
}

func (EnterpriseInfo) TableName() string {
	return "enterprise_info"
}

type EnterpriseCertification struct {
	CertId                 int64  `gorm:"primaryKey" json:"-"`
	UscId                  string `gorm:"column:usc_id"`
	CertificationTitle     string `gorm:"column:certification_title"`
	CertificationCode      string `gorm:"column:certification_code"`
	CertificationLevel     string `gorm:"column:certification_level"`
	CertificationType      string `gorm:"column:certification_type"`
	CertificationSource    string `gorm:"column:certification_source"`
	CertificationDate      *Date  `gorm:"column:certification_date"`
	CertificationTermStart *Date  `gorm:"column:certification_term_start"`
	CertificationTermEnd   *Date  `gorm:"column:certification_term_end"`
	CertificationAuthority string `gorm:"column:certification_authority"`
	StatusCode             int    `gorm:"column:status_code" json:"-"`
	BaseField
}

func (EnterpriseCertification) TableName() string {
	return "enterprise_certification"
}

type Date time.Time

func (dt *Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*dt).Format("2006-01-02"))), nil
}

type EnterpriseRankingList struct {
	UscId                 string `gorm:"column:usc_id"`
	RankingPosition       *int   `gorm:"column:ranking_position"`
	ListTitle             string `gorm:"column:list_title"`
	ListType              string `gorm:"column:list_type"`
	ListSource            string `gorm:"column:list_source"`
	ListParticipantsTotal *int   `gorm:"column:list_participants_total"`
	ListPublishDate       *Date  `gorm:"column:list_published_date"`
	ListUrlQcc            string `gorm:"column:list_url_qcc"`
	ListUrlOrigin         string `gorm:"column:list_url_origin"`
}
