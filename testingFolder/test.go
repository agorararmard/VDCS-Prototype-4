package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"strconv"
)

type PartyInfo struct {
	IP        []byte `json:"IP"`
	Port      int    `json:"Port"`
	PublicKey []byte `json:"PublicKey"`
}

func EncryptAES(encKey []byte, plainText []byte) (ciphertext []byte, ok bool) {

	ok = false //assume failure
	//			encKey = append(encKey, hash)
	c, err := aes.NewCipher(encKey)
	if err != nil {
		//fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		//fmt.Println(err)
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
		//fmt.Println(err)
		return
	}
	ciphertext = gcm.Seal(nonce, nonce, plainText, nil)
	//fmt.Println(ciphertext)
	ok = true

	return
}

func DecryptAES(encKey []byte, cipherText []byte) (plainText []byte, ok bool) {

	ok = false //assume failure

	c, err := aes.NewCipher(encKey)
	if err != nil {
		//fmt.Println(err)
		return
	}

	gcm, err := cipher.NewGCM(c)
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
	plainText, err = gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Println(string(plaintext))
	ok = true
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
	return rsa.EncryptPKCS1v15(rand.Reader, key, data)
}

//RSAPrivateDecrypt decrypts encrypted data with a given rsa.privatekey
func RSAPrivateDecrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key, data)
}

//GenerateRSAKey generates Public/Private Key pair, advised rsaKeySize = 2048
func GenerateRSAKey(rsaKeySize int) (*rsa.PrivateKey, *rsa.PublicKey) {
	if rsaKeySize < 1 {
		rsaKeySize = 2048
	}
	pri, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
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
	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, Convert32BytesToByteStream(SHA256Hash(data)))
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

//GetPartyInfo for a party to extract his own communication info
func GetPartyInfo() (PartyInfo, []byte) {
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
		IP:        ip,
		Port:      port,
		PublicKey: BytesFromRSAPublicKey(pk),
	}
	return pI, BytesFromRSAPrivateKey(sk)
}

func main() {
	//ip := GetMyIP()
	PI, skb := GetPartyInfo()
	var ip net.IP
	ip = PI.IP
	fmt.Println(ip)
	pk := RSAPublicKeyFromBytes(PI.PublicKey)
	sk := RSAPrivateKeyFromBytes(skb)
	//	fmt.Println(pk)
	//fmt.Println(sk)
	msg1 := "Hi sala7"
	hash := SHA256Hash([]byte("This is my key"))
	key := Convert32BytesToByteStream(hash)
	/*key := make([]byte, 32)
	for jk, tmpo := range hash {
		key[jk] = tmpo
	}*/
	//return key
	cipherText, _ := EncryptAES(key, []byte(msg1))

	cipher1, _ := RSAPublicEncrypt(pk, key)
	key1, _ := RSAPrivateDecrypt(sk, cipher1)

	plainText, _ := DecryptAES(key1, cipherText)
	fmt.Println(string(plainText))

	signature, _ := RSAPrivateSign(sk, []byte(msg1))
	err := RSAPublicVerify(pk, signature, []byte(msg1))
	if err == nil {
		fmt.Println("Yes it's his signature")
	}
	signature2, _ := RSAPublicSign(pk, []byte(msg1))
	err = RSAPrivateVerify(sk, signature2, []byte(msg1))
	if err == nil {
		fmt.Println("Yes it's his signature")
	}

	i1 := 1234

	/** converting the i1 variable into a string using Itoa method */
	str1 := strconv.Itoa(i1)
	fmt.Println(str1)

	i2 := 5678
	/** converting the i2 variable into a string using FormatInt method */
	str2 := strconv.FormatInt(int64(i2), 10)
	fmt.Println(str2)
}
