package repositories

import (
	"chicCRM/modules/users/Password/models"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type RepositoryPort interface {
	InitPasswordChicCRMSRepositoris(initPasswordAdditional models.InitPasswordAdditional, initPassword models.InitPassword, cipherUsername string) error
	ChangePasswordChicCRMSRepositoris(changePassword models.ChangePassword, cipherUsernameChangePassword string) error
	RequestResetPasswordChicCRMRepositories(cipherUsernameRequest string) (string, error)
	ResetPasswordChicCRMRepositories(resetPassword models.ResetPassword, cipherUsernameResetPassword string) error
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}

func (r *repositoryAdapter) InitPasswordChicCRMSRepositoris(initPasswordAdditional models.InitPasswordAdditional, initPassword models.InitPassword, cipherUsername string) error {
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(initPassword.Newpassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v \n", err)
		return err
	}
	var organizeMemberUUID string
	err = r.db.QueryRow("SELECT orgmb_id FROM organize_member WHERE orgmb_email = $1", cipherUsername).Scan(&organizeMemberUUID)
	if err != nil {
		log.Printf("Error retrieving orgmb_id(UUID) from organize_member: %v", err)
		return errors.New("user not found")
	}
	_, err = r.db.Exec("UPDATE organize_member_credential SET orgmbcr_password = $1, orgmbcr_reqaction = $2, orgmbcr_editor = $3, orgmbcr_blacklist_token = $4 WHERE orgmbcr_orgmb_id = $5", string(hashedNewPassword), initPassword.Requires_action, organizeMemberUUID, initPasswordAdditional.ExistingToken, organizeMemberUUID)
	if err != nil {
		log.Printf("Error updating password in organize_member_credential: %v\n", err)
		return err
	}

	return nil
}
func (r *repositoryAdapter) ChangePasswordChicCRMSRepositoris(changePassword models.ChangePassword, cipherUsernameChangePassword string) error {
	fmt.Println(cipherUsernameChangePassword, changePassword.Oldpassword, changePassword.Newpassword)
	var hashedOldPassword, organizeMemberUUID string
	err := r.db.QueryRow("SELECT orgmb_id FROM organize_member WHERE orgmb_email = $1", cipherUsernameChangePassword).Scan(&organizeMemberUUID)
	if err != nil {
		log.Printf("Error retriving orgmb_id(UUID) from organize_member: %v", err)
		return errors.New("user not found")
	}

	err = r.db.QueryRow("SELECT orgmbcr_password FROM organize_member_credential WHERE orgmbcr_orgmb_id = $1", organizeMemberUUID).Scan(&hashedOldPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return err
		}
		log.Printf("Error to SELECT orgmbcr_password FROM organize_member_credential. Error: %v\n", err)
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedOldPassword), []byte(changePassword.Oldpassword))
	if err != nil {
		log.Printf("Error Incorrect Old Password. Error: %v\n", err)
		return errors.New("incorrect old password") // better return with information of the error
		// return err
	}
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(changePassword.Newpassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error to hash new password. Error: %v", err)
		return err
	}
	_, err = r.db.Exec("UPDATE organize_member_credential SET orgmbcr_password = $1, orgmbcr_reqaction = $2, orgmbcr_editor = $3 WHERE orgmbcr_orgmb_id = $4", string(hashedNewPassword), changePassword.Requires_action, organizeMemberUUID, organizeMemberUUID)
	if err != nil {
		log.Printf("Error updating password in database: %v", err)
		return err
	}
	return nil
}
func (r *repositoryAdapter) RequestResetPasswordChicCRMRepositories(cipherUsernameRequest string) (string, error) {
	var organizeMemberUUID string
	// fmt.Println(cipherUsernameRequest)
	err := r.db.QueryRow("SELECT orgmb_id FROM organize_member WHERE orgmb_email = $1", cipherUsernameRequest).Scan(&organizeMemberUUID)
	if err != nil {
		log.Printf("Error retriving orgmb_id(UUID) from organize_member: %v", err)
		return "", errors.New("user not found")
	}
	var tokenForResetPassword string
	err = r.db.QueryRow("SELECT orgmbcr_blacklist_token FROM organize_member_credential WHERE orgmbcr_orgmb_id = $1", organizeMemberUUID).Scan(&tokenForResetPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return "", err
		}
		log.Printf("Error to SELECT orgmbcr_blacklist_token FROM organize_member_credential. Error: %v\n", err)
		return "", err
	}
	return tokenForResetPassword, nil
}

func (r *repositoryAdapter) ResetPasswordChicCRMRepositories(resetPassword models.ResetPassword, cipherUsernameResetPassword string) error {
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(resetPassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v \n", err)
		return err
	}
	var organizeMemberUUID string
	err = r.db.QueryRow("SELECT orgmb_id FROM organize_member WHERE orgmb_email = $1", cipherUsernameResetPassword).Scan(&organizeMemberUUID)
	if err != nil {
		log.Printf("Error to retriving orgmb_id(UUID) from organize_member: %v", err)
		return errors.New("user not found")
	}
	_, err = r.db.Exec("UPDATE organize_member_credential SET orgmbcr_password = $1, orgmbcr_editor = $2 WHERE orgmbcr_orgmb_id = $3", string(hashedNewPassword), organizeMemberUUID, organizeMemberUUID)
	if err != nil {
		log.Printf("Error updating password in organize_member_credential(ResetPasswordAPI): %v", err)
		return err
	}
	return nil
}
