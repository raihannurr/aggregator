package utils

import (
	"io/ioutil"
	"net/http"
)

func ReadBody(res *http.Response) string {
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	Panic(err)

	return string(bodyBytes)
}
