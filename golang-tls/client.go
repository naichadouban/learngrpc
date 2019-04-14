package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	roots := x509.NewCertPool()
	pem, err := ioutil.ReadFile("./conf/server.pem")
	if err != nil {
		log.Printf("read crt file error:%v\n", err)
	}
	ok := roots.AppendCertsFromPEM(pem)
	if !ok {
		panic("failed to parse root certificate")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: roots,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:8012/hello")
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	content, _ := ioutil.ReadAll(resp.Body)
	s := strings.TrimSpace(string(content))

	fmt.Println(s)
}
