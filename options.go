package requests

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
)

// Options request
type Options struct {
	Method  string
	URL     string
	Path    []string
	Params  map[string]interface{}
	Headers map[string]string
	Cookies map[string]string
	body    interface{}
	reader  io.Reader
	Form    url.Values
	Timeout int
	Retry   int
	Trace   bool
}

// Option func
type Option func(*Options)

// NewOptions new request
func newOptions(opts ...Option) Options {
	opt := Options{
		Method:  "GET",
		Params:  make(map[string]interface{}),
		Headers: make(map[string]string),
		Cookies: make(map[string]string),
		Form:    make(url.Values),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

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

// Reader set reader body
func Reader(reader io.Reader) Option {
	return func(o *Options) {
		o.reader = reader
	}
}

// Form set form, content-type is
func Form(k, v string) Option {
	return func(o *Options) {
		o.Headers["content-type"] = "application/x-www-form-urlencoded"
		o.Form.Add(k, v)
	}
}

// Header header
func Header(k, v string) Option {
	return func(o *Options) {
		o.Headers[k] = v
	}
}

// Headers headers
func Headers(kv map[string]string) Option {
	return func(o *Options) {
		for k, v := range kv {
			o.Headers[k] = v
		}
	}
}

// Cookie cookie
func Cookie(k, v string) Option {
	return func(o *Options) {
		o.Cookies[k] = v
	}
}

// Cookies cookie
func Cookies(kv map[string]string) Option {
	return func(o *Options) {
		for k, v := range kv {
			o.Cookies[k] = v
		}
	}

}

// Auth base auth
func Auth(user, pass string) Option {
	return Header("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))

}

// Retry set retry
func Retry(retry int) Option {
	return func(o *Options) {
		o.Retry = retry
	}
}

// Timeout timeout
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

// MergeIn merge r into req
func (opt *Options) MergeIn(o Options) {
	for k, v := range o.Params {
		if _, ok := opt.Params[k]; !ok {
			opt.Params[k] = v
		}
	}
	for k, v := range o.Headers {
		if _, ok := opt.Headers[k]; !ok {
			opt.Headers[k] = v
		}
	}
	for k, v := range o.Cookies {
		if _, ok := opt.Cookies[k]; !ok {
			opt.Cookies[k] = v
		}
	}
	if opt.Retry == 0 {
		opt.Retry = o.Retry
	}
	if opt.Timeout == 0 {
		opt.Timeout = o.Timeout
	}
	if opt.Trace == false {
		opt.Trace = o.Trace
	}
}

// Request request
func (opt *Options) Request() (*http.Request, error) {
	return Request(*opt)
}
