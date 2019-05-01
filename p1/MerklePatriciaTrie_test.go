package p1

import (
	"fmt"
	"reflect"
	"testing"
)

// Insert ("a", "apple") Insert ("b", "banana")

func TestMerklePatriciaTrie_Get1(t *testing.T) {
	mpt := makeTree1()
	res, err := mpt.Get("a")

	fmt.Println(err)
	fmt.Println("Get a:", res)
	if !reflect.DeepEqual(res, "apple") {
		t.Errorf("get() want %x -> %x", res, "apple")
	}
}

//Insert ("a", "apple"), Insert ("p", "banana"), Insert ("abc", "new")

func TestMerklePatriciaTrie_Get2(t *testing.T) {

	mpt := makeTree2()
	val1, err1 := mpt.Get("a")
	val2, err2 := mpt.Get("p")
	val3, err3 := mpt.Get("abc")

	if err1 != nil {
		fmt.Println(err1)
	}
	if err2 != nil {
		fmt.Println(err2)
	}
	if err3 != nil {
		fmt.Println(err3)
	}

	fmt.Println(val1, val2, val3)

}

func makeTree1() MerklePatriciaTrie {
	leaf1, leaf2 := Node{node_type: 2, flag_value: Flag_value{encoded_prefix: compact_encode([]uint8{16}), value: "apple"}}, Node{node_type: 2, flag_value: Flag_value{encoded_prefix: compact_encode([]uint8{16}), value: "banana"}}
	branch := Node{node_type: 1, branch_value: [17]string{}}
	branch.branch_value[1] = leaf1.hash_node()
	branch.branch_value[2] = leaf2.hash_node()

	extension := Node{node_type: 2, flag_value: Flag_value{encoded_prefix: compact_encode([]uint8{6}), value: branch.hash_node()}}
	//extension.node_type = 1
	mpt := MerklePatriciaTrie{db: map[string]Node{}}
	mpt.root = extension.hash_node()
	mpt.db[leaf1.hash_node()] = leaf1
	mpt.db[leaf2.hash_node()] = leaf2
	mpt.db[branch.hash_node()] = branch
	mpt.db[extension.hash_node()] = extension
	return mpt
}

//Insert ("a", "apple"), Insert ("p", "banana"), Insert ("abc", "new")
func makeTree2() MerklePatriciaTrie {
	lv1branch1 := Node{node_type: 1, branch_value: [17]string{}}
	lv3branch1 := Node{node_type: 1, branch_value: [17]string{}}
	lv2leaf1 := Node{node_type: 2, flag_value: Flag_value{value: "banana", encoded_prefix: compact_encode([]uint8{0, 16})}}
	lv4leaf1 := Node{node_type: 2, flag_value: Flag_value{value: "new", encoded_prefix: compact_encode([]uint8{2, 6, 3, 16})}}
	lv2ext1 := Node{node_type: 2, flag_value: Flag_value{encoded_prefix: compact_encode([]uint8{1})}}

	lv3branch1.branch_value[6] = lv4leaf1.hash_node()
	lv3branch1.branch_value[16] = "apple"
	lv2ext1.flag_value.value = lv3branch1.hash_node()
	lv1branch1.branch_value[6] = lv2ext1.hash_node()
	lv1branch1.branch_value[7] = lv2leaf1.hash_node()

	mpt := MerklePatriciaTrie{db: map[string]Node{}}
	mpt.root = lv1branch1.hash_node()
	mpt.db[lv1branch1.hash_node()] = lv1branch1
	mpt.db[lv2ext1.hash_node()] = lv2ext1
	mpt.db[lv2leaf1.hash_node()] = lv2leaf1
	mpt.db[lv3branch1.hash_node()] = lv3branch1
	mpt.db[lv4leaf1.hash_node()] = lv4leaf1
	return mpt
}

func TestMerklePatriciaTrie_Get3(t *testing.T) {

	mpt := makeTree3()
	// tests get
	errArr := [4]error{}
	valArr := [4]string{}
	valArr[0], errArr[0] = mpt.Get("p")
	valArr[1], errArr[1] = mpt.Get("aaaaa")
	valArr[2], errArr[2] = mpt.Get("aaaap")
	valArr[3], errArr[3] = mpt.Get("aa")

	for i, err := range errArr {
		if err != nil {
			fmt.Println("i = ", i, " err = ", err)
		}
		fmt.Println("val = ", valArr[i])
	}
}

