package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	switch v := body.(type) {
	case []byte:
		return bytes.NewReader(v), nil
	case string:
		return strings.NewReader(v), nil
	case *bytes.Buffer:
		return bytes.NewReader(v.Bytes()), nil
	case io.Reader, *bytes.Reader, io.ReadSeeker:
		return body.(io.Reader), nil
	case url.Values:
		return strings.NewReader(v.Encode()), nil
	default:
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(b), nil
	}
}

// NewRequestWithContext request
func NewRequestWithContext(ctx context.Context, opt Options) (*http.Request, error) {
	reader, err := getBodyReader(opt.body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, opt.Method, opt.URL, reader)
	if err != nil {
		return nil, err
	}

	// req.URL.Path = path.Join(req.URL.Path, path.Join(opt.Path...))
	for _, path := range opt.Path {
		req.URL.Path += path
	}

	for k, v := range opt.Params {
		if req.URL.RawQuery != "" {
			req.URL.RawQuery += "&"
		}
		req.URL.RawQuery += k + "=" + url.QueryEscape(fmt.Sprintf("%v", v))
	}

	req.Header = opt.Header

	for _, cookie := range opt.Cookies {
		req.AddCookie(&cookie)
	}
	return req, nil
}
