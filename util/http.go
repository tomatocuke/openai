package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	httpDefaultTimeout = time.Second * 5
)

type httpBuilder struct {
	Url     string
	Method  string
	Timeout time.Duration
	Header  map[string]string
	Body    io.Reader
}

func HttpGet(url string, params map[string]string) *httpBuilder {
	var r httpBuilder
	r.Method = http.MethodGet
	if len(params) > 0 {
		r.Url = url + "?" + string(r.buildQuery(params))
	} else {
		r.Url = url
	}
	return &r
}

func HttpPostJson(url string, bodyData interface{}) *httpBuilder {
	var r httpBuilder
	r.Method = http.MethodPost
	r.Url = url
	r.Header = map[string]string{"Content-Type": "chatgptlication/json"}
	b, _ := json.Marshal(bodyData)
	r.Body = bytes.NewReader(b)
	return &r
}

// func HttpPostForm(url string, bodyData map[string]string) *httpBuilder {
// 	var r httpBuilder
// 	r.Method = http.MethodPost
// 	r.Url = url
// 	r.Header = map[string]string{"Content-Type": "chatgptlication/x-www-form-urlencoded"}
// 	bodyDataFormat := r.buildQuery(bodyData)
// 	r.Body = bytes.NewReader(bodyDataFormat)
// 	return &r
// }

func (r *httpBuilder) buildQuery(params map[string]string) []byte {
	if len(params) == 0 {
		return nil
	}

	var buffer bytes.Buffer
	var i int
	for k := range params {
		buffer.WriteString(k)
		buffer.WriteByte('=')
		buffer.WriteString(params[k])
		if i != len(params)-1 {
			buffer.WriteByte('&')
		}
		i++
	}

	return buffer.Bytes()
}

func (r *httpBuilder) SetTimeout(timeout time.Duration) *httpBuilder {
	r.Timeout = timeout
	return r
}

func (r *httpBuilder) AddHeader(key, val string) *httpBuilder {
	if len(r.Header) == 0 {
		r.Header = make(map[string]string)
	}
	r.Header[key] = val
	return r
}

func (r *httpBuilder) Do() ([]byte, error) {
	req, err := http.NewRequest(r.Method, r.Url, r.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Header {
		req.Header.Set(k, v)
	}
	if r.Timeout == 0 {
		r.Timeout = httpDefaultTimeout
	}

	// 发送请求
	client := &http.Client{Timeout: r.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (r *httpBuilder) DoTo(to interface{}) error {
	bs, err := r.Do()
	if err == nil && to != nil {
		err = json.Unmarshal(bs, to)
	}
	return err
}
