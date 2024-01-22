package repositories

import (
	"chicCRM/modules/users/login/models"
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type RepositoryPort interface {
	LoginChicCRMSRepositoris(loginRequest models.LoginRequest, cipherUsernameLogin string) (models.LoginResponse, error)
	LoginTeamleadSecuredocRepositories(loginTeamleadRequest models.LoginRequestTeamlead, cipherTeamleadUsername string) (models.LoginResponseTeamlead, error)
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}

func (r *repositoryAdapter) LoginChicCRMSRepositoris(loginRequest models.LoginRequest, cipherUsernameLogin string) (models.LoginResponse, error) {
	var loginResponse models.LoginResponse
	var organizeMemberUUID string
	err := r.db.QueryRow("SELECT orgmb_id FROM organize_member WHERE orgmb_email = $1", cipherUsernameLogin).Scan(&organizeMemberUUID)
	if err != nil {
		return models.LoginResponse{}, errors.New("user not found")
	}
	var hashedPassword string
	err = r.db.QueryRow("SELECT orgmbcr_password FROM organize_member_credential WHERE  orgmbcr_orgmb_id = $1", organizeMemberUUID).Scan(&hashedPassword)
	if err != nil {
		return models.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password))
	if err != nil {
		log.Printf("Password not match :%v\n", err)
		return models.LoginResponse{}, errors.New("username or password is invalid")
	}
	var roleUUID string
	err = r.db.QueryRow("SELECT orgmb_id, orgmb_holder, orgmb_title, orgmb_email, orgmb_name, orgmb_surname, orgmb_role, orgmb_mobile FROM organize_member WHERE orgmb_id = $1", organizeMemberUUID).Scan(
		&loginResponse.ID, &loginResponse.CompanyID, &loginResponse.TitleToken, &loginResponse.UsernameToken, &loginResponse.FirstnameTokenEN, &loginResponse.SurnameTokenEN,
		&roleUUID, &loginResponse.MobilePhoneToken)
	if err != nil {
		log.Printf("Error retriving data from organize_member: %v", err)
		return models.LoginResponse{}, err
	}
	err = r.db.QueryRow("SELECT org_logo, org_name_en, org_domain FROM organize_master WHERE org_id = $1", loginResponse.CompanyID).Scan(&loginResponse.CompanyLogoToken, &loginResponse.CompanynameToken, &loginResponse.DomainToken)
	if err != nil {
		log.Printf("Error retriving data from organize_master: %v", err)
		return models.LoginResponse{}, err
	}
	err = r.db.QueryRow("SELECT orgrl_name_en FROM organize_role WHERE orgrl_id = $1", roleUUID).Scan(&loginResponse.JobTitleToken)
	if err != nil {
		log.Printf("Error retriving data from organize_role: %v", err)
		return models.LoginResponse{}, err
	}
	err = r.db.QueryRow("SELECT orgmbcr_level, orgmbcr_reqaction FROM organize_member_credential WHERE orgmbcr_orgmb_id = $1", organizeMemberUUID).Scan(&loginResponse.Role, &loginResponse.Requires_action)
	if err != nil {
		log.Printf("Error retriving data from organize_member_credential: %v", err)
		return models.LoginResponse{}, err
	}
	// var organizeMemberAuthorizationUUID string // not use in this moment and decrease response speed
	// err = r.db.QueryRow("SELECT orgmbat_id FROM organize_member_authorization WHERE orgmbat_orgmbid = $1", organizeMemberUUID).Scan(&organizeMemberAuthorizationUUID)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		log.Printf("No authorization data for user: %v\n", organizeMemberUUID)
	// 	} else {
	// 		log.Printf("Error retrieving data from organize_member_authorization: %v\n", err)
	// 		return models.LoginResponse{}, err
	// 	}
	// }
	// if organizeMemberAuthorizationUUID != "" {
	// 	rows, err := r.db.Query("SELECT scdact_id FROM securedoc_activity WHERE scdact_actor = $1", organizeMemberAuthorizationUUID)
	// 	if err != nil {
	// 		log.Printf("Error querying securedoc_activity: %v\n", err)
	// 		return models.LoginResponse{}, err
	// 	}
	// 	defer rows.Close()
	// 	for rows.Next() {
	// 		var requestID string
	// 		if err := rows.Scan(&requestID); err != nil {
	// 			log.Printf("Error scanning securedoc_activity row: %v\n", err)
	// 			return models.LoginResponse{}, err
	// 		}
	// 		loginResponse.RequestID = append(loginResponse.RequestID, requestID)
	// 	}
	// 	if err = rows.Err(); err != nil {
	// 		log.Printf("Error iterating through securedoc_activity rows: %v\n", err)
	// 		return models.LoginResponse{}, err
	// 	}
	// }
	return loginResponse, nil
}

func (r *repositoryAdapter) LoginTeamleadSecuredocRepositories(loginTeamleadRequest models.LoginRequestTeamlead, cipherTeamleadUsername string) (models.LoginResponseTeamlead, error) {
	var loginResponseTeamlead models.LoginResponseTeamlead
	err := r.db.QueryRow("SELECT orgtl_trac_id, orgtl_trac_username, orgtl_trac_phone, orgtl_trac_firstname, orgtl_trac_lastname, orgtl_trac_role, orgtl_level FROM organize_teamlead_trac WHERE orgtl_trac_username = $1 and orgtl_trac_password = $2", cipherTeamleadUsername, loginTeamleadRequest.Password).Scan(
		&loginResponseTeamlead.TeamleadID, &loginResponseTeamlead.UsernameToken, &loginResponseTeamlead.Mobile_phoneToken, &loginResponseTeamlead.FirstnameToken, &loginResponseTeamlead.SurnameToken, &loginResponseTeamlead.Role, &loginResponseTeamlead.Management_Level)
	if err != nil {
		log.Printf("username or password does not match. Error:%v", err)
		return models.LoginResponseTeamlead{}, errors.New("username or password does not match")
	}
	// fmt.Println(loginResponseTeamlead.TeamleadID)
	return loginResponseTeamlead, nil
}
