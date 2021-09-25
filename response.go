package requests

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	"encoding/json"
)

// Stat stats
type Stat struct {
	StartAt string `json:"StartAt"`
	Cost    int64  `json:"Cost"`
	Request struct {
		Method string            `json:"Method"`
		Header map[string]string `json:"Header"`
		URL    string            `json:"URL"`
		Body   interface{}       `json:"Body"`
	} `json:"Request"`
	Response struct {
		Header        map[string]string `json:"Header"`
		Body          interface{}       `json:"Body"`
		StatusCode    int               `json:"StatusCode"`
		ContentLength int64             `json:"ContentLength"`
	} `json:"Response"`
	Err string `json:"Err"`
}

func (stat Stat) String() string {
	b, _ := json.Marshal(stat)
	return string(b)
}

// Response wrap std response
type Response struct {
	StartAt time.Time
	Cost    time.Duration
	*http.Request
	*http.Response
	once sync.Once
	body *bytes.Buffer
	err  error
}

// newResponse newResponse
func newResponse() *Response {
	return &Response{
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
func WarpResponse(start time.Time, req *http.Request, resp *http.Response, err error) *Response {
	resp2 := newResponse()
	resp2.StartAt = start
	resp2.Cost = time.Since(start)
	resp2.Request = req
	resp2.err = err
	resp2.Response = resp
	resp2.getBody()
	return resp2
}

// Stat stat
func (resp Response) Stat() Stat {

	stat := Stat{
		StartAt: resp.StartAt.Format("2006-01-02 15:04:05.000"),
		Cost:    resp.Cost.Milliseconds(),
	}

	if resp.Response != nil {
		m2 := make(map[string]interface{})

		if err := json.Unmarshal(resp.body.Bytes(), &m2); err != nil {
			stat.Response.Body = resp.body.String()
		} else {
			stat.Response.Body = m2
		}

		stat.Response.Header = make(map[string]string)

		for k, v := range resp.Response.Header {
			stat.Response.Header[k] = v[0]
		}
		stat.Response.ContentLength = resp.Response.ContentLength
		stat.Response.StatusCode = resp.StatusCode

		if resp.Request == nil {
			resp.Request = resp.Response.Request
		}

	}

	if resp.Request != nil {
		stat.Request.Method = resp.Request.Method
		stat.Request.URL = resp.Request.URL.String()
		if resp.Request.GetBody != nil {
			body, _ := resp.Request.GetBody()

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(body)

			m := make(map[string]interface{})

			if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
				stat.Request.Body = buf.String()
			} else {
				stat.Request.Body = m
			}
		}

		stat.Request.Header = make(map[string]string)

		for k, v := range resp.Request.Header {
			stat.Request.Header[k] = v[0]
		}
	}

	if resp.err != nil {
		stat.Err = resp.err.Error()
	}

	return stat
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

// Bytes bytes
func (resp *Response) Bytes() []byte {
	return resp.body.Bytes()
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
