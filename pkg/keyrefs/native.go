package keyrefs

import (
	"bytes"
	"sort"
)

type GoMap map[string]*ByteRange

func NewGoMap() Map { return GoMap{} }

func (m GoMap) IsLexOrdered() bool                  { return false }
func (m GoMap) Set(key []byte, br *ByteRange) error { m[string(key)] = br; return nil }
func (m GoMap) Delete(key []byte) error             { delete(m, string(key)); return nil }
func (m GoMap) Get(key []byte) *ByteRange           { return m[string(key)] }
func (m GoMap) Walk(cb WalkFunc) {
	for k, v := range m {
		ok := cb([]byte(k), v)
		if !ok {
			return
		}
	}
}

type GoSliceMap []*goSliceMapItem

type goSliceMapItem struct {
	k []byte
	v *ByteRange
}

func NewGoSliceMap() Map { return &GoSliceMap{} }

func (m *GoSliceMap) IsLexOrdered() bool { return true }

func (m *GoSliceMap) Set(key []byte, br *ByteRange) error {
	for i, pair := range *m {
		if bytes.Equal(key, pair.k) {
			(*m)[i].v = br // found key, replace value
			return nil
		}
	}
	*m = append(*m, &goSliceMapItem{k: key, v: br})
	sort.Sort(m)
	return nil
}

func (m *GoSliceMap) Len() int           { return len(*m) }
func (m *GoSliceMap) Swap(i, j int)      { (*m)[i], (*m)[j] = (*m)[j], (*m)[i] }
func (m *GoSliceMap) Less(i, j int) bool { return bytes.Compare((*m)[i].k, (*m)[j].k) == -1 }

func (m *GoSliceMap) Delete(key []byte) error {
	for i, pair := range *m {
		if bytes.Equal(pair.k, key) {
			*m = append((*m)[:i], (*m)[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m GoSliceMap) Get(key []byte) *ByteRange {
	for _, pair := range m {
		if bytes.Equal(pair.k, key) {
			return pair.v
		}
	}
	return nil
}

func (m GoSliceMap) Walk(cb WalkFunc) {
	for _, pair := range m {
		ok := cb(pair.k, pair.v)
		if !ok {
			return
		}
	}
}
