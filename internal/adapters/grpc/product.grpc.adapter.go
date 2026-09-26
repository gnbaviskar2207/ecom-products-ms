package grpc

import (
	"context"
	"fmt"
	"log/slog"

	errs "github.com/gnbaviskar2207/ecom-common/pkg/err"
	productsV1 "github.com/gnbaviskar2207/ecom-products-ms/gen/products"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/dto"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/ports"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/transform"
)

type ProductAdapter struct {
	productsV1.UnimplementedProductServiceServer
	productService ports.ProductService
	transform      transform.Converter
	logger         *slog.Logger
}

func New(logger *slog.Logger, productService ports.ProductService, transform transform.Converter) *ProductAdapter {
	return &ProductAdapter{productService: productService, logger: logger, transform: transform}
}

func (p *ProductAdapter) FindOneByPid(ctx context.Context, request *productsV1.FindOneByPidRequest) (*productsV1.FindOneByPidResponse, error) {

	if err := request.Validate(); err != nil {
		p.logger.ErrorContext(ctx, " validation failed", slog.String("method", "grpc.adapter.FindOneByPid"), slog.Any("error", err))
		return nil, fmt.Errorf("%w: %v", errs.ErrInvalidArgument, err)
	}
	product, err := p.productService.FindOneByPid(ctx, request.Pid)
	if err != nil {
		return nil, err
	}
	return &productsV1.FindOneByPidResponse{Product: p.transform.ToProductPb(product)}, nil
}

func (p *ProductAdapter) ListProducts(ctx context.Context, req *productsV1.ListProductsRequest) (*productsV1.ListProductsResponse, error) {

	if err := req.Validate(); err != nil {
		p.logger.ErrorContext(ctx, " validation failed", slog.String("method", "grpc.adapter.ListProducts"), slog.Any("error", err))
		return nil, fmt.Errorf("%w: %v", errs.ErrInvalidArgument, err)
	}
	products, err := p.productService.ListProducts(ctx, &dto.ListProductsRequestDTO{
		NextCursor: req.NextCursor,
		Limit:      req.Limit,
	})
	if err != nil {
		return nil, err
	}
	return &productsV1.ListProductsResponse{
		NextCursor: products.NextCursor,
		HasMore:    products.HasMore,
		Products:   p.transform.ToProductsPb(products.Data),
	}, nil
}
