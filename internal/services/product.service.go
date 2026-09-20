package services

import (
	"context"

	"github.com/gnbaviskar2207/ecom-products-ms/internal/domain"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/ports"
)

type ProductService struct {
	repo ports.ProductRepository
}

func ProductServiceNew(repo ports.ProductRepository) ports.ProductService {
	return &ProductService{repo: repo}
}

func (ps *ProductService) FindOneByPid(ctx context.Context, pid string) (domain.Product, error) {
	return ps.repo.FindOneByPid(ctx, pid)
}
