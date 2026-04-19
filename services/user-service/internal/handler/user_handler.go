package handler

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/split/services/user-service/internal/model"
	"github.com/split/services/user-service/internal/service"
	pb "github.com/split/services/user-service/proto"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Email == "" || req.Name == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email, name, and password are required")
	}

	user, accessToken, refreshToken, err := h.svc.Register(ctx, req.Email, req.Name, req.Password, req.Currency)
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		log.Printf("register error: %v", err)
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &pb.RegisterResponse{
		User:         userToProto(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	user, accessToken, refreshToken, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		log.Printf("login error: %v", err)
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &pb.LoginResponse{
		User:         userToProto(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	user, err := h.svc.GetUser(ctx, req.Id)
	if err != nil {
		if err == service.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Printf("get user error: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		User: userToProto(user),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	user, err := h.svc.UpdateUser(ctx, req.Id, req.Name, req.AvatarUrl, req.Currency)
	if err != nil {
		if err == service.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Printf("update user error: %v", err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &pb.UpdateUserResponse{
		User: userToProto(user),
	}, nil
}

func (h *UserHandler) GetUsersByIds(ctx context.Context, req *pb.GetUsersByIdsRequest) (*pb.GetUsersByIdsResponse, error) {
	users, err := h.svc.GetUsersByIDs(ctx, req.Ids)
	if err != nil {
		log.Printf("get users by ids error: %v", err)
		return nil, status.Error(codes.Internal, "failed to get users")
	}

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, userToProto(u))
	}

	return &pb.GetUsersByIdsResponse{
		Users: pbUsers,
	}, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	accessToken, refreshToken, err := h.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == service.ErrInvalidToken {
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		}
		log.Printf("refresh token error: %v", err)
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func userToProto(u *model.User) *pb.User {
	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		AvatarUrl: u.AvatarURL,
		Currency:  u.Currency,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
