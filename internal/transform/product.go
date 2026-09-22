package transform

import (
	productsV1 "github.com/gnbaviskar2207/ecom-products-ms/gen/products"
	"github.com/gnbaviskar2207/ecom-products-ms/internal/domain"
)

// goverter:converter
// goverter:useZeroValueOnPointerInconsistency
type Converter interface {
	// goverter:ignore state unknownFields sizeCache
	// goverter:map ProductNameEn Name
	// goverter:map ProductSku Sku
	// goverter:map ProductImage Image
	// goverter:map ProductWeight Weight
	// goverter:map ProductType Type
	// goverter:map CategoryName CategoryName
	// goverter:map SellPrice SellPrice
	// goverter:map Remark Remark
	// goverter:map IsFreeShipping IsFreeShipping
	// goverter:map IsVideo IsVideo
	// goverter:map SaleStatus SaleStatus
	// goverter:map ListedNum ListedNum
	// goverter:map SupplierName SupplierName
	// goverter:map SupplierId SupplierId
	// goverter:map CategoryId CategoryId
	// goverter:map SourceFrom SourceFrom
	// goverter:map ShippingCountryCodes ShippingCountryCodes
	// goverter:map ThreeCategoryName ThreeCategoryName
	// goverter:map TwoCategoryId TwoCategoryId
	// goverter:map TwoCategoryName TwoCategoryName
	// goverter:map OneCategoryId OneCategoryId
	// goverter:map OneCategoryName OneCategoryName
	// goverter:map CustomizationVersion CustomizationVersion
	// goverter:map IsTestProduct IsTestProduct
	ToProductPb(p *domain.Product) *productsV1.Product

	ToProductsPb(products []*domain.Product) []*productsV1.Product
}
