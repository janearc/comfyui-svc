package main

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

func TestHealthAndMetricsEndpoints(t *testing.T) {
	targetURL, _ := url.Parse("http://localhost:8188")
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	mux := setupMux(proxy)

	// Test Unhealthy State
	atomic.StoreInt32(&comfyuiHealthy, 0)
	
	reqHealth := httptest.NewRequest(http.MethodGet, "/health", nil)
	rrHealth := httptest.NewRecorder()
	mux.ServeHTTP(rrHealth, reqHealth)
	
	if rrHealth.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 for /health, got %v", rrHealth.Code)
	}

	reqMetrics := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rrMetrics := httptest.NewRecorder()
	mux.ServeHTTP(rrMetrics, reqMetrics)

	if rrMetrics.Code != http.StatusOK {
		t.Errorf("Expected 200 for /metrics, got %v", rrMetrics.Code)
	}

	// Test Healthy State
	atomic.StoreInt32(&comfyuiHealthy, 1)

	reqHealthHealthy := httptest.NewRequest(http.MethodGet, "/health", nil)
	rrHealthHealthy := httptest.NewRecorder()
	mux.ServeHTTP(rrHealthHealthy, reqHealthHealthy)

	if rrHealthHealthy.Code != http.StatusOK {
		t.Errorf("Expected 200 for /health when healthy, got %v", rrHealthHealthy.Code)
	}
}

func TestMonitorComfyUI(t *testing.T) {
	// Create mock ComfyUI server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/system_stats" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	atomic.StoreInt32(&comfyuiHealthy, 0)
	go monitorComfyUI(mockServer.URL)
	time.Sleep(100 * time.Millisecond)

	// Test failure path
	mockServerFail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServerFail.Close()

	go monitorComfyUI(mockServerFail.URL)
	time.Sleep(100 * time.Millisecond)
}

func TestRun(t *testing.T) {
	// Test invalid URL
	if err := run(":0", "://invalid-url"); err == nil {
		t.Error("Expected error for invalid target URL")
	}

	// Test valid run, it should start server but we don't want it to block forever,
	// so we use a random port, let it spin up, and then proceed. 
	// We can't cleanly stop http.ListenAndServe without a server object, so we'll 
	// just spawn it and let it hang in the background for this coverage test.
	go run(":0", "http://localhost:8188")
	time.Sleep(100 * time.Millisecond)
}
