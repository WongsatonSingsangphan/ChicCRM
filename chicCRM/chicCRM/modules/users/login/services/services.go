package services

import (
	"chicCRM/modules/users/login/models"
	"chicCRM/modules/users/login/repositories"
	"chicCRM/pkg/auth"
	"chicCRM/pkg/utilts/decrypt"
	"chicCRM/pkg/utilts/encrypt"
	"encoding/base64"
	"errors"
	"log"
	"strings"
	"time"
)

type ServicePort interface {
	LoginChicCRMServices(loginRequest models.LoginRequest) (string, error)
	LoginTeamleadSecuredocServices(loginRequestTeamlead models.LoginRequestTeamlead) (string, error)
}

type serviceAdapter struct {
	r repositories.RepositoryPort
}

func NewServiceAdapter(r repositories.RepositoryPort) ServicePort {
	return &serviceAdapter{r: r}
}

const (
	keyUsername = "e230944a-e25d-4674-83c7-436f7085086e"
	keyPassword = "LLodeHVjIDV-N13thflXkjWZuu1y4rCo723BGOLQ8RYGAalYETJz5HmsYx5MXwfH3mgTXw93UxtzVfPgzGNYCw"
)

func (s *serviceAdapter) LoginChicCRMServices(loginRequest models.LoginRequest) (string, error) {
	if !strings.Contains(loginRequest.Username, "@") || strings.HasPrefix(loginRequest.Username, "@") || strings.HasSuffix(loginRequest.Username, "@") || strings.Count(loginRequest.Username, "@") != 1 {
		return "", errors.New("username must be a valid email address")
	}
	cipherUsernameLogin, err := encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(loginRequest.Username, keyUsername, keyPassword)
	if err != nil {
		log.Println(err)
		return "", err
	}
	loginResponse, err := s.r.LoginChicCRMSRepositoris(loginRequest, cipherUsernameLogin)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var (
		titleDecrypt       []byte
		usernameDecrypt    []byte
		firstnameDecrypt   []byte
		surnameDecrypt     []byte
		mobileDecrypt      []byte
		companynameDecrypt []byte
		jobtitleDecrypt    []byte
		domainDecrypt      []byte
		companyLogoDecrypt []byte
		errChan            = make(chan error, 9)
	)
	go func() {
		var err error
		detokenizeTitle, err := decrypt.Detokenize(loginResponse.TitleToken)
		if err != nil {
			errChan <- err
			return
		}
		titleDecrypt, err = base64.StdEncoding.DecodeString(detokenizeTitle)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeUsername, err := decrypt.DetokenizationEmailForMasking(loginResponse.UsernameToken)
		if err != nil {
			errChan <- err
			return
		}
		usernameDecrypt, err = base64.StdEncoding.DecodeString(detokenizeUsername)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizefirstname, err := decrypt.Detokenize(loginResponse.FirstnameTokenEN)
		if err != nil {
			errChan <- err
			return
		}
		firstnameDecrypt, err = base64.StdEncoding.DecodeString(detokenizefirstname)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizaSurname, err := decrypt.Detokenize(loginResponse.SurnameTokenEN)
		if err != nil {
			errChan <- err
			return
		}
		surnameDecrypt, err = base64.StdEncoding.DecodeString(detokenizaSurname)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeMobile, err := decrypt.DetokenizationPhoneForMasking(loginResponse.MobilePhoneToken)
		if err != nil {
			errChan <- err
			return
		}
		mobileDecrypt, err = base64.StdEncoding.DecodeString(detokenizeMobile)
		errChan <- err
	}()
	go func() {
		detokenizeCompanyname, err := decrypt.Detokenize(loginResponse.CompanynameToken)
		if err != nil {
			errChan <- err
			return
		}
		companynameDecrypt, err = base64.StdEncoding.DecodeString(detokenizeCompanyname)
		errChan <- err
	}()
	go func() {
		detokenizJobtitle, err := decrypt.Detokenize(loginResponse.JobTitleToken)
		if err != nil {
			errChan <- err
			return
		}
		jobtitleDecrypt, err = base64.StdEncoding.DecodeString(detokenizJobtitle)
		errChan <- err
	}()
	go func() {
		detokenizeDomain, err := decrypt.Detokenize(loginResponse.DomainToken)
		if err != nil {
			errChan <- err
			return
		}
		domainDecrypt, err = base64.StdEncoding.DecodeString(detokenizeDomain)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeCompanyLogo, err := decrypt.Detokenize(loginResponse.CompanyLogoToken)
		if err != nil {
			errChan <- err
			return
		}
		companyLogoDecrypt, err = base64.StdEncoding.DecodeString(detokenizeCompanyLogo)
		errChan <- err
	}()

	for i := 0; i < 9; i++ {
		if err := <-errChan; err != nil {
			log.Println(err)
			return "", err
		}
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := models.JwtResponse{
		ID:                   loginResponse.ID,
		TitleOriginal:        string(titleDecrypt),
		UsernameToken:        loginResponse.UsernameToken,
		UsernameOriginal:     string(usernameDecrypt),
		FirstnameTokenEN:     loginResponse.FirstnameTokenEN,
		FirstnameOriginal:    string(firstnameDecrypt),
		SurnameTokenEN:       loginResponse.SurnameTokenEN,
		SurnameTokenOriginal: string(surnameDecrypt),
		MobilePhoneOriginal:  string(mobileDecrypt),
		Requires_action:      loginResponse.Requires_action,
		CompanyID:            loginResponse.CompanyID,
		CompanynameToken:     loginResponse.CompanynameToken,
		CompanynameOrginal:   string(companynameDecrypt),
		CompanyLogoOriginal:  string(companyLogoDecrypt),
		JobTitleToken:        loginResponse.JobTitleToken,
		JobTitleOriginal:     string(jobtitleDecrypt),
		DomainToken:          loginResponse.DomainToken,
		DomainOriginal:       string(domainDecrypt),
		Role:                 loginResponse.Role,
		Exp:                  expirationTime.Unix(),
		// RequestID:            loginResponse.RequestID,
	}
	tokenJWT, err := auth.CreateJWT(claims)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenJWT, nil
}

func (s *serviceAdapter) LoginTeamleadSecuredocServices(loginRequestTeamlead models.LoginRequestTeamlead) (string, error) {
	if !strings.Contains(loginRequestTeamlead.Username, "@") || strings.HasPrefix(loginRequestTeamlead.Username, "@") || strings.HasSuffix(loginRequestTeamlead.Username, "@") || strings.Count(loginRequestTeamlead.Username, "@") != 1 {
		return "", errors.New("username must be a valid email address")
	}
	cipherTeamleadUsername, err := encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(loginRequestTeamlead.Username, keyUsername, keyPassword)
	if err != nil {
		log.Printf("Failed to encrypt teamleadUsername. Error:%v", err)
		return "", err
	}
	loginResponseTeamlead, err := s.r.LoginTeamleadSecuredocRepositories(loginRequestTeamlead, cipherTeamleadUsername)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var (
		usernameTeamleadDecrypt    []byte
		firstnameTeamleadDecrypt   []byte
		surnameTeamleadDecrypt     []byte
		mobilePhoneTeamleadDecrypt []byte
		errChan                    = make(chan error, 4)
	)
	go func() {
		var err error
		detokenizeUsernameTeamlead, err := decrypt.DetokenizationEmailForMasking(loginResponseTeamlead.UsernameToken)
		if err != nil {
			errChan <- err
			return
		}
		usernameTeamleadDecrypt, err = base64.StdEncoding.DecodeString(detokenizeUsernameTeamlead)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeFirstnameTeamlead, err := decrypt.Detokenize(loginResponseTeamlead.FirstnameToken)
		if err != nil {
			errChan <- err
			return
		}
		firstnameTeamleadDecrypt, err = base64.StdEncoding.DecodeString(detokenizeFirstnameTeamlead)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizeSurnameTeamlead, err := decrypt.Detokenize(loginResponseTeamlead.SurnameToken)
		if err != nil {
			errChan <- err
			return
		}
		surnameTeamleadDecrypt, err = base64.StdEncoding.DecodeString(detokenizeSurnameTeamlead)
		errChan <- err
	}()
	go func() {
		var err error
		detokenizemobilePhoneTeamlead, err := decrypt.DetokenizationPhoneForMasking(loginResponseTeamlead.Mobile_phoneToken)
		if err != nil {
			errChan <- err
			return
		}
		mobilePhoneTeamleadDecrypt, err = base64.StdEncoding.DecodeString(detokenizemobilePhoneTeamlead)
		errChan <- err
	}()
	for i := 0; i < 4; i++ {
		if err := <-errChan; err != nil {
			log.Println(err)
			return "", err
		}
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := models.JwtResponseTeamleadSecuredog{
		TeamleadID:        loginResponseTeamlead.TeamleadID,
		UsernameToken:     loginResponseTeamlead.UsernameToken,
		Username:          string(usernameTeamleadDecrypt),
		FirstnameToken:    loginResponseTeamlead.FirstnameToken,
		Firstname:         string(firstnameTeamleadDecrypt),
		SurnameToken:      loginResponseTeamlead.SurnameToken,
		Surname:           string(surnameTeamleadDecrypt),
		Mobile_phoneToken: loginResponseTeamlead.Mobile_phoneToken,
		Mobile_phone:      string(mobilePhoneTeamleadDecrypt),
		Role:              loginResponseTeamlead.Role,
		Management_Level:  loginResponseTeamlead.Management_Level,
		Exp:               expirationTime.Unix(),
	}
	tokenJWT, err := auth.CreateJWTTeamlead(claims)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return tokenJWT, nil
}
