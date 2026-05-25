package grpc

import (
	"context"
	"math"
	"time"

	"go-api/internal/repository"
	usersvc "go-api/internal/service"
	"go-api/pkg/constants"
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
	userv1 "github.com/lamquangmanh/protobuf/gen/go/proto/user/v1"
	"google.golang.org/grpc/codes"
)

// UserHandler is the thin gRPC transport layer for user operations.
// It maps proto request/response types to/from the UserService business layer.
type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	userService *usersvc.UserService
}

// NewUserHandler creates a new UserHandler wrapping the given UserService.
func NewUserHandler(userService *usersvc.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// userRepoToProto converts a repository.User to a userpb.User protobuf message.
func userRepoToProto(u *repository.User) *userv1.User {
	phone := ""
	if u.Phone.Valid {
		phone = u.Phone.String
	}
	avatar := ""
	if u.Avatar.Valid {
		avatar = u.Avatar.String
	}
	createdUserID := ""
	if u.CreatedUserID.Valid {
		createdUserID = u.CreatedUserID.String
	}
	updatedAt := ""
	if u.UpdatedAt.Valid {
		updatedAt = u.UpdatedAt.Time.UTC().Format(time.RFC3339)
	}
	updatedUserID := ""
	if u.UpdatedUserID.Valid {
		updatedUserID = u.UpdatedUserID.String
	}
	deletedAt := ""
	if u.DeletedAt.Valid {
		deletedAt = u.DeletedAt.Time.UTC().Format(time.RFC3339)
	}
	deletedUserID := ""
	if u.DeletedUserID.Valid {
		deletedUserID = u.DeletedUserID.String
	}
	statusValue := userv1.UserStatus_USER_STATUS_UNSPECIFIED
	switch u.Status {
	case repository.UserStatusACTIVE:
		statusValue = userv1.UserStatus_USER_STATUS_ACTIVE
	case repository.UserStatusDEACTIVATED:
		statusValue = userv1.UserStatus_USER_STATUS_DEACTIVATED
	case repository.UserStatusDELETED:
		statusValue = userv1.UserStatus_USER_STATUS_DELETED
	}

	return &userv1.User{
		UserId:        u.UserID.String(),
		Username:      u.UserName,
		Email:         u.Email,
		Phone:         phone,
		Avatar:        avatar,
		Status:        statusValue,
		CreatedAt:     u.CreatedAt.Time.UTC().Format(time.RFC3339),
		CreatedUserId: createdUserID,
		UpdatedAt:     updatedAt,
		UpdatedUserId: updatedUserID,
		DeletedAt:     deletedAt,
		DeletedUserId: deletedUserID,
	}
}

// CreateUser delegates to UserService and maps validation errors into the
// response body (soft errors), returning a gRPC error only for system failures.
func (s *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if req.GetUser() == nil {
		logger.Error("CreateUser request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUserPayloadRequired)})
		return nil, err
	}
	u, errors := s.userService.CreateUser(ctx, usersvc.CreateUserInput{
		Username:    req.GetUser().GetUsername(),
		Email:       req.GetUser().GetEmail(),
		Password:    req.GetUser().GetPassword(),
		Phone:       req.GetUser().GetPhone(),
		Avatar:      req.GetUser().GetAvatar(),
		RoleIDs:     req.GetUser().GetRoleIds(),
		Status:      req.GetUser().GetStatus(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		logger.Error("CreateUser request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrCreateUserInternal.Messagef(errors))})
		return nil, err
	}

	return &userv1.CreateUserResponse{User: userRepoToProto(u)}, nil
}

// GetUser delegates to UserService and maps the result to a User proto message.
func (s *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	u, errors := s.userService.GetUser(ctx, req.GetUserId())
	if errors != nil {
		logger.Error("GetUser request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetUserInternal.Messagef(errors))})
		return nil, err
	}
	roles, _, _, _, err := s.userService.ListUserRoles(ctx, req.GetUserId(), 1000, 0)
	if err == nil {
		roleIDs := make([]string, 0, len(roles))
		for _, role := range roles {
			roleIDs = append(roleIDs, role.RoleID.String())
		}
		protoUser := userRepoToProto(u)
		protoUser.RoleIds = roleIDs
		return &userv1.GetUserResponse{User: protoUser}, nil
	}
	return &userv1.GetUserResponse{User: userRepoToProto(u)}, nil
}

// UpdateUser delegates to UserService and maps the updated user to a User proto message.
func (s *UserHandler) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	if req.GetUser() == nil {
		logger.Error("UpdateUser request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUserPayloadRequired)})
		return nil, err
	}
	_, errors := s.userService.UpdateUser(ctx, usersvc.UpdateUserInput{
		ID:          req.GetUser().GetUserId(),
		Username:    req.GetUser().GetUsername(),
		Email:       req.GetUser().GetEmail(),
		Phone:       req.GetUser().GetPhone(),
		Avatar:      req.GetUser().GetAvatar(),
		Status:      req.GetUser().GetStatus(),
		RoleIDs:     req.GetUser().GetRoleIds(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		logger.Error("UpdateUser request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUpdateUserInternal.Messagef(errors))})
		return nil, err
	}
	return &userv1.UpdateUserResponse{Success: true}, nil
}

// DeleteUser delegates to UserService and returns a DeleteUserResponse.
func (s *UserHandler) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	if errors := s.userService.DeleteUser(ctx, req.GetUserId(), req.GetDeletedUserId()); errors != nil {
		logger.Error("DeleteUser request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrDeleteUserInternal.Messagef(errors))})
		return nil, err
	}
	return &userv1.DeleteUserResponse{Success: true}, nil
}

// ListUsers delegates to UserService and maps the results to a ListUsersResponse.
func (s *UserHandler) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	limit := int32(20)
	page := int32(1)
	if p := req.GetPagination(); p != nil {
		if p.GetLimit() > 0 {
			limit = p.GetLimit()
		}
		if p.GetPage() > 0 {
			page = p.GetPage()
		}
	}
	offset := (page - 1) * limit

	sorts := make([]*basev1.Sort, 0, len(req.GetSorts()))
	sorts = append(sorts, req.GetSorts()...)
	filters := make([]*basev1.Filter, 0, len(req.GetFilters()))
	filters = append(filters, req.GetFilters()...)

	rows, total, appliedLimit, appliedOffset, err := s.userService.ListUsers(ctx, usersvc.ListUsersInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   sorts,
		Filters: filters,
	})
	if err != nil {
		logger.Error("ListUsers request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrListUsersInternal.Messagef(err))})
		return nil, err
	}

	data := make([]*userv1.User, 0, len(rows))
	for _, row := range rows {
		data = append(data, userRepoToProto(row))
	}
	totalItems := int32(total)
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &userv1.GetUsersResponse{
		Data: data,
		Pagination: &basev1.PaginationResponse{
			Page:       (appliedOffset / appliedLimit) + 1,
			Limit:      appliedLimit,
			TotalItems: totalItems,
			TotalPages: totalPages,
			ItemCount:  int32(len(data)),
		},
	}, nil
}

func (s *UserHandler) ChangePassword(ctx context.Context, req *userv1.ChangePasswordRequest) (*userv1.ChangePasswordResponse, error) {
	if errors := s.userService.ChangePassword(ctx, req.GetUserId(), req.GetPassword()); errors != nil {
		logger.Error("ChangePassword request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &userv1.ChangePasswordResponse{Success: true}, nil
}
