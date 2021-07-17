package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func DefaultClient(proxyUrl string) (*http.Client, error) {
	transport := &http.Transport{DisableCompression: true}
	if proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}
	return &http.Client{Transport: transport}, nil
}

type HttpError struct {
	Status string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("http error: %s", e.Status)
}

func Request(method string, client *http.Client, url string, headers map[string]string, body io.Reader) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, &HttpError{Status: res.Status}
	}
	return res.Body, nil
}
