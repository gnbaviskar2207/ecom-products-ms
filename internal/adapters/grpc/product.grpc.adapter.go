package grpc

import (
	"log/slog"

	productsV1 "github.com/gnbaviskar2207/ecom-products-ms/gen/products"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/ports"
)

type ProductHandler struct {
	productsV1.UnimplementedProductServiceServer
	productService ports.ProductService
	logger         *slog.Logger
}

func New(logger *slog.Logger, productService ports.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService, logger: logger}
}
