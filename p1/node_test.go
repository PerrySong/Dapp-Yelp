package p1

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNode_CreateBranch(t *testing.T) {
	branch := CreateBranch("apple")
	if branch.branch_value[16] != "apple" {
		t.Error("Create branch error")
	}

}

func TestNode_CreateExtension(t *testing.T) {
	ext := CreateExtension([]uint8{1, 2, 3, 4})
	if !bytes.Equal(ext.flag_value.encoded_prefix, []uint8{0, 18, 52}) {
		t.Errorf("Create ext error, encodedprefix should be [0, 18, 52] but %d", ext.flag_value.encoded_prefix)
	}
	fmt.Println(compact_encode([]uint8{1, 2, 3, 4}))
}

func TestNode_CreateLeaf(t *testing.T) {
	leaf := CreateLeaf([]uint8{6, 1}, "apple")
	if !bytes.Equal(compact_decode(leaf.flag_value.encoded_prefix), []uint8{6, 1}) {
		t.Errorf("Create ext error, encodedprefix should be [32, 18, 52] but %d", leaf.flag_value.encoded_prefix)
	}
	if leaf.flag_value.value != "apple" {
		t.Errorf("Create ext error, val should be apple but %s", leaf.flag_value.value)
	}
}
