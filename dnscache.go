package dnscache

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/rs/dnscache"
)

var (
	r                  = &dnscache.Resolver{}
	idx                = uint64(0)
	defaultDialContext = (&net.Dialer{Timeout: 30 * time.Second,
		KeepAlive: 30 * time.Second}).DialContext
)

func init() {
	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(true)
		}
	}()

}
func DialContext(ctx context.Context, network string, addr string) (conn net.Conn, err error) {

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ips, err := r.LookupHost(ctx, host)
	if err != nil {
		return nil, err
	}
	l := uint64(len(ips))
	id := atomic.AddUint64(&idx, 1)
	for i := uint64(0); i < l; i++ {
		ip := ips[(i+id)%l]
		conn, err = defaultDialContext(ctx, network, net.JoinHostPort(ip, port))
		if err == nil {
			return conn, err
		}
	}
	return defaultDialContext(ctx, network, addr)

}
