/*
Copyright (C) 2017, Christos Katsakioris
All rights reserved.

This software may be modified and distributed under the terms
of the BSD 2-Clause License. See the LICENSE file for details.
*/

/*
Plain execution of the tests:
	$ go test -cover
or:
	$ go test -v -cover
for a verbose output.

Create a coverprofile and produce a html file:
	$ go test -coverprofile=cov.out
	$ go tool cover -html=cov.out -o cov.htm
*/

package goavl

import (
	"math/rand"
	"sort"
	"testing"
)

// AUXILIARY TYPES

type Integer int

func (i Integer) Equal(j Item) bool {
	return i == j.(Integer)
}
func (i Integer) Less(j Item) bool {
	return i < j.(Integer)
}

// Compile time check that Integer satisfies the Item interface.
var _ Item = Integer(42)

// AUXILIARY FUNCTIONS

func preOrder(t *testing.T, n *treeNode) []Integer {
	t.Helper()
	if n == nil { // case n is leaf
		return nil
	}
	results := []Integer{n.key.(Integer)}
	results = append(results, preOrder(t, n.left)...)
	results = append(results, preOrder(t, n.right)...)
	return results
}

func inOrder(t *testing.T, n *treeNode) []Integer {
	t.Helper()
	if n == nil {
		return nil
	}
	results := []Integer{}
	results = append(results, inOrder(t, n.left)...)
	results = append(results, n.key.(Integer))
	results = append(results, inOrder(t, n.right)...)
	return results
}

