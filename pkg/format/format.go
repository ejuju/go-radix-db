package format

import (
	"io"

	"github.com/ejuju/go-trie-db/pkg/keyrefs"
)

// Format represents an encoding format for storing rows in a file.
type Format interface {
	EncodePutRow(k, v []byte) (int, []byte, error) // int is value start index in output
	EncodeDeleteRow(k []byte) ([]byte, error)

	// Iterate over reader to find rows, set those rows in the map
	// and return the number of bytes read offset.
	Extract(io.Reader, keyrefs.Map) (int, error)
}
