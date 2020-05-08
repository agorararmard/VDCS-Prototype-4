package vdcs

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"./elgamal"
)

//Wire wire abstraction
type Wire struct {
	WireID    []byte `json:"WireID"`
	WireLabel []byte `json:"WireLabel"`
}

//Gate gate abstraction
type Gate struct {
	GateID     []byte   `json:"GateID"`
	GateInputs [][]byte `json:"GateInputs"`
}

//CircuitGate a gate in a boolean circuit
type CircuitGate struct {
	Gate
	TruthTable []bool `json:"TruthTable"`
}

//GarbledGate a gate in a garbled circuit
type GarbledGate struct {
	Gate
	KeyY            [][]byte `json:"EncryptedY"`
	GarbledValuesC1 [][]byte `json:"GarbledValuesC1"`
	GarbledValuesC2 [][]byte `json:"GarbledValuesC2"`

	GarbledValues [][]byte `json:"GarbledValues"` //to be removed later.
}

//ComID computation ID abstraction
type ComID struct {
	CID []byte `json:"ComID"`
}

//Circuit circuit abstraction
type Circuit struct {
	InputGates  []CircuitGate `json:"CircuitInputGates"`
	MiddleGates []CircuitGate `json:"CircuitMiddleGates"`
	OutputGates []CircuitGate `json:"CircuitOutputGates"`
}

//Randomness container for randomness
type Randomness struct {
	Rin       int64 `json:"Rin"`
	Rout      int64 `json:"Rout"`
	Rmask     int64 `json:"Rmask"`
	Rgc       int64 `json:"Rgc"`
	LblLength int   `json:"LblLength"`
}

//CircuitMessage a complete circuit message
type CircuitMessage struct {
	Circuit
	ComID
	Randomness
}

//GarbledCircuit garbled circuit abstraction
type GarbledCircuit struct {
	InputGates  []GarbledGate `json:"GarbledInputGates"`
	MiddleGates []GarbledGate `json:"GarbledMiddleGates"`
	OutputGates []GarbledGate `json:"GarbledOutputGates"`
	ComID
}

//GarbledMessage complete garbled circuit message
type GarbledMessage struct {
	InputWires []Wire `json:"GarbledInputWires"`
	GarbledCircuit
	OutputWires []Wire `json:"GarbledOutputWires"`
}

//ResEval evaluation result abstraction
type ResEval struct {
	Res [][]byte `json:"Result"`
	ComID
}

//PartyInfo container for general information about a node
type PartyInfo struct {
	IP        []byte `json:"IP"`
	Port      int    `json:"Port"`
	PublicKey []byte `json:"PublicKey"`
	UserName  []byte `json:"UserName"`
}

//MyInfo container for general and private information about a node
type MyInfo struct {
	PartyInfo
	CleosKey   []byte `json:"CleosKey"`
	PrivateKey []byte `json:"PrivateKey"`
}

//ServerCapabilities server capabilities abstraction
type ServerCapabilities struct {
	NumberOfGates int     `json:"NumberOfGates"`
	FeePerGate    float64 `json:"FeePerGate"`
}

//Token a token container for the ease of message passing
type Token struct {
	TokenGen []byte `json:"TokenGen"`
}

//ServerInfo container for server relevant info in Directory of Service
type ServerInfo struct {
	PartyInfo
	ServerCapabilities
}

//ClientInfo container for client relevant info in Directory of Service
type ClientInfo struct {
	PartyInfo
}

//RegisterationMessage a complete registration message
type RegisterationMessage struct {
	Type   []byte     `json:"Type"` //Server, Client
	Server ServerInfo `json:"ServerInfo"`
}

//FunctionInfo a container for function requirements
type FunctionInfo struct {
	Token
	NumberOfServers    int `json:"NumberOfServers"`
	ServerCapabilities     //in this case we describe the capabilities needed to compute the circuit
}

//CycleRequestMessage Wrapping In case we needed to add new request types for failure handling
type CycleRequestMessage struct {
	FunctionInfo
}

//Cycle cycle wrapper
type Cycle struct {
	ServersCycle []PartyInfo `json:"ServersCycle"`
}

//CycleMessage a complete cycle message reply
type CycleMessage struct {
	Cycle
	TotalFee int `json:"TotalFee"`
}

//Message passed through cycle
type Message struct {
	Type []byte `json:"Type"` //Garble, Rerand, Eval
	Circuit
	GarbledMessage
	InputWires []Wire `json:"GeneralInputWires"`
	Randomness
	ComID
	NextServer PartyInfo `json:"NextServer"`
}

//MessageArray container of messages
type MessageArray struct {
	Array []Message `json:"Array"`
	Keys  [][]byte  `json:"Keys"`
}

//ChannelContainer contains what is passed through message channels within the client code
type ChannelContainer struct {
	InputWires  []Wire `json:"InputWires"`
	OutputWires []Wire `json:"OutputWires"`
	PartyInfo
	Keys [][]byte `json:"Keys"`
}

//local Gate gate abstraction
type localgate struct {
	GateID     string   `json:"GateID"`
	GateInputs []string `json:"GateInputs"`
}

type localcircuitgate struct {
	localgate
	TruthTable []bool `json:"TruthTable"`
}
type localcircuit struct {
	InputGates  []localcircuitgate `json:"InputGates"`
	MiddleGates []localcircuitgate `json:"MiddleGates"`
	OutputGates []localcircuitgate `json:"OutputGates"`
}

//DirctoryInfo Global Variable to store Directory communication info
var DirctoryInfo = struct {
	Port int
	IP   []byte
}{
	Port: 0,
	IP:   []byte(""),
}

//DecentralizedDirectoryInfo Global Variable to store Directory communication info
var DecentralizedDirectoryInfo = struct {
	URL            string
	ActionAccount  string
	PasswordWallet string
}{
	URL:            "",
	ActionAccount:  "",
	PasswordWallet: "",
}

//MyOwnInfo personal info container
var MyOwnInfo MyInfo

//MyToken holds directory sent token
var MyToken Token

//Decentralization indicates whether the central or decentralized directory of service is used
var Decentralization bool

//ReadyFlag is a simulation for channels between the post handler and the eval function
var ReadyFlag bool

//ReadyMutex is a simulation for channels between the post handler and the eval function
var ReadyMutex = sync.RWMutex{}

//MyResult is a simulation for channels between the post handler and the eval function
var MyResult ResEval

// mutex for local com
var cMutex = sync.RWMutex{}

//SetMyInfo sets the info of the current node
func SetMyInfo(username string, cleosKey string) {
	pI, sk := GetPartyInfo(username)
	MyOwnInfo = MyInfo{
		PartyInfo:  pI,
		CleosKey:   []byte(cleosKey),
		PrivateKey: sk,
	}
}

//SetDirectoryInfo to set the dircotry info
func SetDirectoryInfo(ip []byte, port int) {
	DirctoryInfo.Port = port
	DirctoryInfo.IP = ip
}

//SetDecentralizedDirectoryInfo to set the decentralized dircotry info
func SetDecentralizedDirectoryInfo(url string, actionAccount string, passwordWallet string) {
	DecentralizedDirectoryInfo.URL = url
	DecentralizedDirectoryInfo.ActionAccount = actionAccount
	DecentralizedDirectoryInfo.PasswordWallet = passwordWallet
}

//GetCircuitSize get the number of gates in a circuit
func GetCircuitSize(circ Circuit) int {
	return len(circ.InputGates) + len(circ.MiddleGates) + len(circ.OutputGates)
}

//GetInputSizeOutputSize returns the number of inputs and outputs of a given circuit
func GetInputSizeOutputSize(circ Circuit) (inputSize int, outputSize int) {
	inputSize = len(circ.InputGates) * 2
	outputSize = len(circ.OutputGates)
	return
}

//convertLocalToGlobal converts local context circuits into global context
func convertLocalToGlobal(lc localcircuit) (c Circuit) {
	for _, val := range lc.InputGates {
		tmp := CircuitGate{
			Gate: Gate{
				GateID: []byte(val.GateID),
			},
			TruthTable: val.TruthTable,
		}

		for _, val2 := range val.GateInputs {
			tmp.GateInputs = append(tmp.GateInputs, []byte(val2))
		}
		c.InputGates = append(c.InputGates, tmp)
	}

	for _, val := range lc.MiddleGates {
		tmp := CircuitGate{
			Gate: Gate{
				GateID: []byte(val.GateID),
			},
			TruthTable: val.TruthTable,
		}

		for _, val2 := range val.GateInputs {
			tmp.GateInputs = append(tmp.GateInputs, []byte(val2))
		}
		c.MiddleGates = append(c.MiddleGates, tmp)
	}
	for _, val := range lc.OutputGates {
		tmp := CircuitGate{
			Gate: Gate{
				GateID: []byte(val.GateID),
			},
			TruthTable: val.TruthTable,
		}
		for _, val2 := range val.GateInputs {
			tmp.GateInputs = append(tmp.GateInputs, []byte(val2))
		}
		c.OutputGates = append(c.OutputGates, tmp)
	}
	return
}

//ClientRegister registers a client to directory of service
func ClientRegister() {
	SetMyInfo("", "")
	regMsg := RegisterationMessage{
		Type: []byte("Client"),
		Server: ServerInfo{
			PartyInfo: MyOwnInfo.PartyInfo,
			ServerCapabilities: ServerCapabilities{
				NumberOfGates: 0,
				FeePerGate:    0,
			},
		},
	}
	for !SendToDirectory(regMsg, DirctoryInfo.IP, DirctoryInfo.Port) {
	}
}

//ClientRegisterDecentralized registers a client to the decentralized directory of service
func ClientRegisterDecentralized(username string, cleosKey string) {
	SetMyInfo(username, cleosKey)
	regMsg := RegisterationMessage{
		Type: []byte("Client"),
		Server: ServerInfo{
			PartyInfo: MyOwnInfo.PartyInfo,
			ServerCapabilities: ServerCapabilities{
				NumberOfGates: 0,
				FeePerGate:    0,
			},
		},
	}
	err := UnlockWallet(DecentralizedDirectoryInfo.URL, DecentralizedDirectoryInfo.PasswordWallet)
	CreateAccount(DecentralizedDirectoryInfo.URL, regMsg)
	err = RegisterOnDecentralizedDS(DecentralizedDirectoryInfo.URL, DecentralizedDirectoryInfo.ActionAccount, regMsg)
	if err != nil {
		panic(err)
	}
}

