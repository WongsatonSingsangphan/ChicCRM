package models

type LoginResponse struct {
	ID               string //
	TitleToken       string //
	UsernameToken    string //
	FirstnameTokenEN string //
	SurnameTokenEN   string //
	MobilePhoneToken string //
	Requires_action  string //
	CompanyID        string //
	CompanynameToken string //
	JobTitleToken    string //
	DomainToken      string //
	Role             string //
	CompanyLogoToken string //
	RequestID        []string
	// UsernameOriginal           string
	// FirstnameOriginal          string
	// SurnameTokenOriginal       string
	// CompanynameOrginal      string
}

type LoginResponseTeamlead struct {
	TeamleadID        string
	UsernameToken     string
	FirstnameToken    string
	SurnameToken      string
	Mobile_phoneToken string
	Management_Level  int
	Role              string
}
