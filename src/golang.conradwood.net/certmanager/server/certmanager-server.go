package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/go-acme/lego/v3/lego"
	"net"
	"strings"
	//	au "golang.conradwood.net/apis/auth"
	pb "golang.conradwood.net/apis/certmanager"
	"golang.conradwood.net/apis/common"
	"golang.conradwood.net/apis/h2gproxy"
	"golang.conradwood.net/certmanager/db"
	"golang.conradwood.net/go-easyops/auth"
	"golang.conradwood.net/go-easyops/errors"
	"golang.conradwood.net/go-easyops/server"
	"golang.conradwood.net/go-easyops/sql"
	"golang.conradwood.net/go-easyops/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"os"
	"time"
)

var (
	//	cert_retrieve_service = flag.String("cert_allowed_service", "37", "service id of the service serving the certs (usually h2gproxy)")
	debug      = flag.Bool("debug", false, "debug mode")
	startlego  = flag.Bool("startlego", true, "if false do not start lego and do not support creating public certificates")
	port       = flag.Int("port", 4100, "The grpc server port")
	legoClient *lego.Client
	certStore  *db.DBCertificate
	authStore  *db.DBStoreAuth
	allow_all  = flag.Bool("allow_all", false, "if true disable access control and always allow")
	reqChannel = make(chan *requestCertificate, 50)

	psql *sql.DB
)

type CertServer struct {
}

func main() {
	var err error
	flag.Parse()
	fmt.Printf("Starting CertManagerServer...\n")
	psql, err = sql.Open()
	utils.Bail("failed to open database", err)
	certStore = db.DefaultDBCertificate()
	authStore = db.DefaultDBStoreAuth()
	if *startlego {
		go requestWorker()
		go func() {
			// lego client needs some undetermined amount to initialise
			// it's useful to delay a bit so there is some chance to hit
			// at the first instance
			time.Sleep(time.Duration(5) * time.Second)
			refresher()
		}()
		go func() {
			for {
				legoClient, err = getLego()
				if err != nil {
					fmt.Printf("failed to get lego. will retry\n")
				} else {
					fmt.Printf("Lego client initialised\n")
					break
				}
				time.Sleep(time.Duration(5) * time.Second)
			}
		}()
	}

	sd := server.NewServerDef()
	sd.SetPort(*port)
	sd.Register = server.Register(
		func(server *grpc.Server) error {
			e := new(CertServer)
			pb.RegisterCertManagerServer(server, e)
			return nil
		},
	)
	fmt.Printf("Initialisation ready.\n")
	err = server.ServerStartup(sd)
	utils.Bail("Unable to start server", err)
	os.Exit(0)
}

/************************************
* grpc functions
************************************/
func (e *CertServer) ImportPublicCertificate(ctx context.Context, req *pb.ImportRequest) (*pb.Certificate, error) {
	err := errors.NeedsRoot(ctx)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(req.PemCertificate))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("Cannot parse import certificate: %s\n", err)
		return nil, err
	}
	hostname := ""
	fmt.Printf("Subject: %s\n", cert.Subject)
	for _, d := range cert.DNSNames {
		fmt.Printf("DNSName: %s\n", d)
		hostname = d
	}
	dbcl, err := certStore.ByHost(ctx, hostname)
	if err != nil {
		return nil, err
	}
	dbcl = filter_public_only(dbcl)
	if len(dbcl) != 0 {
		return nil, errors.Error(ctx, codes.AlreadyExists, "certificate already exists", "certificate for host %s already exists", hostname)
	}
	cu := ""
	u := auth.GetUser(ctx)
	if u != nil {
		cu = u.ID
	}
	cs := ""
	u = auth.GetUser(ctx)
	if u != nil {
		cs = u.ID
	}
	exp := uint32(cert.NotAfter.Unix())
	dbc := &pb.Certificate{
		Host:           hostname,
		PemCertificate: req.PemCertificate,
		PemPrivateKey:  req.PemPrivateKey,
		Created:        uint32(time.Now().Unix()),
		Expiry:         exp,
		CreatorUser:    cu,
		CreatorService: cs,
	}
	_, err = certStore.Save(ctx, dbc)
	if err != nil {
		return nil, err
	}
	return dbc, nil
}
func (e *CertServer) ListPublicCertificates(ctx context.Context, req *common.Void) (*pb.CertNameList, error) {
	dbc, err := certStore.All(ctx)
	if err != nil {
		return nil, err
	}
	dbc = filter_public_only(dbc)
	res := &pb.CertNameList{}
	for _, db := range dbc {

		ci := &pb.CertInfo{
			Hostname:    db.Host,
			Created:     db.Created,
			Expiry:      db.Expiry,
			LastRenewed: db.LastAttempt,
		}
		res.Certificates = append(res.Certificates, ci)
	}
	return res, nil
}