//SolveToken recieves a token challenge and solves it
func SolveToken(token Token) Token {
	ans, err := RSAPrivateDecrypt(RSAPrivateKeyFromBytes(MyOwnInfo.PrivateKey), token.TokenGen)
	if err != nil {
		panic("Wrong Token!")
	}
	return Token{TokenGen: ans}
}

//GetHandlerClient recieves a token challenge and solves it
func GetHandlerClient(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var x Token
		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error reading", err)
		}
		err = json.Unmarshal(jsn, &x)
		if err != nil {
			log.Fatal("bad decode", err)
		}
		ret := SolveToken(x)
		MyToken = ret
		responseJSON, err := json.Marshal(ret)
		if err != nil {
			fmt.Fprintf(w, "error %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)

	}
}

//PostHandlerClient recieves the result of evaluation
func PostHandlerClient(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var x ResEval
		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error reading", err)
		}
		err = json.Unmarshal(jsn, &x)
		if err != nil {
			log.Fatal("bad decode", err)
		}
		ReadyMutex.Lock()
		ReadyFlag = true
		MyResult = x
		ReadyMutex.Unlock()
		//Pass the result to the interested eval function
	}
}

//ClientHTTP Client listeners
func ClientHTTP() {
	http.HandleFunc("/post", PostHandlerClient)
	http.HandleFunc("/get", GetHandlerClient)
	http.ListenAndServe(":"+strconv.Itoa(MyOwnInfo.Port), nil)
}

//Comm basically, the channel will need to send the input/output mapping as well
func Comm(cir string, cID int64, numberOfServers int, feePerGate float64, chVDCSCommCircRes chan<- ChannelContainer) {

	file, _ := ioutil.ReadFile(cir + ".json")
	localmCirc := localcircuit{}
	err := json.Unmarshal([]byte(file), &localmCirc) //POSSIBLE BUG
	if err != nil {
		log.Fatal(err)
	}
	//	rand.Seed(int64(cID))

	mCirc := convertLocalToGlobal(localmCirc)

	circuitSize := GetCircuitSize(mCirc)
	cycleRequestMessage := CycleRequestMessage{
		FunctionInfo{
			Token:           MyToken,
			NumberOfServers: numberOfServers,
			ServerCapabilities: ServerCapabilities{
				NumberOfGates: circuitSize,
				FeePerGate:    feePerGate,
			},
		},
	}
	cycleMessage := CycleMessage{}
	ok := false
	if Decentralization == true {
		ok, cycleMessage = FetchCycleDecentralized(DecentralizedDirectoryInfo.URL, DecentralizedDirectoryInfo.ActionAccount, cycleRequestMessage)
	} else {
		cycleMessage, ok = GetFromDirectory(cycleRequestMessage, DirctoryInfo.IP, DirctoryInfo.Port)
		for ok == false {
			cycleMessage, ok = GetFromDirectory(cycleRequestMessage, DirctoryInfo.IP, DirctoryInfo.Port)
		}
	}

	msgArray, randNess, keys := GenerateMessageArray(cycleMessage, cID, mCirc)

	//fmt.Println(cycleMessage)
	//fmt.Println(keys) //store the keys somewhere for recovery or pass on channel

	ipS1 := cycleMessage.ServersCycle[0].IP
	portS1 := cycleMessage.ServersCycle[0].Port

	for !SendToServer(msgArray, ipS1, portS1) {

	}

	//Generate input wires
	arrIn, arrOut := GenerateInputWiresValidate(mCirc, randNess, cID)

	//Send Circuit to channel
	var cc ChannelContainer
	for _, val := range arrIn {
		cc.InputWires = append(cc.InputWires, Wire{WireLabel: val})
	}
	for _, val := range arrOut {
		cc.OutputWires = append(cc.OutputWires, Wire{WireLabel: val})
	}
	cc.PartyInfo = cycleMessage.ServersCycle[numberOfServers-1]
	cc.Keys = keys
	chVDCSCommCircRes <- cc

}

//GenerateMessageArray Takes a CycleMessage, a cID, and a circuit and creates a message array encrypted and returns it with the corresponding randomness for the user to use
func GenerateMessageArray(cycleMessage CycleMessage, cID int64, circ Circuit) (mArr MessageArray, rArr []Randomness, keys [][]byte) {
	cMutex.Lock()
	numberOfServers := len(cycleMessage.ServersCycle)

	rArr = GenerateRandomness(numberOfServers, cID)

	message := Message{
		Type:       []byte("Garble"),
		Circuit:    circ,
		Randomness: rArr[0],
		ComID:      ComID{CID: []byte(strconv.FormatInt(cID, 10))},
		NextServer: cycleMessage.ServersCycle[1],
	}
	k1 := RandomSymmKeyGen()
	messageEnc := EncryptMessageAES(k1, message)

	keys = append(keys, k1)

	k1, err := RSAPublicEncrypt(RSAPublicKeyFromBytes(cycleMessage.ServersCycle[0].PublicKey), k1)
	if err != nil {
		panic("Invalid PublicKey")
	}
	mArr = MessageArray{
		Array: append(mArr.Array, messageEnc),
		Keys:  append(mArr.Keys, k1),
	}

	for i := 1; i < numberOfServers-1; i++ {

		message = Message{
			Type:       []byte("ReRand"),
			Randomness: rArr[i],
			ComID:      ComID{CID: []byte(strconv.FormatInt(cID, 10))},
			NextServer: cycleMessage.ServersCycle[i+1],
		}

		k1 = RandomSymmKeyGen()
		messageEnc = EncryptMessageAES(k1, message)

		keys = append(keys, k1)

		k1, err = RSAPublicEncrypt(RSAPublicKeyFromBytes(cycleMessage.ServersCycle[i].PublicKey), k1)
		if err != nil {
			panic("Invalid PublicKey")
		}
		mArr = MessageArray{
			Array: append(mArr.Array, messageEnc),
			Keys:  append(mArr.Keys, k1),
		}

	}

	message = Message{
		Type:       []byte("SEval"),
		ComID:      ComID{CID: []byte(strconv.FormatInt(cID, 10))},
		NextServer: MyOwnInfo.PartyInfo,
	}
	k1 = RandomSymmKeyGen()
	messageEnc = EncryptMessageAES(k1, message)

	keys = append(keys, k1)

	k1, err = RSAPublicEncrypt(RSAPublicKeyFromBytes(cycleMessage.ServersCycle[numberOfServers-1].PublicKey), k1)
	if err != nil {
		panic("Invalid PublicKey")
	}
	mArr = MessageArray{
		Array: append(mArr.Array, messageEnc),
		Keys:  append(mArr.Keys, k1),
	}

	cMutex.Unlock()
	return
}

//EncryptCircuitGatesAES encrypts an array of circuit gates with a given symmetric key using AES algorithm
func EncryptCircuitGatesAES(key []byte, gates []CircuitGate) []CircuitGate {
	encGates := gates
	var tmp []byte
	var ok bool
	for k, val := range gates {
		//Encrypt gateID
		tmp, ok = EncryptAES(key, []byte(val.GateID))
		if !ok {
			panic("!ok message encryption")
		}
		encGates[k].GateID = tmp
		//Encrypt gate inputs
		var concat [][]byte
		for _, val2 := range val.GateInputs {
			tmp, ok = EncryptAES(key, []byte(val2))
			if !ok {
				panic("!ok message encryption")
			}
			concat = append(concat, tmp)
		}
		encGates[k].GateInputs = concat
		//Encrypt truth table
		//Left for now for further discussion
	}
	return encGates
}

//DecryptCircuitGatesAES decrypts an array of circuit gates with a given symmetric key using AES algorithm
func DecryptCircuitGatesAES(key []byte, gates []CircuitGate) []CircuitGate {
	decGates := gates
	var tmp []byte
	var ok bool
	for k, val := range gates {
		//Decrypt gateID
		tmp, ok = DecryptAES(key, []byte(val.GateID))
		if !ok {
			panic("!ok message decryption")
		}
		decGates[k].GateID = tmp
		//Encrypt gate inputs
		var concat [][]byte
		for _, val2 := range val.GateInputs {
			tmp, ok = DecryptAES(key, []byte(val2))
			if !ok {
				panic("!ok message decryption")
			}
			concat = append(concat, tmp)
		}
		decGates[k].GateInputs = concat
		//decrypt truth table
		//Left for now for further discussion
	}
	return decGates
}

//EncryptGarbledGatesAES encrypts an array of garbled gates with a given symmetric key using AES algorithm
func EncryptGarbledGatesAES(key []byte, gates []GarbledGate) []GarbledGate {
	encGates := gates
	var tmp []byte
	var ok bool
	for k, val := range gates {
		//Encrypt gateID
		tmp, ok = EncryptAES(key, []byte(val.GateID))
		if !ok {
			panic("!ok message encryption")
		}
		encGates[k].GateID = tmp
		//Encrypt gate inputs
		var concat [][]byte
		for _, val2 := range val.GateInputs {
			tmp, ok = EncryptAES(key, []byte(val2))
			if !ok {
				panic("!ok message encryption")
			}
			concat = append(concat, tmp)
		}
		encGates[k].GateInputs = concat
		//Encrypt GarbledTable
		var concat2 [][]byte
		for _, val2 := range val.GarbledValues {
			tmp, ok = EncryptAES(key, val2)
			if !ok {
				panic("!ok message encryption")
			}
			concat2 = append(concat2, tmp)
		}
		encGates[k].GarbledValues = concat2
	}
	return encGates
}

