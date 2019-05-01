package p1

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type MerklePatriciaTrie struct {
	db     map[string]Node
	root   string
	keyVal map[string]string
}

func (mpt *MerklePatriciaTrie) Get(key string) (val string, err error) {

	asciiCode := []uint8(key)
	hexArr := keybytesToHex(asciiCode)
	if mpt == nil {
		val = ""
		err = errors.New("mpt is empty")
		return val, err
	}

	if mpt.db == nil || mpt.root == "" {
		val = ""
		err = errors.New("mpt is empty")
		return val, err
	}

	node := mpt.db[mpt.root]
	return mpt.getHelper(node, hexArr)

}

func (mpt *MerklePatriciaTrie) getHelper(node Node, hexPath []uint8) (string, error) {

	if node.node_type == 0 {
		return "", errors.New("nil node exception")
	}
	nodeType := node.nodeType()

	if nodeType == 1 { // Branch
		if len(hexPath) == 0 {
			return node.branch_value[16], nil
		}
		nextNode := mpt.db[node.branch_value[hexPath[0]]]
		return mpt.getHelper(nextNode, hexPath[1:])
	} else if nodeType == 2 { // Extension
		if len(hexPath) == 0 { // Does not match
			return "", errors.New("could not find the val in mpt 1")
		} else {
			hexNibbles := compact_decode(node.flag_value.encoded_prefix)
			cArrLen := commonArrLen(hexNibbles, hexPath)

			if cArrLen < len(hexNibbles) { //
				return "", errors.New("could not find the val in mpt 2")
			} else {
				nextNode := mpt.db[node.flag_value.value]
				return mpt.getHelper(nextNode, hexPath[cArrLen:])
			}
		}
	} else if nodeType == 3 { // Leaf
		hexNibbles := compact_decode(node.flag_value.encoded_prefix)
		cArrLen := commonArrLen(hexNibbles, hexPath)
		if cArrLen == len(hexNibbles) && cArrLen == len(hexPath) {
			return node.flag_value.value, nil
		} else {
			return "", errors.New("could not find the val in mpt 1")
		}
	}

	return "", errors.New("Unexpected node type: " + string(nodeType))
}

/**
Description: Insert() function takes a pair of <key, value> as arguments.
It will traverse down the Merkle Patricia Trie, find the right place to
insert the value, and do the insertion.

Arguments: key(String), value(String)
Return: None
*/
func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	//if the len(mpt.db) == 0, insert a leaf node as root
	if mpt.keyVal == nil {
		mpt.keyVal = map[string]string{}
	}
	mpt.keyVal[key] = new_value
	if mpt.root == "" {

		newNode := CreateLeaf(keybytesToHex([]uint8(key)), new_value)
		mpt.root = newNode.hash_node()
		if mpt.db == nil {
			mpt.db = make(map[string]Node)
		}
		mpt.db[newNode.hash_node()] = newNode
		return
	}

	hexPath := keybytesToHex([]uint8(key))

	rootNode := mpt.db[mpt.root]
	delete(mpt.db, mpt.root)
	newRoot := mpt.insertHelper(&rootNode, hexPath, new_value)
	mpt.root = newRoot

}

func (mpt *MerklePatriciaTrie) insertHelper(curNode *Node, hexPath []uint8, newVal string) string {

	if curNode == nil {
		newLeaf := CreateLeaf(hexPath, newVal)
		mpt.db[newLeaf.hash_node()] = newLeaf
		return newLeaf.hash_node()
	}

	delete(mpt.db, curNode.hash_node())
	nodeType := curNode.nodeType()

	if nodeType == 1 { // branch
		return mpt.branchInsert(curNode, hexPath, newVal)
	} else if nodeType == 2 { // extension
		return mpt.extensionInsert(curNode, hexPath, newVal)
	} else { // leaf
		return mpt.leafInsert(curNode, hexPath, newVal)
	}
}

