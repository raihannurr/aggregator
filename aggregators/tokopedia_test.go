package aggregators_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/raihannurr/aggregator/aggregators"
	"github.com/stretchr/testify/assert"
)

func Test_Tokopedia(t *testing.T) {
	client := http.Client{}
	tokopedia := aggregators.Tokopedia{
		Host:       "https://gql.tokopedia.com",
		HttpClient: &client,
	}

	params := url.Values{}
	params.Add("keyword", "nintendo switch")
	params.Add("limit", "3")

	products, total := tokopedia.FetchProducts(params)

	assert.NotNil(t, products)
	assert.Equal(t, 3, len(products))
	assert.Greater(t, total, 3)
}
