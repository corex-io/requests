package requests

import (
	"bytes"
	"fmt"
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

const maxTruncateBytes = 1024

func show(b []byte, prompt string) string {
	var buf bytes.Buffer
	for _, line := range bytes.Split(b, []byte("\n")) {
		buf.Write([]byte(prompt))
		buf.Write(bytes.Replace(line, []byte("%"), []byte("%%"), -1))
		buf.WriteString("\n")
	}
	str := buf.String()
	if len(str) > maxTruncateBytes {
		return fmt.Sprintf("%s...[Len=%d, Truncated]", str[:maxTruncateBytes], len(str))
	}
	return str
}
