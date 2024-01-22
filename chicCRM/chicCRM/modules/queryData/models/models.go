package models

type Sex struct {
	ID     int    `json:"id"`
	Sex_th string `json:"sex_th"`
	Sex_en string `json:"sex_en"`
	S_th   string `json:"s_th"`
	S_en   string `json:"s_en"`
}

type Religion struct {
	ID          int    `json:"id"`
	Religion_th string `json:"religion_th"`
	Religion_en string `json:"religion_en"`
}

type Degree struct {
	ID        int    `json:"id"`
	Degree_th string `json:"degree_th"`
	Degree_en string `json:"degree_en"`
}

type Maritalstatus struct {
	ID               int    `json:"maritalstatus_id"`
	Maritalstatus    string `json:"maritalstatus"`
	Maritalstatus_en string `json:"martalstatus_en"`
}

type Militarystatus struct {
	ID                int    `json:"militarystatus_id"`
	Militarystatus    string `json:"militarystatus"`
	Militarystatus_en string `json:"milirarystatus_en"`
}

type UniversalInfoResponse struct {
	Sexes            []Sex            `json:"sexes"`
	Religions        []Religion       `json:"religions"`
	Degrees          []Degree         `json:"degrees"`
	Maritalstatuses  []Maritalstatus  `json:"maritalstatuses"`
	Militarystatuses []Militarystatus `json:"militarystatuses"`
}

type Country struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	TopLevelDomain     string  `json:"topLevelDomain"`
	Alpha2Code         string  `json:"alpha2Code"`
	Alpha3Code         string  `json:"alpha3Code"`
	CallingCodes       string  `json:"callingCodes"`
	Capital            string  `json:"capital"`
	Subregion          string  `json:"subregion"`
	Region             string  `json:"region"`
	Population         int     `json:"population"`
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Demonym            string  `json:"demonym"`
	Area               float64 `json:"area"`
	Gini               float64 `json:"gini"`
	Timezones          string  `json:"timezones"`
	NativeName         string  `json:"nativeName"`
	NumericCode        string  `json:"numericCode"`
	FlagSVG            string  `json:"flag_svg"`
	FlagPNG            string  `json:"flag_png"`
	CurrencyCode       string  `json:"currency_code"`
	CurrencyName       string  `json:"currency_name"`
	CurrencySymbol     string  `json:"currency_symbol"`
	LanguageCode       string  `json:"language_code"`
	LanguageName       string  `json:"language_name"`
	LanguageNativeName string  `json:"language_nativeName"`
	IsIndependent      bool    `json:"isIndependent"`
}
type CountryResponse struct {
	Countries []Country `json:"countries"`
}

type TambonData struct {
	ID                      int    `json:"id"`
	TambonID                string `json:"TambonID"`
	TambonThai              string `json:"TambonThai"`
	TambonEng               string `json:"TambonEng"`
	TambonThaiShort         string `json:"TambonThaiShort"`
	TambonEngShort          string `json:"TambonEngShort"`
	DistrictID              string `json:"DistrictID"`
	DistrictThai            string `json:"DistrictThai"`
	DistrictEng             string `json:"DistrictEng"`
	DistrictThaiShort       string `json:"DistrictThaiShort"`
	DistrictEngShort        string `json:"DistrictEngShort"`
	ProvinceID              string `json:"ProvinceID"`
	ProvinceThai            string `json:"ProvinceThai"`
	ProvinceEng             string `json:"ProvinceEng"`
	OfficialRegions         string `json:"official_regions"`
	FourMainRegions         string `json:"four_main_regions"`
	TouristRegions          string `json:"tourist_regions"`
	GreaterBangkokProvinces string `json:"greater_bangkok_provinces"`
	PostalCodeRemark        string `json:"PostalCodeRemark"`
	PostCodeMain            string `json:"PostCodeMain"`
	PostCodeAll             string `json:"PostCodeAll"`
}
type TambonDataResponse struct {
	ProvinceAmphoeTambonZipcode []TambonData `json:"provinceAmphoeTambonZipcode"`
}
