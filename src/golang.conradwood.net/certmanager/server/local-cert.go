package main

import (
	"context"
	pb "golang.conradwood.net/apis/certmanager"
)

func (e *CertServer) LocalCertificate(ctx context.Context, req *pb.LocalCertificateRequest) (*pb.Certificate, error) {
	res := &pb.Certificate{}
	return res, nil
}
