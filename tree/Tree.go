package tree

import (
	"bytes"
	"encoding/gob"
	"gostore/blob"
)

// A tree of references to persisted objects.
type Tree struct {
	// Child nodes.
	Childs []TreeNode
}

// 
type TreeNode struct {
	// The name of the child.
	Name string
	// The id of the referenced object.
	ObjRef string
	// true iff is a leaf.
	IsLeaf bool
}

// A CAS storage that persists Tree instances.
type TreeStorage interface {

	// Returns a Tree for a given id or an error if a Tree
	// with that id doesn't exists.
	Get(id string) (*Tree, error)

	// Persists a Tree instance and returns its id. 
	Put(*Tree) (string, error)
}

type delegateTreeStorage struct {
	blobStorage blob.BlobStorage
}

func NewTreeStorage(storage blob.BlobStorage) TreeStorage {
	return &delegateTreeStorage{storage}
}

func (self *delegateTreeStorage) Put(tree *Tree) (objRef string, err error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(tree); err != nil {
		return "", err
	}
	if objRef, err = self.blobStorage.Put(buffer.Bytes()); err != nil {
		return "", err
	}
	return objRef, err
}

func (self *delegateTreeStorage) Get(objRef string) (tree *Tree, err error) {
	var content []byte
	content, err = self.blobStorage.Get(objRef)
	decoder := gob.NewDecoder(bytes.NewReader(content))
	if err = decoder.Decode(&tree); err != nil {
		return nil, err
	}
	return tree, nil
}
