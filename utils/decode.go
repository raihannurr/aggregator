package utils

import (
	"encoding/json"
	"net/http"
)

func Decode(res *http.Response, decoded interface{}) error {
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()

	return decoder.Decode(&decoded)
}
