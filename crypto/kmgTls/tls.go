package kmgTls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	//"flag"
	"fmt"
	//"log"
	"math/big"
	//"net"
	//"os"
	//"strings"
	"bytes"
	"crypto/tls"
	"time"
)

func CreateTlsConfig() (*tls.Config, error) {
	cert, err := CreateSelfCert()
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{*cert}}, nil
}

//unity 3d can use ncat cert,but can not use this!!! openssl verify fail!!!
func CreateSelfCert() (*tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %s", err)
	}

	notBefore := time.Date(2001, 12, 31, 23, 59, 59, 0, time.UTC)
	notAfter := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)

	template := x509.Certificate{

		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			Organization: []string{"localhost"},
			CommonName:   "localhost",
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	/*
		hosts := strings.Split(*host, ",")
		for _, h := range hosts {
			if ip := net.ParseIP(h); ip != nil {
				template.IPAddresses = append(template.IPAddresses, ip)
			} else {
				template.DNSNames = append(template.DNSNames, h)
			}
		}

		if *isCA {
			template.IsCA = true
			template.KeyUsage |= x509.KeyUsageCertSign
		}
	*/
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, fmt.Errorf("Failed to create certificate: %s", err)
	}

	certOut := &bytes.Buffer{}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut := &bytes.Buffer{}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	cert, err := tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
