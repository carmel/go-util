package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestEncode(t *testing.T) {
	flag.Parse()

	GenRootCA()

	GenServerCA()

	GenClientCA()
}

func TestDecode(t *testing.T) {
	//读取私钥证书并解码
	pemTmp, err := ioutil.ReadFile("tls/root.key")
	if err != nil {
		fmt.Println(err)
		return
	}

	certBlock, restBlock := pem.Decode(pemTmp)
	if certBlock == nil {
		fmt.Println(restBlock)
		return
	}

	//私钥解析
	certBody, err := x509.ParsePKCS8PrivateKey(certBlock.Bytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v,%s \n", certBody, base64.StdEncoding.EncodeToString(certBlock.Bytes))

	//读取公钥证书并解码
	pemTmp, err = ioutil.ReadFile("tls/root.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	certBlock, restBlock = pem.Decode(pemTmp)
	if certBlock == nil {
		fmt.Println(restBlock)
		return
	}

	//证书解析
	certBody, err = x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	//可以根据证书结构解析
	// fmt.Println(certBody.SignatureAlgorithm)
	// fmt.Println(certBody.PublicKeyAlgorithm)
	fmt.Printf("%+v \n", certBody)
}

func TestEncryptAndDecrypt(t *testing.T) {
	clearText := []byte("I'm clear text")
	bPem, _ := ioutil.ReadFile("tls/server.pem")
	block, _ := pem.Decode(bPem)
	pub, _ := x509.ParsePKCS1PublicKey(block.Bytes)

	fmt.Println(string(bPem))

	cipherText, err := RsaEncrypt(clearText, pub)

	fmt.Println(cipherText, err)

	// 私钥
	bPem, _ = ioutil.ReadFile("tls/server.key")
	block, _ = pem.Decode(bPem)
	prv, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	fmt.Println(RsaDecrypt(cipherText, prv))
}

// RsaEncrypt encrypts data with public key
func RsaEncrypt(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("EncryptOAEP: %v", err)
	}
	return ciphertext, nil
}

// RsaDecrypt decrypts data with private key
func RsaDecrypt(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("DecryptOAEP: %v", err)
	}
	return plaintext, nil
}
