package grpc

import (
	"context"
	"math"
	"time"

	"go-api/internal/repository"
	resourcesvc "go-api/internal/service"

	// "go-api/pkg/api/basepb"
	// "go-api/pkg/api/resourcepb"
	"go-api/pkg/constants"
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	actionv1 "github.com/lamquangmanh/protobuf/gen/go/proto/action/v1"
	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
	resourcev1 "github.com/lamquangmanh/protobuf/gen/go/proto/resource/v1"
	"google.golang.org/grpc/codes"
)

// ResourceHandler implements gRPC ResourceService
type ResourceHandler struct {
	resourcev1.UnimplementedResourceServiceServer
	resourceService *resourcesvc.ResourceService
}

// NewResourceHandler creates a new ResourceHandler
func NewResourceHandler(resourceService *resourcesvc.ResourceService) *ResourceHandler {
	return &ResourceHandler{resourceService: resourceService}
}

func mapRequestTypeFromProto(value actionv1.ActionRequestType) repository.RequestType {
	switch value {
	case actionv1.ActionRequestType_ACTION_REQUEST_TYPE_HTTP:
		return repository.RequestTypeHTTP
	case actionv1.ActionRequestType_ACTION_REQUEST_TYPE_GRAPHQL:
		return repository.RequestTypeGRAPHQL
	case actionv1.ActionRequestType_ACTION_REQUEST_TYPE_GRPC:
		return repository.RequestTypeGRPC
	case actionv1.ActionRequestType_ACTION_REQUEST_TYPE_WEBSOCKET:
		return repository.RequestTypeWEBSOCKET
	default:
		return repository.RequestTypeVIEW
	}
}

func resourceRepoToProto(item *repository.Resource) *resourcev1.Resource {
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

	return &resourcev1.Resource{
		ResourceId:    item.ResourceID.String(),
		Name:          item.Name,
		ModuleId:      item.ModuleID.String(),
		CreatedUserId: createdUserID,
		UpdatedAt:     updatedAt,
		UpdatedUserId: updatedUserID,
		DeletedAt:     deletedAt,
		DeletedUserId: deletedUserID,
		CreatedAt:     item.CreatedAt.Time.UTC().Format(time.RFC3339),
	}
}

// GetResource retrieves a single resource by ID
func (h *ResourceHandler) GetResource(ctx context.Context, req *resourcev1.GetResourceRequest) (*resourcev1.GetResourceResponse, error) {
	item, errors := h.resourceService.GetResource(ctx, req.GetResourceId())
	if errors != nil {
		logger.Error("GetResource request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &resourcev1.GetResourceResponse{Resource: resourceRepoToProto(item)}, nil
}

// GetResources retrieves paginated resource list with filtering
func (h *ResourceHandler) GetResources(ctx context.Context, req *resourcev1.GetResourcesRequest) (*resourcev1.GetResourcesResponse, error) {
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

	items, total, appliedLimit, appliedOffset, err := h.resourceService.ListResources(ctx, resourcesvc.ListResourcesInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   sorts,
		Filters: filters,
	})
	if err != nil {
		logger.Error("ListResources request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrListResourcesInternal.Messagef(err))})
		return nil, err
	}
	data := make([]*resourcev1.Resource, 0, len(items))
	for _, item := range items {
		data = append(data, resourceRepoToProto(item))
	}
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &resourcev1.GetResourcesResponse{
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

// CreateResource creates a new resource
func (h *ResourceHandler) CreateResource(ctx context.Context, req *resourcev1.CreateResourceRequest) (*resourcev1.CreateResourceResponse, error) {
	if req.GetResource() == nil {
		logger.Error("CreateResource request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrResourcePayloadRequired)})
		return nil, err
	}
	actions := make([]resourcesvc.ResourceActionInput, 0, len(req.GetResource().GetActions()))
	for _, action := range req.GetResource().GetActions() {
		actions = append(actions, resourcesvc.ResourceActionInput{
			Name:        action.GetName(),
			Description: action.GetDescription(),
			RequestType: mapRequestTypeFromProto(action.GetRequestType()),
			URL:         action.GetUrl(),
			Method:      action.GetMethod(),
		})
	}
	item, errors := h.resourceService.CreateResource(ctx, resourcesvc.CreateResourceInput{
		Name:        req.GetResource().GetName(),
		ModuleID:    req.GetResource().GetModuleId(),
		Actions:     actions,
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		logger.Error("CreateResource request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrCreateResourceInternal.Messagef(errors))})
		return nil, err
	}
	return &resourcev1.CreateResourceResponse{Resource: resourceRepoToProto(item)}, nil
}

// UpdateResource updates an existing resource
func (h *ResourceHandler) UpdateResource(ctx context.Context, req *resourcev1.UpdateResourceRequest) (*resourcev1.UpdateResourceResponse, error) {
	if req.GetResource() == nil {
		logger.Error("UpdateResource request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrResourcePayloadRequired)})
		return nil, err
	}
	actions := make([]resourcesvc.ResourceActionInput, 0, len(req.GetResource().GetActions()))
	for _, action := range req.GetResource().GetActions() {
		actions = append(actions, resourcesvc.ResourceActionInput{
			Name:        action.GetName(),
			Description: action.GetDescription(),
			RequestType: mapRequestTypeFromProto(action.GetRequestType()),
			URL:         action.GetUrl(),
			Method:      action.GetMethod(),
		})
	}
	_, errors := h.resourceService.UpdateResource(ctx, resourcesvc.UpdateResourceInput{
		ID:          req.GetResource().GetResourceId(),
		Name:        req.GetResource().GetName(),
		ModuleID:    req.GetResource().GetModuleId(),
		Actions:     actions,
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		logger.Error("UpdateResource request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUpdateResourceInternal.Messagef(errors))})
		return nil, err
	}
	return &resourcev1.UpdateResourceResponse{Success: true}, nil
}

// DeleteResource performs soft-delete on a resource
func (h *ResourceHandler) DeleteResource(ctx context.Context, req *resourcev1.DeleteResourceRequest) (*resourcev1.DeleteResourceResponse, error) {
	if errors := h.resourceService.DeleteResource(ctx, req.GetResourceId(), req.GetUserId()); errors != nil {
		logger.Error("DeleteResource request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrDeleteResourceInternal.Messagef(errors))})
		return nil, err
	}
	return &resourcev1.DeleteResourceResponse{Success: true}, nil
}
