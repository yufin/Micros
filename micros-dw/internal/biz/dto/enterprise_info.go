package dto

import (
	"fmt"
	"time"
)

type Date time.Time

func (dt *Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*dt).Format("2006-01-02"))), nil
}

type EnterpriseInfo struct {
	InfoId                        int64     `gorm:"primaryKey" json:"-"`
	UscId                         string    `gorm:"column:usc_id" json:"uscId"`
	EnterpriseTitle               string    `gorm:"column:enterprise_title" json:"enterpriseTitle"`
	EnterpriseTitleEn             string    `gorm:"column:enterprise_title_en" json:"enterpriseTitleEn"`
	BusinessRegistrationNumber    string    `gorm:"column:business_registration_number" json:"businessRegistrationNumber"`
	EstablishedDate               time.Time `gorm:"column:established_date" json:"establishedDate"`
	Region                        string    `gorm:"column:region" json:"region"`
	ApprovedDate                  time.Time `gorm:"column:approved_date" json:"approvedDate"`
	RegisteredAddress             string    `gorm:"column:registered_address" json:"registeredAddress"`
	RegisteredCapital             string    `gorm:"column:registered_capital" json:"registeredCapital"`
	PaidInCapital                 string    `gorm:"column:paid_in_capital" json:"paidInCapital"`
	EnterpriseType                string    `gorm:"column:enterprise_type" json:"enterpriseType"`
	StuffSize                     string    `gorm:"column:stuff_size" json:"stuffSize"`
	StuffInsuredNumber            int       `gorm:"column:stuff_insured_number" json:"stuffInsuredNumber"`
	BusinessScope                 string    `gorm:"column:business_scope" json:"businessScope"`
	ImportExportQualificationCode string    `gorm:"column:import_export_qualification_code" json:"importExportQualificationCode"`
	LegalRepresentative           string    `gorm:"column:legal_representative" json:"legalRepresentative"`
	RegistrationAuthority         string    `gorm:"column:registration_authority" json:"registrationAuthority"`
	RegistrationStatus            string    `gorm:"column:registration_status" json:"registrationStatus"`
	TaxpayerQualification         string    `gorm:"column:taxpayer_qualification" json:"taxpayerQualification"`
	OrganizationCode              string    `gorm:"column:organization_code" json:"organizationCode"`
	UrlQcc                        string    `gorm:"column:url_qcc" json:"urlQcc"`
	UrlHomepage                   string    `gorm:"column:url_homepage" json:"urlHomepage"`
	BusinessTermStart             time.Time `gorm:"column:business_term_start" json:"businessTermStart"`
	BusinessTermEnd               time.Time `gorm:"column:business_term_end" json:"businessTermEnd"`
	StatusCode                    int       `gorm:"column:status_code" json:"-"`
	BaseField
}

func (EnterpriseInfo) TableName() string {
	return "enterprise_info"
}
