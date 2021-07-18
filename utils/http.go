package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func DefaultClient(proxyUrl string) (*http.Client, error) {
	transport := &http.Transport{
		DisableCompression: true,
		IdleConnTimeout:    30 * time.Second,
	}
	if proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}
	return &http.Client{
		Transport: transport,
		Timeout:   5 * time.Minute,
	}, nil
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
