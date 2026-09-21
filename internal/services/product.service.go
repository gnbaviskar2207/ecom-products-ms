package services

import (
	"context"
	"log/slog"

	"github.com/gnbaviskar2207/ecom-products-ms/internal/domain"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/ports"
)

type ProductService struct {
	logger *slog.Logger
	repo   ports.ProductRepository
}

func New(repo ports.ProductRepository) ports.ProductService {
	return &ProductService{repo: repo}
}

func (ps *ProductService) FindOneByPid(ctx context.Context, pid string) (domain.Product, error) {
	return ps.repo.FindOneByPid(ctx, pid)
}
