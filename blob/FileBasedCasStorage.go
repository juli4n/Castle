package blob

import (
  "crypto/sha1"
  "encoding/base32"
  "os"
)

type fileBasedBlobStorage struct {
  basePath string
}

func NewFileBasedBlobStorage(basePath string) BlobStorage {
  //FIXME: base path sanotizing
  return &fileBasedBlobStorage{basePath}
}

func fileAlreadyExists(filename string) bool {
  _, err := os.Open(filename)
  return err == nil
}

func (self fileBasedBlobStorage) Put(content []byte) (id string, err error) {
  var file *os.File
  hash := sha1.New()
  hash.Write(content)
  id = base32.StdEncoding.EncodeToString((hash.Sum([]byte(""))))
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

func (self fileBasedBlobStorage) Get(id string) (data []byte, err error) {
  var file *os.File
  if file, err = os.Open(self.basePath + id); err != nil {
    return nil, err
  }
  defer file.Close()
  stats,_ := file.Stat()
  data = make([]byte, stats.Size())
  if _, err = file.Read(data); err != nil {
    return nil, err
  }
  return data, nil
}
