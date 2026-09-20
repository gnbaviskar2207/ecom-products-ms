package ports

import (
	"context"

	"github.com/gnbaviskar2207/ecom-products-ms/internal/domain"
)

type ProductRepository interface {
	FindOneByPid(ctx context.Context, pid string) (domain.Product, error)
}