func (mpt *MerklePatriciaTrie) branchInsert(curNode *Node, hexPath []uint8, newVal string) (nodeHash string) {
	if len(hexPath) == 0 { // base
		curNode.branch_value[16] = newVal
		mpt.db[curNode.hash_node()] = *curNode
		return curNode.hash_node()
	} else {
		route := hexPath[0]
		// case 1, if the branchVal[route] is empty: curNode -> leaf
		if curNode.branch_value[route] == "" {
			nextLeaf := CreateLeaf(hexPath[1:], newVal)
			curNode.branch_value[route] = nextLeaf.hash_node()

			mpt.db[nextLeaf.hash_node()] = nextLeaf
			mpt.db[curNode.hash_node()] = *curNode
			return curNode.hash_node()
		} else { // case 2, if the branchVal[route] is not empty: nextNode = branchVal[route], insert
			nextNode := mpt.db[curNode.branch_value[route]]
			nextHash := mpt.insertHelper(&nextNode, hexPath[1:], newVal)

			curNode.branch_value[route] = nextHash
			mpt.db[curNode.hash_node()] = *curNode
			return curNode.hash_node()
		}
	}
}

func (mpt *MerklePatriciaTrie) extensionInsert(curNode *Node, hexPath []uint8, newVal string) (nodeHash string) {
	curNodeHex := compact_decode(curNode.flag_value.encoded_prefix)

	commonLen := commonArrLen(curNodeHex, hexPath)
	// case1 len(hexPath) == 0:
	// 1.1 curHexLen == 1: newBranch(newVal) -> curNode.next
	// 1.2 else : newBranch(newVal) -> shortExt -> curNode.next

	//case2 commonLen == len(curNodeHex): -> insertHelper(nextNode, remainHexPath, newVal)

	// case3 commonLen == 0:
	// case3.1: len(curNodeHex) == 1: newBranch -> curNode.next, newBranch -> newLeaf(hexPath[1:], newVal)
	// case3.2: else: newBranch -> newExt(curHex[1:]) -> curNode.next, newBranch -> newLeaf(hexPath[1:], newVal)

	//case4 commonLen < len(curNodeHex)
	// 4.1 commonLen == len(hexPath):
	// 4.1.1 commonLen < len(curNodeHex) - 1: 1. new branch with new_val, 2. new extension1 [:commonLen] -> new branch(newVal) -> new extension2[commonLen + 1:] -> nextNode
	// 4.1.2 else: newExtension1(curNodeHex[:commonLen]) -> newBranch(newVal) -> curNode.next
	// 4.2 commonLen < len(hexPath):
	// 4.2.1 commonLen < len(curNodeHex) - 1: shortExt(curHex[1:]) -> newBranch -> shortExt2(curHex[commonLen + 1:]) -> curNode.next, newBranch -> newLeaf(newVal, hexPath[commonLen + 1:])
	// 4.2.2 else: new extension1 [:commonLen] -> newBranch -> nextNode, newBranch -> newLeaf(newVal, hexPath[commonLen + 1:])

	if len(hexPath) == 0 { //case 1 len(hexPath) == 0:
		// 1.1 curHex == 1: newBranch(newVal) -> curNode.next
		if len(curNodeHex) == 1 {
			newBranch := CreateBranch(newVal)
			nextNode := mpt.db[curNode.flag_value.value]

			newBranch.branch_value[curNodeHex[0]] = nextNode.hash_node()

			mpt.db[newBranch.hash_node()] = newBranch
			return newBranch.hash_node()

		} else { // 1.2 else (curHex > 1) : newBranch(newVal) -> shortExt -> curNode.next
			newBranch := CreateBranch(newVal)
			nextNode := mpt.db[curNode.flag_value.value]
			shortExt := CreateExtension(curNodeHex[1:])

			shortExt.flag_value.value = nextNode.hash_node()
			newBranch.branch_value[curNodeHex[0]] = shortExt.hash_node()

			mpt.db[newBranch.hash_node()] = newBranch
			mpt.db[shortExt.hash_node()] = shortExt
			return newBranch.hash_node()
		}
	} else if commonLen == len(curNodeHex) { // case2: perfect match, -> go to next node -> update cur hash
		nextNode := mpt.db[curNode.flag_value.value]
		nextHash := mpt.insertHelper(&nextNode, hexPath[commonLen:], newVal)
		curNode.flag_value.value = nextHash

		mpt.db[curNode.hash_node()] = *curNode
		return curNode.hash_node()

	} else if commonLen == 0 { //case3 commonLen == 0:

		// case3.1: len(curNodeHex) == 1: newBranch -> curNode.next, newBranch -> newLeaf(hexPath[1:], newVal)
		if len(curNodeHex) == 1 {
			newBranch := CreateBranch("")
			nextNode := mpt.db[curNode.flag_value.value]
			newLeaf := CreateLeaf(hexPath[1:], newVal)

			newBranch.branch_value[curNodeHex[0]] = nextNode.hash_node()
			newBranch.branch_value[hexPath[0]] = newLeaf.hash_node()

			mpt.db[newBranch.hash_node()] = newBranch
			mpt.db[newLeaf.hash_node()] = newLeaf
			return newBranch.hash_node()
		} else { // case3.2: else: newBranch -> newExt(curHex[1:]) -> curNode.next, newBranch -> newLeaf(hexPath[1:], newVal)
			newBranch := CreateBranch("")
			newExt := CreateExtension(curNodeHex[1:])
			nextNode := mpt.db[curNode.flag_value.value]
			newLeaf := CreateLeaf(hexPath[1:], newVal)

			newExt.flag_value.value = nextNode.hash_node()
			newBranch.branch_value[curNodeHex[0]] = newExt.hash_node()
			newBranch.branch_value[hexPath[0]] = newLeaf.hash_node()

			mpt.db[newExt.hash_node()] = newExt
			mpt.db[newLeaf.hash_node()] = newLeaf
			mpt.db[newBranch.hash_node()] = newBranch
			return newBranch.hash_node()
		}

	} else { //case4 commonLen < curNodeHex
		// 4.1 commonLen == len(hexPath):
		if commonLen == len(hexPath) {
			if commonLen < len(curNodeHex)-1 { // 4.1.1 commonLen < len(curNodeHex) - 1: 1. new branch with new_val, 2. new extension1 [:commonLen] -> new branch(newVal) -> new extension2[commonLen + 1:] -> nextNode
				newBranch := CreateBranch(newVal)
				newExt1 := CreateExtension(curNodeHex[:commonLen])
				newExt2 := CreateExtension(curNodeHex[commonLen+1:])
				nextNode := mpt.db[curNode.flag_value.value]

				newExt2.flag_value.value = nextNode.hash_node()
				newBranch.branch_value[curNodeHex[commonLen]] = newExt2.hash_node()
				newExt1.flag_value.value = newBranch.hash_node()

				mpt.db[newExt1.hash_node()] = newExt1
				mpt.db[newBranch.hash_node()] = newBranch
				mpt.db[newExt2.hash_node()] = newExt2
				return newExt1.hash_node()
			} else { // 4.1.2 else: newExtension1(curNodeHex[:commonLen]) -> newBranch(newVal) -> curNode.next
				newExt := CreateExtension(curNodeHex[:commonLen])
				newBranch := CreateBranch(newVal)
				nextNode := mpt.db[curNode.flag_value.value]

				newBranch.branch_value[curNodeHex[commonLen]] = nextNode.hash_node()
				newExt.flag_value.value = newBranch.hash_node()

				mpt.db[newExt.hash_node()] = newExt
				mpt.db[newBranch.hash_node()] = newBranch
				return newExt.hash_node()
			}
		} else { // 4.2 commonLen < hexPath:
			if commonLen < len(curNodeHex)-1 { // 4.2.1 commonLen < len(curNodeHex) - 1: shortExt1(curHex[:commonLen]) -> newBranch -> shortExt2(curHex[commonLen + 1:]) -> curNode.next, newBranch -> newLeaf(newVal, hexPath[1:])

				newExt1 := CreateExtension(curNodeHex[:commonLen])
				newBranch := CreateBranch("")
				newExt2 := CreateExtension(curNodeHex[commonLen+1:])

				nextNode := mpt.db[curNode.flag_value.value]
				newLeaf := CreateLeaf(hexPath[commonLen+1:], newVal)

				newExt2.flag_value.value = nextNode.hash_node()
				//newBranch.branch_value[curNodeHex[commonLen]] = newExt2.hash_node()
				//newBranch.branch_value[hexPath[commonLen]] = newLeaf.hash_node()
				newBranch.branch_value[curNodeHex[commonLen]] = newExt2.hash_node()
				newBranch.branch_value[hexPath[commonLen]] = newLeaf.hash_node()
				newExt1.flag_value.value = newBranch.hash_node()

				mpt.db[newExt1.hash_node()] = newExt1
				mpt.db[newBranch.hash_node()] = newBranch
				mpt.db[newExt2.hash_node()] = newExt2
				mpt.db[newLeaf.hash_node()] = newLeaf

				return newExt1.hash_node()

			} else { // 4.2.2 else: new extension1 [:commonLen] -> newBranch -> nextNode, newBranch -> newLeaf(newVal, hexPath[commonLen + 1:])
				newExt1 := CreateExtension(curNodeHex[:commonLen])
				newBranch := CreateBranch("")
				nextNode := mpt.db[curNode.flag_value.value]
				newLeaf := CreateLeaf(hexPath[commonLen+1:], newVal)

				newBranch.branch_value[curNodeHex[commonLen]] = nextNode.hash_node()
				newBranch.branch_value[hexPath[commonLen]] = newLeaf.hash_node()
				newExt1.flag_value.value = newBranch.hash_node()

				mpt.db[newExt1.hash_node()] = newExt1
				mpt.db[newBranch.hash_node()] = newBranch
				mpt.db[newLeaf.hash_node()] = newLeaf
				return newExt1.hash_node()
			}
		}
	}
}

