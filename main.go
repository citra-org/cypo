package main

import (
	"crypto/sha256"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/citra-org/cypo/crypto"
	"github.com/citra-org/cypo/processing"
	"github.com/citra-org/cypo/storage"
	"golang.org/x/term"
)

func main() {
    fmt.Print("Enter 'e' to encrypt or 'd' to decrypt: ")
    var operation string
    fmt.Scan(&operation)

    if operation == "e" {
        encryptImages()
    } else if operation == "d" {
        decryptImages()
    } else {
        fmt.Println("Invalid operation. Please enter 'e' or 'd'.")
    }
}

func encryptImages() {
    fmt.Print("Enter a password for encryption: ")
    bytePassword, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatalf("Failed to read password: %v", err)
    }
    fmt.Println() 
    key := sha256.Sum256(bytePassword)

    folderPath := "."

    
    encryptedBaseDir := "encrypted"
    err = os.MkdirAll(encryptedBaseDir, os.ModePerm)
    if err != nil {
        log.Fatalf("Failed to create encrypted directory: %v", err)
    }

    err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }


        if !info.IsDir() && isImageFile(path) {
            fmt.Printf("Processing image: %s\n", path)
            shards, err := processing.ShardImage(path, 4) 
            if err != nil {
                return err
            }

            for i, shard := range shards {
                encryptedShard, err := crypto.EncryptShard(shard, key[:])
                if err != nil {
                    return err
                }

                subfolder := fmt.Sprintf("%s/branch_%d", encryptedBaseDir, i)
                err = os.MkdirAll(subfolder, os.ModePerm)
                if err != nil {
                    return err
                }
                err = storage.StoreShard(subfolder, fmt.Sprintf("shard_%d.enc", i), encryptedShard)
                if err != nil {
                    return err
                }
            }
        }

        return nil
    })

    if err != nil {
        log.Fatalf("Error processing folder: %v", err)
    }

    fmt.Println("All images processed, sharded, encrypted, and stored successfully.")
}

func decryptImages() {
    fmt.Print("Enter the password for decryption: ")
    bytePassword, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatalf("Failed to read password: %v", err)
    }
    fmt.Println() 
    key := sha256.Sum256(bytePassword)
    encryptedBaseDir := "encrypted"
    decryptedBaseDir := "decrypted"
    err = os.MkdirAll(decryptedBaseDir, os.ModePerm)
    if err != nil {
        log.Fatalf("Failed to create decrypted directory: %v", err)
    }

    shards := make(map[int][]byte)
    err = filepath.Walk(encryptedBaseDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && strings.HasSuffix(info.Name(), ".enc") {
            fmt.Printf("Decrypting shard: %s\n", path)
            shardData, err := os.ReadFile(path)
            if err != nil {
                return err
            }

            decryptedShard, err := crypto.DecryptShard(shardData, key[:])
            if err != nil {
                return err
            }
            var shardIndex int
            fmt.Sscanf(info.Name(), "shard_%d.enc", &shardIndex)

            shards[shardIndex] = decryptedShard
        }

        return nil
    })

    if err != nil {
        log.Fatalf("Error processing folder: %v", err)
    }

    if len(shards) > 0 {
        reconstructedImage, err := processing.ReconstructImage(shards)
        if err != nil {
            log.Fatalf("Failed to reconstruct image: %v", err)
        }

        outputFilePath := filepath.Join(decryptedBaseDir, "reconstructed_image.png")
        outputFile, err := os.Create(outputFilePath)
        if err != nil {
            log.Fatalf("Failed to create output file: %v", err)
        }
        defer outputFile.Close()

        err = png.Encode(outputFile, reconstructedImage)
        if err != nil {
            log.Fatalf("Failed to save reconstructed image: %v", err)
        }

        fmt.Printf("Image successfully decrypted and saved to: %s\n", outputFilePath)
    } else {
        fmt.Println("No shards found to decrypt.")
    }
}

func isImageFile(filePath string) bool {
    ext := strings.ToLower(filepath.Ext(filePath))
    switch ext {
    case ".jpg", ".jpeg", ".png":
        return true
    }
    return false
}