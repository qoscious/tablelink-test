package repository

import (
	"tablelink_test/internal/domain"

	"gorm.io/gorm"
)

type RoleRepositoryPostgres struct {
	db *gorm.DB
}

func NewRoleRepositoryPostgres(db *gorm.DB) *RoleRepositoryPostgres {
	return &RoleRepositoryPostgres{db}
}

func (r *RoleRepositoryPostgres) GetByID(id string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepositoryPostgres) GetByName(name string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
