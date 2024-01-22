package services

import (
	"chicCRM/modules/migrateData/models"
	"chicCRM/modules/migrateData/repositories"
	"errors"
	"fmt"
)

type ServicePort interface {
	MigrateMemberDataByOrganizeServices(organize models.Organize) error
	MigrateTeamleadTracByOrganizeServies(organize models.Organize) error
}

type serviceAdapter struct {
	r repositories.RepositoryPort
}

func NewServiceAdapter(r repositories.RepositoryPort) ServicePort {
	return &serviceAdapter{r: r}
}

func (s *serviceAdapter) MigrateMemberDataByOrganizeServices(organize models.Organize) error {
	if organize.OrganizeID == "" {
		return errors.New("please insert valid organize ID")
	}
	err := s.r.MigrateMemberDataByOrganizeRepositories(organize)
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceAdapter) MigrateTeamleadTracByOrganizeServies(organize models.Organize) error {
	fmt.Println(organize)
	if organize.OrganizeID == "" {
		return errors.New("please insert valid organize ID")
	}
	err := s.r.MigrateTeamleadTracByOrganizeRepositories(organize)
	if err != nil {
		return err
	}
	return nil
}
