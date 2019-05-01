package data

import (
	"encoding/json"
	"github.com/cs686-blockchain-p3-PerrySong/p1"
	"github.com/cs686-blockchain-p3-PerrySong/p2"
	"math/rand"
)

type HeartBeatData struct {
	IfNewBlock  bool   `json:"ifNewBlock"`
	Id          int32  `json:"id"`
	BlockJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}

func NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string) HeartBeatData {
	return HeartBeatData{
		IfNewBlock:  ifNewBlock,
		Id:          id,
		BlockJson:   blockJson,
		PeerMapJson: peerMapJson,
		Addr:        addr,
		Hops:        3,
	}
}

/*
	For each HeartBeat, a node would randomly decide (this will change in Project 4) if it will
	create a new block. If so, add the block information into HeartBeatData and send the HeartBeatData to others
*/
func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapJson string, addr string) HeartBeatData {
	var newBlock p2.Block
	ifNewBlock := rand.Intn(100) < 50
	if ifNewBlock { // create a new block
		newBlock = sbc.GenBlock(p1.MerklePatriciaTrie{}) /* TODO: Empty for now */
	}

	return NewHeartBeatData(ifNewBlock, selfId, newBlock.EncodeToJson(), peerMapJson, addr)
}

func (heartBeatData *HeartBeatData) ToHeartBeatJsonData() ([]byte, error) {
	return json.Marshal(heartBeatData)
}
