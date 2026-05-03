package grpc

import (
	"context"
	"math"
	"time"

	"go-api/internal/repository"
	resourcesvc "go-api/internal/service"
	actionpb "go-api/pkg/api/actionpb"
	basepb "go-api/pkg/api/basepb"
	resourcepb "go-api/pkg/api/resourcepb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"
)

// ResourceHandler implements gRPC ResourceService
type ResourceHandler struct {
	resourcepb.UnimplementedResourceServiceServer
	resourceService *resourcesvc.ResourceService
}

// NewResourceHandler creates a new ResourceHandler
func NewResourceHandler(resourceService *resourcesvc.ResourceService) *ResourceHandler {
	return &ResourceHandler{resourceService: resourceService}
}

func mapRequestTypeFromProto(value actionpb.ActionRequestType) repository.RequestType {
	switch value {
	case actionpb.ActionRequestType_ACTION_REQUEST_TYPE_HTTP:
		return repository.RequestTypeHTTP
	case actionpb.ActionRequestType_ACTION_REQUEST_TYPE_GRAPHQL:
		return repository.RequestTypeGRAPHQL
	case actionpb.ActionRequestType_ACTION_REQUEST_TYPE_GRPC:
		return repository.RequestTypeGRPC
	case actionpb.ActionRequestType_ACTION_REQUEST_TYPE_WEBSOCKET:
		return repository.RequestTypeWEBSOCKET
	default:
		return repository.RequestTypeVIEW
	}
}

func resourceRepoToProto(item *repository.Resource) *resourcepb.Resource {
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

	return &resourcepb.Resource{
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
func (h *ResourceHandler) GetResource(ctx context.Context, req *resourcepb.GetResourceRequest) (*resourcepb.GetResourceResponse, error) {
	item, errors := h.resourceService.GetResource(ctx, req.GetResourceId())
	if errors != nil {
		return &resourcepb.GetResourceResponse{Errors: errors}, nil
	}
	return &resourcepb.GetResourceResponse{Resource: resourceRepoToProto(item)}, nil
}

// GetResources retrieves paginated resource list with filtering
func (h *ResourceHandler) GetResources(ctx context.Context, req *resourcepb.GetResourcesRequest) (*resourcepb.GetResourcesResponse, error) {
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

	items, total, appliedLimit, appliedOffset, err := h.resourceService.ListResources(ctx, resourcesvc.ListResourcesInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   req.GetSorts(),
		Filters: req.GetFilters(),
	})
	if err != nil {
		return &resourcepb.GetResourcesResponse{Data: nil, Pagination: &basepb.PaginationResponse{}}, nil
	}
	data := make([]*resourcepb.Resource, 0, len(items))
	for _, item := range items {
		data = append(data, resourceRepoToProto(item))
	}
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &resourcepb.GetResourcesResponse{
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

// CreateResource creates a new resource
func (h *ResourceHandler) CreateResource(ctx context.Context, req *resourcepb.CreateResourceRequest) (*resourcepb.CreateSuccess, error) {
	if req.GetResource() == nil {
		return &resourcepb.CreateSuccess{Errors: []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourcePayloadRequired)}}, nil
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
		return &resourcepb.CreateSuccess{Errors: errors}, nil
	}
	return &resourcepb.CreateSuccess{Resource: resourceRepoToProto(item)}, nil
}

// UpdateResource updates an existing resource
func (h *ResourceHandler) UpdateResource(ctx context.Context, req *resourcepb.UpdateResourceRequest) (*basepb.UpdateSuccess, error) {
	if req.GetResource() == nil {
		return utils.UpdateErr(utils.ErrMsg(constants.ErrResourcePayloadRequired)), nil
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
		return utils.UpdateErr(errors...), nil
	}
	return utils.UpdateOK(), nil
}

// DeleteResource performs soft-delete on a resource
func (h *ResourceHandler) DeleteResource(ctx context.Context, req *resourcepb.DeleteResourceRequest) (*basepb.DeleteSuccess, error) {
	if errors := h.resourceService.DeleteResource(ctx, req.GetResourceId(), req.GetUserId()); errors != nil {
		return utils.DeleteErr(errors...), nil
	}
	return utils.DeleteOK(), nil
}
