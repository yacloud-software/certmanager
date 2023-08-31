package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	//	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"math/big"
	"time"
)

var (
	default_ca_subject = flag.String("default_ca", "ca.conradwood.net", "ca subject")
)

func (e *CertServer) LocalCertificate(ctx context.Context, req *pb.LocalCertificateRequest) (*pb.Certificate, error) {
	u := auth.GetUser(ctx)
	if u == nil {
		return nil, errors.Unauthenticated(ctx, "please log in")
	}
	certs, err := certStore.ByHost(ctx, req.Subject)
	if err != nil {
		return nil, err
	}
	exp := time.Now().Unix()
	for _, c := range certs {
		if !c.IsLocalCert {
			continue
		}
		if c.IsLocalCA {
			continue
		}
		if int64(c.Expiry) <= exp {
			continue
		}
		return c, nil
	}

	ca, err := get_ca(*default_ca_subject)
	if err != nil {
		return nil, err
	}
	cert, err := gen_new_cert(ca, req.Subject, u.Email)
	if err != nil {
		return nil, err
	}
	certStore.Save(ctx, cert)
	return cert, nil
}

// get ca certificate or create new one
func get_ca(subject string) (*pb.Certificate, error) {
	certs, err := certStore.ByHost(context.Background(), subject)
	if err != nil {
		return nil, err
	}
	var ca *pb.Certificate
	for _, c := range certs {
		if !c.IsLocalCA {
			continue
		}
		// TODO: check time.After (expiry)
		ca = c
		break
	}
	if ca != nil {
		return ca, nil
	}
	ca, err = create_ca(subject)
	if err != nil {
		return nil, err
	}
	if ca.ID == 0 {
		certStore.Save(context.Background(), ca)
	}
	return ca, nil
}

func gen_new_cert(ca *pb.Certificate, subject string, requestor_email string) (*pb.Certificate, error) {
	fmt.Printf("Creating cert for \"%s\"\n", subject)
	cacert, err := parseCert(ca)
	if err != nil {
		return nil, fmt.Errorf("failure parsing ca cert: %w", err)
	}
	cert := &x509.Certificate{
		EmailAddresses: []string{requestor_email},
		DNSNames:       []string{subject},
		SerialNumber:   big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{subject},
			Country:       []string{"UK"},
			Province:      []string{"Greater London"},
			Locality:      []string{"Hammersmith"},
			StreetAddress: []string{"Lyric Square"},
			PostalCode:    []string{"1"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 3, 0), // 3 months
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cacert.X509Certificate(), &certPrivKey.PublicKey, cacert.PrivateKey())
	if err != nil {
		return nil, err
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
	res := &pb.Certificate{
		Host:           subject,
		PemCertificate: certPEM.String(),
		PemPrivateKey:  certPrivKeyPEM.String(),
		PemCA:          ca.PemCA,
		Created:        uint32(cert.NotBefore.Unix()),
		Expiry:         uint32(cert.NotAfter.Unix()),
		IsLocalCert:    true,
	}
	return res, nil
}

// create a new ca
func create_ca(subject string) (*pb.Certificate, error) {
	fmt.Printf("Creating ca for \"%s\"\n", subject)
	ca := &x509.Certificate{
		DNSNames:              []string{subject},
		IssuingCertificateURL: []string{"https://" + subject + "/ca.html"},
		SerialNumber:          big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{subject},
			Country:       []string{"UK"},
			Province:      []string{"Greater London"},
			Locality:      []string{"Hammersmith"},
			StreetAddress: []string{"Lyric Square"},
			PostalCode:    []string{"1"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}
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
	res := &pb.Certificate{
		Host:           subject,
		PemCertificate: caPEM.String(),
		PemPrivateKey:  caPrivKeyPEM.String(),
		PemCA:          caPEM.String(),
		Created:        uint32(ca.NotBefore.Unix()),
		Expiry:         uint32(ca.NotAfter.Unix()),
		IsLocalCA:      true,
		IsLocalCert:    true,
	}
	return res, nil
}
