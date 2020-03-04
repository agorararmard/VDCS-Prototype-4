package main

import (
	"fmt"
)

func main() {
	var i string = "a"
	var j string = "a"
	//VDCS
	if myEqual(i, j) == true {
		fmt.Println("i == j")
	} else {
		fmt.Println("i != j")
	}

	var z string = "b"
	//VDCS
	if myEqual(i, z) == true {
		fmt.Println("i == z")
	} else {
		fmt.Println("i != z")
	}
}
