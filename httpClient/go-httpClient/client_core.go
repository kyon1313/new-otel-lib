package apw_http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	defaultMaxConnection             = 5
	defaultResponseTimeoutDuration   = time.Second * 5
	defaultConnectionTimeoutDuration = time.Second * 1
)

func (c *httpClient) do(httpMethod, url string, headers http.Header, body any) (*Response, error) {

	fullHeaders := c.getRequestHeader(headers)

	requestBody, err := c.getRequestBody(fullHeaders.Get("Content-Type"), body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.New("unable to create new request")
	}

	request.Header = fullHeaders

	client := c.getHttpClient()

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	finalResponse := &Response{
		status:     resp.Status,
		statusCode: resp.StatusCode,
		headers:    resp.Header,
		body:       respBody,
	}

	return finalResponse, nil
}

func (c *httpClient) getHttpClient() *http.Client {
	c.clientOnce.Do(func() {
		c.client = &http.Client{
			Timeout: c.getConnectionTimeout() + c.getResponseTimeout(),
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   c.getMaxIdleConnections(),
				ResponseHeaderTimeout: c.getResponseTimeout(),
				DialContext: (&net.Dialer{
					Timeout: c.getConnectionTimeout(),
				}).DialContext,
			},
		}
	})

	return c.client
}

func (c *httpClient) getMaxIdleConnections() int {
	if c.builder.maxIdleConnection > 0 {
		return c.builder.maxIdleConnection
	}

	return defaultMaxConnection
}

func (c *httpClient) getResponseTimeout() time.Duration {
	if c.builder.responseTimeout > 0 {
		return c.builder.responseTimeout
	}

	if c.builder.disableTimeout {
		return 0
	}

	return defaultResponseTimeoutDuration
}

func (c *httpClient) getConnectionTimeout() time.Duration {
	if c.builder.connectionTimeout > 0 {
		return c.builder.connectionTimeout
	}

	if c.builder.disableTimeout {
		return 0
	}

	return defaultConnectionTimeoutDuration
}

func (c *httpClient) getRequestHeader(requestHeaders http.Header) http.Header {
	result := make(http.Header)

	for header, value := range c.builder.headers {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}

	for header, value := range requestHeaders {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}
	return result
}

func (v *httpClient) getRequestBody(contentType string, body any) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	switch strings.ToLower(contentType) {
	case "application/json":
		return json.Marshal(body)
	case "application/xml":
		return xml.Marshal(body)
	default:
		return json.Marshal(body)
	}

}
