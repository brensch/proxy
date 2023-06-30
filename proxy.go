package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	TargetHeader = "X-Target"
)

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.Header.Get(TargetHeader)
	if target == "" {
		http.Error(w, "X-Target header is missing", http.StatusBadRequest)
		return
	}

	// Create new URL
	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Invalid X-Target header", http.StatusBadRequest)
		return
	}

	// Create new request
	newRequest := r.Clone(r.Context())
	newRequest.URL.Scheme = targetURL.Scheme
	newRequest.URL.Host = targetURL.Host

	// Calculate the Host header value
	host := targetURL.Host
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Set the Host header
	newRequest.Host = host

	// use a randomised useragent
	newRequest.Header.Del(TargetHeader)
	newRequest.Header.Del("Authorization")
	newRequest.Header.Set("User-Agent", RandomUserAgent())

	// Dump all headers
	dumpHeaders(newRequest)

	resp, err := http.DefaultTransport.RoundTrip(newRequest)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// Copy headers and body
	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func dumpHeaders(r *http.Request) {
	for name, headers := range r.Header {
		for _, h := range headers {
			log.Printf("%v: %v\n", name, h)
		}
	}
}
