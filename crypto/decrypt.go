package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "errors"
)

func DecryptShard(encryptedShard []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    if len(encryptedShard) < aes.BlockSize {
        return nil, errors.New("ciphertext too short")
    }

    iv := encryptedShard[:aes.BlockSize]
    ciphertext := encryptedShard[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(ciphertext, ciphertext)

    return ciphertext, nil
}
