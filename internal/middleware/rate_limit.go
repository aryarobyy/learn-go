package middlewares

import (
	"log"
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type Client struct {
	limiter *rate.Limiter
}

var (
	clients = make(map[string]*Client)
	mu      sync.Mutex
)

func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func getClientLimiter() *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	ip := getLocalIP()

	if client, exists := clients[ip.String()]; exists {
		return client.limiter
	}

	limiter := rate.NewLimiter(5, 1)
	clients[ip.String()] = &Client{limiter: limiter}
	return limiter
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ip := r.RemoteAddr
		limiter := getClientLimiter()

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
