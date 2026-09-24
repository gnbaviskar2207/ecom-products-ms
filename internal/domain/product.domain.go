package domain

type Product struct {
	Pid string `json:"pid" bson:"pid"`
	// ProductName          []string `bson:"productName"`
	ProductNameEn string `json:"name" bson:"productNameEn"`
	ProductSku    string `json:"sku" bson:"sku"`
	ProductImage  string `json:"image" bson:"image"`
	ProductWeight string `json:"weight" bson:"weight"`
	ProductType   string `json:"type" bson:"type"`
	// ProductUnit          string   `bson:"productUnit"`
	CategoryName string `json:"category_name" bson:"categoryName"`
	// ListingCount         int64    `bson:"listingCount"`
	SellPrice string `json:"sell_price" bson:"sellPrice"`
	Remark    string `json:"remark" bson:"remark"`
	// AddMarkStatus        string   `bson:"addMarkStatus"`
	IsFreeShipping       bool     `json:"is_free_shipping" bson:"isFreeShipping"`
	IsVideo              *bool    `json:"is_video" bson:"isVideo"`
	SaleStatus           int64    `json:"sale_status" bson:"saleStatus"`
	ListedNum            int64    `json:"listed_num" bson:"listedNum"`
	SupplierName         *string  `json:"supplier_name" bson:"supplierName"`
	SupplierId           string   `json:"supplier_id" bson:"supplierId"`
	CategoryId           string   `json:"category_id" bson:"categoryId"`
	SourceFrom           string   `json:"source_from" bson:"sourceFrom"`
	ShippingCountryCodes []string `json:"shipping_country_codes" bson:"shippingCountryCodes"`
	ThreeCategoryName    *string  `json:"three_category_name" bson:"threeCategoryName"`
	TwoCategoryId        *string  `json:"two_category_id" bson:"twoCategoryId"`
	TwoCategoryName      *string  `json:"two_category_name" bson:"twoCategoryName"`
	OneCategoryId        *string  `json:"one_category_id" bson:"oneCategoryId"`
	OneCategoryName      *string  `json:"one_category_name" bson:"oneCategoryName"`
	CustomizationVersion int64    `json:"customization_version" bson:"customizationVersion"`
	IsTestProduct        bool     `json:"is_test_product" bson:"isTestProduct"`
	// CreateTime           int64    `bson:"createTime"`
}

type PaginatedResult[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}
