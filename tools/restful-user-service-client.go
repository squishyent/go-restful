package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Service struct {
	httpClient *http.Client
	scheme     string
	host       string
	port       int
}

func New(client *http.Client, scheme string, host string, port int) *Service {
	return &Service{client, scheme, host, port}
}

type User struct {
	Id   string
	Name string
}

// findUser : get a user
// GET /users/{user-id}
func (c Service) findUser(user_id string) (*User, error) {
	var body io.Reader = nil
	uri := NewURIBuilder(c.scheme, c.host, c.port, "/users/{user-id}")
	uri.PathParam("user-id", user_id)
	req, err := http.NewRequest("GET", uri.Build(), body)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	model := new(User)
	if err := json.NewDecoder(resp.Body).Decode(model); err != nil {
		return nil, err
	}
	return model, nil
}

func main() {
	service := New(new(http.Client), "http", "localhost", 8080)
	usr, err := service.findUser("1")
	fmt.Printf("%v,%v\n", usr, err)
}

// URIBuilder is a helper object to construct URIs from template and parameters.
type URIBuilder struct {
	scheme          string
	host            string
	port            int
	template        string
	pathParameters  map[string]string
	queryParameters map[string][]string
}

// NewURIBuilder create a new URIBuilder from the given host,port and template.
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

// PathParam add or replaces the value of a Path parameter
func (u *URIBuilder) PathParam(name string, value string) {
	u.pathParameters[name] = value
}

// QueryParam adds the value of a Query parameter; creates a list for multiple values.
func (u *URIBuilder) QueryParam(name string, value string) {
	list := u.queryParameters[name]
	if len(list) == 0 {
		u.queryParameters[name] = []string{value}
	} else {
		u.queryParameters[name] = append(list, value)
	}
}

// Build returns the URI based on the scheme,host,port,template and parameters.
func (u URIBuilder) Build() string {
	var buf = new(bytes.Buffer)
	buf.WriteString(u.scheme)
	buf.WriteString("://")
	buf.WriteString(u.host)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(u.port)) //fmt.Fprint(buf, "%d", u.port)
	tokens := strings.Split(u.template, "/")
	for _, each := range tokens {
		if len(each) > 0 {
			buf.WriteByte('/')
			if strings.HasPrefix(each, "{") { // substitute
				buf.WriteString(u.pathParameters[each[1:len(each)-1]])
			} else {
				buf.WriteString(each)
			}
		}
	}
	if len(u.queryParameters) > 0 {
		buf.WriteByte('?')
		one := false
		for key, value := range u.queryParameters {
			if one {
				buf.WriteByte('&')
			} else {
				one = true
			}
			for i, elem := range value {
				if i > 0 {
					buf.WriteByte('&')
				}
				buf.WriteString(url.QueryEscape(key))
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(elem))
			}
		}
	}
	return buf.String()
}