func (mpt *MerklePatriciaTrie) leafInsert(curNode *Node, hexPath []uint8, newVal string) (nodeHash string) {
	curNodeHex := compact_decode(curNode.flag_value.encoded_prefix)
	commonLen := commonArrLen(curNodeHex, hexPath)

	// case 1 same hex, update the value
	// case 2 commonLen == 0
	// case 2.1 curNode hex = "" : newBranch(curNodeVal) -> leaf1(newVal)
	// case 2.2 hexPath hex = "" : newBranch(new_val) -> leaf1
	// case 2.3 else: newBranch -> leaf 1, newBranch -> leaf 2
	// case3 commonLen == len(curNodeHex): newExt(common) -> new branch(value = curNode.value) -> new leaf (hexPath[commonLen + 1:], newVal)
	// case4 commonLen == len(hexPath): newExt(common) -> newBranch(value = newVal) -> new leaf (curHex[commonLen + 1:], curVal)
	// case5 commonLen < len curNodeHex && commonLen < len hexPath:  newExt -> new branch -> new leaf1 (hexPath[commonLen:], newVal), new branch -> new leaf2 (curNodeHex[commonLen:], curNode.value)

	if commonLen == len(curNodeHex) && commonLen == len(hexPath) { // case 1 same hex, update the value
		curNode.flag_value.value = newVal

		mpt.db[curNode.hash_node()] = *curNode
		return curNode.hash_node()
	} else if commonLen == 0 { // case2 commonLen == 0

		// case 2.1 curNode hex = "" : newBranch(curNodeVal) -> leaf1(newVal)
		// case 2.2 hexPath hex = "" : newBranch(new_val) -> leaf1
		// case 2.3 else: newBranch -> leaf 1, newBranch -> leaf 2

		if len(curNodeHex) == 0 { // case 2.1 curNode hex = "" : newBranch(curNodeVal) -> leaf1(newVal)
			newBranch := CreateBranch(curNode.flag_value.value)
			nextLeaf := CreateLeaf(hexPath[1:], newVal)
			newBranch.branch_value[hexPath[0]] = nextLeaf.hash_node()

			mpt.db[nextLeaf.hash_node()] = nextLeaf
			mpt.db[newBranch.hash_node()] = newBranch
			return newBranch.hash_node()
		} else if len(hexPath) == 0 { // case 2.2 hexPath hex = "" : newBranch(new_val) -> leaf1
			newBranch := CreateBranch(newVal)
			nextLeaf := CreateLeaf(curNodeHex[1:], curNode.flag_value.value)
			newBranch.branch_value[curNodeHex[0]] = nextLeaf.hash_node()

			mpt.db[nextLeaf.hash_node()] = nextLeaf
			mpt.db[newBranch.hash_node()] = newBranch
			return newBranch.hash_node()
		} else { // case 2.3 else: newBranch -> leaf 1, newBranch -> leaf 2
			newBranch := CreateBranch("")
			nextLeaf1 := CreateLeaf(curNodeHex[1:], curNode.flag_value.value)
			nextLeaf2 := CreateLeaf(hexPath[1:], newVal)

			newBranch.branch_value[curNodeHex[0]] = nextLeaf1.hash_node()
			newBranch.branch_value[hexPath[0]] = nextLeaf2.hash_node()

			mpt.db[nextLeaf1.hash_node()] = nextLeaf1
			mpt.db[nextLeaf2.hash_node()] = nextLeaf2
			mpt.db[newBranch.hash_node()] = newBranch

			return newBranch.hash_node()
		}

	} else if commonLen == len(curNodeHex) { // case4 commonLen == len(curNodeHex): newExt(common) -> new branch(value = curNode.value) -> new leaf (hexPath[commonLen + 1:], newVal)
		newExt := CreateExtension(hexPath[:commonLen])
		newBranch := CreateBranch(curNode.flag_value.value)
		newLeaf := CreateLeaf(hexPath[commonLen+1:], newVal)

		newBranch.branch_value[hexPath[commonLen]] = newLeaf.hash_node()
		newExt.flag_value.value = newBranch.hash_node()

		mpt.db[newLeaf.hash_node()] = newLeaf
		mpt.db[newBranch.hash_node()] = newBranch
		mpt.db[newExt.hash_node()] = newExt

		return newExt.hash_node()

	} else if commonLen == len(hexPath) { // case5 commonLen == len(hexPath): newExt(common) -> newBranch(value = newVal) -> new leaf (curHex[commonLen + 1:], curVal)
		newExt := CreateExtension(hexPath[:commonLen])
		newBranch := CreateBranch(newVal)
		newLeaf := CreateLeaf(curNodeHex[commonLen+1:], curNode.flag_value.value)

		newBranch.branch_value[curNodeHex[commonLen]] = newLeaf.hash_node()
		newExt.flag_value.value = newBranch.hash_node()

		mpt.db[newLeaf.hash_node()] = newLeaf
		mpt.db[newBranch.hash_node()] = newBranch
		mpt.db[newExt.hash_node()] = newExt

		return newExt.hash_node()

	} else { // case6 commonLen < len curNodeHex && commonLen < len hexPath:  newExt -> new branch -> new leaf1 (hexPath[commonLen:], newVal), new branch -> new leaf2 (curNodeHex[commonLen:], curNode.value)
		newExt := CreateExtension(hexPath[:commonLen])
		newBranch := CreateBranch("")
		newLeaf1 := CreateLeaf(curNodeHex[commonLen+1:], curNode.flag_value.value)
		newLeaf2 := CreateLeaf(hexPath[commonLen+1:], newVal)

		newBranch.branch_value[curNodeHex[commonLen]] = newLeaf1.hash_node()
		newBranch.branch_value[hexPath[commonLen]] = newLeaf2.hash_node()
		newExt.flag_value.value = newBranch.hash_node()

		mpt.db[newExt.hash_node()] = newExt
		mpt.db[newBranch.hash_node()] = newBranch
		mpt.db[newLeaf1.hash_node()] = newLeaf1
		mpt.db[newLeaf2.hash_node()] = newLeaf2

		return newExt.hash_node()
	}
}

