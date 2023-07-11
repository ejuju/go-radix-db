package format

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/ejuju/go-trie-db/pkg/keyrefs"
)

/*
Binary encodes rows to a file in a binary (non human-readable) format.

# Pseudo-code representation

	[]byte{op, keyLengthAsUint8 + valueLengthAsBigEndianUint32, key, value}

# Limitations:
  - Maximum key-length is 255
  - Maximum value-length is 4294967295 (=~ 4.2 GB)
*/
type Binary struct {
	PutOp    byte
	DeleteOp byte
}

var DefaultBinaryFormat = Binary{
	PutOp:    'P',
	DeleteOp: 'D',
}

const MaxBinaryKeyLength = 1<<8 - 1
const MaxBinaryValueLength = 1<<32 - 1

var (
	ErrKeyTooLong   = errors.New("key too long")
	ErrValueTooLong = errors.New("value too long")
	ErrUnknownOp    = errors.New("unknown op")
)

func (f Binary) encodeRow(op byte, k, v []byte) ([]byte, error) {
	if len(k) > MaxBinaryKeyLength {
		return nil, fmt.Errorf("%w: %d (max %d)", ErrKeyTooLong, len(k), MaxBinaryKeyLength)
	}
	if len(v) > MaxBinaryValueLength {
		return nil, fmt.Errorf("%w: %d (max %d)", ErrValueTooLong, len(k), MaxBinaryValueLength)
	}
	out := []byte{op}
	out = append(out, uint8(len(k)))
	out = binary.BigEndian.AppendUint32(out, uint32(len(v)))
	out = append(out, k...)
	out = append(out, v...)
	return out, nil
}

func (f Binary) EncodePutRow(k, v []byte) ([]byte, error) { return f.encodeRow(f.PutOp, k, v) }
func (f Binary) EncodeDeleteRow(k []byte) ([]byte, error) { return f.encodeRow(f.DeleteOp, k, nil) }

func (f Binary) Extract(r io.Reader, m keyrefs.Map) (int, error) {
	offset := 0

	for {
		// Read first 6 bytes (1 for op + 1 for keylength + 4 for valuelength)
		header := make([]byte, 6)
		n, err := io.ReadFull(r, header)
		if errors.Is(err, io.EOF) && n == 0 {
			break // Simply exit the loop if EOF and zero bytes were read (row hasn't started)
		}
		if err != nil {
			return n, fmt.Errorf("read row header at offset %d: %w", n, err)
		}
		offset += n

		// Retrieve op (first byte of header) and validate
		op := header[0]
		if !(op == f.PutOp || op == f.DeleteOp) {
			return n, fmt.Errorf("%w: %q", ErrUnknownOp, op)
		}

		kLen := uint8(header[1])                    // Retrieve key-length (second byte of header)
		vLen := binary.BigEndian.Uint32(header[2:]) // Retrieve value-length (last 8 bytes of header)

		// Read key
		k := make([]byte, kLen)
		n, err = io.ReadFull(r, k)
		if err != nil {
			return n, fmt.Errorf("read key at offset %d: %w", n, err)
		}
		offset += n

		// Read value
		offsetBeforeValue := offset
		v := make([]byte, vLen)
		n, err = io.ReadFull(r, v)
		if err != nil {
			return n, fmt.Errorf("read value at offset %d: %w", n, err)
		}
		offset += n

		// If delete-op, any possible past keyref should be deleted, then continue to next row.
		if header[0] == f.DeleteOp {
			err = m.Delete(k)
			if err != nil {
				return n, fmt.Errorf("found delete key %q: %w", k, err)
			}
			continue
		}

		// If set op, set the keyref for this key.
		err = m.Set(k, &keyrefs.ByteRange{Index: offsetBeforeValue, Width: len(v)})
		if err != nil {
			return n, fmt.Errorf("found set key %q: %w", k, err)
		}
	}

	return offset, nil
}
