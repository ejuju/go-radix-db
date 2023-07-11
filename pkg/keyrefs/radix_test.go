package keyrefs

import (
	"testing"
)

func TestRadixMap(t *testing.T) {
	testMap(t, func() Map { return NewRadixMap() }, false)
}

func BenchmarkRadixMap(b *testing.B) {
	benchmarkMap(b, func() Map { return NewRadixMap() })
}
