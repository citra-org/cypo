package main

import (
    "log"
    "net/http"
    "time"
    "github.com/citra-org/cypo/config"
    "github.com/citra-org/cypo/auth"
    "github.com/citra-org/cypo/file"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        if r.URL.Path == "/" {
            auth.HandleOTP(w, r)
        } else if r.URL.Path == "/upload" {
            file.HandleFileUpload(w, r)
        }
        return
    }

    config.Mu.Lock()
    if config.Authenticated {
        config.Mu.Unlock()
        http.ServeFile(w, r, "assets/upload.html")
        return
    }
    config.Mu.Unlock()
    http.ServeFile(w, r, "assets/otp.html")
}

func main() {
    config.OtpIssued = time.Now()

    http.HandleFunc("/", handleRequest)

    log.Println("Server is listening on port", config.Port)
    if err := http.ListenAndServe(config.Port, nil); err != nil {
        log.Fatal("Server failed:", err)
    }
}
