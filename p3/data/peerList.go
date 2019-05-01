package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type PeerList struct {
	selfId    int32
	peerMap   map[string]int32
	maxLength int32
	mux       sync.Mutex
}

// this struct is a helper struct to sort the peerMap
type Pair struct {
	Key   string
	Value int32
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func NewPeerList(id int32, maxLength int32) PeerList {
	return PeerList{
		selfId:    id,
		maxLength: maxLength,
		peerMap:   make(map[string]int32),
		mux:       sync.Mutex{},
	}
}

func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	peers.peerMap[addr] = id
	peers.mux.Unlock()
}

func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	delete(peers.peerMap, addr)
	peers.mux.Unlock()
}

func (peers *PeerList) Rebalance() {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	if int32(len(peers.peerMap)) <= peers.maxLength {
		return
	}

	pl := make(PairList, len(peers.peerMap)+1)
	i := 0
	for k, v := range peers.peerMap {
		pl[i] = Pair{k, v}
		i++
	}
	pl[i] = Pair{Value: peers.selfId}
	sort.Sort(pl)
	fmt.Println(pl)

	var low, high int
	for i, pair := range pl {
		if pair.Value == peers.selfId {
			low = (i - int(peers.maxLength/2) + pl.Len()) % pl.Len()
			high = (i + int(peers.maxLength/2) + pl.Len()) % pl.Len()
		}
	}

	fmt.Println("high = ", high, " low = ", low)
	if high < low {
		for i := high + 1; i < low; i++ {
			delete(peers.peerMap, pl[i].Key)
		}
	} else {
		for i := 0; i < low; i++ {
			delete(peers.peerMap, pl[i].Key)
		}
		for i := high + 1; i < int(peers.maxLength); i++ {
			delete(peers.peerMap, pl[i].Key)
		}
	}

}

func (peers *PeerList) Show() string {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	var sb strings.Builder

	for k, v := range peers.peerMap {
		curInfo := fmt.Sprint("address: ", k, ", id: ", v, "\n")
		sb.WriteString(curInfo)
	}
	return sb.String()
}

func (peers *PeerList) Register(id int32) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	peers.selfId = id
	fmt.Printf("SelfId=%v\n", id)
}

func (peers *PeerList) Copy() map[string]int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.peerMap
}

func (peers *PeerList) GetSelfId() int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.selfId
}

func (peers *PeerList) PeerMapToJson() (string, error) {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	jsonString, err := json.Marshal(peers.peerMap)
	return string(jsonString), err
}

func (peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, senderAddr string, senderId int32) {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	var peerMap map[string]int32
	err := json.Unmarshal([]byte(peerMapJsonStr), &peerMap)
	if err != nil {
		panic(err)
	}
	for k, v := range peerMap {
		if v != peers.selfId {
			peers.peerMap[k] = v
		}
	}
	//fmt.Println("addr ", senderAddr)
	peers.peerMap[senderAddr] = senderId
}

func TestPeerListRebalance() {
	peers := NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	fmt.Println(peers, " ", expected)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = NewPeerList(5, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)

	fmt.Println(peers, " ", expected)
	fmt.Println(reflect.DeepEqual(peers, expected))

	peers = NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)

	fmt.Println(peers, " ", expected)
	fmt.Println(reflect.DeepEqual(peers, expected))
}
