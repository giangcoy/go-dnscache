package dnscache

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/rs/dnscache"
)

var (
	// Resolver holds the DNS cache resolver.
	r = &dnscache.Resolver{}
	// idx is used for load balancing among resolved IP addresses.
	idx = uint64(0)
	// defaultDialContext is the default dial context with a 30-second timeout and keep-alive.
	defaultDialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext
	Unused = false
)

// init initializes a background goroutine to periodically refresh the DNS cache.
func init() {
	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(Unused)
		}
	}()
}

// DialContext resolves the host using the DNS cache and returns a connection.
func DialContext(ctx context.Context, network, addr string) (conn net.Conn, err error) {
	// Split host and port from the address.
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	// Lookup the host in the DNS cache.
	ips, err := r.LookupHost(ctx, host)
	if err != nil {
		return nil, err
	}

	// Calculate the length of the resolved IP addresses and get a unique index.
	l := uint64(len(ips))
	id := atomic.AddUint64(&idx, 1)

	// Iterate over the resolved IP addresses with load balancing.
	for i := uint64(0); i < l; i++ {
		ip := ips[(i+id)%l]
		conn, err = defaultDialContext(ctx, network, net.JoinHostPort(ip, port))
		if err == nil {
			return conn, err
		}
	}

	// If all attempts fail, fallback to the default dial context.
	return defaultDialContext(ctx, network, addr)
}
