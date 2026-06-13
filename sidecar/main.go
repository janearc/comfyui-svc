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

func main() {
	targetStr := os.Getenv("TARGET_URL")
	if targetStr == "" {
		targetStr = "http://localhost:8188"
	}

	targetURL, err := url.Parse(targetStr)
	if err != nil {
		panic(err)
	}

	go monitorComfyUI(targetStr)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

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

	fmt.Println("Sidecar listening on :8080, proxying to", targetStr)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
