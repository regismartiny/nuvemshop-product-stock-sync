package nuvemshop

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL   *url.URL
	apiKey    string
	UserAgent string
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(baseUrl *url.URL, apiKey string, userAgent string) *Client {
	return &Client{
		BaseURL:   baseUrl,
		apiKey:    apiKey,
		UserAgent: userAgent,
	}
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authentication", fmt.Sprintf("bearer %s", c.apiKey))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetProducts(ctx *context.Context, options *ProductsListOptions) ([]Product, error) {
	page := 1
	per_page := 100
	if options != nil {
		page = options.Page
		per_page = options.Per_page
	}

	req, err := http.NewRequestWithContext(*ctx, "GET", fmt.Sprintf("%s/products?page=%d&per_page=%d", c.BaseURL, page, per_page), nil)
	if err != nil {
		return nil, err
	}

	var res []Product
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) UpdateProductStock(ctx *context.Context, productsUpdateStockPrice []ProductsUpdateStockPrice) ([]ProductsUpdateStockPriceResponse, error) {

	out, err := json.Marshal(productsUpdateStockPrice)
	if err != nil {
		log.Println("Error marshaling content:", err)
		return nil, err
	}

	log.Println("Request JSON:", string(out))

	req, err := http.NewRequestWithContext(*ctx, "PATCH", fmt.Sprintf("%s/products/stock-price", c.BaseURL), bytes.NewBuffer(out))
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	var res []ProductsUpdateStockPriceResponse
	if err := c.sendRequest(req, &res); err != nil {
		log.Println("Error sending request:", err)
		return nil, err
	}

	resout, err := json.Marshal(res)
	if err != nil {
		log.Println("Error marshalling content:", err)
		return nil, err
	}
	log.Println("Response JSON:", string(resout))

	return res, nil
}
