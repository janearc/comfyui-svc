package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync/atomic"
	"time"
)

var comfyuiHealthy int32 = 0

func monitorComfyUI(targetStr string) {
	for {
		resp, err := http.Get(targetStr + "/system_stats")
		if err == nil && resp.StatusCode == 200 {
			atomic.StoreInt32(&comfyuiHealthy, 1)
		} else {
			atomic.StoreInt32(&comfyuiHealthy, 0)
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(30 * time.Second)
	}
}

func setupMux(proxy *httputil.ReverseProxy) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if atomic.LoadInt32(&comfyuiHealthy) == 1 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"error", "reason":"ComfyUI API unreachable"}`))
		}
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "# HELP comfyui_up Whether ComfyUI is up\n# TYPE comfyui_up gauge\ncomfyui_up %d\n", atomic.LoadInt32(&comfyuiHealthy))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	return mux
}

func run(addr, targetStr string) error {
	targetURL, err := url.Parse(targetStr)
	if err != nil {
		return err
	}

	go monitorComfyUI(targetStr)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	mux := setupMux(proxy)

	fmt.Println("Sidecar listening on", addr, "proxying to", targetStr)
	return http.ListenAndServe(addr, mux)
}

func main() {
	targetStr := os.Getenv("TARGET_URL")
	if targetStr == "" {
		targetStr = "http://localhost:8188"
	}
	if err := run(":8080", targetStr); err != nil {
		panic(err)
	}
}
