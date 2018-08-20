package new

import (
	"fmt"
	"testing"
)

func BenchmarkSequetiallyGenRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenHTTPReq(data)
	}
}

func BenchmarkGenRequest(b *testing.B) {
	qlen := 200
	hlen := 40
	req := Req{
		Url:  "http://localhost:{{.port}}",
		Body: `{"data":"{{.body}}"}`,
	}
	data := map[string]interface{}{
		"port": 8080,
		"body": "hello world",
	}

	for i := 0; i < qlen; i++ {
		req.Query = append(req.Query, KV{fmt.Sprintf("k%d", i), fmt.Sprintf("{{.qv%d}}", i)})
		data[fmt.Sprintf("qv%d", i)] = fmt.Sprintf("value%d", i)
	}

	for i := 0; i < hlen; i++ {
		req.Headers = append(req.Headers, KV{fmt.Sprintf("k%d", i), fmt.Sprintf("{{.hv%d}}", i)})
		data[fmt.Sprintf("hv%d", i)] = fmt.Sprintf("value%d", i)
	}

	InitTmpls(&req)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenHTTPReq(data)
		}
	})
}

func TestGenHTTPReq(t *testing.T) {
	t.Log("foo")

	qlen := 10
	hlen := 2
	req := Req{
		Url:  "http://localhost:{{.port}}",
		Body: `{"data":"{{.body}}"}`,
	}
	data := map[string]interface{}{
		"port": 8080,
		"body": "hello world",
	}

	for i := 0; i < qlen; i++ {
		req.Query = append(req.Query, KV{fmt.Sprintf("k%d", i), fmt.Sprintf("{{.qv%d}}", i)})
		data[fmt.Sprintf("qv%d", i)] = fmt.Sprintf("value%d", i)
	}

	for i := 0; i < hlen; i++ {
		req.Headers = append(req.Headers, KV{fmt.Sprintf("k%d", i), fmt.Sprintf("{{.hv%d}}", i)})
		data[fmt.Sprintf("hv%d", i)] = fmt.Sprintf("value%d", i)
	}

	InitTmpls(&req)

	request := GenHTTPReq(data)

	t.Log("url:", request.URL)
	t.Log("body:", request.Body)
	t.Log("query:", request.QueryParam)
	t.Log("headers:", request.Header)
}
