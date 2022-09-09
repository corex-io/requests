package requests

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Options request
type Options struct {
	Method    string         `json:"method"`
	URL       string         `json:"url"`
	Path      []string       `json:"path"`
	Params    map[string]any `json:"params"`
	body      any
	Header    http.Header   `json:"headers"`
	Cookies   []http.Cookie `json:"cookies"`
	Timeout   time.Duration `json:"timeout"`
	Trace     bool          `json:"trace"`
	Verify    bool          `json:"verify"`
	Logf      func(ctx context.Context, stat Stat)
	LocalAddr net.Addr
}

// Option func
type Option func(*Options)

// NewOptions new request
func newOptions(opts ...Option) Options {
	opt := Options{
		Method:  "GET",
		Params:  make(map[string]any),
		Header:  make(http.Header),
		Timeout: 30 * time.Second,
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
func Params(query map[string]any) Option {
	return func(o *Options) {
		for k, v := range query {
			o.Params[k] = v
		}
	}
}

// Param params
func Param(k string, v any) Option {
	return func(o *Options) {
		o.Params[k] = v
	}
}

// Body request body
func Body(body any) Option {
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
func Timeout(timeout time.Duration) Option {
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

func Logf(f func(context.Context, Stat)) Option {
	return func(o *Options) {
		o.Logf = f
	}
}

func LocalAddr(addr net.Addr) Option {
	return func(o *Options) {
		o.LocalAddr = addr
	}
}

// Copy copy
func (opt Options) Copy() Options {
	options := Options{
		Method:    opt.Method,
		URL:       opt.URL,
		Path:      opt.Path,
		Params:    opt.Params,
		Header:    opt.Header,
		Cookies:   opt.Cookies,
		body:      opt.body,
		Timeout:   opt.Timeout,
		Trace:     opt.Trace,
		Verify:    opt.Verify,
		Logf:      opt.Logf,
		LocalAddr: opt.LocalAddr,
	}
	return options
}
