package repository

import "tablelink_test/internal/domain"

type RoleRepository interface {
	GetByID(id string) (*domain.Role, error)
	GetByName(name string) (*domain.Role, error)
}
