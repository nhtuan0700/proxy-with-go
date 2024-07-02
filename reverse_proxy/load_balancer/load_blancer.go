package load_balancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type LoadBalancer struct {
	servers []*Server
	current  uint32
}

func NewLoadBalancer(targets []string) (*LoadBalancer, error) {
	var servers []*Server
	for _, target := range targets {
		url, err := url.Parse(target)
		if err != nil {
			return nil, fmt.Errorf("failed to parse url: %v", err)
		}

		server := &Server{
			URL:          url,
			Alive:        1,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		}
		servers = append(servers, server)
	}

	lb := &LoadBalancer{servers: servers}
	go lb.CheckHealth()

	return lb, nil
}

func (lb *LoadBalancer) GetAvailableServer() *Server {
	for range lb.servers {
		current := atomic.AddUint32(&lb.current, 1)
		index := int(current % uint32(len(lb.servers)))
		if lb.servers[index].IsAlive() {
			return lb.servers[index]
		}
	}
	return nil
}

func (lb *LoadBalancer) CheckHealth() {
	for {
		for _, server := range lb.servers {
			go server.checkHealth()
		}
		time.Sleep(10 * time.Second)
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for range lb.servers {
		server := lb.GetAvailableServer()
		if server == nil {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}

		proxyReq := r.Clone(r.Context())
		proxyReq.URL.Scheme = server.URL.Scheme
		proxyReq.URL.Host = server.URL.Host
		proxyReq.URL.Path = r.URL.Path
		proxyReq.Host = server.URL.Host
		proxyReq.Header = r.Header

		// send request to server
		server.ReverseProxy.ServeHTTP(w, proxyReq)

		if !server.checkHealth() {
			continue
		} else {
			return
		}
	}
	http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
}
