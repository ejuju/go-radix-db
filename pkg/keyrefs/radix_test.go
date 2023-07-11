package keyrefs

import (
	"testing"
)

func TestRadixMap(t *testing.T) { testMap(t, NewRadixMap) }

func BenchmarkRadixMap(b *testing.B) { benchmarkMap(b, NewRadixMap) }
