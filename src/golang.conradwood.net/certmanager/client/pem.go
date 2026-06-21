package main

import (
	"crypto/x509"
	"encoding/pem"

	"golang.conradwood.net/go-easyops/utils"
)

func pem_to_certs(pems []byte) []*x509.Certificate {
	var res []*x509.Certificate
	b := pems
	for len(b) != 0 {
		//fmt.Printf("BLOCK:\n%s\n", string(b))
		block, nb := pem.Decode(b)
		if block == nil {
			break
		}
		b = nb
		certs, err := x509.ParseCertificates(block.Bytes)
		utils.Bail("failed to parse certs", err)
		for _, cert := range certs {
			//			fmt.Printf("Cert: %#v\n", cert)

			res = append(res, cert)
		}
		//}

	}
	return res
}
