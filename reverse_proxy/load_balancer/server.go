package load_balancer

import (
	"log"
	"net"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

type Server struct {
	URL          *url.URL
	Alive        uint32
	ReverseProxy *httputil.ReverseProxy
}

func (s *Server) IsAlive() bool {
	return s.Alive == 1
}

func (s *Server) SetAlive(alive bool) {
	var value uint32
	if alive {
		value = 1
	} else {
		value = 0
	}
	atomic.StoreUint32(&s.Alive, value)
}

func (s *Server) checkHealth() bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(s.URL.Hostname(), s.URL.Port()), 10*time.Second)
	if err != nil {
		log.Println("Server is unavailable:", err)
		s.SetAlive(false)
		return false
	}
	defer conn.Close()
	s.SetAlive(true)
	return true
}
