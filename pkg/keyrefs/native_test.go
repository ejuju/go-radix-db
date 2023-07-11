package keyrefs

import "testing"

func TestGoMap(t *testing.T) { testMap(t, NewGoMap) }

func BenchmarkGoMap(b *testing.B) { benchmarkMap(b, NewGoMap) }

func TestGoSliceMap(t *testing.T) { testMap(t, NewGoSliceMap) }

func BenchmarkGoSliceMap(b *testing.B) { benchmarkMap(b, NewGoSliceMap) }
