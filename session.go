package requests

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

// var
var (
	ErrEmptyProxy = errors.New("proxy is empty")
)

// Session httpclient session
// Clients and Transports are safe for concurrent use by multiple goroutines
// for efficiency should only be created once and re-used.
// so, session is also safe for concurrent use by multiple goroutines.
type Session struct {
	*http.Transport
	*http.Client
	opts Options
	// optFunc []Option
	LogFunc func(string, ...interface{})
	errs    chan error
	wg      *sync.Mutex
}

// New new session
func New(opts ...Option) *Session {

	options := newOptions(opts...)

	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 限制建立TCP连接的时间
			KeepAlive: 300 * time.Second,
		}).DialContext,
		MaxIdleConns: 100,
		// TLSHandshakeTimeout:   10 * time.Second, // 限制 TLS握手的时间
		// IdleConnTimeout:       120 * time.Second,
		// ResponseHeaderTimeout: 60 * time.Second, // 限制读取response header的时间
		DisableCompression: true,
		DisableKeepAlives:  false,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !options.Verify,
		},
	}
	jar, _ := cookiejar.New(nil)

	sess := &Session{
		Transport: tr,
		Client: &http.Client{
			Timeout:   time.Duration(options.Timeout) * time.Millisecond,
			Transport: tr,
			Jar:       jar,
		},
		opts: options,
		// optFunc: opts,
		LogFunc: func(format string, v ...interface{}) {
			fmt.Fprintf(os.Stderr, format+"\n", v...)
		},
		errs: make(chan error),
		wg:   new(sync.Mutex),
	}
	return sess
}

// Init init
func (sess *Session) Init(opts ...Option) {
	for _, o := range opts {
		o(&sess.opts)
	}
}

// Load config
func (sess *Session) Load(v interface{}) error {
	return sess.opts.Load(v)
}

// Proxy set proxy addr
// os.Setenv("HTTP_PROXY", "http://127.0.0.1:9743")
// os.Setenv("HTTPS_PROXY", "https://127.0.0.1:9743")
func (sess *Session) Proxy(addr string) error {
	if addr == "" {
		return ErrEmptyProxy
	}
	proxyURL, err := url.Parse(addr)
	if err != nil {
		return err
	}
	switch proxyURL.Scheme {
	case "socks5", "socks4":
		sess.Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, nil, proxy.Direct)
			if err != nil {
				return nil, err
			}
			return dialer.Dial(network, addr)
		}
	default:
		sess.Transport.Proxy = http.ProxyURL(proxyURL)
	}
	return nil
}

// SetLogFunc set log handler
func (sess *Session) SetLogFunc(f func(string, ...interface{})) *Session {
	sess.LogFunc = f
	return sess
}

// Timeout set client timeout
func (sess *Session) Timeout(timeout int) *Session {
	sess.Client.Timeout = time.Duration(timeout) * time.Second
	return sess
}

// SetKeepAlives set transport disableKeepAlives default transport is keepalive,
// if set false, only use the connection to the server for a single HTTP request.
func (sess *Session) SetKeepAlives(keepAlives bool) *Session {
	sess.Transport.DisableKeepAlives = !keepAlives
	return sess
}

// DoRequest send a request and return a response
func (sess Session) DoRequest(ctx context.Context, opts ...Option) (*Response, error) {

	sess.wg.Lock()

	options := sess.opts
	for _, o := range opts {
		o(&options)
	}

	req, err := Request(options)

	sess.wg.Unlock()

	start := time.Now()
	var resp *http.Response

	if err != nil {
		return WarpResponse(start, req, resp, err), fmt.Errorf("Request: %w", err)
	}

	if ctx != nil {
		req = req.WithContext(ctx) // !!! WithContext returns a shallow copy of r with its context changed to ctx
	}
	if options.Trace {
		resp, err = sess.DebugTrace(req)
	} else {
		resp, err = sess.Client.Do(req)
	}

	return WarpResponse(start, req, resp, err), err
}

// Do http request
func (sess *Session) Do(method, url, contentType string, body io.Reader) (*Response, error) {
	return sess.DoRequest(context.Background(), Method(method), URL(url), Header("Content-Type", contentType), Reader(body))
}

// DoWithContext http request
func (sess *Session) DoWithContext(ctx context.Context, method, url, contentType string, body io.Reader) (*Response, error) {
	return sess.DoRequest(ctx, Method(method), URL(url), Header("Content-Type", contentType), Reader(body))
}

// Get send get request
func (sess *Session) Get(url string) (*Response, error) {
	return sess.DoRequest(context.Background(), Method("GET"), URL(url))
}

// Head send head request
func (sess *Session) Head(url string) (*Response, error) {
	return sess.DoRequest(context.Background(), Method("HEAD"), URL(url))
}

