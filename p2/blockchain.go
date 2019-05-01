package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"sort"
)

type BlockChain struct {
	Chain  map[int32][]Block
	Length int32 // Length equals to the highest block height.
}

type BlockchainJson struct {
	blockList []Block
}

func NewBlockChain() BlockChain {
	var blockchain BlockChain
	blockchain.Chain = make(map[int32][]Block)
	blockchain.Length = 0
	return blockchain
}

/**
Description: This function takes a Height as the argument,
returns the list of blocks stored in that Height or None if the Height doesn't exist.
Argument: int32
Return type: []Block
*/
func (blockchain *BlockChain) Get(height int32) []Block {
	if blockchain.Chain == nil {
		return nil
	}
	if len(blockchain.Chain[height]) == 0 {
		return nil
	}
	return blockchain.Chain[height]
}

/**
Description: This function takes a block as the argument, use its Height to find the corresponding list in blockchain's Chain map.
If the list has already contained that block's Hash, ignore it because we don't store duplicate blocks; if not, insert the block into the list.
Argument: block
*/
func (blockchain *BlockChain) Insert(block Block) {

	for _, curblock := range blockchain.Chain[block.Header.Height] {
		if reflect.DeepEqual(curblock.Header.Hash, block.Header.Hash) {
			return
		}
	}

	if block.Header.Height > blockchain.Length {
		blockchain.Length = block.Header.Height
	}
	blockchain.Chain[block.Header.Height] = append(blockchain.Chain[block.Header.Height], block)
}

/**

 */
func (blockchain *BlockChain) EncodeToJson() (string, error) {
	var blockchainJson []BlockJson

	for _, curBlockList := range blockchain.Chain {
		for _, block := range curBlockList {
			blockJson := block.ToBlockJson()
			blockchainJson = append(blockchainJson, blockJson)
		}
	}
	res, err := json.MarshalIndent(blockchainJson, "", "	")
	if err != nil {
		panic(err)
	}

	return string(res), err
}

/**
Description: This function is called upon a blockchain instance. It takes a blockchain JSON string as input, decodes the JSON string back to a list of block JSON strings, decodes each block JSON string back to a block instance, and inserts every block into the blockchain.
Argument: self, string
*/
func DecodeFromJson(jsonStr string) (BlockChain, error) {
	blockChain := NewBlockChain()
	blockList := make([]BlockJson, 0)
	jsonBytesArr := []byte(jsonStr)
	//fmt.Println("Json = ", jsonStr)

	err := json.Unmarshal(jsonBytesArr, &blockList)

	//fmt.Println("blockList = ", blockList)
	for _, blockJson := range blockList {
		block := blockJson.ToBlock()
		blockChain.Insert(block)
	}

	return blockChain, err
}

func (bc *BlockChain) Show() string {
	rs := ""
	var idList []int
	for id := range bc.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range bc.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash+"<="+block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha3.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}

/**
This function returns the list of blocks of height "BlockChain.length".
*/

func (bc *BlockChain) GetLatestBlocks() []Block {
	return bc.Get(bc.Length)
}

/**
This function takes a block as the parameter, and returns its parent block.
*/
func (bc *BlockChain) GetParentBlock(curBlock Block) (Block, bool) {
	parentHash := curBlock.Header.ParentHash
	parentHeight := curBlock.Header.Height - 1

	prevBlockList := bc.Get(parentHeight)

	for _, block := range prevBlockList {
		if reflect.DeepEqual(block.Header.Hash, parentHash) {
			return block, true
		}
	}

	return Block{}, false
}

func (bc *BlockChain) HasBlock(height int32, hash string) bool {
	if height > bc.Length {
		return false
	}
	blockList := bc.Chain[height]
	for _, block := range blockList {
		if hash == block.Header.Hash {
			return true
		}
	}
	return false
}
