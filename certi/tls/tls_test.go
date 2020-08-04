package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"
	"time"
)

type myhandler struct {
}

func (h *myhandler) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprintf(w,
		"Hi, This is an example of http service in golang!\n")
}

func TestMain(m *testing.M) {
	var exitCode int
	defer os.Exit(exitCode)

	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile("root.pem")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	s := &http.Server{
		Addr:    ":8088",
		Handler: &myhandler{},
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	time.AfterFunc(3*time.Second, func() {
		log.Println("run main after")
		exitCode = m.Run()
	})

	err = s.ListenAndServeTLS("server.pem", "server.key")
	if err != nil {
		fmt.Println("ListenAndServeTLS err:", err)
	}
}

func TestTcpClient(t *testing.T) {
	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile("root.pem")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
		return
	}
	config := &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cliCrt},
	}
	conn, err := tls.Dial("tcp", "127.0.0.1:8088", config)
	if err != nil {
		fmt.Printf("client: dial: %s\n", err)
	}
	defer conn.Close()
	fmt.Println("client: connected to: ", conn.RemoteAddr())

	req, _ := http.NewRequest("Get", "/", nil)
	reqBuff, _ := httputil.DumpRequest(req, false)
	conn.Write(reqBuff)

	reply := make([]byte, 256)
	n, err := conn.Read(reply)
	fmt.Printf("client read:(%d bytes)\n%q \n", n, string(reply[:n]))

}

func TestHttpClient(t *testing.T) {
	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile("root.pem")
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:8088")
	if err != nil {
		fmt.Println("Get error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
