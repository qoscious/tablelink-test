package usecase

type AuthUsecase interface {
	Login(email, password string) (string, error)
	Logout(token string) error
}