// insert: (p, apple) (aaaaa, banana) (aaaap, orange) (aa, new)
func makeTree3() MerklePatriciaTrie {
	lv1branch1 := Node{node_type: 1, branch_value: [17]string{}}
	lv3branch1 := Node{node_type: 1, branch_value: [17]string{}}
	lv5branch1 := Node{node_type: 1, branch_value: [17]string{}}

	lv2ext1 := Node{node_type: 2, flag_value: Flag_value{encoded_prefix: compact_encode([]uint8{1, 6, 1})}}
	lv4ext1 := Node{node_type: 2, flag_value: Flag_value{encoded_prefix: compact_encode([]uint8{1, 6, 1})}}

	lv2leaf1 := Node{node_type: 2, flag_value: Flag_value{value: "apple", encoded_prefix: compact_encode([]uint8{0, 16})}}
	lv6leaf1 := Node{node_type: 2, flag_value: Flag_value{value: "banana", encoded_prefix: compact_encode([]uint8{1, 16})}}
	lv6leaf2 := Node{node_type: 2, flag_value: Flag_value{value: "orange", encoded_prefix: compact_encode([]uint8{0, 16})}}

	lv5branch1.branch_value[6] = lv6leaf1.hash_node()
	lv5branch1.branch_value[7] = lv6leaf2.hash_node()

	lv4ext1.flag_value.value = lv5branch1.hash_node()

	lv3branch1.branch_value[16] = "new"
	lv3branch1.branch_value[6] = lv4ext1.hash_node()

	lv2ext1.flag_value.value = lv3branch1.hash_node()
	lv1branch1.branch_value[6] = lv2ext1.hash_node()
	lv1branch1.branch_value[7] = lv2leaf1.hash_node()

	lv1branch1.branch_value[6] = lv2ext1.hash_node()
	lv1branch1.branch_value[7] = lv2leaf1.hash_node()

	mpt := MerklePatriciaTrie{db: map[string]Node{}, root: lv1branch1.hash_node()}
	db := mpt.db

	// Put to db
	db[lv1branch1.hash_node()] = lv1branch1
	db[lv2leaf1.hash_node()] = lv2leaf1
	db[lv2ext1.hash_node()] = lv2ext1
	db[lv3branch1.hash_node()] = lv3branch1
	db[lv4ext1.hash_node()] = lv4ext1
	db[lv5branch1.hash_node()] = lv5branch1
	db[lv6leaf1.hash_node()] = lv6leaf1
	db[lv6leaf2.hash_node()] = lv6leaf2
	return mpt
}

func TestMerklePatriciaTrie_Get4(t *testing.T) {

	mpt := makeTree4()
	// tests get
	errArr := [3]error{}
	valArr := [3]string{}
	valArr[0], errArr[0] = mpt.Get("p")
	valArr[1], errArr[1] = mpt.Get("aa")
	valArr[2], errArr[2] = mpt.Get("ap")

	for i, err := range errArr {
		if err != nil {
			fmt.Println("i = ", i, " err = ", err)
		}
		fmt.Println("val = ", valArr[i])
	}
}

func makeTree4() MerklePatriciaTrie {
	lv1Branch := CreateBranch("")
	lv3Branch := CreateBranch("")
	lv2Ext := CreateExtension([]uint8{1})
	lv2Leaf := CreateLeaf([]uint8{0}, "apple")
	lv4Leaf1 := CreateLeaf([]uint8{1}, "banana")
	lv4Leaf2 := CreateLeaf([]uint8{0}, "orange")

	lv3Branch.branch_value[6] = lv4Leaf1.hash_node()
	lv3Branch.branch_value[7] = lv4Leaf2.hash_node()
	lv2Ext.flag_value.value = lv3Branch.hash_node()
	lv1Branch.branch_value[6] = lv2Ext.hash_node()
	lv1Branch.branch_value[7] = lv2Leaf.hash_node()

	mpt := MerklePatriciaTrie{}
	mpt.db = make(map[string]Node)
	mpt.db[lv4Leaf2.hash_node()] = lv4Leaf2
	mpt.db[lv3Branch.hash_node()] = lv3Branch
	mpt.db[lv4Leaf1.hash_node()] = lv4Leaf1
	mpt.db[lv2Ext.hash_node()] = lv2Ext
	mpt.db[lv2Leaf.hash_node()] = lv2Leaf
	mpt.db[lv1Branch.hash_node()] = lv1Branch
	mpt.root = lv1Branch.hash_node()
	return mpt
}

