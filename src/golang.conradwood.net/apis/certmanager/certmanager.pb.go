// Code generated by protoc-gen-go.
// source: protos/golang.conradwood.net/apis/certmanager/certmanager.proto
// DO NOT EDIT!

/*
Package certmanager is a generated protocol buffer package.

It is generated from these files:
	protos/golang.conradwood.net/apis/certmanager/certmanager.proto

It has these top-level messages:
	StoreAuth
	Certificate
	ProcessedCertificate
	PublicCertRequest
	CertInfo
	CertNameList
	GetVerificationRequest
	ImportRequest
*/
package certmanager

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "golang.conradwood.net/apis/common"
import h2gproxy "golang.conradwood.net/apis/h2gproxy"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type VerifyType int32

const (
	VerifyType_INVALID VerifyType = 0
	VerifyType_DNS     VerifyType = 1
	VerifyType_HTTP    VerifyType = 2
)

var VerifyType_name = map[int32]string{
	0: "INVALID",
	1: "DNS",
	2: "HTTP",
}
var VerifyType_value = map[string]int32{
	"INVALID": 0,
	"DNS":     1,
	"HTTP":    2,
}

func (x VerifyType) String() string {
	return proto.EnumName(VerifyType_name, int32(x))
}
func (VerifyType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// database type to temporarily store tokens by acme
type StoreAuth struct {
	ID      uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Domain  string `protobuf:"bytes,2,opt,name=Domain" json:"Domain,omitempty"`
	Token   string `protobuf:"bytes,3,opt,name=Token" json:"Token,omitempty"`
	KeyAuth string `protobuf:"bytes,4,opt,name=KeyAuth" json:"KeyAuth,omitempty"`
	Created uint32 `protobuf:"varint,5,opt,name=Created" json:"Created,omitempty"`
}

func (m *StoreAuth) Reset()                    { *m = StoreAuth{} }
func (m *StoreAuth) String() string            { return proto.CompactTextString(m) }
func (*StoreAuth) ProtoMessage()               {}
func (*StoreAuth) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StoreAuth) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *StoreAuth) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *StoreAuth) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *StoreAuth) GetKeyAuth() string {
	if m != nil {
		return m.KeyAuth
	}
	return ""
}

func (m *StoreAuth) GetCreated() uint32 {
	if m != nil {
		return m.Created
	}
	return 0
}

