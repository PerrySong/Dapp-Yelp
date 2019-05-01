package p3

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cs686-blockchain-p3-PerrySong/p1"
	"github.com/cs686-blockchain-p3-PerrySong/p2"
	"github.com/cs686-blockchain-p3-PerrySong/p3/data"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var TA_SERVER = "http://localhost:6688"
var REGISTER_SERVER = TA_SERVER + "/peer"
var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var SELF_ADDR string

var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool
var port int
var firstNodePort = 6688

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.

	var portString string
	SBC = data.NewBlockChain()
	ifStarted = false
	if len(os.Args) <= 1 || os.Args[1] == "-test.v" {
		port = 6688
		portString = "6688"
		Download()
	} else {
		var err error
		portString = os.Args[1]
		port, err = strconv.Atoi(portString)
		if err != nil {
			panic(err)
		}
	}
	Peers = data.NewPeerList(int32(port), 32)
	SELF_ADDR = "http://localhost:" + portString

	log.Println(portString)
	if portString != "6688" {
		Peers.Add(TA_SERVER, 6688)
		log.Println(Peers)
	}
}

// Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {

	ifStarted = true
	Download()
	go StartHeartBeat()
	go StartTryingNonces()
	fmt.Fprint(w, `{"start": true}`)
}

// Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
}

// Register to TA's server, get an ID
// TODO: Not implemented for now
func Register() {
	// TODO: Hard code for now
	//Peers.Register(1)
}

// Download blockchain from TA server
func Download() {
	// TODO: Hard code for now
	//selfId := Peers.GetSelfId()
	if port == firstNodePort { // Generate head block
		SBC.Insert(p2.JsonToBlock(`{
    		"Hash":"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48",
    		"TimeStamp":1234567890,
    		"Height":1,
    		"ParentHash":"genesis",
    		"Size":1174,
   		 	"mpt":{
				"charles":"ge",	
        		"hello":"world"
    			}
			}`))
		return
	}
	resp, err := http.Get(BC_DOWNLOAD_SERVER)
	if err != nil {
		log.Println(err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	SBC.UpdateEntireBlockChain(string(bodyBytes))
	if port != firstNodePort {
		Peers.Add(TA_SERVER, int32(firstNodePort))
	}
}

/*
	Return the BlockChain's JSON. And add the remote peer into the PeerMap.
	UploadBlock(): Return the Block's JSON.
*/

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		// TODO:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	fmt.Fprint(w, blockChainJson)
}

/*
	If you have the block, return the JSON string of the specific block;
	if you don't have the block,
 	return HTTP 204: StatusNoContent;
	if there's an error, return HTTP 500: InternalServerError.
*/
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	sList := strings.Split(path, "/")

	if len(sList) != 4 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	height, err := strconv.Atoi(sList[2])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	hash := sList[3]
	//blockList, hasHeight := SBC.Get(int32(height))
	//if !hasHeight {
	//	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
	//}

	block, hasBlock := SBC.GetBlock(int32(height), hash)
	if !hasBlock {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(block.EncodeToJson()))

}

// Received a heartbeat
/*
	Add the remote address, and the PeerMapJSON into local PeerMap. Then check
	if the HeartBeatData contains a new block. If so, do these: (1) check if
	the parent block exists. If not, call AskForBlock() to download the parent block.
	(2) insert the new block from HeartBeatData. (3) HeartBeatData.hops minus one,
	and if it's still bigger than 0, call ForwardHeartBeat() to forward this heartBeat
	to all peers.
*/

// Alter this function so that when it receives a HeartBeatData with a new block, it verifies the nonce as described above.
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var heartBeat data.HeartBeatData
	err := decoder.Decode(&heartBeat)
	if err != nil {
		//fmt.Println("decode heartbeat fail: heartbeat = ", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Inject peerMap
	senderAddress := heartBeat.Addr
	senderId := heartBeat.Id
	// Terminate if the sender id equals to the self id
	if senderId == Peers.GetSelfId() {
		fmt.Println("My self block")
		return
	}

	Peers.InjectPeerMapJson(heartBeat.PeerMapJson, senderAddress, senderId)

	if !heartBeat.IfNewBlock {
		fmt.Println("Not new block ", err)
		return
	}

	// (1) check if the parent block exists. If not, call AskForBlock() to download the parent block.

	// (2) insert the new block from HeartBeatData.
	var bJson p2.BlockJson
	err = json.Unmarshal([]byte(heartBeat.BlockJson), &bJson)
	if err != nil {
		fmt.Println("fail to unmarshal ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newBlock := bJson.ToBlock()
	// verify the proof of work
	if !isValidBlock(newBlock) {
		fmt.Println("Block is not valid ", newBlock)
		http.Error(w, "incoming block is not valid", http.StatusBadRequest)
		return
	}
	SBC.Insert(newBlock)

	// (3) HeartBeatData.hops minus one,
	//	and if it's still bigger than 0, call ForwardHeartBeat() to forward this heartBeat
	//	to all peers.
	if !SBC.CheckParentHash(newBlock) {
		//log.Println("Checking parenthash: ", bJson.ParentHash, " height = ", bJson.Height-1)
		AskForBlock(bJson.Height-1, bJson.ParentHash)
	}

	heartBeat.Hops = heartBeat.Hops - 1
	if heartBeat.Hops > 0 {
		//log.Println("Rendering heartbeat hop = ", heartBeat.Hops)
		ForwardHeartBeat(heartBeat)
	}

}

// Ask another server to return a block of certain height and hash
// Loop through all peers in local PeerMap to download a block.
// As soon as one peer returns the block, stop the loop.
// Update the function to recursively ask for all block
func AskForBlock(height int32, hash string) {

	if height < 1 || SBC.HasBlock(height, hash) {
		return
	}
	peerMap := Peers.Copy()
	myId := Peers.GetSelfId()
	for address, _ := range peerMap {
		req := fmt.Sprint(address, "/block", myId, "/", height, "/", hash)
		resp, err := http.Get(req)
		if err != nil {
			// TODO: Handle error
			continue
		} else {
			if resp.StatusCode == http.StatusOK {
				bodyBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					// TODO: handle error
				}
				bodyString := string(bodyBytes)
				block := p2.JsonToBlock(bodyString)
				SBC.Insert(block)
				// Recursively fetch parent's block
				AskForBlock(block.Header.Height-1, block.Header.ParentHash)
				return
			}
		}
		resp.Body.Close()
	}
}

//TODO
func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
	peerMap := Peers.Copy()

	for address, _ := range peerMap {
		body, err := heartBeatData.ToHeartBeatJsonData()
		if err != nil {
			// Handle error
		}

		_, err = http.Post(address+"/heartbeat/receive", "J", bytes.NewBuffer(body))
		if err != nil {
			// Handle error
		}
	}
}

