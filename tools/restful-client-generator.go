package main

// Takes a Swagger JSON file and generates a Service Http client in Go

import (
	"encoding/json"
	"flag"
	"github.com/emicklei/go-restful/swagger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// go run restful-client-generator.go -url http://localhost:8080/apidocs.json

var apidocsUrl string

func main() {
	flag.StringVar(&apidocsUrl, "url", "", "endpoint of a REST service (e.g. http://myservice/apidocs.json)")
	flag.Parse()
	if len(apidocsUrl) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.Print("fetching Swagger docs")

	listing := new(swagger.ResourceListing)
	fetch(apidocsUrl, &listing)
	log.Printf("Api version:%s, Swagger version:%s", listing.ApiVersion, listing.SwaggerVersion)

	for _, each := range listing.Apis {
		generateForApi(each)
	}
}

func generateForApi(api swagger.Api) {
	declaration := new(swagger.ApiDeclaration)
	fetch(apidocsUrl+api.Path, &declaration)
	for _, each := range declaration.Apis {
		log.Printf("api:%v", each.Path)
		for _, op := range each.Operations {
			log.Printf("operation:%v", op.HttpMethod)
		}
	}
}

func fetch(url string, model interface{}) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Unable to fetch Swagger specification:%v", err)
	}
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read Swagger specification:%v", err)
	}
	json.Unmarshal(buffer, model)
}
