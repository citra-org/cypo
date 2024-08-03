package processing

import (
    "bytes"
    "errors"
    "image"
    "image/jpeg"
    "os"
)

func ShardImage(imagePath string, numShards int) ([][]byte, error) {
    file, err := os.Open(imagePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }

    shards := make([][]byte, numShards)
    buffer := new(bytes.Buffer)
    err = jpeg.Encode(buffer, img, nil)
    if err != nil {
        return nil, err
    }

    shardSize := buffer.Len() / numShards
    for i := 0; i < numShards; i++ {
        if i == numShards-1 {
            shards[i] = buffer.Bytes()[i*shardSize:]
        } else {
            shards[i] = buffer.Bytes()[i*shardSize : (i+1)*shardSize]
        }
    }

    return shards, nil
}

func ReconstructImage(shards map[int][]byte) (image.Image, error) {
    var fullImage []byte
    for i := 0; i < len(shards); i++ {
        shard, ok := shards[i]
        if !ok {
            return nil, errors.New("missing shard " + string(i))
        }
        fullImage = append(fullImage, shard...)
    }

    img, _, err := image.Decode(bytes.NewReader(fullImage))
    if err != nil {
        return nil, err
    }

    return img, nil
}
