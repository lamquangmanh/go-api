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
	userpb "go-api/pkg/api/userpb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"
)

type UserService struct {
	q    *repository.Queries
	pool *pgxpool.Pool
}

func NewUserService(q *repository.Queries, pool *pgxpool.Pool) *UserService {
	return &UserService{q: q, pool: pool}
}

type CreateUserInput struct {
	Username    string
	Email       string
	Password    string
	Phone       string
	Avatar      string
	Status      userpb.UserStatus
	RoleIDs     []string
	ActorUserID string
}

type UpdateUserInput struct {
	ID          string
	Username    string
	Email       string
	Phone       string
	Avatar      string
	Status      userpb.UserStatus
	RoleIDs     []string
	ActorUserID string
}

type ListUsersInput struct {
	IDs     []string
	Limit   int32
	Offset  int32
	Sorts   []*basepb.Sort
	Filters []*basepb.Filter
}

func (s *UserService) CreateUser(ctx context.Context, in CreateUserInput) (*repository.User, []*basepb.ErrorMessage) {
	username := strings.TrimSpace(in.Username)
	email := strings.TrimSpace(in.Email)
	errMsgs := make([]*basepb.ErrorMessage, 0, 2)

	if username == "" {
		errMsgs = append(errMsgs, utils.ErrMsg(constants.ErrUsernameRequired))
	} else if !utils.ValidateUsername(username) {
		errMsgs = append(errMsgs, utils.ErrMsg(constants.ErrUsernameInvalid))
	}
	if email == "" {
		errMsgs = append(errMsgs, utils.ErrMsg(constants.ErrEmailRequired))
	} else if !utils.ValidateEmail(email) {
		errMsgs = append(errMsgs, utils.ErrMsg(constants.ErrEmailInvalid))
	}
	if len(errMsgs) > 0 {
		return nil, errMsgs
	}

	if u, err := s.q.GetUserByEmail(ctx, strings.ToLower(email)); err == nil && u != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrEmailAlreadyExists)}
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrGetUserInternal.Messagef(err.Error()))}
	}
	if u, err := s.q.GetUserByUsername(ctx, username); err == nil && u != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUsernameAlreadyExists)}
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrGetUserInternal.Messagef(err.Error()))}
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrBeginTransaction.Messagef(err.Error()))}
	}
	defer tx.Rollback(ctx)
	q := s.q.WithTx(tx)

	roleUUIDs := make([]uuid.UUID, 0, len(in.RoleIDs))
	for _, r := range in.RoleIDs {
		rid, err := uuid.Parse(strings.TrimSpace(r))
		if err != nil {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidRoleIDFormatIn.Messagef(r))}
		}
		roleUUIDs = append(roleUUIDs, rid)
	}
	if len(roleUUIDs) > 0 {
		existing, err := q.GetRolesByIDs(ctx, roleUUIDs)
		if err != nil {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrValidateRoles.Messagef(err.Error()))}
		}
		if len(existing) != len(roleUUIDs) {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrRoleIDsNotExist)}
		}
	}

	user, err := q.InsertUser(ctx, repository.InsertUserParams{
		UserName:      username,
		Email:         strings.ToLower(email),
		Password:      in.Password,
		Phone:         pgtype.Text{String: strings.TrimSpace(in.Phone), Valid: strings.TrimSpace(in.Phone) != ""},
		Avatar:        pgtype.Text{String: strings.TrimSpace(in.Avatar), Valid: strings.TrimSpace(in.Avatar) != ""},
		Status:        toRepositoryUserStatus(in.Status),
		CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrCreateUserInternal.Messagef(err.Error()))}
	}

	for _, rid := range roleUUIDs {
		if _, err := q.InsertUserRole(ctx, repository.InsertUserRoleParams{
			UserID:        user.UserID,
			RoleID:        rid,
			CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
			UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		}); err != nil {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrAddUserRoleMapping.Messagef(err.Error()))}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrCommitTransaction.Messagef(err.Error()))}
	}
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*repository.User, []*basepb.ErrorMessage) {
	uid, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidUserIDFormat)}
	}
	u, err := s.q.GetUser(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUserNotFound.Messagef(id))}
		}
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrGetUserInternal.Messagef(err.Error()))}
	}
	return u, nil
}

