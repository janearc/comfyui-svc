package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	targetStr := os.Getenv("TARGET_URL")
	if targetStr == "" {
		targetStr = "http://localhost:8188"
	}

	targetURL, err := url.Parse(targetStr)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# HELP comfyui_up Whether ComfyUI is up\n# TYPE comfyui_up gauge\ncomfyui_up 1\n"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	fmt.Println("Sidecar listening on :8080, proxying to", targetStr)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
