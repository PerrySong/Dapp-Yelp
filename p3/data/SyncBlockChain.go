package data

import (
	"github.com/cs686-blockchain-p3-PerrySong/p1"
	"github.com/cs686-blockchain-p3-PerrySong/p2"
	"sync"
)

type SyncBlockChain struct {
	bc  p2.BlockChain
	mux sync.Mutex
}

func NewBlockChain() SyncBlockChain {
	return SyncBlockChain{
		bc:  p2.NewBlockChain(),
		mux: sync.Mutex{},
	}
}

func NewDummyBlockChain() SyncBlockChain {

	jsonBlockChain := "[{\"Hash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"TimeStamp\": 1234567890, \"Height\": 1, \"ParentHash\": \"genesis\", \"Size\": 1174, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}, {\"Hash\": \"24cf2c336f02ccd526a03683b522bfca8c3c19aed8a1bed1bbc23c33cd8d1159\", \"TimeStamp\": 1234567890, \"Height\": 2, \"ParentHash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"Size\": 1231, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}]"
	bc, _ := p2.DecodeFromJson(jsonBlockChain)
	return SyncBlockChain{
		bc:  bc,
		mux: sync.Mutex{},
	}
}

func (sbc *SyncBlockChain) Get(height int32) ([]p2.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(height) /* TODO change the return bool */, true
}

// TODO: what does hash mean?
func (sbc *SyncBlockChain) GetBlock(height int32, hash string) (p2.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	blockList := sbc.bc.Get(height)
	// TODO
	for _, block := range blockList {
		if block.Header.Hash == hash {
			return block, true
		}
	}
	return p2.Block{}, false
}

func (sbc *SyncBlockChain) Insert(block p2.Block) {
	sbc.mux.Lock()
	sbc.bc.Insert(block)
	sbc.mux.Unlock()
}

/**
CheckParentHash() is used to check if the parent hash(or parent block)
exist in the current blockchain when you want to insert a new block sent by others.
*/
func (sbc *SyncBlockChain) CheckParentHash(insertBlock p2.Block) bool {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	parentHash := insertBlock.Header.ParentHash
	height := insertBlock.Header.Height - 1
	blockList := sbc.bc.Get(height)
	for _, block := range blockList {
		if parentHash == block.Header.Hash {
			return true
		}
	}
	return false
}

func (sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	blockchain, err := p2.DecodeFromJson(blockChainJson)
	if err != nil {
		// How to handle the error?
		panic(err)
	}
	sbc.bc = blockchain

}

func (sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.EncodeToJson()
}

/*
NewBlock(height int32, previousHash string, value p1.MerklePatriciaTrie)
*/
func (sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie) p2.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	prevHash := sbc.bc.Get(sbc.bc.Length)[0].Header.Hash // TODO: How to choose parent
	return *p2.NewBlock(sbc.bc.Length+1, prevHash, "...", mpt)
}

func (sbc *SyncBlockChain) Show() string {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Show()
}

/**
This function returns the list of blocks of height "BlockChain.length".
*/
func (sbc *SyncBlockChain) GetLatestBlock() []p2.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetLatestBlocks()
}

/**
This function takes a block as the parameter, and returns its parent block.
*/
func (sbc *SyncBlockChain) GetParentBlock(block p2.Block) (p2.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetParentBlock(block)
}

func (sbc *SyncBlockChain) HasBlock(height int32, hash string) bool {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.HasBlock(height, hash)
}

func (sbc *SyncBlockChain) GetLength() int32 {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Length
}
