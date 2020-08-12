package test

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	certi "github.com/goUtil/certi"
)

func TestEncode(t *testing.T) {
	flag.Parse()

	certi.GenRootCA()

	certi.GenServerCA()

	certi.GenClientCA()
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
