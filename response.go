package requests

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"

	"encoding/json"
)

// Response wrap std response
type Response struct {
	*http.Request
	*http.Response
	once *sync.Once
	body *bytes.Buffer
	err  error
}

// newResponse newResponse
func newResponse() *Response {
	return &Response{
		once: new(sync.Once),
		body: new(bytes.Buffer),
	}
}

func (resp *Response) getBody() {
	resp.once.Do(func() {
		if resp.Response == nil || resp.Response.Body == nil {
			return
		}
		defer resp.Response.Body.Close()
		var n int64
		n, resp.err = resp.body.ReadFrom(resp.Response.Body)
		if resp.Response.ContentLength <= 0 {
			resp.Response.ContentLength = n
		}
	})
}

// WarpResponse warp response
func WarpResponse(resp *http.Response, req ...*http.Request) *Response {
	resp2 := newResponse()
	resp2.Response = resp
	if len(req) != 0 {
		resp2.Request = req[0]
	}
	resp2.getBody()
	return resp2
}

func (resp Response) String() string {
	return resp.Text()
}

func (resp Response) Error() string {
	if resp.err == nil {
		return ""
	}
	return resp.err.Error()
}

// StdLib return net/http.Response
func (resp *Response) StdLib() *http.Response {
	return resp.Response
}

// Text parse parse to string
func (resp *Response) Text() string {
	return resp.body.String()
}

// Download parse response to a file
func (resp *Response) Download(name string) (int, error) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(resp.body.Bytes())
}

// JSON parse response
func (resp *Response) JSON(v interface{}) error {
	return json.Unmarshal(resp.body.Bytes(), v)
}

// Dump returns the given request in its HTTP/1.x wire representation.
func (resp *Response) Dump() ([]byte, error) {
	return httputil.DumpResponse(resp.Response, true)
}

// Copy deep copy response
func (resp *Response) Copy() *Response {
	resp2 := &Response{}
	return resp2
}
