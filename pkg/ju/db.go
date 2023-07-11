package ju

import (
	"os"
	"sync"

	"github.com/ejuju/go-trie-db/pkg/keyrefs"
)

type File struct {
	mu     sync.RWMutex
	refs   keyrefs.Map
	r, w   *os.File
	offset int // current write offset for new rows
}

func OpenFile(fpath string) (*File, error) {
	file := &File{refs: keyrefs.NewRadixMap()}
	var err error

	// Open read-only file handle (create file if needed) and write-only file handle (in append mode).
	if file.r, err = os.OpenFile(fpath, os.O_RDONLY|os.O_CREATE, os.ModePerm); err != nil {
		return nil, err
	}
	if file.w, err = os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND, os.ModePerm); err != nil {
		return nil, err
	}

	// Todo: Extract existing file data and store in memory

	return file, nil
}

func (f *File) Set(k, v []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	panic("not implemented yet")
}
