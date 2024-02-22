package main

import (
	"context"
	"fmt"
	"time"

	"github.com/regismartiny/nuvemshop-product-stock-sync/configs"
	"github.com/regismartiny/nuvemshop-product-stock-sync/internal/nuvemshop"
)

func main() {
	config, _ := configs.LoadConfig(".")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	baseUrl := config.NuvemshopAPIBaseURL + "/" + config.NuvemshopStoreID
	client := nuvemshop.NewClient(baseUrl, config.NuvemshopAPIToken, config.NuvemshopUserAgent)

	products, err := client.GetProducts(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(products)
}
