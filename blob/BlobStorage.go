package blob

type Batchable interface {
	Start()

	End()
}

// An interface for content addressable blob storages
type BlobStorage interface {

	// Adds the content to the store and returns
	// the id of the new blob.
	Put(content []byte) (string, error)

	PutWithId(id string, content []byte) error

	// Returns a stored blob for a given id or error
	// if it doesn't exists.
	Get(id string) ([]byte, error)
}
