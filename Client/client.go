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
	_VDCSEqual1Ch0 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 0, 4, 1, _VDCSEqual1Ch0)
	_VDCSEqual1Ch1 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 1, 4, 1, _VDCSEqual1Ch1)
	_VDCSEqual1Ch2 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 2, 4, 1, _VDCSEqual1Ch2)
	_VDCSEqual1Ch3 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 3, 4, 1, _VDCSEqual1Ch3)
	_VDCSEqual1Ch4 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 4, 4, 1, _VDCSEqual1Ch4)
	_VDCSEqual1Ch5 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 5, 4, 1, _VDCSEqual1Ch5)
	_VDCSEqual1Ch6 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 6, 4, 1, _VDCSEqual1Ch6)
	_VDCSEqual1Ch7 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 7, 4, 1, _VDCSEqual1Ch7)
	_VDCSEqual1Ch8 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 8, 4, 1, _VDCSEqual1Ch8)
	_VDCSEqual1Ch9 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 9, 4, 1, _VDCSEqual1Ch9)
	_VDCSEqual1Ch10 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 10, 4, 1, _VDCSEqual1Ch10)
	_VDCSEqual1Ch11 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 11, 4, 1, _VDCSEqual1Ch11)
	_VDCSEqual1Ch12 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 12, 4, 1, _VDCSEqual1Ch12)
	_VDCSEqual1Ch13 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 13, 4, 1, _VDCSEqual1Ch13)
	_VDCSEqual1Ch14 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 14, 4, 1, _VDCSEqual1Ch14)
	_VDCSEqual1Ch15 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 15, 4, 1, _VDCSEqual1Ch15)
	_VDCSEqual1Ch16 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 16, 4, 1, _VDCSEqual1Ch16)
	_VDCSEqual1Ch17 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 17, 4, 1, _VDCSEqual1Ch17)
	_VDCSEqual1Ch18 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 18, 4, 1, _VDCSEqual1Ch18)
	_VDCSEqual1Ch19 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 19, 4, 1, _VDCSEqual1Ch19)
	_VDCSEqual1Ch20 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 20, 4, 1, _VDCSEqual1Ch20)
	_VDCSEqual1Ch21 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 21, 4, 1, _VDCSEqual1Ch21)
	_VDCSEqual1Ch22 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 22, 4, 1, _VDCSEqual1Ch22)
	_VDCSEqual1Ch23 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 23, 4, 1, _VDCSEqual1Ch23)
	_VDCSEqual1Ch24 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 24, 4, 1, _VDCSEqual1Ch24)
	_VDCSEqual1Ch25 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 25, 4, 1, _VDCSEqual1Ch25)
	_VDCSEqual1Ch26 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 26, 4, 1, _VDCSEqual1Ch26)
	_VDCSEqual1Ch27 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 27, 4, 1, _VDCSEqual1Ch27)
	_VDCSEqual1Ch28 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 28, 4, 1, _VDCSEqual1Ch28)
	_VDCSEqual1Ch29 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 29, 4, 1, _VDCSEqual1Ch29)
	_VDCSEqual1Ch30 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 30, 4, 1, _VDCSEqual1Ch30)
	_VDCSEqual1Ch31 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 31, 4, 1, _VDCSEqual1Ch31)
	_VDCSEqual1Ch32 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 32, 4, 1, _VDCSEqual1Ch32)
	_VDCSEqual1Ch33 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 33, 4, 1, _VDCSEqual1Ch33)
	_VDCSEqual1Ch34 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 34, 4, 1, _VDCSEqual1Ch34)
	_VDCSEqual1Ch35 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 35, 4, 1, _VDCSEqual1Ch35)
	_VDCSEqual1Ch36 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 36, 4, 1, _VDCSEqual1Ch36)
	_VDCSEqual1Ch37 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 37, 4, 1, _VDCSEqual1Ch37)
	_VDCSEqual1Ch38 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 38, 4, 1, _VDCSEqual1Ch38)
	_VDCSEqual1Ch39 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 39, 4, 1, _VDCSEqual1Ch39)
	_VDCSEqual1Ch40 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 40, 4, 1, _VDCSEqual1Ch40)
	_VDCSEqual1Ch41 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 41, 4, 1, _VDCSEqual1Ch41)
	_VDCSEqual1Ch42 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 42, 4, 1, _VDCSEqual1Ch42)
	_VDCSEqual1Ch43 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 43, 4, 1, _VDCSEqual1Ch43)
	_VDCSEqual1Ch44 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 44, 4, 1, _VDCSEqual1Ch44)
	_VDCSEqual1Ch45 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 45, 4, 1, _VDCSEqual1Ch45)
	_VDCSEqual1Ch46 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 46, 4, 1, _VDCSEqual1Ch46)
	_VDCSEqual1Ch47 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 47, 4, 1, _VDCSEqual1Ch47)
	_VDCSEqual1Ch48 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 48, 4, 1, _VDCSEqual1Ch48)
	_VDCSEqual1Ch49 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 49, 4, 1, _VDCSEqual1Ch49)
	_VDCSEqual1Ch50 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 50, 4, 1, _VDCSEqual1Ch50)
	_VDCSEqual1Ch51 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 51, 4, 1, _VDCSEqual1Ch51)
	_VDCSEqual1Ch52 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 52, 4, 1, _VDCSEqual1Ch52)
	_VDCSEqual1Ch53 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 53, 4, 1, _VDCSEqual1Ch53)
	_VDCSEqual1Ch54 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 54, 4, 1, _VDCSEqual1Ch54)
	_VDCSEqual1Ch55 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 55, 4, 1, _VDCSEqual1Ch55)
	_VDCSEqual1Ch56 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 56, 4, 1, _VDCSEqual1Ch56)
	_VDCSEqual1Ch57 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 57, 4, 1, _VDCSEqual1Ch57)
	_VDCSEqual1Ch58 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 58, 4, 1, _VDCSEqual1Ch58)
	_VDCSEqual1Ch59 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 59, 4, 1, _VDCSEqual1Ch59)
	_VDCSEqual1Ch60 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 60, 4, 1, _VDCSEqual1Ch60)
	_VDCSEqual1Ch61 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 61, 4, 1, _VDCSEqual1Ch61)
	_VDCSEqual1Ch62 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 62, 4, 1, _VDCSEqual1Ch62)
	_VDCSEqual1Ch63 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 63, 4, 1, _VDCSEqual1Ch63)
	_VDCSEqual1Ch64 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 64, 4, 1, _VDCSEqual1Ch64)
	_VDCSEqual1Ch65 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 65, 4, 1, _VDCSEqual1Ch65)
	_VDCSEqual1Ch66 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 66, 4, 1, _VDCSEqual1Ch66)
	_VDCSEqual1Ch67 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 67, 4, 1, _VDCSEqual1Ch67)
	_VDCSEqual1Ch68 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 68, 4, 1, _VDCSEqual1Ch68)
	_VDCSEqual1Ch69 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 69, 4, 1, _VDCSEqual1Ch69)
	_VDCSEqual1Ch70 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 70, 4, 1, _VDCSEqual1Ch70)
	_VDCSEqual1Ch71 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 71, 4, 1, _VDCSEqual1Ch71)
	_VDCSEqual1Ch72 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 72, 4, 1, _VDCSEqual1Ch72)
	_VDCSEqual1Ch73 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 73, 4, 1, _VDCSEqual1Ch73)
	_VDCSEqual1Ch74 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 74, 4, 1, _VDCSEqual1Ch74)
	_VDCSEqual1Ch75 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 75, 4, 1, _VDCSEqual1Ch75)
	_VDCSEqual1Ch76 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 76, 4, 1, _VDCSEqual1Ch76)
	_VDCSEqual1Ch77 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 77, 4, 1, _VDCSEqual1Ch77)
	_VDCSEqual1Ch78 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 78, 4, 1, _VDCSEqual1Ch78)
	_VDCSEqual1Ch79 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 79, 4, 1, _VDCSEqual1Ch79)
	_VDCSEqual1Ch80 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 80, 4, 1, _VDCSEqual1Ch80)
	_VDCSEqual1Ch81 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 81, 4, 1, _VDCSEqual1Ch81)
	_VDCSEqual1Ch82 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 82, 4, 1, _VDCSEqual1Ch82)
	_VDCSEqual1Ch83 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 83, 4, 1, _VDCSEqual1Ch83)
	_VDCSEqual1Ch84 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 84, 4, 1, _VDCSEqual1Ch84)
	_VDCSEqual1Ch85 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 85, 4, 1, _VDCSEqual1Ch85)
	_VDCSEqual1Ch86 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 86, 4, 1, _VDCSEqual1Ch86)
	_VDCSEqual1Ch87 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 87, 4, 1, _VDCSEqual1Ch87)
	_VDCSEqual1Ch88 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 88, 4, 1, _VDCSEqual1Ch88)
	_VDCSEqual1Ch89 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 89, 4, 1, _VDCSEqual1Ch89)
	_VDCSEqual1Ch90 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 90, 4, 1, _VDCSEqual1Ch90)
	_VDCSEqual1Ch91 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 91, 4, 1, _VDCSEqual1Ch91)
	_VDCSEqual1Ch92 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 92, 4, 1, _VDCSEqual1Ch92)
	_VDCSEqual1Ch93 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 93, 4, 1, _VDCSEqual1Ch93)
	_VDCSEqual1Ch94 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 94, 4, 1, _VDCSEqual1Ch94)
	_VDCSEqual1Ch95 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 95, 4, 1, _VDCSEqual1Ch95)
	_VDCSEqual1Ch96 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 96, 4, 1, _VDCSEqual1Ch96)
	_VDCSEqual1Ch97 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 97, 4, 1, _VDCSEqual1Ch97)
	_VDCSEqual1Ch98 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 98, 4, 1, _VDCSEqual1Ch98)
	_VDCSEqual1Ch99 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual1", 99, 4, 1, _VDCSEqual1Ch99)
	_VDCSEqual2kCh100 := make(chan vdcs.ChannelContainer)
   	go vdcs.Comm("VDCSEqual2k", 100, 4, 1, _VDCSEqual2kCh100)

	//USER PROGRAM:
	a := "a"
	b := "b"
	idx := 0
	result1 := eval(a, b, 0, _VDCSEqual1Ch0)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 1, _VDCSEqual1Ch1)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 2, _VDCSEqual1Ch2)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 3, _VDCSEqual1Ch3)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 4, _VDCSEqual1Ch4)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 5, _VDCSEqual1Ch5)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 6, _VDCSEqual1Ch6)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 7, _VDCSEqual1Ch7)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 8, _VDCSEqual1Ch8)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 9, _VDCSEqual1Ch9)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 10, _VDCSEqual1Ch10)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 11, _VDCSEqual1Ch11)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 12, _VDCSEqual1Ch12)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 13, _VDCSEqual1Ch13)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 14, _VDCSEqual1Ch14)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 15, _VDCSEqual1Ch15)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 16, _VDCSEqual1Ch16)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 17, _VDCSEqual1Ch17)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 18, _VDCSEqual1Ch18)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 19, _VDCSEqual1Ch19)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 20, _VDCSEqual1Ch20)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 21, _VDCSEqual1Ch21)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 22, _VDCSEqual1Ch22)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 23, _VDCSEqual1Ch23)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 24, _VDCSEqual1Ch24)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 25, _VDCSEqual1Ch25)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 26, _VDCSEqual1Ch26)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 27, _VDCSEqual1Ch27)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 28, _VDCSEqual1Ch28)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 29, _VDCSEqual1Ch29)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 30, _VDCSEqual1Ch30)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 31, _VDCSEqual1Ch31)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 32, _VDCSEqual1Ch32)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 33, _VDCSEqual1Ch33)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 34, _VDCSEqual1Ch34)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 35, _VDCSEqual1Ch35)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 36, _VDCSEqual1Ch36)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 37, _VDCSEqual1Ch37)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 38, _VDCSEqual1Ch38)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 39, _VDCSEqual1Ch39)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 40, _VDCSEqual1Ch40)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 41, _VDCSEqual1Ch41)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 42, _VDCSEqual1Ch42)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 43, _VDCSEqual1Ch43)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 44, _VDCSEqual1Ch44)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 45, _VDCSEqual1Ch45)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 46, _VDCSEqual1Ch46)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 47, _VDCSEqual1Ch47)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 48, _VDCSEqual1Ch48)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 49, _VDCSEqual1Ch49)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 50, _VDCSEqual1Ch50)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 51, _VDCSEqual1Ch51)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 52, _VDCSEqual1Ch52)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 53, _VDCSEqual1Ch53)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 54, _VDCSEqual1Ch54)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 55, _VDCSEqual1Ch55)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 56, _VDCSEqual1Ch56)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 57, _VDCSEqual1Ch57)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 58, _VDCSEqual1Ch58)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 59, _VDCSEqual1Ch59)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 60, _VDCSEqual1Ch60)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 61, _VDCSEqual1Ch61)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 62, _VDCSEqual1Ch62)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 63, _VDCSEqual1Ch63)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 64, _VDCSEqual1Ch64)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 65, _VDCSEqual1Ch65)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 66, _VDCSEqual1Ch66)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 67, _VDCSEqual1Ch67)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 68, _VDCSEqual1Ch68)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 69, _VDCSEqual1Ch69)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 70, _VDCSEqual1Ch70)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 71, _VDCSEqual1Ch71)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 72, _VDCSEqual1Ch72)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 73, _VDCSEqual1Ch73)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 74, _VDCSEqual1Ch74)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 75, _VDCSEqual1Ch75)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 76, _VDCSEqual1Ch76)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 77, _VDCSEqual1Ch77)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 78, _VDCSEqual1Ch78)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 79, _VDCSEqual1Ch79)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 80, _VDCSEqual1Ch80)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 81, _VDCSEqual1Ch81)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 82, _VDCSEqual1Ch82)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 83, _VDCSEqual1Ch83)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 84, _VDCSEqual1Ch84)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 85, _VDCSEqual1Ch85)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 86, _VDCSEqual1Ch86)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 87, _VDCSEqual1Ch87)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 88, _VDCSEqual1Ch88)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 89, _VDCSEqual1Ch89)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 90, _VDCSEqual1Ch90)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 91, _VDCSEqual1Ch91)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 92, _VDCSEqual1Ch92)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 93, _VDCSEqual1Ch93)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 94, _VDCSEqual1Ch94)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 95, _VDCSEqual1Ch95)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 96, _VDCSEqual1Ch96)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 97, _VDCSEqual1Ch97)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 98, _VDCSEqual1Ch98)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	result1  = eval(a, b, 99, _VDCSEqual1Ch99)
	fmt.Println("EvaluationRequest:")
	fmt.Println(idx)
	if result1 {
	fmt.Println("VDCS_RESULT:a==b")
	} else {
	fmt.Println("VDCS_RESULT:a!=b")
	}
	idx  = idx+1
	longA := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	longB := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbc"
	fmt.Println("===========================================================")
	result2 := eval(longA, longB, 100, _VDCSEqual2kCh100)
	fmt.Println("FinalEvaluationRequest-2kB:")
	fmt.Println(result2)

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