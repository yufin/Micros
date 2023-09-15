package dto

import "time"

type EnterpriseCertification struct {
	CertId                 int64     `gorm:"primaryKey" json:"-"`
	UscId                  string    `gorm:"column:usc_id" json:"uscId"`
	CertificationTitle     string    `gorm:"column:certification_title" json:"certificationTitle"`
	CertificationCode      string    `gorm:"column:certification_code" json:"certificationCode"`
	CertificationLevel     string    `gorm:"column:certification_level" json:"certificationLevel"`
	CertificationType      string    `gorm:"column:certification_type" json:"certificationType"`
	CertificationSource    string    `gorm:"column:certification_source" json:"certificationSource"`
	CertificationDate      time.Time `gorm:"column:certification_date" json:"certificationDate"`
	CertificationTermStart time.Time `gorm:"column:certification_term_start" json:"certificationTermStart"`
	CertificationTermEnd   time.Time `gorm:"column:certification_term_end" json:"certificationTermEnd"`
	CertificationAuthority string    `gorm:"column:certification_authority" json:"certificationAuthority"`
	StatusCode             int       `gorm:"column:status_code" json:"-"`
	BaseField
}

func (EnterpriseCertification) TableName() string {
	return "enterprise_certification"
}
