package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func calcPeriod(start string, day int) (notBefore time.Time, notAfter time.Time) {
	var err error
	if start == "" {
		notBefore = time.Now()
	} else {
		if notBefore, err = time.Parse("2006-01-02", start); err != nil {
			log.Fatalln("time parse failed.")
		} else {
			notAfter = notBefore.Add(time.Duration(day*24) * time.Hour)
		}
	}
	return
}

func genPrivateKey() interface{} {
	var err error
	var priv interface{}
	switch *ecdsaCurve {
	case "":
		if *useEd25519Key {
			_, priv, err = ed25519.GenerateKey(rand.Reader)
		} else {
			priv, err = rsa.GenerateKey(rand.Reader, *rsaBits)
		}
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		log.Fatalf("Unrecognized elliptic curve: %q", *ecdsaCurve)
	}

	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	return priv
}

func saveFile(path, filename string, data []byte, perm os.FileMode) error {

	if _, err := os.Stat(path); err == nil {
		if err = ioutil.WriteFile(fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), filename), data, perm); err != nil {
			return errors.Wrap(err, "Write file failed.")
		}
	} else if os.IsNotExist(err) {
		return errors.Wrap(err, "Directory not exist.")
	} else {
		return err
	}
	return nil
}

func createCAFile(name string, cert *x509.Certificate, key interface{}, caCert *x509.Certificate, caKey interface{}) {

	var pub interface{}
	switch k := key.(type) {
	case *rsa.PrivateKey:
		pub = &k.PublicKey
	case *ecdsa.PrivateKey:
		pub = &k.PublicKey
	case ed25519.PrivateKey:
		pub = k.Public().(ed25519.PublicKey)
	default:
		pub = nil
	}

	var privPm interface{}
	if caKey != nil {
		privPm = caKey
	} else {
		privPm = key
	}
	ca_b, err := x509.CreateCertificate(rand.Reader, cert, caCert, pub, privPm)
	if err != nil {
		log.Println("create failed", err)
		return
	}

	ca_b64 := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE",
		// Headers: map[string]string{},
		Bytes: ca_b})
	// fmt.Printf("%s private key： %s \n", name, ca_b)

	if err = saveFile(*dir, name+".pem", ca_b64, 0777); err != nil {
		log.Fatalln(err)
	}

	priv_b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		log.Fatalln("marshal PKCS8 private key failed")
	}
	// ioutil.WriteFile(name+".key", priv_b, 0777)

	priv_b64 := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: priv_b})

	if err = saveFile(*dir, name+".key", priv_b64, 0777); err != nil {
		log.Fatalln(err)
	}
}

// 公钥证书解析
func getCA() (*x509.Certificate, error) {
	//读取公钥证书并解码
	if pemTmp, err := ioutil.ReadFile(fmt.Sprintf("%s%sroot.pem", *dir, string(os.PathSeparator))); err != nil {
		return nil, errors.Wrap(err, "root.pem: Failed to read file.")
	} else {
		certBlock, restBlock := pem.Decode(pemTmp)
		if certBlock == nil {
			return nil, errors.Errorf("root.pem: Cert block is nil, rest block is %s.", restBlock)
		}

		return x509.ParseCertificate(certBlock.Bytes)
	}
}

// 私钥解析
func getPrivateKey() (key interface{}, err error) {
	if pemTmp, err := ioutil.ReadFile(fmt.Sprintf("%s%sroot.key", *dir, string(os.PathSeparator))); err != nil {
		return nil, errors.Wrap(err, "root.key: Failed to read file.")
	} else {
		certBlock, restBlock := pem.Decode(pemTmp)
		if certBlock == nil {
			return nil, errors.Errorf("root.key: Cert block is nil, rest block is %s.", restBlock)
		}

		return x509.ParsePKCS8PrivateKey(certBlock.Bytes)
	}
}

func genSerialNumber() *big.Int {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	if serialNumber, err := rand.Int(rand.Reader, serialNumberLimit); err != nil {
		log.Fatalln("generate serial number failed")
		return nil
	} else {
		return serialNumber
	}
}

func GenRootCA() {

	s, e := calcPeriod(*rootStartDate, *rootExpire)
	ca := &x509.Certificate{
		SerialNumber: genSerialNumber(),
		Subject: pkix.Name{
			// Country:            []string{"China"},
			Organization: []string{*rootOrg},
			// OrganizationalUnit: []string{"Shit company Unit"},
		},
		NotBefore:             s,
		NotAfter:              e,
		IsCA:                  true,
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	privCa := genPrivateKey()

	createCAFile("root", ca, privCa, ca, nil)

}

func GenServerCA() {

	s, e := calcPeriod(*serverStartDate, *serverExpire)
	server := &x509.Certificate{
		SerialNumber: genSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{*serverOrg},
		},
		NotBefore: s,
		NotAfter:  e,
		// SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	hosts := strings.Split(*host, ",")

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			server.IPAddresses = append(server.IPAddresses, ip)
		} else {
			server.DNSNames = append(server.DNSNames, h)
		}
	}
	privSer := genPrivateKey()

	var ca *x509.Certificate
	var err error
	var privCa interface{}

	ca, err = getCA()
	if err != nil {
		log.Fatalln(err)
	}

	privCa, err = getPrivateKey()
	if err != nil {
		log.Fatalln(err)
	}

	createCAFile("server", server, privSer, ca, privCa)
}

func GenClientCA() {

	s, e := calcPeriod(*clientStartDate, *clientExpire)
	client := &x509.Certificate{
		SerialNumber: genSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{*clientOrg},
		},
		NotBefore: s,
		NotAfter:  e,
		// SubjectKeyId: []byte{1, 2, 3, 4, 7},
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}
	privCli := genPrivateKey()

	var ca *x509.Certificate
	var err error
	var privCa interface{}

	ca, err = getCA()
	if err != nil {
		log.Fatalln(err)
	}

	privCa, err = getPrivateKey()
	if err != nil {
		log.Fatalln(err)
	}

	createCAFile("client", client, privCli, ca, privCa)

}
