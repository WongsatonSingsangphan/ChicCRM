package services

import (
	"chicCRM/modules/queryData/models"
	"chicCRM/modules/queryData/repositories"
)

// type ServicePort interface {
// 	GetCountriesServices() ([]getQueryModels.CountryResponse, error)
// }

type ServicePort interface {
	GetCountriesServices() (*models.CountryResponse, error)
	GetProvinceAmphoeTambonZipcodeServices() (*models.TambonDataResponse, error)
	GetUniversalInfoServices() (*models.UniversalInfoResponse, error)
}

type serviceAdapter struct {
	r repositories.RepositoryPort
}

func NewServiceAdapter(r repositories.RepositoryPort) ServicePort {
	return &serviceAdapter{r: r}
}

// func (s serviceAdapter) GetCountriesServices() ([]getQueryModels.CountryResponse, error) { ///////// return slice in case ,easy to iterate and return  list of data , flexible.

// 	countries, err := s.r.GetCountries()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// ไม่จำเป็นต้องสร้าง slice ของ CountryResponse ใหม่ ใช้ CountryResponse ที่มี field เป็น slice ของ Country
// 	countriesResponse := getQueryModels.CountryResponse{
// 		Countries: countries, // เพิ่ม countries slice ไปยัง field Countries
// 	}

// 	return []getQueryModels.CountryResponse{countriesResponse}, nil // รีเทิร์นเป็น slice ของ CountryResponse

// }

func (s *serviceAdapter) GetCountriesServices() (*models.CountryResponse, error) { /////// return single obeject complex and not flexible but very fast , the value pass references not copy into the funciton.

	countries, err := s.r.GetCountriesRepositories()
	if err != nil {
		return nil, err
	}
	countriesResponse := models.CountryResponse{
		Countries: countries,
	}
	return &countriesResponse, nil

}

func (s *serviceAdapter) GetProvinceAmphoeTambonZipcodeServices() (*models.TambonDataResponse, error) {
	tambonData, err := s.r.GetProvinceAmphoeTambonZipcodeRepositories()
	if err != nil {
		return nil, err
	}
	tambonDataResponse := models.TambonDataResponse{
		ProvinceAmphoeTambonZipcode: tambonData,
	}
	return &tambonDataResponse, nil
}
func (s *serviceAdapter) GetUniversalInfoServices() (*models.UniversalInfoResponse, error) {
	sexes, religions, degrees, maritalstatuses, militarystatuses, err := s.r.GetUniversalInfoRepositories()
	if err != nil {
		return nil, err
	}
	universalInfoResponse := models.UniversalInfoResponse{
		Sexes:            sexes,
		Religions:        religions,
		Degrees:          degrees,
		Maritalstatuses:  maritalstatuses,
		Militarystatuses: militarystatuses,
	}
	return &universalInfoResponse, nil
}
