package object

import (
	"bytes"
	"castle/blob"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
)

type ObjectStorage interface {
	Get(string, interface{}) error
	Put(interface{}) (string, error)
}

type Encoder interface {
	Encode(interface{}) error
}

type EncoderFactory func(io.Writer) Encoder

type Decoder interface {
	Decode(interface{}) error
}

type DecoderFactory func(io.Reader) Decoder

type delegateObjectStorage struct {
	blobStorage    blob.BlobStorage
	encoderFactory EncoderFactory
	decoderFactory DecoderFactory
}

func NewStorage(s blob.BlobStorage, enc EncoderFactory, dec DecoderFactory) ObjectStorage {
	return &delegateObjectStorage{s, enc, dec}
}

func NewGobStorage(s blob.BlobStorage) ObjectStorage {
	return &delegateObjectStorage{s,
		func(w io.Writer) Encoder { return gob.NewEncoder(w) },
		func(r io.Reader) Decoder { return gob.NewDecoder(r) }}
}

func NewXMLStorage(s blob.BlobStorage) ObjectStorage {
	return &delegateObjectStorage{s,
		func(w io.Writer) Encoder { return xml.NewEncoder(w) },
		func(r io.Reader) Decoder { return xml.NewDecoder(r) }}
}

func NewJSONStorage(s blob.BlobStorage) ObjectStorage {
	return &delegateObjectStorage{s,
		func(w io.Writer) Encoder { return json.NewEncoder(w) },
		func(r io.Reader) Decoder { return json.NewDecoder(r) }}
}

func (self *delegateObjectStorage) Put(value interface{}) (objRef string, err error) {
	buffer := &bytes.Buffer{}
	encoder := self.encoderFactory(buffer)
	if err := encoder.Encode(value); err != nil {
		return "", err
	}
	if objRef, err = self.blobStorage.Put(buffer.Bytes()); err != nil {
		return "", err
	}
	return objRef, err
}

func (self *delegateObjectStorage) Get(objRef string, value interface{}) (err error) {
	var content []byte
	content, err = self.blobStorage.Get(objRef)
	decoder := self.decoderFactory(bytes.NewReader(content))
	if err = decoder.Decode(value); err != nil {
		return err
	}
	return nil
}
