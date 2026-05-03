package grpc

import (
	"context"
	"math"
	"time"

	"go-api/internal/repository"
	productsvc "go-api/internal/service"
	basepb "go-api/pkg/api/basepb"
	productpb "go-api/pkg/api/productpb"
	"go-api/pkg/constants"
	"go-api/pkg/utils"
)

type ProductHandler struct {
	productpb.UnimplementedProductServiceServer
	productService *productsvc.ProductService
}

// NewProductHandler creates a gRPC transport handler for ProductService.
func NewProductHandler(productService *productsvc.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// productRepoToProto maps repository product model to protobuf response model.
func productRepoToProto(item *repository.Product) *productpb.Product {
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

	return &productpb.Product{
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
func (h *ProductHandler) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	item, errors := h.productService.GetProduct(ctx, req.GetProductId())
	if errors != nil {
		return &productpb.GetProductResponse{
			Errors: errors,
		}, nil
	}

	return &productpb.GetProductResponse{
		Product: productRepoToProto(item),
	}, nil
}

// GetProducts returns paginated products with validated sorts and filters.
func (h *ProductHandler) GetProducts(ctx context.Context, req *productpb.GetProductsRequest) (*productpb.GetProductsResponse, error) {
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
	sorts := make([]*basepb.Sort, 0, len(req.GetSorts()))
	sorts = append(sorts, req.GetSorts()...)

	filters := make([]*basepb.Filter, 0, len(req.GetFilters()))
	filters = append(filters, req.GetFilters()...)

	items, total, appliedLimit, appliedOffset, err := h.productService.ListProducts(ctx, productsvc.ListProductsInput{
		Limit:   limit,
		Offset:  offset,
		Sorts:   sorts,
		Filters: filters,
	})
	if err != nil {
		return &productpb.GetProductsResponse{Data: nil, Pagination: &basepb.PaginationResponse{}}, nil
	}
	data := make([]*productpb.Product, 0, len(items))
	for _, item := range items {
		data = append(data, productRepoToProto(item))
	}
	totalPages := int32(math.Ceil(float64(total) / float64(appliedLimit)))
	if appliedLimit <= 0 {
		totalPages = 0
	}

	return &productpb.GetProductsResponse{
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

// CreateProduct validates payload and creates a new product.
func (h *ProductHandler) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateSuccess, error) {
	if req.GetProduct() == nil {
		return &productpb.CreateSuccess{
			Errors: []*basepb.ErrorMessage{utils.ErrMsg(constants.ErrProductPayloadRequired)},
		}, nil
	}
	item, errors := h.productService.CreateProduct(ctx, productsvc.CreateProductInput{
		Name:        req.GetProduct().GetName(),
		Description: req.GetProduct().GetDescription(),
		Url:         req.GetProduct().GetUrl(),
		ActorUserID: req.GetUserId(),
	})
	if errors != nil {
		return &productpb.CreateSuccess{
			Errors: errors,
		}, nil
	}
	return &productpb.CreateSuccess{
		Product: productRepoToProto(item),
	}, nil
}

// UpdateProduct updates an existing product.
func (h *ProductHandler) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*basepb.UpdateSuccess, error) {
	if req.GetProduct() == nil {
		return utils.UpdateErr(utils.ErrMsg(constants.ErrProductPayloadRequired)), nil
	}
	_, err := h.productService.UpdateProduct(ctx, productsvc.UpdateProductInput{
		ID:          req.GetProduct().GetProductId(),
		Name:        req.GetProduct().GetName(),
		Description: req.GetProduct().GetDescription(),
		Url:         req.GetProduct().GetUrl(),
		ActorUserID: req.GetUserId(),
	})
	if err != nil {
		return utils.UpdateErr(err...), nil
	}
	return utils.UpdateOK(), nil
}

// DeleteProduct performs a soft-delete for a product.
func (h *ProductHandler) DeleteProduct(ctx context.Context, req *productpb.DeleteProductRequest) (*basepb.DeleteSuccess, error) {
	if err := h.productService.DeleteProduct(ctx, req.GetProductId(), req.GetUserId()); err != nil {
		return utils.DeleteErr(err...), nil
	}
	return utils.DeleteOK(), nil
}