func (mpt *MerklePatriciaTrie) branchConnect(branch Node, node Node, index uint8) {

	node.branch_value[index] = node.hash_node()
}

func (mpt *MerklePatriciaTrie) extensionConnect(ext Node, node Node) {
	node.flag_value.value = node.hash_node()
}

// return the common array len

func commonArrLen(arr1 []uint8, arr2 []uint8) (index int) {
	for i := 0; i < min(len(arr1), len(arr2)); i++ {
		if arr1[i] != arr2[i] {
			return i
		}
	}
	return min(len(arr1), len(arr2))
}

func min(i int, j int) int {
	if i > j {
		return j
	} else {
		return i
	}
}

func (mpt *MerklePatriciaTrie) Delete(key string) string {

	oldRoot := mpt.db[mpt.root]
	if mpt.root == "" {
		return "path_not_found"
	}
	if mpt.db == nil || len(mpt.db) == 0 {
		return "path_not_found"
	}
	binaryArr := []uint8(key)
	rootNode := mpt.db[mpt.root]

	newRoot := mpt.deleteHelper(&rootNode, keybytesToHex(binaryArr))
	mpt.db[newRoot.hash_node()] = *newRoot
	mpt.root = newRoot.hash_node()
	if oldRoot.hash_node() == newRoot.hash_node() {
		return "path_not_found"
	} else {
		return ""
	}
}

