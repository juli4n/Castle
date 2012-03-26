package blob

import (
	"crypto/sha1"
	"encoding/base32"
	"os"
)

// A simple implementation of BlobStorage that
// stores blobs in the filesystem.
// All blobs are stored in a single flat directory.
type fileBasedBlobStorage struct {
	basePath string
}

// A factory for fileBasedBlobStorage instances
func NewFileBasedBlobStorage(basePath string) BlobStorage {
	return &fileBasedBlobStorage{basePath}
}

func (self fileBasedBlobStorage) Put(content []byte) (id string, err error) {
	var file *os.File
	hash := sha1.New()
	hash.Write(content)
	id = base32.StdEncoding.EncodeToString(hash.Sum([]byte{}))
	filename := self.basePath + id
	if fileAlreadyExists(filename) {
		return id, nil
	}
	if file, err = os.Create(self.basePath + id); err != nil {
		return "", err
	}
	defer file.Close()
	if _, err = file.Write(content); err != nil {
		return "", err
	}
	return id, nil
}

func fileAlreadyExists(filename string) bool {
	_, err := os.Open(filename)
	return err == nil
}

func (self fileBasedBlobStorage) Get(id string) (data []byte, err error) {
	var file *os.File
	if file, err = os.Open(self.basePath + id); err != nil {
		return nil, err
	}
	defer file.Close()
	stats, _ := file.Stat()
	data = make([]byte, stats.Size())
	if _, err = file.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}
