package main

import (
	"crypto/x509"
	"fmt"

	"golang.conradwood.net/go-easyops/utils"
)

type CertPrinter struct {
	certs []*certmetaprinter
}
type certmetaprinter struct {
	cert   *x509.Certificate
	source string
}

func (cp *CertPrinter) AddCert(source string, cert *x509.Certificate) {
	cp.certs = append(cp.certs, &certmetaprinter{cert: cert, source: source})
}
func (cp *CertPrinter) AddCerts(source string, cert []*x509.Certificate) {
	for _, c := range cert {
		cp.AddCert(source, c)
	}
}
func (cp *CertPrinter) AddCertsAsBytes(source string, cert_bytes [][]byte) {
	for _, b := range cert_bytes {
		certs := pem_to_certs(b)
		cp.AddCerts(source, certs)
	}
}
func (cp *CertPrinter) ToString() string {
	t := &utils.Table{}
	t.AddHeaders("subject", "commonname", "issuer", "source")
	for _, x := range cp.certs {
		c := x.cert
		t.AddString(fmt.Sprintf("%s", c.Subject))
		t.AddString(c.Subject.CommonName)
		t.AddString(fmt.Sprintf("%s", c.Issuer))
		t.AddString(x.source)
		t.NewRow()
	}
	return t.ToPrettyString()
}
