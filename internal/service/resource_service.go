package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-api/internal/repository"
	basepb "go-api/pkg/api/basepb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"
)

type ResourceService struct {
	q    *repository.Queries
	pool *pgxpool.Pool
}

func NewResourceService(q *repository.Queries, pool *pgxpool.Pool) *ResourceService {
	return &ResourceService{q: q, pool: pool}
}

type ResourceActionInput struct {
	Name        string
	Description string
	RequestType repository.RequestType
	URL         string
	Method      string
}

type CreateResourceInput struct {
	Name        string
	ModuleID    string
	Actions     []ResourceActionInput
	ActorUserID string
}

type UpdateResourceInput struct {
	ID          string
	Name        string
	ModuleID    string
	Actions     []ResourceActionInput
	ActorUserID string
}

type ListResourcesInput struct {
	Limit   int32
	Offset  int32
	Sorts   []*basepb.Sort
	Filters []*basepb.Filter
}

func (s *ResourceService) CreateResource(ctx context.Context, in CreateResourceInput) (*repository.Resource, []*basepb.ErrorMessage) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourceNameRequired)}
	}
	mid, err := uuid.Parse(strings.TrimSpace(in.ModuleID))
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidResourceModuleID)}
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrBeginTransaction.Messagef(err.Error()))}
	}
	defer tx.Rollback(ctx)
	q := s.q.WithTx(tx)

	item, err := q.InsertResource(ctx, repository.InsertResourceParams{
		Name:          name,
		ModuleID:      mid,
		CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrCreateResourceInternal.Messagef(err.Error()))}
	}

	for _, action := range in.Actions {
		aName := strings.TrimSpace(action.Name)
		if aName == "" {
			continue
		}
		reqType := action.RequestType
		if reqType == "" {
			reqType = repository.RequestTypeVIEW
		}
		desc := strings.TrimSpace(action.Description)
		if _, err := q.InsertAction(ctx, repository.InsertActionParams{
			ResourceID:    item.ResourceID,
			Name:          aName,
			Description:   pgtype.Text{String: desc, Valid: desc != ""},
			RequestType:   reqType,
			Url:           strings.TrimSpace(action.URL),
			Method:        strings.TrimSpace(action.Method),
			CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
			UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		}); err != nil {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrCreateResourceActionInternal.Messagef(err.Error()))}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrCommitTransaction.Messagef(err.Error()))}
	}
	return item, nil
}

func (s *ResourceService) GetResource(ctx context.Context, id string) (*repository.Resource, []*basepb.ErrorMessage) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourceIDRequired)}
	}
	rid, err := uuid.Parse(id)
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidResourceIDFormat)}
	}
	item, err := s.q.GetResource(ctx, rid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourceNotFound.Messagef(id))}
		}
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrGetResourceInternal.Messagef(err.Error()))}
	}
	return item, nil
}

func (s *ResourceService) UpdateResource(ctx context.Context, in UpdateResourceInput) (*repository.Resource, []*basepb.ErrorMessage) {
	in.ID = strings.TrimSpace(in.ID)
	if in.ID == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourceIDRequired)}
	}
	rid, err := uuid.Parse(in.ID)
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidResourceIDFormat)}
	}
	mid, err := uuid.Parse(strings.TrimSpace(in.ModuleID))
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidResourceModuleID)}
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourceNameRequired)}
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrBeginTransaction.Messagef(err.Error()))}
	}
	defer tx.Rollback(ctx)
	q := s.q.WithTx(tx)

	item, err := q.UpdateResource(ctx, repository.UpdateResourceParams{
		ResourceID:    rid,
		Name:          name,
		ModuleID:      mid,
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourceNotFound.Messagef(in.ID))}
		}
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUpdateResourceInternal.Messagef(err.Error()))}
	}

	oldActions, err := q.ListActionsByResource(ctx, repository.ListActionsByResourceParams{ResourceID: rid, Limit: 10000, Offset: 0})
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrListOldActionsInternal.Messagef(err.Error()))}
	}
	for _, oldAction := range oldActions {
		if _, err := q.DeleteAction(ctx, repository.DeleteActionParams{
			ActionID:      oldAction.ActionID,
			DeletedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		}); err != nil {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrDeleteOldActionInternal.Messagef(err.Error()))}
		}
	}

	for _, action := range in.Actions {
		aName := strings.TrimSpace(action.Name)
		if aName == "" {
			continue
		}
		reqType := action.RequestType
		if reqType == "" {
			reqType = repository.RequestTypeVIEW
		}
		desc := strings.TrimSpace(action.Description)
		if _, err := q.InsertAction(ctx, repository.InsertActionParams{
			ResourceID:    rid,
			Name:          aName,
			Description:   pgtype.Text{String: desc, Valid: desc != ""},
			RequestType:   reqType,
			Url:           strings.TrimSpace(action.URL),
			Method:        strings.TrimSpace(action.Method),
			CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
			UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		}); err != nil {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInsertActionInternal.Messagef(err.Error()))}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrCommitTransaction.Messagef(err.Error()))}
	}
	return item, nil
}

func (s *ResourceService) DeleteResource(ctx context.Context, id string, actorUserID string) []*basepb.ErrorMessage {
	id = strings.TrimSpace(id)
	if id == "" {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourceIDRequired)}
	}
	rid, err := uuid.Parse(id)
	if err != nil {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidResourceIDFormat)}
	}
	if _, err := s.q.DeleteResource(ctx, repository.DeleteResourceParams{
		ResourceID:    rid,
		DeletedUserID: pgtype.Text{String: actorUserID, Valid: strings.TrimSpace(actorUserID) != ""},
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrResourceNotFound.Messagef(id))}
		}
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrDeleteResourceInternal.Messagef(err.Error()))}
	}
	return nil
}

// ListResources returns paginated resources with reusable sort/filter validation.
func (s *ResourceService) ListResources(ctx context.Context, in ListResourcesInput) ([]*repository.Resource, int64, int32, int32, error) {
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}

	cfg, ok := utils.GetTableQueryConfig(utils.TableResources)
	if !ok {
		return nil, 0, 0, 0, constants.ErrMissingResourcesQueryConfig.Status()
	}
	sorts, err := utils.NormalizeSorts(in.Sorts, cfg, 5)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidResourceSort.Statusf(err)
	}
	filters, err := utils.NormalizeFilters(in.Filters, cfg, 10)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidResourceFilter.Statusf(err)
	}

	listSQL, listArgs, countSQL, countArgs := utils.BuildListAndCountQueries(cfg, sorts, filters, limit, offset)
	items, total, err := s.q.ListResourcesByDynamicQuery(ctx, repository.DynamicResourceListParams{
		ListSQL:   listSQL,
		ListArgs:  listArgs,
		CountSQL:  countSQL,
		CountArgs: countArgs,
	})
	if err != nil {
		return nil, 0, 0, 0, constants.ErrListResourcesInternal.Statusf(err)
	}

	return items, total, limit, offset, nil
}
