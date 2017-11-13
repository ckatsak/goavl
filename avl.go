/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 2-Clause License. See the LICENSE file for details.
*/

// Package goavl implements the AVL data structure; heavily based on the
// description at http://www.geeksforgeeks.org/avl-tree-set-1-insertion/.
package goavl

// AVLNode represents a single node in the AVL tree.
type AVLNode struct {
	key    int
	left   *AVLNode
	right  *AVLNode
	height int
}

// NewAVLNode creates, initializes and returns the address of a new AVLNode.
func NewAVLNode(key int) *AVLNode {
	return &AVLNode{
		key:    key,
		height: 1, // initially inserted as leaf
	}
}

// Height returns the height of the subtree rooted with n.
func (n *AVLNode) Height() int {
	if n == nil {
		return 0
	}
	return n.height
}

// SubtreeRotateRight performs a right rotation of the subtree rooted with n, and
// returns a pointer to an AVLNode, which is the new root of the subtree.
func (n *AVLNode) SubtreeRotateRight() *AVLNode {
	m := n.left
	t2 := m.right

	// rotation
	m.right = n
	n.left = t2

	// update heights
	n.height = 1 + max(n.left.Height(), n.right.Height())
	m.height = 1 + max(m.left.Height(), m.right.Height())

	return m
}

// SubtreeRotateLeft performs a left rotation of the subtree rooted with n, and
// returns a pointer to an AVLNode, which is the new root of the subtree.
func (n *AVLNode) SubtreeRotateLeft() *AVLNode {
	m := n.right
	t2 := m.left

	// rotation
	m.left = n
	n.right = t2

	// update heights
	n.height = 1 + max(n.left.Height(), n.right.Height())
	m.height = 1 + max(m.left.Height(), m.right.Height())

	return m
}

// Balance returns the "balance factor" of AVLNode n.
func (n *AVLNode) Balance() int {
	if n == nil {
		return 0
	}
	return n.left.Height() - n.right.Height()
}

// SubtreeInsert inserts key as a new node in the AVL subtree rooted with n.
func (n *AVLNode) SubtreeInsert(key int) *AVLNode {
	// Step 1: Normal BST insertion
	if n == nil {
		return NewAVLNode(key)
	}

	if key < n.key {
		n.left = n.left.SubtreeInsert(key)
	} else if key > n.key {
		n.right = n.right.SubtreeInsert(key)
	} else {
		return n // no duplicate nodes
	}

	// Step 2: Update the height of this ancestor node
	n.height = 1 + max(n.left.Height(), n.right.Height())

	// Step 3: Check if the node is now unbalanced;
	//         if it is, handle the 4 possible cases.
	bal := n.Balance()
	switch {
	case bal > 1:
		switch {
		case key < n.left.key: // case left left
			return n.SubtreeRotateRight()
		case key > n.left.key: // case left right
			n.left = n.left.SubtreeRotateLeft()
			return n.SubtreeRotateRight()
		}
	case bal < -1:
		switch {
		case key > n.right.key: // case right right
			return n.SubtreeRotateLeft()
		case key < n.right.key: // case right left
			n.right = n.right.SubtreeRotateRight()
			return n.SubtreeRotateLeft()
		}
	}

	return n
}

// SubtreeDeleteNode deletes the node associated with key from the AVL subtree
// rooted with n.
func (n *AVLNode) SubtreeDeleteNode(key int) *AVLNode {
	// Step 1: Normal BST deletion
	if n == nil {
		return nil
	}

	if key < n.key {
		n.left = n.left.SubtreeDeleteNode(key)
	} else if key > n.key {
		n.right = n.right.SubtreeDeleteNode(key)
	} else { // this is the AVLNode to be deleted
		if n.left == nil || n.right == nil { // case of having < 2 children
			var tmp *AVLNode
			if n.left == nil {
				tmp = n.right
			} else {
				tmp = n.left
			}

			if tmp == nil { // case of no child at all
				tmp = n
				n = nil
			} else { // case of 1 child
				n = tmp
			}
		} else { // case of having exactly 2 children
			// get the inorder successor (smallest in the right subtree):
			tmp := n.right.SubtreeMin()
			// copy its data to us:
			n.key = tmp.key
			// delete the inorder successor:
			n.right = n.right.SubtreeDeleteNode(tmp.key)
		}
	}
	// If the tree had only 1 node, then return
	if n == nil {
		return n
	}

	// Step 2: Update the height of the node
	n.height = 1 + max(n.left.Height(), n.right.Height())

	// Step 3: Check if the node is now unbalanced;
	//         if it is, handle the 4 possible cases.
	bal := n.Balance()
	switch {
	case bal > 1:
		if n.left.Balance() >= 0 { // case left left
			return n.SubtreeRotateRight()
		} else { // case left right
			n.left = n.left.SubtreeRotateLeft()
			return n.SubtreeRotateRight()
		}
	case bal < -1:
		if n.right.Balance() <= 0 { // case right right
			return n.SubtreeRotateLeft()
		} else { // case right left
			n.right = n.right.SubtreeRotateRight()
			return n.SubtreeRotateLeft()
		}
	}

	return n
}

// SubtreeMin returns the AVLNode associated with the minimum key currently in
// the AVLTree.
func (n *AVLNode) SubtreeMin() *AVLNode {
	curr := n
	for curr.left != nil {
		curr = curr.left
	}
	return curr
}

// SubtreeMax returns the AVLNode associated with the maximum key currently in
// the AVLTree.
func (n *AVLNode) SubtreeMax() *AVLNode {
	curr := n
	for curr.right != nil {
		curr = curr.right
	}
	return curr
}

// AVLTree is a simple interface to interact with the AVL tree.
type AVLTree struct {
	root *AVLNode
}

// NewAVLTree creates a new empty AVL tree.
func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

// Insert inserts key into the AVL tree.
func (t *AVLTree) Insert(key int) {
	t.root = t.root.SubtreeInsert(key)
}

// Delete removes key from the AVL tree.
func (t *AVLTree) Delete(key int) {
	t.root = t.root.SubtreeDeleteNode(key)
}

// Min returns the minimum key currently in the AVL tree.
func (t *AVLTree) Min() int {
	return t.root.SubtreeMin().key
}

// Max returns the maximum key currently in the AVL tree.
func (t *AVLTree) Max() int {
	return t.root.SubtreeMax().key
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
