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

// RoleService provides business logic for roles.
type RoleService struct {
	q *repository.Queries
}

func NewRoleService(q *repository.Queries) *RoleService {
	return &RoleService{q: q}
}

type CreateRoleInput struct {
	Name        string
	Description string
	ModuleID    string
	ActorUserID string
}

type UpdateRoleInput struct {
	ID          string
	Name        string
	Description string
	ModuleID    string
	ActorUserID string
}

type ListRolesInput struct {
	Limit   int32
	Offset  int32
	Sorts   []*basev1.Sort
	Filters []*basev1.Filter
}

func (s *RoleService) CreateRole(ctx context.Context, in CreateRoleInput) (*repository.Role, []*errorv1.ErrorItem) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleNameRequired)}
	}
	if !utils.ValidateRoleName(name) {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleNameInvalid)}
	}

	if existing, err := s.q.GetRoleByName(ctx, name); err == nil && existing != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleNameAlreadyExists)}
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetRoleInternal.Messagef(err.Error()))}
	}

	moduleID := uuid.Nil
	if strings.TrimSpace(in.ModuleID) != "" {
		parsedModuleID, err := uuid.Parse(strings.TrimSpace(in.ModuleID))
		if err != nil {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidRoleIDFormat)}
		}
		moduleID = parsedModuleID
	}

	role, err := s.q.InsertRole(ctx, repository.InsertRoleParams{
		Name:          name,
		Description:   pgtype.Text{String: strings.TrimSpace(in.Description), Valid: strings.TrimSpace(in.Description) != ""},
		ModuleID:      moduleID,
		CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrCreateRoleInternal.Messagef(err.Error()))}
	}
	return role, nil
}

func (s *RoleService) GetRole(ctx context.Context, id string) (*repository.Role, []*errorv1.ErrorItem) {
	rid, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidRoleIDFormat)}
	}
	r, err := s.q.GetRole(ctx, rid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleNotFound.Messagef(id))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetRoleInternal.Messagef(err.Error()))}
	}
	return r, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, in UpdateRoleInput) (*repository.Role, []*errorv1.ErrorItem) {
	if strings.TrimSpace(in.ID) == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleIDRequired)}
	}
	rid, err := uuid.Parse(strings.TrimSpace(in.ID))
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidRoleIDFormat)}
	}

	name := strings.TrimSpace(in.Name)
	description := strings.TrimSpace(in.Description)
	moduleID := uuid.Nil
	if strings.TrimSpace(in.ModuleID) != "" {
		parsedModuleID, err := uuid.Parse(strings.TrimSpace(in.ModuleID))
		if err != nil {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidRoleIDFormat)}
		}
		moduleID = parsedModuleID
	}

	if name != "" {
		if !utils.ValidateRoleName(name) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleNameInvalid)}
		}
		if existing, err := s.q.GetRoleByName(ctx, name); err == nil && existing.RoleID != rid {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleNameAlreadyExists)}
		} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetRoleInternal.Messagef(err.Error()))}
		}
	}

	role, err := s.q.UpdateRole(ctx, repository.UpdateRoleParams{
		RoleID:        rid,
		Name:          name,
		Description:   pgtype.Text{String: description, Valid: description != ""},
		ModuleID:      moduleID,
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleNotFound.Messagef(in.ID))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUpdateRoleInternal.Messagef(err.Error()))}
	}
	return role, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id string, actorUserID string) []*errorv1.ErrorItem {
	rid, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidRoleIDFormat)}
	}
	if _, err := s.q.GetRole(ctx, rid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrRoleNotFound.Messagef(id))}
		}
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetRoleInternal.Messagef(err.Error()))}
	}
	if _, err := s.q.DeleteRole(ctx, repository.DeleteRoleParams{
		RoleID:        rid,
		DeletedUserID: pgtype.Text{String: actorUserID, Valid: strings.TrimSpace(actorUserID) != ""},
	}); err != nil {
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrDeleteRoleInternal.Messagef(err.Error()))}
	}
	return nil
}

// ListRoles returns paginated roles with reusable sort/filter validation.
func (s *RoleService) ListRoles(ctx context.Context, in ListRolesInput) ([]*repository.Role, int64, int32, int32, error) {
	limit := int32(20)
	offset := int32(0)
	if in.Limit > 0 {
		limit = in.Limit
	}
	if in.Offset >= 0 {
		offset = in.Offset
	}

	cfg, ok := utils.GetTableQueryConfig(utils.TableRoles)
	if !ok {
		return nil, 0, 0, 0, constants.ErrMissingRolesQueryConfig.Status()
	}
	sorts, err := utils.NormalizeSorts(in.Sorts, cfg, 5)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidRoleSort.Statusf(err)
	}
	filters, err := utils.NormalizeFilters(in.Filters, cfg, 10)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidRoleFilter.Statusf(err)
	}

	listSQL, listArgs, countSQL, countArgs := utils.BuildListAndCountQueries(cfg, sorts, filters, limit, offset)
	roles, total, err := s.q.ListRolesByDynamicQuery(ctx, repository.DynamicRoleListParams{
		ListSQL:   listSQL,
		ListArgs:  listArgs,
		CountSQL:  countSQL,
		CountArgs: countArgs,
	})
	if err != nil {
		return nil, 0, 0, 0, constants.ErrListRolesInternal.Statusf(err)
	}

	return roles, total, limit, offset, nil
}
