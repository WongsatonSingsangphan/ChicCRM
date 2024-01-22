package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtResponse struct {
	ID                   string
	TitleOriginal        string
	UsernameToken        string
	UsernameOriginal     string
	FirstnameTokenEN     string
	FirstnameOriginal    string
	SurnameTokenEN       string
	SurnameTokenOriginal string
	MobilePhoneOriginal  string
	Requires_action      string
	CompanyID            string
	CompanynameToken     string
	CompanynameOrginal   string
	CompanyLogoOriginal  string
	JobTitleToken        string
	JobTitleOriginal     string
	DomainToken          string
	DomainOriginal       string
	Role                 string
	Exp                  int64
	// RequestID            []string
}

type JwtResponseTeamleadSecuredog struct {
	TeamleadID        string
	Username          string
	UsernameToken     string
	Firstname         string
	FirstnameToken    string
	Surname           string
	SurnameToken      string
	Mobile_phone      string
	Mobile_phoneToken string
	Role              string
	Management_Level  int
	Exp               int64
}

func (c JwtResponse) Valid() error {
	if c.Exp < time.Now().Unix() {
		return jwt.ErrTokenExpired
	}
	// คุณสามารถเพิ่มการตรวจสอบเงื่อนไขอื่นๆ ที่จำเป็น
	return nil
}

func (c JwtResponseTeamleadSecuredog) Valid() error {
	if c.Exp < time.Now().Unix() {
		return jwt.ErrTokenExpired
	}
	return nil
}