//DecryptGarbledGatesAES decrypts an array of garbled gates with a given symmetric key using AES algorithm
func DecryptGarbledGatesAES(key []byte, gates []GarbledGate) []GarbledGate {
	decGates := gates
	var tmp []byte
	var ok bool
	for k, val := range gates {
		//Decrypt gateID
		tmp, ok = DecryptAES(key, []byte(val.GateID))
		if !ok {
			panic("!ok message decryption")
		}
		decGates[k].GateID = tmp
		//Dcrypt gate inputs
		var concat [][]byte
		for _, val2 := range val.GateInputs {
			tmp, ok = DecryptAES(key, []byte(val2))
			if !ok {
				panic("!ok message decryption")
			}
			concat = append(concat, tmp)
		}
		decGates[k].GateInputs = concat
		//Decrypt GarbledTable
		var concat2 [][]byte
		for _, val2 := range val.GarbledValues {
			tmp, ok = DecryptAES(key, val2)
			if !ok {
				panic("!ok message decryption")
			}
			concat2 = append(concat2, tmp)
		}
		decGates[k].GarbledValues = concat2
	}
	return decGates
}

//EncryptWiresAES encrypts an array of wires with a given key using AES Algorithm
func EncryptWiresAES(key []byte, wArr []Wire) []Wire {
	nWArr := wArr
	var ok bool
	for k, val := range wArr {
		//Encrypt wireID

		//Encrypt WireLabel
		nWArr[k].WireLabel, ok = EncryptAES(key, val.WireLabel)
		if !ok {
			panic("!ok message encryption")
		}
	}
	return nWArr
}

//DecryptWiresAES decrypts an array of wires with a given key using AES Algorithm
func DecryptWiresAES(key []byte, wArr []Wire) []Wire {
	nWArr := wArr
	var ok bool
	for k, val := range wArr {
		//Encrypt wireID

		//Encrypt WireLabel
		nWArr[k].WireLabel, ok = DecryptAES(key, val.WireLabel)
		if !ok {
			panic("!ok message decryption")
		}
	}
	return nWArr
}

//EncryptRandomnessAES encrypts a randomness container with a given key using AES Algorithm
func EncryptRandomnessAES(key []byte, rArr Randomness) Randomness {
	nRArr := rArr
	//Everything has to be converted into byte arrays.. message wise
	return nRArr
}

//DecryptRandomnessAES decrypts a randomness container with a given key using AES Algorithm
func DecryptRandomnessAES(key []byte, rArr Randomness) Randomness {
	nRArr := rArr
	//Everything has to be converted into byte arrays.. message wise
	return nRArr
}

//EncryptPartyInfoAES encrypts PartyInfo container with a given key using AES Algorithm
func EncryptPartyInfoAES(key []byte, pI PartyInfo) (nPI PartyInfo) {
	var ok bool
	//Encrypt IP
	nPI.IP, ok = EncryptAES(key, pI.IP)
	if !ok {
		panic("!ok message encryption")
	}
	//Encrypt Port
	nPI.Port = pI.Port
	//Should be converted into byte array
	//Encrypt PublicKey
	nPI.PublicKey, ok = EncryptAES(key, pI.PublicKey)
	if !ok {
		panic("!ok message encryption")
	}
	return
}

//DecryptPartyInfoAES decrypts PartyInfo container with a given key using AES Algorithm
func DecryptPartyInfoAES(key []byte, pI PartyInfo) (nPI PartyInfo) {
	var ok bool
	//Encrypt IP
	nPI.IP, ok = DecryptAES(key, pI.IP)
	if !ok {
		panic("!ok message decryption")
	}
	//Decrypt Port
	nPI.Port = pI.Port
	//Should be converted into byte array
	//Decrypt PublicKey
	nPI.PublicKey, ok = DecryptAES(key, pI.PublicKey)
	if !ok {
		panic("!ok message decryption")
	}
	return
}

//EncryptMessageAES takes a symmetric key and message, and encrypts the message using that key
func EncryptMessageAES(key []byte, msg Message) (nMsg Message) {
	nMsg = msg
	var ok bool
	var tmp []byte
	if string(msg.Type) == "Garble" {
		//Encrypt input gates
		nMsg.Circuit.InputGates = EncryptCircuitGatesAES(key, msg.Circuit.InputGates)
		//Encrypt Middle Gates
		nMsg.Circuit.MiddleGates = EncryptCircuitGatesAES(key, msg.Circuit.MiddleGates)
		//Encrypt Output Gates
		nMsg.Circuit.OutputGates = EncryptCircuitGatesAES(key, msg.Circuit.OutputGates)
		//Encrypt Randomness
		nMsg.Randomness = EncryptRandomnessAES(key, msg.Randomness)
		//Encrypt NextServer Info
		nMsg.NextServer = EncryptPartyInfoAES(key, msg.NextServer)
	} else if string(msg.Type) == "ReRand" {
		//Encrypt input gates
		nMsg.GarbledMessage.InputGates = EncryptGarbledGatesAES(key, msg.GarbledMessage.InputGates)
		//Encrypt middle gates
		nMsg.GarbledMessage.MiddleGates = EncryptGarbledGatesAES(key, msg.GarbledMessage.MiddleGates)
		//Encrypt output gates
		nMsg.GarbledMessage.OutputGates = EncryptGarbledGatesAES(key, msg.GarbledMessage.OutputGates)
		//Encrypt GarbledMessage Input wires
		nMsg.GarbledMessage.InputWires = EncryptWiresAES(key, msg.GarbledMessage.InputWires)
		//Encrypt GarbledMessage Output wires
		nMsg.GarbledMessage.OutputWires = EncryptWiresAES(key, msg.GarbledMessage.OutputWires)
		//Encrypt Randomness
		nMsg.Randomness = EncryptRandomnessAES(key, msg.Randomness)
		//Encrypt NextServer Info
		nMsg.NextServer = EncryptPartyInfoAES(key, msg.NextServer)
	} else if string(msg.Type) == "SEval" {
		if len(msg.GarbledMessage.InputGates) != 0 {
			//Encrypt input gates
			nMsg.GarbledMessage.InputGates = EncryptGarbledGatesAES(key, msg.GarbledMessage.InputGates)
			//Encrypt middle gates
			nMsg.GarbledMessage.MiddleGates = EncryptGarbledGatesAES(key, msg.GarbledMessage.MiddleGates)
			//Encrypt output gates
			nMsg.GarbledMessage.OutputGates = EncryptGarbledGatesAES(key, msg.GarbledMessage.OutputGates)

		}
		//Encrypt NextServer Info
		nMsg.NextServer = EncryptPartyInfoAES(key, msg.NextServer)
	} else if string(msg.Type) == "CEval" {
		//Encrypt InputWires
		nMsg.InputWires = EncryptWiresAES(key, msg.InputWires)
		//Encrypt NextServer Info
		nMsg.NextServer = EncryptPartyInfoAES(key, msg.NextServer)
	}

	//Encrypt the type
	tmp, ok = EncryptAES(key, []byte(msg.Type))
	if !ok {
		panic("!ok message encryption")
	}
	nMsg.Type = tmp

	return nMsg
}

//DecryptMessageAES takes a symmetric key and message, and decrypts the message using that key
func DecryptMessageAES(key []byte, msg Message) (nMsg Message) {
	nMsg = msg
	var ok bool
	var tmp []byte

	//Decrypt the type
	tmp, ok = DecryptAES(key, []byte(msg.Type))
	if !ok {
		panic("!ok message encryption")
	}
	nMsg.Type = tmp

	if string(nMsg.Type) == "Garble" {
		//Decrypt input gates
		nMsg.Circuit.InputGates = DecryptCircuitGatesAES(key, msg.Circuit.InputGates)
		//Decrypt Middle Gates
		nMsg.Circuit.MiddleGates = DecryptCircuitGatesAES(key, msg.Circuit.MiddleGates)
		//Decrypt Output Gates
		nMsg.Circuit.OutputGates = DecryptCircuitGatesAES(key, msg.Circuit.OutputGates)
		//Decrypt Randomness
		nMsg.Randomness = DecryptRandomnessAES(key, msg.Randomness)
		//Decrypt NextServer Info
		nMsg.NextServer = DecryptPartyInfoAES(key, msg.NextServer)
	} else if string(nMsg.Type) == "ReRand" {
		//Decrypt input gates
		nMsg.GarbledMessage.InputGates = DecryptGarbledGatesAES(key, msg.GarbledMessage.InputGates)
		//Decrypt middle gates
		nMsg.GarbledMessage.MiddleGates = DecryptGarbledGatesAES(key, msg.GarbledMessage.MiddleGates)
		//Decrypt output gates
		nMsg.GarbledMessage.OutputGates = DecryptGarbledGatesAES(key, msg.GarbledMessage.OutputGates)
		//Decrypt GarbledMessage Input wires
		nMsg.GarbledMessage.InputWires = DecryptWiresAES(key, msg.GarbledMessage.InputWires)
		//Decrypt GarbledMessage Output wires
		nMsg.GarbledMessage.OutputWires = DecryptWiresAES(key, msg.GarbledMessage.OutputWires)
		//Decrypt Randomness
		nMsg.Randomness = DecryptRandomnessAES(key, msg.Randomness)
		//Decrypt NextServer Info
		nMsg.NextServer = DecryptPartyInfoAES(key, msg.NextServer)
	} else if string(nMsg.Type) == "SEval" {
		if len(msg.GarbledMessage.InputGates) != 0 {
			//Decrypt input gates
			nMsg.GarbledMessage.InputGates = DecryptGarbledGatesAES(key, msg.GarbledMessage.InputGates)
			//Decrypt middle gates
			nMsg.GarbledMessage.MiddleGates = DecryptGarbledGatesAES(key, msg.GarbledMessage.MiddleGates)
			//Decrypt output gates
			nMsg.GarbledMessage.OutputGates = DecryptGarbledGatesAES(key, msg.GarbledMessage.OutputGates)

		}
		//Decrypt NextServer Info
		nMsg.NextServer = DecryptPartyInfoAES(key, msg.NextServer)

	} else if string(nMsg.Type) == "CEval" {
		//Decrypt InputWires
		nMsg.InputWires = DecryptWiresAES(key, msg.InputWires)
		//Decrypt NextServer Info
		nMsg.NextServer = DecryptPartyInfoAES(key, msg.NextServer)
	}

	return nMsg
}

//RandomSymmKeyGen Generates a random key for the AES algorithm
func RandomSymmKeyGen() (key []byte) {
	key = make([]byte, 32)

	_, err := cryptoRand.Read(key)
	if err != nil {
		panic("Error generating random symmetric key")
	}
	return
}

