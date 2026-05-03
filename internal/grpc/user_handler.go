package grpc

import (
	"context"
	"math"
	"time"

	"go-api/internal/repository"
	usersvc "go-api/internal/service"
	basepb "go-api/pkg/api/basepb"
	userpb "go-api/pkg/api/userpb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"
)

// UserHandler is the thin gRPC transport layer for user operations.
// It maps proto request/response types to/from the UserService business layer.
type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userService *usersvc.UserService
}

// NewUserHandler creates a new UserHandler wrapping the given UserService.
func NewUserHandler(userService *usersvc.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// userRepoToProto converts a repository.User to a userpb.User protobuf message.
func userRepoToProto(u *repository.User) *userpb.User {
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
	statusValue := userpb.UserStatus_USER_STATUS_UNSPECIFIED
	switch u.Status {
	case repository.UserStatusACTIVE:
		statusValue = userpb.UserStatus_USER_STATUS_ACTIVE
	case repository.UserStatusDEACTIVATED:
		statusValue = userpb.UserStatus_USER_STATUS_DEACTIVATED
	case repository.UserStatusDELETED:
		statusValue = userpb.UserStatus_USER_STATUS_DELETED
	}

	return &userpb.User{
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
func (s *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateSuccess, error) {
	if req.GetUser() == nil {
		return &userpb.CreateSuccess{Errors: []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUserPayloadRequired)}}, nil
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
		return &userpb.CreateSuccess{Errors: errors}, nil
	}

	return &userpb.CreateSuccess{User: userRepoToProto(u)}, nil
}

// GetUser delegates to UserService and maps the result to a User proto message.
func (s *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	u, errors := s.userService.GetUser(ctx, req.GetUserId())
	if errors != nil {
		return &userpb.GetUserResponse{Errors: errors}, nil
	}
	roles, _, _, _, err := s.userService.ListUserRoles(ctx, req.GetUserId(), 1000, 0)
	if err == nil {
		roleIDs := make([]string, 0, len(roles))
		for _, role := range roles {
			roleIDs = append(roleIDs, role.RoleID.String())
		}
		protoUser := userRepoToProto(u)
		protoUser.RoleIds = roleIDs
		return &userpb.GetUserResponse{User: protoUser}, nil
	}
	return &userpb.GetUserResponse{User: userRepoToProto(u)}, nil
}

// UpdateUser delegates to UserService and maps the updated user to a User proto message.
func (s *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*basepb.UpdateSuccess, error) {
	if req.GetUser() == nil {
		return utils.UpdateErr(utils.ErrMsg(constants.ErrUserPayloadRequired)), nil
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
		return utils.UpdateErr(errors...), nil
	}
	return utils.UpdateOK(), nil
}

// DeleteUser delegates to UserService and returns a DeleteUserResponse.
func (s *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*basepb.DeleteSuccess, error) {
	if errors := s.userService.DeleteUser(ctx, req.GetUserId(), req.GetDeletedUserId()); errors != nil {
		return utils.DeleteErr(errors...), nil
	}
	return utils.DeleteOK(), nil
}

// ListUsers delegates to UserService and maps the results to a ListUsersResponse.
func (s *UserHandler) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
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

	rows, total, appliedLimit, appliedOffset, err := s.userService.ListUsers(ctx, usersvc.ListUsersInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   req.GetSorts(),
		Filters: req.GetFilters(),
	})
	if err != nil {
		return &userpb.GetUsersResponse{Data: nil, Pagination: &basepb.PaginationResponse{}}, nil
	}

	data := make([]*userpb.User, 0, len(rows))
	for _, row := range rows {
		data = append(data, userRepoToProto(row))
	}
	totalItems := int32(total)
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &userpb.GetUsersResponse{
		Data: data,
		Pagination: &basepb.PaginationResponse{
			Page:       (appliedOffset / appliedLimit) + 1,
			Limit:      appliedLimit,
			TotalItems: totalItems,
			TotalPages: totalPages,
			ItemCount:  int32(len(data)),
		},
	}, nil
}

func (s *UserHandler) ChangePassword(ctx context.Context, req *userpb.ChangePasswordRequest) (*basepb.UpdateSuccess, error) {
	if errors := s.userService.ChangePassword(ctx, req.GetUserId(), req.GetPassword()); errors != nil {
		return utils.UpdateErr(errors...), nil
	}
	return utils.UpdateOK(), nil
}
