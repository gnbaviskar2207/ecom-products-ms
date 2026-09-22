package ports

import (
	"context"

	"github.com/gnbaviskar2207/ecom-products-ms/internal/domain"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/dto"
)

type ProductRepository interface {
	FindOneByPid(ctx context.Context, pid string) (*domain.Product, error)
	ListProducts(ctx context.Context, req *dto.ListProductsRequestDTO) (*domain.PaginatedProducts, error)
}
