package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	"./vdcs"
)

func main() {
	vdcs.ReadyMutex.Lock()
	vdcs.ReadyFlag = false
	vdcs.ReadyMutex.Unlock()

	username := os.Args[1]
	cleosKey := os.Args[2]
	actionAccount := os.Args[3]
	passwordWallet := os.Args[4]

	vdcs.SetDecentralizedDirectoryInfo("http://127.0.0.1:8888", actionAccount, passwordWallet)
	vdcs.ClientRegisterDecentralized(username, cleosKey)

	go vdcs.ClientHTTP()

	//VDCS Communications:
	_myEqual_string_1_string_1Ch0 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("myEqual_string_1_string_1", 0, 6, 1, _myEqual_string_1_string_1Ch0)
	_myEqual_string_1_string_1Ch1 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("myEqual_string_1_string_1", 1, 6, 1, _myEqual_string_1_string_1Ch1)
	_myEqual_string_1_string_1Ch2 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("myEqual_string_1_string_1", 2, 6, 1, _myEqual_string_1_string_1Ch2)

	//USER PROGRAM:
	i := "abcdefghij"
	j := "aaaaaaaaaa"
	k := "aaaaaaaaaa"
	result := eval(i, j, 0, _myEqual_string_1_string_1Ch0)
	fmt.Println("Result_i==j?:")
	fmt.Println(result)
	result  = eval(i, k, 1, _myEqual_string_1_string_1Ch1)
	fmt.Println("Result_i==k?:")
	fmt.Println(result)
	result  = eval(j, k, 2, _myEqual_string_1_string_1Ch2)
	fmt.Println("Result_j==k?:")
	fmt.Println(result)

}
func eval (a string, b string, cID int64, chVDCSEvalCircRes <-chan vdcs.ChannelContainer) (bool){
   _inWire0:=[]byte(a)
   _inWire1:=[]byte(b)
   //generate input wires for given inputs
   k := <-chVDCSEvalCircRes
    myInWires := make([]vdcs.Wire, len(_inWire0)*8*2)
    for idxByte := 0; idxByte < len(_inWire0); idxByte++ {
      for idxBit := 0; idxBit < 8; idxBit++ {
        contA := (_inWire0[idxByte] >> idxBit) & 1
        myInWires[(idxBit+idxByte*8)*2] = k.InputWires[(idxBit+idxByte*8)*4+int(contA)]
        contB := (_inWire1[idxByte] >> idxBit) & 1
        myInWires[(idxBit+idxByte*8)*2+1] = k.InputWires[(idxBit+idxByte*8)*4+2+int(contB)]
      }
    }
    message := vdcs.Message{
      Type:       []byte("CEval"),
      ComID:      vdcs.ComID{CID: []byte(strconv.FormatInt(cID, 10))},
      InputWires: myInWires,
      NextServer: vdcs.MyOwnInfo.PartyInfo,
    }
    key := vdcs.RandomSymmKeyGen()
    messageEnc := vdcs.EncryptMessageAES(key, message)
    nkey, err := vdcs.RSAPublicEncrypt(vdcs.RSAPublicKeyFromBytes(k.PublicKey), key)
    if err != nil {
      panic("Invalid PublicKey")
    }
    mTmp := make([]vdcs.Message, 1)
    mTmp[0] = messageEnc
    kTmp := make([][]byte, 1)
    kTmp[0] = nkey
    msgArr := vdcs.MessageArray{
      Array: mTmp,
      Keys:  kTmp,
    }
    for ok := vdcs.SendToServer(msgArr, k.IP, k.Port); !ok; {
    }
    var res vdcs.ResEval
    for true {
      vdcs.ReadyMutex.RLock()
      tmpflag := vdcs.ReadyFlag
      vdcs.ReadyMutex.RUnlock()
      if tmpflag == true {
        break
      }
      time.Sleep(1 * time.Second)
    }
    vdcs.ReadyMutex.RLock()
    res = vdcs.MyResult
    vdcs.ReadyMutex.RUnlock()
    vdcs.ReadyMutex.Lock()
    vdcs.ReadyFlag = false
    vdcs.ReadyMutex.Unlock()
    //validate and decode res
    if bytes.Compare(res.Res[0], k.OutputWires[0].WireLabel) == 0 {
      return false
    } else if bytes.Compare(res.Res[0], k.OutputWires[1].WireLabel) == 0 {
      return true
    } else {
      panic("The server cheated me while evaluating")
    }
 }