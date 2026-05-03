package grpc

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"go-api/internal/repository"
	modulesvc "go-api/internal/service"
	productsvc "go-api/internal/service"
	basepb "go-api/pkg/api/basepb"
	modulepb "go-api/pkg/api/modulepb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"
)

type ModuleHandler struct {
	modulepb.UnimplementedModuleServiceServer
	moduleService  *modulesvc.ModuleService
	productService *productsvc.ProductService
}

// NewModuleHandler creates a gRPC transport handler for ModuleService.
func NewModuleHandler(moduleService *modulesvc.ModuleService, productService *productsvc.ProductService) *ModuleHandler {
	return &ModuleHandler{moduleService: moduleService, productService: productService}
}

// moduleRepoToProto maps repository module model to protobuf response model.
func moduleRepoToProto(item *repository.Module) *modulepb.Module {
	description := ""
	if item.Description.Valid {
		description = item.Description.String
	}
	icon := ""
	if item.Icon.Valid {
		icon = item.Icon.String
	}
	url := ""
	if item.Url.Valid {
		url = item.Url.String
	}
	productID := ""
	if item.ProductID.Valid {
		productID = uuid.UUID(item.ProductID.Bytes).String()
	}
	createdUserID := ""
	if item.CreatedUserID.Valid {
		createdUserID = item.CreatedUserID.String
	}
	updatedAt := ""
	if item.UpdatedAt.Valid {
		updatedAt = item.UpdatedAt.Time.UTC().Format(time.RFC3339)
	}
	updatedUserID := ""
	if item.UpdatedUserID.Valid {
		updatedUserID = item.UpdatedUserID.String
	}
	deletedAt := ""
	if item.DeletedAt.Valid {
		deletedAt = item.DeletedAt.Time.UTC().Format(time.RFC3339)
	}
	deletedUserID := ""
	if item.DeletedUserID.Valid {
		deletedUserID = item.DeletedUserID.String
	}

	return &modulepb.Module{
		ModuleId:      item.ModuleID.String(),
		Name:          item.Name,
		Description:   description,
		CreatedAt:     item.CreatedAt.Time.UTC().Format(time.RFC3339),
		CreatedUserId: createdUserID,
		UpdatedAt:     updatedAt,
		UpdatedUserId: updatedUserID,
		DeletedAt:     deletedAt,
		DeletedUserId: deletedUserID,
		ProductId:     productID,
		Icon:          &icon,
		Url:           &url,
	}
}

// GetModule returns a single module by ID.
func (h *ModuleHandler) GetModule(ctx context.Context, req *modulepb.GetModuleRequest) (*modulepb.GetModuleResponse, error) {
	item, errors := h.moduleService.GetModule(ctx, req.GetModuleId())
	if errors != nil {
		return &modulepb.GetModuleResponse{Errors: errors}, nil
	}
	return &modulepb.GetModuleResponse{Module: moduleRepoToProto(item)}, nil
}

// GetModules returns paginated modules with validated sorts and filters.
func (h *ModuleHandler) GetModules(ctx context.Context, req *modulepb.GetModulesRequest) (*modulepb.GetModulesResponse, error) {
	limit := int32(20)
	page := int32(1)
	if req.GetPagination() != nil {
		if req.GetPagination().GetLimit() > 0 {
			limit = req.GetPagination().GetLimit()
		}
		if req.GetPagination().GetPage() > 0 {
			page = req.GetPagination().GetPage()
		}
	}
	offset := (page - 1) * limit

	items, total, appliedLimit, appliedOffset, err := h.moduleService.ListModules(ctx, modulesvc.ListModulesInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   req.GetSorts(),
		Filters: req.GetFilters(),
	})
	if err != nil {
		return &modulepb.GetModulesResponse{Data: nil, Pagination: &basepb.PaginationResponse{}}, nil
	}
	data := make([]*modulepb.Module, 0, len(items))
	for _, item := range items {
		data = append(data, moduleRepoToProto(item))
	}
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &modulepb.GetModulesResponse{
		Data: data,
		Pagination: &basepb.PaginationResponse{
			Page:       (appliedOffset / appliedLimit) + 1,
			Limit:      appliedLimit,
			TotalItems: int32(total),
			TotalPages: totalPages,
			ItemCount:  int32(len(data)),
		},
	}, nil
}

// CreateModule validates payload and creates a new module.
func (h *ModuleHandler) CreateModule(ctx context.Context, req *modulepb.CreateModuleRequest) (*modulepb.CreateSuccess, error) {
	if req.GetModule() == nil {
		return &modulepb.CreateSuccess{Errors: []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModulePayloadRequired)}}, nil
	}
	icon := ""
	if req.GetModule().Icon != nil {
		icon = req.GetModule().GetIcon()
	}
	url := ""
	if req.GetModule().Url != nil {
		url = req.GetModule().GetUrl()
	}

	if _, productErrors := h.productService.GetProduct(ctx, req.GetModule().GetProductId()); productErrors != nil {
		return &modulepb.CreateSuccess{Errors: []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidModuleProductID)}}, nil
	}

	item, errors := h.moduleService.CreateModule(ctx, modulesvc.CreateModuleInput{
		Name:        req.GetModule().GetName(),
		Description: req.GetModule().GetDescription(),
		ProductID:   req.GetModule().GetProductId(),
		Icon:        icon,
		URL:         url,
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		return &modulepb.CreateSuccess{Errors: errors}, nil
	}
	return &modulepb.CreateSuccess{Module: moduleRepoToProto(item)}, nil
}

// UpdateModule updates an existing module.
func (h *ModuleHandler) UpdateModule(ctx context.Context, req *modulepb.UpdateModuleRequest) (*basepb.UpdateSuccess, error) {
	if req.GetModule() == nil {
		return utils.UpdateErr(utils.ErrMsg(constants.ErrModulePayloadRequired)), nil
	}
	icon := ""
	if req.GetModule().Icon != nil {
		icon = req.GetModule().GetIcon()
	}
	url := ""
	if req.GetModule().Url != nil {
		url = req.GetModule().GetUrl()
	}

	if _, productErrors := h.productService.GetProduct(ctx, req.GetModule().GetProductId()); productErrors != nil {
		return utils.UpdateErr(utils.ErrMsg(constants.ErrInvalidModuleProductID)), nil
	}

	_, errors := h.moduleService.UpdateModule(ctx, modulesvc.UpdateModuleInput{
		ID:          req.GetModule().GetModuleId(),
		Name:        req.GetModule().GetName(),
		Description: req.GetModule().GetDescription(),
		ProductID:   req.GetModule().GetProductId(),
		Icon:        icon,
		URL:         url,
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		return utils.UpdateErr(errors...), nil
	}
	return utils.UpdateOK(), nil
}

// DeleteModule performs a soft-delete for a module.
func (h *ModuleHandler) DeleteModule(ctx context.Context, req *modulepb.DeleteModuleRequest) (*basepb.DeleteSuccess, error) {
	if errors := h.moduleService.DeleteModule(ctx, req.GetModuleId(), req.GetUserId()); errors != nil {
		return utils.DeleteErr(errors...), nil
	}
	return utils.DeleteOK(), nil
}
