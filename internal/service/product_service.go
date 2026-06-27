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
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
)

type ProductService struct {
	q *repository.Queries
}

// NewProductService creates a business service for product operations.
func NewProductService(q *repository.Queries) *ProductService {
	return &ProductService{q: q}
}

type CreateProductInput struct {
	Name        string
	Description string
	Url         string
	ActorUserID string
}

type UpdateProductInput struct {
	ID          string
	Name        string
	Description string
	Url         string
	ActorUserID string
}

// CreateProduct validates and creates a product.
func (s *ProductService) CreateProduct(ctx context.Context, in CreateProductInput) (*repository.Product, []*errorv1.ErrorItem) {
	name := strings.TrimSpace(in.Name)
	description := strings.TrimSpace(in.Description)
	url := strings.TrimSpace(in.Url)
	if name == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductNameRequired)}

	}
	item, err := s.q.InsertProduct(ctx, repository.InsertProductParams{
		Name:          name,
		Description:   pgtype.Text{String: description, Valid: description != ""},
		Url:           pgtype.Text{String: url, Valid: url != ""},
		CreatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrCreateProductInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// GetProduct returns a product by ID.
func (s *ProductService) GetProduct(ctx context.Context, id string) (*repository.Product, []*errorv1.ErrorItem) {
	logger.Info("GetProduct called with ID")
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductIDRequired, map[string]any{"product_id": "a non-empty string"})}
		// return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductIDRequired)}
	}

	pid, err := uuid.Parse(id)
	if err != nil {
		// return nil, constants.ErrInvalidProductIDFormat.Status()
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidProductIDFormat)}
	}
	item, err := s.q.GetProduct(ctx, pid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductNotFound.Messagef(id))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrGetProductInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// UpdateProduct validates input and updates a product.
func (s *ProductService) UpdateProduct(ctx context.Context, in UpdateProductInput) (*repository.Product, []*errorv1.ErrorItem) {
	in.ID = strings.TrimSpace(in.ID)
	if in.ID == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductIDRequired)}
	}
	pid, err := uuid.Parse(in.ID)
	if err != nil {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidProductIDFormat)}
	}
	name := strings.TrimSpace(in.Name)
	description := strings.TrimSpace(in.Description)
	url := strings.TrimSpace(in.Url)
	if name == "" {
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductNameRequired)}
	}
	item, err := s.q.UpdateProduct(ctx, repository.UpdateProductParams{
		ProductID:     pid,
		Name:          name,
		Description:   pgtype.Text{String: description, Valid: description != ""},
		Url:           pgtype.Text{String: url, Valid: url != ""},
		UpdatedUserID: pgtype.Text{String: in.ActorUserID, Valid: strings.TrimSpace(in.ActorUserID) != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductNotFound.Messagef(in.ID))}
		}
		return nil, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrUpdateProductInternal.Messagef(err.Error()))}
	}
	return item, nil
}

// DeleteProduct performs soft-delete for a product.
func (s *ProductService) DeleteProduct(ctx context.Context, id string, actorUserID string) []*errorv1.ErrorItem {
	id = strings.TrimSpace(id)
	if id == "" {
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductIDRequired)}
	}
	pid, err := uuid.Parse(id)
	if err != nil {
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrInvalidProductIDFormat)}
	}
	if _, err := s.q.DeleteProduct(ctx, repository.DeleteProductParams{
		ProductID:     pid,
		DeletedUserID: pgtype.Text{String: actorUserID, Valid: strings.TrimSpace(actorUserID) != ""},
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductNotFound.Messagef(id))}
		}
		return []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrDeleteProductInternal.Messagef(err.Error()))}
	}
	return nil
}

type ListProductsInput struct {
	Limit   int32
	Offset  int32
	Sorts   []*basev1.Sort
	Filters []*basev1.Filter
}

// ListProducts returns a paginated product list using centralized query metadata
// and reusable sort/filter normalization.
//
// Behavior overview:
//   - Applies defaults:
//   - limit defaults to 20 when input limit <= 0
//   - offset defaults to 0 when input offset < 0
//   - Loads table query config via utils.GetTableQueryConfig(TableProducts).
//   - Validates and normalizes sorting rules via utils.NormalizeSorts:
//   - rejects unknown or non-sortable fields
//   - enforces maximum 5 sort fields
//   - Validates and normalizes filters via utils.NormalizeFilters:
//   - rejects unknown fields or unsupported operators
//   - parses typed values (text/uuid/timestamp)
//   - enforces maximum 10 filter fields
//   - Builds dynamic list/count SQL using utils.BuildListAndCountQueries.
//   - Executes dynamic list/count query through repository layer.
//
// Return values:
//   - items: matched products for the requested page
//   - total: total number of matched rows (without limit/offset)
//   - limit: effective limit after defaulting
//   - offset: effective offset after defaulting
//   - error: gRPC status error (InvalidArgument or Internal)
//
// Sample input:
//
//	ListProductsInput{
//	  Limit: 10,
//	  Offset: 0,
//	  Sorts:   [{Field:"name", Order:ASC}],
//	  Filters: [{Field:"name", Operator:LIKE, StringValue:"crm"}],
//	}
//
// Sample output:
//
//	items=[{product_id:..., name:"CRM", ...}], total=1, limit=10, offset=0, err=nil
func (s *ProductService) ListProducts(ctx context.Context, in ListProductsInput) ([]*repository.Product, int64, int32, int32, error) {
	limit := in.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := in.Offset
	if offset < 0 {
		offset = 0
	}
	cfg, ok := utils.GetTableQueryConfig(utils.TableProducts)
	if !ok {
		logger.Error("missing query config for products table")
		return nil, 0, 0, 0, constants.ErrMissingProductsQueryConfig.Status()
	}

	sorts, err := utils.NormalizeSorts(in.Sorts, cfg, 5)
	if err != nil {
		logger.Error("invalid product sort configuration: %v", err)
		return nil, 0, 0, 0, constants.ErrInvalidProductSort.Statusf(err)
	}

	filters, err := utils.NormalizeFilters(in.Filters, cfg, 10)
	if err != nil {
		logger.Error("invalid product filter configuration: %v", err)
		return nil, 0, 0, 0, constants.ErrInvalidProductFilter.Statusf(err)
	}

	// logger.Info("Executing filters: %s", filters)
	listSQL, listArgs, countSQL, countArgs := utils.BuildListAndCountQueries(cfg, sorts, filters, limit, offset)
	// logger.Info("Executing dynamic product list query: %s", listSQL)
	items, total, err := s.q.ListProductsByDynamicQuery(ctx, repository.DynamicProductListParams{
		ListSQL:   listSQL,
		ListArgs:  listArgs,
		CountSQL:  countSQL,
		CountArgs: countArgs,
	})
	if err != nil {
		return nil, 0, 0, 0, constants.ErrListProductsInternal.Statusf(err)
	}

	return items, total, limit, offset, nil
}
