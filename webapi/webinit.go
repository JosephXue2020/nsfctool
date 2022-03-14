package webapi

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0"

var queryUrl = "https://isisn.nsfc.gov.cn/egrantindex/funcindex/prjsearch-list"
var headers = map[string]string{
	"User-Agent": UA,
}

func init() {
	resp, err := get(queryUrl, headers, nil, 3)
	if err != nil {
		panic("与网站连接初始化失败.")
	}
	cookie := resp.Header.Get("Set-Cookie")
	if cookie == "" {
		panic("与网站连接初始化失败.")
	}
	headers["Cookie"] = cookie
	// fmt.Println(headers)
}

// start a request by GET method with headers
func get(url string, headers map[string]string, param map[string]string, timeout int) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if param != nil {
		q := req.URL.Query()
		for k, v := range param {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// start a request by POST method with headers
func postForm(url string, headers map[string]string, data url.Values, timeout int) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
