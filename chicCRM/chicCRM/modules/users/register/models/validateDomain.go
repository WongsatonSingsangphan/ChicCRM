package models

type ValidateDomainRequest struct {
	Username string `json:"username"`
}

type ValidateDomainResponse struct {
	CompanyID    string
	CompanyLogo  string
	Companyname  string
	Domain       string
	Alpha2Code   string
	LanguageName string
	CompanyAlias string
	Currency     string
	AddressNo    string
	Address1En   string
	Geolocation  string
	CompanyPhone string
	Zipcode      string
	Province     string
	District     string
	SubDistrict  string
	Timezone     string
	Country      string
}
