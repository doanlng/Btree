package btree

import (
	"testing"
)

func TestPrint(t *testing.T) {
	// Create a B-tree with degree 3
	tree := NewBTree(3)

	// Test printing empty tree
	t.Log("Empty tree:")
	tree.Print()

	// Insert some values
	values := []int{10, 20, 5, 6, 12, 30, 7, 17}

	for _, val := range values {
		tree.Insert(val)
	}

	// Print the tree structure
	t.Log("Tree after insertions:")
	tree.Print()
}

func TestTreeOperations(t *testing.T) {
	// Create a B-tree with degree 3
	tree := NewBTree(2)

	// Insert values and print at each step
	values := []int{10, 20, 5, 6, 12, 30, 40, 15, 26, 18, 500, 19}

	for i, val := range values {
		tree.Insert(val)
		t.Logf("After inserting %d (step %d):", val, i+1)
		tree.Print()
		t.Log("---")
	}
}
