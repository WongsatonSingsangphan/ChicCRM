package repositories

import (
	"chicCRM/modules/migrateData/models"
	"database/sql"
	"errors"
	"fmt"
)

type RepositoryPort interface {
	MigrateMemberDataByOrganizeRepositories(organize models.Organize) error
	MigrateTeamleadTracByOrganizeRepositories(organize models.Organize) error
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}

func (r *repositoryAdapter) MigrateMemberDataByOrganizeRepositories(organize models.Organize) error {
	var organizeUUID string
	err := r.db.QueryRow("SELECT org_id FROM organize_master WHERE org_id = $1", organize.OrganizeID).Scan(&organizeUUID)
	if err != nil {
		fmt.Println(err)
		return errors.New("organize_id not found")
	}
	rows, err := r.db.Query("SELECT orgmb_id FROM organize_member WHERE orgmb_holder = $1", organizeUUID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for rows.Next() {
		var orgmbID string
		if err := rows.Scan(&orgmbID); err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		_, err := tx.Exec("INSERT INTO organize_member_authorization (orgmbat_orgmbid, orgmbat_orgid) VALUES ($1, $2) ON CONFLICT (orgmbat_orgmbid) DO NOTHING", orgmbID, organizeUUID)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *repositoryAdapter) MigrateTeamleadTracByOrganizeRepositories(organize models.Organize) error {
	var organizeUUID string
	err := r.db.QueryRow("SELECT orgtl_org_id FROM organize_teamlead_trac WHERE orgtl_org_id = $1", organize.OrganizeID).Scan(&organizeUUID)
	if err != nil {
		fmt.Println(err)
		return errors.New("organize_id not found")
	}
	fmt.Println(organize)
	rows, err := r.db.Query("SELECT orgtl_trac_id, orgtl_level FROM organize_teamlead_trac WHERE orgtl_org_id = $1", organizeUUID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(organizeUUID)
	defer rows.Close()
	tx, err := r.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for rows.Next() {
		var orgmbID, orgmbLV string
		if err := rows.Scan(&orgmbID, &orgmbLV); err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
		_, err := tx.Exec("INSERT INTO organize_member_authorization (orgmbat_orgmbid, orgmbat_orgid, orgmbat_level) VALUES ($1, $2, $3) ON CONFLICT (orgmbat_orgmbid) DO NOTHING", orgmbID, organizeUUID, orgmbLV)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
