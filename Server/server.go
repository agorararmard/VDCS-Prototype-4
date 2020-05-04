package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"./vdcs"
)

var gm_pendingEval = make(map[string]vdcs.GarbledMessage)
var in_pendingEval = make(map[string]vdcs.GarbledMessage)

var pendingRepo = make(map[string]bool)

var mutexE = sync.RWMutex{}
var op = sync.RWMutex{}


var wg = sync.WaitGroup{}

func main() {

	server()
}

func server() {
	initServer()
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	port := ":" + strconv.Itoa(vdcs.MyOwnInfo.PartyInfo.Port)
	print(port)
	print(vdcs.MyOwnInfo.PartyInfo.PublicKey)
	http.ListenAndServe(port, nil)
}

func initServer() {
	//set whatever to the directory
	username := os.Args[1]
	cleosKey := os.Args[2]
	actionAccount := os.Args[3]
	passwordWallet := os.Args[4]

	vdcs.SetDecentralizedDirectoryInfo("http://127.0.0.1:8888", actionAccount, passwordWallet)

	//register now
	vdcs.ServerRegisterDecentralized(username, cleosKey, 32000, 2)

}

func ServerRegister(numberOfGates int, feePerGate float64) {

	vdcs.SetMyInfo("", "")
	regMsg := vdcs.RegisterationMessage{
		Type: []byte("Server"),
		Server: vdcs.ServerInfo{
			PartyInfo: vdcs.MyOwnInfo.PartyInfo,
			ServerCapabilities: vdcs.ServerCapabilities{
				NumberOfGates: numberOfGates,
				FeePerGate:    feePerGate,
			},
		},
	}
	//fmt.Println(regMsg)
	for !vdcs.SendToDirectory(regMsg, vdcs.DirctoryInfo.IP, vdcs.DirctoryInfo.Port) {
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("I'm solving the token right now!")
	if r.Method == "GET" {
		var x vdcs.Token
		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error reading", err)
		}
		err = json.Unmarshal(jsn, &x)
		if err != nil {
			log.Fatal("bad decode", err)
		}
		ret := vdcs.SolveToken(x)
		vdcs.MyToken = ret
		responseJSON, err := json.Marshal(ret)
		if err != nil {
			fmt.Fprintf(w, "error %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}
}

func handlePostRequest(x vdcs.MessageArray) {
	//fmt.Println("Message Array Size: ", len(x.Array))
	//Decryption
	sk := vdcs.RSAPrivateKeyFromBytes(vdcs.MyOwnInfo.PrivateKey)

	//fmt.Println("symmetric key encrypted for me: ", x.Keys[0])
	k, err := vdcs.RSAPrivateDecrypt(sk, x.Keys[0])
	if err != nil {
		log.Fatal("Error Decrypting the key", err)
	}
	//fmt.Println("symmetric key for me: ", k)
	//fmt.Println("Encrypted Message Type for me:", x.Array[0].Type) //store the keys somewhere for recovery or pass on channel

	//fmt.Println("Encrypted Message: ", x.Array[0])
	//saving it so I won't have to decrypt it again in each thread
	x.Array[0] = vdcs.DecryptMessageAES(k, x.Array[0])
	//fmt.Println("The message to meeee: ", x.Array[0])

	//fmt.Println("The message to meeee from previous server probably: ", x.Array[len(x.Array)-1])

	//checking the type
	reqType := x.Array[0].Type
	////fmt.Println("msg: ", x.Array[0])
	//fmt.Println("reqType: ", string(reqType), x.Array[0].ComID.CID)

	if string(reqType) == "Garble" {
		//the garbling thread
		go garbleLogic(x)

	} else if string(reqType) == "ReRand" {
		//the rerand thread
		go rerandLogic(x)
	} else if string(reqType) == "SEval" {
		//the eval thread
		k, err := vdcs.RSAPrivateDecrypt(sk, x.Keys[1])
		if err != nil {
			log.Fatal("Error Decrypting the key", err)
		}
		x.Array[1] = vdcs.DecryptMessageAES(k, x.Array[1])

		msg := vdcs.Message{
			Type:           x.Array[0].Type,
			GarbledMessage: x.Array[1].GarbledMessage,
			ComID:          x.Array[0].ComID,
			NextServer:     x.Array[0].NextServer,
		}
		//fmt.Println("IP of the client before Handling logic: ", string(x.Array[0].NextServer.IP))
		//fmt.Println("Public key of the client before Handling logic: ", x.Array[0].NextServer.PublicKey)
		//fmt.Println("Length of the array: ", len(x.Array), "Length of the keys: ", len(x.Keys))
		go evalLogic(msg, string(reqType))
	} else if string(reqType) == "CEval" {
		//the thread for the client requesting the result
		go evalLogic(x.Array[0], string(reqType))
	}

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Post is invoked server!")
	if r.Method == "POST" {

		//getting the array of messages
		var x vdcs.MessageArray
		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal("Error reading", err)
		}
		err = json.Unmarshal(jsn, &x)
		if err != nil {
			log.Fatal("bad decode", err)
		}
		go handlePostRequest(x)
	}
}

//ServerRegister registers a client to directory of service
//should work fine after gouhar fix the issue of naming

func garbleLogic(arr vdcs.MessageArray) {
	op.Lock()

	//access the first message which is aleardy decrypted in the post handler
	x := arr.Array[0]
	//fmt.Println("first message: ", x)
	//create the circuit message to garble it
	circM := vdcs.CircuitMessage{
		Circuit: vdcs.Circuit{
			InputGates:  x.Circuit.InputGates,
			MiddleGates: x.Circuit.MiddleGates,
			OutputGates: x.Circuit.OutputGates,
		},
		ComID: vdcs.ComID{
			CID: x.CID,
		},
		Randomness: vdcs.Randomness{
			Rin:       x.Rin,
			Rout:      x.Rout,
			Rgc:       x.Rgc,
			LblLength: x.LblLength,
		},
	}

	//garbling
	gm := vdcs.Garble(circM)

	//the message to be sent
	mess := vdcs.Message{
		GarbledMessage: gm,
		ComID: vdcs.ComID{
			CID: x.CID,
		},
	}
	//fmt.Println("next server: ", mess.NextServer)

	//removing the first one
	arr.Array = append(arr.Array[:0], arr.Array[1:]...)
	//fmt.Println("message array: ", arr.Array)

	//appending the new message
	arr.Array = append(arr.Array, mess)
	//fmt.Println("message array: ", arr.Array)

	//setting the request type for the next one
	if len(arr.Array) > 2 {
		arr.Array[len(arr.Array)-1].Type = []byte("ReRand")
	} else {
		arr.Array[len(arr.Array)-1].Type = []byte("SEval")
	}

	//encrypting the message by generating a new key first then using it
	k := vdcs.RandomSymmKeyGen()
	arr.Array[len(arr.Array)-1] = vdcs.EncryptMessageAES(k, arr.Array[len(arr.Array)-1])

	//encreypting the key used in previous line
	pk := vdcs.RSAPublicKeyFromBytes(x.NextServer.PublicKey)
	key, err := vdcs.RSAPublicEncrypt(pk, k)
	if err != nil {
		log.Fatal("Error in decrypting", err)
	}
	//appending the new key
	arr.Keys = append(arr.Keys, key)
	//removing the first one
	arr.Keys = append(arr.Keys[:0], arr.Keys[1:]...)

	//send to the next server
	//fmt.Println("Next Server IP ", string(x.NextServer.IP))
	//fmt.Println("Next Server port", x.NextServer.Port)

	//fmt.Println("input gates length sent: ", len(mess.GarbledMessage.InputGates))
	//fmt.Println("middle gates length sent: ", len(mess.GarbledMessage.MiddleGates))
	//fmt.Println("output gates length sent: ", len(mess.GarbledMessage.OutputGates))

	//fmt.Println("input gates length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.InputGates))
	//fmt.Println("middle gates length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.MiddleGates))
	//fmt.Println("output gates length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.OutputGates))

	//fmt.Println("input wires length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.InputWires))
	//fmt.Println("output wires length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.OutputWires))

	vdcs.SendToServer(arr, x.NextServer.IP, x.NextServer.Port)
	op.Unlock()
}

func rerandLogic(arr vdcs.MessageArray) {
	op.Lock()
	//fmt.Println("Message Array Size inside rerandLogic: ", len(arr.Array))
	//fmt.Println("The message sent to me encrypted: ", arr.Array[len(arr.Array)-1])
	//fmt.Println("input gates length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.InputGates))
	//fmt.Println("middle gates length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.MiddleGates))
	//fmt.Println("output gates length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.OutputGates))

	//fmt.Println("input wires length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.InputWires))
	//fmt.Println("output wires length sent encrypted: ", len(arr.Array[len(arr.Array)-1].GarbledMessage.OutputWires))

	//fmt.Println("CID before encrypt: ", arr.Array[len(arr.Array)-1].ComID)
	//the first message already decrypted
	//fmt.Println("input gates length sent encrypted first message: ", len(arr.Array[0].GarbledMessage.InputGates))
	//fmt.Println("middle gates length sent encrypted first message: ", len(arr.Array[0].GarbledMessage.MiddleGates))
	//fmt.Println("output gates length sent encrypted first message: ", len(arr.Array[0].GarbledMessage.OutputGates))

	//the info for the next server
	// variable-array consistency potential problem
	next := arr.Array[0].NextServer

	//get the last & first Message
	x0 := arr.Array[0] //Already decrypted earlier
	x1 := arr.Array[len(arr.Array)-1]

	//fmt.Println("input gates length received: ", len(x1.GarbledMessage.InputGates))
	//fmt.Println("middle gates length received: ", len(x1.GarbledMessage.MiddleGates))
	//fmt.Println("output gates length received: ", len(x1.GarbledMessage.OutputGates))

	//decrypting x1
	sk := vdcs.RSAPrivateKeyFromBytes(vdcs.MyOwnInfo.PrivateKey)
	k, err := vdcs.RSAPrivateDecrypt(sk, arr.Keys[len(arr.Keys)-1])
	if err != nil {
		log.Fatal("Error Decrypting the key", err)
	}
	x1 = vdcs.DecryptMessageAES(k, x1) //now x1& x0 are decrypted
	//fmt.Println("The message sent to me decrypted: ", x1)
	//fmt.Println("Type after decrypt: ", string(x1.Type))

	//fmt.Println("CID after decrypt: ", x1.ComID)

	// getting the garble message from x1 and the nextserver from x0
	//Forming a single message out of the two to work with
	mess := vdcs.Message{
		//from x1
		GarbledMessage: vdcs.GarbledMessage{
			InputWires: x1.GarbledMessage.InputWires,
			GarbledCircuit: vdcs.GarbledCircuit{
				InputGates:  x1.GarbledCircuit.InputGates,
				MiddleGates: x1.GarbledCircuit.MiddleGates,
				OutputGates: x1.GarbledCircuit.OutputGates,
				ComID: vdcs.ComID{
					CID: x1.ComID.CID,
				},
			},
			OutputWires: x1.GarbledMessage.OutputWires,
		},

		//from x0
		NextServer: x0.NextServer,

		//from x0
		Randomness: vdcs.Randomness{
			Rin:       x0.Rin,
			Rout:      x0.Rout,
			Rgc:       x0.Rgc,
			LblLength: x0.LblLength,
		},

		ComID: vdcs.ComID{
			CID: x0.CID,
		},
	}
	//fmt.Println("input gates length sent before rerand: ", len(mess.GarbledMessage.InputGates))
	//fmt.Println("middle gates length sent before rerand: ", len(mess.GarbledMessage.MiddleGates))
	//fmt.Println("output gates length sent before rerand: ", len(mess.GarbledMessage.OutputGates))

	//reranding the garbled circuit
	ngcm := vdcs.ReRand(mess.GarbledMessage, mess.Randomness)
	//newMessage to append
	nMessage := vdcs.Message{
		//from reRand
		GarbledMessage: ngcm,

		ComID: vdcs.ComID{
			CID: mess.CID,
		},
	}

	//fmt.Println("input gates length sent after rerand: ", len(nMessage.GarbledMessage.InputGates))
	//fmt.Println("middle gates length sent after rerand: ", len(nMessage.GarbledMessage.MiddleGates))
	//fmt.Println("output gates length sent after rerand: ", len(nMessage.GarbledMessage.OutputGates))

	//removing the first one
	arr.Array = append(arr.Array[:0], arr.Array[1:]...)
	//remove the last one
	arr.Array = arr.Array[:len(arr.Array)-1]
	//appending the new message
	arr.Array = append(arr.Array, nMessage)
	//setting the type of the new message
	if len(arr.Array) > 2 {
		arr.Array[len(arr.Array)-1].Type = []byte("ReRand")
	} else {
		arr.Array[len(arr.Array)-1].Type = []byte("SEval")
	}

	//encrypting the message by generating a new key first then using it
	kn := vdcs.RandomSymmKeyGen()
	arr.Array[len(arr.Array)-1] = vdcs.EncryptMessageAES(kn, arr.Array[len(arr.Array)-1])

	//encreypting the key used in previous line using public key
	pk := vdcs.RSAPublicKeyFromBytes(mess.NextServer.PublicKey)
	key, err := vdcs.RSAPublicEncrypt(pk, kn)
	if err != nil {
		log.Fatal("Error in decrypting", err)
	}

	//removing the first one
	arr.Keys = append(arr.Keys[:0], arr.Keys[1:]...)
	//remove the last one
	arr.Keys = arr.Keys[:len(arr.Keys)-1]
	//appending the new message
	arr.Keys = append(arr.Keys, key)

	//send to the next server
	//fmt.Println("Next Server IP ", string(next.IP))
	//fmt.Println("Next Server port", next.Port)
	//fmt.Println("input gates length sent: ", len(nMessage.GarbledMessage.InputGates))
	//fmt.Println("middle gates length sent: ", len(nMessage.GarbledMessage.MiddleGates))
	//fmt.Println("output gates length sent: ", len(nMessage.GarbledMessage.OutputGates))

	//send it to the next server.... (from the first message)
	vdcs.SendToServer(arr, next.IP, next.Port)
	op.Unlock()
}

func evalLogic(mess vdcs.Message, reqType string) {
	op.Lock()
	//the first one is already decrypted
	gm := vdcs.GarbledMessage{

		InputWires: mess.InputWires,

		OutputWires: mess.OutputWires,

		GarbledCircuit: vdcs.GarbledCircuit{

			InputGates: mess.GarbledCircuit.InputGates,

			OutputGates: mess.GarbledCircuit.OutputGates,

			MiddleGates: mess.GarbledCircuit.MiddleGates,

			ComID: vdcs.ComID{
				CID: mess.ComID.CID,
			},
		},
	}

	//check whether this ComID have any pending wires OR Circuts
	mutexE.Lock()
	if _, ok := pendingRepo[string(gm.CID)]; ok {

		if reqType == "SEval" {

			evalGm := vdcs.GarbledMessage{

				InputWires: in_pendingEval[string(gm.CID)].InputWires,

				OutputWires: gm.OutputWires,

				GarbledCircuit: vdcs.GarbledCircuit{

					InputGates: gm.GarbledCircuit.InputGates,

					OutputGates: gm.GarbledCircuit.OutputGates,

					MiddleGates: gm.GarbledCircuit.MiddleGates,

					ComID: vdcs.ComID{
						CID: gm.CID,
					},
				},
			}

			//remove the pending from the map
			delete(pendingRepo, string(gm.CID))
			delete(in_pendingEval, string(gm.CID))
			mutexE.Unlock()

			//send them
			res := vdcs.Evaluate(evalGm)

			//send to the client
			//fmt.Println("Next Server IP ", string(mess.NextServer.IP))
			//fmt.Println("Next Server port", mess.NextServer.Port)
			//fmt.Println("result sent is: ", res)
			//fmt.Println("input wires received: ", evalGm.InputWires)

			//fmt.Println("input gates length received: ", len(evalGm.InputGates))
			//fmt.Println("middle gates length received: ", len(evalGm.MiddleGates))
			//fmt.Println("output gates length received: ", len(evalGm.OutputGates))

			vdcs.SendToClient(res, mess.NextServer.IP, mess.NextServer.Port)

			//if the client send the input wires
		} else {

			evalGm := vdcs.GarbledMessage{
				InputWires:  gm.InputWires,
				OutputWires: gm_pendingEval[string(gm.CID)].OutputWires,
				GarbledCircuit: vdcs.GarbledCircuit{
					InputGates:  gm_pendingEval[string(gm.CID)].InputGates,
					OutputGates: gm_pendingEval[string(gm.CID)].OutputGates,
					MiddleGates: gm_pendingEval[string(gm.CID)].MiddleGates,
					ComID: vdcs.ComID{
						CID: gm_pendingEval[string(gm.CID)].CID,
					},
				},
			}

			delete(pendingRepo, string(gm.CID))
			delete(gm_pendingEval, string(gm.CID))
			mutexE.Unlock()

			//send them
			res := vdcs.Evaluate(evalGm)

			//send to the client
			//fmt.Println("Next Server IP ", string(mess.NextServer.IP))
			//fmt.Println("Next Server port", mess.NextServer.Port)
			//fmt.Println("result sent is: ", res)
			//fmt.Println("input wires received: ", evalGm.InputWires)

			//fmt.Println("input gates length received: ", len(evalGm.InputGates))
			//fmt.Println("middle gates length received: ", len(evalGm.MiddleGates))
			//fmt.Println("output gates length received: ", len(evalGm.OutputGates))

			vdcs.SendToClient(res, mess.NextServer.IP, mess.NextServer.Port)

		}

	} else {
		// cid potential problem
		pendingRepo[string(gm.CID)] = true

		if reqType == "SEval" {
			gm_pendingEval[string(gm.CID)] = gm
		} else {
			in_pendingEval[string(gm.CID)] = gm
		}
		mutexE.Unlock()
	}
	op.Unlock()
}
