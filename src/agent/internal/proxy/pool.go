package proxy

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Proxy struct {
	URL      string
	Username string
	Password string
	LastUsed time.Time
	Failures int
}

type ProxyPool struct {
	proxies     []*Proxy
	current     int
	mu          sync.RWMutex
	maxFailures int
}

func NewProxyPool() *ProxyPool {
	return &ProxyPool{
		proxies:     []*Proxy{},
		current:     0,
		maxFailures: 3,
	}
}

func (p *ProxyPool) LoadFromEnv() {
	proxyList := os.Getenv("HTTP_PROXIES")
	if proxyList == "" {
		proxyList = os.Getenv("http_proxy")
	}
	if proxyList == "" {
		slog.Info("no proxies configured")
		return
	}

	urls := splitProxyList(proxyList)
	for _, url := range urls {
		proxy := &Proxy{URL: url}
		p.proxies = append(p.proxies, proxy)
	}

	slog.Info("loaded proxies", "count", len(p.proxies))
}

func splitProxyList(proxyList string) []string {
	var urls []string
	for _, p := range []string{",", ";", " "} {
		if contains(proxyList, p) {
			urls = split(proxyList, p)
			break
		}
	}
	if len(urls) == 0 {
		urls = []string{proxyList}
	}
	return urls
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func split(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + len(sep)
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

func (p *ProxyPool) Get() *Proxy {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.proxies) == 0 {
		return nil
	}

	proxy := p.proxies[p.current]
	p.current = (p.current + 1) % len(p.proxies)
	proxy.LastUsed = time.Now()

	return proxy
}

func (p *ProxyPool) GetWorking() *Proxy {
	p.mu.Lock()
	defer p.mu.Unlock()

	var working []*Proxy
	for _, proxy := range p.proxies {
		if proxy.Failures < p.maxFailures {
			working = append(working, proxy)
		}
	}

	if len(working) == 0 {
		// Reset failures and return random
		for _, proxy := range p.proxies {
			proxy.Failures = 0
		}
		return p.proxies[rand.Intn(len(p.proxies))]
	}

	proxy := working[rand.Intn(len(working))]
	proxy.LastUsed = time.Now()
	return proxy
}

func (p *ProxyPool) MarkFailure(proxy *Proxy) {
	p.mu.Lock()
	defer p.mu.Unlock()

	proxy.Failures++
	slog.Warn("proxy marked failed", "url", proxy.URL, "failures", proxy.Failures)
}

func (p *ProxyPool) MarkSuccess(proxy *Proxy) {
	p.mu.Lock()
	defer p.mu.Unlock()

	proxy.Failures = 0
}

func (p *ProxyPool) Add(proxy *Proxy) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.proxies = append(p.proxies, proxy)
}

func (p *ProxyPool) Remove(url string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, proxy := range p.proxies {
		if proxy.URL == url {
			p.proxies = append(p.proxies[:i], p.proxies[i+1:]...)
			break
		}
	}
}

func (p *ProxyPool) Count() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.proxies)
}

func (p *ProxyPool) RoundTrip(ctx context.Context, req *http.Request) (*http.Response, error) {
	proxy := p.GetWorking()
	if proxy == nil {
		return nil, fmt.Errorf("no working proxies available")
	}

	slog.Debug("using proxy", "url", proxy.URL)

	transport := &http.Transport{
		Proxy: func(*http.Request) (*url.URL, error) {
			return url.Parse(proxy.URL)
		},
	}

	client := &http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		p.MarkFailure(proxy)
		return nil, err
	}

	p.MarkSuccess(proxy)
	return resp, nil
}
