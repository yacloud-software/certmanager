package main

import (
	"flag"
	"fmt"
	pb "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/go-easyops/authremote"
	"golang.conradwood.net/go-easyops/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	certClient pb.CertManagerClient
	im         = flag.String("importdir", "", "`directory` containing certificate.pem and key.pem")
	hostname   = flag.String("hostname", "", "hostname to operate on")
	get        = flag.Bool("get", false, "get certificate")
	request    = flag.Bool("request", false, "request new certificate")
	reqlist    = flag.Bool("list", false, "request list of certs")
	save_pem   = flag.String("save_pem", "/tmp/certs/", "pem `directory` to save certificate to")
	local      = flag.Bool("local", false, "if true, generate locally signed certificate")
)

func main() {
	flag.Parse()

	certClient = pb.GetCertManagerClient()
	if *im != "" {
		doimport()
		os.Exit(0)
	}
	if *reqlist {
		list()
	}
	if *request {
		utils.Bail("failed to request cert", requestCert(*hostname))
	}
	if *get {
		doget(*hostname)
	}
	fmt.Printf("Done.\n")
	os.Exit(0)
}
func doget(host string) {
	if len(*hostname) < 3 {
		fmt.Printf("Invalid hostname \"%s\" (missing -hostname option?)\n", *hostname)
		os.Exit(10)
	}
	if len(*save_pem) < 3 {
		fmt.Printf("Invalid save directory \"%s\" (missing -save_pem option?)\n", *im)
		os.Exit(10)
	}
	ctx := authremote.Context()
	pcr := &pb.PublicCertRequest{Hostname: host, VerifyType: pb.VerifyType_HTTP}
	cert, err := certClient.GetPublicCertificate(ctx, pcr)
	utils.Bail("Failed to get cert", err)
	fmt.Printf("Certificate for %s\n", host)
	fmt.Printf("   Expires    : %s\n", utils.TimestampString(cert.Cert.Expiry))
	save(cert.Cert)
}
func save(cert *pb.Certificate) {
	if *save_pem == "" {
		fmt.Printf("Not saving (missing -save_pem option)\n")
		return
	}
	if cert.Host == "" {
		fmt.Printf("Warning - certificate has no hostname\n")
	}

	d := *save_pem + "/" + cert.Host
	err := os.MkdirAll(d, 0777)
	utils.Bail("failed to mkdir", err)
	bcert := fixup(cert.PemCertificate)
	bkey := fixup(cert.PemPrivateKey)
	bca := fixup(cert.PemCA)
	saveFile(d+"/ca.pem", []byte(bca))
	saveFile(d+"/cert-and-key.pem", []byte(bcert+bca+bkey))
	saveFile(d+"/certificate.pem", []byte(bcert+bca))
	saveFile(d+"/key.pem", []byte(bkey))
	if cert.PemPublicKey != "" {
		saveFile(d+"/public_key.pem", []byte(cert.PemPublicKey))
		dns_txt := to_dns_txt("v=DKIM1; k=rsa; p=", cert.PemPublicKey)
		saveFile(d+"/public_key.dns_txt", []byte(dns_txt))
	}
}

// returns a bind zonefile line
func to_dns_txt(prefix, pem string) string {
	pem = strings.ReplaceAll(pem, "\n", "")
	pem = strings.ReplaceAll(pem, "-----BEGIN PUBLIC KEY-----", "")
	pem = strings.ReplaceAll(pem, "-----END PUBLIC KEY-----", "")
	// pem = strings.TrimSuffix(pem, "==") // need that!
	pem = prefix + pem
	pem = strings.ReplaceAll(pem, ";", "\\;")
	//pem = pem + "\n"
	var entries []string
	offset := 0
	repeat := true
	for repeat {
		size := 255
		if len(pem) < size+offset {
			size = len(pem) - offset
			repeat = false
		}
		e := pem[offset : offset+size]
		entries = append(entries, e)
		offset = offset + size
	}
	res := ""
	deli := ""
	for _, e := range entries {
		res = res + deli + fmt.Sprintf("\"%s\"", e)
		deli = " "
		//fmt.Printf("%03d <<"+e+">>\n", len(e))
	}
	return res
}

func fixup(pemthing string) string {
	res := strings.TrimSuffix(pemthing, "\n")
	res = strings.TrimPrefix(res, "\n")
	return res + "\n"
}
func doimport() {
	cb, err := utils.ReadFile(*im + "/certificate.pem")
	utils.Bail("failed to read certificate", err)
	kb, err := utils.ReadFile(*im + "/key.pem")
	utils.Bail("failed to read key", err)
	ctx := authremote.Context()
	ir := &pb.ImportRequest{
		Hostname:       *hostname,
		PemCertificate: string(cb),
		PemPrivateKey:  string(kb),
	}
	cert, err := certClient.ImportPublicCertificate(ctx, ir)
	utils.Bail("failed to import", err)
	fmt.Printf("Imported certificate #%d: %s\n", cert.ID, cert.Host)
}

func list() {
	ctx := authremote.Context()
	response, err := certClient.ListPublicCertificates(ctx, &common.Void{})
	utils.Bail("Failed to ping server", err)
	t := utils.Table{}
	t.AddHeaders("Hostname", "Created", "Expiry", "In days")
	for _, c := range response.Certificates {
		mr := time.Unix(int64(c.Expiry), 0)
		days := mr.Sub(time.Now()).Hours() / 24
		fmt.Printf("Hostname: %s, Expiry: %v (%f days)\n", c.Hostname, mr, days)
		t.AddString(c.Hostname)
		t.AddTimestamp(c.Created)
		e := ""
		if days < 0 {
			e = " EXPIRED "
		}
		t.AddString(fmt.Sprintf("%v%s", mr, e))
		if days < 0 {
			t.AddInt(0)
		} else {
			t.AddInt(int(days))
		}
		t.NewRow()
	}
	fmt.Printf(t.ToPrettyString())
}

func requestCert(host string) error {
	if *local {
		return requestLocalCert(host)
	}
	ctx := authremote.Context()
	pcr := &pb.PublicCertRequest{Hostname: host, VerifyType: pb.VerifyType_DNS}
	response, err := certClient.RequestPublicCertificate(ctx, pcr)
	if err != nil {
		return err
	}
	fmt.Printf("Response to ping: %v\n", response)
	return nil
}

func requestLocalCert(host string) error {
	ctx := authremote.Context()
	req := &pb.LocalCertificateRequest{Subject: host}
	response, err := certClient.GetLocalCertificate(ctx, req)
	if err != nil {
		return err
	}
	fmt.Printf("Server responded with Certificate #%d\n", response.ID)
	save(response)
	return nil
}
func saveFile(filename string, content []byte) {
	filename, err := filepath.Abs(filename)
	utils.Bail("failed to make filename absolute", err)
	err = utils.WriteFile(filename, content)
	if err != nil {
		fmt.Printf("I/O error: %s\n", err)
		return
	}
	fmt.Printf("Saved \"%s\"\n", filename)
}
