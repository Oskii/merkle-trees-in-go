package merkletree

import (
	"fmt"
	"testing"
)

//Happy path test for TestAggregatedProof
func TestAggregatedProof(t *testing.T) {
    elements := []string{"some", "test", "elements", "for", "testing", "aggregated", "proofs"}
    mt, err := NewMerkleTree(elements)
    if err != nil {
        t.Error(err)
    }

    proof, err := mt.GetAggregatedProof(1, 3)
    if err != nil {
        t.Error(err)
    }

    if !VerifyAggregatedProof(mt.getRoot(), proof, mt) {
        t.Error("invalid proof")
    }
	
}

//Test the entire merkle tree as an aggregated list of nodes to proove (edge case)
func TestEntireMerkleTreeUsingAggregatedProof(t *testing.T) {
    elements := []string{"some", "test", "elements", "for", "testing", "aggregated", "proofs"}
    mt, err := NewMerkleTree(elements)
    if err != nil {
        t.Error(err)
    }

	len:= uint64(len(elements))

    proof, err := mt.GetAggregatedProof(0, len-1)
    if err != nil {
        t.Error(err)
    }

    if !VerifyAggregatedProof(mt.getRoot(), proof, mt) {
        t.Error("invalid proof")
    }
	
}

//Test the smallest aggregated proof possible (edge case)
//I'm being opinionated here and saying that 2 elements is the smallest aggregated proof possible
//because 1 element is just a proof so you would use the other function for that 
func TestSmallestAggregatedProof(t *testing.T) {
	elements := []string{"one", "two"}
	mt, err := NewMerkleTree(elements)
	if err != nil {
		t.Error(err)
	}

	proof, err := mt.GetAggregatedProof(0, 1)
	if err != nil {
		t.Error(err)
	}

	if !VerifyAggregatedProof(mt.getRoot(), proof, mt) {
		t.Error("invalid proof")
	}
}

func TestInvalidlySmallSetOfElementsAggregatedProof(t *testing.T) {
	elements := []string{"one"}
	mt, err := NewMerkleTree(elements)
	if err != nil {
		t.Error(err)
	}

	_, err = mt.GetAggregatedProof(0, 0)
	if err == nil {
		t.Error("Expected an error for invalidly small set of elements, but got none")
	}
}
	

//Start index is greater than end index
func TestInvalidAggregatedProofIndexes(t *testing.T) {
    elements := []string{"some", "test", "elements", "for", "testing", "aggregated", "proofs"}
    mt, err := NewMerkleTree(elements)
    if err != nil {
        t.Error(err)
    }

    // Test case with invalid indexes (start >= end)
    _, err = mt.GetAggregatedProof(3, 1)
    if err == nil {
        t.Error("Expected an error for invalid indexes, but got none")
    }
}

//End index is greater than the number of elements
func TestInvalidAggregatedProofEndIndex(t *testing.T) {
    elements := []string{"some", "test", "elements", "for", "testing", "aggregated", "proofs"}
    mt, err := NewMerkleTree(elements)
    if err != nil {
        t.Error(err)
    }

    // Test case with an invalid endIndex greater than the number of elements
    _, err = mt.GetAggregatedProof(1, 10)
    if err == nil {
        t.Error("Expected an error for invalid endIndex, but got none")
    }
}

func TestRoot(t *testing.T) {
	elements := []string{"some", "test", "elements"}
	expectedRoot := hashNode(
		hashNode(hashLeaf("some"), hashLeaf("test")),
		hashNode(hashLeaf("elements"), hashLeaf("")),
	)
	testname := fmt.Sprintf("computes correct root")
	t.Run(testname, func(t *testing.T) {
		mt, err := NewMerkleTree(elements)
		if err != nil {
			t.Error(err)
		}
		if mt.getRoot() != expectedRoot {
			t.Errorf("got %s, want %s", mt.getRoot(), expectedRoot)
		}
	})
}

func TestProof(t *testing.T) {
	elements := []string{"some", "test", "elements"}
	mt, err := NewMerkleTree(elements)
	if err != nil {
		t.Error(err)
	}

	for i, elem := range elements {
		testname := fmt.Sprintf("valid proof for element: %d", i)
		t.Run(testname, func(t *testing.T) {
			proof, err := mt.GetProof(uint64(i))
			if err != nil {
				t.Error(err)
			}
			if !VerifyProof(mt.getRoot(), proof) {
				t.Error("invalid proof")
			}
			if hashLeaf(elem) != proof.hElement {
				t.Errorf("got %s, want %s", elem, proof.hElement)
			}
		})
	}
}

func TestProofOnSingleElement(t *testing.T) {
	elements := []string{"one"}
	mt, err := NewMerkleTree(elements)
	if err != nil {
		t.Error(err)
	}

	for i, elem := range elements {
		testname := fmt.Sprintf("valid proof for element: %d", i)
		t.Run(testname, func(t *testing.T) {
			proof, err := mt.GetProof(uint64(i))
			if err != nil {
				t.Error(err)
			}
			if !VerifyProof(mt.getRoot(), proof) {
				t.Error("invalid proof")
			}
			if hashLeaf(elem) != proof.hElement {
				t.Errorf("got %s, want %s", elem, proof.hElement)
			}
		})
	}
}

func TestInvalidProof(t *testing.T) {
	elements := []string{"some", "test", "elements"}
	mt, err := NewMerkleTree(elements)
	if err != nil {
		t.Error(err)
	}

	// Test case with an invalid index
	_, err = mt.GetProof(100)
	if err == nil {
		t.Error("Expected an error for invalid index, but got none")
	}
}

func TestUpdateElement(t *testing.T) {
    elements := []string{"alpha", "beta", "gamma"}
    mt, _ := NewMerkleTree(elements)

    // Update a leaf
    newElement := "delta"
    err := mt.UpdateElement(1, newElement) // Updating the second element
    if err != nil {
        t.Error("Failed to update element:", err)
    }

    // Verify the update
    updatedProof, _ := mt.GetProof(1)
    if updatedProof.hElement != hashLeaf(newElement) {
        t.Errorf("UpdateElement failed, got %s, want %s", updatedProof.hElement, hashLeaf(newElement))
    }
}

