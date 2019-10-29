package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Request request
func Request(opt Options, opts ...Option) (*http.Request, error) {
	for _, o := range opts {
		o(&opt)
	}

	var reader io.Reader
	switch {
	case opt.reader != nil:
		reader = opt.reader
	case opt.body != nil:
		b, err := json.Marshal(opt.body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	case opt.Form != nil:
		reader = strings.NewReader(opt.Form.Encode())
	}

	req, err := http.NewRequest(opt.Method, opt.URL+path.Join(opt.Path...), reader)
	if err != nil {
		return nil, err
	}

	for k, v := range opt.Params {
		if req.URL.RawQuery != "" {
			req.URL.RawQuery += "&"
		}
		req.URL.RawQuery += k + "=" + url.QueryEscape(fmt.Sprintf("%v", v))
	}
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}
	for k, v := range opt.Cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	return req, nil
}
