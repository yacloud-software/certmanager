package main

import (
	"crypto/x509"
	"fmt"

	"golang.conradwood.net/go-easyops/utils"
)

type CertPrinter struct {
	certs []*x509.Certificate
}

func (cp *CertPrinter) AddCert(cert *x509.Certificate) {
	cp.certs = append(cp.certs, cert)
}
func (cp *CertPrinter) AddCerts(cert []*x509.Certificate) {
	cp.certs = append(cp.certs, cert...)
}
func (cp *CertPrinter) AddCertsAsBytes(cert_bytes [][]byte) {
	for _, b := range cert_bytes {
		certs := pem_to_certs(b)
		cp.AddCerts(certs)
	}
}
func (cp *CertPrinter) ToString() string {
	t := &utils.Table{}
	t.AddHeaders("subject", "commonname", "issuer")
	for _, c := range cp.certs {
		t.AddString(fmt.Sprintf("%s", c.Subject))
		t.AddString(c.Subject.CommonName)
		t.AddString(fmt.Sprintf("%s", c.Issuer))
		t.NewRow()
	}
	return t.ToPrettyString()
}
