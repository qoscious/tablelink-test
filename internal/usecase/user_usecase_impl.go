package usecase

import (
	"context"
	"errors"
	"tablelink_test/internal/domain"
	"tablelink_test/internal/repository"
	"tablelink_test/pkg"
)

type UserUsecaseImpl struct {
	userRepo      repository.UserRepository
	roleRepo      repository.RoleRepository
	roleRightRepo repository.RoleRightRepository
	sessionRepo   *repository.SessionRepositoryRedis
	jwtManager    *pkg.JWTManager
}

func NewUserUsecase(userRepo repository.UserRepository, roleRepo repository.RoleRepository, roleRightRepo repository.RoleRightRepository, sessionRepo *repository.SessionRepositoryRedis, jwtManager *pkg.JWTManager) *UserUsecaseImpl {
	return &UserUsecaseImpl{userRepo, roleRepo, roleRightRepo, sessionRepo, jwtManager}
}

func (u *UserUsecaseImpl) validateAccess(token, section, route string, action string) (string, error) {
	claims, err := u.jwtManager.Verify(token)
	if err != nil {
		return "", errors.New("invalid token")
	}
	_, err = u.sessionRepo.CheckSession(context.Background(), token)
	if err != nil {
		return "", errors.New("session expired")
	}
	rr, err := u.roleRightRepo.GetByRoleAndRoute(claims.RoleID, section, route)
	if err != nil {
		return "", errors.New("no access")
	}
	switch action {
	case "read":
		if rr.RRead != 1 {
			return "", errors.New("no read access")
		}
	case "create":
		if rr.RCreate != 1 {
			return "", errors.New("no create access")
		}
	case "update":
		if rr.RUpdate != 1 {
			return "", errors.New("no update access")
		}
	case "delete":
		if rr.RDelete != 1 {
			return "", errors.New("no delete access")
		}
	}
	return claims.UserID, nil
}

func (u *UserUsecaseImpl) GetAll(token, section, route string) ([]*domain.User, error) {
	_, err := u.validateAccess(token, section, route, "read")
	if err != nil {
		return nil, err
	}
	return u.userRepo.GetAll()
}

func (u *UserUsecaseImpl) Create(user *domain.User, token, section, route string) error {
	_, err := u.validateAccess(token, section, route, "create")
	if err != nil {
		return err
	}
	hash, err := pkg.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return u.userRepo.Create(user)
}

func (u *UserUsecaseImpl) Update(user *domain.User, token, section, route string) error {
	_, err := u.validateAccess(token, section, route, "update")
	if err != nil {
		return err
	}
	return u.userRepo.Update(user)
}

func (u *UserUsecaseImpl) Delete(userID, token, section, route string) error {
	_, err := u.validateAccess(token, section, route, "delete")
	if err != nil {
		return err
	}
	return u.userRepo.Delete(userID)
}
