package requests

import (
	"bytes"
	"net/http"
	"net/http/httputil"
)

// DumpRequest returns the given request in its HTTP/1.x wire representation.
func DumpRequest(req *http.Request) ([]byte, error) {
	return httputil.DumpRequestOut(req, true)
}

// DumpRequestIndent warp Dump
func DumpRequestIndent(req *http.Request) string {
	dump, _ := DumpRequest(req)
	var b bytes.Buffer
	for _, line := range bytes.Split(dump, []byte("\n")) {
		b.Write([]byte("> "))
		b.Write(line)
		b.WriteString("\n")
	}
	return b.String()
}
