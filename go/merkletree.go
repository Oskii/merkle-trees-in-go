package merkletree

import (
	"crypto/sha256"
	"fmt"
)

// Verifies a Merkle proof against a known root.
// Checks the direction to combine the hashes correctly
func VerifyProof(root string, proof MerkleProof) bool {
    currentHash := proof.hElement

    for i, siblingHash := range proof.siblings {
        if proof.directions[i] {
            currentHash = hashNode(siblingHash, currentHash)
        } else {
            currentHash = hashNode(currentHash, siblingHash)
        }
    }

    return currentHash == root
}

// Hash function to hash the leaves of the merkle tree
func hashLeaf(leaf string) string {
	h := sha256.New()
	h.Write([]byte(leaf))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Hash function to be used for the construction of the merkle tree nodes
func hashNode(a string, b string) string {
	h := sha256.New()
	h.Write([]byte(a))
	h.Write([]byte(b))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Generates a Merkle proof of the inclusion of the element at the given index.
// If the index is out of bounds, an error is returned.
//
// Example:
// proof for index 2 (marked with `h`), return the nodes marked `*` at each layer.
//
// tree:
// d0:                                   [ R ]
// d1:                [*]                                     [*]
// d2:      [*]                 [*]                 [ ]                 [ ]
// d3: [ ]       [ ]       [h]       [*]       [ ]       [ ]       [ ]       [ ]
//
// proof:
// hElement   = h
// siblings   = [d3-3, d2-0, d1-1]
// directions = [false, true, false]
func (t *MerkleTree) GetProof(index uint64) (MerkleProof, error) {
    if int(index) >= len(t.leaves) {
        return MerkleProof{}, fmt.Errorf("index out of bounds")
    }
    var siblings []string
    var directions []bool
    node := t.leaves[index]
    for node != t.root {
        parent := node.parent
        if parent.left == node {
            siblings = append(siblings, parent.right.hash)
            directions = append(directions, false)
        } else {
            siblings = append(siblings, parent.left.hash)
            directions = append(directions, true)
        }
        node = parent
    }

    return MerkleProof{
        hElement:   t.leaves[index].hash,
        siblings:   siblings,
        directions: directions,
    }, nil
}

// Creates a merkle tree from a list of elements.
// The tree should have the minimum height needed to contain all elements.
// Empty slots should be filled with an empty string.
func NewMerkleTree(elements []string) (*MerkleTree, error) {
    if len(elements) == 0 {
        return nil, fmt.Errorf("no elements to create a Merkle Tree")
    }

    leaves := make([]*Node, len(elements))

    for i, elem := range elements {
        leaves[i] = &Node{hash: hashLeaf(elem)}
    }

    requiredLeaves := getNextPowerOfTwo(len(elements))
    for len(leaves) < requiredLeaves {
        leaves = append(leaves, &Node{hash: hashLeaf("")})
    }

	// make some memory on the heap for the nodes
    allNodes := make([]*Node, 0, len(leaves)*2) 

	// put the leaves and nodes in the slots
    allNodes = append(allNodes, leaves...)

    for len(leaves) > 1 {
        
		var nextLevel []*Node
        for i := 0; i < len(leaves); i += 2 {
            left := leaves[i]
            right := leaves[i+1]
            parent := &Node{
                left:  left,
                right: right,
                hash:  hashNode(left.hash, right.hash),
            }
            left.parent = parent
            right.parent = parent
            nextLevel = append(nextLevel, parent)
        }
        allNodes = append(allNodes, nextLevel...)
        leaves = nextLevel
    }

    return &MerkleTree{root: leaves[0], leaves: allNodes[:requiredLeaves]}, nil
}

func (n *Node) getSibling() *Node {
    if n.parent == nil {
        return nil // Root node has no sibling
    }
    if n == n.parent.left {
        return n.parent.right
    }
    return n.parent.left
}

func (n *Node) isLeftChild() bool {
    if n.parent == nil {
        return false // Root node is not a left child
    }
    return n == n.parent.left
}

// ** BONUS (optional - easy) **
// Updates the Merkle tree (from leaf to root) to include the new element at index.
// For simplicity, the index must be within the bounds of the original vector size.
// If it is not, return an error.
func (t *MerkleTree) UpdateElement(index uint64, element string) error {
    if int(index) >= len(t.leaves) {
        return fmt.Errorf("index out of bounds")
    }

    // Update the leaf node
    t.leaves[index].hash = hashLeaf(element)

    // Update all ancestor hashes up to the root
    node := t.leaves[index]
    for node != t.root {
        parent := node.parent
        if parent.left == node {
            parent.hash = hashNode(node.hash, parent.right.hash)
        } else {
            parent.hash = hashNode(parent.left.hash, node.hash)
        }
        node = parent
    }

    return nil
}

// ** BONUS (optional - hard) **
// Generates a Merkle proof of the inclusion of contiguous elements,
// starting at startIndex (inclusive) and ending at endIndex (exclusive).
// If the indexes are out of bounds or startIndex >= endIndex, an error is returned.
//
// Note: modify the method signature to return your proof type.
// Implement a separate VerifyAggregatedProof for this type.
//
// The aggregated proof size should generally be smaller than
// that of the naive approach (calling GetProof for every index).
func (t *MerkleTree) GetAggregatedProof(startIndex, endIndex uint64) (*AggregatedMerkleProof, error) {
    if startIndex >= endIndex || int(endIndex) > len(t.leaves) {
		return nil, fmt.Errorf("index out of bounds")
	}

	if endIndex - startIndex < 1 {
		return nil, fmt.Errorf("an aggregated proof must contain at least two elements")
	}

	// Collect path from leaf to root
	var siblings []string
	var directions []bool
	node := t.leaves[startIndex]
	for node != t.root {
		parent := node.parent
		if parent.left == node {
			siblings = append(siblings, parent.right.hash)
			directions = append(directions, false)
		} else {
			siblings = append(siblings, parent.left.hash)
			directions = append(directions, true)
		}
		node = parent
	}

	return &AggregatedMerkleProof{
		start:      startIndex,
		end:        endIndex,
		siblings:   siblings,
		directions: directions,
	}, nil
}

// Verifies an aggregated Merkle proof against a known root.
func VerifyAggregatedProof(root string, aggProof *AggregatedMerkleProof, tree *MerkleTree) bool {
    if aggProof.start >= aggProof.end || aggProof.end > uint64(len(tree.leaves)) {
        return false
    }

    currentHash := tree.leaves[aggProof.start].hash

    for i := 0; i < len(aggProof.siblings); i++ {
        siblingHash := aggProof.siblings[i]
        // Check the direction to combine the hashes correctly
        if aggProof.directions[i] {
            // Sibling is on the left
            currentHash = hashNode(siblingHash, currentHash)
        } else {
            // Sibling is on the right
            currentHash = hashNode(currentHash, siblingHash)
        }
    }

    // Compare the computed hash with the root hash
	
    return currentHash == root
}