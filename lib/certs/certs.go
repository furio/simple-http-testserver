package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/emirpasic/gods/sets/hashset"
)

func GenerateCerts(dsnCert string) (serverTLS []tls.Certificate, rootTLS []byte, err error) {

	set := hashset.New("localhost")
	for dns := range strings.Split(dsnCert, ",") {
		set.Add(dns)
	}

	dnsNames := make([]string, set.Size())
	for i, v := range set.Values() {
		dnsNames[i] = fmt.Sprint(v)
	}

	return certsetup(dnsNames)
}

func certsetup(dnsList []string) (serverTLS []tls.Certificate, rootTLS []byte, err error) {
	// set up our CA certificate
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			Organization:  []string{"Furiosoft"},
			Country:       []string{"IT"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{"Italy"},
			PostalCode:    []string{"90210"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	// set up our server certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Furiosoft"},
			Country:       []string{"IT"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{"Italy"},
			PostalCode:    []string{"901210"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     dnsList,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	serverCert, err := tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
	if err != nil {
		return nil, nil, err
	}

	serverTLS = []tls.Certificate{serverCert}
	rootTLS = caPEM.Bytes()

	return
}

func LoadCerts(caFile, keyFile, certFile string) (serverTLS []tls.Certificate, rootTLS []byte, err error) {

	caPEM, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, nil, err
	}

	certPrivKey, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	certPEM, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}

	serverCert, err := tls.X509KeyPair(certPEM, certPrivKey)
	if err != nil {
		return nil, nil, err
	}

	serverTLS = []tls.Certificate{serverCert}
	rootTLS = caPEM

	return serverTLS, rootTLS, nil
}
