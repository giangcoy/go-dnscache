# go-dnscache
# Install

Install using the "go get" command:

```
go get -u github.com/giangcoy/go-dnscache
```

# Usage

Create httpClient with dnscache.DialContext

```go
c := &http.Client{Transport: &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext:           dnscache.DialContext,
	}}
```
