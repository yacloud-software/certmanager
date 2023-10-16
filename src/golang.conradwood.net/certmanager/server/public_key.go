package main

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	pb "golang.conradwood.net/apis/certmanager"
)

func add_public_key(cert *pb.Certificate) error {
	priv := cert.PemPrivateKey

	block, _ := pem.Decode([]byte(priv))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return fmt.Errorf("[publickey] failed to decode PEM block containing public key")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("[publickey] failed to parse key: %s", err)
	}

	publicKeyDer, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return fmt.Errorf("[publickey] Failed to marshal public key: %s\n", err)
	}
	pubKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	pubKeyPem := string(pem.EncodeToMemory(&pubKeyBlock))
	cert.PemPublicKey = pubKeyPem
	return nil
}
