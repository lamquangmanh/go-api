package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"go-api/internal/repository"
	"go-api/pkg/constants"
	"go-api/pkg/utils"

	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
)

type ActionService struct {
	q *repository.Queries
}

// NewActionService creates a business service for action operations.
func NewActionService(q *repository.Queries) *ActionService {
	return &ActionService{q: q}
}

type CreateActionInput struct {
	ResourceID  string
	Name        string
	Description string
	RequestType repository.RequestType
	URL         string
	Method      string
	ActorUserID string
}

type UpdateActionInput struct {
	ID          string
	ResourceID  string
	Name        string
	Description string
	RequestType repository.RequestType
	URL         string
	Method      string
	ActorUserID string
}

type ListActionsInput struct {
	Limit   int32
	Offset  int32
	Sorts   []*basev1.Sort
	Filters []*basev1.Filter
}

// CreateAction validates and creates an action.
func (s *ActionService) CreateAction(ctx context.Context, in CreateActionInput) (*repository.Action, []*errorv1.ErrorItem) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionNameRequired)}
	}
	rid, err := uuid.Parse(strings.TrimSpace(in.ResourceID))
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidActionResourceID)}
	}
	if in.RequestType == "" {
		in.RequestType = repository.RequestTypeVIEW
	}
	description := strings.TrimSpace(in.Description)
	item, err := s.q.InsertAction(ctx, repository.InsertActionParams{
		ResourceID:    rid,
		Name:          name,
		Description:   pgtype.Text{String: description, Valid: description != ""},
		RequestType:   in.RequestType,
		Url:           strings.TrimSpace(in.URL),
		Method:        strings.TrimSpace(in.Method),
		CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrCreateActionInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// GetAction returns an action by ID.
func (s *ActionService) GetAction(ctx context.Context, id string) (*repository.Action, []*errorv1.ErrorItem) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionIDRequired)}
	}
	aid, err := uuid.Parse(id)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidActionIDFormat)}
	}
	item, err := s.q.GetAction(ctx, aid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionNotFound.Messagef(id))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetActionInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// UpdateAction validates input and updates an action.
func (s *ActionService) UpdateAction(ctx context.Context, in UpdateActionInput) (*repository.Action, []*errorv1.ErrorItem) {
	in.ID = strings.TrimSpace(in.ID)
	if in.ID == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionIDRequired)}
	}
	aid, err := uuid.Parse(in.ID)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidActionIDFormat)}
	}
	rid, err := uuid.Parse(strings.TrimSpace(in.ResourceID))
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidActionResourceID)}
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionNameRequired)}
	}
	if in.RequestType == "" {
		in.RequestType = repository.RequestTypeVIEW
	}
	description := strings.TrimSpace(in.Description)
	item, err := s.q.UpdateAction(ctx, repository.UpdateActionParams{
		ActionID:      aid,
		ResourceID:    rid,
		Name:          name,
		Description:   pgtype.Text{String: description, Valid: description != ""},
		RequestType:   in.RequestType,
		Url:           strings.TrimSpace(in.URL),
		Method:        strings.TrimSpace(in.Method),
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionNotFound.Messagef(in.ID))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUpdateActionInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// DeleteAction performs soft-delete for an action.
func (s *ActionService) DeleteAction(ctx context.Context, id string, actorUserID string) []*errorv1.ErrorItem {
	id = strings.TrimSpace(id)
	if id == "" {
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionIDRequired)}
	}
	aid, err := uuid.Parse(id)
	if err != nil {
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidActionIDFormat)}
	}
	if _, err := s.q.DeleteAction(ctx, repository.DeleteActionParams{
		ActionID:      aid,
		DeletedUserID: pgtype.Text{String: actorUserID, Valid: strings.TrimSpace(actorUserID) != ""},
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrActionNotFound.Messagef(id))}
		}
			return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrDeleteActionInternal.Messagef(err.Error()))}
	}
	return nil
}

// ListActions returns paginated actions with reusable sort/filter validation.
func (s *ActionService) ListActions(ctx context.Context, in ListActionsInput) ([]*repository.Action, int64, int32, int32, error) {
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}

	cfg, ok := utils.GetTableQueryConfig(utils.TableActions)
	if !ok {
		return nil, 0, 0, 0, constants.ErrMissingActionsQueryConfig.Status()
	}
	sorts, err := utils.NormalizeSorts(in.Sorts, cfg, 5)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidActionSort.Statusf(err)
	}
	filters, err := utils.NormalizeFilters(in.Filters, cfg, 10)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidActionFilter.Statusf(err)
	}

	listSQL, listArgs, countSQL, countArgs := utils.BuildListAndCountQueries(cfg, sorts, filters, limit, offset)
	items, total, err := s.q.ListActionsByDynamicQuery(ctx, repository.DynamicActionListParams{
		ListSQL:   listSQL,
		ListArgs:  listArgs,
		CountSQL:  countSQL,
		CountArgs: countArgs,
	})
	if err != nil {
		return nil, 0, 0, 0, constants.ErrListActionsInternal.Statusf(err)
	}

	return items, total, limit, offset, nil
}
