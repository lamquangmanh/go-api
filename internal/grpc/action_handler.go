package grpc

import (
	"context"
	"math"
	"time"

	"go-api/internal/repository"
	actionsvc "go-api/internal/service"
	actionpb "go-api/pkg/api/actionpb"
	basepb "go-api/pkg/api/basepb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"
)

type ActionHandler struct {
	actionpb.UnimplementedActionServiceServer
	actionService *actionsvc.ActionService
}

// NewActionHandler creates a gRPC transport handler for ActionService.
func NewActionHandler(actionService *actionsvc.ActionService) *ActionHandler {
	return &ActionHandler{actionService: actionService}
}

// mapRequestTypeToProto converts repository request type into protobuf enum.
func mapRequestTypeToProto(value repository.RequestType) actionpb.ActionRequestType {
	switch value {
	case repository.RequestTypeHTTP:
		return actionpb.ActionRequestType_ACTION_REQUEST_TYPE_HTTP
	case repository.RequestTypeGRAPHQL:
		return actionpb.ActionRequestType_ACTION_REQUEST_TYPE_GRAPHQL
	case repository.RequestTypeGRPC:
		return actionpb.ActionRequestType_ACTION_REQUEST_TYPE_GRPC
	case repository.RequestTypeWEBSOCKET:
		return actionpb.ActionRequestType_ACTION_REQUEST_TYPE_WEBSOCKET
	default:
		return actionpb.ActionRequestType_ACTION_REQUEST_TYPE_VIEW
	}
}

// mapRequestTypeFromProtoAction converts protobuf request type into repository enum.
func mapRequestTypeFromProtoAction(value actionpb.ActionRequestType) repository.RequestType {
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

// actionRepoToProto maps repository action model to protobuf response model.
func actionRepoToProto(item *repository.Action) *actionpb.Action {
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

	return &actionpb.Action{
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
func (h *ActionHandler) GetAction(ctx context.Context, req *actionpb.GetActionRequest) (*actionpb.GetActionResponse, error) {
	item, errors := h.actionService.GetAction(ctx, req.GetActionId())
	if errors != nil {
		return &actionpb.GetActionResponse{Errors: errors}, nil
	}
	return &actionpb.GetActionResponse{Action: actionRepoToProto(item)}, nil
}

// GetActions returns paginated actions with validated sorts and filters.
func (h *ActionHandler) GetActions(ctx context.Context, req *actionpb.GetActionsRequest) (*actionpb.GetActionsResponse, error) {
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

	items, total, appliedLimit, appliedOffset, err := h.actionService.ListActions(ctx, actionsvc.ListActionsInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   req.GetSorts(),
		Filters: req.GetFilters(),
	})
	if err != nil {
		return &actionpb.GetActionsResponse{Data: nil, Pagination: &basepb.PaginationResponse{}}, nil
	}
	data := make([]*actionpb.Action, 0, len(items))
	for _, item := range items {
		data = append(data, actionRepoToProto(item))
	}
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &actionpb.GetActionsResponse{
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

// CreateAction validates payload and creates a new action.
func (h *ActionHandler) CreateAction(ctx context.Context, req *actionpb.CreateActionRequest) (*actionpb.CreateSuccess, error) {
	if req.GetAction() == nil {
		return &actionpb.CreateSuccess{Errors: []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrActionPayloadRequired)}}, nil
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
		return &actionpb.CreateSuccess{Errors: errors}, nil
	}
	return &actionpb.CreateSuccess{Action: actionRepoToProto(item)}, nil
}

// UpdateAction updates an existing action.
func (h *ActionHandler) UpdateAction(ctx context.Context, req *actionpb.UpdateActionRequest) (*basepb.UpdateSuccess, error) {
	if req.GetAction() == nil {
		return utils.UpdateErr(utils.ErrMsg(constants.ErrActionPayloadRequired)), nil
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
		return utils.UpdateErr(errors...), nil
	}
	return utils.UpdateOK(), nil
}

// DeleteAction performs a soft-delete for an action.
func (h *ActionHandler) DeleteAction(ctx context.Context, req *actionpb.DeleteActionRequest) (*basepb.DeleteSuccess, error) {
	if errors := h.actionService.DeleteAction(ctx, req.GetActionId(), req.GetUserId()); errors != nil {
		return utils.DeleteErr(errors...), nil
	}
	return utils.DeleteOK(), nil
}
