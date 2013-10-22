package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type UserServiceClient struct {
	httpClient *http.Client
	scheme     string
	host       string
	port       int
}

func NewUserServiceClient(client *http.Client, scheme string, host string, port int) *UserServiceClient {
	return &UserServiceClient{client, scheme, host, port}
}

func (c UserServiceClient) newURIBuilder(template string) *URIBuilder {
	return NewURIBuilder(c.scheme, c.host, c.port, template)
}

type URIBuilder struct {
	scheme           string
	host             string
	port             int
	template         string
	pathParameters   map[string]string
	queryParameters  map[string][]string
	headerParameters map[string]string
}

func NewURIBuilder(scheme string, host string, port int, template string) *URIBuilder {
	return &URIBuilder{
		scheme: scheme, host: host, port: port,
		template:         template,
		pathParameters:   map[string]string{},
		queryParameters:  map[string][]string{},
		headerParameters: map[string]string{},
	}
}

func (u *URIBuilder) PathParam(name string, value string) {
	u.pathParameters[name] = value
}

func (u URIBuilder) Build() string {
	return ""
}

type User struct {
	Id   string
	Name string
}

// findUser : get a user
// GET /users/{user-id}
func (c UserServiceClient) findUser(user_id string) (*User, error) {
	uri := c.newURIBuilder("/users/{user-id}")
	uri.PathParam("user-id", user_id)
	req, err := http.NewRequest("GET", uri.Build(), nil)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err = json.Unmarshal(buffer, model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func main() {}
