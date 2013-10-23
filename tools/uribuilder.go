package main

import (
	"bytes"
	"strconv"
	"strings"
)

type URIBuilder struct {
	scheme          string
	host            string
	port            int
	template        string
	pathParameters  map[string]string
	queryParameters map[string][]string
}

func NewURIBuilder(scheme string, host string, port int, template string) *URIBuilder {
	return &URIBuilder{
		scheme:          scheme,
		host:            host,
		port:            port,
		template:        template,
		pathParameters:  map[string]string{},
		queryParameters: map[string][]string{},
	}
}

func (u *URIBuilder) PathParam(name string, value string) {
	u.pathParameters[name] = value
}

func (u *URIBuilder) QueryParam(name string, value string) {
	u.queryParameters[name] = value
}

func (u URIBuilder) Build() string {
	var buf = new(bytes.Buffer)
	buf.WriteString(u.scheme)
	buf.WriteString("://")
	buf.WriteString(u.host)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(u.port)) //fmt.Fprint(buf, "%d", u.port)
	// check / prefix
	// replace if needed
	tokens := strings.Split(u.template, "{}")

	return buf.String()
}

func main() {
	uri := NewURIBuilder("http", "localhost", 9999, "/hier/{daar}/maar")
	uri.PathParam("daar", "xyz")
	println(uri.Build())
}
