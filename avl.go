/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 2-Clause License. See the LICENSE file for details.
*/

// Package goavl provides a generic implementation of the AVL Tree data
// structure.
//
// Based on the description found at GeeksforGeeks.
package goavl

import "fmt"

// Item is the interface required to be satisfied by any type to be able to
// populate the AVL tree.
type Item interface {
	Equal(to Item) bool
	Less(than Item) bool
}

// treeNode represents a single node in the AVL tree.
type treeNode struct {
	key         Item
	left, right *treeNode
	h           int
}

// newNode allocates, initializes and returns the address of a new treeNode.
func newNode(key Item) *treeNode {
	return &treeNode{
		key: key,
		h:   1, // initially inserted as a leaf
	}
}

// height returns the height of the subtree rooted with n.
func (n *treeNode) height() int {
	if n == nil {
		return 0
	}
	return n.h
}

// subtreeRotateRight performs a right rotation of the subtree rooted with n, and
// returns a pointer to a treeNode, which is the new root of the subtree.
func (n *treeNode) subtreeRotateRight() *treeNode {
	m := n.left
	t2 := m.right

	// rotation
	m.right = n
	n.left = t2

	// update heights
	n.h = 1 + max(n.left.height(), n.right.height())
	m.h = 1 + max(m.left.height(), m.right.height())

	return m
}

// subtreeRotateLeft performs a left rotation of the subtree rooted with n, and
// returns a pointer to a treeNode, which is the new root of the subtree.
func (n *treeNode) subtreeRotateLeft() *treeNode {
	m := n.right
	t2 := m.left

	// rotation
	m.left = n
	n.right = t2

	// update heights
	n.h = 1 + max(n.left.height(), n.right.height())
	m.h = 1 + max(m.left.height(), m.right.height())

	return m
}

// balanceFactor returns the "balance factor" of treeNode n.
func (n *treeNode) balanceFactor() int {
	if n == nil {
		// NOTE: This is probably unreachable, but anyway.
		return 0
	}
	return n.left.height() - n.right.height()
}

// subtreeInsertNode inserts key as a new node in the AVL subtree rooted with n.
func (n *treeNode) subtreeInsertNode(key Item) (*treeNode, error) {
	var err error

	// Step 1: Normal BST insertion
	if n == nil {
		return newNode(key), nil
	}

	if key.Less(n.key) {
		n.left, err = n.left.subtreeInsertNode(key)
	} else if key.Equal(n.key) {
		return n, fmt.Errorf("Key already in the tree: %v", key) // no duplicate nodes
	} else { // if key.Greater(n.key) {
		n.right, err = n.right.subtreeInsertNode(key)
	}

	// Step 2: Update the height of this ancestor node
	n.h = 1 + max(n.left.height(), n.right.height())

	// Step 3: Check if the node is now unbalanced;
	//         if it is, handle the 4 possible cases.
	bal := n.balanceFactor()
	switch {
	case bal > 1:
		if key.Less(n.left.key) { // case left left
			return n.subtreeRotateRight(), err
		}
		// else if key.Greater(n.left.key): // case left right
		n.left = n.left.subtreeRotateLeft()
		return n.subtreeRotateRight(), err
	case bal < -1:
		if key.Less(n.right.key) { // case right left
			n.right = n.right.subtreeRotateRight()
			return n.subtreeRotateLeft(), err
		}
		// else if key.Greater(n.right.key): // case right right
		return n.subtreeRotateLeft(), err
	}

	return n, err
}

