package services

import (
	"chicCRM/modules/users/register/models"
	"chicCRM/modules/users/register/repositories"
	"chicCRM/pkg/auth"
	"chicCRM/pkg/utilts/decrypt"
	"chicCRM/pkg/utilts/encrypt"
	"chicCRM/pkg/utilts/generate"
	"chicCRM/pkg/utilts/sendEmailFunctions"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type ServicePort interface {
	// RegisterChicCRMServices(loginData models.RegisterRequest) error
	RegisterChicCRMServices(loginData models.RegisterRequest) (models.RegisterResponses, error)
	ValidateDomainChicCRMServices(validateDomainRequest models.ValidateDomainRequest) (models.ValidateDomainResponse, error)
}

type serviceAdapter struct {
	r repositories.RepositoryPort
}

func NewServiceAdapter(r repositories.RepositoryPort) ServicePort {
	return &serviceAdapter{r: r}
}

func (s *serviceAdapter) RegisterChicCRMServices(loginData models.RegisterRequest) (models.RegisterResponses, error) {
	// func (s *serviceAdapter) RegisterChicCRMServices(loginData models.RegisterRequest) error {
	// responses, err := s.r.RegisterChicCRMSRepositoris(loginData)
	const (
		keyUsername = "e230944a-e25d-4674-83c7-436f7085086e"
		keyPassword = "LLodeHVjIDV-N13thflXkjWZuu1y4rCo723BGOLQ8RYGAalYETJz5HmsYx5MXwfH3mgTXw93UxtzVfPgzGNYCw"
	)

	if len(loginData.Mobile_phone) != 10 {
		log.Printf(`{"status": "Error", "message": "Mobile phone must be 10 digits 089-XXX-XXXX"}`)
		return models.RegisterResponses{}, errors.New("mobile phone must be 10 digits 089-XXX-XXXX")
	}

	generatedPassword, err := generate.GenerateRandomPassword(8)
	if err != nil {
		log.Printf("Failed to generate password. Error:%v\n", err)
		return models.RegisterResponses{}, err
	}
	fmt.Println(generatedPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password. Error:%v\n", err)
		return models.RegisterResponses{}, err
	}
	hashedPasswordString := string(hashedPassword)

	if !strings.Contains(loginData.Username, "@") || strings.HasPrefix(loginData.Username, "@") || strings.HasSuffix(loginData.Username, "@") || strings.Count(loginData.Username, "@") != 1 {
		return models.RegisterResponses{}, errors.New("username must be a valid email address")
	}
	emailDomain := loginData.Username
	splitEmail := strings.Split(emailDomain, "@")
	domain := "@" + splitEmail[1]
	// fmt.Println(splitEmail, splitEmail[1])

	token, err := auth.CreateTokenI(loginData.Username)
	if err != nil {
		log.Printf("Failed to create JWT for send email")
		return models.RegisterResponses{}, err
	}
	registrationComplete := make(chan bool)
	// defer close(registrationComplete) // do not use

	var (
		cipherProvince, cipherDistrict, cipherSubdistrict, cipherZipcode, cipherCreateLocation, cipherUrlLogo, cipherCompanyNameEN, cipherCompanyDomain, cipherCompanyMobile, cipherCompanyAlias, cipherTitle, cipherUsername, cipherFirstnameEn, cipherSurnamEn, cipherMobilePhone, cipherAddressNo, cipherAddress1En, cipherCompanyGeolo, cipherJobTitle string
		errChan                                                                                                                                                                                                                                                                                                                                            = make(chan error, 19)
	)
	go func() {
		<-registrationComplete
		to := loginData.Username
		subject := "Welcome! You have successfully registered."
		body := "Please click the link provided below to Login<br>" +
			"Email: " + loginData.Username + "<br>" +
			// "<a href='https://partnerdemo.tracthai.com/resetpassword/?token=" + token + "'>Confirm Link</a><br>"
			"<a href='http://localhost:3000/ResetPassword/?token=" + token + "'>Confirm Link</a><br>"
		if err := sendEmailFunctions.SendEmailRegister(to, subject, body); err != nil {
			log.Printf("เกิดข้อผิดพลาดในการส่งอีเมล: %s", err.Error())
			// notifyAdmin("เกิดข้อผิดพลาดในการส่งอีเมล: " + err.Error()) // ฟังก์ชันสมมุติเพื่อส่งการแจ้งเตือน
			return
		}
	}()
	// func notifyAdmin(message string) {
	// 	// โค้ดสำหรับส่งข้อความไปยัง Slack, email, SMS หรือระบบแจ้งเตือนอื่นๆ
	// }
	go func() {
		var err error
		cipherProvince, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Province, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherDistrict, err = encrypt.SendToFortanixSDKMSTokenization(loginData.District, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherSubdistrict, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Sub_district, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherZipcode, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Zipcode, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherCreateLocation, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Create_location, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherUrlLogo, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Url_logo, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherCompanyNameEN, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Company_name_en, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherCompanyDomain, err = encrypt.SendToFortanixSDKMSTokenization(domain, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherCompanyMobile, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Company_mobile, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherCompanyAlias, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Company_alias, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherTitle, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Title, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherUsername, err = encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(loginData.Username, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherFirstnameEn, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Firstname_en, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherSurnamEn, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Surname_en, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherMobilePhone, err = encrypt.SendToFortanixSDKMSTokenizationPhoneForMasking(loginData.Mobile_phone, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherAddressNo, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Address_no, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherAddress1En, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Address1_en, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherCompanyGeolo, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Company_geolo, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherJobTitle, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Job_title, keyUsername, keyPassword)
		errChan <- err
	}()
	// go func() { // not use yet remain for references for update
	// 	var err error
	// 	cipherUpdateLocation, err = encrypt.SendToFortanixSDKMSTokenization(loginData.Update_location, keyUsername, keyPassword)
	// 	errChan <- err
	// }()
	for i := 0; i < 19; i++ {
		if err := <-errChan; err != nil {
			log.Println(err)
			return models.RegisterResponses{}, err
		}
	}
	encryptData := models.EncryptedRegisterRequest{
		// CipherUpdate_Location: cipherUpdateLocation, // not use yet
		CipherProvince:        cipherProvince,
		CipherDistrict:        cipherDistrict,
		CipherSub_district:    cipherSubdistrict,
		CipherZipcode:         cipherZipcode,
		CipherCreate_Location: cipherCreateLocation,
		CipherUrl_logo:        cipherUrlLogo,
		CipherCompany_name_en: cipherCompanyNameEN,
		CipherCompany_domain:  cipherCompanyDomain,
		CipherCompany_mobile:  cipherCompanyMobile,
		CipherCompany_alias:   cipherCompanyAlias,
		CipherTitle:           cipherTitle,
		CipherUsername:        cipherUsername,
		CipherFirstname_en:    cipherFirstnameEn,
		CipherSurname_en:      cipherSurnamEn,
		CipherMobile_phone:    cipherMobilePhone,
		CipherAddress_no:      cipherAddressNo,
		CipherAddress1_en:     cipherAddress1En,
		CipherCompany_geolo:   cipherCompanyGeolo,
		CipherJob_title:       cipherJobTitle,
		HashPassword:          hashedPasswordString,
	}
	responses, err := s.r.RegisterChicCRMSRepositoris(encryptData, loginData)
	if err != nil {
		return responses, err
	}
	registrationComplete <- true
	return responses, nil

}

func (s *serviceAdapter) ValidateDomainChicCRMServices(validateDomainRequest models.ValidateDomainRequest) (models.ValidateDomainResponse, error) {
	if !strings.Contains(validateDomainRequest.Username, "@") || strings.HasPrefix(validateDomainRequest.Username, "@") || strings.HasSuffix(validateDomainRequest.Username, "@") || strings.Count(validateDomainRequest.Username, "@") != 1 {
		return models.ValidateDomainResponse{}, errors.New("username must be a valid email address")
	}
	emailDomainNotAllowed := map[string]bool{ //** ignored for dev
		"@gmail.com":   true,
		"@yahoo.com":   true,
		"@hotmail.com": true,
		"@outlook.com": true,
	}
	emailParts := strings.Split(validateDomainRequest.Username, "@")
	if len(emailParts) == 2 && emailDomainNotAllowed["@"+emailParts[1]] {
		log.Println("Error to check valid email gmail, yahoo, hotmail, outlook")
		return models.ValidateDomainResponse{}, errors.New("only company email allowed")
	}
	email := validateDomainRequest.Username
	splitEmail := strings.Split(email, "@")

	// // ตรวจสอบว่าได้สองส่วนหรือไม่
	// if len(splitEmail) != 2 {
	// 	log.Println("Error to splitEmail for validate Domain")
	// 	return models.ValidateDomainResponse{}, errors.New("invalid email format")
	// }

	// ใช้ส่วนที่สองของอีเมล (คือส่วนของโดเมน)
	emailDomain := "@" + splitEmail[1]
	const (
		keyUsername = "e230944a-e25d-4674-83c7-436f7085086e"
		keyPassword = "LLodeHVjIDV-N13thflXkjWZuu1y4rCo723BGOLQ8RYGAalYETJz5HmsYx5MXwfH3mgTXw93UxtzVfPgzGNYCw"
	)
	var (
		cipherUsernameToken, cipherUsernameSuffix string
		errChan1                                  = make(chan error, 2)
	)
	go func() {
		var err error
		cipherUsernameToken, err = encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(validateDomainRequest.Username, keyUsername, keyPassword)
		errChan1 <- err
	}()
	go func() {
		var err error
		cipherUsernameSuffix, err = encrypt.SendToFortanixSDKMSTokenization(emailDomain, keyUsername, keyPassword)
		errChan1 <- err
	}()
	for i := 0; i < 2; i++ {
		if err := <-errChan1; err != nil {
			log.Println(err)
			return models.ValidateDomainResponse{}, err
		}
	}
	// fmt.Println(cipherUsernameSuffix, cipherUsernameToken)

	validateResponse, err := s.r.ValidateDomainChicCRMSRepositoris(validateDomainRequest, cipherUsernameSuffix, cipherUsernameToken)
	if err != nil {
		return models.ValidateDomainResponse{}, err
	}

	var (
		companyLogoDecrypt, companynameDecrypt, domainDecrypt, companyAliasDecrypt, addressNoDecrypt, address1EnDecrypt, geolocationDecrypt, companyPhoneDecrypt, zipcodeDecrypt, provinceDecrypt, districtDecrypt, subdistrictDecrypt, timezoneDecrypt, countryDecrypt []byte
		errChan                                                                                                                                                                                                                                                         = make(chan error, 14)
	)
	go func() {
		var err error
		detokenizeLogo, err := decrypt.Detokenize(validateResponse.CompanyLogo)
		if err != nil {
			errChan <- err
			return
		}
		companyLogoDecrypt, err = base64.StdEncoding.DecodeString(detokenizeLogo)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeCompanyname, err := decrypt.Detokenize(validateResponse.Companyname)
		if err != nil {
			errChan <- err
			return
		}
		companynameDecrypt, err = base64.StdEncoding.DecodeString(detokenizeCompanyname)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeDomain, err := decrypt.Detokenize(validateResponse.Domain)
		if err != nil {
			errChan <- err
			return
		}
		domainDecrypt, err = base64.StdEncoding.DecodeString(detokenizeDomain)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeCompanyAlias, err := decrypt.Detokenize(validateResponse.CompanyAlias)
		if err != nil {
			errChan <- err
			return
		}
		companyAliasDecrypt, err = base64.StdEncoding.DecodeString(detokenizeCompanyAlias)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeAddressNo, err := decrypt.Detokenize(validateResponse.AddressNo)
		if err != nil {
			errChan <- err
			return
		}
		addressNoDecrypt, err = base64.StdEncoding.DecodeString(detokenizeAddressNo)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeAddress1En, err := decrypt.Detokenize(validateResponse.Address1En)
		if err != nil {
			errChan <- err
			return
		}
		address1EnDecrypt, err = base64.StdEncoding.DecodeString(detokenizeAddress1En)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeGeolocation, err := decrypt.Detokenize(validateResponse.Geolocation)
		if err != nil {
			errChan <- err
			return
		}
		geolocationDecrypt, err = base64.StdEncoding.DecodeString(detokenizeGeolocation)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeCompanyPhone, err := decrypt.Detokenize(validateResponse.CompanyPhone)
		if err != nil {
			errChan <- err
			return
		}
		companyPhoneDecrypt, err = base64.StdEncoding.DecodeString(detokenizeCompanyPhone)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeZipcode, err := decrypt.Detokenize(validateResponse.Zipcode)
		if err != nil {
			errChan <- err
			return
		}
		zipcodeDecrypt, err = base64.StdEncoding.DecodeString(detokenizeZipcode)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeProvince, err := decrypt.Detokenize(validateResponse.Province)
		if err != nil {
			errChan <- err
			return
		}
		provinceDecrypt, err = base64.StdEncoding.DecodeString(detokenizeProvince)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeDistrict, err := decrypt.Detokenize(validateResponse.District)
		if err != nil {
			errChan <- err
			return
		}
		districtDecrypt, err = base64.StdEncoding.DecodeString(detokenizeDistrict)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeSubdistrict, err := decrypt.Detokenize(validateResponse.SubDistrict)
		if err != nil {
			errChan <- err
			return
		}
		subdistrictDecrypt, err = base64.StdEncoding.DecodeString(detokenizeSubdistrict)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeTimeZone, err := decrypt.Detokenize(validateResponse.Timezone)
		if err != nil {
			errChan <- err
			return
		}
		timezoneDecrypt, err = base64.StdEncoding.DecodeString(detokenizeTimeZone)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeCountry, err := decrypt.Detokenize(validateResponse.Country)
		if err != nil {
			errChan <- err
			return
		}
		countryDecrypt, err = base64.StdEncoding.DecodeString(detokenizeCountry)
		errChan <- err
	}()
	for i := 0; i < 14; i++ {
		if err := <-errChan; err != nil {
			log.Println(err)
			return models.ValidateDomainResponse{}, err
		}
	}
	validateResponse = models.ValidateDomainResponse{
		CompanyID:    validateResponse.CompanyID,
		CompanyLogo:  string(companyLogoDecrypt),
		Companyname:  string(companynameDecrypt),
		Domain:       string(domainDecrypt),
		Alpha2Code:   validateResponse.Alpha2Code,
		LanguageName: validateResponse.LanguageName,
		CompanyAlias: string(companyAliasDecrypt),
		Currency:     validateResponse.Currency,
		AddressNo:    string(addressNoDecrypt),
		Address1En:   string(address1EnDecrypt),
		Geolocation:  string(geolocationDecrypt),
		CompanyPhone: string(companyPhoneDecrypt),
		Zipcode:      string(zipcodeDecrypt),
		Province:     string(provinceDecrypt),
		District:     string(districtDecrypt),
		SubDistrict:  string(subdistrictDecrypt),
		Timezone:     string(timezoneDecrypt),
		Country:      string(countryDecrypt),
	}

	return validateResponse, nil
}