func verifyTraversal(t *testing.T, traversal []Integer, sortedRands []int) {
	t.Helper()
	for i := 0; i < len(sortedRands); i++ {
		if traversal[i] != Integer(sortedRands[i]) {
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
		if err := tree.Insert(Integer(r)); err != nil {
			t.Errorf("\t%v\n", err)
		}
		rands = append(rands, r)
	}
	return rands
}

// TEST FUNCTIONS

func TestSimplePreorder(t *testing.T) {
	tree := NewTree()

	for _, key := range []Integer{9, 5, 10, 0, 6, 11, -1, 1, 2} {
		if err := tree.Insert(key); err != nil {
			t.Errorf("\t%v\n", err)
		}
	}
	t.Logf("Preorder before deletion of 10: %v\n", preOrder(t, tree.root))

	if err := tree.Delete(Integer(10)); err != nil {
		t.Errorf("\t%v\n", err)
	}
	t.Logf("Preorder after deletion of 10: %v\n", preOrder(t, tree.root))
}

func TestInsertExisting(t *testing.T) {
	tree := NewTree()
	var err error

	t.Logf("Preorder initial: %v\n", preOrder(t, tree.root))

	err = tree.Insert(Integer(42))
	t.Logf("Preorder after inserting 42: %v\n", preOrder(t, tree.root))
	if err != nil {
		t.Errorf("\t%v\n", err)
	} else {
		t.Logf("\tNo error value returned, as expected.\n")
	}

	err = tree.Insert(Integer(42))
	t.Logf("Preorder after re-inserting 42: %v\n", preOrder(t, tree.root))
	if err == nil {
		t.Errorf("\tExpected an error!\n")
	} else {
		t.Logf("\tError value returned, as expected: \"%v\"\n", err)
	}

	err = tree.Insert(Integer(42))
	t.Logf("Preorder after re-inserting 42: %v\n", preOrder(t, tree.root))
	if err == nil {
		t.Errorf("\tExpected an error!\n")
	} else {
		t.Logf("\tError value returned, as expected: \"%v\"\n", err)
	}
}

func TestDeleteNonExisting(t *testing.T) {
	tree := NewTree()
	var err error

	t.Logf("Preorder initial: %v\n", preOrder(t, tree.root))

	err = tree.Delete(Integer(42))
	t.Logf("Preorder after deleting 42: %v\n", preOrder(t, tree.root))
	if err == nil {
		t.Errorf("\tExpected an error!\n")
	} else {
		t.Logf("\tError value returned, as expected: \"%v\"\n", err)
	}

	if err = tree.Insert(Integer(24)); err != nil {
		t.Errorf("\t%v\n", err)
	}
	t.Logf("Preorder after inserting 24: %v\n", preOrder(t, tree.root))

	err = tree.Delete(Integer(42))
	t.Logf("Preorder after re-deleting 42: %v\n", preOrder(t, tree.root))
	if err == nil {
		t.Errorf("\tExpected an error!\n")
	} else {
		t.Logf("\tError value returned, as expected: \"%v\"\n", err)
	}

	if err = tree.Insert(Integer(42)); err != nil {
		t.Errorf("\t%v\n", err)
	}
	t.Logf("Preorder after inserting 42: %v\n", preOrder(t, tree.root))

	err = tree.Delete(Integer(42))
	t.Logf("Preorder after re-deleting 42: %v\n", preOrder(t, tree.root))
	if err != nil {
		t.Errorf("\t%v\n", err)
	} else {
		t.Logf("\tNo error value returned, as expected\n")
	}
}

func TestInsertInOrder(t *testing.T) {
	tree := NewTree()

	// Create a slice of random integers
	rands := populateTreeAndSlice(t, tree, 1<<20)

	// Create the inorder traversal of the tree
	traversal := inOrder(t, tree.root)
	//if !sort.IntsAreSorted(traversal) {
	//	t.Errorf("In-order traversal resulted in unsorted set.")
	//}

	// Sort the slice of random integers and compare it against the inorder traversal
	sortedRands := append([]int{}, rands...)
	sort.Ints(sortedRands)

	verifyTraversal(t, traversal, sortedRands)
}

func TestDeleteInOrder(t *testing.T) {
	tree := NewTree()

	// Create a slice of random integers
	rands := populateTreeAndSlice(t, tree, 1<<20)

	indicesToRemove := []int{}
	for i := 0; i < 1<<11; i++ {
		r := rand.Intn((1 << 20) - i)
		indicesToRemove = append(indicesToRemove, r)

		if err := tree.Delete(Integer(rands[r])); err != nil {
			t.Errorf("\t%v\n", err)
		}
		rands[r] = rands[len(rands)-1]
		rands = rands[:len(rands)-1]
	}

	// Sort the slice of random integers and compare it against the inorder traversal
	sortedRands := append([]int{}, rands...)
	sort.Ints(sortedRands)

	// Create the inorder traversal of the tree
	traversal := inOrder(t, tree.root)
	//if !sort.IntsAreSorted(traversal) {
	//	t.Errorf("In-order traversal resulted in unsorted set.\n")
	//}

	verifyTraversal(t, traversal, sortedRands)
}

func TestEmptyMinMax(t *testing.T) {
	tree := NewTree()
	if _, err := tree.Min(); err != nil {
		t.Logf("\tError value returned, as expected: \"%v\"\n", err)
	} else {
		t.Errorf("\tExpected an error!\n")
	}
	if _, err := tree.Max(); err != nil {
		t.Logf("\tError value returned, as expected: \"%v\"\n", err)
	} else {
		t.Errorf("\tExpected an error!\n")
	}
}

func TestMinDelete(t *testing.T) {
	tree := NewTree()

	// Create a slice of random integers
	size := uint(1 << 20)
	rands := populateTreeAndSlice(t, tree, size)

	sort.Ints(rands)
	for i := uint(0); i < size; i++ {
		listMin := rands[0]
		treeMin, err := tree.Min()
		if err != nil {
			t.Errorf("\t%v\n", err)
		}
		if Integer(listMin) != treeMin {
			t.Errorf("listMin = %d, treeMin = %d\n", listMin, treeMin)
		}
		rands = rands[1:]
		if err := tree.Delete(treeMin); err != nil {
			t.Errorf("\t%v\n", err)
		}
	}
}

func TestMaxDelete(t *testing.T) {
	tree := NewTree()

	// Create a slice of random integers
	size := uint(1 << 20)
	rands := populateTreeAndSlice(t, tree, size)

	sort.Ints(rands)
	for i := uint(0); i < size; i++ {
		listMax := rands[len(rands)-1]
		treeMax, err := tree.Max()
		if err != nil {
			t.Errorf("\t%v\n", err)
		}
		if Integer(listMax) != treeMax {
			t.Errorf("listMax = %d, treeMax = %d\n", listMax, treeMax)
		}
		rands = rands[:len(rands)-1]
		if err := tree.Delete(treeMax); err != nil {
			t.Errorf("\t%v\n", err)
		}
	}
}

func TestHeight(t *testing.T) {
	tree := NewTree()

	t.Logf("Height for no keys: %d\n\n", tree.Height())

	if err := tree.Insert(Integer(0)); err != nil {
		t.Errorf("\t%v\n", err)
	}
	t.Logf("Height for 1 key: %d\n\n", tree.Height())

	// for exp=29, more than 14G of memory are required
	// for exp=28, ~ 8.5G - 10.5G of memory are required (I didn't notice the exact amount)
	for exp := uint(1); exp < 24; exp++ {
		// Insert new keys from range [2**(e-1), (2**e)-2] --> 2**(e-1)-2 new keys.
		for i := 1 << (exp - 1); i < (1<<exp)-1; i++ {
			if err := tree.Insert(Integer(i)); err != nil {
				t.Errorf("\t%v\n", err)
			}
		}
		t.Logf("Height for %d keys: %d\n", (1<<exp)-1, tree.Height())
		//t.Logf("\tPreorder: %v\n", preOrder(t, tree.root))
		if tree.Height() != int(exp) {
			t.Errorf("\tHeight for %d keys is expected to be %d.\n", (1<<exp)-1, exp)
		}

		// Insert 2**e -th key, which should increase tree's height by 1.
		if err := tree.Insert(Integer((1 << exp) - 1)); err != nil {
			t.Errorf("\t%v\n", err)
		}
		t.Logf("Height for %d keys: %d\n", 1<<exp, tree.Height())
		//t.Logf("\tPreorder: %v\n", preOrder(t, tree.root))
		if tree.Height() != int(exp+1) {
			t.Errorf("\tHeight for %d keys is expected to be %d.\n", 1<<exp, exp+1)
		}

		// Insert a 2**e+1 -th key, which shouldn't increase tree's height, and then remove it again.
		if err := tree.Insert(Integer(-42)); err != nil {
			t.Errorf("\t%v\n", err)
		}
		t.Logf("Height for %d keys: %d\n", (1<<exp)+1, tree.Height())
		//t.Logf("\tPreorder: %v\n", preOrder(t, tree.root))
		if tree.Height() != int(exp+1) {
			t.Errorf("\tHeight for %d keys is expected to be %d.\n", (1<<exp)+1, exp+1)
		}
		if err := tree.Delete(Integer(-42)); err != nil {
			t.Errorf("\t%v\n", err)
		}
		t.Logf("\n")
	}
}

func TestSize(t *testing.T) {
	tree := NewTree()

	size := 1 << 20
	for i := 0; i < size; i++ {
		if tree.Size() != i {
			t.Errorf("\ttree.Size() returned %d; expected %d\n", tree.Size(), i)
			t.Errorf("\t ^ Inorder: %v\n", inOrder(t, tree.root))
		}
		if err := tree.Insert(Integer(i)); err != nil {
			t.Fatalf("\t%v\n", err)
		}
	}
	for i := 0; i < size; i += 2 {
		if err := tree.Delete(Integer(i)); err != nil {
			t.Fatalf("\t%v\n", err)
		}
	}
	if tree.Size() != size/2 {
		t.Errorf("\ttree.Size() returned %d; expected %d\n", tree.Size(), size/2)
		t.Logf("\t ^ Inorder: %v\n", inOrder(t, tree.root))
	}
}
