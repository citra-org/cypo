package storage

import (
	"os"
	"path/filepath"
)

func StoreShard(folderPath string, fileName string, data []byte) error {
    filePath := filepath.Join(folderPath, fileName)
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.Write(data)
    if err != nil {
        return err
    }

    return nil
}
