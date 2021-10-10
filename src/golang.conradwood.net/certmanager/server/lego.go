package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/go-acme/lego/v3/certificate"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/registration"
	//	au "golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/utils"
	"time"
)

var (
	last_attempt      = make(map[string]*requestAttempt)
	letsencrypt_email = flag.String("letsencrypt_email", "", "email for letsencryptaccount")
	letsencrypt_key   = flag.String("letsencrypt_key", "", "key for letsencryptaccount")
)

// You'll need a user or account type that implements acme.User
type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
func getLego() (*lego.Client, error) {
	var privateKey *ecdsa.PrivateKey
	var err error
	b, err := utils.ReadFile("configs/letsencrypt/private.key")
	if err != nil {
		fmt.Printf("Creating new key - Failed to read key (%s)\n", err)
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
	} else {
		block, _ := pem.Decode([]byte(b))
		x509Encoded := block.Bytes
		privateKey, err = x509.ParseECPrivateKey(x509Encoded)
		if err != nil {
			return nil, err
		}

	}
	//	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	//	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	//	fmt.Printf("Private Key: %s\n", string(pemEncoded))

	myUser := MyUser{
		Email: *letsencrypt_email,
		key:   privateKey,
	}

	config := lego.NewConfig(&myUser)
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}
	prov := &DBProvider{}
	err = client.Challenge.SetHTTP01Provider(prov)

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	myUser.Registration = reg
	fmt.Printf("Registration: %#v\n", reg)

	return client, nil
}

type DBProvider struct {
}

func (d *DBProvider) CleanUp(domain, token, keyAuth string) error {
	return nil
}
func (d *DBProvider) Present(domain, token, keyAuth string) error {
	store(domain, token, keyAuth)
	fmt.Printf("Presenting token %s for domain %s (key=%s)\n", token, domain, keyAuth)
	return nil
}

type requestCertificate struct {
	req       *pb.PublicCertRequest
	serviceid string
	userid    string
}

func request(r *requestCertificate) error {
	if !isValid(r.req.Hostname) {
		ctx := context.Background()
		return errors.InvalidArgs(ctx, "hostname invalid", "hostname \"%s\" too short, or invalid", r.req.Hostname)
	}

	if legoClient == nil {
		return fmt.Errorf("let's Encrypt unavailable. try later")
	}
	req := r.req
	hostname := req.Hostname
	if busy(hostname) {
		return fmt.Errorf("certificate for %s in progress already or blocked", hostname)
	}
	fmt.Printf("Requesting certificate for %s from acme\n", hostname)
	ctx := context.Background()
	legorequest := certificate.ObtainRequest{
		Domains: []string{hostname},
		Bundle:  true,
	}
	createAttempt(hostname)
	certificates, err := legoClient.Certificate.Obtain(legorequest)
	if err != nil {
		failAttempt(ctx, hostname)
		return err
	}

	if err != nil {
		return err
	}
	var cur *pb.Certificate
	dbc, err := certStore.ByHost(ctx, hostname)
	if err != nil {
		return err
	}
	if len(dbc) != 0 {
		cur = dbc[0]
		err = updateCert(r, cur, certificates)
		if err != nil {
			return err
		}
		cur.LastAttempt = uint32(time.Now().Unix())
		err = certStore.Update(ctx, cur)
		if err != nil {
			return err
		}
	} else {
		cur = &pb.Certificate{}
		err = updateCert(r, cur, certificates)
		if err != nil {
			return err
		}
		cur.LastAttempt = uint32(time.Now().Unix())
		_, err = certStore.Save(ctx, cur)
		if err != nil {
			return err
		}
	}

	clearStore(ctx, hostname)
	return nil
}
func failAttempt(ctx context.Context, hostname string) {
	var cur *pb.Certificate
	dbc, err := certStore.ByHost(ctx, hostname)
	if err != nil {
		fmt.Printf("Unable to load from db for cert for host %s: %s\n", hostname, err)
		return
	}
	if len(dbc) == 0 {
		return
	}
	cur = dbc[0]
	cur.LastAttempt = uint32(time.Now().Unix())
	err = certStore.Update(ctx, cur)
	if err != nil {
		fmt.Printf("Unable to update db for cert for host %s: %s\n", hostname, err)
		return
	}
}
func updateCert(r *requestCertificate, target *pb.Certificate, lego *certificate.Resource) error {
	block, _ := pem.Decode(lego.Certificate)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("Cannot parse lego certificate: %s\n", err)
		return err
	}
	exp := cert.NotAfter
	expts := uint32(exp.Unix())
	b := &bytes.Buffer{}
	err = pem.Encode(b, block)
	if err != nil {
		return err
	}
	clientCert := b.Bytes()
	target.Host = lego.Domain
	target.PemCertificate = string(clientCert)
	target.PemPrivateKey = string(lego.PrivateKey)
	target.PemCA = string(lego.IssuerCertificate)
	target.Created = uint32(time.Now().Unix())
	target.Expiry = expts
	target.CreatorUser = r.userid
	target.CreatorService = r.serviceid

	return nil
}

type requestAttempt struct {
	created time.Time
}

func createAttempt(hostname string) {
	ra := &requestAttempt{created: time.Now()}
	last_attempt[hostname] = ra
}
func busy(hostname string) bool {
	ra := last_attempt[hostname]
	if ra == nil {
		return false
	}
	if time.Since(ra.created) < time.Duration(30)*time.Minute {
		return true
	}
	return false

}
