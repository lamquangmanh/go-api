package service

import (
	"context"
	"errors"
	"strings"

	basepb "go-api/pkg/api/basepb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"go-api/internal/repository"
)

type ModuleService struct {
	q *repository.Queries
}

// NewModuleService creates a business service for module operations.
func NewModuleService(q *repository.Queries) *ModuleService {
	return &ModuleService{q: q}
}

type CreateModuleInput struct {
	Name        string
	Description string
	ProductID   string
	Icon        string
	URL         string
	ActorUserID string
}

type UpdateModuleInput struct {
	ID          string
	Name        string
	Description string
	ProductID   string
	Icon        string
	URL         string
	ActorUserID string
}

type ListModulesInput struct {
	Limit   int32
	Offset  int32
	Sorts   []*basepb.Sort
	Filters []*basepb.Filter
}

func parseOptionalUUID(v string) (pgtype.UUID, error) {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return pgtype.UUID{Valid: false}, nil
	}
	id, err := uuid.Parse(trimmed)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: id, Valid: true}, nil
}

// CreateModule validates and creates a module.
func (s *ModuleService) CreateModule(ctx context.Context, in CreateModuleInput) (*repository.Module, []*basepb.ErrorMessage) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModuleNameRequired)}
	}
	productID, err := parseOptionalUUID(in.ProductID)
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidModuleProductID)}
	}
	description := strings.TrimSpace(in.Description)
	icon := strings.TrimSpace(in.Icon)
	url := strings.TrimSpace(in.URL)
	item, err := s.q.InsertModule(ctx, repository.InsertModuleParams{
		Name:          name,
		Description:   pgtype.Text{String: description, Valid: description != ""},
		Url:           pgtype.Text{String: url, Valid: url != ""},
		Icon:          pgtype.Text{String: icon, Valid: icon != ""},
		ProductID:     productID,
		CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrCreateModuleInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// GetModule returns a module by ID.
func (s *ModuleService) GetModule(ctx context.Context, id string) (*repository.Module, []*basepb.ErrorMessage) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModuleIDRequired)}
	}
	mid, err := uuid.Parse(id)
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidModuleIDFormat)}
	}
	item, err := s.q.GetModule(ctx, mid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModuleNotFound.Messagef(id))}
		}
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrGetModuleInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// UpdateModule validates input and updates a module.
func (s *ModuleService) UpdateModule(ctx context.Context, in UpdateModuleInput) (*repository.Module, []*basepb.ErrorMessage) {
	in.ID = strings.TrimSpace(in.ID)
	if in.ID == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModuleIDRequired)}
	}
	mid, err := uuid.Parse(in.ID)
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidModuleIDFormat)}
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModuleNameRequired)}
	}
	productID, err := parseOptionalUUID(in.ProductID)
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidModuleProductID)}
	}
	description := strings.TrimSpace(in.Description)
	icon := strings.TrimSpace(in.Icon)
	url := strings.TrimSpace(in.URL)
	item, err := s.q.UpdateModule(ctx, repository.UpdateModuleParams{
		ModuleID:      mid,
		Name:          name,
		Description:   pgtype.Text{String: description, Valid: description != ""},
		Url:           pgtype.Text{String: url, Valid: url != ""},
		Icon:          pgtype.Text{String: icon, Valid: icon != ""},
		ProductID:     productID,
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModuleNotFound.Messagef(in.ID))}
		}
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUpdateModuleInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// DeleteModule performs soft-delete for a module.
func (s *ModuleService) DeleteModule(ctx context.Context, id string, actorUserID string) []*basepb.ErrorMessage {
	id = strings.TrimSpace(id)
	if id == "" {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModuleIDRequired)}
	}
	mid, err := uuid.Parse(id)
	if err != nil {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidModuleIDFormat)}
	}
	if _, err := s.q.DeleteModule(ctx, repository.DeleteModuleParams{
		ModuleID:      mid,
		DeletedUserID: pgtype.Text{String: actorUserID, Valid: strings.TrimSpace(actorUserID) != ""},
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrModuleNotFound.Messagef(id))}
		}
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrDeleteModuleInternal.Messagef(err.Error()))}
	}
	return nil
}

// ListModules returns paginated modules with reusable sort/filter validation.
func (s *ModuleService) ListModules(ctx context.Context, in ListModulesInput) ([]*repository.Module, int64, int32, int32, error) {
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}

	cfg, ok := utils.GetTableQueryConfig(utils.TableModules)
	if !ok {
		return nil, 0, 0, 0, constants.ErrMissingModulesQueryConfig.Status()
	}
	sorts, err := utils.NormalizeSorts(in.Sorts, cfg, 5)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidModuleSort.Statusf(err)
	}
	filters, err := utils.NormalizeFilters(in.Filters, cfg, 10)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidModuleFilter.Statusf(err)
	}

	listSQL, listArgs, countSQL, countArgs := utils.BuildListAndCountQueries(cfg, sorts, filters, limit, offset)
	items, total, err := s.q.ListModulesByDynamicQuery(ctx, repository.DynamicModuleListParams{
		ListSQL:   listSQL,
		ListArgs:  listArgs,
		CountSQL:  countSQL,
		CountArgs: countArgs,
	})
	if err != nil {
		return nil, 0, 0, 0, constants.ErrListModulesInternal.Statusf(err)
	}

	return items, total, limit, offset, nil
}
