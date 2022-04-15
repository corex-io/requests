package requests

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
)

// Options request
type Options struct {
	Method  string                 `json:"method"`
	URL     string                 `json:"url"`
	Path    []string               `json:"path"`
	Params  map[string]interface{} `json:"params"`
	Header  http.Header            `json:"headers"`
	Cookies []http.Cookie          `json:"cookies"`
	body    interface{}
	Timeout int  `json:"timeout"`
	Trace   bool `json:"trace"`
	Verify  bool `json:"verify"`
}

// Option func
type Option func(*Options)

// NewOptions new request
func newOptions(opts ...Option) Options {
	opt := Options{
		Method:  "GET",
		Params:  make(map[string]interface{}),
		Header:  make(http.Header),
		Timeout: 30000, // 30s
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Method http method
var (
	MethodGet  = Method("GET")
	MethodPost = Method("POST")
)

// Method set method
func Method(method string) Option {
	return func(o *Options) {
		o.Method = method
	}
}

// URL set url
func URL(url string) Option {
	return func(o *Options) {
		o.URL = url
	}
}

// Path set path
func Path(path string) Option {
	return func(o *Options) {
		o.Path = append(o.Path, path)
	}
}

// Params add query args
func Params(query map[string]interface{}) Option {
	return func(o *Options) {
		for k, v := range query {
			o.Params[k] = v
		}
	}
}

// Param params
func Param(k string, v interface{}) Option {
	return func(o *Options) {
		o.Params[k] = v
	}
}

// Body request body
func Body(body interface{}) Option {
	return func(o *Options) {
		o.body = body
	}
}

// Form set form, content-type is
func Form(form url.Values) Option {
	return func(o *Options) {
		o.Header.Add("content-type", "application/x-www-form-urlencoded")
		o.body = form
	}
}

// Header header
func Header(k, v string) Option {
	return func(o *Options) {
		o.Header.Add(k, v)
	}
}

// Headers headers
func Headers(kv map[string]string) Option {
	return func(o *Options) {
		for k, v := range kv {
			o.Header.Add(k, v)
		}
	}
}

// Cookie cookie
func Cookie(cookie http.Cookie) Option {
	return func(o *Options) {
		o.Cookies = append(o.Cookies, cookie)
	}
}

// Cookies cookies
func Cookies(cookies ...http.Cookie) Option {
	return func(o *Options) {
		o.Cookies = append(o.Cookies, cookies...)
	}
}

// BasicAuth base auth
func BasicAuth(user, pass string) Option {
	return Header("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))

}

// Timeout timeout, Millisecond 毫秒
func Timeout(timeout int) Option {
	return func(o *Options) {
		o.Timeout = timeout
	}
}

// Trace Trace
func Trace(trace bool) Option {
	return func(o *Options) {
		o.Trace = trace
	}
}

// Verify verify
func Verify(verify bool) Option {
	return func(o *Options) {
		o.Verify = verify
	}
}

// Copy copy
func (opt *Options) Copy() (Options, error) {
	var options Options
	b, err := json.Marshal(opt)
	if err != nil {
		return options, err
	}
	return options, json.Unmarshal(b, &options)
}
