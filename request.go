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
func Request(opt Options) (*http.Request, error) {

	var reader io.Reader
	switch {
	case opt.body != nil:
		switch v := opt.body.(type) {
		case []byte:
			reader = bytes.NewReader(v)
		case string:
			reader = bytes.NewReader([]byte(v))
		default:
			b, err := json.Marshal(opt.body)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewReader(b)
		}
	case opt.reader != nil:
		reader = opt.reader
	case len(opt.Form) != 0:
		reader = strings.NewReader(opt.Form.Encode())
	}

	req, err := http.NewRequest(opt.Method, opt.URL, reader)
	if err != nil {
		return nil, err
	}
	req.URL.Path = path.Join(req.URL.Path, path.Join(opt.Path...))
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
