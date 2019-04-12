package util

import (
	"crypto/tls"
	"io/ioutil"

	"golang.org/x/net/http2"
)

func GetTlsConfig(certPemPath, certKeyPath string) (*tls.Config, error) {
	cert, err := ioutil.ReadFile(certPemPath)
	if err != nil {
		return nil, err
	}
	key, err := ioutil.ReadFile(certKeyPath)
	if err != nil {
		return nil, err
	}
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{pair},
		NextProtos:   []string{http2.NextProtoTLS},
	}, nil

}
