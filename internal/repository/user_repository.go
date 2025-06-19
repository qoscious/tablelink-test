package repository

import "tablelink_test/internal/domain"

type UserRepository interface {
	GetByEmail(email string) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id string) error
}
