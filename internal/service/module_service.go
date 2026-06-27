package service

import (
	"context"
	"errors"
	"strings"

	"go-api/pkg/constants"
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"

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
	Sorts   []*basev1.Sort
	Filters []*basev1.Filter
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
func (s *ModuleService) CreateModule(ctx context.Context, in CreateModuleInput) (*repository.Module, []*errorv1.ErrorItem) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModuleNameRequired)}
	}
	productID, err := parseOptionalUUID(in.ProductID)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidModuleProductID)}
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
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrCreateModuleInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// GetModule returns a module by ID.
func (s *ModuleService) GetModule(ctx context.Context, id string) (*repository.Module, []*errorv1.ErrorItem) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModuleIDRequired)}
	}
	mid, err := uuid.Parse(id)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidModuleIDFormat)}
	}
	item, err := s.q.GetModule(ctx, mid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModuleNotFound.Messagef(id))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetModuleInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// UpdateModule validates input and updates a module.
func (s *ModuleService) UpdateModule(ctx context.Context, in UpdateModuleInput) (*repository.Module, []*errorv1.ErrorItem) {
	in.ID = strings.TrimSpace(in.ID)
	if in.ID == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModuleIDRequired)}
	}
	mid, err := uuid.Parse(in.ID)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidModuleIDFormat)}
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModuleNameRequired)}
	}
	productID, err := parseOptionalUUID(in.ProductID)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidModuleProductID)}
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
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModuleNotFound.Messagef(in.ID))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUpdateModuleInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// DeleteModule performs soft-delete for a module.
func (s *ModuleService) DeleteModule(ctx context.Context, id string, actorUserID string) []*errorv1.ErrorItem {
	id = strings.TrimSpace(id)
	if id == "" {
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModuleIDRequired)}
	}
	mid, err := uuid.Parse(id)
	if err != nil {
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidModuleIDFormat)}
	}
	if _, err := s.q.DeleteModule(ctx, repository.DeleteModuleParams{
		ModuleID:      mid,
		DeletedUserID: pgtype.Text{String: actorUserID, Valid: strings.TrimSpace(actorUserID) != ""},
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrModuleNotFound.Messagef(id))}
		}
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrDeleteModuleInternal.Messagef(err.Error()))}
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
		logger.Error("cfg: cfg=%v", cfg)
		return nil, 0, 0, 0, constants.ErrMissingModulesQueryConfig.Status()
	}
	sorts, err := utils.NormalizeSorts(in.Sorts, cfg, 5)
	if err != nil {
		logger.Error("sorts: sorts=%v, err=%v", sorts, err)
		return nil, 0, 0, 0, constants.ErrInvalidModuleSort.Statusf(err)
	}
	filters, err := utils.NormalizeFilters(in.Filters, cfg, 10)
	if err != nil {
		logger.Error("Filters: filters=%v, err=%v", filters, err)
		return nil, 0, 0, 0, constants.ErrInvalidModuleFilter.Statusf(err)
	}
	logger.Info("Before run query")
	listSQL, listArgs, countSQL, countArgs := utils.BuildListAndCountQueries(cfg, sorts, filters, limit, offset)
	logger.Info("ListModules: listSQL=%s, listArgs=%v, countSQL=%s, countArgs=%v",
		listSQL, listArgs, countSQL, countArgs,
	)
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
