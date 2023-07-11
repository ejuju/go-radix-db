package keyrefs

type ByteRange struct {
	Index int
	Width int
}

// Map is responsible for managing the variou byte range (on file) associated with a given keys.
// It mimicks how a Go map would behave except that the keyrefs.Map sorts keys in a lexicographical order.
type Map interface {
	Set(key []byte, br *ByteRange) error
	Delete(key []byte) error
	Get(key []byte) *ByteRange
	Walk(WalkFunc)
	IsLexOrdered() bool
}

type WalkFunc func(key []byte, br *ByteRange) bool
