package grpc

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	"go-api/internal/repository"
	modulesvc "go-api/internal/service"
	productsvc "go-api/internal/service"
	"go-api/pkg/constants"
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
	modulev1 "github.com/lamquangmanh/protobuf/gen/go/proto/module/v1"
)

type ModuleHandler struct {
	modulev1.UnimplementedModuleServiceServer
	moduleService  *modulesvc.ModuleService
	productService *productsvc.ProductService
}

// NewModuleHandler creates a gRPC transport handler for ModuleService.
func NewModuleHandler(moduleService *modulesvc.ModuleService, productService *productsvc.ProductService) *ModuleHandler {
	return &ModuleHandler{moduleService: moduleService, productService: productService}
}

// moduleRepoToProto maps repository module model to protobuf response model.
func moduleRepoToProto(item *repository.Module) *modulev1.Module {
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

	return &modulev1.Module{
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
func (h *ModuleHandler) GetModule(ctx context.Context, req *modulev1.GetModuleRequest) (*modulev1.GetModuleResponse, error) {
	item, errors := h.moduleService.GetModule(ctx, req.GetModuleId())
	if errors != nil {
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &modulev1.GetModuleResponse{Module: moduleRepoToProto(item)}, nil
}

// GetModules returns paginated modules with validated sorts and filters.
func (h *ModuleHandler) GetModules(ctx context.Context, req *modulev1.GetModulesRequest) (*modulev1.GetModulesResponse, error) {
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

	sorts := make([]*basev1.Sort, 0, len(req.GetSorts()))
	sorts = append(sorts, req.GetSorts()...)
	filters := make([]*basev1.Filter, 0, len(req.GetFilters()))
	filters = append(filters, req.GetFilters()...)
	logger.Info("filters: filters=%d",
		filters,
	)
	items, total, appliedLimit, appliedOffset, err := h.moduleService.ListModules(ctx, modulesvc.ListModulesInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   sorts,
		Filters: filters,
	})
	if err != nil {
		logger.Error("error: ",
			err,
		)
		return &modulev1.GetModulesResponse{Data: nil, Pagination: &basev1.PaginationResponse{}}, nil
	}
	data := make([]*modulev1.Module, 0, len(items))
	for _, item := range items {
		data = append(data, moduleRepoToProto(item))
	}
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	logger.Info("GetModules: total=%d, appliedLimit=%d, appliedOffset=%d, totalPages=%d, returnedItems=%d",
		total, appliedLimit, appliedOffset, totalPages, len(data),
	)

	return &modulev1.GetModulesResponse{
		Data: data,
		Pagination: &basev1.PaginationResponse{
			Page:       (appliedOffset / appliedLimit) + 1,
			Limit:      appliedLimit,
			TotalItems: int32(total),
			TotalPages: totalPages,
			ItemCount:  int32(len(data)),
		},
	}, nil
}

// CreateModule validates payload and creates a new module.
func (h *ModuleHandler) CreateModule(ctx context.Context, req *modulev1.CreateModuleRequest) (*modulev1.CreateModuleResponse, error) {
	if req.GetModule() == nil {
		logger.Error("CreateModule request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModulePayloadRequired)})
		return nil, err
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
		logger.Error("CreateModule request has invalid product ID")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidModuleProductID)})
		return nil, err
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
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &modulev1.CreateModuleResponse{Module: moduleRepoToProto(item)}, nil
}

// UpdateModule updates an existing module.
func (h *ModuleHandler) UpdateModule(ctx context.Context, req *modulev1.UpdateModuleRequest) (*modulev1.UpdateModuleResponse, error) {
	if req.GetModule() == nil {
		logger.Error("UpdateModule request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModulePayloadRequired)})
		return nil, err
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
		logger.Error("UpdateModule request has invalid product ID")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidModuleProductID)})
		return nil, err
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
		logger.Error("UpdateModule request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &modulev1.UpdateModuleResponse{}, nil
}

// DeleteModule performs a soft-delete for a module.
func (h *ModuleHandler) DeleteModule(ctx context.Context, req *modulev1.DeleteModuleRequest) (*modulev1.DeleteModuleResponse, error) {
	if errors := h.moduleService.DeleteModule(ctx, req.GetModuleId(), req.GetUserId()); errors != nil {
		logger.Error("DeleteModule request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &modulev1.DeleteModuleResponse{}, nil
}
