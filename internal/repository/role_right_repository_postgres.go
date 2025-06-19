package repository

import (
	"tablelink_test/internal/domain"

	"gorm.io/gorm"
)

type RoleRightRepositoryPostgres struct {
	db *gorm.DB
}

func NewRoleRightRepositoryPostgres(db *gorm.DB) *RoleRightRepositoryPostgres {
	return &RoleRightRepositoryPostgres{db}
}

func (r *RoleRightRepositoryPostgres) GetByRoleAndRoute(roleID, section, route string) (*domain.RoleRight, error) {
	var rr domain.RoleRight
	if err := r.db.Where("role_id = ? AND section = ? AND route = ?", roleID, section, route).First(&rr).Error; err != nil {
		return nil, err
	}
	return &rr, nil
}
