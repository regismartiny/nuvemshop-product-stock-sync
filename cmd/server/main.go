package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/regismartiny/nuvemshop-product-stock-sync/configs"
	"github.com/regismartiny/nuvemshop-product-stock-sync/internal/db"
	"github.com/regismartiny/nuvemshop-product-stock-sync/internal/nuvemshop"
)

func main() {
	config, _ := configs.LoadConfig(".")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	baseUrl, err := url.Parse(config.NuvemshopAPIBaseURL)
	if err != nil {
		log.Fatal(err)
	}

	baseUrl.Path = path.Join(baseUrl.Path, config.NuvemshopStoreID)
	client := nuvemshop.NewClient(baseUrl, config.NuvemshopAPIToken, config.NuvemshopUserAgent)

	updateProducts(&ctx, client)
}

func updateProducts(ctx *context.Context, client *nuvemshop.Client) {

	products, err := client.GetProducts(ctx, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Products in Nuvemshop")
	log.Println(products)

	productsUpdateStockPrice := make([]nuvemshop.ProductsUpdateStockPrice, 0)

	for _, product := range products {

		for _, variant := range product.Variants {

			variantsUpdateStockPrice := make([]nuvemshop.VariantsUpdateStockPrice, 0)

			variantsUpdateStockPrice = verifyVariantToUpdateStockPrice(variant, variantsUpdateStockPrice)

			productUpdateStockPrice := nuvemshop.ProductsUpdateStockPrice{
				ID:       product.ID,
				Variants: variantsUpdateStockPrice,
			}

			productsUpdateStockPrice = append(productsUpdateStockPrice, productUpdateStockPrice)
		}
	}

	log.Println("Products to update")
	log.Println(productsUpdateStockPrice)

	if len(productsUpdateStockPrice) == 0 {
		return
	}

	updateProductsInNuvemshop(ctx, client, productsUpdateStockPrice)
}

func verifyVariantToUpdateStockPrice(nuvemshopVariant nuvemshop.Variant, variantsUpdateStockPrice []nuvemshop.VariantsUpdateStockPrice) []nuvemshop.VariantsUpdateStockPrice {
	dbVariantStock := db.GetProductVariantStockBySKU(nuvemshopVariant.SKU)

	if dbVariantStock != nuvemshopVariant.Stock {

		priceFloat, err := strconv.ParseFloat(nuvemshopVariant.Price, 64)
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing price of variant with SKU %s:", nuvemshopVariant.SKU), err)
			return variantsUpdateStockPrice
		}

		variantUpdateStockPrice := nuvemshop.VariantsUpdateStockPrice{
			ID:    nuvemshopVariant.ID,
			Price: priceFloat,
			InventoryLevels: []nuvemshop.InventoryLevels{
				{Stock: dbVariantStock},
			},
		}

		variantsUpdateStockPrice = append(variantsUpdateStockPrice, variantUpdateStockPrice)
	}

	return variantsUpdateStockPrice
}

func updateProductsInNuvemshop(ctx *context.Context, client *nuvemshop.Client, productsUpdateStockPrice []nuvemshop.ProductsUpdateStockPrice) {
	log.Println("Updating Products in Nuvemshop")

	updateResponse, err := client.UpdateProductStock(ctx, productsUpdateStockPrice)

	if err != nil {
		log.Println("Error updating products:", err)
		return
	}

	validateProductsUpdate(updateResponse)
}

func validateProductsUpdate(productUpdateResponse []nuvemshop.ProductsUpdateStockPriceResponse) {

	for _, productUpdate := range productUpdateResponse {
		for _, variant := range productUpdate.Variants {
			if !variant.Success {
				log.Println("error updating variant of id ", variant.ID)

			}
		}
	}

}
