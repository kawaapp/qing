package httpx

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
	"time"
)

func Get(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var (
			statusCode = resp.StatusCode
			message, _ = ioutil.ReadAll(resp.Body)
		)
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v, message=%v", uri, statusCode, string(message))
	}
	return ioutil.ReadAll(resp.Body)
}

func Post(uri string, data []byte) ([]byte, error) {
	resp, err := http.Post(uri, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var (
			statusCode = resp.StatusCode
			message, _ = ioutil.ReadAll(resp.Body)
		)
		return nil, fmt.Errorf("http post error : uri=%v , statusCode=%v, message=%v", uri, statusCode, string(message))
	}
	return ioutil.ReadAll(resp.Body)
}

func GetHeader(uri string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var (
			statusCode = resp.StatusCode
			message, _ = ioutil.ReadAll(resp.Body)
		)
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v, message=%v", uri, statusCode, string(message))
	}
	return ioutil.ReadAll(resp.Body)
}

func PostHeader(uri string, header map[string]string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var (
			statusCode = resp.StatusCode
			message, _ = ioutil.ReadAll(resp.Body)
		)
		return nil, fmt.Errorf("http post error : uri=%v , statusCode=%v, message=%v", uri, statusCode, string(message))
	}
	return ioutil.ReadAll(resp.Body)
}


var netClient = &http.Client{
	Timeout: time.Second * 10,
}