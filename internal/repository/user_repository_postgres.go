package repository

import (
	"tablelink_test/internal/domain"

	"gorm.io/gorm"
)

type UserRepositoryPostgres struct {
	db *gorm.DB
}

func NewUserRepositoryPostgres(db *gorm.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db}
}

func (r *UserRepositoryPostgres) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) GetByID(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) GetAll() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryPostgres) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepositoryPostgres) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepositoryPostgres) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}
