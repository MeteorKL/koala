package koala

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetRequest(URL string) (int, []byte) {
	resp, err := http.Get(URL)
	if err != nil {
		return resp.StatusCode, []byte(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, []byte(err.Error())
	}
	defer resp.Body.Close()
	return resp.StatusCode, body
}

func PostRequest(URL string, param map[string]string) (int, []byte) {
	query := url.Values{}
	for k, v := range param {
		query.Set(k, v)
	}
	resp, err := http.PostForm(URL, query)
	if err != nil {
		return resp.StatusCode, []byte(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, []byte(err.Error())
	}
	defer resp.Body.Close()
	return resp.StatusCode, body
}

var client = http.Client{}

func Request(method, url string, param string) []byte {
	req, err := http.NewRequest(method, url, strings.NewReader(param))
	if err != nil {
		return []byte(err.Error())
	}
	if method != "GET" && method != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	if err != nil {
		return []byte(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []byte(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(err.Error())
	}
	return b
}
