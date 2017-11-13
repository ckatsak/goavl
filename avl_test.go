/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 2-Clause License. See the LICENSE file for details.
*/

package goavl

import (
	"math/rand"
	"sort"
	"testing"
)

func preOrder(t *testing.T, n *treeNode) []int {
	t.Helper()
	if n == nil { // case n is leaf
		return nil
	}
	results := []int{n.key}
	results = append(results, preOrder(t, n.left)...)
	results = append(results, preOrder(t, n.right)...)
	return results
}

func TestSimplePreorder(t *testing.T) {
	tree := NewTree()

	tree.Insert(9)
	tree.Insert(5)
	tree.Insert(10)
	tree.Insert(0)
	tree.Insert(6)
	tree.Insert(11)
	tree.Insert(-1)
	tree.Insert(1)
	tree.Insert(2)

	t.Logf("Preorder before deletion of 10: %v\n", preOrder(t, tree.root))
	tree.Delete(10)
	t.Logf("Preorder after deletion of 10: %v\n", preOrder(t, tree.root))
}

func TestInsertExisting(t *testing.T) {
	tree := NewTree()

	t.Logf("Preorder initial: %v\n", preOrder(t, tree.root))
	tree.Insert(42)
	t.Logf("Preorder after inserting 42: %v\n", preOrder(t, tree.root))
	tree.Insert(42)
	t.Logf("Preorder after re-inserting 42: %v\n", preOrder(t, tree.root))
	tree.Insert(42)
	t.Logf("Preorder after re-inserting 42: %v\n", preOrder(t, tree.root))
}

func TestDeleteNonExisting(t *testing.T) {
	tree := NewTree()

	t.Logf("Preorder initial: %v\n", preOrder(t, tree.root))
	tree.Delete(42)
	t.Logf("Preorder after deleting 42: %v\n", preOrder(t, tree.root))
	tree.Insert(24)
	t.Logf("Preorder after inserting 24: %v\n", preOrder(t, tree.root))
	tree.Delete(42)
	t.Logf("Preorder after re-deleting 42: %v\n", preOrder(t, tree.root))
}

func inOrder(t *testing.T, n *treeNode) []int {
	t.Helper()
	if n == nil {
		return nil
	}
	results := []int{}
	results = append(results, inOrder(t, n.left)...)
	results = append(results, n.key)
	results = append(results, inOrder(t, n.right)...)
	return results
}

func verifyTraversal(t *testing.T, traversal, sortedRands []int) {
	t.Helper()
	for i := 0; i < len(sortedRands); i++ {
		if traversal[i] != sortedRands[i] {
			if sortedRands[i] == sortedRands[i-1] {
				t.Error("Unlucky: duplicate random number spotted.")
				t.FailNow()
			}
			t.Errorf("traversal[%d] is %d, should be %d\n", i, traversal[i], sortedRands[i])
		}
	}
}

func populateTreeAndSlice(t *testing.T, tree *Tree, size uint) []int {
	t.Helper()
	rands := []int{}
	for i := uint(0); i < size; i++ {
		r := rand.Int()
		tree.Insert(r)
		rands = append(rands, r)
	}
	return rands
}

func TestInsertInOrder(t *testing.T) {
	tree := NewTree()

	// Create a slice of random integers
	rands := populateTreeAndSlice(t, tree, 1<<21)

	// Create the inorder traversal of the tree
	traversal := inOrder(t, tree.root)
	if !sort.IntsAreSorted(traversal) {
		t.Errorf("In-order traversal resulted in unsorted set.")
	}

	// Sort the slice of random integers and compare it against the inorder traversal
	sortedRands := append([]int{}, rands...)
	sort.Ints(sortedRands)

	verifyTraversal(t, traversal, sortedRands)
}

