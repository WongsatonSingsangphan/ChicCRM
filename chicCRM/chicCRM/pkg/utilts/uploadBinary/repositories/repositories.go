package repositories

import (
	"chicCRM/pkg/utilts/uploadBinary/models"
	"database/sql"
	"errors"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/lib/pq"
)

type RepositoryPort interface {
	UploadBinaryChicCRMSRepositoris(requestOrganizeBinary models.RequestOrganizeBinary, file []*multipart.FileHeader) error
	GetBinaryChicCRMRepositories(organizeID string) ([]byte, string, error)
	EditLogoProfileChiCCRMRepositories(editLogoRequest models.EditLogoRequest, file []*multipart.FileHeader) error
	EditPersonalProfileChicCRMRepositories(editPersonalProfileRequest models.EditPersonalProfileRequest, file []*multipart.FileHeader) error
	GetPersonalProfileChicCRMRepositiries(personalID string) ([]byte, string, error)
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}

func (r *repositoryAdapter) UploadBinaryChicCRMSRepositoris(requestOrganizeBinary models.RequestOrganizeBinary, file []*multipart.FileHeader) error {
	var orgLogoBinary []byte
	err := r.db.QueryRow("SELECT org_logo_binary FROM organize_master WHERE org_id = $1", requestOrganizeBinary.OrganizeID).Scan(&orgLogoBinary)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("organize_master UUID not match")
		}
		return err
	}
	if orgLogoBinary == nil {
		_, err = r.db.Exec(`UPDATE organize_master SET org_logo_binary = $2 WHERE org_id = $1`, requestOrganizeBinary.OrganizeID, pq.Array(requestOrganizeBinary.FileBytes))
		if err != nil {
			log.Printf("Failed to update org_logo_binary: %v\n", err)
			return err
		}
	} else {
		log.Printf("org_logo_binary already exists for org_id: %s\n", requestOrganizeBinary.OrganizeID)
	}

	return nil
}
func (r *repositoryAdapter) GetBinaryChicCRMRepositories(organizeID string) ([]byte, string, error) {
	var uploadedFiles [][]byte
	err := r.db.QueryRow("SELECT org_logo_binary FROM organize_master WHERE org_id = $1", organizeID).Scan(pq.Array(&uploadedFiles))
	if err != nil {
		return nil, "", err
	}
	if len(uploadedFiles) == 0 || len(uploadedFiles[0]) == 0 {
		return nil, "", sql.ErrNoRows
	}
	fileType := http.DetectContentType(uploadedFiles[0])
	return uploadedFiles[0], fileType, nil
}

func (r *repositoryAdapter) EditLogoProfileChiCCRMRepositories(editLogoRequest models.EditLogoRequest, file []*multipart.FileHeader) error {
	_, err := r.db.Exec("UPDATE organize_master SET org_logo_binary = $2 WHERE org_id = $1", editLogoRequest.OrganizeID, pq.Array(editLogoRequest.FileBytes))
	if err != nil {
		log.Printf("Failed to update org_logo_binary: %v\n", err)
		return err
	}
	return err
}

func (r *repositoryAdapter) EditPersonalProfileChicCRMRepositories(editPersonalProfileRequest models.EditPersonalProfileRequest, file []*multipart.FileHeader) error {
	_, err := r.db.Exec("UPDATE organize_member SET orgmb_profile = $2 WHERE orgmb_id = $1", editPersonalProfileRequest.PersonalID, pq.Array(editPersonalProfileRequest.FileBytes))
	if err != nil {
		log.Printf("Failed to update org_logo_binary: %v\n", err)
		return err
	}
	return err
}

func (r *repositoryAdapter) GetPersonalProfileChicCRMRepositiries(personalID string) ([]byte, string, error) {
	var uploadedFiles [][]byte
	// fmt.Println(uploadedFiles)
	// fmt.Println(personalID)
	err := r.db.QueryRow("SELECT orgmb_profile FROM organize_member WHERE orgmb_id = $1", personalID).Scan(pq.Array(&uploadedFiles))
	if err != nil {
		return nil, "", errors.New("personalID does not match")
	}
	if len(uploadedFiles) == 0 || len(uploadedFiles[0]) == 0 {
		return nil, "", sql.ErrNoRows
	}
	fileType := http.DetectContentType(uploadedFiles[0])
	return uploadedFiles[0], fileType, nil
}
