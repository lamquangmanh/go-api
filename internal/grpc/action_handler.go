package grpc

import (
	"context"
	"math"
	"time"

	"go-api/internal/repository"
	actionsvc "go-api/internal/service"
	"go-api/pkg/constants"
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	actionv1 "github.com/lamquangmanh/protobuf/gen/go/proto/action/v1"
	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
	"google.golang.org/grpc/codes"
)

type ActionHandler struct {
	actionv1.UnimplementedActionServiceServer
	actionService *actionsvc.ActionService
}

// NewActionHandler creates a gRPC transport handler for ActionService.
func NewActionHandler(actionService *actionsvc.ActionService) *ActionHandler {
	return &ActionHandler{actionService: actionService}
}

// mapRequestTypeToProto converts repository request type into protobuf enum.
func mapRequestTypeToProto(value repository.RequestType) actionv1.ActionRequestType {
	switch value {
	case repository.RequestTypeHTTP:
		return actionv1.ActionRequestType_ACTION_REQUEST_TYPE_HTTP
	case repository.RequestTypeGRAPHQL:
		return actionv1.ActionRequestType_ACTION_REQUEST_TYPE_GRAPHQL
	case repository.RequestTypeGRPC:
		return actionv1.ActionRequestType_ACTION_REQUEST_TYPE_GRPC
	case repository.RequestTypeWEBSOCKET:
		return actionv1.ActionRequestType_ACTION_REQUEST_TYPE_WEBSOCKET
	default:
		return actionv1.ActionRequestType_ACTION_REQUEST_TYPE_VIEW
	}
}

// mapRequestTypeFromProtoAction converts protobuf request type into repository enum.
func mapRequestTypeFromProtoAction(value actionv1.ActionRequestType) repository.RequestType {
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

// actionRepoToProto maps repository action model to protobuf response model.
func actionRepoToProto(item *repository.Action) *actionv1.Action {
	description := ""
	if item.Description.Valid {
		description = item.Description.String
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

	return &actionv1.Action{
		ActionId:      item.ActionID.String(),
		ResourceId:    item.ResourceID.String(),
		Name:          item.Name,
		Description:   description,
		RequestType:   mapRequestTypeToProto(item.RequestType),
		Url:           item.Url,
		Method:        item.Method,
		CreatedAt:     item.CreatedAt.Time.UTC().Format(time.RFC3339),
		CreatedUserId: createdUserID,
		UpdatedAt:     updatedAt,
		UpdatedUserId: updatedUserID,
		DeletedAt:     deletedAt,
		DeletedUserId: deletedUserID,
	}
}

// GetAction returns a single action by ID.
func (h *ActionHandler) GetAction(ctx context.Context, req *actionv1.GetActionRequest) (*actionv1.GetActionResponse, error) {
	item, errors := h.actionService.GetAction(ctx, req.GetActionId())
	if errors != nil {
		logger.Error("GetAction request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &actionv1.GetActionResponse{Action: actionRepoToProto(item)}, nil
}

// GetActions returns paginated actions with validated sorts and filters.
func (h *ActionHandler) GetActions(ctx context.Context, req *actionv1.GetActionsRequest) (*actionv1.GetActionsResponse, error) {
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

	items, total, appliedLimit, appliedOffset, err := h.actionService.ListActions(ctx, actionsvc.ListActionsInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   sorts,
		Filters: filters,
	})
	if err != nil {
		return &actionv1.GetActionsResponse{Data: nil, Pagination: &basev1.PaginationResponse{}}, nil
	}
	data := make([]*actionv1.Action, 0, len(items))
	for _, item := range items {
		data = append(data, actionRepoToProto(item))
	}
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &actionv1.GetActionsResponse{
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

// CreateAction validates payload and creates a new action.
func (h *ActionHandler) CreateAction(ctx context.Context, req *actionv1.CreateActionRequest) (*actionv1.CreateActionResponse, error) {
	if req.GetAction() == nil {
		logger.Error("CreateAction request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionPayloadRequired)})
		return nil, err
	}
	item, errors := h.actionService.CreateAction(ctx, actionsvc.CreateActionInput{
		ResourceID:  req.GetAction().GetResourceId(),
		Name:        req.GetAction().GetName(),
		Description: req.GetAction().GetDescription(),
		RequestType: mapRequestTypeFromProtoAction(req.GetAction().GetRequestType()),
		URL:         req.GetAction().GetUrl(),
		Method:      req.GetAction().GetMethod(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &actionv1.CreateActionResponse{Action: actionRepoToProto(item)}, nil
}

// UpdateAction updates an existing action.
func (h *ActionHandler) UpdateAction(ctx context.Context, req *actionv1.UpdateActionRequest) (*actionv1.UpdateActionResponse, error) {
	if req.GetAction() == nil {
		logger.Error("UpdateAction request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionPayloadRequired)})
		return nil, err
	}
	_, errors := h.actionService.UpdateAction(ctx, actionsvc.UpdateActionInput{
		ID:          req.GetAction().GetActionId(),
		ResourceID:  req.GetAction().GetResourceId(),
		Name:        req.GetAction().GetName(),
		Description: req.GetAction().GetDescription(),
		RequestType: mapRequestTypeFromProtoAction(req.GetAction().GetRequestType()),
		URL:         req.GetAction().GetUrl(),
		Method:      req.GetAction().GetMethod(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &actionv1.UpdateActionResponse{}, nil
}

// DeleteAction performs a soft-delete for an action.
func (h *ActionHandler) DeleteAction(ctx context.Context, req *actionv1.DeleteActionRequest) (*actionv1.DeleteActionResponse, error) {
	if errors := h.actionService.DeleteAction(ctx, req.GetActionId(), req.GetUserId()); errors != nil {
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &actionv1.DeleteActionResponse{}, nil
}