func (mpt *MerklePatriciaTrie) deleteHelper(curNode *Node, hexPath []uint8) (node *Node) {

	// Remove curNode in db
	// 1. curNode is a leaf
	// 1. curHex matches hexPath => return ""
	// 2. curHex does not match hexPath => return curNode

	// 2. curNode is a branch:
	// 2.1. hexPath len == 0: remove the val in branch_val[16]

	// 2.1.1 if curNode.branch_val len == 1:
	// Delete the remain in db
	// 2.1.1.1 remained node is a branch: newExt(index) -> remain branch
	// 2.1.1.2 remained node is a leaf: newLeaf(index + remainHex , leaf.val)
	// 2.1.1.3 remained node is an extension: newExt(index + remain) -> remain.next return newExt
	// 2.1.2 else: delete	curNode.branchVal[16]

	// 2.2. hexPath len != 0:
	// returned node = deleteHelper(nextNode, hexPath[1:]) ->
	// remove the branch.val[hexPath[0]]
	// 2.2.1. if the returned node.type == 0 // nil
	// 2.2.1.1. : if the curBranch.branch_val[16] == "" && curBranch.branch_val only has 1 path
	// 2.2.1.1.1 remain node is a branch: newExt(index) -> nextNode
	// 2.2.1.1.2 remain node is a extension: 1. rm extension in db 2. newExt(index + extension.prefix) -> nextNode.next
	// 2.2.1.1.3 remain node is a leaf: 1. rm leaf 2. newLeaf(index + leaf.prefix, leaf.val)
	// 2.2.1.2. : else if curBranch.branch_val only has 0 path and which implies curBranch.branch_val[16] != "":  newCurNode = newLeaf(branch_val[16]), return newCurNode.hash_node
	// -> newLeaf(branch_val[16])
	// 2.2.1.3 : else (space is enough for keeping the branch): return curNode

	// 2.2.3. else (the returned flag_value.type != 0): curNode.branchVal[hexPath[0]] = returned node.hash_node

	// 3. curNode is an extension:
	// 3.1 commonLen == curNodeHexLen
	// 3.1.1 returned node is nil
	// 3.1.2 returned node is branch: curNode.next = returnedBranch
	// 3.1.3 returned node is extension: newExt(curExtension, returnedExtension)
	// 3.1.4 returned node is leaf: newLeaf(curExtension, returnedLeaf)
	// 3.2 else (commonLen != curNodeHexLen): return

	// Remove curNode in db
	// 1. curNode is a leaf
	if curNode.nodeType() == 0 {
		return curNode
	} else if curNode.nodeType() == 3 {
		return mpt.leafDelete(curNode, hexPath)
	} else if curNode.nodeType() == 1 {
		res := mpt.branchDelete(curNode, hexPath)
		return res
	} else if curNode.nodeType() == 2 {
		return mpt.extensionDelete(curNode, hexPath)
	}
	return &Node{}
}

