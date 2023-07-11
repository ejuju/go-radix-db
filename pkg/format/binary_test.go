package format

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ejuju/go-trie-db/pkg/keyrefs"
)

func TestBinaryFormat(t *testing.T) {
	t.Run("can encode a put-row", func(t *testing.T) {
		row, err := DefaultBinaryFormat.EncodePutRow([]byte("MyK"), []byte("MyVal"))
		if err != nil {
			panic(err)
		}

		wantRow := []byte{
			'P',        // op
			3,          // key-length
			0, 0, 0, 5, // value-length (big-endian uint32)
			'M', 'y', 'K', // key
			'M', 'y', 'V', 'a', 'l', // value
		}
		if !bytes.Equal(row, wantRow) {
			t.Fatalf("got row %q but wanted %q", row, wantRow)
		}
	})

	t.Run("can encode a delete-row", func(t *testing.T) {
		row, err := DefaultBinaryFormat.EncodeDeleteRow([]byte("MyK"))
		if err != nil {
			panic(err)
		}

		wantRow := []byte{ //
			'D',        // op
			3,          // key-length
			0, 0, 0, 0, // value-length (big-endian uint32)
			'M', 'y', 'K', // key
			// no value
		}
		if !bytes.Equal(row, wantRow) {
			t.Fatalf("got row %q but wanted %q", row, wantRow)
		}
	})

	t.Run("can extract row refs", func(t *testing.T) {
		// Populate test file
		buf := []byte{}

		buf = append(buf, 'P', 1, 0, 0, 0, 3, '1', 'v', 'a', 'l')      // put key 1
		buf = append(buf, 'P', 1, 0, 0, 0, 3, '2', 'v', 'a', 'l')      // put key 2
		buf = append(buf, 'D', 1, 0, 0, 0, 0, '2')                     // delete key 2
		buf = append(buf, 'P', 1, 0, 0, 0, 1, '3', 'v')                // put key 3
		buf = append(buf, 'P', 1, 0, 0, 0, 4, '3', 'w', 'x', 'y', 'z') // update key 3

		// Extract values from test file to the keyrefs map.
		m := keyrefs.NewGoSliceMap()
		offset, err := DefaultBinaryFormat.Extract(bytes.NewReader(buf), m)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Ensure number of bytes read is correct
		if offset != len(buf) {
			t.Fatalf("got offset %d but wanted %d", offset, len(buf))
		}

		// Check key 1
		key1ref := m.Get([]byte{'1'})
		wantKey1ref := &keyrefs.ByteRange{Index: 7, Width: 3} // key 1 value position
		if key1ref == nil {
			t.Fatalf("should have ref for key '1'")
		}
		if !reflect.DeepEqual(key1ref, wantKey1ref) {
			t.Fatalf("got ref %+v but wanted %+v", key1ref, wantKey1ref)
		}

		// Ensure key 2 is deleted
		key2ref := m.Get([]byte{'2'})
		if key2ref != nil {
			t.Fatalf("should not get deleted ref for key '2' but got %+v", key2ref)
		}

		// Check key 3
		key3ref := m.Get([]byte{'3'})
		wantKey3ref := &keyrefs.ByteRange{Index: len(buf) - 4, Width: 4} // key 3 value position
		if key1ref == nil {
			t.Fatalf("should have ref for key '3'")
		}
		if !reflect.DeepEqual(key3ref, wantKey3ref) {
			t.Fatalf("got ref %+v but wanted %+v", key3ref, wantKey3ref)
		}
	})
}