//GenerateInputWiresValidate Given circuit and randomness generate the input wires corresponding to server n-1
func GenerateInputWiresValidate(circ Circuit, rArr []Randomness, cID int64) (in [][]byte, out [][]byte) {
	cMutex.Lock()
	inputSize, outputSize := GetInputSizeOutputSize(circ)
	in = YaoGarbledCkt_in(rArr[0].Rin, rArr[0].LblLength, inputSize)
	out = YaoGarbledCkt_out(rArr[0].Rout, rArr[0].LblLength, outputSize)

	for i := 1; i < len(rArr)-1; i++ {
		randIn := genRandomR(len(circ.InputGates), 1, rArr[i].Rin)
		randOut := genRandomR(len(circ.OutputGates), 1, rArr[i].Rout)
		in = reRandWires(randIn, in, true)
		out = reRandWires(randOut, out, false)
	}
	cMutex.Unlock()
	return
}

func reRandWires(Randoms []*big.Int, wires [][]byte, in bool) [][]byte {
	for j := 0; j < len(Randoms); j++ {
		R := Randoms[j]
		if in {
			newWire := new(big.Int).SetBytes(wires[j*4])
			newWire = newWire.Mul(newWire, R)
			newWire = newWire.Mod(newWire, fromHex(primeHex))
			wires[j*4] = newWire.Bytes()

			newWire = new(big.Int).SetBytes(wires[j*4+1])
			newWire = newWire.Mul(newWire, R)
			newWire = newWire.Mod(newWire, fromHex(primeHex))
			wires[j*4+1] = newWire.Bytes()

			newWire = new(big.Int).SetBytes(wires[j*4+2])
			newWire = newWire.Mul(newWire, R)
			newWire = newWire.Mod(newWire, fromHex(primeHex))
			wires[j*4+2] = newWire.Bytes()

			newWire = new(big.Int).SetBytes(wires[j*4+3])
			newWire = newWire.Mul(newWire, R)
			newWire = newWire.Mod(newWire, fromHex(primeHex))
			wires[j*4+3] = newWire.Bytes()
		} else {
			newWire := new(big.Int).SetBytes(wires[j*2])
			newWire = newWire.Mul(newWire, R)
			newWire = newWire.Mod(newWire, fromHex(primeHex))
			wires[j*2] = newWire.Bytes()

			newWire = new(big.Int).SetBytes(wires[j*2+1])
			newWire = newWire.Mul(newWire, R)
			newWire = newWire.Mod(newWire, fromHex(primeHex))
			wires[j*2+1] = newWire.Bytes()
		}
	}
	return wires
}

//GenerateRandomness generates randomness array corresponding to NumberOfServers with a certain computation ID
func GenerateRandomness(numberOfServers int, cID int64) []Randomness {
	rArr := make([]Randomness, numberOfServers)
	rand.Seed(cID)

	for k := range rArr {
		rArr[k] = Randomness{
			Rin:       rand.Int63(),
			Rout:      rand.Int63(),
			Rgc:       rand.Int63(),
			LblLength: 16, //Should be rand.Int()%16 + 16
		}
	}
	return rArr
}

//CompareWires Takes a garbled circuit and compares wires to input,output wires provided by the user
func CompareWires(gcm GarbledMessage, arrIn [][]byte, arrOut [][]byte) bool {
	for k, val := range gcm.InputWires {
		if bytes.Compare(arrIn[k], val.WireLabel) != 0 {
			fmt.Println("I was cheated on this: ", arrIn[k], val.WireLabel)
			//			panic("The server has cheated me") //redo the process, by recovering from panic by recalling comm
			return false
		}
	}
	for k, val := range gcm.OutputWires {
		if bytes.Compare(arrOut[k], val.WireLabel) != 0 {

			fmt.Println("I was cheated on this: ", arrOut[k], val.WireLabel)
			//panic("The server has cheated me") //redo the process, by recovering from panic by recalling comm
			return false
		}
	}
	return true
}

