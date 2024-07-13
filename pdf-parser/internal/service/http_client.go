package service

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClientWrapper struct {
	client *http.Client
}

const TIMEOUT = 10

func NewHttpClient() *HttpClientWrapper {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Second * time.Duration(TIMEOUT),
	}
	return &HttpClientWrapper{
		client: client,
	}
}

func (hc *HttpClientWrapper) DownloadPDF(url string) (*bytes.Buffer, error) {
	resp, err := hc.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download PDF: received status code %d", resp.StatusCode)
	}

	var buffer bytes.Buffer

	// Copy the response body to the buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy PDF content to buffer: %w", err)
	}

	return &buffer, nil
}