func (mpt MerklePatriciaTrie) leafDelete(curNode *Node, hexPath []uint8) *Node {
	// 1.1. curHex matches hexPath => return ""
	delete(mpt.db, curNode.hash_node())
	curNodeHex := compact_decode(curNode.flag_value.encoded_prefix)
	commonLen := commonArrLen(hexPath, curNodeHex)

	if commonLen == len(hexPath) && commonLen == len(curNodeHex) {
		res := Node{node_type: 0}
		return &res
	} else { // 1.2. curHex does not match hexPath => return curNode
		mpt.db[curNode.hash_node()] = *curNode
		return curNode
	}
}

// TODO tests branchDelete
func (mpt MerklePatriciaTrie) branchDelete(curNode *Node, hexPath []uint8) *Node {
	// 2. curNode is a branch:
	// 	2.1. hexPath len == 0: remove the val in branch_val[16]

	// 		2.1.1 if curNode.branch_val len == 1:
	// 		Delete the remain in db
	// 			2.1.1.1 remained node is a branch: newExt(index) -> remain branch
	// 			2.1.1.2 remained node is a leaf: newLeaf(index + remainHex , leaf.val)
	// 			2.1.1.3 remained node is an extension: newExt(index + remain) -> remain.next return newExt
	// 		2.1.2 else: delete	curNode.branchVal[16]

	// 	2.2. hexPath len != 0:
	// 	Check if we can find the next node, if not, return cur node
	// 	returned node = deleteHelper(nextNode, hexPath[1:]) ->
	// 	remove the branch.val[hexPath[0]]
	// 		2.2.1. if the returned node.type == 0 // nil
	// 			2.2.1.1. : if the curBranch.branch_val[16] == "" && curBranch.branch_val only has 1 path
	// 				2.2.1.1.1 remain node is a branch: newExt(index) -> nextNode
	// 				2.2.1.1.2 remain node is a extension: 1. rm extension in db 2. newExt(index + extension.prefix) -> nextNode.next
	// 				2.2.1.1.3 remain node is a leaf: 1. rm leaf 2. newLeaf(index + leaf.prefix, leaf.val)
	// 			2.2.1.2. : else if curBranch.branch_val only has 0 path and which implies curBranch.branch_val[16] != "":  newCurNode = newLeaf(branch_val[16]), return newCurNode.hash_node
	// 			-> newLeaf(branch_val[16])
	// 			2.2.1.3 : else (space is enough for keeping the branch): return curNode

	// 		2.2.3. else (the returned flag_value.type != 0): curNode.branchVal[hexPath[0]] = returned node.hash_node

	delete(mpt.db, curNode.hash_node())
	if len(hexPath) == 0 { // 2.1. hexPath len == 0: remove the val in branch_val[16]

		if curNode.branchSize() == 1 { // 2.1.1 if curNode.branch_val len == 1:
			index, remainHash := curNode.branchFirstHashPath()
			remain := mpt.db[remainHash]
			// Delete the remain in db
			delete(mpt.db, remainHash)
			if remain.nodeType() == 1 { // 2.1.1.1 remained node is a branch: newExt(index) -> remain branch
				newExt := CreateExtension([]uint8{uint8(index)})
				newExt.flag_value.value = remain.hash_node()

				mpt.db[newExt.hash_node()] = newExt
				mpt.db[remain.hash_node()] = remain
				return &newExt
			} else if remain.nodeType() == 3 { // 2.1.1.2 remained node is a leaf: newLeaf(index + remainHex , leaf.val)
				remainHex := compact_decode(remain.flag_value.encoded_prefix)
				newLeafHex := append([]uint8{uint8(index)}, remainHex...)
				newLeaf := CreateLeaf(newLeafHex, remain.flag_value.value)

				mpt.db[newLeaf.hash_node()] = newLeaf
				return &newLeaf
			} else { // 2.1.1.3 remained node is an extension: newExt(index + remain) -> remain.next return newExt
				remainHex := compact_decode(remain.flag_value.encoded_prefix)

				newExtHex := append([]uint8{uint8(index)}, remainHex...)
				newExt := CreateExtension(newExtHex)
				nextNode := mpt.db[remain.flag_value.value]
				newExt.flag_value.value = nextNode.hash_node()

				mpt.db[newExt.hash_node()] = newExt
				return &newExt
			}
		} else { // 2.1.2 else: delete	curNode.branchVal[16]
			curNode.branch_value[16] = ""
			mpt.db[curNode.hash_node()] = *curNode
			return curNode
		}
	} else { // 2.2. hexPath len != 0:
		nextNode, ok := mpt.db[curNode.branch_value[hexPath[0]]] // Check if we can find the next node
		if !ok {
			mpt.db[curNode.hash_node()] = *curNode
			return curNode
		}
		// returned node = deleteHelper(nextNode, hexPath[1:]) ->

		returnedNode := mpt.deleteHelper(&nextNode, hexPath[1:])
		// remove the branch.val[hexPath[0]]
		curNode.branch_value[hexPath[0]] = ""

		if returnedNode.node_type == 0 { // 2.2.1. if the returned node.type == 0 // nil

			// 2.2.1.1. : if the curBranch.branch_val[16] == "" && curBranch.branch_val only has 1 path:
			if curNode.branch_value[16] == "" && curNode.branchSize() == 1 {

				remainIndex, remainNodeHash := curNode.branchFirstHashPath()
				remainNode := mpt.db[remainNodeHash]
				if remainNode.nodeType() == 1 { // 2.2.1.1.1 remain node is a branch: newExt(index) -> remainNode
					newExt := CreateExtension([]uint8{uint8(remainIndex)})
					newExt.flag_value.value = remainNode.hash_node()

					mpt.db[newExt.hash_node()] = newExt
					mpt.db[remainNode.hash_node()] = remainNode
					return &newExt
				} else if remainNode.nodeType() == 2 { // 2.2.1.1.2 remain node is a extension: 1. rm extension in db 2. newExt(index + extension.prefix) -> remainNode.next
					delete(mpt.db, remainNodeHash)
					remainHex := compact_decode(remainNode.flag_value.encoded_prefix)
					newExt := CreateExtension(append([]uint8{uint8(remainIndex)}, remainHex...))
					remainNodeNext := mpt.db[remainNode.flag_value.value]
					newExt.flag_value.value = remainNodeNext.hash_node()
					mpt.db[newExt.hash_node()] = newExt

					return &newExt
				} else { // 2.2.1.1.3 remain node is a leaf: 1. rm leaf 2. newLeaf(index + leaf.prefix, leaf.val)
					delete(mpt.db, remainNodeHash)
					newLeafHex := append([]uint8{uint8(remainIndex)}, compact_decode(remainNode.flag_value.encoded_prefix)...)
					newLeaf := CreateLeaf(newLeafHex, remainNode.flag_value.value)

					mpt.db[newLeaf.hash_node()] = newLeaf
					return &newLeaf
				}

			} else if curNode.branchSize() == 0 { // 2.2.1.2. : else if curBranch.branch_val only has 0 path and which implies curBranch.branch_val[16] != "":
				// newCurNode = newLeaf(branch_val[16]), return newCurNode
				// -> newLeaf(branch_val[16])
				newLeaf := CreateLeaf([]uint8{}, curNode.branch_value[16])
				mpt.db[newLeaf.hash_node()] = newLeaf
				return &newLeaf
			} else { // 2.2.1.3 : else (space is enough for keeping the branch): return curNode
				mpt.db[curNode.hash_node()] = *curNode
				return curNode
			}

		} else { // 2.2.3. else (the returned flag_value.type != 0): curNode.branchVal[hexPath[0]] = returned node.hash_node
			curNode.branch_value[hexPath[0]] = returnedNode.hash_node()

			mpt.db[curNode.hash_node()] = *curNode
			return curNode
		}
	}
}

