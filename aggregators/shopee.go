package aggregators

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/raihannurr/aggregator/utils"
)

type shopeeProductSearchResponse struct {
	Adjust struct {
		Count int `json:"count"`
	} `json:"adjust"`
	Items []shopeeItem `json:"items"`
}

type shopeeItem struct {
	ItemBasic struct {
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Image  string `json:"image"`
		ShopID int    `json:"shopid"`
		ItemID int    `json:"itemid"`
	} `json:"item_basic"`
}

type Shopee struct {
	HttpClient *http.Client
	Host       string // https://shopee.co.id
}

func (s Shopee) FetchProducts(params url.Values) ([]Product, int) {
	params.Add("by", "relevancy")
	params.Add("newest", params.Get("offset"))

	url := s.Host + "/api/v4/search/search_items?" + params.Encode()

	req, err := http.NewRequest("GET", url, nil)
	utils.Panic(err)

	res, err := s.HttpClient.Do(req)
	utils.Panic(err)

	var parsedResponse shopeeProductSearchResponse
	utils.Decode(res, &parsedResponse)

	return s.parseProducts(parsedResponse.Items), parsedResponse.Adjust.Count
}

func (s Shopee) parseProducts(items []shopeeItem) []Product {
	products := []Product{}

	for _, item := range items {
		product := Product{
			Name:  item.ItemBasic.Name,
			Price: item.ItemBasic.Price / 100000,
			Image: fmt.Sprintf("https://cf.shopee.co.id/file/%s", item.ItemBasic.Image),
			URL:   fmt.Sprintf("https://shopee.co.id/product-i.%d.%d", item.ItemBasic.ShopID, item.ItemBasic.ItemID),
		}

		products = append(products, product)
	}
	return products
}
