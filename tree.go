package btree

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

		node.keys[i] = key
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

func (t *Tree) splitChild(parent *TreeNode, idx int) {
	// Retrieve the child that is about to be split
	fullChild := parent.children[idx]
	// Create a new node that will store keys greater than the median
	newChild := &TreeNode{
		isLeaf: fullChild.isLeaf,
	}

	// Step 1: Move the right half of the keys to the new child
	newChild.keys = append(newChild.keys, fullChild.keys[idx:]...)
	// Step 2: Keep only the left half of the keys in the original child
	// (keys before the median)
	medianKey := fullChild.keys[idx-1]
	fullChild.keys = fullChild.keys[:idx-1]

	// Step 3: If the node is not a leaf, we also need to move the right half of the children
	if !fullChild.isLeaf {
		newChild.children = append(newChild.children, fullChild.children[:idx]...)
		fullChild.children = fullChild.children[:idx]
	}
	// Step 4: Insert the new child into the parent’s children list
	// (insert right after the full child we just split)
	parent.children = append(parent.children[:idx+1],
		append([]*TreeNode{newChild}, parent.children[idx+1:]...)...)
	// Step 5: Insert the median key into the parent

	parent.keys = append(parent.keys[:idx+1],
		append([]int{medianKey}, parent.keys[idx+1:]...)...)
}
