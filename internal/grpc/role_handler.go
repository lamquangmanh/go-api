package grpc

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"go-api/internal/repository"
	rolesvc "go-api/internal/service"
	basepb "go-api/pkg/api/basepb"
	rolepb "go-api/pkg/api/rolepb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"
)

type RoleHandler struct {
	rolepb.UnimplementedRoleServiceServer
	roleService *rolesvc.RoleService
}

// NewRoleHandler creates a thin gRPC handler that delegates business logic to RoleService.
func NewRoleHandler(roleService *rolesvc.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// roleRepoToProto converts repository.Role to protobuf Role response format.
func roleRepoToProto(role *repository.Role) *rolepb.Role {
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

	return &rolepb.Role{
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
func (s *RoleHandler) CreateRole(ctx context.Context, req *rolepb.CreateRoleRequest) (*rolepb.CreateSuccess, error) {
	if req.GetRole() == nil {
		return &rolepb.CreateSuccess{Errors: []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrRolePayloadRequired)}}, nil
	}
	role, errors := s.roleService.CreateRole(ctx, rolesvc.CreateRoleInput{
		Name:        req.GetRole().GetName(),
		Description: req.GetRole().GetDescription(),
		ModuleID:    req.GetRole().GetModuleId(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		return &rolepb.CreateSuccess{Errors: errors}, nil
	}

	return &rolepb.CreateSuccess{Role: roleRepoToProto(role)}, nil
}

// GetRole fetches a role by ID and maps it into protobuf response.
func (s *RoleHandler) GetRole(ctx context.Context, req *rolepb.GetRoleRequest) (*rolepb.GetRoleResponse, error) {
	role, errors := s.roleService.GetRole(ctx, req.GetRoleId())
	if errors != nil {
		return &rolepb.GetRoleResponse{Errors: errors}, nil
	}

	return &rolepb.GetRoleResponse{Role: roleRepoToProto(role)}, nil
}

// UpdateRole maps gRPC request to partial-update service input.
func (s *RoleHandler) UpdateRole(ctx context.Context, req *rolepb.UpdateRoleRequest) (*basepb.UpdateSuccess, error) {
	if req.GetRole() == nil {
		return utils.UpdateErr(utils.ErrMsg(constants.ErrRolePayloadRequired)), nil
	}
	_, errors := s.roleService.UpdateRole(ctx, rolesvc.UpdateRoleInput{
		ID:          req.GetRole().GetRoleId(),
		Name:        req.GetRole().GetName(),
		Description: req.GetRole().GetDescription(),
		ModuleID:    req.GetRole().GetModuleId(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		return utils.UpdateErr(errors...), nil
	}
	return utils.UpdateOK(), nil
}

// DeleteRole triggers role soft-delete and returns deletion status.
func (s *RoleHandler) DeleteRole(ctx context.Context, req *rolepb.DeleteRoleRequest) (*basepb.DeleteSuccess, error) {
	if errors := s.roleService.DeleteRole(ctx, req.GetRoleId(), req.GetUserId()); errors != nil {
		return utils.DeleteErr(errors...), nil
	}
	return utils.DeleteOK(), nil
}

// ListRoles returns filtered and paginated roles plus pagination metadata.
func (s *RoleHandler) GetRoles(ctx context.Context, req *rolepb.GetRolesRequest) (*rolepb.GetRolesResponse, error) {
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

	roles, total, appliedLimit, appliedOffset, err := s.roleService.ListRoles(ctx, rolesvc.ListRolesInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   req.GetSorts(),
		Filters: req.GetFilters(),
	})
	if err != nil {
		return &rolepb.GetRolesResponse{Data: nil, Pagination: &basepb.PaginationResponse{}}, nil
	}

	roleItems := make([]*rolepb.Role, 0, len(roles))
	for _, role := range roles {
		roleItems = append(roleItems, roleRepoToProto(role))
	}

	totalItems := int32(total)
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &rolepb.GetRolesResponse{
		Data: roleItems,
		Pagination: &basepb.PaginationResponse{
			Page:       (appliedOffset / appliedLimit) + 1,
			Limit:      appliedLimit,
			TotalItems: totalItems,
			TotalPages: totalPages,
			ItemCount:  int32(len(roleItems)),
		},
	}, nil
}