func TestDeleteInOrder(t *testing.T) {
	tree := NewTree()

	// Create a slice of random integers
	rands := populateTreeAndSlice(t, tree, 1<<21)

	indicesToRemove := []int{}
	for i := 0; i < 1<<11; i++ {
		r := rand.Intn((1 << 21) - i)
		indicesToRemove = append(indicesToRemove, r)

		tree.Delete(rands[r])
		rands[r] = rands[len(rands)-1]
		rands = rands[:len(rands)-1]
	}

	// Sort the slice of random integers and compare it against the inorder traversal
	sortedRands := append([]int{}, rands...)
	sort.Ints(sortedRands)

	// Create the inorder traversal of the tree
	traversal := inOrder(t, tree.root)
	if !sort.IntsAreSorted(traversal) {
		t.Errorf("In-order traversal resulted in unsorted set.\n")
	}

	verifyTraversal(t, traversal, sortedRands)
}

func TestMinDelete(t *testing.T) {
	tree := NewTree()

	// Create a slice of random integers
	size := uint(1 << 21)
	rands := populateTreeAndSlice(t, tree, size)

	sort.Ints(rands)
	for i := uint(0); i < size; i++ {
		listMin := rands[0]
		treeMin := tree.Min()
		if listMin != treeMin {
			t.Errorf("listMin = %d, treeMin = %d\n", listMin, treeMin)
		}
		rands = rands[1:]
		tree.Delete(treeMin)
	}
}

func TestMaxDelete(t *testing.T) {
	tree := NewTree()

	// Create a slice of random integers
	size := uint(1 << 21)
	rands := populateTreeAndSlice(t, tree, size)

	sort.Ints(rands)
	for i := uint(0); i < size; i++ {
		listMax := rands[len(rands)-1]
		treeMax := tree.Max()
		if listMax != treeMax {
			t.Errorf("listMax = %d, treeMax = %d\n", listMax, treeMax)
		}
		rands = rands[:len(rands)-1]
		tree.Delete(treeMax)
	}
}

func TestHeight(t *testing.T) {
	tree := NewTree()

	t.Logf("Height for no keys: %d\n\n", tree.Height())

	tree.Insert(0)
	t.Logf("Height for 1 key: %d\n\n", tree.Height())

	// for exp=29, more than 14G of memory are required
	// for exp=28, ~ 8.5G - 10.5G of memory are required (I didn't notice the exact amount)
	for exp := uint(1); exp < 25; exp++ {
		// Insert new keys from range [2**(e-1), (2**e)-2] --> 2**(e-1)-2 new keys.
		for i := 1 << (exp - 1); i < (1<<exp)-1; i++ {
			tree.Insert(i)
		}
		t.Logf("Height for %d keys: %d\n", (1<<exp)-1, tree.Height())
		//t.Logf("\tPreorder: %v\n", preOrder(t, tree.root))
		if tree.Height() != int(exp) {
			t.Errorf("\tHeight for %d keys is expected to be %d.\n", (1<<exp)-1, exp)
		}

		// Insert 2**e -th key, which should increase tree's height by 1.
		tree.Insert((1 << exp) - 1)
		t.Logf("Height for %d keys: %d\n", 1<<exp, tree.Height())
		//t.Logf("\tPreorder: %v\n", preOrder(t, tree.root))
		if tree.Height() != int(exp+1) {
			t.Errorf("\tHeight for %d keys is expected to be %d.\n", 1<<exp, exp+1)
		}

		// Insert a 2**e+1 -th key, which shouldn't increase tree's height, and then remove it again.
		tree.Insert(-42)
		t.Logf("Height for %d keys: %d\n", (1<<exp)+1, tree.Height())
		//t.Logf("\tPreorder: %v\n", preOrder(t, tree.root))
		if tree.Height() != int(exp+1) {
			t.Errorf("\tHeight for %d keys is expected to be %d.\n", (1<<exp)+1, exp+1)
		}
		tree.Delete(-42)
		t.Logf("\n")
	}
}