// Insert ("a", "apple") Insert ("b", "banana")
func TestMerklePatriciaTrie_Insert1(t *testing.T) {
	mpt1 := makeTree1()
	mpt2 := MerklePatriciaTrie{}
	mpt2.Insert("a", "apple")
	mpt2.Insert("b", "banana")

	fmt.Println(mpt2)
	if str, err := mpt2.Get("a"); str != "apple" || err != nil {
		t.Errorf("wnat apple but %s", str)
	}

	//fmt.Println("root = ", mpt.db[mpt.root])

	if str, err := mpt2.Get("b"); str != "banana" || err != nil {
		t.Errorf("wnat banana but %s", str)
	}

	if !reflect.DeepEqual(mpt1, mpt2) {
		t.Errorf("want %+v \n, but %+v", mpt1, mpt2)
	}

}

//Insert ("a", "apple"), Insert ("p", "banana"), Insert ("abc", "new")
func TestMerklePatriciaTrie_Insert2(t *testing.T) {

	mpt1 := makeTree2()
	mpt2 := MerklePatriciaTrie{}

	mpt2.Insert("a", "apple")
	mpt2.Insert("p", "banana")
	mpt2.Insert("abc", "new")

	if !reflect.DeepEqual(mpt1, mpt2) {
		t.Errorf("want %+v \n, but %+v", mpt1, mpt2)
	}
}

// insert: (p, apple) (aaaaa, banana) (aaaap, orange) (aa, new)
func TestMerklePatriciaTrie_Insert3(t *testing.T) {
	mpt1 := makeTree3()
	mpt2 := MerklePatriciaTrie{}
	mpt2.Insert("p", "apple")
	mpt2.Insert("aaaaa", "banana")
	mpt2.Insert("aaaap", "orange")
	mpt2.Insert("aa", "new")

	//obj1, _ := mpt2.Get("p")
	//obj2, _ := mpt2.Get("aaaaa")
	//obj3, _ := mpt2.Get("aaaap")
	//fmt.Printf("Val = %v", obj1)
	//fmt.Printf("Val = %v,", obj2)
	//fmt.Printf("Val = %v,", obj3)
	//fmt.Println(mpt2)
	if !reflect.DeepEqual(mpt1, mpt2) {
		t.Errorf("want %+v \n, but %+v", mpt1, mpt2)
	}
}

func TestMerklePatriciaTrie_Insert4(t *testing.T) {
	mpt1 := makeTree4()
	mpt2 := MerklePatriciaTrie{}
	mpt2.Insert("p", "apple")
	mpt2.Insert("aa", "banana")
	mpt2.Insert("ap", "orange")

	//obj1, _ := mpt2.Get("p")
	//obj2, _ := mpt2.Get("aa")
	//obj3, _ := mpt2.Get("ap")
	//fmt.Printf("Val = %v", obj1)
	//fmt.Printf("Val = %v,", obj2)
	//fmt.Printf("Val = %v,", obj3)
	//
	//obj4, _ := mpt1.Get("p")
	//obj5, _ := mpt1.Get("aa")
	//obj6, _ := mpt1.Get("ap")
	//fmt.Printf("Val = %v", obj4)
	//fmt.Printf("Val = %v,", obj5)
	//fmt.Printf("Val = %v,", obj6)

	fmt.Println(mpt2)
	fmt.Println(mpt1)
	if !reflect.DeepEqual(mpt1, mpt2) {
		t.Errorf("want %+v \n, but %+v", mpt1, mpt2)
	}
}

func TestMerklePatriciaTrie_Delete1(t *testing.T) {
	mpt1 := makeTree1()
	mpt2 := MerklePatriciaTrie{}

	mpt2.Insert("a", "apple")

	mpt1.Delete("b")

	fmt.Println(mpt2.Get("a"))
	fmt.Println(mpt2.Get("b"))
	if !reflect.DeepEqual(mpt1, mpt2) {
		t.Errorf("want %+v \n, but %+v", mpt2, mpt1)
	}
}

//func sameTree(mpt1 MerklePatriciaTrie, mpt2, trie MerklePatriciaTrie) bool {
//	root1 :=
//	root2 :=
//}
