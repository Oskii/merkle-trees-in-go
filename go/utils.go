
package merkletree

//Optimized power of 2 helper function that uses bitwise operations in log(n) time
func getNextPowerOfTwo(n int) int {
    if n <= 0 {
        return 1
    }
    n--
    n |= n >> 1
    n |= n >> 2
    n |= n >> 4
    n |= n >> 8
    n |= n >> 16
    return n + 1
}

func (t *MerkleTree) getRoot() string {
    if t.root == nil {
        return ""
    }
    return t.root.hash
}