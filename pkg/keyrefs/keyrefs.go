package keyrefs

type ByteRange struct {
	Index int
	Width int
}

// Map is responsible for managing the various byte range (on file) associated with a given keys.
// It mimicks how a Go map would behave (except that the keyrefs.Map can sorts keys in a lexicographical order).
type Map interface {
	Set(key []byte, br *ByteRange) error
	Delete(key []byte) error
	Get(key []byte) *ByteRange
	Walk(WalkFunc)
	NumKeys() int
	IsLexOrdered() bool
}

type WalkFunc func(key []byte, br *ByteRange) bool
