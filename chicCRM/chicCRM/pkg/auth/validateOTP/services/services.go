package services

import (
	"chicCRM/pkg/auth/validateOTP/models"
	"chicCRM/pkg/utilts/generate"
	"chicCRM/pkg/utilts/sendEmailFunctions"
	"chicCRM/pkg/utilts/utility"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type ServicePort interface {
	RequestEmailForValidateOTPChicCRMServices(email string) (string, error)
	ValidateOTPFromRequestEmailChicCRMServices(receivedOTP, receivedReferenceID string) error
	QrTOTPChicCRMServices(accountName string, value int) (string, bool, error)
	ValidateQrTOTPChicCRMServices(validateQrTOTP models.ValidateQrTOTP) (bool, error)
	DeleteKeyQrTOTPChicCRMServices(deleteKeyQrTOTP string) error
}

type serviceAdapter struct {
	otpStore   map[string]int
	otpKeys    map[string]*otp.Key
	refIDStore map[string]string
}

func NewServiceAdapter() ServicePort {
	return &serviceAdapter{
		otpStore:   make(map[string]int), // int input for validate
		otpKeys:    make(map[string]*otp.Key),
		refIDStore: make(map[string]string),
	}
}

func (s *serviceAdapter) RequestEmailForValidateOTPChicCRMServices(email string) (string, error) {
	otp := generate.GenerateOTP(6)
	referenceID := generate.GenerateReferenceID(6)
	s.otpStore[email] = otp
	s.refIDStore[email] = referenceID
	// subject := "Verify your identity with an OTP."
	// body := fmt.Sprintf("Your OTP code is: %d", otp) // %d for integer
	subject := "Verify your identity with an OTP."
	body := fmt.Sprintf("Please Use OTP provided below to verify your identity<br>Your OTP code is: %d<br>Your Reference No is: %s", otp, referenceID)

	if err := sendEmailFunctions.SendEmailOTP(email, subject, body); err != nil {
		log.Printf("Failed to send OTP via Email: %v\n", err)
		return "", err
	}
	return referenceID, nil
}

// func (s *serviceAdapter) ValidateOTPFromRequestEmailChicCRMServices(email string, receivedOTP, receivedReferenceID string) error { << use email, refID, OTP to validate
//
//		expectedOTP, exists := s.otpStore[email]
//		expectedReferenceID, existsRefID := s.refIDStore[email]
//		if !exists || !existsRefID {
//			return errors.New("OTP or reference ID was not requested for this email")
//		}
//		if receivedOTP == strconv.Itoa(expectedOTP) && receivedReferenceID == expectedReferenceID {
//			delete(s.otpStore, email)
//			delete(s.refIDStore, email)
//			return nil
//		}
//		return errors.New("OTP or reference ID is invalid")
//	}
func (s *serviceAdapter) ValidateOTPFromRequestEmailChicCRMServices(receivedOTP, receivedReferenceID string) error {
	for email, expectedOTP := range s.otpStore {
		expectedReferenceID, existsRefID := s.refIDStore[email]

		if existsRefID && receivedOTP == strconv.Itoa(expectedOTP) && receivedReferenceID == expectedReferenceID {
			delete(s.otpStore, email)
			delete(s.refIDStore, email)
			return nil
		}
	}
	return errors.New("OTP or reference ID is invalid")
}
func (s *serviceAdapter) QrTOTPChicCRMServices(accountName string, value int) (string, bool, error) {
	if value != 1 {
		return "", false, nil
	}
	if _, found := s.otpKeys[accountName]; found {
		return "", false, errors.New("AccountName already exists")
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Chiccrm",
		AccountName: accountName,
		Algorithm:   otp.AlgorithmSHA512,
		SecretSize:  32,
		Period:      30,
	})
	if err != nil {
		return "", false, err
	}

	s.otpKeys[accountName] = key
	qrCodeURL := utility.GenerateQRCodeURL(key)

	return qrCodeURL, true, nil
}

func (s *serviceAdapter) ValidateQrTOTPChicCRMServices(validateQrTOTP models.ValidateQrTOTP) (bool, error) {
	key, found := s.otpKeys[validateQrTOTP.AccountName]
	if !found {
		return false, errors.New("no OTP key available, account name does not match")
	}

	serverOTP, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return false, err
	}

	return validateQrTOTP.OTP == serverOTP, nil
}

func (s *serviceAdapter) DeleteKeyQrTOTPChicCRMServices(deleteKeyQrTOTP string) error {
	if _, found := s.otpKeys[deleteKeyQrTOTP]; found {
		delete(s.otpKeys, deleteKeyQrTOTP)
		return nil
	}
	return errors.New("OTP key for " + deleteKeyQrTOTP + " not found")
}
