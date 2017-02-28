package services

import (
	"models"
	"log"
	"repository"
	"models/roles"
)

type RoleService struct {
	repo *repository.Repository
}

func NewRoleService() *RoleService {
	return &RoleService{repo:repository.NewRepo()}
}

func (s *RoleService) GetRole(r roles.Role) (*models.Role, error) {
	var role models.Role
	if err := s.repo.FindOne(&models.Role{Name: r.String()}, &role); err != nil {
		log.Printf("Error finding Role %s : %v", r.String(), err)
		return nil, err
	}

	return &role, nil
}