func checkAccess(ctx context.Context, cert *pb.Certificate) error {
	if *allow_all {
		return nil
	}
	service := auth.GetService(ctx)
	ret_service := auth.GetServiceIDByName("h2gproxy.H2GProxyService")
	if service != nil && service.ID == ret_service {
		return nil
	}

	u := auth.GetUser(ctx)
	if u != nil && u.ID == cert.CreatorUser {
		return nil
	}
	err := errors.NeedsRoot(ctx)
	if err != nil {
		return err
	}
	return nil
}
func rewrite_host_name(host string) string {
	hostname := host
	if strings.Contains(strings.ToLower(hostname), ".proxy.conradwood") {
		hostname = "proxy.conradwood.net"
		fmt.Printf("rewritten hostname %s to be exactly '%s'\n", host, hostname)
	}
	return hostname

}
func (e *CertServer) GetPublicCertificate(ctx context.Context, req *pb.PublicCertRequest) (*pb.ProcessedCertificate, error) {
	hostname := req.Hostname
	hostname = rewrite_host_name(hostname)
	dbc, err := certStore.ByHost(ctx, hostname)
	if err != nil {
		return nil, err
	}
	dbc = filter_public_only(dbc)
	if len(dbc) == 0 {
		return nil, errors.NotFound(ctx, "no certificate for \"%s\"\n", hostname)
	}
	cert := dbc[0]

	err = checkAccess(ctx, cert)
	if err != nil {
		// we return "notfound" so not to disclose that a cert exists but user has no access to it
		fmt.Printf("access error for \"%s\": %s\n", hostname, utils.ErrorString(err))
		return nil, errors.NotFound(ctx, "access denied (%s): %s", hostname, utils.ErrorString(err))
	}
	res := &pb.ProcessedCertificate{Cert: cert}
	tc, err := tls.X509KeyPair([]byte(cert.PemCertificate), []byte(cert.PemPrivateKey))
	if err != nil {
		fmt.Printf("Failed to parse cert %s: %s\n", hostname, err)
		return nil, err
	}
	// add the ca:
	block, _ := pem.Decode([]byte(cert.PemCA))
	if block == nil {
		fmt.Printf("certificate %s has no CA certificate\n", cert.Host)
	} else {
		_, xerr := x509.ParseCertificate(block.Bytes)
		if xerr != nil {
			fmt.Printf("Cannot parse certificate %s: %s\n", cert.Host, err)
			return nil, err
		}

		b := &bytes.Buffer{}
		err = pem.Encode(b, block)
		if err != nil {
			return nil, err
		}
		tc.Certificate = append(tc.Certificate, block.Bytes)
	}
	res.TLSCerts = append(res.TLSCerts, tc.Certificate...)
	return res, nil
}
func (e *CertServer) RequestPublicCertificate(ctx context.Context, req *pb.PublicCertRequest) (*common.Void, error) {
	hostname := req.Hostname
	if !isValid(hostname) {
		return nil, errors.InvalidArgs(ctx, "hostname invalid", "hostname \"%s\" too short, or invalid", hostname)
	}
	hostname = rewrite_host_name(hostname)
	fmt.Printf("Request by %s to get certificate for %s\n", auth.Description(auth.GetUser(ctx)), hostname)

	dbc, err := certStore.ByHost(ctx, hostname)
	if err != nil {
		return nil, err
	}
	dbc = filter_public_only(dbc)
	mostRecent := uint32(0)
	for _, db := range dbc {
		if mostRecent == 0 || db.Expiry < mostRecent {
			mostRecent = db.Expiry
		}
	}
	if mostRecent != 0 {
		mr := time.Unix(int64(mostRecent), 0)
		days := mr.Sub(time.Now()).Hours() / 24
		fmt.Printf("Certificate for %s is available until %v (%f days)\n", hostname, mr, days)
		if days > 30 {
			return nil, errors.Error(ctx, codes.AlreadyExists, "certificate already exists", "certificate for host %s already exists and is valid until %v (%f days)", hostname, mr, days)
		}
	}

	// annoyingly lego only does synchronous retrieval and it might take forever...
	// this should be in a channel, entirely async
	r := &requestCertificate{req: req}
	user := auth.GetUser(ctx)
	if user != nil {
		r.userid = user.ID
	}
	user = auth.GetService(ctx)
	if user != nil {
		r.serviceid = user.ID
	}

	reqChannel <- r
	res := &common.Void{}
	return res, nil
}

