package keyrefs

type RadixMap struct{ root *radixNode }

func NewRadixMap() Map { return &RadixMap{root: &radixNode{}} }

type radixNode struct {
	children [256]*radixNode // 256 children (all values a byte/uint8 can have)
	br       *ByteRange      // not nil if leaf node
}

func (m *RadixMap) Set(key []byte, ref *ByteRange) error {
	currNode := m.root
	for i := 0; i < len(key); i++ {
		currKeyByte := key[i]

		// If children already exists for this char, update current node
		// and continue to next character
		if currNode.children[currKeyByte] != nil {
			currNode = currNode.children[currKeyByte]
			continue
		}

		// Allocate a new node for this char if not defined yet and continue
		newNode := &radixNode{}
		currNode.children[currKeyByte] = newNode
		currNode = newNode
	}
	// Set ref on last node
	currNode.br = ref
	return nil
}

func (m *RadixMap) Delete(key []byte) error {
	currNode := m.root
	for i := 0; i < len(key); i++ {
		currKeyByte := key[i]
		if currNode.children[currKeyByte] == nil {
			return nil // Stop and return (delete is no-op if key does not exist)
		}
		currNode = currNode.children[currKeyByte]
	}
	currNode.br = nil // Remove ref on last node (so it is not treated as leaf node anymore)
	return nil
}

func (m *RadixMap) Get(key []byte) *ByteRange {
	currNode := m.root
	for i := 0; i < len(key); i++ {
		currKeyByte := key[i]
		if currNode.children[currKeyByte] == nil {
			return nil // Return nil if unknown char (no children, no ref)
		}
		currNode = currNode.children[currKeyByte]
	}
	return currNode.br
}

func (m *RadixMap) Walk(cb WalkFunc) { m.root.walkRecurse([]byte{}, cb) }

func (n *radixNode) walkRecurse(k []byte, callback WalkFunc) {
	if n.br != nil {
		ok := callback(k, n.br)
		if !ok {
			return
		}
	}
	for i, child := range n.children {
		if child == nil {
			continue
		}
		child.walkRecurse(append(k, byte(i)), callback)
	}
}

func (m *RadixMap) IsLexOrdered() bool { return true }

func (m *RadixMap) NumKeys() int {
	count := 0
	m.Walk(func(_ []byte, _ *ByteRange) bool {
		count++
		return true
	})
	return count
}
