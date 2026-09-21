package domain

type Product struct {
	Pid string `bson:"pid"`
	// ProductName          []string `bson:"productName"`
	ProductNameEn string `bson:"productNameEn"`
	ProductSku    string `bson:"productSku"`
	ProductImage  string `bson:"productImage"`
	ProductWeight string `bson:"productWeight"`
	ProductType   string `bson:"productType"`
	// ProductUnit          string   `bson:"productUnit"`
	CategoryName string `bson:"categoryName"`
	// ListingCount         int64    `bson:"listingCount"`
	SellPrice string `bson:"sellPrice"`
	Remark    string `bson:"remark"`
	// AddMarkStatus        string   `bson:"addMarkStatus"`
	IsFreeShipping       bool     `bson:"isFreeShipping"`
	IsVideo              *bool    `bson:"isVideo"`
	SaleStatus           int64    `bson:"saleStatus"`
	ListedNum            int64    `bson:"listedNum"`
	SupplierName         *string  `bson:"supplierName"`
	SupplierId           string   `bson:"supplierId"`
	CategoryId           string   `bson:"categoryId"`
	SourceFrom           string   `bson:"sourceFrom"`
	ShippingCountryCodes []string `bson:"shippingCountryCodes"`
	ThreeCategoryName    *string  `bson:"threeCategoryName"`
	TwoCategoryId        *string  `bson:"twoCategoryId"`
	TwoCategoryName      *string  `bson:"twoCategoryName"`
	OneCategoryId        *string  `bson:"oneCategoryId"`
	OneCategoryName      *string  `bson:"oneCategoryName"`
	CustomizationVersion int64    `bson:"customizationVersion"`
	IsTestProduct        bool     `bson:"isTestProduct"`
	// CreateTime           int64    `bson:"createTime"`
}
