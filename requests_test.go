package requests

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// func Test_Download(t *testing.T) {
// 	t.Log("Testing Download")
// 	err := DownloadFile("https://github.com/prometheus/prometheus/releases/download/v2.12.0/prometheus-2.12.0.linux-amd64.tar.gz", true)
// 	t.Log(err)
// }

func Test_Basic(t *testing.T) {
	resp, _ := Get("http://httpbin.org/get")
	t.Log(resp.Text())
	resp, _ = Post("http://httpbin.org/post", "application/json", strings.NewReader(`{"a": "b"}`))
	t.Log(resp.Text())
}

func Test_Get(t *testing.T) {
	t.Log("Testing get request")
	sess := New(
		Retry(3),
		Header("a", "b"),
		Cookie("username", "golang"),
		BasicAuth("user", "123456"),
		Timeout(1),
	)

	// req.SetParam("uid", 1).SetCookie("username", "000000")
	resp, err := sess.DoRequest(context.Background(), Method("GET"), URL("http://httpbin.org/get"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.StatusCode, err)
	t.Log(resp.Text())
}

func Test_Post(t *testing.T) {
	t.Log("Testing get request")
	sess := New(BasicAuth("user", "123456"))
	resp, err := sess.DoRequest(context.Background(),
		Method("POST"),
		URL("http://httpbin.org/post"),
		Params(map[string]interface{}{
			"a": "b",
			"c": 3,
			"d": []int{1, 2, 3},
		}),
		Body(`{"body":"QWER"}`),
		Retry(3),
		Header("hello", "world"),
		Trace(true),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.StatusCode, err, resp.Response.ContentLength, resp.Request.ContentLength)
	t.Log(resp.Text())
	t.Log(resp.Stat())
}

func Test_FormPost(t *testing.T) {
	t.Log("Testing get request")

	sess := New()
	resp, err := sess.DoRequest(context.Background(),
		Method("POST"),
		URL("http://httpbin.org/post"),
		Retry(3),
		Form("name", "12.com"),
		Params(map[string]interface{}{
			"a": "b",
			"c": 3,
			"d": []int{1, 2, 3},
		}),
		Trace(true),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	t.Log(resp.StatusCode, err, resp.Response.ContentLength, resp.Request.ContentLength)

}

func Test_PostForm2(t *testing.T) {
	res, err := http.PostForm("http://httpbin.org/post", url.Values{
		"key":   {"this is url key"},
		"value": {"this is url value"},
	})
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	resp := WarpResponse(time.Now(), nil, res, err)
	t.Log("$$$$$$4", resp.Stat())
}

func Test_Race(t *testing.T) {
	opts := Options{}
	ctx := context.Background()
	t.Logf("%#v", opts)
	sess := New(URL("http://httpbin.org/post")) //, Auth("user", "123456"))
	for i := 0; i < 10; i++ {
		go func() {
			sess.DoRequest(ctx, MethodPost, Body(`{"a":"b"}`), Params(map[string]interface{}{"1": "2"})) // nolint: errcheck
		}()
	}
	time.Sleep(3 * time.Second)
}
