package old

import (
	"bytes"
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

var qTmpls, hTmpls []*template.Template
var uTmpl, bTmpl *template.Template

// InitTmpls initialize templates
func InitTmpls(r *Req) {
	loadTmpl := func(s string) *template.Template {
		return template.Must(template.New("").Parse(s))
	}

	for _, q := range r.Query {
		qTmpls = append(qTmpls, loadTmpl(q.Value))
	}

	for _, h := range r.Headers {
		hTmpls = append(hTmpls, loadTmpl(h.Value))
	}

	uTmpl = loadTmpl(r.Url)
	bTmpl = loadTmpl(r.Body)
}

// GenHTTPReq generate a HTTP request
func GenHTTPReq(r *Req, data interface{}) *resty.Request {
	req := resty.R()
	getParamValueFunc := func(def string, t *template.Template) (string, error) {
		if t == nil {
			return def, nil
		}
		buf := bytes.NewBuffer(nil)
		if err := t.Execute(buf, data); nil != err {
			return def, fmt.Errorf("Render parameter '%s' error: %v", def, err)
		}

		return buf.String(), nil
	}

	for i, kv := range r.Query {
		val, _ := getParamValueFunc(kv.Value, qTmpls[i])
		req.SetQueryParam(kv.Key, val)
	}

	for i, kv := range r.Headers {
		val, _ := getParamValueFunc(kv.Value, hTmpls[i])
		req.SetHeader(kv.Key, val)
	}

	req.URL, _ = getParamValueFunc(r.Url, uTmpl)
	b, _ := getParamValueFunc(r.Body, bTmpl)
	req.SetBody(b)

	return req
}
