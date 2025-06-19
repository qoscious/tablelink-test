package repository

import "tablelink_test/internal/domain"

type RoleRightRepository interface {
	GetByRoleAndRoute(roleID, section, route string) (*domain.RoleRight, error)
}
