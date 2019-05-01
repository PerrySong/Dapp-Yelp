package p2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var blockchain = NewBlockChain()
var block1Json = `{
        "Hash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
        "TimeStamp":1234567890,
        "Height":1,
        "ParentHash":"genesis",
        "Size":1174,
        "mpt":{
            "hello":"world",
            "charles":"ge"
        }
	}`
var block2Json = `{
        "Hash":"24cf2c336f02ccd526a03683b522bfca8c3c19aed8a1bed1bbc23c33cd8d1159",
        "TimeStamp":1234567890,
        "Height":2,
        "ParentHash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
        "Size":1231,
        "mpt":{
            "hello":"world",
            "charles":"ge"
        }
    }`

var block1 = JsonToBlock(block1Json)
var block2 = JsonToBlock(block2Json)

func TestNewBlockChain(t *testing.T) {
	blockchain := NewBlockChain()
	if blockchain.Length != 0 {
		t.Errorf("BlockChain length incorrect")
	}
}

func TestBlockchain_Get(t *testing.T) {
	makeBlockChain()
	blockList := blockchain.Get(1)
	if block1.Header.Hash != blockList[0].Header.Hash {
		t.Errorf("Want %v\n but %v", block1.Header.Hash, blockList[0].Header.Hash)
	}

}

func makeBlockChain() {
	blockchain.Chain[1] = []Block{block1}
	blockchain.Chain[2] = []Block{block2}
	blockchain.Length = 2
}

func TestBlockchain_Insert(t *testing.T) {
	bc := NewBlockChain()
	makeBlockChain()
	bc.Insert(block1)
	bc.Insert(block2)

	if !reflect.DeepEqual(bc, blockchain) {
		t.Errorf("Want %v\n but %v", bc, blockchain)
	}
}

func TestBlockchain_EncodeToJson(t *testing.T) {
	makeBlockChain()

	//fmt.Println("I am here", blockchain.EncodeToJson())
}

func TestBlockchain_DecodeFromJson(t *testing.T) {
	makeBlockChain()
	fmt.Println("hahah")

}

/**

 */
func TestBlockChain(t *testing.T) {
	jsonBlockChain := "[{\"Hash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"TimeStamp\": 1234567890, \"Height\": 1, \"ParentHash\": \"genesis\", \"Size\": 1174, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}, {\"Hash\": \"24cf2c336f02ccd526a03683b522bfca8c3c19aed8a1bed1bbc23c33cd8d1159\", \"TimeStamp\": 1234567890, \"Height\": 2, \"ParentHash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"Size\": 1231, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}]"
	bc, err := DecodeFromJson(jsonBlockChain)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println("bc = ", bc)
	jsonNew, err := bc.EncodeToJson()
	//fmt.Println(jsonNew)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	var realValue []BlockJson
	var expectedValue []BlockJson

	err = json.Unmarshal([]byte(jsonNew), &realValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	err = json.Unmarshal([]byte(jsonBlockChain), &expectedValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if !reflect.DeepEqual(realValue, expectedValue) {
		fmt.Println("=========Real=========")
		fmt.Println(realValue)
		fmt.Println("=========Expcected=========")
		fmt.Println(expectedValue)
		t.Fail()
	}

}

func TestBlockChain2(t *testing.T) {
	jsonBlockChain := "[{\"Height\":1,\"TimeStamp\":1551025401,\"Hash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"ParentHash\":\"genesis\",\"Size\":2089,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}},{\"Height\":2,\"TimeStamp\":1551025401,\"Hash\":\"944eb943b05caba08e89a613097ac5ac7d373d863224d17b1958541088dc20e2\",\"ParentHash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"Size\":2146,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}},{\"Height\":2,\"TimeStamp\":1551025401,\"Hash\":\"f8af68feadf25a635bc6e81c08f81c6740bbe1fb2514c1b4c56fe1d957c7448d\",\"ParentHash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"Size\":707,\"mpt\":{\"ge\":\"Charles\"}},{\"Height\":3,\"TimeStamp\":1551025401,\"Hash\":\"f367b7f59c651e69be7e756298aad62fb82fddbfeda26cb06bfd8adf9c8aa094\",\"ParentHash\":\"f8af68feadf25a635bc6e81c08f81c6740bbe1fb2514c1b4c56fe1d957c7448d\",\"Size\":707,\"mpt\":{\"ge\":\"Charles\"}},{\"Height\":3,\"TimeStamp\":1551025401,\"Hash\":\"05ac44dd82b6cc398a5e9664add21856ae19d107d9035af5fc54c9b0ffdef336\",\"ParentHash\":\"944eb943b05caba08e89a613097ac5ac7d373d863224d17b1958541088dc20e2\",\"Size\":2146,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}}]"
	bc, err := DecodeFromJson(jsonBlockChain)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println("bc len = ", len(bc.Chain[2]))

	jsonNew, err := bc.EncodeToJson()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println("jsonBlockChain: ", jsonBlockChain)
	fmt.Println("bc       :", bc.Chain[2])

	var realValue []BlockJson
	var expectedValue []BlockJson

	err = json.Unmarshal([]byte(jsonNew), &realValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	err = json.Unmarshal([]byte(jsonBlockChain), &expectedValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if !reflect.DeepEqual(realValue, expectedValue) {
		fmt.Println("=========Real=========")
		fmt.Println(realValue)
		fmt.Println("=========Expcected=========")
		fmt.Println(expectedValue)
		t.Fail()
	}
}
