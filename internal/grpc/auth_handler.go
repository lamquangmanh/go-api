package grpc

import (
	"context"

	"go-api/internal/service"
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	authv1 "github.com/lamquangmanh/protobuf/gen/go/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthHandler implements gRPC AuthService
type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user authentication
func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	// Basic validation
	if req == nil {
		logger.Error("Login request is nil")
		return nil, status.Error(codes.InvalidArgument, "login request is required")
	}

	// Delegate to auth service
	if h.authService == nil {
		logger.Error("Auth service is not configured")
		return nil, status.Error(codes.Internal, "auth service is not configured")
	}
	response, errors := h.authService.Login(ctx, service.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if errors != nil {
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return response, nil
}

// GetMe retrieves current user information
func (h *AuthHandler) GetMe(ctx context.Context, req *authv1.GetMeRequest) (*authv1.GetMeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "get me request is required")
	}
	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service is not configured")
	}

	user, errors := h.authService.GetMe(ctx, req.GetUserId())
	if errors != nil {
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}

	return &authv1.GetMeResponse{
		User:   user,
		Errors: nil,
	}, nil
}

// Verify validates authentication tokens
func (h *AuthHandler) Verify(ctx context.Context, req *authv1.VerifyRequest) (*authv1.VerifyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "verify request is required")
	}
	if h.authService == nil {
		return nil, status.Error(codes.Internal, "auth service is not configured")
	}

	_, errors := h.authService.Verify(ctx, req.GetToken())
	if errors != nil {
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}

	return &authv1.VerifyResponse{
		Success: true,
		Errors:  nil,
	}, nil
}
