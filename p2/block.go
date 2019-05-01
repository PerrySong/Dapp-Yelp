package p2

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cs686-blockchain-p3-PerrySong/p1"
	"golang.org/x/crypto/sha3"
	"time"
)

type Header struct {
	Nonce      string
	Height     int32
	TimeStamp  int64
	Hash       string
	ParentHash string
	Size       int32
}

type Block struct {
	Header Header
	value  p1.MerklePatriciaTrie
}

type BlockJson struct {
	Nonce      string            `json:"Nonce"`
	Height     int32             `json:"Height"`
	Timestamp  int64             `json:"TimeStamp"`
	Hash       string            `json:"Hash"`
	ParentHash string            `json:"ParentHash"`
	Size       int32             `json:"Size"`
	MPT        map[string]string `json:"mpt"`
}

func NewBlock(height int32, previousHash string, nonce string, value p1.MerklePatriciaTrie) *Block {
	var header Header
	var block Block

	block.value = value

	header.Height = height
	header.TimeStamp = int64(time.Now().Unix())

	hashStr := string(header.Height) + string(header.TimeStamp) + header.ParentHash + value.GetRoot() + string(header.Size)
	sum := sha3.Sum256([]byte(hashStr))
	header.Hash = hex.EncodeToString(sum[:])

	header.ParentHash = previousHash
	header.Size = int32(len(encodeToBytesArr(block)))
	header.Nonce = nonce
	block.Header = header

	return &block
}

/**
This function update the current block's Hash value, when block insert to the blockchain, it's ParentHash will change, hence, it should change its value
*/
func (block *Block) updateHash() {
	block.Header.Hash = string(block.Header.Height) + string(block.Header.TimeStamp) + block.Header.ParentHash + block.value.GetRoot() + string(block.Header.Size)
}

func encodeToBytesArr(block Block) []byte {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	if err := enc.Encode(block.value); err != nil {
		fmt.Println(err)
	}
	return buffer.Bytes()
}

func JsonToBlock(jsonString string) Block {
	var block Block
	var blockJson BlockJson
	jsonByteArr := []byte(jsonString)
	err := json.Unmarshal(jsonByteArr, &blockJson)

	if err != nil {
		fmt.Println("Unable to unmarshal jsonString")
		panic(err)
	}

	mpt := mapToMerkleTree(blockJson.MPT)
	header := Header{Height: blockJson.Height, Hash: blockJson.Hash, TimeStamp: blockJson.Timestamp, ParentHash: blockJson.ParentHash}
	block = Block{Header: header, value: *mpt}
	//block = *NewBlock(blockJson.Height, blockJson.ParentHash, *mpt)

	//fmt.Println(block)
	return block
}

func mapToMerkleTree(mptMap map[string]string) *p1.MerklePatriciaTrie {
	mpt := new(p1.MerklePatriciaTrie)
	for key, val := range mptMap {
		mpt.Insert(key, val)
	}
	return mpt
}

func (block *Block) ToBlockJson() BlockJson {
	blockJson := BlockJson{
		Nonce:      block.Header.Nonce,
		Height:     block.Header.Height,
		Hash:       block.Header.Hash,
		Size:       block.Header.Size,
		ParentHash: block.Header.ParentHash,
		Timestamp:  block.Header.TimeStamp,
		MPT:        block.value.ToMap(),
	}
	return blockJson
}

func (block *Block) EncodeToJson() string {
	blockJson := block.ToBlockJson()

	res, err := json.Marshal(blockJson)
	if err != nil {
		fmt.Println("Cannot convert block to json string")
	}
	return string(res)
}

/**

 */
func (blockJson *BlockJson) ToBlock() Block {
	header := Header{Nonce: blockJson.Nonce, Height: blockJson.Height, Hash: blockJson.Hash, Size: blockJson.Size, ParentHash: blockJson.ParentHash, TimeStamp: blockJson.Timestamp}
	mpt := mapToMerkleTree(blockJson.MPT)
	block := Block{Header: header, value: *mpt}
	return block
}

/**
Return a digest for a block's summary
*/
func (block Block) Digest() string {
	res := fmt.Sprintf("height=%d, timeStamp=%d, hash=%s, parentHash=%s, size=%d\n", block.Header.Height, block.Header.TimeStamp, block.Header.Hash, block.Header.ParentHash, block.Header.Size)
	return res
}

func (block Block) GetMPTRoot() string {
	return block.value.GetRoot()
}
