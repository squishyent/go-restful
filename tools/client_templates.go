package main

import (
	"fmt"
)

func import_source() string {
	return `
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
`
}

func service_source() string {
	return `
type Service struct {
	httpClient *http.Client
	scheme     string
	host       string
	port       int
}

func New(client *http.Client, scheme string, host string, port int) *Service {
	return &Service{client, scheme, host, port}
}
`
}

func dorequest_source() string {
	return `resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}`
}

func decoderesponse_source(typeName string) string {
	return fmt.Sprintf("model := new(%s)
	if err := json.NewDecoder(resp.Body).Decode(model); err != nil {
		return nil, err
	}
	return model, nil
	", typeName)
}

func createrequest_source(method string) string {
	return fmt.Sprintf("req, err := http.NewRequest(\"%s\", uri.Build(), body)", method)
}