func (s *UserService) UpdateUser(ctx context.Context, in UpdateUserInput) (*repository.User, []*basepb.ErrorMessage) {
	if strings.TrimSpace(in.ID) == "" {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUserIDRequired)}
	}
	uid, err := uuid.Parse(strings.TrimSpace(in.ID))
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidUserIDFormat)}
	}

	roleUUIDs := make([]uuid.UUID, 0, len(in.RoleIDs))
	for _, r := range in.RoleIDs {
		rid, err := uuid.Parse(strings.TrimSpace(r))
		if err != nil {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidRoleIDFormatIn.Messagef(r))}
		}
		roleUUIDs = append(roleUUIDs, rid)
	}

	if strings.TrimSpace(in.Username) != "" && !utils.ValidateUsername(in.Username) {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUsernameInvalid)}
	}
	if strings.TrimSpace(in.Email) != "" && !utils.ValidateEmail(in.Email) {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrEmailInvalid)}
	}

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrBeginTransaction.Messagef(err.Error()))}
	}
	defer tx.Rollback(ctx)
	q := s.q.WithTx(tx)

	existing, err := q.GetUser(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUserNotFound.Messagef(in.ID))}
		}
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrGetUserInternal.Messagef(err.Error()))}
	}

	newUsername := existing.UserName
	if strings.TrimSpace(in.Username) != "" {
		newUsername = strings.TrimSpace(in.Username)
	}
	newEmail := existing.Email
	if strings.TrimSpace(in.Email) != "" {
		newEmail = strings.ToLower(strings.TrimSpace(in.Email))
	}
	newPhone := existing.Phone
	if strings.TrimSpace(in.Phone) != "" {
		newPhone = pgtype.Text{String: strings.TrimSpace(in.Phone), Valid: true}
	}
	newAvatar := existing.Avatar
	if strings.TrimSpace(in.Avatar) != "" {
		newAvatar = pgtype.Text{String: strings.TrimSpace(in.Avatar), Valid: true}
	}
	newStatus := existing.Status
	if in.Status != userpb.UserStatus_USER_STATUS_UNSPECIFIED {
		newStatus = toRepositoryUserStatus(in.Status)
	}

	user, err := q.UpdateUser(ctx, repository.UpdateUserParams{
		UserID:        uid,
		UserName:      newUsername,
		Email:         newEmail,
		Password:      existing.Password,
		Phone:         newPhone,
		Avatar:        newAvatar,
		Status:        newStatus,
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUserNotFound.Messagef(in.ID))}
		}
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUpdateUserInternal.Messagef(err.Error()))}
	}

	if err := q.DeleteUserRoleMappingsByUserID(ctx, uid); err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrClearUserRoleMappings.Messagef(err.Error()))}
	}
	for _, rid := range roleUUIDs {
		if _, err := q.InsertUserRole(ctx, repository.InsertUserRoleParams{
			UserID:        uid,
			RoleID:        rid,
			CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
			UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		}); err != nil {
			return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrAddUserRoleMapping.Messagef(err.Error()))}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrCommitTransaction.Messagef(err.Error()))}
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string, actorUserID string) []*basepb.ErrorMessage {
	uid, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidUserIDFormat)}
	}
	if _, err := s.q.GetUser(ctx, uid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUserNotFound.Messagef(id))}
		}
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrGetUserInternal.Messagef(err.Error()))}
	}
	if _, err := s.q.DeleteUser(ctx, repository.DeleteUserParams{
		UserID:        uid,
		DeletedUserID: pgtype.Text{String: actorUserID, Valid: strings.TrimSpace(actorUserID) != ""},
	}); err != nil {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrDeleteUserInternal.Messagef(err.Error()))}
	}
	return nil
}

// ListUsers returns paginated users with reusable sort/filter validation.
func (s *UserService) ListUsers(ctx context.Context, in ListUsersInput) ([]*repository.User, int64, int32, int32, error) {
	limit := int32(20)
	offset := int32(0)
	if in.Limit > 0 {
		limit = in.Limit
	}
	if in.Offset >= 0 {
		offset = in.Offset
	}

	cfg, ok := utils.GetTableQueryConfig(utils.TableUsers)
	if !ok {
		return nil, 0, 0, 0, constants.ErrMissingUsersQueryConfig.Status()
	}
	sorts, err := utils.NormalizeSorts(in.Sorts, cfg, 5)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidUserSort.Statusf(err)
	}
	filters, err := utils.NormalizeFilters(in.Filters, cfg, 10)
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidUserFilter.Statusf(err)
	}

	listSQL, listArgs, countSQL, countArgs := utils.BuildListAndCountQueries(cfg, sorts, filters, limit, offset)
	users, total, err := s.q.ListUsersByDynamicQuery(ctx, repository.DynamicUserListParams{
		ListSQL:   listSQL,
		ListArgs:  listArgs,
		CountSQL:  countSQL,
		CountArgs: countArgs,
	})
	if err != nil {
		return nil, 0, 0, 0, constants.ErrListUsersInternal.Statusf(err)
	}

	return users, total, limit, offset, nil
}

func (s *UserService) ListUserRoles(ctx context.Context, userID string, limit int32, offset int32) ([]*repository.Role, int64, int32, int32, error) {
	uid, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return nil, 0, 0, 0, constants.ErrInvalidUserIDFormat.Status()
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	roles, err := s.q.ListRolesByUser(ctx, repository.ListRolesByUserParams{UserID: uid, Limit: limit, Offset: offset})
	if err != nil {
		return nil, 0, 0, 0, constants.ErrGetUserInternal.Statusf(err)
	}
	total := int64(len(roles))
	return roles, total, limit, offset, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID string, password string) []*basepb.ErrorMessage {
	uid, err := uuid.Parse(strings.TrimSpace(userID))
	if err != nil {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrInvalidUserIDFormat)}
	}
	existing, err := s.q.GetUser(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUserNotFound.Messagef(userID))}
		}
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrGetUserInternal.Messagef(err.Error()))}
	}
	if strings.TrimSpace(password) == "" {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrPasswordRequired)}
	}
	if _, err := s.q.UpdateUser(ctx, repository.UpdateUserParams{
		UserID:        uid,
		UserName:      existing.UserName,
		Email:         existing.Email,
		Password:      password,
		Phone:         existing.Phone,
		Avatar:        existing.Avatar,
		Status:        existing.Status,
		UpdatedUserID: pgtype.Text{String: userID, Valid: true},
	}); err != nil {
		return []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrUpdateUserInternal.Messagef(err.Error()))}
	}
	return nil
}

func toRepositoryUserStatus(statusValue userpb.UserStatus) repository.UserStatus {
	switch statusValue {
	case userpb.UserStatus_USER_STATUS_DEACTIVATED:
		return repository.UserStatusDEACTIVATED
	case userpb.UserStatus_USER_STATUS_DELETED:
		return repository.UserStatusDELETED
	default:
		return repository.UserStatusACTIVE
	}
}