// database
type Certificate struct {
	ID             uint64 `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Host           string `protobuf:"bytes,2,opt,name=Host" json:"Host,omitempty"`
	PemCertificate string `protobuf:"bytes,3,opt,name=PemCertificate" json:"PemCertificate,omitempty"`
	PemPrivateKey  string `protobuf:"bytes,4,opt,name=PemPrivateKey" json:"PemPrivateKey,omitempty"`
	PemCA          string `protobuf:"bytes,5,opt,name=PemCA" json:"PemCA,omitempty"`
	Created        uint32 `protobuf:"varint,6,opt,name=Created" json:"Created,omitempty"`
	Expiry         uint32 `protobuf:"varint,7,opt,name=Expiry" json:"Expiry,omitempty"`
	CreatorUser    string `protobuf:"bytes,8,opt,name=CreatorUser" json:"CreatorUser,omitempty"`
	CreatorService string `protobuf:"bytes,9,opt,name=CreatorService" json:"CreatorService,omitempty"`
	LastAttempt    uint32 `protobuf:"varint,10,opt,name=LastAttempt" json:"LastAttempt,omitempty"`
	LastError      string `protobuf:"bytes,11,opt,name=LastError" json:"LastError,omitempty"`
}

func (m *Certificate) Reset()                    { *m = Certificate{} }
func (m *Certificate) String() string            { return proto.CompactTextString(m) }
func (*Certificate) ProtoMessage()               {}
func (*Certificate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Certificate) GetID() uint64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Certificate) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Certificate) GetPemCertificate() string {
	if m != nil {
		return m.PemCertificate
	}
	return ""
}

func (m *Certificate) GetPemPrivateKey() string {
	if m != nil {
		return m.PemPrivateKey
	}
	return ""
}

func (m *Certificate) GetPemCA() string {
	if m != nil {
		return m.PemCA
	}
	return ""
}

func (m *Certificate) GetCreated() uint32 {
	if m != nil {
		return m.Created
	}
	return 0
}

func (m *Certificate) GetExpiry() uint32 {
	if m != nil {
		return m.Expiry
	}
	return 0
}

func (m *Certificate) GetCreatorUser() string {
	if m != nil {
		return m.CreatorUser
	}
	return ""
}

func (m *Certificate) GetCreatorService() string {
	if m != nil {
		return m.CreatorService
	}
	return ""
}

func (m *Certificate) GetLastAttempt() uint32 {
	if m != nil {
		return m.LastAttempt
	}
	return 0
}

func (m *Certificate) GetLastError() string {
	if m != nil {
		return m.LastError
	}
	return ""
}

// derived certificate, pre-processed
type ProcessedCertificate struct {
	Cert *Certificate `protobuf:"bytes,1,opt,name=Cert" json:"Cert,omitempty"`
	// ready to add to tls.Certificate
	TLSCerts [][]byte `protobuf:"bytes,2,rep,name=TLSCerts,proto3" json:"TLSCerts,omitempty"`
}

func (m *ProcessedCertificate) Reset()                    { *m = ProcessedCertificate{} }
func (m *ProcessedCertificate) String() string            { return proto.CompactTextString(m) }
func (*ProcessedCertificate) ProtoMessage()               {}
func (*ProcessedCertificate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ProcessedCertificate) GetCert() *Certificate {
	if m != nil {
		return m.Cert
	}
	return nil
}

func (m *ProcessedCertificate) GetTLSCerts() [][]byte {
	if m != nil {
		return m.TLSCerts
	}
	return nil
}

type PublicCertRequest struct {
	Hostname   string     `protobuf:"bytes,1,opt,name=Hostname" json:"Hostname,omitempty"`
	VerifyType VerifyType `protobuf:"varint,2,opt,name=VerifyType,enum=certmanager.VerifyType" json:"VerifyType,omitempty"`
}

func (m *PublicCertRequest) Reset()                    { *m = PublicCertRequest{} }
func (m *PublicCertRequest) String() string            { return proto.CompactTextString(m) }
func (*PublicCertRequest) ProtoMessage()               {}
func (*PublicCertRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PublicCertRequest) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *PublicCertRequest) GetVerifyType() VerifyType {
	if m != nil {
		return m.VerifyType
	}
	return VerifyType_INVALID
}

type CertInfo struct {
	Hostname string `protobuf:"bytes,1,opt,name=Hostname" json:"Hostname,omitempty"`
	Created  uint32 `protobuf:"varint,2,opt,name=Created" json:"Created,omitempty"`
	Expiry   uint32 `protobuf:"varint,3,opt,name=Expiry" json:"Expiry,omitempty"`
}

func (m *CertInfo) Reset()                    { *m = CertInfo{} }
func (m *CertInfo) String() string            { return proto.CompactTextString(m) }
func (*CertInfo) ProtoMessage()               {}
func (*CertInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *CertInfo) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *CertInfo) GetCreated() uint32 {
	if m != nil {
		return m.Created
	}
	return 0
}

func (m *CertInfo) GetExpiry() uint32 {
	if m != nil {
		return m.Expiry
	}
	return 0
}

type CertNameList struct {
	Certificates []*CertInfo `protobuf:"bytes,1,rep,name=Certificates" json:"Certificates,omitempty"`
}

func (m *CertNameList) Reset()                    { *m = CertNameList{} }
func (m *CertNameList) String() string            { return proto.CompactTextString(m) }
func (*CertNameList) ProtoMessage()               {}
func (*CertNameList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *CertNameList) GetCertificates() []*CertInfo {
	if m != nil {
		return m.Certificates
	}
	return nil
}

type GetVerificationRequest struct {
	Hostname string `protobuf:"bytes,1,opt,name=Hostname" json:"Hostname,omitempty"`
}

func (m *GetVerificationRequest) Reset()                    { *m = GetVerificationRequest{} }
func (m *GetVerificationRequest) String() string            { return proto.CompactTextString(m) }
func (*GetVerificationRequest) ProtoMessage()               {}
func (*GetVerificationRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *GetVerificationRequest) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

type ImportRequest struct {
	Hostname       string `protobuf:"bytes,1,opt,name=Hostname" json:"Hostname,omitempty"`
	PemCertificate string `protobuf:"bytes,2,opt,name=PemCertificate" json:"PemCertificate,omitempty"`
	PemPrivateKey  string `protobuf:"bytes,3,opt,name=PemPrivateKey" json:"PemPrivateKey,omitempty"`
}

func (m *ImportRequest) Reset()                    { *m = ImportRequest{} }
func (m *ImportRequest) String() string            { return proto.CompactTextString(m) }
func (*ImportRequest) ProtoMessage()               {}
func (*ImportRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ImportRequest) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *ImportRequest) GetPemCertificate() string {
	if m != nil {
		return m.PemCertificate
	}
	return ""
}

func (m *ImportRequest) GetPemPrivateKey() string {
	if m != nil {
		return m.PemPrivateKey
	}
	return ""
}

func init() {
	proto.RegisterType((*StoreAuth)(nil), "certmanager.StoreAuth")
	proto.RegisterType((*Certificate)(nil), "certmanager.Certificate")
	proto.RegisterType((*ProcessedCertificate)(nil), "certmanager.ProcessedCertificate")
	proto.RegisterType((*PublicCertRequest)(nil), "certmanager.PublicCertRequest")
	proto.RegisterType((*CertInfo)(nil), "certmanager.CertInfo")
	proto.RegisterType((*CertNameList)(nil), "certmanager.CertNameList")
	proto.RegisterType((*GetVerificationRequest)(nil), "certmanager.GetVerificationRequest")
	proto.RegisterType((*ImportRequest)(nil), "certmanager.ImportRequest")
	proto.RegisterEnum("certmanager.VerifyType", VerifyType_name, VerifyType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for CertManager service

type CertManagerClient interface {
	// import a public certificate
	ImportPublicCertificate(ctx context.Context, in *ImportRequest, opts ...grpc.CallOption) (*Certificate, error)
	// list all certs
	ListPublicCertificates(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*CertNameList, error)
	// request a certificate. The lego/acme library is designed to synchronous execution.
	// (Its primary usecase seems to be command line tools)
	// this means we can only trigger it but do not get a response.
	// we have to poll GetPublicCertificate after requesting it until it is available
	RequestPublicCertificate(ctx context.Context, in *PublicCertRequest, opts ...grpc.CallOption) (*common.Void, error)
	// get a certificate from the certificate store. If none is available it
	// will return an error
	GetPublicCertificate(ctx context.Context, in *PublicCertRequest, opts ...grpc.CallOption) (*ProcessedCertificate, error)
	// serve acme URLs (permission from h2gproxy only)
	ServeHTML(ctx context.Context, in *h2gproxy.ServeRequest, opts ...grpc.CallOption) (*h2gproxy.ServeResponse, error)
}

type certManagerClient struct {
	cc *grpc.ClientConn
}

func NewCertManagerClient(cc *grpc.ClientConn) CertManagerClient {
	return &certManagerClient{cc}
}

func (c *certManagerClient) ImportPublicCertificate(ctx context.Context, in *ImportRequest, opts ...grpc.CallOption) (*Certificate, error) {
	out := new(Certificate)
	err := grpc.Invoke(ctx, "/certmanager.CertManager/ImportPublicCertificate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certManagerClient) ListPublicCertificates(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*CertNameList, error) {
	out := new(CertNameList)
	err := grpc.Invoke(ctx, "/certmanager.CertManager/ListPublicCertificates", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certManagerClient) RequestPublicCertificate(ctx context.Context, in *PublicCertRequest, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/certmanager.CertManager/RequestPublicCertificate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certManagerClient) GetPublicCertificate(ctx context.Context, in *PublicCertRequest, opts ...grpc.CallOption) (*ProcessedCertificate, error) {
	out := new(ProcessedCertificate)
	err := grpc.Invoke(ctx, "/certmanager.CertManager/GetPublicCertificate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *certManagerClient) ServeHTML(ctx context.Context, in *h2gproxy.ServeRequest, opts ...grpc.CallOption) (*h2gproxy.ServeResponse, error) {
	out := new(h2gproxy.ServeResponse)
	err := grpc.Invoke(ctx, "/certmanager.CertManager/ServeHTML", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CertManager service

type CertManagerServer interface {
	// import a public certificate
	ImportPublicCertificate(context.Context, *ImportRequest) (*Certificate, error)
	// list all certs
	ListPublicCertificates(context.Context, *common.Void) (*CertNameList, error)
	// request a certificate. The lego/acme library is designed to synchronous execution.
	// (Its primary usecase seems to be command line tools)
	// this means we can only trigger it but do not get a response.
	// we have to poll GetPublicCertificate after requesting it until it is available
	RequestPublicCertificate(context.Context, *PublicCertRequest) (*common.Void, error)
	// get a certificate from the certificate store. If none is available it
	// will return an error
	GetPublicCertificate(context.Context, *PublicCertRequest) (*ProcessedCertificate, error)
	// serve acme URLs (permission from h2gproxy only)
	ServeHTML(context.Context, *h2gproxy.ServeRequest) (*h2gproxy.ServeResponse, error)
}

func RegisterCertManagerServer(s *grpc.Server, srv CertManagerServer) {
	s.RegisterService(&_CertManager_serviceDesc, srv)
}

func _CertManager_ImportPublicCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertManagerServer).ImportPublicCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/certmanager.CertManager/ImportPublicCertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertManagerServer).ImportPublicCertificate(ctx, req.(*ImportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CertManager_ListPublicCertificates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertManagerServer).ListPublicCertificates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/certmanager.CertManager/ListPublicCertificates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertManagerServer).ListPublicCertificates(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _CertManager_RequestPublicCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicCertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertManagerServer).RequestPublicCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/certmanager.CertManager/RequestPublicCertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertManagerServer).RequestPublicCertificate(ctx, req.(*PublicCertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CertManager_GetPublicCertificate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublicCertRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertManagerServer).GetPublicCertificate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/certmanager.CertManager/GetPublicCertificate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertManagerServer).GetPublicCertificate(ctx, req.(*PublicCertRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CertManager_ServeHTML_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(h2gproxy.ServeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CertManagerServer).ServeHTML(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/certmanager.CertManager/ServeHTML",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CertManagerServer).ServeHTML(ctx, req.(*h2gproxy.ServeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CertManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "certmanager.CertManager",
	HandlerType: (*CertManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ImportPublicCertificate",
			Handler:    _CertManager_ImportPublicCertificate_Handler,
		},
		{
			MethodName: "ListPublicCertificates",
			Handler:    _CertManager_ListPublicCertificates_Handler,
		},
		{
			MethodName: "RequestPublicCertificate",
			Handler:    _CertManager_RequestPublicCertificate_Handler,
		},
		{
			MethodName: "GetPublicCertificate",
			Handler:    _CertManager_GetPublicCertificate_Handler,
		},
		{
			MethodName: "ServeHTML",
			Handler:    _CertManager_ServeHTML_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/golang.conradwood.net/apis/certmanager/certmanager.proto",
}

func init() {
	proto.RegisterFile("protos/golang.conradwood.net/apis/certmanager/certmanager.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 686 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x54, 0x5d, 0x4f, 0xdb, 0x4a,
	0x10, 0xbd, 0x71, 0x02, 0x49, 0x26, 0x80, 0xb8, 0x2b, 0x6e, 0xd8, 0x1b, 0x55, 0x55, 0x6a, 0x55,
	0x28, 0xaa, 0x90, 0x91, 0xd2, 0x4a, 0x55, 0xa5, 0x4a, 0x55, 0x4a, 0x10, 0xb1, 0x08, 0x34, 0x72,
	0x52, 0xd4, 0xbe, 0xd5, 0x24, 0x43, 0xb0, 0x8a, 0xbd, 0xee, 0xee, 0x42, 0xf1, 0x43, 0xdf, 0xfa,
	0x07, 0xfb, 0x8f, 0xaa, 0x5d, 0x3b, 0x61, 0x5d, 0x03, 0x42, 0x7d, 0xf2, 0x9e, 0xb3, 0xf3, 0x71,
	0x3c, 0x33, 0x3b, 0xf0, 0x2e, 0xe6, 0x4c, 0x32, 0xb1, 0x37, 0x67, 0x97, 0x7e, 0x34, 0x77, 0xa6,
	0x2c, 0xe2, 0xfe, 0xec, 0x3b, 0x63, 0x33, 0x27, 0x42, 0xb9, 0xe7, 0xc7, 0x81, 0xd8, 0x9b, 0x22,
	0x97, 0xa1, 0x1f, 0xf9, 0x73, 0xe4, 0xe6, 0xd9, 0xd1, 0x9e, 0xa4, 0x61, 0x50, 0x2d, 0xe7, 0xa1,
	0x30, 0x2c, 0x0c, 0x59, 0x94, 0x7d, 0x52, 0xe7, 0x56, 0xf7, 0x01, 0xfb, 0x8b, 0xee, 0x3c, 0xe6,
	0xec, 0x26, 0x59, 0x1e, 0x52, 0x1f, 0xfb, 0x07, 0xd4, 0xc7, 0x92, 0x71, 0xec, 0x5d, 0xc9, 0x0b,
	0xb2, 0x01, 0x96, 0xdb, 0xa7, 0xa5, 0x76, 0xa9, 0x53, 0xf1, 0x2c, 0xb7, 0x4f, 0x9a, 0xb0, 0xda,
	0x67, 0xa1, 0x1f, 0x44, 0xd4, 0x6a, 0x97, 0x3a, 0x75, 0x2f, 0x43, 0x64, 0x0b, 0x56, 0x26, 0xec,
	0x2b, 0x46, 0xb4, 0xac, 0xe9, 0x14, 0x10, 0x0a, 0xd5, 0x23, 0x4c, 0x54, 0x20, 0x5a, 0xd1, 0xfc,
	0x02, 0xaa, 0x9b, 0x7d, 0x8e, 0xbe, 0xc4, 0x19, 0x5d, 0x69, 0x97, 0x3a, 0xeb, 0xde, 0x02, 0xda,
	0xbf, 0x2c, 0x68, 0xec, 0x23, 0x97, 0xc1, 0x79, 0x30, 0xf5, 0x25, 0x16, 0x14, 0x10, 0xa8, 0x0c,
	0x98, 0x90, 0x59, 0x7e, 0x7d, 0x26, 0x3b, 0xb0, 0x31, 0xc2, 0xd0, 0xf0, 0xca, 0x64, 0xfc, 0xc1,
	0x92, 0xe7, 0xb0, 0x3e, 0xc2, 0x70, 0xc4, 0x83, 0x6b, 0x5f, 0xe2, 0x11, 0x26, 0x99, 0xaa, 0x3c,
	0xa9, 0xfe, 0x45, 0xf9, 0xf5, 0xb4, 0xb2, 0xba, 0x97, 0x02, 0x53, 0xf1, 0x6a, 0x4e, 0xb1, 0xaa,
	0xc9, 0xc1, 0x4d, 0x1c, 0xf0, 0x84, 0x56, 0xf5, 0x45, 0x86, 0x48, 0x1b, 0x1a, 0xda, 0x84, 0xf1,
	0x8f, 0x02, 0x39, 0xad, 0xe9, 0x68, 0x26, 0xa5, 0x74, 0x67, 0x70, 0x8c, 0xfc, 0x3a, 0x98, 0x22,
	0xad, 0xa7, 0xba, 0xf3, 0xac, 0x8a, 0x34, 0xf4, 0x85, 0xec, 0x49, 0x89, 0x61, 0x2c, 0x29, 0xe8,
	0x34, 0x26, 0x45, 0x9e, 0x40, 0x5d, 0xc1, 0x03, 0xce, 0x19, 0xa7, 0x0d, 0x1d, 0xe4, 0x96, 0xb0,
	0xbf, 0xc0, 0xd6, 0x88, 0xb3, 0x29, 0x0a, 0x81, 0x33, 0xb3, 0x1e, 0xbb, 0x50, 0x51, 0x50, 0x57,
	0xb7, 0xd1, 0xa5, 0x8e, 0x39, 0x7d, 0x86, 0x9d, 0xa7, 0xad, 0x48, 0x0b, 0x6a, 0x93, 0xe1, 0x58,
	0x1d, 0x05, 0xb5, 0xda, 0xe5, 0xce, 0x9a, 0xb7, 0xc4, 0xf6, 0x05, 0xfc, 0x3b, 0xba, 0x3a, 0xbb,
	0x0c, 0xa6, 0x0a, 0x7a, 0xf8, 0xed, 0x0a, 0x85, 0x76, 0x50, 0xed, 0x89, 0xfc, 0x10, 0x75, 0x8a,
	0xba, 0xb7, 0xc4, 0xe4, 0x35, 0xc0, 0x29, 0xf2, 0xe0, 0x3c, 0x99, 0x24, 0x31, 0xea, 0x66, 0x6e,
	0x74, 0xb7, 0x73, 0x02, 0x6e, 0xaf, 0x3d, 0xc3, 0xd4, 0xfe, 0x04, 0x35, 0x95, 0xc3, 0x8d, 0xce,
	0xd9, 0x83, 0x09, 0x8c, 0x7e, 0x59, 0xf7, 0xf5, 0xab, 0x6c, 0xf6, 0xcb, 0x76, 0x61, 0x4d, 0x45,
	0x3e, 0xf1, 0x43, 0x1c, 0x06, 0x42, 0x92, 0x37, 0x29, 0xce, 0x8a, 0x20, 0x68, 0xa9, 0x5d, 0xee,
	0x34, 0xba, 0xff, 0x15, 0xaa, 0xa4, 0xa4, 0x78, 0x39, 0x53, 0xfb, 0x15, 0x34, 0x0f, 0x51, 0x6a,
	0xd5, 0x8a, 0x09, 0x58, 0xf4, 0x88, 0x9a, 0xd8, 0x09, 0xac, 0xbb, 0x61, 0xcc, 0x1e, 0x57, 0xc0,
	0xe2, 0xcc, 0x5b, 0x8f, 0x9b, 0xf9, 0xf2, 0x1d, 0x33, 0xff, 0x62, 0xd7, 0x6c, 0x07, 0x69, 0x40,
	0xd5, 0x3d, 0x39, 0xed, 0x0d, 0xdd, 0xfe, 0xe6, 0x3f, 0xa4, 0x0a, 0xe5, 0xfe, 0xc9, 0x78, 0xb3,
	0x44, 0x6a, 0x50, 0x19, 0x4c, 0x26, 0xa3, 0x4d, 0xab, 0xfb, 0xb3, 0x9c, 0xbe, 0xd1, 0xe3, 0xb4,
	0x0a, 0xe4, 0x03, 0x6c, 0xa7, 0xc2, 0x6f, 0x67, 0x60, 0x91, 0xbe, 0x95, 0x2b, 0x57, 0xee, 0xf7,
	0x5a, 0xf7, 0x0e, 0x1c, 0xe9, 0x41, 0x53, 0xb5, 0xa0, 0x10, 0x4e, 0x90, 0x35, 0x27, 0x5b, 0x70,
	0xa7, 0x2c, 0x98, 0xb5, 0xfe, 0x2f, 0x44, 0x58, 0x76, 0x6f, 0x00, 0x34, 0xcb, 0x53, 0x14, 0xf5,
	0x34, 0xe7, 0x56, 0x18, 0xdc, 0x56, 0x2e, 0x09, 0xf9, 0x0c, 0x5b, 0x87, 0xf8, 0x17, 0x51, 0x9e,
	0xe5, 0xef, 0xef, 0x7a, 0x80, 0x6f, 0xa1, 0xae, 0xde, 0x38, 0x0e, 0x26, 0xc7, 0x43, 0xd2, 0x74,
	0x96, 0x9b, 0x58, 0x93, 0x8b, 0x38, 0xdb, 0x05, 0x5e, 0xc4, 0x2c, 0x12, 0xf8, 0xbe, 0x03, 0x3b,
	0x11, 0x4a, 0x73, 0xb9, 0x67, 0xeb, 0x5e, 0xed, 0x77, 0x33, 0xf9, 0xd9, 0xaa, 0x5e, 0xed, 0x2f,
	0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0x44, 0x78, 0xa0, 0xe0, 0x8e, 0x06, 0x00, 0x00,
}
