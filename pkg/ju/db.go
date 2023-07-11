package ju

import (
	"fmt"
	"os"
	"sync"

	"github.com/ejuju/go-trie-db/pkg/format"
	"github.com/ejuju/go-trie-db/pkg/keyrefs"
)

type File struct {
	mu      sync.RWMutex
	refs    keyrefs.Map
	fformat format.Format // file format
	r, w    *os.File
	offset  int // current write offset for new rows
}

func OpenFile(fpath string) (*File, error) {
	file := &File{refs: keyrefs.NewRadixMap(), fformat: format.DefaultBinaryFormat}
	var err error

	// Open read-only file handle (create file if needed) and write-only file handle (in append mode).
	if file.r, err = os.OpenFile(fpath, os.O_RDONLY|os.O_CREATE, os.ModePerm); err != nil {
		return nil, err
	}
	if file.w, err = os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND, os.ModePerm); err != nil {
		return nil, err
	}

	// Extract existing file data and store in memory
	file.offset, err = file.fformat.Extract(file.r, file.refs)
	if err != nil {
		return nil, fmt.Errorf("extract refs from file %s: %w", fpath, err)
	}

	return file, nil
}

func (f *File) Put(k, v []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Encode row
	vIndex, row, err := f.fformat.EncodePutRow(k, v)
	if err != nil {
		return fmt.Errorf("encode put-row: %w", err)
	}

	// Write row to file
	n, err := f.w.Write(row)
	if err != nil {
		if n == len(row) {
			// Todo: mitigate risk of a row being completely written to file but the write fails
			// and the row is not set in the keyrefs.
			panic(fmt.Errorf("corruption, wrote complete row bytes but failed write: %w", err))
		}
		if n > 0 {
			// Todo: mitigate risk of a row being partially written to file
			panic(fmt.Errorf("corruption, wrote %d row bytes but failed write: %w", n, err))
		}
		return err
	}

	// Set ref in-memory
	err = f.refs.Set(k, &keyrefs.ByteRange{Index: vIndex, Width: len(v)})
	if err != nil {
		// Todo: mitigate risk of row being written to file but no recorded in memory.
		panic(fmt.Errorf("corruption, wrote row to file but failed to set keyref: %w", err))
	}

	return nil
}
