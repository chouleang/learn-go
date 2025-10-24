package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        hostname, _ := os.Hostname()
        fmt.Fprintf(w, "Hello, DevOps World!\n")
        fmt.Fprintf(w, "Container: %s\n", hostname)
        fmt.Fprintf(w, "🕐 Version: 2.0.0 - Automated GKE Deployment\n")
        fmt.Fprintf(w, "⏰ Server Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
        fmt.Fprintf(w, "🌍 Region: Singapore GKE\n")
        fmt.Fprintf(w, "🎯 Commit: Automated CI/CD Pipeline\n")
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    fmt.Printf("Server running (port=8080)\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
