package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gnbaviskar2207/ecom-products-ms/internal/domain"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/ports"
)

type ProductService struct {
	logger *slog.Logger
	repo   ports.ProductRepository
}

func New(repo ports.ProductRepository, logger *slog.Logger) ports.ProductService {
	return &ProductService{repo: repo, logger: logger}
}

func (ps *ProductService) FindOneByPid(ctx context.Context, pid string) (domain.Product, error) {
	products, err := ps.repo.ListProducts(ctx, nil)
	if err != nil {
		ps.logger.ErrorContext(ctx, "failed to get paginated products", slog.String("error", err.Error()))
	} else {
		ps.logger.InfoContext(ctx, "products retrieved", slog.Any("products", products))
		fmt.Println("products", &products)
	}
	return ps.repo.FindOneByPid(ctx, pid)
}
