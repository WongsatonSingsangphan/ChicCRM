package services

import (
	"chicCRM/pkg/utilts/uploadBinary/models"
	"chicCRM/pkg/utilts/uploadBinary/repositories"
	"fmt"
	"io"
	"log"
	"mime/multipart"
)

type ServicePort interface {
	UploadBinaryChicCRMServices(requestOrganizeBinary models.RequestOrganizeBinary, file []*multipart.FileHeader) error
	GetBinaryChicCRMServices(organizeID string) ([]byte, string, error)
	EditLogoProfileChiCCRMServices(editLogoRequest models.EditLogoRequest, file []*multipart.FileHeader) error
	EditPersonalProfileChicCRMServices(editPersonalProfileRequest models.EditPersonalProfileRequest, file []*multipart.FileHeader) error
	GetPersonalProfileChicCRMServices(personalID string) ([]byte, string, error)
}

type serviceAdapter struct {
	r repositories.RepositoryPort
}

func NewServiceAdapter(r repositories.RepositoryPort) ServicePort {
	return &serviceAdapter{r: r}
}

func (s *serviceAdapter) UploadBinaryChicCRMServices(requestOrganizeBinary models.RequestOrganizeBinary, file []*multipart.FileHeader) error {
	var fileBytesArray [][]byte
	for _, file := range file {
		openedFile, err := file.Open()
		if err != nil {
			return err
		}
		fileBytes, err := io.ReadAll(openedFile)
		openedFile.Close()
		if err != nil {
			return err
		}
		fileBytesArray = append(fileBytesArray, fileBytes)
	}
	requestOrganizeBinary.FileBytes = fileBytesArray
	err := s.r.UploadBinaryChicCRMSRepositoris(requestOrganizeBinary, file)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *serviceAdapter) GetBinaryChicCRMServices(organizeID string) ([]byte, string, error) {
	fileBytes, fileType, err := s.r.GetBinaryChicCRMRepositories(organizeID)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	return fileBytes, fileType, nil
}

func (s *serviceAdapter) EditLogoProfileChiCCRMServices(editLogoRequest models.EditLogoRequest, file []*multipart.FileHeader) error {
	var fileBytesArray [][]byte
	for _, fileHeader := range file {
		openedFile, err := fileHeader.Open()
		if err != nil {
			log.Println(err)
			return fmt.Errorf("failed to open file: %w", err)
		}
		fileBytes, err := io.ReadAll(openedFile)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("failed to read file: %w", err)
		}
		openedFile.Close()
		fileBytesArray = append(fileBytesArray, fileBytes)
	}
	editLogoRequest.FileBytes = fileBytesArray
	if err := s.r.EditLogoProfileChiCCRMRepositories(editLogoRequest, file); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *serviceAdapter) EditPersonalProfileChicCRMServices(editPersonalProfileRequest models.EditPersonalProfileRequest, file []*multipart.FileHeader) error {
	var fileBytesArray [][]byte
	for _, fileHeader := range file {
		openedFile, err := fileHeader.Open()
		if err != nil {
			log.Println(err)
			return fmt.Errorf("failed to open file: %w", err)
		}
		fileBytes, err := io.ReadAll(openedFile)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("failed to read file: %w", err)
		}
		openedFile.Close()
		fileBytesArray = append(fileBytesArray, fileBytes)
	}
	editPersonalProfileRequest.FileBytes = fileBytesArray
	if err := s.r.EditPersonalProfileChicCRMRepositories(editPersonalProfileRequest, file); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *serviceAdapter) GetPersonalProfileChicCRMServices(personalID string) ([]byte, string, error) {
	fileBytes, fileType, err := s.r.GetPersonalProfileChicCRMRepositiries(personalID)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	return fileBytes, fileType, nil
}
