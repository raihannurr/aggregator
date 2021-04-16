package aggregators

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/raihannurr/aggregator/utils"
)

type tokopediaProductSearchResponse struct {
	Data struct {
		AceSearchProductV4 struct {
			Data struct {
				Products []tokopediaItem `json:"products"`
			} `json:"data"`
			Header struct {
				TotalData int `json:"totalData"`
			} `json:"header"`
		} `json:"ace_search_product_v4"`
	} `json:"data"`
}

type tokopediaItem struct {
	Name  string `json:"name"`
	Price string `json:"price"`
	URL   string `json:"url"`
	Image string `json:"imageUrl"`
}

type searchPayload struct {
	OperationName string          `json:"operationName"`
	Query         string          `json:"query"`
	Variables     searchVariables `json:"variables"`
}

type searchVariables struct {
	Params string `json:"params"`
}

type Tokopedia struct {
	HttpClient *http.Client
	Host       string // gql.tokopedia.com
}

func (t Tokopedia) FetchProducts(params url.Values) ([]Product, int) {
	url := t.Host
	req, err := http.NewRequest("POST", url, t.searchPayload(params))
	utils.Panic(err)

	req.Header.Add("origin", "https://www.tokopedia.com")
	req.Header.Add("content-type", "application/json")

	res, err := t.HttpClient.Do(req)
	utils.Panic(err)

	var parsedResponse tokopediaProductSearchResponse
	utils.Decode(res, &parsedResponse)

	return t.parseProducts(parsedResponse.Data.AceSearchProductV4.Data.Products), parsedResponse.Data.AceSearchProductV4.Header.TotalData
}

func (t Tokopedia) searchPayload(params url.Values) *bytes.Buffer {
	body := new(bytes.Buffer)
	variables := url.Values{}
	variables.Add("device", "desktop")
	variables.Add("source", "universe")
	variables.Add("page", params.Get("page"))
	variables.Add("rows", params.Get("limit"))
	variables.Add("start", params.Get("offset"))
	variables.Add("q", params.Get("keyword"))

	err := json.NewEncoder(body).Encode(searchPayload{
		OperationName: "SearchProductQueryV4",
		Query:         searchQuery,
		Variables: searchVariables{
			Params: variables.Encode(),
		},
	})
	utils.Panic(err)

	return body
}

func (t Tokopedia) parseProducts(items []tokopediaItem) []Product {
	products := []Product{}

	for _, item := range items {
		product := Product{
			Name:  item.Name,
			Price: t.convertPrice(item.Price),
			Image: item.Image,
			URL:   item.URL,
		}

		products = append(products, product)
	}
	return products
}

func (t Tokopedia) convertPrice(price string) int {
	price = strings.Replace(price, "Rp", "", -1)
	price = strings.Replace(price, ".", "", -1)

	priceInt, err := strconv.Atoi(price)
	utils.Panic(err)

	return priceInt
}

const searchQuery = "query SearchProductQueryV4($params: String!) {\n  ace_search_product_v4(params: $params) {\n    header {\n      totalData\n      totalDataText\n      processTime\n      responseCode\n      errorMessage\n      additionalParams\n      keywordProcess\n      __typename\n    }\n    data {\n      isQuerySafe\n      ticker {\n        text\n        query\n        typeId\n        __typename\n      }\n      redirection {\n        redirectUrl\n        departmentId\n        __typename\n      }\n      related {\n        relatedKeyword\n        otherRelated {\n          keyword\n          url\n          product {\n            id\n            name\n            price\n            imageUrl\n            rating\n            countReview\n            url\n            priceStr\n            wishlist\n            shop {\n              city\n              isOfficial\n              isPowerBadge\n              __typename\n            }\n            ads {\n              adsId: id\n              productClickUrl\n              productWishlistUrl\n              shopClickUrl\n              productViewUrl\n              __typename\n            }\n            ratingAverage\n            labelGroups {\n              position\n              type\n              title\n              __typename\n            }\n            __typename\n          }\n          __typename\n        }\n        __typename\n      }\n      suggestion {\n        currentKeyword\n        suggestion\n        suggestionCount\n        instead\n        insteadCount\n        query\n        text\n        __typename\n      }\n      products {\n        id\n        name\n        ads {\n          adsId: id\n          productClickUrl\n          productWishlistUrl\n          productViewUrl\n          __typename\n        }\n        badges {\n          title\n          imageUrl\n          show\n          __typename\n        }\n        category: departmentId\n        categoryBreadcrumb\n        categoryId\n        categoryName\n        countReview\n        discountPercentage\n        gaKey\n        imageUrl\n        labelGroups {\n          position\n          title\n          type\n          __typename\n        }\n        originalPrice\n        price\n        priceRange\n        rating\n        ratingAverage\n        shop {\n          id\n          name\n          url\n          city\n          isOfficial\n          isPowerBadge\n          __typename\n        }\n        url\n        wishlist\n        sourceEngine: source_engine\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n"
