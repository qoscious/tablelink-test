package grpc

import (
	"context"
	"tablelink_test/internal/domain"
	"tablelink_test/internal/usecase"
	authpb "tablelink_test/proto/auth"
	userpb "tablelink_test/proto/user"
)

type AuthServiceServer struct {
	authpb.UnimplementedAuthServiceServer
	AuthUC usecase.AuthUsecase
}

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	UserUC usecase.UserUsecase
}

// AuthService
func (s *AuthServiceServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	token, err := s.AuthUC.Login(req.Email, req.Password)
	if err != nil {
		return &authpb.LoginResponse{
			Status:  false,
			Message: "Login failed: " + err.Error(),
		}, nil
	}
	return &authpb.LoginResponse{
		Status:  true,
		Message: "Successfully",
		Data:    &authpb.LoginResponse_Data{AccessToken: token},
	}, nil
}

func (s *AuthServiceServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	err := s.AuthUC.Logout(req.AccessToken)
	if err != nil {
		return &authpb.LogoutResponse{Status: false, Message: "Logout failed: " + err.Error()}, nil
	}
	return &authpb.LogoutResponse{Status: true, Message: "Successfully"}, nil
}

// UserService
func (s *UserServiceServer) GetAllUser(ctx context.Context, req *userpb.GetAllUserRequest) (*userpb.GetAllUserResponse, error) {
	token, section := extractTokenAndSection(ctx)
	users, err := s.UserUC.GetAll(token, section, "/users/user")
	if err != nil {
		return &userpb.GetAllUserResponse{Status: false, Message: err.Error()}, nil
	}
	var pbUsers []*userpb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &userpb.User{
			RoleId:     u.RoleID,
			RoleName:   "", // Optionally fetch role name
			Name:       u.Name,
			Email:      u.Email,
			LastAccess: u.LastAccess,
		})
	}
	return &userpb.GetAllUserResponse{
		Status:  true,
		Message: "Successfully",
		Data:    &userpb.GetAllUserResponse_Data{User: pbUsers},
	}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user := &domain.User{
		RoleID:   req.RoleId,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	token, section := extractTokenAndSection(ctx)
	err := s.UserUC.Create(user, token, section, "/users/user")
	if err != nil {
		return &userpb.CreateUserResponse{Status: false, Message: err.Error()}, nil
	}
	return &userpb.CreateUserResponse{Status: true, Message: "Successfully"}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	user := &domain.User{
		ID:   req.UserId,
		Name: req.Name,
	}
	token, section := extractTokenAndSection(ctx)
	err := s.UserUC.Update(user, token, section, "/users/user")
	if err != nil {
		return &userpb.UpdateUserResponse{Status: false, Message: err.Error()}, nil
	}
	return &userpb.UpdateUserResponse{Status: true, Message: "Successfully"}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	token, section := extractTokenAndSection(ctx)
	err := s.UserUC.Delete(req.UserId, token, section, "/users/user")
	if err != nil {
		return &userpb.DeleteUserResponse{Status: false, Message: err.Error()}, nil
	}
	return &userpb.DeleteUserResponse{Status: true, Message: "Successfully"}, nil
}

// Helper untuk ambil token dan section dari metadata
func extractTokenAndSection(ctx context.Context) (string, string) {
	md, _ := getMetadata(ctx)
	token := ""
	section := ""
	if vals, ok := md["authorization"]; ok && len(vals) > 0 {
		token = vals[0]
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}
	if vals, ok := md["x-link-service"]; ok && len(vals) > 0 {
		section = vals[0]
	}
	return token, section
}

func getMetadata(ctx context.Context) (map[string][]string, bool) {
	md, ok := ctx.Value("grpc.metadata").(map[string][]string)
	if !ok {
		return map[string][]string{}, false
	}
	return md, true
}
