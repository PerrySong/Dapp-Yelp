package p1

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

type Flag_value struct {
	encoded_prefix []uint8 // Store ASCII num
	value          string  // val, could be nil for branch
}

type Node struct {
	node_type    int        // 0: Null, 1: Branch, 2: Ext or Leaf -> if leaf, append 16 in the end
	branch_value [17]string // store hash and value
	flag_value   Flag_value // Only useful for root, leaf and branch
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

// Return 1: Branch, 2: Extension, 3: Leaf
func (node *Node) nodeType() int {
	if node.node_type == 0 {
		return 0
	}
	if node.node_type == 1 { // Branch
		return 1
	} else {
		prefix := keybytesToHex(node.flag_value.encoded_prefix)[0]
		if prefix&2 == 0 {
			return 2 // Extension
		} else {
			return 3
		}
	}
}

func CreateBranch(val string) Node {
	var branchVal [17]string
	branchVal[16] = val
	return Node{node_type: 1, branch_value: branchVal}
}

func CreateExtension(hexArr []uint8) Node {
	return Node{node_type: 2, flag_value: Flag_value{encoded_prefix: compact_encode(hexArr)}}
}

func CreateLeaf(hexArr []uint8, val string) Node {
	return Node{node_type: 2, flag_value: Flag_value{value: val, encoded_prefix: compact_encode(append(hexArr, 16))}}
}

/**
return the number of path branch_val has
*/
func (node *Node) branchSize() int {
	res := 0
	branchVal := node.branch_value
	for i, element := range branchVal {
		if i != 16 && element != "" {
			res++
		}
	}
	return res
}

/**
return the remain index and hash_node
*/
func (node *Node) branchFirstHashPath() (int, string) {
	for i, hash := range node.branch_value {
		if hash != "" {
			return i, hash
		}
	}
	return 0, ""
}