// GetWithContext http request
func (sess *Session) GetWithContext(ctx context.Context, url string) (*Response, error) {
	return sess.DoRequest(ctx, Method("GET"), URL(url))
}

// Post send post request
func (sess *Session) Post(url, contentType string, body io.Reader) (*Response, error) {
	return sess.Do("POST", url, contentType, body)
}

// PostWithContext send post request
func (sess *Session) PostWithContext(ctx context.Context, url, contentType string, body io.Reader) (*Response, error) {
	return sess.DoWithContext(ctx, "POST", url, contentType, body)
}

// PostForm post form request
func (sess *Session) PostForm(url string, data url.Values) (*Response, error) {
	return sess.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// PostFormWithContext post form request
func (sess *Session) PostFormWithContext(ctx context.Context, url string, data url.Values) (*Response, error) {
	return sess.PostWithContext(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// Put send put request
func (sess *Session) Put(url, contentType string, body io.Reader) (*Response, error) {
	return sess.Do("PUT", url, contentType, body)
}

// PutWithContext send put request
func (sess *Session) PutWithContext(ctx context.Context, url, contentType string, body io.Reader) (*Response, error) {
	return sess.DoWithContext(ctx, "PUT", url, contentType, body)
}

// Delete send delete request
func (sess *Session) Delete(url, contentType string, body io.Reader) (resp *Response, err error) {
	return sess.Do("DELETE", url, contentType, body)
}

// DeleteWithContext send delete request
func (sess *Session) DeleteWithContext(ctx context.Context, url, contentType string, body io.Reader) (*Response, error) {
	return sess.DoWithContext(ctx, "DELETE", url, contentType, body)
}

// Upload upload file
func (sess *Session) Upload(url, file string) (*Response, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return sess.Post(url, "binary/octet-stream", f)
}

// Uploadmultipart upload with multipart form
func (sess *Session) Uploadmultipart(url, file string, fields map[string]string) (*Response, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", fields["filename"])
	if err != nil {
		return nil, fmt.Errorf("CreateFormFile %v", err)
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, fmt.Errorf("copying fileWriter %v", err)
	}
	for k, v := range fields {
		if err = writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}

	err = writer.Close() // close writer before POST request
	if err != nil {
		return nil, fmt.Errorf("writerClose: %v", err)
	}

	return sess.Post(url, writer.FormDataContentType(), body)
}

// DebugTrace trace a request
// ** BUG ** 显示URI中的‘/’参数会被%转义, 字符串%有特殊含义. 正确输入为:namespace=aws/ec2 实际输出为: namespace=aws%!F(MISSING)ec2
// 不影响使用, 展示问题
func (sess *Session) DebugTrace(req *http.Request) (*http.Response, error) {
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			sess.LogFunc("* Connect: %v", hostPort)
		},
		ConnectStart: func(network, addr string) {
			sess.LogFunc("* Trying %v %v...", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			sess.LogFunc("* Completed connection: %v %v, err=%v", network, addr, err)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			sess.LogFunc("* Got Conn: %v -> %v", connInfo.Conn.LocalAddr(), connInfo.Conn.RemoteAddr())
		},
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			sess.LogFunc("* Resolved Host: %v", dnsInfo.Host)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			var ipaddrs []string
			for _, ipaddr := range dnsInfo.Addrs {
				ipaddrs = append(ipaddrs, ipaddr.String())
			}
			sess.LogFunc("* Resolved DNS: %v, Coalesced: %v, err=%v", ipaddrs, dnsInfo.Coalesced, dnsInfo.Err)
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			sess.LogFunc("* SSL HandshakeComplete: %v", state.HandshakeComplete)
		},
		WroteRequest: func(reqInfo httptrace.WroteRequestInfo) {
		},
	}

	ctx := httptrace.WithClientTrace(req.Context(), trace)
	req2 := req.WithContext(ctx)
	reqLog, err := DumpRequest(req2)
	if err != nil {
		sess.LogFunc("request error: %w", err)
		return nil, err
	}
	resp, err := sess.Transport.RoundTrip(req2)
	sess.LogFunc(show(reqLog, "> "))
	if err != nil {
		sess.LogFunc("response error: %w", err)
		return nil, err
	}

	respLog, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	sess.LogFunc(show(respLog, "< "))
	return resp, nil
}

func show(b []byte, prompt string) string {
	var maxTruncate = 9999
	var buf bytes.Buffer
	for _, line := range bytes.Split(b, []byte("\n")) {
		buf.Write([]byte(prompt))
		buf.Write(bytes.Replace(line, []byte("%"), []byte("%%"), -1))
		buf.WriteString("\n")
	}
	str := buf.String()
	if len(str) > maxTruncate {
		return fmt.Sprintf("%s...[Len=%d, Truncated]", str[:maxTruncate], len(str))
	}
	return str
}
