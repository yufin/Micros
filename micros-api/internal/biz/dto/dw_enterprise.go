package dto

import (
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"
	"time"
)

type Date time.Time

func (dt *Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*dt).Format("2006-01-02"))), nil
}

type EnterpriseInfo struct {
	UscId                         string `json:"uscId"`
	EnterpriseTitle               string `json:"enterpriseTitle"`
	EnterpriseTitleEn             string `json:"enterpriseTitleEn"`
	BusinessRegistrationNumber    string `json:"businessRegistrationNumber"`
	EstablishedDate               string `json:"establishedDate"`
	Region                        string `json:"region"`
	ApprovedDate                  string `json:"approvedDate"`
	RegisteredAddress             string `json:"registeredAddress"`
	RegisteredCapital             string `json:"registeredCapital"`
	PaidInCapital                 string `json:"paidInCapital"`
	EnterpriseType                string `json:"enterpriseType"`
	StuffSize                     string `json:"stuffSize"`
	StuffInsuredNumber            int    `json:"stuffInsuredNumber"`
	BusinessScope                 string `json:"businessScope"`
	ImportExportQualificationCode string `json:"importExportQualificationCode"`
	LegalRepresentative           string `json:"legalRepresentative"`
	RegistrationAuthority         string `json:"registrationAuthority"`
	RegistrationStatus            string `json:"registrationStatus"`
	TaxpayerQualification         string `json:"taxpayerQualification"`
	OrganizationCode              string `json:"organizationCode"`
	UrlQcc                        string `json:"urlQcc"`
	UrlHomepage                   string `json:"urlHomepage"`
	BusinessTermStart             string `json:"businessTermStart"`
	BusinessTermEnd               string `json:"businessTermEnd"`
}

type EnterpriseCertification struct {
	UscId                  string `json:"uscId"`
	CertificationTitle     string `json:"certificationTitle"`
	CertificationCode      string `json:"certificationCode"`
	CertificationLevel     string `json:"certificationLevel"`
	CertificationType      string `json:"certificationType"`
	CertificationSource    string `json:"certificationSource"`
	CertificationDate      string `json:"certificationDate"`
	CertificationTermStart string `json:"certificationTermStart"`
	CertificationTermEnd   string `json:"certificationTermEnd"`
	CertificationAuthority string `json:"certificationAuthority"`
}

type EnterpriseRankingList struct {
	UscId                 string `json:"uscId"`
	RankingPosition       int    `json:"rankingPosition"`
	ListTitle             string `json:"listTitle"`
	ListType              string `json:"listType"`
	ListSource            string `json:"listSource"`
	ListParticipantsTotal int    `json:"listParticipantsTotal"`
	ListPublishedDate     string `json:"listPublishDate"`
	ListUrlQcc            string `json:"listUrlQcc"`
	ListUrlOrigin         string `json:"listUrlOrigin"`
}

type EnterpriseEquityTransparency struct {
	UscId      string
	Conclusion string
	Name       string
	Data       []*structpb.Struct
}
