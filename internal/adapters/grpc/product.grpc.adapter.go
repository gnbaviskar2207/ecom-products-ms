package grpc

import (
	"context"
	"log/slog"

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

func (p *ProductAdapter) FindOneByPid(ctx context.Context, request *productsV1.FindOneByPidRequest) (*productsV1.Product, error) {
	p.logger.DebugContext(ctx, "FindOneByPid request received", slog.String("method", "grpc.adapter.FindOneByPid"), slog.Any("request", request))

	product, err := p.productService.FindOneByPid(ctx, request.Pid)
	if err != nil {
		return nil, err
	}
	return p.transform.ToProductPb(product), nil
}

func (p *ProductAdapter) ListProducts(ctx context.Context, req *productsV1.ListProductsRequest) (*productsV1.ListProductsResponse, error) {
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
		Products:   p.transform.ToProductsPb(products.Products),
	}, nil
	// return p.transform.ToProductsResponsePb(products), nil
}
