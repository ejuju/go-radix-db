package keyrefs

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func testMap(t *testing.T, newMap func() Map) {
	t.Helper()

	t.Run("can set and get back data", func(t *testing.T) {
		k := []byte("mykey")
		v := &ByteRange{Index: 13, Width: 13}

		m := newMap()
		err := m.Set(k, v)
		if err != nil {
			panic(err)
		}

		gotV := m.Get(k)
		if !reflect.DeepEqual(gotV, v) {
			t.Fatalf("got ref %+v but wanted %+v", gotV, v)
		}
	})

	t.Run("can delete a key", func(t *testing.T) {
		k := []byte("mykey")
		m := newMap()

		err := m.Set(k, &ByteRange{Index: 13, Width: 13})
		if err != nil {
			panic(err)
		}

		err = m.Delete(k)
		if err != nil {
			panic(err)
		}

		gotV := m.Get(k)
		if gotV != nil {
			t.Fatalf("got ref %+v but wanted nil", gotV)
		}
	})

	t.Run("stores keys in lexicographical order", func(t *testing.T) {
		m := newMap()
		if !m.IsLexOrdered() {
			t.SkipNow()
		}

		setKeys := []string{"00", "10", "1", "2", "0"}  // unordered
		wantKeys := []string{"0", "00", "1", "10", "2"} // ordered
		for _, k := range setKeys {
			err := m.Set([]byte(k), &ByteRange{}) // file ref data is irrelevant here
			if err != nil {
				panic(err)
			}
		}

		gotKeys := []string{}
		m.Walk(func(key []byte, br *ByteRange) bool {
			gotKeys = append(gotKeys, string(key))
			return true
		})
		if !reflect.DeepEqual(gotKeys, wantKeys) {
			t.Fatalf("got keys %q but wanted ordered keys %q", gotKeys, wantKeys)
		}
	})

	t.Run("can count keys", func(t *testing.T) {
		m := newMap()
		setKeys := []string{"1", "2", "3"}
		for i, k := range setKeys {
			err := m.Set([]byte(k), &ByteRange{Index: i + 1})
			if err != nil {
				panic(err)
			}
		}

		gotCount := m.NumKeys()
		if gotCount != len(setKeys) {
			t.Fatalf("got count %d but wanted %d", gotCount, len(setKeys))
		}
	})
}

func benchmarkMap(b *testing.B, newMap func() Map) {
	b.Helper()

	b.Run("set 1000 new keys (length 16)", func(b *testing.B) {
		m := newMap()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			for i := 0; i < 1000; i++ {
				k := []byte(fmt.Sprintf("%016s", strconv.Itoa(i)))
				br := &ByteRange{Index: i}
				err := m.Set(k, br)
				if err != nil {
					panic(err)
				}
			}
		}
	})

	b.Run("get a known key (length 16)", func(b *testing.B) {
		m := newMap()
		key := []byte("1234567890123456")
		err := m.Set(key, &ByteRange{})
		if err != nil {
			panic(err)
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = m.Get(key)
		}
	})

	b.Run("get an unknown key (length 16)", func(b *testing.B) {
		key := []byte("1234567890123456")
		m := newMap()

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = m.Get(key)
		}
	})

	b.Run("walk 1_000 keys", func(b *testing.B) {
		m := newMap()
		for i := 0; i < 1_000; i++ {
			err := m.Set([]byte(fmt.Sprintf("%05s", strconv.Itoa(i))), &ByteRange{Index: i})
			if err != nil {
				panic(err)
			}
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			m.Walk(func(_ []byte, _ *ByteRange) bool { return true })
		}
	})
}
