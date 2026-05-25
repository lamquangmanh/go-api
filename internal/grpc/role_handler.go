package grpc

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	"go-api/internal/repository"
	rolesvc "go-api/internal/service"
	"go-api/pkg/constants"
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
	rolev1 "github.com/lamquangmanh/protobuf/gen/go/proto/role/v1"
)

type RoleHandler struct {
	rolev1.UnimplementedRoleServiceServer
	roleService *rolesvc.RoleService
}

// NewRoleHandler creates a thin gRPC handler that delegates business logic to RoleService.
func NewRoleHandler(roleService *rolesvc.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// roleRepoToProto converts repository.Role to protobuf Role response format.
func roleRepoToProto(role *repository.Role) *rolev1.Role {
	description := ""
	if role.Description.Valid {
		description = role.Description.String
	}
	moduleID := ""
	if role.ModuleID != uuid.Nil {
		moduleID = role.ModuleID.String()
	}
	createdUserID := ""
	if role.CreatedUserID.Valid {
		createdUserID = role.CreatedUserID.String
	}
	updatedAt := ""
	if role.UpdatedAt.Valid {
		updatedAt = role.UpdatedAt.Time.UTC().Format(time.RFC3339)
	}
	updatedUserID := ""
	if role.UpdatedUserID.Valid {
		updatedUserID = role.UpdatedUserID.String
	}
	deletedAt := ""
	if role.DeletedAt.Valid {
		deletedAt = role.DeletedAt.Time.UTC().Format(time.RFC3339)
	}
	deletedUserID := ""
	if role.DeletedUserID.Valid {
		deletedUserID = role.DeletedUserID.String
	}

	return &rolev1.Role{
		RoleId:        role.RoleID.String(),
		Name:          role.Name,
		Description:   description,
		ModuleId:      moduleID,
		CreatedUserId: createdUserID,
		UpdatedAt:     updatedAt,
		UpdatedUserId: updatedUserID,
		DeletedAt:     deletedAt,
		DeletedUserId: deletedUserID,
		CreatedAt:     role.CreatedAt.Time.UTC().Format(time.RFC3339),
	}
}

// CreateRole maps gRPC request to service input and returns created role.
func (s *RoleHandler) CreateRole(ctx context.Context, req *rolev1.CreateRoleRequest) (*rolev1.CreateRoleResponse, error) {
	if req.GetRole() == nil {
		logger.Error("CreateRole request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRolePayloadRequired)})
		return nil, err
	}
	role, errors := s.roleService.CreateRole(ctx, rolesvc.CreateRoleInput{
		Name:        req.GetRole().GetName(),
		Description: req.GetRole().GetDescription(),
		ModuleID:    req.GetRole().GetModuleId(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		logger.Error("CreateRole request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrCreateRoleInternal.Messagef(errors))})
		return nil, err
	}

	return &rolev1.CreateRoleResponse{Role: roleRepoToProto(role)}, nil
}

// GetRole fetches a role by ID and maps it into protobuf response.
func (s *RoleHandler) GetRole(ctx context.Context, req *rolev1.GetRoleRequest) (*rolev1.GetRoleResponse, error) {
	role, errors := s.roleService.GetRole(ctx, req.GetRoleId())
	if errors != nil {
		logger.Error("GetRole request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetRoleInternal.Messagef(errors))})
		return nil, err
	}

	return &rolev1.GetRoleResponse{Role: roleRepoToProto(role)}, nil
}

// UpdateRole maps gRPC request to partial-update service input.
func (s *RoleHandler) UpdateRole(ctx context.Context, req *rolev1.UpdateRoleRequest) (*rolev1.UpdateRoleResponse, error) {
	if req.GetRole() == nil {
		logger.Error("UpdateRole request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRolePayloadRequired)})
		return nil, err
	}
	_, errors := s.roleService.UpdateRole(ctx, rolesvc.UpdateRoleInput{
		ID:          req.GetRole().GetRoleId(),
		Name:        req.GetRole().GetName(),
		Description: req.GetRole().GetDescription(),
		ModuleID:    req.GetRole().GetModuleId(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		logger.Error("UpdateRole request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUpdateRoleInternal.Messagef(errors))})
		return nil, err
	}
	return &rolev1.UpdateRoleResponse{Success: true}, nil
}

// DeleteRole triggers role soft-delete and returns deletion status.
func (s *RoleHandler) DeleteRole(ctx context.Context, req *rolev1.DeleteRoleRequest) (*rolev1.DeleteRoleResponse, error) {
	if errors := s.roleService.DeleteRole(ctx, req.GetRoleId(), req.GetUserId()); errors != nil {
		logger.Error("DeleteRole request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrDeleteRoleInternal.Messagef(errors))})
		return nil, err
	}
	return &rolev1.DeleteRoleResponse{Success: true}, nil
}

// ListRoles returns filtered and paginated roles plus pagination metadata.
func (s *RoleHandler) GetRoles(ctx context.Context, req *rolev1.GetRolesRequest) (*rolev1.GetRolesResponse, error) {
	var limit int32 = 20
	var page int32 = 1
	if req.GetPagination() != nil {
		if req.GetPagination().GetLimit() > 0 {
			limit = req.GetPagination().GetLimit()
		}
		if req.GetPagination().GetPage() > 0 {
			page = req.GetPagination().GetPage()
		}
	}
	offset := (page - 1) * limit

	sorts := make([]*basev1.Sort, 0, len(req.GetSorts()))
	sorts = append(sorts, req.GetSorts()...)
	filters := make([]*basev1.Filter, 0, len(req.GetFilters()))
	filters = append(filters, req.GetFilters()...)

	roles, total, appliedLimit, appliedOffset, err := s.roleService.ListRoles(ctx, rolesvc.ListRolesInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   sorts,
		Filters: filters,
	})
	if err != nil {
		logger.Error("ListRoles request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrListRolesInternal.Messagef(err))})
		return nil, err
	}

	roleItems := make([]*rolev1.Role, 0, len(roles))
	for _, role := range roles {
		roleItems = append(roleItems, roleRepoToProto(role))
	}

	totalItems := int32(total)
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &rolev1.GetRolesResponse{
		Data: roleItems,
		Pagination: &basev1.PaginationResponse{
			Page:       (appliedOffset / appliedLimit) + 1,
			Limit:      appliedLimit,
			TotalItems: totalItems,
			TotalPages: totalPages,
			ItemCount:  int32(len(roleItems)),
		},
	}, nil
}
