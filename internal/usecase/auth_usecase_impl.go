package usecase

import (
	"context"
	"errors"
	"tablelink_test/internal/repository"
	"tablelink_test/pkg"
	"time"
)

type AuthUsecaseImpl struct {
	userRepo    repository.UserRepository
	sessionRepo *repository.SessionRepositoryRedis
	jwtManager  *pkg.JWTManager
}

func NewAuthUsecase(userRepo repository.UserRepository, sessionRepo *repository.SessionRepositoryRedis, jwtManager *pkg.JWTManager) *AuthUsecaseImpl {
	return &AuthUsecaseImpl{userRepo, sessionRepo, jwtManager}
}

func (u *AuthUsecaseImpl) Login(email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if !pkg.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}
	token, err := u.jwtManager.Generate(user.ID, user.RoleID)
	if err != nil {
		return "", err
	}
	err = u.sessionRepo.SetSession(context.Background(), token, user.ID, 24*time.Hour)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *AuthUsecaseImpl) Logout(token string) error {
	return u.sessionRepo.DeleteSession(context.Background(), token)
}
