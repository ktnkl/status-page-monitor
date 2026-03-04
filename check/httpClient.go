package check

import (
	"fmt"
	"net/http"
	"time"
)

type ServerChecker struct {
	client  *http.Client
	timeout time.Duration
}

type ServerStatus struct {
	URL          string
	IsActive     bool
	StatusCode   int
	ResponseTime time.Duration
	Error        string
}

func NewServerChecker(timeout time.Duration) *ServerChecker {
	return &ServerChecker{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     30 * time.Second,
				DisableKeepAlives:   false,
			},
		},
		timeout: timeout,
	}
}

func (sc *ServerChecker) CheckServer(url string) (*ServerStatus, error) {
	start := time.Now()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create resuest: %w", err)
	}

	req.Header.Set("User-Agent", "StatusPageMonitor")
	req.Header.Set("Accept", "*/*")

	resp, err := sc.client.Do(req)
	if err != nil {
		return &ServerStatus{
			URL:          url,
			IsActive:     false,
			Error:        err.Error(),
			ResponseTime: time.Since(start),
		}, err
	}

	isUp := resp.StatusCode >= 200 && resp.StatusCode < 300

	status := &ServerStatus{
		URL:          url,
		IsActive:     isUp,
		StatusCode:   resp.StatusCode,
		ResponseTime: time.Since(start),
	}

	return status, nil
}
