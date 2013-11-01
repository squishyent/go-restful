package main

// Takes a Swagger JSON file and generates a Service Http client in Go

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/emicklei/go-restful/swagger"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// go run *.go -url http://localhost:8080/apidocs.json -pkg users && cat /tmp/service.go

var apidocsUrl string
var packageName string

func main() {
	flag.StringVar(&apidocsUrl, "url", "", "endpoint of a REST service (e.g. http://myservice/apidocs.json)")
	flag.StringVar(&packageName, "pkg", "", "name of the package for the generated Service.")
	flag.Parse()
	if len(apidocsUrl) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
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
	out, _ := os.Create("/tmp/service.go")
	defer out.Close()
	io.WriteString(out, package_source(packageName))
	io.WriteString(out, import_service_new_source())

	for _, each := range declaration.Apis {
		log.Printf("api:%v\n", each.Path)
		for _, op := range each.Operations {
			generateForOperation(each.Path, op, out)
		}
	}
	io.WriteString(out, uribuilder_source())
}

func generateForOperation(path string, op swagger.Operation, out io.Writer) {
	io.WriteString(out, "func (c Service) ")
	io.WriteString(out, op.Nickname)
	io.WriteString(out, "(")
	for i, each := range op.Parameters {
		if i > 0 {
			io.WriteString(out, ", ")
		}
		writeParameterSignature(each, out)
	}
	io.WriteString(out, ") ")
	if "void" == op.Type {
		io.WriteString(out, "(interface{}, error)")
	} else {
		io.WriteString(out, "("+noPkg(op.Type)+", error)")
	}
	io.WriteString(out, " {\n")
	io.WriteString(out, newuri_source(path))
	for _, each := range op.Parameters {
		writeSetUriParameter(each, out)
	}
	io.WriteString(out, createrequest_source(op.HttpMethod))
	io.WriteString(out, accept_source(op.Consumes[0])) // TODO list all
	io.WriteString(out, dorequest_source())
	io.WriteString(out, "\n}\n")
}

func writeSetUriParameter(param swagger.Parameter, out io.Writer) {
	if "path" == param.ParamType {
		io.WriteString(out, fmt.Sprintf("\turi.pathParam(\"%s\",%s)\n", param.Name, toVar(param.Name)))
	} else if "query" == param.ParamType {
		io.WriteString(out, fmt.Sprintf("\turi.queryParam(\"%s\",%s)\n", param.Name, toVar(param.Name)))
	}
}

func writeParameterSignature(param swagger.Parameter, out io.Writer) {
	fmt.Fprintf(out, "%s %s", toVar(noPkg(param.Name)), noPkg(param.Type))
}

func toVar(varName string) string {
	return strings.ToLower(strings.NewReplacer("-", "_").Replace(varName))
}

func noPkg(name string) string {
	dot := strings.LastIndex(name, ".")
	if dot == -1 {
		return name
	} else {
		return name[dot+1:]
	}
}

func fetch(url string, model interface{}) {
	log.Printf("fetching %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Unable to fetch Swagger specification:%v", err)
	}
	defer resp.Body.Close()
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read Swagger specification:%v", err)
	}
	json.Unmarshal(buffer, model)
}