//SendToServer Invokes the post method on the server
func SendToServer(k MessageArray, ip []byte, port int) bool {
	circuitJSON, err := json.Marshal(k)
	req, err := http.NewRequest("POST", "http://"+string(ip)+":"+strconv.Itoa(port)+"/post", bytes.NewBuffer(circuitJSON))
	if err != nil {
		fmt.Println("generating request failed")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
	if err != nil {
		//log.Fatal(err)
		return false
	}
	return true
}

//GetFromServer Invokes the get method on the server
func GetFromServer(tokenChallenge Token, ip []byte, port int) (token Token, ok bool) {
	ok = false //assume failure
	iDJSON, err := json.Marshal(tokenChallenge)
	req, err := http.NewRequest("GET", "http://"+string(ip)+":"+strconv.Itoa(port)+"/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &token)
	if err != nil {
		return
	}
	resp.Body.Close()
	ok = true
	return
}

//SendToDirectory Invokes the post method on the directory
func SendToDirectory(k RegisterationMessage, ip []byte, port int) bool {
	Decentralization = false

	circuitJSON, err := json.Marshal(k)
	req, err := http.NewRequest("POST", "http://"+string(ip)+":"+strconv.Itoa(port)+"/post", bytes.NewBuffer(circuitJSON))
	if err != nil {
		fmt.Println("generating request failed")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
	if err != nil {
		//log.Fatal(err)
		return false
	}
	return true
}

//GetFromDirectory Invokes the get method on the directory
func GetFromDirectory(k CycleRequestMessage, ip []byte, port int) (cyc CycleMessage, ok bool) {
	ok = false //assume failure
	iDJSON, err := json.Marshal(k)
	req, err := http.NewRequest("GET", "http://"+string(ip)+":"+strconv.Itoa(port)+"/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &cyc)
	if err != nil {
		return
	}
	resp.Body.Close()
	ok = true
	return
}

//SendToClient Invokes the post method on the server
func SendToClient(res ResEval, ip []byte, port int) bool {
	circuitJSON, err := json.Marshal(res)
	req, err := http.NewRequest("POST", "http://"+string(ip)+":"+strconv.Itoa(port)+"/post", bytes.NewBuffer(circuitJSON))
	if err != nil {
		fmt.Println("generating request failed")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
	if err != nil {
		//log.Fatal(err)
		return false
	}
	return true
}

//GetFromClient Invokes the get method on the client
func GetFromClient(tokenChallenge Token, ip []byte, port int) (token Token, ok bool) {
	ok = false //assume failure
	iDJSON, err := json.Marshal(tokenChallenge)
	req, err := http.NewRequest("GET", "http://"+string(ip)+":"+strconv.Itoa(port)+"/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &token)
	if err != nil {
		return
	}
	resp.Body.Close()
	ok = true
	return
}

//SendToServerGarble used in pt2
func SendToServerGarble(k CircuitMessage) bool {
	circuitJSON, err := json.Marshal(k)
	req, err := http.NewRequest("POST", "http://localhost:8080/post", bytes.NewBuffer(circuitJSON))
	if err != nil {
		fmt.Println("generating request failed")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
	if err != nil {
		//log.Fatal(err)
		return false
	}
	return true
}

//GetFromServerGarble used in pt2
func GetFromServerGarble(id string) (k GarbledMessage, ok bool) {
	ok = false //assume failure
	iDJSON, err := json.Marshal(ComID{CID: []byte(id)})
	req, err := http.NewRequest("GET", "http://localhost:8080/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &k)
	if err != nil {
		return
	}
	resp.Body.Close()
	if string(k.CID) != id {
		panic("The server sent me the wrong circuit") //replace with a request repeat.
	}
	ok = true
	return
}

//SendToServerEval used in pt2
func SendToServerEval(k GarbledMessage) bool {
	circuitJSON, err := json.Marshal(k)
	req, err := http.NewRequest("POST", "http://localhost:8081/post", bytes.NewBuffer(circuitJSON))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

//GetFromServerEval used in pt2
func GetFromServerEval(id string) (res [][]byte, ok bool) {
	ok = false // assume failure
	iDJSON, err := json.Marshal(ComID{CID: []byte(id)})
	req, err := http.NewRequest("GET", "http://localhost:8081/get", bytes.NewBuffer(iDJSON))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	var k ResEval
	err = json.Unmarshal(body, &k)
	if err != nil {
		return
	}
	resp.Body.Close()
	if string(k.CID) != id {
		panic("The server sent me the wrong circuit") //replace with a request repeat.
	}
	res = k.Res
	//fmt.Println("Result Returned", k.Res)
	ok = true
	return
}

//GenNRandNumbers generating random byte arrays
func GenNRandNumbers(n int, length int, r int64, considerR bool) [][]byte {
	if considerR {
		rand.Seed(r)
	}
	seeds := make([][]byte, n)
	for i := 0; i < n; i++ {
		seeds[i] = make([]byte, length)
		_, err := rand.Read(seeds[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	return seeds
}

//YaoGarbledCkt_in input wire garbling
func YaoGarbledCkt_in(rIn int64, length int, inputSize int) [][]byte {
	return GenNRandNumbers(2*inputSize, length, rIn, true)
}

//YaoGarbledCkt_out output wire garbling
func YaoGarbledCkt_out(rOut int64, length int, outputSize int) [][]byte {
	// only one output bit for now
	return GenNRandNumbers(2*outputSize, length, rOut, true)
}

//YaoGarbledCkt_mask mask for garbling and rerandomization
func YaoGarbledCkt_mask(rMask int64, length int, outputSize int) [][]byte {
	return GenNRandNumbers(2*outputSize, length, rMask, true)
}

//EncryptAES symmetric encryption using AES algorithm
func EncryptAES(encKey []byte, plainText []byte) (ciphertext []byte, ok bool) {

	ok = false //assume failure
	//			encKey = append(encKey, hash)
	c, err := aes.NewCipher(encKey)
	//fmt.Println("cipher enc: ", c)
	if err != nil {
		//fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	//fmt.Println("gcm enc: ", gcm)
	if err != nil {
		//fmt.Println(err)
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	//fmt.Println("nonce enc: ", nonce)
	if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
		//fmt.Println(err)
		return
	}
	ciphertext = gcm.Seal(nonce, nonce, plainText, nil)
	//fmt.Println("ciphertext enc: ", ciphertext)
	//fmt.Println(ciphertext)
	ok = true

	return
}

//DecryptAES symmetric decryption using AES algorithm
func DecryptAES(encKey []byte, cipherText []byte) (plainText []byte, ok bool) {

	ok = false //assume failure

	c, err := aes.NewCipher(encKey)
	//fmt.Println("cipher dec: ", c)
	if err != nil {
		//fmt.Println(err)
		return
	}

	gcm, err := cipher.NewGCM(c)
	//fmt.Println("gcm dec: ", gcm)

	if err != nil {
		//fmt.Println(err)
		return
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		//fmt.Println(err)
		return
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	//fmt.Println("ciphertext dec: ", cipherText)
	//fmt.Println("nonce dec: ", nonce)

	plainText, err = gcm.Open(nil, nonce, cipherText, nil)
	//fmt.Println("plain text dec: ", plainText)

	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(plaintext))
	ok = true
	return
}

//Garble circuit garbling. This now supports only 2-input gates
func Garble(circ CircuitMessage) GarbledMessage {
	fmt.Println("In Garble")
	//Input and output size extraction... This should be improved to adobt non-2-input gates
	inputSize := len(circ.InputGates) * 2
	outputSize := len(circ.OutputGates)

	//Generating array of input wires, output wires
	arrIn := YaoGarbledCkt_in(circ.Rin, circ.LblLength, inputSize)
	arrOut := YaoGarbledCkt_out(circ.Rout, circ.LblLength, outputSize)
	//arrMask := YaoGarbledCkt_mask(circ.Rmask, circ.LblLength, outputSize)

	//Maps to keep track of wires in and out of gates throughout the topological traversal of the circuit
	inWires := make(map[string][]Wire)
	outWires := make(map[string][]Wire)

	//The randomization seed is now Rgc
	rand.Seed(circ.Rgc)

	//Preparing containers for garbled circuit and input/output wires
	var gc GarbledCircuit
	inputWiresGC := []Wire{}
	outputWiresGC := []Wire{}

	//Setting computation ID
	gc.CID = circ.CID

	// Input Gates Garbling
	var wInCnt int = 0
	fmt.Println("Garbling Input Gates")
	for k, val := range circ.InputGates {
		gc.InputGates = append(gc.InputGates, GarbledGate{
			Gate: Gate{
				GateID: val.GateID,
			},
		})

		gc.InputGates[k].GateInputs = val.GateInputs

		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		if len(val.GateInputs) > 2 || inCnt > 2 {
			panic("non-2-input Gate")
		}

		//fmt.Printf("%v, %T\n", val.GateID, val.GateID)

		inWires[string(val.GateID)] = []Wire{}
		//Fetching input wire values from the arrIn array
		for i := 0; i < inCnt; i++ {
			inWires[string(val.GateID)] = append(inWires[string(val.GateID)], Wire{
				WireLabel: arrIn[wInCnt],
			}, Wire{
				WireLabel: arrIn[wInCnt+1],
			})
			inputWiresGC = append(inputWiresGC, Wire{
				WireLabel: arrIn[wInCnt],
			}, Wire{
				WireLabel: arrIn[wInCnt+1],
			})
			wInCnt += 2
		}
		//Generating output wire labels and dumbing them into the big outWires container
		outWires[string(val.GateID)] = []Wire{}
		outWire := GenNRandNumbers(2, circ.LblLength, 0, false)
		outWires[string(val.GateID)] = append(outWires[string(val.GateID)], Wire{
			WireLabel: outWire[0],
		}, Wire{
			WireLabel: outWire[1],
		})
		//in1:	0	0	1	1
		//in0:	0	1	0	1
		//out:	1	0	0	1

		//fmt.Println("Here we getting inWires: \n")

		//Preparing containers for the garbled values, the Cs used in elgamal encryption, and the mask
		lengthVal := len(val.TruthTable)
		gc.InputGates[k].GarbledValues = make([][]byte, lengthVal) //To be removed

		gc.InputGates[k].GarbledValuesC1 = make([][]byte, lengthVal)
		gc.InputGates[k].GarbledValuesC2 = make([][]byte, lengthVal)
		gc.InputGates[k].KeyY = make([][]byte, lengthVal)

		for key, value := range val.TruthTable {
			//The loop below supported multi input. Since it won't affect performance it was left this way
			var firstInput []byte
			var finalInput []byte
			//This loop basically: treats the index in the truth table as the carrier of the input values.
			//i.e: index 1 for a 3 input gate will carry 001 which means the first input is 1, while the rest is 0
			//The loop extracts the value of the wire from the index and fetches the corresponding wire label from the container descriped above
			//The concatenation of the n-1 input wires is used to encrypt
			for i := 0; i < inCnt-1; /*You just added this -1*/ i++ {
				idx := ((key >> i) & (1))
				firstInput = append(firstInput, inWires[string(val.GateID)][(i*2)+idx].WireLabel...)
			} //Now it only extracts the first input wire in a 2-input or a 1-input gate/ no masks to encrypt

			//The final input wire is extracted
			finalInput = inWires[string(val.GateID)][((inCnt-1)*2)+((key>>(inCnt-1))&(1))].WireLabel

			//ElGamal Encryption Process:

			//Generate Keys:
			privOutput := GenerateElGamalKey(ByteSliceMul(firstInput, finalInput))

			//fetch output label
			var idxOut int
			if value {
				idxOut = 1
			}
			outKey := outWires[string(val.GateID)][int(idxOut)]

			//The encrypted value is now the X (private) part of the Gamal Key generated by the label.
			//After this point the label is obsolete, however it is stored to generate the key of the next gate (Input wire multiplication).
			outC1, outC2, err := elgamal.Encrypt(cryptoRand.Reader, &privOutput.PublicKey, false, outKey.WireLabel)

			gc.InputGates[k].GarbledValuesC1[key] = outC1.Bytes()
			gc.InputGates[k].GarbledValuesC2[key] = outC2.Bytes()
			gc.InputGates[k].KeyY[key] = privOutput.Y.Bytes()
			if err != nil {
				panic("Out Label encryption Error in Input Gates")
			}
		}

		//fmt.Println("\nwe got'em inWires \n")

	}

	fmt.Println("Garbling Middle Gates")
	//Middle Gates Garbling
	for k, val := range circ.MiddleGates {

		//Adding the gate container
		gc.MiddleGates = append(gc.MiddleGates, GarbledGate{
			Gate: Gate{
				GateID: val.GateID,
			},
		})

		//Setting gate inputs
		gc.MiddleGates[k].GateInputs = val.GateInputs

		//Calculating the input wires count
		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		if len(val.GateInputs) > 2 || inCnt > 2 {
			panic("non-2-input Gate")
		}

		//Creating a place holder for the wires in the map
		inWires[string(val.GateID)] = []Wire{}

		//Extracting the input wire labels from the map
		for _, j := range val.GateInputs {
			inWires[string(val.GateID)] = append(inWires[string(val.GateID)], outWires[string(j)][0])
			inWires[string(val.GateID)] = append(inWires[string(val.GateID)], outWires[string(j)][1])
			//wInCnt++
		}

		//Generating the output wire labels and storing them into the map
		outWires[string(val.GateID)] = []Wire{}
		outWire := GenNRandNumbers(2, circ.LblLength, 0, false)
		outWires[string(val.GateID)] = append(outWires[string(val.GateID)], Wire{
			WireLabel: outWire[0],
		}, Wire{
			WireLabel: outWire[1],
		})

		//gc.MiddleGates[k].GarbledValues = make([][]byte, len(val.TruthTable))

		//Preparing containers for the garbled values, the Cs used in elgamal encryption, and the mask
		lengthVal := len(val.TruthTable)
		gc.MiddleGates[k].GarbledValues = make([][]byte, lengthVal) //to be removed

		gc.MiddleGates[k].GarbledValuesC1 = make([][]byte, lengthVal)
		gc.MiddleGates[k].GarbledValuesC2 = make([][]byte, lengthVal)
		gc.MiddleGates[k].KeyY = make([][]byte, lengthVal)

		for key, value := range val.TruthTable {
			//The loop below supported multi input. Since it won't affect performance it was left this way
			var firstInput []byte
			var finalInput []byte
			//This loop basically: treats the index in the truth table as the carrier of the input values.
			//i.e: index 1 for a 3 input gate will carry 001 which means the first input is 1, while the rest is 0
			//The loop extracts the value of the wire from the index and fetches the corresponding wire label from the container descriped above
			//The concatenation of the n-1 input wires is used to encrypt
			for i := 0; i < inCnt-1; /*You just added this -1*/ i++ {
				idx := ((key >> i) & (1))
				firstInput = append(firstInput, inWires[string(val.GateID)][(i*2)+idx].WireLabel...)
			} //Now it only extracts the first input wire in a 2-input or a 1-input gate/ no masks to encrypt

			//The final input wire is extracted
			finalInput = inWires[string(val.GateID)][((inCnt-1)*2)+((key>>(inCnt-1))&(1))].WireLabel

			//ElGamal Encryption Process:
			//Generate Keys:
			privOutput := GenerateElGamalKey(ByteSliceMul(firstInput, finalInput))

			//fetch output label
			var idxOut int = 0
			if value {
				idxOut = 1
			}
			outKey := outWires[string(val.GateID)][int(idxOut)]

			//The encrypted value is now the X (private) part of the Gamal Key generated by the label.
			//After this point the label is obsolete, however it is stored to generate the key of the next gate (Input wire multiplication).
			outC1, outC2, err := elgamal.Encrypt(cryptoRand.Reader, &privOutput.PublicKey, false, outKey.WireLabel)

			gc.MiddleGates[k].GarbledValuesC1[key] = outC1.Bytes()
			gc.MiddleGates[k].GarbledValuesC2[key] = outC2.Bytes()
			gc.MiddleGates[k].KeyY[key] = privOutput.Y.Bytes()

			if err != nil {
				panic("Out Label encryption Error in Middle Gates")
			}

		}

	}

	fmt.Println("Garbling Output Gates")
	//Output Gates Garbling

	//A count for the output wires
	wOutCnt := 0
	for k, val := range circ.OutputGates {
		//container initialization
		gc.OutputGates = append(gc.OutputGates, GarbledGate{
			Gate: Gate{
				GateID: val.GateID,
			},
		})

		//Setting the gate inputs
		gc.OutputGates[k].GateInputs = val.GateInputs

		//Calculating the number of input entries
		inCnt := int(math.Log2(float64(len(val.TruthTable))))

		if len(val.GateInputs) > 2 || inCnt > 2 {
			panic("non-2-input Gate")
		}

		//Fetching the input wire labels from the map
		inWires[string(val.GateID)] = []Wire{}
		for _, j := range val.GateInputs {
			inWires[string(val.GateID)] = append(inWires[string(val.GateID)], outWires[string(j)][0])
			inWires[string(val.GateID)] = append(inWires[string(val.GateID)], outWires[string(j)][1])

			//wInCnt++
		}

		//Adding the output wire labels to the map
		outWires[string(val.GateID)] = []Wire{}
		outWires[string(val.GateID)] = append(outWires[string(val.GateID)], Wire{
			WireLabel: arrOut[wOutCnt],
		}, Wire{
			WireLabel: arrOut[wOutCnt+1],
		})

		//Storing the output wire labels independently
		outputWiresGC = append(outputWiresGC, Wire{
			WireLabel: arrOut[wOutCnt],
		}, Wire{
			WireLabel: arrOut[wOutCnt+1],
		})
		wOutCnt += 2

		//Preparing containers for the garbled values, the Cs used in elgamal encryption, and the mask
		lengthVal := len(val.TruthTable)
		gc.OutputGates[k].GarbledValues = make([][]byte, lengthVal) //to be removed

		gc.OutputGates[k].GarbledValuesC1 = make([][]byte, lengthVal)
		gc.OutputGates[k].GarbledValuesC2 = make([][]byte, lengthVal)
		gc.OutputGates[k].KeyY = make([][]byte, lengthVal)

		for key, value := range val.TruthTable {
			//The loop below supported multi input. Since it won't affect performance it was left this way
			var firstInput []byte
			var finalInput []byte
			//This loop basically: treats the index in the truth table as the carrier of the input values.
			//i.e: index 1 for a 3 input gate will carry 001 which means the first input is 1, while the rest is 0
			//The loop extracts the value of the wire from the index and fetches the corresponding wire label from the container descriped above
			//The concatenation of the n-1 input wires is used to encrypt
			for i := 0; i < inCnt-1; /*You just added this -1*/ i++ {
				idx := ((key >> i) & (1))
				firstInput = append(firstInput, inWires[string(val.GateID)][(i*2)+idx].WireLabel...)
			} //Now it only extracts the first input wire in a 2-input or a 1-input gate/ no masks to encrypt

			//The final input wire is used to encrypt the mask
			finalInput = inWires[string(val.GateID)][((inCnt-1)*2)+((key>>(inCnt-1))&(1))].WireLabel

			//ElGamal Encryption Process:
			//Generate Keys:
			privOutput := GenerateElGamalKey(ByteSliceMul(firstInput, finalInput))
			//fetch output label
			var idxOut int = 0
			if value {
				idxOut = 1
			}
			outKey := outWires[string(val.GateID)][int(idxOut)]

			//The encrypted value is now the X (private) part of the Gamal Key generated by the label.
			//After this point the label is obsolete, however it is stored to generate the key of the next gate (Input wire multiplication).
			outC1, outC2, err := elgamal.Encrypt(cryptoRand.Reader, &privOutput.PublicKey, false, outKey.WireLabel)

			gc.OutputGates[k].GarbledValuesC1[key] = outC1.Bytes()
			gc.OutputGates[k].GarbledValuesC2[key] = outC2.Bytes()
			gc.OutputGates[k].KeyY[key] = privOutput.Y.Bytes()

			if err != nil {
				panic("Out Label encryption Error in Output Gates")
			}
		}

	}

	//fmt.Println(arrIn)
	//fmt.Println(arrOut)
	//fmt.Println("Input Wires GC:", inWires)
	//fmt.Println("Output Wires GC:", outWires)
	//fmt.Println("GC: ", gc)
	gm := GarbledMessage{
		InputWires:     inputWiresGC,
		GarbledCircuit: gc,
		OutputWires:    outputWiresGC,
	}
	return gm
}

func genRandomR(n int, factor int, seed int64) (rands []*big.Int) {
	rands = make([]*big.Int, n*factor)
	rand.Seed(seed)
	one := new(big.Int).SetInt64(1)
	for i := 0; i < n*factor; i++ {
		R1 := new(big.Int).SetInt64(rand.Int63())
		for one.Cmp(new(big.Int).GCD(new(big.Int), new(big.Int), R1, new(big.Int).Sub(fromHex(primeHex), one))) != 0 {
			R1 = new(big.Int).SetInt64(rand.Int63())

		}
		rands[i] = R1
	}
	return
}

//ReRand does things for a certain circuit view
func ReRand(gcm GarbledMessage, r Randomness) GarbledMessage {
	fmt.Println("In Rerand")
	one := new(big.Int).SetInt64(1)

	randIn := genRandomR(len(gcm.InputGates), 1, r.Rin)
	randOut := genRandomR(len(gcm.OutputGates), 1, r.Rout)
	randGC := genRandomR(1, 1, r.Rgc)

	fmt.Println("Rerandomizing Input Gates")
	for i, val := range gcm.InputGates {
		R1 := randIn[i]
		R2 := randIn[i]
		R := randGC[0]
		//fmt.Println(i)
		for k := range val.GarbledValuesC1 {
			c1 := new(big.Int).SetBytes(val.GarbledValuesC1[k])
			c2 := new(big.Int).SetBytes(val.GarbledValuesC2[k])
			y := new(big.Int).SetBytes(val.KeyY[k])
			// R1 := randIn[4*i + k/2]
			// R2 := randIn[4*i + (k%2 - 1) + 3]
			Rsquared := new(big.Int).Mul(R1, R2)
			Rsquared = new(big.Int).Mod(Rsquared, fromHex(primeHex))
			Rinverse := new(big.Int).ModInverse(Rsquared, new(big.Int).Sub(fromHex(primeHex), one))
			newY := new(big.Int).Exp(y, Rsquared, fromHex(primeHex))
			newC1 := new(big.Int).Exp(c1, Rinverse, fromHex(primeHex))
			newC2 := new(big.Int).Mul(c2, R)
			newC2 = new(big.Int).Mod(newC2, fromHex(primeHex))
			key := elgamal.PublicKey{
				G: fromHex(generatorHex),
				P: fromHex(primeHex),
				Y: newY,
			}
			hider1, hider2, _ := elgamal.Encrypt(cryptoRand.Reader, &key, false, one.Bytes())
			newC1 = new(big.Int).Mul(newC1, hider1)
			newC1 = new(big.Int).Mod(newC1, fromHex(primeHex))
			newC2 = new(big.Int).Mul(newC2, hider2)
			newC2 = new(big.Int).Mod(newC2, fromHex(primeHex))
			val.GarbledValuesC1[k] = newC1.Bytes()
			val.GarbledValuesC2[k] = newC2.Bytes()
			val.KeyY[k] = newY.Bytes()
		}
	}
	fmt.Println("Rerandomizing Middle Gates")
	// layerCounter := 0
	// totalLayersVisited := len(gcm.InputGates)/2
	// thisLayerLength := len(gcm.InputGates)/2
	for _, val := range gcm.MiddleGates {
		// if (i == totalLayersVisited){
		// 	thisLayerLength /= 2
		// 	totalLayersVisited += thisLayerLength
		// 	layerCounter = 0
		// }
		R1 := randGC[0] // randGC[2*i]
		R2 := randGC[0] // randGC[2*i+1]
		R := randGC[0]  // randGC[len(gcm.InputGates) + i]
		// fmt.Println(2*i)
		// fmt.Println(2*i+1)
		// fmt.Println(len(gcm.InputGates) + i)

		for k := range val.GarbledValuesC1 {
			c1 := new(big.Int).SetBytes(val.GarbledValuesC1[k])
			c2 := new(big.Int).SetBytes(val.GarbledValuesC2[k])
			y := new(big.Int).SetBytes(val.KeyY[k])
			//fmt.Println(R1)
			//fmt.Println(R2)
			Rsquared := new(big.Int).Mul(R1, R2)
			Rsquared = new(big.Int).Mod(Rsquared, fromHex(primeHex))
			Rinverse := new(big.Int).ModInverse(Rsquared, new(big.Int).Sub(fromHex(primeHex), one))
			newY := new(big.Int).Exp(y, Rsquared, fromHex(primeHex))
			newC1 := new(big.Int).Exp(c1, Rinverse, fromHex(primeHex))
			newC2 := new(big.Int).Mul(c2, R)
			newC2 = new(big.Int).Mod(newC2, fromHex(primeHex))
			key := elgamal.PublicKey{
				G: fromHex(generatorHex),
				P: fromHex(primeHex),
				Y: newY,
			}
			hider1, hider2, _ := elgamal.Encrypt(cryptoRand.Reader, &key, false, one.Bytes())
			newC1 = new(big.Int).Mul(newC1, hider1)
			newC1 = new(big.Int).Mod(newC1, fromHex(primeHex))
			newC2 = new(big.Int).Mul(newC2, hider2)
			newC2 = new(big.Int).Mod(newC2, fromHex(primeHex))
			val.GarbledValuesC1[k] = newC1.Bytes()
			val.GarbledValuesC2[k] = newC2.Bytes()
			val.KeyY[k] = newY.Bytes()
		}
	}
	fmt.Println("Rerandomizing Output Gates")
	for i, val := range gcm.OutputGates {

		R1 := randGC[0] // randGC[len(randGC)- 2*len(gcm.OutputGates)+2*i]
		R2 := randGC[0] // randGC[len(randGC)- 2*len(gcm.OutputGates)+2*i+1]
		R := randOut[i]
		// fmt.Println(len(randGC)- 2*len(gcm.OutputGates)+2*i)
		// fmt.Println(len(randGC)- 2*len(gcm.OutputGates)+2*i+1)
		// fmt.Println(i)
		for k := range val.GarbledValuesC1 {
			c1 := new(big.Int).SetBytes(val.GarbledValuesC1[k])
			c2 := new(big.Int).SetBytes(val.GarbledValuesC2[k])
			y := new(big.Int).SetBytes(val.KeyY[k])
			Rsquared := new(big.Int).Mul(R1, R2)
			Rsquared = new(big.Int).Mod(Rsquared, fromHex(primeHex))
			Rinverse := new(big.Int).ModInverse(Rsquared, new(big.Int).Sub(fromHex(primeHex), one))
			newY := new(big.Int).Exp(y, Rsquared, fromHex(primeHex))
			newC1 := new(big.Int).Exp(c1, Rinverse, fromHex(primeHex))
			newC2 := new(big.Int).Mul(c2, R)
			newC2 = new(big.Int).Mod(newC2, fromHex(primeHex))
			key := elgamal.PublicKey{
				G: fromHex(generatorHex),
				P: fromHex(primeHex),
				Y: newY,
			}
			hider1, hider2, _ := elgamal.Encrypt(cryptoRand.Reader, &key, false, one.Bytes())
			newC1 = new(big.Int).Mul(newC1, hider1)
			newC1 = new(big.Int).Mod(newC1, fromHex(primeHex))
			newC2 = new(big.Int).Mul(newC2, hider2)
			newC2 = new(big.Int).Mod(newC2, fromHex(primeHex))
			val.GarbledValuesC1[k] = newC1.Bytes()
			val.GarbledValuesC2[k] = newC2.Bytes()
			val.KeyY[k] = newY.Bytes()
		}
	}
	return gcm
}

//Evaluate evaluate a garbled circuit
func Evaluate(gc GarbledMessage) (result ResEval) {
	fmt.Println("In Eval")
	//Setting computation ID
	result.CID = gc.CID
	//A map for storing the output wires for the ease of topological traversal
	outWires := make(map[string]Wire)
	var wInCnt int
	// Rk := new(big.Int).SetInt64(2669985732393126063)
	// R := Rk.Bytes()
	//Traversing Input Gates
	fmt.Println("Evaluating Input Gates")
	for _, val := range gc.InputGates {

		//Calculating the number of input entries
		inCnt := int(math.Log2(float64(len(val.GarbledValues))))

		//The loop below supported multi input. Since it won't affect performance it was left this way
		var firstInput []byte
		var finalInput []byte

		for i := 0; i < inCnt-1; /*You just added this -1*/ i++ {
			firstInput = append(firstInput, gc.InputWires[wInCnt].WireLabel...)
			wInCnt++
		}

		//fetching the final input
		finalInput = gc.InputWires[wInCnt].WireLabel
		wInCnt++

		//Generate Keys:
		privOutput := GenerateElGamalKey(ByteSliceMul(firstInput, finalInput))

		//A place holder for the output wires
		outWires[string(val.GateID)] = Wire{}

		//Loop over entries
		for garbledKey := range val.GarbledValuesC1 {

			//Decrypt the masked output label
			outLabel, err := elgamal.Decrypt(privOutput, false, new(big.Int).SetBytes(val.GarbledValuesC1[garbledKey]), new(big.Int).SetBytes(val.GarbledValuesC2[garbledKey]))

			//if decrypted
			if err == nil {
				//Add the wire label to the map to continue the traversal
				outWires[string(val.GateID)] = Wire{
					WireLabel: outLabel,
				}
				// outWires[string(val.GateID)] = Wire{
				// 	WireLabel: ByteSliceAdd(outLabel, R),
				// }
				//break out of looping over entries since you already found the output of the gate
				break
			}

		}

		//if no output label was found, then the evaluation failed
		if (bytes.Compare(outWires[string(val.GateID)].WireLabel, Wire{}.WireLabel)) == 0 {
			fmt.Println("Fail Evaluation Input Gate")
		}

		/*else {
			fmt.Println("\n\nYaaay\nGate ", val.GateID, " Now has an output wire: \n", outWires[val.GateID].WireLabel, "\n\n")
		}*/
	}
	fmt.Println("Evaluating Middle Gates")
	//Traversing middle gates
	for _, val := range gc.MiddleGates {

		//inCnt := len(val.GateInputs)
		var firstInput []byte
		var finalInput []byte
		for kGate, preGate := range val.GateInputs {
			if kGate == len(val.GateInputs)-1 {
				//The final input label
				finalInput = outWires[string(preGate)].WireLabel
			} else {
				//The concatentation of all previous labels
				firstInput = append(firstInput, outWires[string(preGate)].WireLabel...)
			}
			//wInCnt++
		}

		//Generate Keys:
		privOutput := GenerateElGamalKey(ByteSliceMul(firstInput, finalInput))

		//A place holder for the output wires
		outWires[string(val.GateID)] = Wire{}

		//Loop over entries
		for garbledKey := range val.GarbledValuesC1 {
			//Decrypt the masked output label
			outLabel, err := elgamal.Decrypt(privOutput, false, new(big.Int).SetBytes(val.GarbledValuesC1[garbledKey]), new(big.Int).SetBytes(val.GarbledValuesC2[garbledKey]))
			//if decrypted
			if err == nil {
				//Add the wire label to the map to continue the traversal
				outWires[string(val.GateID)] = Wire{
					WireLabel: outLabel,
				}
				// outWires[string(val.GateID)] = Wire{
				// 	WireLabel: ByteSliceAdd(outLabel, R),
				// }
				//break out of looping over entries since you already found the output of the gate
				break
			}

		}

		//if no output label was found, then the evaluation failed
		if (bytes.Compare(outWires[string(val.GateID)].WireLabel, Wire{}.WireLabel)) == 0 {
			fmt.Println("Fail Evaluation Middle Gate")
		}
		/*else {
			fmt.Println("\n\nYaaay\nGate ", val.GateID, " Now has an output wire: \n", outWires[val.GateID].WireLabel, "\n\n")
		}*/
	}
	fmt.Println("Evaluating Output Gates")
	//Traversing outputGates
	for _, val := range gc.OutputGates {

		//inCnt := len(val.GateInputs)
		var firstInput []byte
		var finalInput []byte
		for kGate, preGate := range val.GateInputs {
			if kGate == len(val.GateInputs)-1 {
				//The final input label
				finalInput = outWires[string(preGate)].WireLabel
			} else {
				//The concatentation of all previous labels
				firstInput = append(firstInput, outWires[string(preGate)].WireLabel...)
			}
			//wInCnt++
		}
		//Generate Keys:
		privOutput := GenerateElGamalKey(ByteSliceMul(firstInput, finalInput))
		//A place holder for the output wires
		outWires[string(val.GateID)] = Wire{}
		//Loop over entries
		for garbledKey := range val.GarbledValuesC1 {

			//Decrypt the output label
			outLabel, err := elgamal.Decrypt(privOutput, false, new(big.Int).SetBytes(val.GarbledValuesC1[garbledKey]), new(big.Int).SetBytes(val.GarbledValuesC2[garbledKey]))
			//if decrypted
			if err == nil {
				//Add the wire label to the map to continue the traversal
				outWires[string(val.GateID)] = Wire{
					WireLabel: outLabel,
				}

				//Appending to the result message sent back
				result.Res = append(result.Res, outLabel)
				//break out of looping over entries since you already found the output of the gate
				break
			}
			// outWires[string(val.GateID)] = Wire{
			// 	WireLabel: ByteSliceAdd(outLabel, R),
			// }

		}

		//if no output label was found, then the evaluation failed
		if (bytes.Compare(outWires[string(val.GateID)].WireLabel, Wire{}.WireLabel)) == 0 {
			fmt.Println("Fail Evaluation Output Gate")
		}
		/*else {
			fmt.Println("\n\nYaaay\nGate ", val.GateID, " Now has an output wire: \n", outWires[val.GateID].WireLabel, "\n\n")
		}*/
	}

	return
}

//Convert32BytesToByteStream receives a byte array returns the first 32 bytes from it
func Convert32BytesToByteStream(msg [32]byte) []byte {
	key := make([]byte, 32)
	for jk, tmpo := range msg {
		key[jk] = tmpo
	}
	return key
}

//SHA256Hash Hashes a byte array using sha256
func SHA256Hash(msg []byte) [32]byte {
	return sha256.Sum256(msg)
}

// GetIP getting The IP
func GetIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, err
}

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

//RSAPublicEncrypt encrypts data with a given rsa.publickey
func RSAPublicEncrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(cryptoRand.Reader, key, data)
}

//RSAPrivateDecrypt decrypts encrypted data with a given rsa.privatekey
func RSAPrivateDecrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(cryptoRand.Reader, key, data)
}

//GenerateRSAKey generates Public/Private Key pair, advised rsaKeySize = 2048
func GenerateRSAKey(rsaKeySize int) (*rsa.PrivateKey, *rsa.PublicKey) {
	if rsaKeySize < 1 {
		rsaKeySize = 2048
	}
	pri, err := rsa.GenerateKey(cryptoRand.Reader, rsaKeySize)
	if err != nil {
		panic(err)
	}
	return pri, &pri.PublicKey
}

//RSAPublicKeyFromBytes extracts rsa.publickey from its byte array encoding
func RSAPublicKeyFromBytes(key []byte) *rsa.PublicKey {
	pk, err := x509.ParsePKCS1PublicKey(key)
	if err != nil {
		panic(err)
	}
	return pk
}

//BytesFromRSAPublicKey returns byte array encoding from an rsa.publickey
func BytesFromRSAPublicKey(pk *rsa.PublicKey) []byte {
	pubBytes := x509.MarshalPKCS1PublicKey(pk)
	return pubBytes
}

//BytesFromRSAPrivateKey returns byte array encoding from an rsa.privatekey
func BytesFromRSAPrivateKey(sk *rsa.PrivateKey) []byte {
	priBytes, err := x509.MarshalPKCS8PrivateKey(sk)
	if err != nil {
		panic(err)
	}
	return priBytes
}

//RSAPrivateKeyFromBytes extracts rsa.privatekey from its byte array encoding
func RSAPrivateKeyFromBytes(key []byte) *rsa.PrivateKey {
	pri, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		panic(err)
	}
	p, ok := pri.(*rsa.PrivateKey)
	if !ok {
		panic("Invalid Key type")
	}
	return p
}

//RSAPrivateSign makes a signature with a private key
func RSAPrivateSign(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(cryptoRand.Reader, key, crypto.SHA256, Convert32BytesToByteStream(SHA256Hash(data)))
}

//RSAPrivateVerify verifies a signature made with a private key
func RSAPrivateVerify(key *rsa.PrivateKey, sign, data []byte) error {
	h, err := RSAPrivateDecrypt(key, sign)
	if err != nil {
		return err
	}
	if !bytes.Equal(h, Convert32BytesToByteStream(SHA256Hash(data))) {
		return rsa.ErrVerification
	}
	return nil
}

//RSAPublicSign makes a signature with a public key
func RSAPublicSign(key *rsa.PublicKey, data []byte) ([]byte, error) {
	return RSAPublicEncrypt(key, Convert32BytesToByteStream(SHA256Hash(data)))
}

//RSAPublicVerify verifies a signature made with a public key
func RSAPublicVerify(key *rsa.PublicKey, sign, data []byte) error {
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, Convert32BytesToByteStream(SHA256Hash(data)), sign)
}

//IPtoProperByte puts the IP in its proper formatting
func IPtoProperByte(ip net.IP) []byte {
	var iN0 int = int(ip[0])
	var iN1 int = int(ip[1])
	var iN2 int = int(ip[2])
	var iN3 int = int(ip[3])

	ret := []byte(strconv.Itoa(iN0) + "." + strconv.Itoa(iN1) + "." + strconv.Itoa(iN2) + "." + strconv.Itoa(iN3))

	return ret
}

//GetPartyInfo for a party to extract his own communication info
func GetPartyInfo(username string) (PartyInfo, []byte) {
	port, err := GetFreePort()
	if err != nil {
		panic(err)
	}
	sk, pk := GenerateRSAKey(0)
	if err != nil {
		panic(err)
	}
	ip, err := GetIP()
	if err != nil {
		panic(err)
	}
	pI := PartyInfo{
		IP:        IPtoProperByte(ip),
		Port:      port,
		PublicKey: BytesFromRSAPublicKey(pk),
		UserName:  []byte(username),
	}
	return pI, BytesFromRSAPrivateKey(sk)
}

func fromHex(hex string) *big.Int {
	n, ok := new(big.Int).SetString(hex, 16)
	if !ok {
		panic("failed to parse hex number")
	}
	return n
}

// This is the 1024-bit MODP group from RFC 5114, section 2.1:
const primeHex = "B10B8F96A080E01DDE92DE5EAE5D54EC52C99FBCFB06A3C69A6A9DCA52D23B616073E28675A23D189838EF1E2EE652C013ECB4AEA906112324975C3CD49B83BFACCBDD7D90C4BD7098488E9C219A73724EFFD6FAE5644738FAA31A4FF55BCCC0A151AF5F0DC8B4BD45BF37DF365C1A65E68CFDA76D4DA708DF1FB2BC2E4A4371"

const generatorHex = "A4D1CBD5C3FD34126765A442EFB99905F8104DD258AC507FD6406CFF14266D31266FEA1E5C41564B777E690F5504F213160217B4B01B886A5E91547F9E2749F4D7FBD7D3B9A92EE1909D0D2263F80A76A6A24C087A091F531DBF0A0169B6A28AD662A4D18E73AFA32D779D5918D08BC8858F4DCEF97C2A24855E6EEB22B3B2E5"

//GenerateElGamalKey takes a hexadecimal number and generates a private key
func GenerateElGamalKey(hexS []byte) *elgamal.PrivateKey {

	priv := &elgamal.PrivateKey{
		PublicKey: elgamal.PublicKey{
			G: fromHex(generatorHex),
			P: fromHex(primeHex),
		},
		X: new(big.Int).SetBytes(hexS), //fromHex(hex.EncodeToString(hexS)),
	}

	priv.Y = new(big.Int).Exp(priv.G, priv.X, priv.P)

	return priv
}

//P is a false prime with a probability of 1/4^primeProb
func requestElgamalParam(n int, primeProb int, randomG bool) (P *big.Int, G *big.Int) {
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)

	P = new(big.Int)
	q := new(big.Int)
	flag := false

	// get p and q such that p is prime and q = (p-1)/2 is also prime
	for !flag {
		P, _ = cryptoRand.Prime(cryptoRand.Reader, n)
		q = new(big.Int).Sub(P, one)
		q = new(big.Int).Div(q, two)

		if P.ProbablyPrime(primeProb) && q.ProbablyPrime(primeProb) {
			flag = true
		}
	}

	if randomG {
		PminusOne := new(big.Int).Sub(P, one)

		h, _ := cryptoRand.Int(cryptoRand.Reader, PminusOne)
		exponent := new(big.Int).Div(PminusOne, q)
		G = new(big.Int).Exp(h, exponent, P)

		// G must be > 1 because G shouldn't equal 1 and when compared with zero it returned true vkhbhtk
		for G.Cmp(one) != 1 {
			h, _ = cryptoRand.Int(cryptoRand.Reader, PminusOne)
			exponent = new(big.Int).Div(PminusOne, q)
			G = new(big.Int).Exp(h, exponent, P)
		}
	} else {
		G = new(big.Int).SetInt64(2)
	}
	return
}

func byteSliceXOR(A []byte, B []byte) (C []byte) {
	C = []byte{}
	for key, val := range A {
		C = append(C, val^B[key])
	}
	return
}

//ByteSliceAdd performs addition of (A + B) % P by converting them into Big Ints and storing the byte slice value into X
func ByteSliceAdd(A []byte, B []byte) (X []byte) {
	C := (new(big.Int).Add(new(big.Int).SetBytes(A), new(big.Int).SetBytes(B)))
	C = C.Mod(C, fromHex(primeHex))
	X = C.Bytes()
	return
}

//ByteSliceMul performs multiplication of (A * B) % P by converting them into Big Ints and storing the byte slice value into X
func ByteSliceMul(A []byte, B []byte) (X []byte) {

	C := (new(big.Int).Mul(new(big.Int).SetBytes(A), new(big.Int).SetBytes(B)))
	C = C.Mod(C, fromHex(primeHex))
	X = C.Bytes()
	return
}

//SYS CALL PRINT FUNCTIONS
func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

//UnlockWallet for a node to unlock its eosio wallet
func UnlockWallet(url string, walletKey string) error {
	cmd := exec.Command("cleos", "-u", url, "wallet", "unlock", "--password", walletKey)
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	return err
}

//RegisterOnDecentralizedDS to register in the decentralized Directory of Service
func RegisterOnDecentralizedDS(URL string, ActionAccount string, r RegisterationMessage) error {
	Decentralization = true

	cmd := exec.Command("cleos", "-u", URL, "push", "action", ActionAccount, "login", "[\""+string(r.Server.PartyInfo.UserName)+"\",\""+string(r.Type)+"\",\""+string(r.Server.PartyInfo.IP)+"\",\""+hex.EncodeToString(r.Server.PartyInfo.PublicKey)+"\",\""+strconv.Itoa(r.Server.ServerCapabilities.NumberOfGates)+"\",\""+strconv.Itoa(r.Server.PartyInfo.Port)+"\"]", "-p", string(r.Server.PartyInfo.UserName)+"@active", "-f")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	return err

}

//FetchCycleDecentralized fetches a cycle from the decentralized directory of service
func FetchCycleDecentralized(URL string, ActionAccount string, c CycleRequestMessage) (bool, CycleMessage) {
	cmd := exec.Command("cleos", "-u", URL, "--verbose", "push", "action", ActionAccount, "fetchcycle", "[\""+strconv.Itoa(c.FunctionInfo.NumberOfServers)+"\",\""+strconv.Itoa(c.FunctionInfo.ServerCapabilities.NumberOfGates)+"\"]", "-p", ActionAccount+"@active", "-f")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
	return ConstructCycleStruct(output, c.FunctionInfo.NumberOfServers)
}

//ConstructCycleStruct construct the cycle returned from the transaction
func ConstructCycleStruct(outs []byte, size int) (bool, CycleMessage) {
	var cm CycleMessage
	cm.Cycle.ServersCycle = make([]PartyInfo, size)
	str := string(outs[:])
	if strings.Contains(str, "fetch cycle success") {
		scanner := bufio.NewScanner(strings.NewReader(str))
		i := 0
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), ">>") && !(strings.Contains(scanner.Text(), "fetch cycle success")) {
				s := strings.Split(scanner.Text(), " ")
				cm.Cycle.ServersCycle[i].IP = []byte(s[2])
				cm.Cycle.ServersCycle[i].PublicKey, _ = hex.DecodeString(s[3])
				cm.Cycle.ServersCycle[i].Port, _ = strconv.Atoi(s[4])
				i++
			}
		}
		return true, cm
	}
	return false, cm

}

//ServerRegisterDecentralized registers a client to the decentralized directory of service
func ServerRegisterDecentralized(username string, cleosKey string, numberOfGates int, feePerGate float64) {

	SetMyInfo(username, cleosKey)
	regMsg := RegisterationMessage{
		Type: []byte("Server"),
		Server: ServerInfo{
			PartyInfo: MyOwnInfo.PartyInfo,
			ServerCapabilities: ServerCapabilities{
				NumberOfGates: numberOfGates,
				FeePerGate:    feePerGate,
			},
		},
	}
	err := UnlockWallet(DecentralizedDirectoryInfo.URL, DecentralizedDirectoryInfo.PasswordWallet)

	CreateAccount(DecentralizedDirectoryInfo.URL, regMsg)
	err = RegisterOnDecentralizedDS(DecentralizedDirectoryInfo.URL, DecentralizedDirectoryInfo.ActionAccount, regMsg)
	if err != nil {
		panic(err)
	}

}

//CreateAccount to create a new account in the blockchain.
func CreateAccount(URL string, r RegisterationMessage) {
	cmd := exec.Command("cleos", "-u", URL, "create", "account", "eosio", string(r.Server.PartyInfo.UserName), "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printError(err)
	printOutput(output)
}
