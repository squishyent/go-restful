package main

import (
	"fmt"
)

func package_source(pkg string) string {
	return fmt.Sprintf("package %s\n\n", pkg)
}

func import_service_new_source() string {
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

func newuri_source(path string) string {
	return fmt.Sprintf(`	var body io.Reader = nil
	uri := NewURIBuilder(c.scheme, c.host, c.port, "%s")
`, path)
}

func dorequest_source() string {
	return `	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}`
}

func decoderesponse_source(typeName string) string {
	return fmt.Sprintf(`model := new(%s)
	if err := json.NewDecoder(resp.Body).Decode(model); err != nil {
		return nil, err
	}
	return model, nil
	`, typeName)
}

func createrequest_source(method string) string {
	return fmt.Sprintf("	req, err := http.NewRequest(\"%s\", uri.Build(), body)\n", method)
}

func contenttype_source(mime string) string {
	return fmt.Sprintf("	req.Header.Set(\"Content-Type\", \"%s\")\n", mime)
}

func accept_source(mime string) string {
	return fmt.Sprintf("	req.Header.Set(\"Accept\", \"%s\")\n", mime)
}
