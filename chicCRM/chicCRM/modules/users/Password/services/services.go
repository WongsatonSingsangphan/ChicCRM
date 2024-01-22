package services

import (
	"chicCRM/modules/users/Password/models"
	"chicCRM/modules/users/Password/repositories"
	"chicCRM/pkg/utilts/encrypt"
	"chicCRM/pkg/utilts/sendEmailFunctions"
	"errors"
	"fmt"
	"log"
)

type ServicePort interface {
	InitPasswordChicCRMServices(initPasswordAdditional models.InitPasswordAdditional, initPassword models.InitPassword) error
	ChangePasswordChicCRMServices(changePassword models.ChangePassword) error
	RequestResetPasswordChicCRMServices(requestResetPassword models.RequestResetPassword) error
	ResetPasswordChicCRMServices(resetPassword models.ResetPassword, ResetPasswordAdditional models.InitPasswordAdditional) error
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

func (s *serviceAdapter) InitPasswordChicCRMServices(initPasswordAdditional models.InitPasswordAdditional, initPassword models.InitPassword) error {
	cipherUsername, err := encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(initPassword.Username, keyUsername, keyPassword)
	if err != nil {
		log.Printf("Failed to tokenize username. Error:%v\n", err)
		return err
	}
	if initPassword.Newpassword == "" || len(initPassword.Newpassword) < 8 {
		return errors.New("password must not be empty and must be at least 8 characters long")
	}

	err = s.r.InitPasswordChicCRMSRepositoris(initPasswordAdditional, initPassword, cipherUsername)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *serviceAdapter) ChangePasswordChicCRMServices(changePassword models.ChangePassword) error {
	cipherUsernameChangePassword, err := encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(changePassword.Username, keyUsername, keyPassword)
	if err != nil {
		log.Printf("Failed to tokenize username. Error:%v\n", err)
		return err
	}
	fmt.Println(cipherUsernameChangePassword, changePassword.Oldpassword, changePassword.Newpassword)
	if changePassword.Oldpassword == "" || changePassword.Newpassword == "" || len(changePassword.Oldpassword) < 8 || len(changePassword.Newpassword) < 8 {
		return errors.New("password must not be empty and must be at least 8 characters long")

	}
	err = s.r.ChangePasswordChicCRMSRepositoris(changePassword, cipherUsernameChangePassword)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (s *serviceAdapter) RequestResetPasswordChicCRMServices(requestResetPassword models.RequestResetPassword) error {
	cipherUsernameRequest, err := encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(requestResetPassword.Username, keyUsername, keyPassword)
	if err != nil {
		log.Printf("Failed to tokenize username. Error:%v\n", err)
		return err
	}
	// fmt.Println(cipherUsernameRequest)
	token, err := s.r.RequestResetPasswordChicCRMRepositories(cipherUsernameRequest)
	if err != nil {
		log.Println(err)
		return err
	}
	// fmt.Println(cipherUsernameRequest)
	to := requestResetPassword.Username
	subject := "Reset your password."
	body := "Please click the link provided below to reset password.<br>" +
		"Email: " + requestResetPassword.Username + "<br>" +
		// "<a href='https://partnerdemo.tracthai.com/resetpassword/?token=" + token + "'>Confirm Link</a><br>"
		// "<a href='http://localhost:3000/ResetPassword/?token=" + token + "'>Reset Password</a><br>"
		"<a href='http://localhost:3000/Changepassword/?token=" + token + "'>Reset Password</a><br>"
	if err := sendEmailFunctions.SendEmailRegister(to, subject, body); err != nil {
		log.Printf("เกิดข้อผิดพลาดในการส่งอีเมล: %s", err.Error())
		// notifyAdmin("เกิดข้อผิดพลาดในการส่งอีเมล: " + err.Error()) // ฟังก์ชันสมมุติเพื่อส่งการแจ้งเตือน
		return err
	}
	return nil
}

func (s *serviceAdapter) ResetPasswordChicCRMServices(resetPassword models.ResetPassword, ResetPasswordAdditional models.InitPasswordAdditional) error {
	cipherUsernameResetPassword, err := encrypt.SendToFortanixSDKMSTokenizationEmailForMasking(resetPassword.Username, keyUsername, keyPassword)
	if err != nil {
		log.Printf("Failed to tokenize username. Error:%v\n", err)
		return err
	}
	if resetPassword.NewPassword == "" || len(resetPassword.NewPassword) < 8 {
		return errors.New("password must not be empty and must be at least 8 characters long")
	}
	err = s.r.ResetPasswordChicCRMRepositories(resetPassword, cipherUsernameResetPassword)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