func requestWorker() {
	for {
		r := <-reqChannel
		err := request(r)
		if err != nil {
			fmt.Printf("Failed to get certificate \"%s\": %s\n", r.req.Hostname, err)
		}
	}
}

func isValid(hostname string) bool {
	// is it an ipaddress?
	if net.ParseIP(hostname) != nil {
		return false
	}
	if len(hostname) < 4 {
		return false
	}
	if strings.Index(hostname, ".") == -1 {
		return false
	}
	if strings.HasSuffix(hostname, "local") {
		return false
	}
	if strings.HasSuffix(hostname, "localdomain") {
		return false
	}
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		// no such host error
		return false
	}
	if len(addrs) == 0 {
		// resolved to 0 ip addresses
		return false
	}
	return true
}

func (e *CertServer) ServeHTML(ctx context.Context, req *h2gproxy.ServeRequest) (*h2gproxy.ServeResponse, error) {
	service := auth.GetService(ctx)
	if service == nil {
		return nil, errors.Unauthenticated(ctx, "access by h2gproxy only")
	}
	CUT := `/.well-known/acme-challenge/`
	chidx := strings.Index(req.Path, CUT)
	if chidx == -1 {
		fmt.Printf("ACME Challenge path invalid: \"%s\"\n", req.Path)
		return nil, errors.InvalidArgs(ctx, "acme challenge invalid", "acme challenge (%s) for host \"%s\" invalid", req.Path, req.Host)
	}
	challenge := req.Path[chidx+len(CUT):]

	fmt.Printf("ACME Challenge by %s: %s, Host: \"%s\", challenge=\"%s\"\n", service.ID, req.Path, req.Host, challenge)
	sa, err := getFromStoreByChallenge(ctx, challenge)
	if err != nil {
		fmt.Printf("ACME Challenge: unable to get requested challenge: %s\n", err)
	}
	if sa == nil {
		sa, err = getFromStore(ctx, req.Host)
		if err != nil {
			fmt.Printf("ACME Challenge \"%s\" could not be served: %s\n", req.Host, err)
			return nil, err
		}
	}
	if sa == nil {
		fmt.Printf("ACME Challenge not found: %s\n", req.Host)
		return nil, errors.NotFound(ctx, "acme challenge not found")
	}
	res := &h2gproxy.ServeResponse{
		HTTPResponseCode: 200,
		MimeType:         "text/plain",
		Body:             []byte(sa.KeyAuth),
	}
	fmt.Printf("ACME Challenge served for %s (%s)\n", req.Host, sa.Domain)
	return res, nil
}
