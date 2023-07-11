package keyrefs

import (
	"bytes"
)

type GoMap map[string]*ByteRange

func NewGoMap() Map { return GoMap{} }

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

func (m GoMap) IsLexOrdered() bool { return false }
func (m GoMap) NumKeys() int       { return len(m) }

type GoSliceMap []*goSliceMapItem

type goSliceMapItem struct {
	k []byte
	v *ByteRange
}

func NewGoSliceMap() Map { return &GoSliceMap{} }

func (m *GoSliceMap) Set(key []byte, br *ByteRange) error {
	for i, pair := range *m {
		switch bytes.Compare(key, pair.k) {
		case 0:
			(*m)[i].v = br // found key, replace value
			return nil
		case -1:
			*m = append((*m)[:i], append([]*goSliceMapItem{{k: key, v: br}}, (*m)[i:]...)...)
			return nil
		}
	}
	*m = append(*m, &goSliceMapItem{k: key, v: br})
	return nil
}

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

func (m GoSliceMap) IsLexOrdered() bool { return true }
func (m GoSliceMap) NumKeys() int       { return len(m) }
