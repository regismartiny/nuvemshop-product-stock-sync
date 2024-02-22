package nuvemshop

type Variant struct {
	ID    int    `json:"id"`
	SKU   string `json:"sku"`
	Stock int    `json:"stock"`
}

type Product struct {
	ID       int       `json:"id"`
	Variants []Variant `json:"variants"`
}

type ProductsListOptions struct {
	Page     int
	Per_page int
}
