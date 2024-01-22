package repositories

import (
	"chicCRM/modules/users/register/models"
	"chicCRM/pkg/utilts/encrypt"
	"database/sql"
	"errors"
	"log"
)

type RepositoryPort interface {
	// RegisterChicCRMSRepositoris(loginData models.RegisterRequest) (models.RegisterResponses, error)
	// RegisterChicCRMSRepositoris(loginData models.RegisterRequest) error
	RegisterChicCRMSRepositoris(encryptedData models.EncryptedRegisterRequest, loginData models.RegisterRequest) (models.RegisterResponses, error)
	ValidateDomainChicCRMSRepositoris(validateDomainRequest models.ValidateDomainRequest, cipherUsernameSuffix, cipherUsernameToken string) (models.ValidateDomainResponse, error)
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}

func (r *repositoryAdapter) RegisterChicCRMSRepositoris(encryptedData models.EncryptedRegisterRequest, loginData models.RegisterRequest) (models.RegisterResponses, error) {
	// func (r *repositoryAdapter) RegisterChicCRMSRepositoris(loginData models.RegisterRequest) error {
	var registerResponses models.RegisterResponses
	var checkMobilePhoneExists, checkEmailExists string
	err := r.db.QueryRow("SELECT orgmb_email FROM organize_member WHERE orgmb_email = $1", encryptedData.CipherUsername).Scan(&checkEmailExists)
	if err == nil {
		log.Printf("Email already exists")
		return models.RegisterResponses{}, errors.New("email already exists")
	}
	err = r.db.QueryRow("SELECT orgmb_mobile FROM organize_member WHERE orgmb_mobile = $1", encryptedData.CipherMobile_phone).Scan(&checkMobilePhoneExists)
	if err == nil {
		log.Printf("Mobile Phone already exists")
		return models.RegisterResponses{}, errors.New("mobile Phone already exists")
	}

	var ctyName, ctyID, timeZone, alpha2Code, languageNativeName, currencyName string
	// err := r.db.QueryRow("SELECT name, id FROM countries WHERE name = $1", encryptedData.CipherCountry).Scan(&ctyName, &ctyID)
	err = r.db.QueryRow("SELECT name, id, timezones, alpha2code, languagenativename, currency_name FROM countries WHERE name = $1", loginData.Country).Scan(&ctyName, &ctyID, &timeZone, &alpha2Code, &languageNativeName, &currencyName)
	if err != nil {
		log.Println(err)
		return models.RegisterResponses{}, err
		// return err
	}
	// fmt.Println(alpha2Code, languageNativeName, currencyName)
	const (
		keyUsername = "e230944a-e25d-4674-83c7-436f7085086e"
		keyPassword = "LLodeHVjIDV-N13thflXkjWZuu1y4rCo723BGOLQ8RYGAalYETJz5HmsYx5MXwfH3mgTXw93UxtzVfPgzGNYCw"
	)

	var (
		cipherNameToken, cipherTimezone string
		errChan                         = make(chan error, 2)
	)
	go func() {
		var err error
		cipherNameToken, err = encrypt.SendToFortanixSDKMSTokenization(ctyName, keyUsername, keyPassword)
		errChan <- err
	}()
	go func() {
		var err error
		cipherTimezone, err = encrypt.SendToFortanixSDKMSTokenization(timeZone, keyUsername, keyPassword)
		errChan <- err
	}()
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			log.Println(err)
			return models.RegisterResponses{}, err
		}
	}

	var countryUUID string
	err = r.db.QueryRow("SELECT id FROM country WHERE cty_alias = $1 and cty_id = $2", cipherNameToken, ctyID).Scan(&countryUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = r.db.QueryRow("INSERT INTO country (cty_alias, cty_id, cty_createlocation) VALUES ($1, $2, $3) RETURNING id", cipherNameToken, ctyID, encryptedData.CipherCreate_Location).Scan(&countryUUID)
			if err != nil {
				log.Println(err)
				return models.RegisterResponses{}, err
				// return err
			}
		} else {
			log.Println(err)
			return models.RegisterResponses{}, err
			// return err
		}
	}
	tx, err := r.db.Begin()
	if err != nil {
		log.Println(err)
		return models.RegisterResponses{}, err
		// return err
	}
	// defer tx.Rollback()
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			if rollbackErr != sql.ErrTxDone {
				log.Printf("could not rollback: %v\n", rollbackErr)
			}
		}
	}()

	var localAddressUUID string
	err = r.db.QueryRow("SELECT lca_id FROM local_address WHERE lca_country = $1 AND lca_zipcode = $2 AND lca_province_en = $3 AND lca_district_en = $4 and lca_subdistrict_en = $5", cipherNameToken, encryptedData.CipherZipcode, encryptedData.CipherProvince, encryptedData.CipherDistrict, encryptedData.CipherSub_district).Scan(&localAddressUUID)
	// fmt.Println(localAddressUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = tx.QueryRow("INSERT INTO local_address (lca_country, lca_zipcode, lca_province_en, lca_district_en, lca_subdistrict_en, lca_timezone, lca_createlocation) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING lca_id",
				// ctyNameToken, encryptedData.CipherZipcode, encryptedData.CipherProvince, encryptedData.CipherDistrict, encryptedData.CipherSub_district).Scan(&localAddressUUID) do go routine so need to change ctyNameToken into cipherNameToken
				cipherNameToken, encryptedData.CipherZipcode, encryptedData.CipherProvince, encryptedData.CipherDistrict, encryptedData.CipherSub_district, cipherTimezone, encryptedData.CipherCreate_Location).Scan(&localAddressUUID)

			if err != nil {
				log.Println(err)
				return models.RegisterResponses{}, err
				// return err
			}
		} else {
			log.Println(err)
			return models.RegisterResponses{}, err
			// return err
		}
	}
	// fmt.Println(localAddressUUID)
	var organizeMasterUUID string
	err = r.db.QueryRow("SELECT org_id FROM organize_master WHERE org_name_en = $1", encryptedData.CipherCompany_name_en).Scan(&organizeMasterUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = tx.QueryRow("INSERT INTO organize_master (org_logo, org_country, org_language, org_name_en, org_alias, org_currency, org_domain, org_createlocation) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING org_id",
				encryptedData.CipherUrl_logo, alpha2Code, languageNativeName, encryptedData.CipherCompany_name_en, encryptedData.CipherCompany_alias, currencyName, encryptedData.CipherCompany_domain, encryptedData.CipherCreate_Location).Scan(&organizeMasterUUID)
			if err != nil {
				log.Printf("Failed to insert data into organize_master table and return id. Error:%v\n", err)
				return models.RegisterResponses{}, err
			}
		} else {
			log.Printf("Failed to check for existing corporate entry. Error:%v", err)
			return models.RegisterResponses{}, err
		}
	}
	// fmt.Println(organizeMasterUUID)
	var organizeMemberUUID string
	err = r.db.QueryRow("SELECT orgmb_id FROM organize_member WHERE orgmb_email = $1", encryptedData.CipherUsername).Scan(&organizeMemberUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = tx.QueryRow("INSERT INTO organize_member (orgmb_holder, orgmb_title, orgmb_email, orgmb_name, orgmb_surname, orgmb_mobile, orgmb_createlocation) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING orgmb_id", organizeMasterUUID, encryptedData.CipherTitle,
				encryptedData.CipherUsername, encryptedData.CipherFirstname_en, encryptedData.CipherSurname_en, encryptedData.CipherMobile_phone, encryptedData.CipherCreate_Location).Scan(&organizeMemberUUID)
			if err != nil {
				log.Printf("Failed to insert data into organize_member table and return id. Error:%v\n", err)
				return models.RegisterResponses{}, err
			}
		} else {
			log.Printf("Failed to check for existing oganize_member entry. Error:%v", err)
			return models.RegisterResponses{}, err
		}
	}
	const requires_action = "change_password"
	var organizeMemberCredentialUUID string
	err = r.db.QueryRow("SELECT orgmbcr_id FROM organize_member_credential WHERE orgmbcr_orgmb_id = $1", organizeMemberUUID).Scan(&organizeMemberCredentialUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = tx.QueryRow("INSERT INTO organize_member_credential (orgmbcr_orgmb_id, orgmbcr_password, orgmbcr_level, orgmbcr_reqaction, orgmbcr_creator, orgmbcr_createlocation) VALUES ($1, $2, $3, $4, $5, $6) RETURNING orgmbcr_id",
				organizeMemberUUID, encryptedData.HashPassword, loginData.Role, requires_action, organizeMemberUUID, encryptedData.CipherCreate_Location).Scan(&organizeMemberCredentialUUID)
			if err != nil {
				log.Printf("Failed to insert data into organize_member_credential table and return id. Error:%v\n", err)
				return models.RegisterResponses{}, err
			}
		} else {
			log.Printf("Failed to check existing organize_member_credential entry. Error:%v\n", err)
			return models.RegisterResponses{}, err
		}
	}
	err = r.db.QueryRow("SELECT org_id FROM organize_master WHERE org_id = $1", organizeMasterUUID).Scan(&organizeMasterUUID) //// important use update in case need update not insert
	if err != nil {
		_, err = tx.Exec("UPDATE organize_master SET org_creator = $1 WHERE org_id = $2", organizeMemberUUID, organizeMasterUUID)
		if err != nil {
			log.Printf("Failed to update creator in organize_master. Error:%v\n", err)
			return models.RegisterResponses{}, err
		}
	}

	var organizeLocationUUID string
	err = r.db.QueryRow("SELECT orglo_id FROM organize_location WHERE orglo_org_id = $1", organizeMasterUUID).Scan(&organizeLocationUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = tx.QueryRow("INSERT INTO organize_location (orglo_org_id, orglo_addrno, orglo_addr1_en, orglo_lca_id, orglo_geolo, orglo_phone, orglo_createlocation, orglo_creator) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING orglo_id",
				organizeMasterUUID, encryptedData.CipherAddress_no, encryptedData.CipherAddress1_en, localAddressUUID, encryptedData.CipherCompany_geolo, encryptedData.CipherCompany_mobile, encryptedData.CipherCreate_Location, organizeMemberUUID).Scan(&organizeLocationUUID)
			if err != nil {
				log.Printf("Failed to insert data into organize_location table and return id. Error:%v\n", err)
				return models.RegisterResponses{}, err
			}
		} else {
			log.Printf("Failed to check for existing organize_location entry. Error:%v\n", err)
			return models.RegisterResponses{}, err
		}
	}

	var organizeRoleUUID string
	err = r.db.QueryRow("SELECT orgrl_id FROM organize_role WHERE orgrl_name_en = $1", encryptedData.CipherJob_title).Scan(&organizeRoleUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = tx.QueryRow("INSERT INTO organize_role (orgrl_name_en, orgrl_creator, orgrl_createlocation) VALUES ($1, $2, $3) RETURNING orgrl_id", encryptedData.CipherJob_title, organizeMemberUUID, encryptedData.CipherCreate_Location).Scan(&organizeRoleUUID)
			if err != nil {
				log.Printf("Failed to insert data into organize_role table and return id. Error:%v\n", err)
				return models.RegisterResponses{}, err
			}
		} else {
			log.Printf("Failed to check for existing organize_role entry. Error:%v\n", err)
			return models.RegisterResponses{}, err
		}
	}
	err = r.db.QueryRow("SELECT orgmb_id FROM organize_member WHERE orgmb_id = $1", organizeMemberUUID).Scan(&organizeMemberUUID)
	if err != nil {
		_, err = tx.Exec("UPDATE organize_member SET orgmb_role= $1, orgmb_workplace = $2, orgmb_creator = $3 WHERE orgmb_id = $4", organizeRoleUUID, organizeLocationUUID, organizeMemberUUID, organizeMemberUUID)
		if err != nil {
			log.Printf("Failed to update creator in organize_master. Error:%v\n", err)
			return models.RegisterResponses{}, err
		}
	}
	registerResponses.CompanyID = organizeMasterUUID
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return models.RegisterResponses{}, err
	}
	return registerResponses, nil
}

