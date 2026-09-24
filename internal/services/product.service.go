package services

import (
	"context"
	"log/slog"

	"github.com/gnbaviskar2207/ecom-products-ms/internal/domain"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/dto"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/ports"
)

type ProductService struct {
	logger *slog.Logger
	repo   ports.ProductRepository
}

func New(repo ports.ProductRepository, logger *slog.Logger) ports.ProductService {
	return &ProductService{repo: repo, logger: logger}
}

func (p *ProductService) FindOneByPid(ctx context.Context, pid string) (*domain.Product, error) {
	return p.repo.FindOneByPid(ctx, pid)
}

func (p *ProductService) ListProducts(ctx context.Context, req *dto.ListProductsRequestDTO) (*domain.PaginatedProducts, error) {
	products, err := p.repo.ListProducts(ctx, req)
	if err != nil {
		p.logger.ErrorContext(ctx, "failed to get paginated products", slog.String("error", err.Error()))
		return nil, err
	}
	return products, nil
}
