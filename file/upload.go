package file

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "github.com/citra-org/cypo/config"
)

func HandleFileUpload(w http.ResponseWriter, r *http.Request) {
    config.Mu.Lock()
    if !config.Authenticated {
        config.Mu.Unlock()
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    config.Mu.Unlock()

    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }

    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    file, fileHeader, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Failed to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    filePath := "./" + fileHeader.Filename
    outFile, err := os.Create(filePath)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }
    defer outFile.Close()
    
    _, err = io.Copy(outFile, file)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }
    

    fmt.Fprintln(w, "File uploaded successfully.")
}