func (mpt MerklePatriciaTrie) extensionDelete(curNode *Node, hexPath []uint8) *Node {
	// 3. curNode is an extension:
	// 3.1 commonLen == curNodeHexLen
	// 3.1.1 returned node is nil
	// 3.1.2 returned node is branch: curNode.next = returnedBranch
	// 3.1.3 returned node is extension: newExt(curExtension, returnedExtension)
	// 3.1.4 returned node is leaf: newLeaf(curExtension, returnedLeaf)
	// 3.2 else (commonLen != curNodeHexLen): return

	delete(mpt.db, curNode.hash_node())
	curNodeHex := compact_decode(curNode.flag_value.encoded_prefix)
	commonLen := commonArrLen(hexPath, curNodeHex)

	if commonLen == len(curNodeHex) { // 3.1 commonLen == curNodeHexLen
		nextNode := mpt.db[curNode.flag_value.value]
		returnedNode := mpt.deleteHelper(&nextNode, hexPath[commonLen:])
		rNodeType := returnedNode.nodeType()
		if rNodeType == 0 { // 3.1.1 returned node is nil
			mpt.db[curNode.hash_node()] = *curNode
			return curNode
		} else if rNodeType == 1 { // 3.1.2 returned node is branch: curNode.next = returnedBranch
			curNode.flag_value.value = returnedNode.hash_node()
			mpt.db[curNode.hash_node()] = *curNode
			return curNode
		} else if rNodeType == 2 { // 3.1.3 returned node is extension: newExt(curExtension, returnedExtension) -> returnNode.next
			delete(mpt.db, returnedNode.hash_node())
			newExt := CreateExtension(append(compact_decode(curNode.flag_value.encoded_prefix), compact_decode(returnedNode.flag_value.encoded_prefix)...))
			newExt.flag_value.value = returnedNode.flag_value.value
			mpt.db[newExt.hash_node()] = newExt
			return &newExt
		} else { // 3.1.4 returned node is leaf: newLeaf(curExtension, returnedLeaf)
			delete(mpt.db, returnedNode.hash_node())
			newLeaf := CreateLeaf(append(compact_decode(curNode.flag_value.encoded_prefix), compact_decode(returnedNode.flag_value.encoded_prefix)...), returnedNode.flag_value.value)

			mpt.db[newLeaf.hash_node()] = newLeaf
			return &newLeaf
		}

	} else { // 3.2 else (commonLen != curNodeHexLen): return
		mpt.db[curNode.hash_node()] = *curNode
		return curNode
	}

}

func (mpt *MerklePatriciaTrie) ToMap() map[string]string {
	return mpt.keyVal
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart"+cur_hash+"HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart"+cur_hash+"HashEnd", fmt.Sprintf("hash%v", i), -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.db = make(map[string]Node)
	mpt.root = ""
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) GetRoot() string {
	return mpt.root
}

func (mpt *MerklePatriciaTrie) GetDB() map[string]Node {
	return mpt.db
}
