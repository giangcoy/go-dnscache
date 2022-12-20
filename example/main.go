package main

import (
	"net/http"
	"time"

	"github.com/giangcoy/go-dnscache"
)

func main() {
	url := "https://www.google.com/"
	c := &http.Client{Transport: &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext:           dnscache.DialContext,
	}}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
	//c.Do(req)

}