func (r *repositoryAdapter) ValidateDomainChicCRMSRepositoris(validateDomainRequest models.ValidateDomainRequest, cipherUsernameSuffix, cipherUsernameToken string) (models.ValidateDomainResponse, error) {
	var validateDomainResponse models.ValidateDomainResponse
	var UsernameExists string
	err := r.db.QueryRow("SELECT orgmb_email FROM organize_member WHERE orgmb_email = $1", cipherUsernameToken).Scan(&UsernameExists)
	if err == nil {
		return models.ValidateDomainResponse{}, errors.New("username already exists")
	} else if err != sql.ErrNoRows {
		return models.ValidateDomainResponse{}, err
	}
	// fmt.Println(cipherUsernameToken, cipherUsernameSuffix)
	err = r.db.QueryRow(`SELECT
	org.org_id,
	org.org_logo,
	org.org_country,
	org.org_language,
	org.org_name_en,
	org.org_alias,
	org.org_currency,
	org.org_domain,
	orglo.orglo_addrno,
	orglo.orglo_addr1_en,
	orglo.orglo_geolo,
	orglo.orglo_phone,
	lca.lca_zipcode,
	lca.lca_province_en,
	lca.lca_district_en,
	lca.lca_subdistrict_en,
	lca.lca_timezone,
	lca.lca_country
	FROM 
		organize_master org
	LEFT JOIN
		organize_location orglo ON org.org_id = orglo.orglo_org_id
	LEFT JOIN
		local_address lca ON orglo.orglo_lca_id = lca.lca_id
	WHERE
		org.org_domain = $1;
	`, cipherUsernameSuffix).Scan(&validateDomainResponse.CompanyID, &validateDomainResponse.CompanyLogo, &validateDomainResponse.Alpha2Code, &validateDomainResponse.LanguageName,
		&validateDomainResponse.Companyname, &validateDomainResponse.CompanyAlias, &validateDomainResponse.Currency, &validateDomainResponse.Domain, &validateDomainResponse.AddressNo,
		&validateDomainResponse.Address1En, &validateDomainResponse.Geolocation, &validateDomainResponse.CompanyPhone, &validateDomainResponse.Zipcode, &validateDomainResponse.Province, &validateDomainResponse.District,
		&validateDomainResponse.SubDistrict, &validateDomainResponse.Timezone, &validateDomainResponse.Country)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Domain does not match")
			return models.ValidateDomainResponse{}, errors.New("domain does not match. To proceed, please check your email")
		} else {
			log.Printf("Database query error: %v", err)
			return models.ValidateDomainResponse{}, err
		}
	}
	return validateDomainResponse, nil
}
