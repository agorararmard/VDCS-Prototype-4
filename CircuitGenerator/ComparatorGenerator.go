package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

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

//Xnor gate definition for an XNOR gate
var Xnor localcircuitgate = localcircuitgate{
	TruthTable: []bool{true, false, false, true},
}

//And gate definition for an AND gate
var And localcircuitgate = localcircuitgate{
	TruthTable: []bool{false, false, false, true},
}

func main() {
	bytesNumber, err := strconv.ParseInt(os.Args[1], 10, 32)
	if err != nil {
		panic("Invalid parameter for number of bytes")
	}

	var queue []string
	//Input Gates
	var generatedCircuit localcircuit
	gateCounter := 0
	for ; gateCounter < int(bytesNumber)*8; gateCounter++ {
		tmpGate := localcircuitgate{
			localgate: localgate{
				GateID: strconv.Itoa(gateCounter),
			},
			TruthTable: Xnor.TruthTable,
		}
		queue = append(queue, tmpGate.GateID)
		//		fmt.Println(tmpGate)
		generatedCircuit.InputGates = append(generatedCircuit.InputGates, tmpGate)

	}

	//Middle Gates
	for len(queue) > 2 {
		tmpGate := localcircuitgate{
			localgate: localgate{
				GateID:     strconv.Itoa(gateCounter),
				GateInputs: []string{queue[0], queue[1]},
			},
			TruthTable: And.TruthTable,
		}
		queue = queue[2:]
		queue = append(queue, tmpGate.GateID)
		//		fmt.Println(tmpGate)
		generatedCircuit.MiddleGates = append(generatedCircuit.MiddleGates, tmpGate)
		gateCounter++
	}

	//OutputGates
	tmpGate := localcircuitgate{
		localgate: localgate{
			GateID:     strconv.Itoa(gateCounter),
			GateInputs: []string{queue[0], queue[1]},
		},
		TruthTable: And.TruthTable,
	}
	queue = queue[2:]
	queue = append(queue, tmpGate.GateID)
	generatedCircuit.OutputGates = append(generatedCircuit.OutputGates, tmpGate)
	gateCounter++

	//fmt.Println(queue)

	//fmt.Println(generatedCircuit)

	file, _ := json.Marshal(generatedCircuit)

	_ = ioutil.WriteFile("myEqual_string_1_string_1.json", file, 0644)
}
