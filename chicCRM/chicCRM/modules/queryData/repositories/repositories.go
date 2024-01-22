package repositories

import (
	"chicCRM/modules/queryData/models"
	"database/sql"
	"log"
)

type RepositoryPort interface {
	GetCountriesRepositories() ([]models.Country, error)
	GetProvinceAmphoeTambonZipcodeRepositories() ([]models.TambonData, error)
	GetUniversalInfoRepositories() ([]models.Sex, []models.Religion, []models.Degree, []models.Maritalstatus, []models.Militarystatus, error)
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}

func (r *repositoryAdapter) GetCountriesRepositories() ([]models.Country, error) {
	var countries []models.Country
	rows, err := r.db.Query("SELECT * FROM countries")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var country models.Country
		err := rows.Scan(&country.ID, &country.Name, &country.TopLevelDomain, &country.Alpha2Code, &country.Alpha3Code, &country.CallingCodes, &country.Capital, &country.Subregion, &country.Region, &country.Population, &country.Latitude, &country.Longitude, &country.Demonym, &country.Area, &country.Gini, &country.Timezones, &country.NativeName, &country.NumericCode, &country.FlagSVG, &country.FlagPNG, &country.CurrencyCode, &country.CurrencyName, &country.CurrencySymbol, &country.LanguageCode, &country.LanguageName, &country.LanguageNativeName, &country.IsIndependent)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func (r *repositoryAdapter) GetProvinceAmphoeTambonZipcodeRepositories() ([]models.TambonData, error) {
	var tambonData []models.TambonData

	rows, err := r.db.Query("SELECT * FROM \"thailandTambon\"") // even uppercase also need \"\"
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data models.TambonData
		if err := rows.Scan(
			&data.ID, &data.TambonID, &data.TambonThai, &data.TambonEng, &data.TambonThaiShort, &data.TambonEngShort,
			&data.DistrictID, &data.DistrictThai, &data.DistrictEng, &data.DistrictThaiShort, &data.DistrictEngShort,
			&data.ProvinceID, &data.ProvinceThai, &data.ProvinceEng, &data.OfficialRegions, &data.FourMainRegions,
			&data.TouristRegions, &data.GreaterBangkokProvinces, &data.PostalCodeRemark, &data.PostCodeMain, &data.PostCodeAll,
		); err != nil {
			log.Println(err)
			return nil, err
		}
		tambonData = append(tambonData, data)
	}
	return tambonData, nil
}

func (r *repositoryAdapter) GetUniversalInfoRepositories() ([]models.Sex, []models.Religion, []models.Degree, []models.Maritalstatus, []models.Militarystatus, error) {
	sexRows, err := r.db.Query("SELECT * FROM genders")
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, nil, err
	}
	defer sexRows.Close()
	var sexes []models.Sex

	for sexRows.Next() {
		var sex models.Sex
		err := sexRows.Scan(&sex.ID, &sex.Sex_th, &sex.Sex_en, &sex.S_th, &sex.S_en)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, nil, nil, err
		}
		sexes = append(sexes, sex)
	}

	religionRows, err := r.db.Query("SELECT * FROM religions")
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, nil, err
	}
	defer religionRows.Close()
	var religions []models.Religion

	for religionRows.Next() {
		var religion models.Religion
		err := religionRows.Scan(&religion.ID, &religion.Religion_th, &religion.Religion_en)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, nil, nil, err
		}
		religions = append(religions, religion)
	}
	degreeRows, err := r.db.Query("SELECT * FROM \"levelDegrees\"")
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, nil, err
	}
	defer degreeRows.Close()
	var degrees []models.Degree

	for degreeRows.Next() {
		var degree models.Degree
		err := degreeRows.Scan(&degree.ID, &degree.Degree_th, &degree.Degree_en)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, nil, nil, err
		}
		degrees = append(degrees, degree)
	}

	maritalstatusRows, err := r.db.Query("SELECT * FROM \"maritalStatuses\"")
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, nil, err
	}
	defer maritalstatusRows.Close()
	var maritalstatuses []models.Maritalstatus

	for maritalstatusRows.Next() {
		var maritalstatus models.Maritalstatus
		err := maritalstatusRows.Scan(&maritalstatus.ID, &maritalstatus.Maritalstatus, &maritalstatus.Maritalstatus_en)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, nil, nil, err
		}
		maritalstatuses = append(maritalstatuses, maritalstatus)
	}

	militarystatusRows, err := r.db.Query("SELECT * FROM \"militaryStatuses\"")
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, nil, err
	}
	defer militarystatusRows.Close()
	var militarystatuses []models.Militarystatus

	for militarystatusRows.Next() {
		var militarystatus models.Militarystatus
		err := militarystatusRows.Scan(&militarystatus.ID, &militarystatus.Militarystatus, &militarystatus.Militarystatus_en)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, nil, nil, err
		}
		militarystatuses = append(militarystatuses, militarystatus)
	}
	return sexes, religions, degrees, maritalstatuses, militarystatuses, nil
}
