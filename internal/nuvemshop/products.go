package nuvemshop

type Variant struct {
	ID    int    `json:"id"`
	SKU   string `json:"sku"`
	Stock int    `json:"stock"`
	Price string `json:"price"`
}

type Product struct {
	ID       int       `json:"id"`
	Variants []Variant `json:"variants"`
}

type ProductsListOptions struct {
	Page     int
	Per_page int
}

type InventoryLevels struct {
	Stock int `json:"stock"`
}

type VariantsUpdateStockPrice struct {
	ID              int               `json:"id"`
	Price           float64           `json:"price"`
	InventoryLevels []InventoryLevels `json:"inventory_levels"`
}

type ProductsUpdateStockPrice struct {
	ID       int                        `json:"id"`
	Variants []VariantsUpdateStockPrice `json:"variants"`
}

type ProductsUpdateStockPriceResponse struct {
	ID       int                                `json:"id"`
	Variants []VariantsUpdateStockPriceResponse `json:"variants"`
}

type VariantsUpdateStockPriceResponse struct {
	ID      int  `json:"id"`
	Success bool `json:"success"`
}
