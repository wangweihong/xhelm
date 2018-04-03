package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func Get(url string, token string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("New request for url [%v] fail for %v", url, err)
	}
	//测试使用
	//	request.Header.Add("token", "1234567890987654321")
	request.Header.Add("token", token)

	client := http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request url [%v] fail for %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("request url [%v] fail for %v", url, resp.Status)
		} else {
			return nil, fmt.Errorf("request url [%v] fail for %v", url, string(body))
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse body from url [%v] request's body fail for %v", url, err)
	}
	return body, nil
}

func Post(url string, token string, data string) ([]byte, error) {
	dio := strings.NewReader(data)
	request, err := http.NewRequest("Post", url, dio)
	if err != nil {
		return nil, fmt.Errorf("New request for url [%v] fail for %v", url, err)
	}
	//测试使用
	//	request.Header.Add("token", "1234567890987654321")
	request.Header.Add("token", token)

	client := http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request url [%v] fail for %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("request url [%v] fail for %v", url, resp.Status)
		} else {
			return nil, fmt.Errorf("request url [%v] fail for %v", url, string(body))
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse body from url [%v] request's body fail for %v", url, err)
	}
	return body, nil
}
