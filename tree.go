package btree

import "fmt"

type Tree struct {
	root   *TreeNode
	degree int
}

type TreeNode struct {
	keys     []int
	children []*TreeNode
	isLeaf   bool
}

func NewBTree(degree int) *Tree {
	if degree < 2 { // every node except root should have between degree - 1 and 2 * degree - 1 keys
		panic("minimum degree t must be >= 2")
	}
	root := TreeNode{
		keys:     []int{},
		children: []*TreeNode{},
		isLeaf:   true,
	}
	return &Tree{
		root:   &root,
		degree: degree,
	}
}

func (t *Tree) SearchTree(val int) (*TreeNode, int) {
	node, idx, found := t.root.search(val)
	if !found {
		return nil, -1 // missing element
	}
	return node, idx
}

func (tn *TreeNode) search(val int) (*TreeNode, int, bool) {
	i := 0
	for i < len(tn.keys) && tn.keys[i] < val {
		i += 1
	}

	if i < len(tn.keys) && tn.keys[i] == val {
		return tn, i, true
	}

	if tn.isLeaf {
		return nil, -1, false
	}

	return tn.children[i].search(val)

}

func (t *Tree) Insert(val int) bool {
	node := t.root
	if len(node.keys) == 2*t.degree-1 { // the root is full
		s := &TreeNode{
			keys:     []int{},
			children: []*TreeNode{node}, // put the root as the child of the promoted node
			isLeaf:   false,
		} // create the promoted node
		t.root = s // set the just promoted node as the new root

		// Split the old root into two children.
		t.splitChild(s, 0)

		// The root should have non-full children now and we can insert as normal
		t.insertNonFull(s, val)
	} else {
		t.insertNonFull(node, val)
	}

	return true
}

func (t *Tree) insertNonFull(node *TreeNode, key int) {
	i := len(node.keys) - 1
	if node.isLeaf { // case 1: the node is a leaf
		// find where to insert the key in its keys
		node.keys = append(node.keys, 0)

		for i >= 0 && node.keys[i] > key { // find the smallest place to insert the key
			node.keys[i+1] = node.keys[i] // slide the key over one to make room
			i -= 1
		}

		node.keys[i+1] = key
		return
	}

	// case 2: the node is not a leaf and requires us to recursively descend
	for i >= 0 && node.keys[i] > key {
		i -= 1
	}

	i += 1 // the child to descend into will be at i + 1
	if len(node.children[i].keys) == 2*t.degree-1 {
		//split before we descend
		t.splitChild(node, i)

		if key > node.keys[i] {
			i += 1
		}
	}

	t.insertNonFull(node.children[i], key)
}

/*
takes a parent tree node and the index of its full child
*/
func (t *Tree) splitChild(parent *TreeNode, idx int) {
	// Retrieve the child that is about to be split
	fullChild := parent.children[idx]
	// Create a new node that will store keys greater than the median
	newChild := &TreeNode{
		isLeaf: fullChild.isLeaf,
	}

	// The median is at position t.degree - 1
	medianIdx := t.degree - 1

	// Step 1: Move the right half of the full childs keys into the new node
	newChild.keys = append(newChild.keys, fullChild.keys[medianIdx+1:]...)
	// Step 2: Keep only the left half of the keys in the original child
	// (keys before the median)
	medianKey := fullChild.keys[medianIdx]
	fullChild.keys = fullChild.keys[:medianIdx]

	// Step 3: If the node is not a leaf, we also need to move the right half of the children
	if !fullChild.isLeaf {
		newChild.children = append(newChild.children, fullChild.children[medianIdx+1:]...)
		fullChild.children = fullChild.children[:medianIdx+1]
	}
	// Step 4: Insert the new child into the parent's children list
	// (insert right after the full child we just split)
	parent.children = append(parent.children[:idx+1],
		append([]*TreeNode{newChild}, parent.children[idx+1:]...)...)
	// Step 5: Insert the median key into the parent

	parent.keys = append(parent.keys[:idx],
		append([]int{medianKey}, parent.keys[idx:]...)...)
}

// Print displays the B-tree structure in a readable format
func (t *Tree) Print() {
	if t.root == nil {
		fmt.Println("Empty tree")
		return
	}
	fmt.Printf("B-tree (degree %d):\n", t.degree)
	t.printNode(t.root, "", true)
}

// printNode recursively prints the tree structure with proper indentation
func (t *Tree) printNode(node *TreeNode, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Print current node
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	nodeType := "Internal"
	if node.isLeaf {
		nodeType = "Leaf"
	}

	fmt.Printf("%s%s%s: %v\n", prefix, connector, nodeType, node.keys)

	// Print children if this is not a leaf
	if !node.isLeaf {
		childPrefix := prefix
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}

		for i, child := range node.children {
			isLastChild := i == len(node.children)-1
			t.printNode(child, childPrefix, isLastChild)
		}
	}
}
