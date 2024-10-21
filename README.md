## Summary and Asymptotic Complexity 

My approach was to prioritize efficiency, and minimize computational complexity. I avoided any greedy algorithms for the computation of Merkle tree generations, proof generations, or proof verifications. The approach I adopted achieved the desired results with a time complexity of O(log n), making it well-suited for handling large Merkle trees.

### Building the Merkle tree
The NewMerkleTree function has a time complexity of O(n), where n is the number of elements added to the tree.

### GetProof 
The GetProof function has a time complexity of O(log n), where n is the number of elements. More specifically it has to interact with O(2logn) nodes, because at every step it has to hash the two nodes together to get to the next level. The height (number of levels) of the Merkle tree is logarithmic in its number of elements at every level, so the traversal complexity is O(log n).

### UpdateElement 
The UpdateElement function has a time complexity of O(log n), where n is the number of elements. Updating an element requires recalculating hashes from the updated leaf node to the root. This involves traversing the tree from the leaf to the root, ending up with a log time complexity, due to the fact that the number of nodes half in size at every next level.

### GetAggregatedProof 

The GetAggregatedProof function has a time complexity of O(log n), where n is the number of elements. The function traverses the tree from the leaf node at startIndex to the root. In the worst case, this traversal involves going up the height of the tree, which is logarithmic in the number of elements (O(log n)).

At each step of the traversal, the function appends the sibling hash and direction to their respective slices. Appending to a slice is an O(1) operation on average because Go slices have amortized constant-time append operations
### VerifyAggregatedProof

The VerifyAggregatedProof function has a time complexity of O(log n), where n is the number of elements.

## Potential Further Work
One cool thing to implement would be multi-range proofs for checking numerous chunks of the merkle tree simultaneously like they did here
https://eprint.iacr.org/2021/038.pdf

## Rationale
A core cryptographic primitive in Ethereum is the vector commitment, which commits to an ordered vector of elements and can be opened (i.e. proved that an element exists in the vector committed to) at any position. A Merkle tree is a vector commitment, represented by a simple binary tree structure where the root node is the hash of its two child hashes, concatenated. Given a collision-resistant hash function and any vector of elements represented by the leaves, this leads to a unique root. Proof of existence of an element can be represented by a hash path, a list of pairwise child hashes at each layer, from leaf to root.

## Building and running
```
# clone the repo
git clone git@github.com:SpecularL2/interview-test.git

# run the test for go
go test .

# run the tests for rust
cargo test
```
