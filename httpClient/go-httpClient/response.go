package apw_http

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	status     string
	statusCode int
	headers    http.Header
	body       []byte
}

func (r *Response) Status() string {
	return r.status
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Headers() http.Header {
	return r.headers
}

func (r *Response) BodyString() string {
	return string(r.body)
}

func (r *Response) BodyBytes() []byte {
	return r.body
}

func (r *Response) UnmarshalJson(target any) error {
	return json.Unmarshal(r.BodyBytes(), target)
}
