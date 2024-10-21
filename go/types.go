package merkletree

type Node struct {
    hash  string // hash of the node
    left  *Node	 // left child
    right *Node	 // right child
	parent *Node // parent node
}

type MerkleProof struct {
	hElement   string   // hash of element for which we want to prove inclusion
	siblings   []string // path of siblings from the element up to the root
	directions []bool   // signal if this sibling at this index of the directions array is on the left or right of its parent
}

type MerkleTree struct {
    root   *Node		// root of the tree
    leaves []*Node		// leaves of the tree
}

type AggregatedMerkleProof struct {
    start       uint64  // index of the first element in the proof
    end         uint64  // index of the last element in the proof
    siblings    []string  // path of siblings from the element up to the root
    directions  []bool    // signal if this sibling at this index of the directions array is on the left or right of its parent
}