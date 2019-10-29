package requests_test

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"requests"
)

// func Test_Download(t *testing.T) {
// 	t.Log("Testing Download")
// 	err := requests.DownloadFile("https://github.com/prometheus/prometheus/releases/download/v2.12.0/prometheus-2.12.0.linux-amd64.tar.gz", true)
// 	t.Log(err)
// }

func Test_Basic(t *testing.T) {
	resp, _ := requests.Get("http://httpbin.org/get")
	t.Log(resp.Text())
	resp, _ = requests.Post("http://httpbin.org/post", "application/json", strings.NewReader(`{"a": "b"}`))
	t.Log(resp.Text())
}

func Test_Get(t *testing.T) {
	t.Log("Testing get request")
	sess := requests.New(
		requests.Retry(3),
		requests.Header("a", "b"),
		requests.Cookie("username", "golang"),
		requests.Auth("user", "123456"))

	// req.SetParam("uid", 1).SetCookie("username", "000000")
	resp, err := sess.DoRequest(nil, requests.Method("GET"), requests.URL("http://httpbin.org/get"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.StatusCode, err)
	t.Log(resp.Text())
}

func Test_Post(t *testing.T) {
	t.Log("Testing get request")
	sess := requests.New(requests.Auth("user", "123456"))
	resp, err := sess.DoRequest(nil,
		requests.Method("POST"),
		requests.URL("http://httpbin.org/post"),
		requests.Params(map[string]interface{}{
			"a": "b",
			"c": 3,
			"d": []int{1, 2, 3},
		}),
		requests.Body(`{"body":"QWER"}`),
		requests.Retry(3),
		requests.Header("hello", "world"),
		requests.Trace(true),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.StatusCode, err, resp.Response.ContentLength, resp.Request.ContentLength)
	t.Log(resp.Text())
}

func Test_FormPost(t *testing.T) {
	t.Log("Testing get request")

	sess := requests.New()
	resp, err := sess.DoRequest(nil,
		requests.Method("POST"),
		requests.URL("http://httpbin.org/post"),
		requests.Retry(3),
		requests.Form("name", "12.com"),
		requests.Params(map[string]interface{}{
			"a": "b",
			"c": 3,
			"d": []int{1, 2, 3},
		}),
		requests.Trace(true),
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
	if err != nil {
		t.Error(err)
		return
	}
	resp := requests.WarpResponse(res)
	t.Log(resp.Text())
}