//
func StartHeartBeat() {
	for ifStarted {
		time.Sleep(2 * time.Second)
		peersMapJson, _ := Peers.PeerMapToJson()
		heartBeatData := data.NewHeartBeatData(true, Peers.GetSelfId(), "", peersMapJson, SELF_ADDR)
		ForwardHeartBeat(heartBeatData)
	}
}

/**
	This function starts a new thread that tries different nonces to generate new blocks. Nonce is a string of 16 hexes such as
	"1f7b169c846f218a". Initialize the rand when you start a new node with something unique about each node, such as the current
	time or the port number. Here's the workflow of generating blocks:

    (1) Start a while loop.
    (2) Get the latest block or one of the latest blocks to use as a parent block.
    (3) Create an MPT.
    (4) Randomly generate the first nonce, verify it with simple PoW algorithm to see
	if SHA3(parentHash + nonce + mptRootHash) starts with 10 0's (or the number you modified into).
	Since we use one laptop to try different nonces, six to seven 0's could be enough. If the nonce failed the verification,
	increment it by 1 and try the next nonce.
    (6) If a nonce is found and the next block is generated, forward that block to all peers with an HeartBeatData;
    (7) If someone else found a nonce first, and you received the new block through your function ReceiveHeartBeat(),
	stop trying nonce on the current block, continue to the while loop by jumping to the step(2).
*/
func StartTryingNonces() {
	tryNonces()
}

func tryNonces() {
	nonceByte := newNonce()
	// Build dummy mpt
	mpt := p1.MerklePatriciaTrie{}
	mpt.Insert("id", string(Peers.GetSelfId()))
	rootHash := mpt.GetRoot()
	for {
		parentBlock := SBC.GetLatestBlock()
		parentHash := parentBlock[0].Header.Hash
		nonce := string(nonceByte)

		proof := parentHash + nonce + rootHash

		sum := sha3.Sum256([]byte(proof))
		hash := hex.EncodeToString(sum[:])

		if validProof(hash) {
			fmt.Println(nonce)
			block := p2.NewBlock(parentBlock[0].Header.Height+1, parentHash, nonce, mpt)
			blockJson := block.EncodeToJson()
			pMapJson, err := Peers.PeerMapToJson()
			if err != nil {
				continue
			}
			heartBeatData := data.NewHeartBeatData(true, Peers.GetSelfId(), blockJson, pMapJson, SELF_ADDR)
			SBC.Insert(*block)
			ForwardHeartBeat(heartBeatData)
		}
		nonceByte = incrementNonce(nonceByte)
	}
}

func validProof(hash string) bool {
	return strings.HasPrefix(hash, "00000")
}

func isValidBlock(block p2.Block) bool {
	parentHash := block.Header.ParentHash
	nonce := block.Header.Nonce
	rootHash := block.GetMPTRoot()

	proof := parentHash + string(nonce) + rootHash
	sum := sha3.Sum256([]byte(proof))

	hash := hex.EncodeToString(sum[:])
	return validProof(hash)
}

func incrementNonce(nonce []byte) []byte {
	carry := byte(1)
	for i := len(nonce) - 1; i >= 0; i-- {
		num := carry + nonce[i]
		nonce[i] = num % 16
		carry = num / 16
		if carry == 0 {
			return nonce
		}
	}
	nonce = make([]byte, 16)
	return nonce
}

func newNonce() []byte {
	res := make([]byte, 16)
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	for i := 0; i < 16; i++ {
		res[i] = byte(random.Int() % 16)
	}
	return res
}

/**
	This function prints the current canonical chain, and chains of all forks if there are forks.
	Note that all forks should end at the same height (otherwise there wouldn't be a fork).
  	Example of the output of Canonical() function: You can have a different format,
	but it should be clean and clear.
*/

func Canonical(w http.ResponseWriter, r *http.Request) {
	blockList := SBC.GetLatestBlock()
	canonicalStr := ""
	for id, block := range blockList {
		chain, hasChain := traverseChain(block)
		if hasChain {
			canonicalStr += "Blockchain " + strconv.Itoa(id) + "\n" + chain
		}
	}
	fmt.Fprintf(w, "%s\n", canonicalStr)
}

func traverseChain(block p2.Block) (string, bool) {
	res := ""
	for block.Header.Height > 1 {
		res += block.Digest()
		prevBlock, hasParent := SBC.GetParentBlock(block)
		fmt.Println(prevBlock)
		block = prevBlock
		if !hasParent {
			return res, false
		}
	}
	res += block.Digest()
	return res, true
}
