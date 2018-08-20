package new

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"gopkg.in/resty.v1"
)

var (
	qlen = 200
	hlen = 40
	req  = Req{
		Url:  "http://localhost:{{.port}}",
		Body: `{"data":"{{.body}}"}`,
	}

	data = map[string]interface{}{
		"port": 8080,
		"body": "hello world",
	}
)

func init() {
	for i := 0; i < qlen; i++ {
		req.Query = append(req.Query, KV{fmt.Sprintf("k%d", i), fmt.Sprintf("{{.qv%d}}", i)})
		data[fmt.Sprintf("qv%d", i)] = fmt.Sprintf("value%d", i)
	}

	for i := 0; i < hlen; i++ {
		req.Headers = append(req.Headers, KV{fmt.Sprintf("k%d", i), fmt.Sprintf("{{.hv%d}}", i)})
		data[fmt.Sprintf("hv%d", i)] = fmt.Sprintf("value%d", i)
	}

	InitTmpls(&req)
}

// Req is HTTP request settings
type Req struct {
	Query   []KV
	Headers []KV
	Url     string
	Body    string
}

type KV struct {
	Key   string
	Value string
}

var rtmpl *template.Template

// InitTmpls initialize templates
func InitTmpls(r *Req) {
	tmplBytes, _ := json.Marshal(r)
	rtmpl = template.Must(template.New("").Parse(string(tmplBytes)))
}

// GenHTTPReq generate a HTTP request
func GenHTTPReq(data interface{}) *resty.Request {
	req := resty.R()

	buf := bytes.NewBuffer(nil)
	if err := rtmpl.Execute(buf, data); nil != err {
		fmt.Println("render template error:", err)
		return nil
	}

	var r Req
	if err := json.Unmarshal(buf.Bytes(), &r); nil != err {
		fmt.Println("Decode request settings error:", err)
		return nil
	}

	for _, kv := range r.Query {
		req.SetQueryParam(kv.Key, kv.Value)
	}

	for _, kv := range r.Headers {
		req.SetHeader(kv.Key, kv.Value)
	}

	req.URL = r.Url
	req.SetBody(r.Body)

	return req
}
