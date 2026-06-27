package grpc

import (
	"context"
	"fmt"
	"math"
	"time"

	"go-api/internal/repository"
	productsvc "go-api/internal/service"
	"go-api/pkg/constants"
	"go-api/pkg/logger"
	"go-api/pkg/utils"

	basev1 "github.com/lamquangmanh/protobuf/gen/go/proto/base/v1"
	errorv1 "github.com/lamquangmanh/protobuf/gen/go/proto/error/v1"
	productv1 "github.com/lamquangmanh/protobuf/gen/go/proto/product/v1"
	"google.golang.org/grpc/codes"
)

type ProductHandler struct {
	productv1.UnimplementedProductServiceServer
	productService *productsvc.ProductService
}

// NewProductHandler creates a gRPC transport handler for ProductService.
func NewProductHandler(productService *productsvc.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// productRepoToProto maps repository product model to protobuf response model.
func productRepoToProto(item *repository.Product) *productv1.Product {
	url := ""
	if item.Url.Valid {
		url = item.Url.String
	}
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

	return &productv1.Product{
		ProductId:     item.ProductID.String(),
		Name:          item.Name,
		Description:   description,
		Url:           url,
		CreatedAt:     item.CreatedAt.Time.UTC().Format(time.RFC3339),
		CreatedUserId: createdUserID,
		UpdatedAt:     updatedAt,
		UpdatedUserId: updatedUserID,
		DeletedAt:     deletedAt,
		DeletedUserId: deletedUserID,
	}
}

// GetProduct returns a single product by ID.
func (h *ProductHandler) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
	item, errors := h.productService.GetProduct(ctx, req.GetProductId())
	if errors != nil {
		logger.Error("GetProduct request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}

	return &productv1.GetProductResponse{
		Product: productRepoToProto(item),
	}, nil
}

// GetProducts returns paginated products with validated sorts and filters.
func (h *ProductHandler) GetProducts(ctx context.Context, req *productv1.GetProductsRequest) (*productv1.GetProductsResponse, error) {
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

	// Clone proto sorts and filters directly
	sorts := make([]*basev1.Sort, 0, len(req.GetSorts()))
	sorts = append(sorts, req.GetSorts()...)

	filters := make([]*basev1.Filter, 0, len(req.GetFilters()))
	filters = append(filters, req.GetFilters()...)
	logger.Info("GetProducts request: %d", req.GetFilters())
	logger.Info("Received GetProducts request: limit=%d, page=%d, offset=%d, sorts=%v, filters=%v", limit, page, offset, sorts, filters)

	items, total, appliedLimit, appliedOffset, err := h.productService.ListProducts(ctx, productsvc.ListProductsInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   sorts,
		Filters: filters,
	})
	if err != nil {
		return &productv1.GetProductsResponse{Data: nil, Pagination: &basev1.PaginationResponse{}}, nil
	}
	data := make([]*productv1.Product, 0, len(items))
	for _, item := range items {
		data = append(data, productRepoToProto(item))
	}
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}
	logger.Info("ListProducts: total=%d, limit=%d, offset=%d, totalPages=%d, returnedItems=%d", total, appliedLimit, appliedOffset, totalPages, len(data))

	response := &productv1.GetProductsResponse{
		Data: data,
		Pagination: &basev1.PaginationResponse{
			Page:       (appliedOffset / appliedLimit) + 1,
			Limit:      appliedLimit,
			TotalItems: int32(total),
			TotalPages: totalPages,
			ItemCount:  int32(len(data)),
		},
	}
	logger.Info(fmt.Sprintf("GetProducts response: %+v", response))
	return response, nil
}

// CreateProduct validates payload and creates a new product.
func (h *ProductHandler) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductResponse, error) {
	if req.GetProduct() == nil {
		logger.Error("CreateProduct request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductPayloadRequired)})
		return nil, err
	}
	item, errors := h.productService.CreateProduct(ctx, productsvc.CreateProductInput{
		Name:        req.GetProduct().GetName(),
		Description: req.GetProduct().GetDescription(),
		Url:         req.GetProduct().GetUrl(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		logger.Error("CreateProduct request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &productv1.CreateProductResponse{
		Product: productRepoToProto(item),
	}, nil
}

// UpdateProduct updates an existing product.
func (h *ProductHandler) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.UpdateProductResponse, error) {
	if req.GetProduct() == nil {
		logger.Error("UpdateProduct request is nil")
		_, err := utils.ResponseError(codes.InvalidArgument, []*errorv1.ErrorItem{utils.ErrMsg(constants.ErrProductPayloadRequired)})
		return nil, err
	}
	_, errors := h.productService.UpdateProduct(ctx, productsvc.UpdateProductInput{
		ID:          req.GetProduct().GetProductId(),
		Name:        req.GetProduct().GetName(),
		Description: req.GetProduct().GetDescription(),
		Url:         req.GetProduct().GetUrl(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		logger.Error("UpdateProduct request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &productv1.UpdateProductResponse{
		Success: true,
	}, nil
}

// DeleteProduct performs a soft-delete for a product.
func (h *ProductHandler) DeleteProduct(ctx context.Context, req *productv1.DeleteProductRequest) (*productv1.DeleteProductResponse, error) {
	if errors := h.productService.DeleteProduct(ctx, req.GetProductId(), req.GetUserId()); errors != nil {
		logger.Error("DeleteProduct request has errors")
		_, err := utils.ResponseError(codes.InvalidArgument, errors)
		return nil, err
	}
	return &productv1.DeleteProductResponse{
		Success: true,
	}, nil
}
