package usecase

import "tablelink_test/internal/domain"

type UserUsecase interface {
	GetAll(token, section, route string) ([]*domain.User, error)
	Create(user *domain.User, token, section, route string) error
	Update(user *domain.User, token, section, route string) error
	Delete(userID, token, section, route string) error
}
