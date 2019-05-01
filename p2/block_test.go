package p2

import (
	"encoding/hex"
	"fmt"
	"github.com/cs686-blockchain-p3-PerrySong/p1"
	"strings"
	"testing"
)

func TestNewBlock(t *testing.T) {
	block := NewBlock(1, "prevhash", "nonce+", p1.MerklePatriciaTrie{})
	fmt.Println(block.Header.Nonce)
	w := []byte{1, 3, 5}
	fmt.Println(hex.EncodeToString(w), "what")
}

func TestBlockChainBasic(t *testing.T) {
	//jsonBlockChain := "[{\"Hash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"TimeStamp\": 1234567890, \"Height\": 1, \"ParentHash\": \"genesis\", \"Size\": 1174, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}, {\"Hash\": \"24cf2c336f02ccd526a03683b522bfca8c3c19aed8a1bed1bbc23c33cd8d1159\", \"TimeStamp\": 1234567890, \"Height\": 2, \"ParentHash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"Size\": 1231, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}]"
	//bc, err := DecodeJsonToBlockChain(jsonBlockChain)
	//if err != nil {
	//	fmt.Println(err)
	//	t.Fail()
	//}
	//jsonNew, err := bc.EncodeToJson()
	//if err != nil {
	//	fmt.Println(err)
	//	t.Fail()
	//}
	//var realValue []BlockJson
	//var expectedValue []BlockJson
	//err = json.Unmarshal([]byte(jsonNew), &realValue)
	//if err != nil {
	//	fmt.Println(err)
	//	t.Fail()
	//}
	//err = json.Unmarshal([]byte(jsonBlockChain), &expectedValue)
	//if err != nil {
	//	fmt.Println(err)
	//	t.Fail()
	//}
	//if !reflect.DeepEqual(realValue, expectedValue) {
	//	fmt.Println("=========Real=========")
	//	fmt.Println(realValue)
	//	fmt.Println("=========Expcected=========")
	//	fmt.Println(expectedValue)
	//	t.Fail()
	//}
}

func TestJsonToBlock(t *testing.T) {
	//JsonToBlock()
	json := `{
    "Hash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
    "TimeStamp":1234567890,
    "Height":1,
    "ParentHash":"genesis",
    "Size":1174,
    "mpt":{
        "charles":"ge",
        "hello":"world"
    	}
	}`
	block := JsonToBlock(json)
	if block.Header.Hash != "3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48" {
		t.Errorf("Incorrect Header.Hash want %s but %s", "3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48", block.Header.Hash)
	}
}

func TestBlock_EncodeToJson(t *testing.T) {
	json := `{
    "Hash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
    "TimeStamp":1234567890,
    "Height":1,
    "ParentHash":"genesis",
    "Size":1174,
    "mpt":{
        "charles":"ge",
        "hello":"world"
    	}
	}`
	json = strings.Replace(json, " \n", "", -1)
	block := JsonToBlock(json)
	res := block.EncodeToJson()
	if !strings.Contains(res, "Hash") || !strings.Contains(res, "Size") {
		t.Errorf("Fail to encode to json, want %v but %v", json, res)
	}
}
