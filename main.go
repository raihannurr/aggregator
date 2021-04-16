package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/raihannurr/aggregator/aggregators"
	"github.com/raihannurr/aggregator/utils"
)

type Application struct {
	Shopee    aggregators.Shopee
	Tokopedia aggregators.Tokopedia
}

type Response struct {
	Data interface{}            `json:"data"`
	Meta map[string]interface{} `json:"meta"`
}

var app Application

func init() {
	client := http.Client{}

	app.Shopee = aggregators.Shopee{
		Host:       "https://shopee.co.id",
		HttpClient: &client,
	}

	app.Tokopedia = aggregators.Tokopedia{
		Host:       "https://gql.tokopedia.com",
		HttpClient: &client,
	}
}

func shopeeHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	products, total := app.Shopee.FetchProducts(params)

	response, err := json.Marshal(Response{
		Data: products,
		Meta: map[string]interface{}{
			"total":  total,
			"limit":  params.Get("limit"),
			"offset": params.Get("offset"),
		},
	})
	utils.Panic(err)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(response))
}

func tokopediaHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	products, total := app.Tokopedia.FetchProducts(params)

	response, err := json.Marshal(Response{
		Data: products,
		Meta: map[string]interface{}{
			"total":  total,
			"limit":  params.Get("limit"),
			"offset": params.Get("offset"),
		},
	})
	utils.Panic(err)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(response))
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	tmplPath := filepath.Join("templates", "home.html")
	tmpl := template.Must(template.ParseFiles(tmplPath))
	data := "La Chartreuse"
	tmpl.Execute(w, data)
}

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/shopee", shopeeHandler)
	http.HandleFunc("/tokped", tokopediaHandler)
	http.HandleFunc("/", homeHandler)
	log.Println("Listing for requests at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
