// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	MATRAND "math/rand"
	"os"

	"./elgamal"
	//"golang.org/x/crypto/openpgp/elgamal"
)

// This is the 1024-bit MODP group from RFC 5114, section 2.1:
const primeHex = "B10B8F96A080E01DDE92DE5EAE5D54EC52C99FBCFB06A3C69A6A9DCA52D23B616073E28675A23D189838EF1E2EE652C013ECB4AEA906112324975C3CD49B83BFACCBDD7D90C4BD7098488E9C219A73724EFFD6FAE5644738FAA31A4FF55BCCC0A151AF5F0DC8B4BD45BF37DF365C1A65E68CFDA76D4DA708DF1FB2BC2E4A4371"

const generatorHex = "A4D1CBD5C3FD34126765A442EFB99905F8104DD258AC507FD6406CFF14266D31266FEA1E5C41564B777E690F5504F213160217B4B01B886A5E91547F9E2749F4D7FBD7D3B9A92EE1909D0D2263F80A76A6A24C087A091F531DBF0A0169B6A28AD662A4D18E73AFA32D779D5918D08BC8858F4DCEF97C2A24855E6EEB22B3B2E5"

func fromHex(hex string) *big.Int {
	n, ok := new(big.Int).SetString(hex, 16)
	if !ok {
		panic("failed to parse hex number")
	}
	return n
}

func GenerateKey(hexS []byte) *elgamal.PrivateKey {

	priv := &elgamal.PrivateKey{
		PublicKey: elgamal.PublicKey{
			G: fromHex(generatorHex),
			P: fromHex(primeHex),
		},
		X: fromHex(hex.EncodeToString(hexS)),
	}

	priv.Y = new(big.Int).Exp(priv.G, priv.X, priv.P)

	return priv
}

func TestEncryptDecrypt() {

	priv := GenerateKey([]byte("a4"))
	//priv1 := GenerateKey("45")
	message := []byte("hello world")

	c1, c2, err := elgamal.Encrypt(rand.Reader, &priv.PublicKey, message)
	y := c1.Bytes()
	x := new(big.Int).SetBytes(y)
	if err != nil {
		fmt.Println("error encrypting: ")
		panic(err)
	}

	message2, err := elgamal.Decrypt(priv, x, c2)
	//message3, err := elgamal.Decrypt(priv1, c1, c2)

	if err != nil {
		fmt.Println("error decrypting: ")
		panic(err)
	}
	if !bytes.Equal(message2, message) {
		panic("decryption failed, got: " + string(message2) + ", want: " + string(message))
	}
	fmt.Println(message)
	fmt.Println(message2)
	//	fmt.Println(message3)
}

func TestDecryptBadKey() {
	priv := &elgamal.PrivateKey{
		PublicKey: elgamal.PublicKey{
			G: fromHex(generatorHex),
			P: fromHex("2"),
		},
		X: fromHex("42"),
	}
	priv.Y = new(big.Int).Exp(priv.G, priv.X, priv.P)
	c1, c2 := fromHex("8"), fromHex("8")
	if _, err := elgamal.Decrypt(priv, c1, c2); err == nil {
		panic("unexpected success decrypting")
	}
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

//GenNRandNumbers generating random byte arrays
func GenNRandNumbers(n int, length int, r int64, considerR bool) [][]byte {
	if considerR {
		MATRAND.Seed(r)
	}
	seeds := make([][]byte, n)
	for i := 0; i < n; i++ {
		seeds[i] = make([]byte, length)
		_, err := MATRAND.Read(seeds[i])
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	return seeds
}

func byteSliceXOR(A []byte, B []byte) (C []byte) {
	C = []byte{}
	for key, val := range A {
		C = append(C, val^B[key])
	}
	return
}

func parseByteToHex(arr []byte) (hexS string) {
	/*hexS = ""
	for _, val := range arr {

		hex += strconv.Itoa((int(val)))
	}*/
	hexS = hex.EncodeToString(arr)
	return
}

func main() {
	//TestEncryptDecrypt()
	//TestDecryptBadKey()

	arrIn := YaoGarbledCkt_in(MATRAND.Int63(), 16, 2)
	arrOut := YaoGarbledCkt_out(MATRAND.Int63(), 16, 1)

	MATRAND.Seed(MATRAND.Int63())
	fmt.Println(arrIn[0])
	fmt.Println((arrIn[1]))

	fmt.Println(parseByteToHex(arrIn[1]))
	fmt.Println("///////")
	fmt.Println(arrOut[0])

	//Generate the mask
	mask := GenNRandNumbers(1, 16, 0, false)

	fmt.Println(string(mask[0]))
	//Generate Keys:
	privMask := GenerateKey(arrIn[0])
	privOutput := GenerateKey(arrIn[1])
	fmt.Println("My little test")
	fmt.Println(privOutput.X.Bytes())
	fmt.Println(arrIn[1])

	maskC1, maskC2, err := elgamal.Encrypt(rand.Reader, &privMask.PublicKey, mask[0])
	tmpSlice := byteSliceXOR(mask[0], arrOut[0])
	if err != nil {
		panic("Mask encryption Error in Input Gates")
	}
	outC1, outC2, err := elgamal.Encrypt(rand.Reader, &privOutput.PublicKey, tmpSlice)
	if err != nil {
		panic("Label encryption Error in Input Gates")
	}

	outMask, err := elgamal.Decrypt(privMask, maskC1, maskC2)
	if err != nil {
		panic("Mask decryption Error in Input Gates")
	}

	outLabelMasked, err := elgamal.Decrypt(privOutput, outC1, outC2)
	if err != nil {
		panic("Label decryption Error in Input Gates")
	}

	trueOutLabel := byteSliceXOR(outMask, outLabelMasked)

	fmt.Println(trueOutLabel)

	if bytes.Compare(trueOutLabel, arrOut[0]) == 0 {
		fmt.Println("Success")
	}
}
