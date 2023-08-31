package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	pb "golang.conradwood.net/apis/certmanager"
)

// a helper struct to work with certificate protos (which are essentially pem encoded blobs)
type Cert struct {
	x509Cert   *x509.Certificate
	x509CA     *x509.Certificate
	privateKey crypto.Signer
}

func parseCert(cert *pb.Certificate) (*Cert, error) {
	res := &Cert{}
	block, _ := pem.Decode([]byte(cert.PemCertificate))
	c, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("failure parsing import certificate: %s\n", err)
		return nil, err
	}
	res.x509Cert = c

	block, _ = pem.Decode([]byte(cert.PemCA))
	c, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("failure parsing import certificate: %s\n", err)
		return nil, err
	}
	res.x509CA = c

	block, _ = pem.Decode([]byte(cert.PemPrivateKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Printf("failure parsing private key for cert %s: %s\n", cert.Host, err)
		return nil, err
	}
	res.privateKey = key
	return res, nil
}
func (c *Cert) PrivateKey() crypto.Signer {
	return c.privateKey
}
func (c *Cert) X509Certificate() *x509.Certificate {
	return c.x509Cert
}
func (c *Cert) X509CA() *x509.Certificate {
	return c.x509CA
}