// subtreeDeleteNode deletes the node associated with key from the AVL subtree
// rooted with n.
func (n *treeNode) subtreeDeleteNode(key Item) (*treeNode, error) {
	var err error

	// Step 1: Normal BST deletion
	if n == nil {
		return nil, fmt.Errorf("Key not found in the tree: %v", key)
	}

	if key.Less(n.key) {
		n.left, err = n.left.subtreeDeleteNode(key)
	} else if key.Equal(n.key) { // this is the treeNode to be deleted
		if n.left == nil || n.right == nil { // case of having < 2 children
			var tmp *treeNode
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
			tmp := n.right.subtreeMin()
			// copy its data to us:
			n.key = tmp.key
			// delete the inorder successor:
			n.right, err = n.right.subtreeDeleteNode(tmp.key)
		}
	} else { // if key.Greater(n.key) {
		n.right, err = n.right.subtreeDeleteNode(key)
	}
	// If the tree had only 1 node, then return
	if n == nil {
		return n, err
	}

	// Step 2: Update the height of the node
	n.h = 1 + max(n.left.height(), n.right.height())

	// Step 3: Check if the node is now unbalanced;
	//         if it is, handle the 4 possible cases.
	bal := n.balanceFactor()
	switch {
	case bal > 1:
		if n.left.balanceFactor() >= 0 { // case left left
			return n.subtreeRotateRight(), err
		}
		// else if n.left.balanceFactor() < 0: // case left right
		n.left = n.left.subtreeRotateLeft()
		return n.subtreeRotateRight(), err
	case bal < -1:
		if n.right.balanceFactor() <= 0 { // case right right
			return n.subtreeRotateLeft(), err
		}
		// else if n.right.balanceFactor() > 0: // case right left
		n.right = n.right.subtreeRotateRight()
		return n.subtreeRotateLeft(), err
	}

	return n, err
}

// subtreeMin returns the treeNode associated with the minimum key currently in
// the AVL tree.
func (n *treeNode) subtreeMin() *treeNode {
	curr := n
	for curr.left != nil {
		curr = curr.left
	}
	return curr
}

// subtreeMax returns the treeNode associated with the maximum key currently in
// the AVL tree.
func (n *treeNode) subtreeMax() *treeNode {
	curr := n
	for curr.right != nil {
		curr = curr.right
	}
	return curr
}

// subtreeInOrder returns a slice of all Items currently in the AVL sub-tree
// rooted by n, by performing an in-order traversal of its nodes.
func (n *treeNode) subtreeInOrder() []Item {
	if n == nil {
		return nil
	}
	ret := []Item{}
	ret = append(ret, n.left.subtreeInOrder()...)
	ret = append(ret, n.key)
	ret = append(ret, n.right.subtreeInOrder()...)
	return ret
}

// subtreePreOrder returns a slice of all Items currently in the AVL sub-tree
// rooted by n, by performing a pre-order traversal of its nodes.
func (n *treeNode) subtreePreOrder() []Item {
	if n == nil {
		return nil
	}
	ret := []Item{n.key}
	ret = append(ret, n.left.subtreePreOrder()...)
	ret = append(ret, n.right.subtreePreOrder()...)
	return ret
}

// Tree is the exported struct for interacting with the AVL tree.
type Tree struct {
	root *treeNode
	size int
}

// NewTree creates a new empty AVL tree.
func NewTree() *Tree {
	return &Tree{}
}

// Size returns the current number of keys in the AVL tree.
func (t *Tree) Size() int {
	return t.size
}

// Insert inserts a key into the AVL tree and returns an error value, which is
// non-nil if the key already exists in the tree (i.e. duplicate keys are not
// supported).
func (t *Tree) Insert(key Item) (err error) {
	if t.root, err = t.root.subtreeInsertNode(key); err == nil {
		t.size++
	}
	return
}

// Delete removes a key from the AVL tree and returns an error value, which is
// non-nil if the key doesn't exist in the tree.
func (t *Tree) Delete(key Item) (err error) {
	if t.root, err = t.root.subtreeDeleteNode(key); err == nil {
		t.size--
	}
	return
}

// Min returns the minimum key in the AVL tree and an error value. If the tree
// is empty, the error value is non-nil and the result should not be trusted.
func (t *Tree) Min() (Item, error) {
	if t.root == nil {
		return nil, fmt.Errorf("Empty tree")
	}
	return t.root.subtreeMin().key, nil
}

// Max returns the maximum key in the AVL tree and an error value. If the tree
// is empty, the error value is non-nil and the result should not be trusted.
func (t *Tree) Max() (Item, error) {
	if t.root == nil {
		return nil, fmt.Errorf("Empty tree")
	}
	return t.root.subtreeMax().key, nil
}

// Height returns the current height of the AVL tree.
func (t *Tree) Height() int {
	return t.root.height()
}

// InOrder returns a slice of all Items currently in the AVL Tree by performing
// an in-order traversal of its nodes.
func (t *Tree) InOrder() []Item {
	return t.root.subtreeInOrder()
}

// PreOrder returns a slice of all Items currently in the AVL Tree by
// performing a pre-order traversal of its nodes.
func (t *Tree) PreOrder() []Item {
	return t.root.subtreePreOrder()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
