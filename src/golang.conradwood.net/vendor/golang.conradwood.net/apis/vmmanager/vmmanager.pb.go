// Code generated by protoc-gen-go.
// source: golang.conradwood.net/apis/vmmanager/vmmanager.proto
// DO NOT EDIT!

/*
Package vmmanager is a generated protocol buffer package.

It is generated from these files:
	golang.conradwood.net/apis/vmmanager/vmmanager.proto

It has these top-level messages:
	Nic
	Disk
	Server
	ServerList
	Image
	CloneRequest
	CloneResponse
	ImageList
	ServerRequest
*/
package vmmanager

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "golang.conradwood.net/apis/common"

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

type Nic struct {
	Lan  uint32 `protobuf:"varint,1,opt,name=Lan" json:"Lan,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
	MAC  string `protobuf:"bytes,3,opt,name=MAC" json:"MAC,omitempty"`
}

func (m *Nic) Reset()                    { *m = Nic{} }
func (m *Nic) String() string            { return proto.CompactTextString(m) }
func (*Nic) ProtoMessage()               {}
func (*Nic) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Nic) GetLan() uint32 {
	if m != nil {
		return m.Lan
	}
	return 0
}

func (m *Nic) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Nic) GetMAC() string {
	if m != nil {
		return m.MAC
	}
	return ""
}

type Disk struct {
	ID      string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
	HDDSize uint64 `protobuf:"varint,3,opt,name=HDDSize" json:"HDDSize,omitempty"`
	Created uint64 `protobuf:"varint,4,opt,name=Created" json:"Created,omitempty"`
}

func (m *Disk) Reset()                    { *m = Disk{} }
func (m *Disk) String() string            { return proto.CompactTextString(m) }
func (*Disk) ProtoMessage()               {}
func (*Disk) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Disk) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Disk) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Disk) GetHDDSize() uint64 {
	if m != nil {
		return m.HDDSize
	}
	return 0
}

func (m *Disk) GetCreated() uint64 {
	if m != nil {
		return m.Created
	}
	return 0
}

type Server struct {
	ID      string  `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Name    string  `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
	Nics    []*Nic  `protobuf:"bytes,3,rep,name=Nics" json:"Nics,omitempty"`
	Disks   []*Disk `protobuf:"bytes,4,rep,name=Disks" json:"Disks,omitempty"`
	Cores   uint32  `protobuf:"varint,5,opt,name=Cores" json:"Cores,omitempty"`
	Ram     uint64  `protobuf:"varint,6,opt,name=Ram" json:"Ram,omitempty"`
	Created uint64  `protobuf:"varint,7,opt,name=Created" json:"Created,omitempty"`
}

func (m *Server) Reset()                    { *m = Server{} }
func (m *Server) String() string            { return proto.CompactTextString(m) }
func (*Server) ProtoMessage()               {}
func (*Server) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Server) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Server) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Server) GetNics() []*Nic {
	if m != nil {
		return m.Nics
	}
	return nil
}

func (m *Server) GetDisks() []*Disk {
	if m != nil {
		return m.Disks
	}
	return nil
}

func (m *Server) GetCores() uint32 {
	if m != nil {
		return m.Cores
	}
	return 0
}

func (m *Server) GetRam() uint64 {
	if m != nil {
		return m.Ram
	}
	return 0
}

func (m *Server) GetCreated() uint64 {
	if m != nil {
		return m.Created
	}
	return 0
}

type ServerList struct {
	Servers []*Server `protobuf:"bytes,1,rep,name=Servers" json:"Servers,omitempty"`
}

func (m *ServerList) Reset()                    { *m = ServerList{} }
func (m *ServerList) String() string            { return proto.CompactTextString(m) }
func (*ServerList) ProtoMessage()               {}
func (*ServerList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ServerList) GetServers() []*Server {
	if m != nil {
		return m.Servers
	}
	return nil
}

type Image struct {
	ID      string `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
	Size    uint64 `protobuf:"varint,3,opt,name=Size" json:"Size,omitempty"`
	Created uint64 `protobuf:"varint,4,opt,name=Created" json:"Created,omitempty"`
}

func (m *Image) Reset()                    { *m = Image{} }
func (m *Image) String() string            { return proto.CompactTextString(m) }
func (*Image) ProtoMessage()               {}
func (*Image) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Image) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Image) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Image) GetSize() uint64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *Image) GetCreated() uint64 {
	if m != nil {
		return m.Created
	}
	return 0
}

type CloneRequest struct {
	ServerID string `protobuf:"bytes,1,opt,name=ServerID" json:"ServerID,omitempty"`
}

func (m *CloneRequest) Reset()                    { *m = CloneRequest{} }
func (m *CloneRequest) String() string            { return proto.CompactTextString(m) }
func (*CloneRequest) ProtoMessage()               {}
func (*CloneRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *CloneRequest) GetServerID() string {
	if m != nil {
		return m.ServerID
	}
	return ""
}

type CloneResponse struct {
	NewID   string `protobuf:"bytes,1,opt,name=NewID" json:"NewID,omitempty"`
	NewName string `protobuf:"bytes,2,opt,name=NewName" json:"NewName,omitempty"`
}

func (m *CloneResponse) Reset()                    { *m = CloneResponse{} }
func (m *CloneResponse) String() string            { return proto.CompactTextString(m) }
func (*CloneResponse) ProtoMessage()               {}
func (*CloneResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CloneResponse) GetNewID() string {
	if m != nil {
		return m.NewID
	}
	return ""
}

func (m *CloneResponse) GetNewName() string {
	if m != nil {
		return m.NewName
	}
	return ""
}

type ImageList struct {
	Images []*Image `protobuf:"bytes,1,rep,name=Images" json:"Images,omitempty"`
}

func (m *ImageList) Reset()                    { *m = ImageList{} }
func (m *ImageList) String() string            { return proto.CompactTextString(m) }
func (*ImageList) ProtoMessage()               {}
func (*ImageList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ImageList) GetImages() []*Image {
	if m != nil {
		return m.Images
	}
	return nil
}

type ServerRequest struct {
	NewName string `protobuf:"bytes,1,opt,name=NewName" json:"NewName,omitempty"`
	DiskGB  uint32 `protobuf:"varint,2,opt,name=DiskGB" json:"DiskGB,omitempty"`
	RamGB   uint32 `protobuf:"varint,3,opt,name=RamGB" json:"RamGB,omitempty"`
	Cores   uint32 `protobuf:"varint,4,opt,name=Cores" json:"Cores,omitempty"`
}

func (m *ServerRequest) Reset()                    { *m = ServerRequest{} }
func (m *ServerRequest) String() string            { return proto.CompactTextString(m) }
func (*ServerRequest) ProtoMessage()               {}
func (*ServerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ServerRequest) GetNewName() string {
	if m != nil {
		return m.NewName
	}
	return ""
}

func (m *ServerRequest) GetDiskGB() uint32 {
	if m != nil {
		return m.DiskGB
	}
	return 0
}

func (m *ServerRequest) GetRamGB() uint32 {
	if m != nil {
		return m.RamGB
	}
	return 0
}

func (m *ServerRequest) GetCores() uint32 {
	if m != nil {
		return m.Cores
	}
	return 0
}

func init() {
	proto.RegisterType((*Nic)(nil), "vmmanager.Nic")
	proto.RegisterType((*Disk)(nil), "vmmanager.Disk")
	proto.RegisterType((*Server)(nil), "vmmanager.Server")
	proto.RegisterType((*ServerList)(nil), "vmmanager.ServerList")
	proto.RegisterType((*Image)(nil), "vmmanager.Image")
	proto.RegisterType((*CloneRequest)(nil), "vmmanager.CloneRequest")
	proto.RegisterType((*CloneResponse)(nil), "vmmanager.CloneResponse")
	proto.RegisterType((*ImageList)(nil), "vmmanager.ImageList")
	proto.RegisterType((*ServerRequest)(nil), "vmmanager.ServerRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for VMManager service

type VMManagerClient interface {
	RequestCitusServer(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error)
	RequestDefaultServer(ctx context.Context, in *ServerRequest, opts ...grpc.CallOption) (*common.Void, error)
	GetCitusServers(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*ServerList, error)
	ListImages(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*ImageList, error)
	ListServers(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*ServerList, error)
	RebootServers(ctx context.Context, in *ServerList, opts ...grpc.CallOption) (*common.Void, error)
	StartServers(ctx context.Context, in *ServerList, opts ...grpc.CallOption) (*common.Void, error)
	StopServers(ctx context.Context, in *ServerList, opts ...grpc.CallOption) (*common.Void, error)
	TriggerDHCPSync(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error)
	AttachStuffToSomeOtherStuff(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error)
	CloneVM(ctx context.Context, in *CloneRequest, opts ...grpc.CallOption) (*CloneResponse, error)
}

type vMManagerClient struct {
	cc *grpc.ClientConn
}

func NewVMManagerClient(cc *grpc.ClientConn) VMManagerClient {
	return &vMManagerClient{cc}
}

func (c *vMManagerClient) RequestCitusServer(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/RequestCitusServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) RequestDefaultServer(ctx context.Context, in *ServerRequest, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/RequestDefaultServer", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) GetCitusServers(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*ServerList, error) {
	out := new(ServerList)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/GetCitusServers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) ListImages(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*ImageList, error) {
	out := new(ImageList)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/ListImages", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) ListServers(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*ServerList, error) {
	out := new(ServerList)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/ListServers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) RebootServers(ctx context.Context, in *ServerList, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/RebootServers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) StartServers(ctx context.Context, in *ServerList, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/StartServers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) StopServers(ctx context.Context, in *ServerList, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/StopServers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) TriggerDHCPSync(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/TriggerDHCPSync", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) AttachStuffToSomeOtherStuff(ctx context.Context, in *common.Void, opts ...grpc.CallOption) (*common.Void, error) {
	out := new(common.Void)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/AttachStuffToSomeOtherStuff", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vMManagerClient) CloneVM(ctx context.Context, in *CloneRequest, opts ...grpc.CallOption) (*CloneResponse, error) {
	out := new(CloneResponse)
	err := grpc.Invoke(ctx, "/vmmanager.VMManager/CloneVM", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for VMManager service

type VMManagerServer interface {
	RequestCitusServer(context.Context, *common.Void) (*common.Void, error)
	RequestDefaultServer(context.Context, *ServerRequest) (*common.Void, error)
	GetCitusServers(context.Context, *common.Void) (*ServerList, error)
	ListImages(context.Context, *common.Void) (*ImageList, error)
	ListServers(context.Context, *common.Void) (*ServerList, error)
	RebootServers(context.Context, *ServerList) (*common.Void, error)
	StartServers(context.Context, *ServerList) (*common.Void, error)
	StopServers(context.Context, *ServerList) (*common.Void, error)
	TriggerDHCPSync(context.Context, *common.Void) (*common.Void, error)
	AttachStuffToSomeOtherStuff(context.Context, *common.Void) (*common.Void, error)
	CloneVM(context.Context, *CloneRequest) (*CloneResponse, error)
}

func RegisterVMManagerServer(s *grpc.Server, srv VMManagerServer) {
	s.RegisterService(&_VMManager_serviceDesc, srv)
}

func _VMManager_RequestCitusServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).RequestCitusServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/RequestCitusServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).RequestCitusServer(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_RequestDefaultServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).RequestDefaultServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/RequestDefaultServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).RequestDefaultServer(ctx, req.(*ServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_GetCitusServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).GetCitusServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/GetCitusServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).GetCitusServers(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_ListImages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).ListImages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/ListImages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).ListImages(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_ListServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).ListServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/ListServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).ListServers(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_RebootServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).RebootServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/RebootServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).RebootServers(ctx, req.(*ServerList))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_StartServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).StartServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/StartServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).StartServers(ctx, req.(*ServerList))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_StopServers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServerList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).StopServers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/StopServers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).StopServers(ctx, req.(*ServerList))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_TriggerDHCPSync_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).TriggerDHCPSync(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/TriggerDHCPSync",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).TriggerDHCPSync(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_AttachStuffToSomeOtherStuff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).AttachStuffToSomeOtherStuff(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/AttachStuffToSomeOtherStuff",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).AttachStuffToSomeOtherStuff(ctx, req.(*common.Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _VMManager_CloneVM_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloneRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VMManagerServer).CloneVM(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vmmanager.VMManager/CloneVM",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VMManagerServer).CloneVM(ctx, req.(*CloneRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _VMManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vmmanager.VMManager",
	HandlerType: (*VMManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestCitusServer",
			Handler:    _VMManager_RequestCitusServer_Handler,
		},
		{
			MethodName: "RequestDefaultServer",
			Handler:    _VMManager_RequestDefaultServer_Handler,
		},
		{
			MethodName: "GetCitusServers",
			Handler:    _VMManager_GetCitusServers_Handler,
		},
		{
			MethodName: "ListImages",
			Handler:    _VMManager_ListImages_Handler,
		},
		{
			MethodName: "ListServers",
			Handler:    _VMManager_ListServers_Handler,
		},
		{
			MethodName: "RebootServers",
			Handler:    _VMManager_RebootServers_Handler,
		},
		{
			MethodName: "StartServers",
			Handler:    _VMManager_StartServers_Handler,
		},
		{
			MethodName: "StopServers",
			Handler:    _VMManager_StopServers_Handler,
		},
		{
			MethodName: "TriggerDHCPSync",
			Handler:    _VMManager_TriggerDHCPSync_Handler,
		},
		{
			MethodName: "AttachStuffToSomeOtherStuff",
			Handler:    _VMManager_AttachStuffToSomeOtherStuff_Handler,
		},
		{
			MethodName: "CloneVM",
			Handler:    _VMManager_CloneVM_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "golang.conradwood.net/apis/vmmanager/vmmanager.proto",
}

func init() {
	proto.RegisterFile("golang.conradwood.net/apis/vmmanager/vmmanager.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 621 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x54, 0xdd, 0x4e, 0xdb, 0x30,
	0x18, 0x55, 0x69, 0xda, 0xae, 0x1f, 0x14, 0x98, 0xc5, 0x36, 0xab, 0xbb, 0x41, 0xd1, 0x36, 0x55,
	0x9b, 0x16, 0x50, 0x61, 0x93, 0x90, 0xf6, 0x23, 0x48, 0x24, 0x40, 0xa2, 0xd9, 0xe4, 0x20, 0xa4,
	0xdd, 0x4c, 0x32, 0xa9, 0x09, 0xd1, 0x48, 0xdc, 0xc5, 0x2e, 0x68, 0x7b, 0xb3, 0xbd, 0xce, 0x9e,
	0x64, 0xf2, 0x4f, 0x5b, 0x53, 0x10, 0x83, 0xab, 0xf8, 0xf8, 0xfb, 0x8e, 0xcf, 0xf1, 0x71, 0x6c,
	0xd8, 0xce, 0xf8, 0x05, 0x2d, 0xb3, 0x20, 0xe5, 0x65, 0x45, 0x87, 0x57, 0x9c, 0x0f, 0x83, 0x92,
	0xc9, 0x0d, 0x3a, 0xca, 0xc5, 0xc6, 0x65, 0x51, 0xd0, 0x92, 0x66, 0xac, 0x9a, 0x8d, 0x82, 0x51,
	0xc5, 0x25, 0x47, 0xed, 0xe9, 0x44, 0x37, 0xb8, 0x63, 0x81, 0x94, 0x17, 0x05, 0x2f, 0xed, 0xc7,
	0x50, 0xfd, 0x8f, 0x50, 0x8f, 0xf3, 0x14, 0xad, 0x42, 0xfd, 0x88, 0x96, 0xb8, 0xb6, 0x5e, 0xeb,
	0x75, 0x88, 0x1a, 0x22, 0x04, 0x5e, 0x4c, 0x0b, 0x86, 0x17, 0xd6, 0x6b, 0xbd, 0x36, 0xd1, 0x63,
	0xd5, 0x35, 0xd8, 0x0d, 0x71, 0x5d, 0x4f, 0xa9, 0xa1, 0xff, 0x1d, 0xbc, 0x28, 0x17, 0x3f, 0xd0,
	0x32, 0x2c, 0x1c, 0x46, 0x9a, 0xde, 0x26, 0x0b, 0x87, 0xd1, 0xad, 0x6c, 0x0c, 0xad, 0x83, 0x28,
	0x4a, 0xf2, 0xdf, 0x4c, 0xaf, 0xe0, 0x91, 0x09, 0x54, 0x95, 0xb0, 0x62, 0x54, 0xb2, 0x21, 0xf6,
	0x4c, 0xc5, 0x42, 0xff, 0x4f, 0x0d, 0x9a, 0x09, 0xab, 0x2e, 0x59, 0x75, 0x2f, 0x09, 0x1f, 0xbc,
	0x38, 0x4f, 0x05, 0xae, 0xaf, 0xd7, 0x7b, 0x8b, 0xfd, 0xe5, 0x60, 0x16, 0x54, 0x9c, 0xa7, 0x44,
	0xd7, 0xd0, 0x4b, 0x68, 0x28, 0xcb, 0x02, 0x7b, 0xba, 0x69, 0xc5, 0x69, 0x52, 0xf3, 0xc4, 0x54,
	0xd1, 0x1a, 0x34, 0x42, 0x5e, 0x31, 0x81, 0x1b, 0x3a, 0x13, 0x03, 0x54, 0x02, 0x84, 0x16, 0xb8,
	0xa9, 0x5d, 0xaa, 0xa1, 0xeb, 0xbd, 0x75, 0xdd, 0xfb, 0x0e, 0x80, 0xb1, 0x7e, 0x94, 0x0b, 0x89,
	0xde, 0x40, 0xcb, 0x20, 0x81, 0x6b, 0x5a, 0xf8, 0xb1, 0x23, 0x6c, 0x2a, 0x64, 0xd2, 0xe1, 0x7f,
	0x83, 0xc6, 0x61, 0x41, 0x33, 0x76, 0xaf, 0x4d, 0x23, 0xf0, 0x9c, 0x50, 0xbd, 0xff, 0x24, 0xfa,
	0x1a, 0x96, 0xc2, 0x0b, 0x5e, 0x32, 0xc2, 0x7e, 0x8e, 0x99, 0x90, 0xa8, 0x0b, 0x8f, 0x8c, 0xea,
	0x54, 0x67, 0x8a, 0xfd, 0xcf, 0xd0, 0xb1, 0xbd, 0x62, 0xc4, 0x4b, 0xc1, 0x54, 0x28, 0x31, 0xbb,
	0x9a, 0x76, 0x1a, 0xa0, 0xc4, 0x62, 0x76, 0xe5, 0xf8, 0x9a, 0x40, 0xff, 0x1d, 0xb4, 0xf5, 0x3e,
	0x74, 0x02, 0x3d, 0x68, 0x6a, 0x30, 0x09, 0x60, 0xd5, 0x09, 0x40, 0x17, 0x88, 0xad, 0xfb, 0x05,
	0x74, 0x6c, 0x22, 0xd6, 0xa4, 0xa3, 0x50, 0xbb, 0xa6, 0x80, 0x9e, 0x42, 0x53, 0x9d, 0xd7, 0xfe,
	0x9e, 0x96, 0xee, 0x10, 0x8b, 0x94, 0x53, 0x42, 0x8b, 0xfd, 0x3d, 0x9d, 0x4a, 0x87, 0x18, 0x30,
	0x3b, 0x54, 0xcf, 0x39, 0xd4, 0xfe, 0x5f, 0x0f, 0xda, 0x27, 0x83, 0x81, 0xb1, 0x82, 0x36, 0x01,
	0x59, 0xd9, 0x30, 0x97, 0x63, 0x61, 0xff, 0xbe, 0xa5, 0xc0, 0x5e, 0x9b, 0x13, 0x9e, 0x0f, 0xbb,
	0xd7, 0x10, 0xfa, 0x04, 0x6b, 0x96, 0x11, 0xb1, 0x33, 0x3a, 0xbe, 0x90, 0x96, 0x83, 0x6f, 0x9e,
	0xb0, 0x69, 0x9b, 0xe3, 0xbf, 0x87, 0x95, 0x7d, 0xe6, 0xaa, 0x89, 0x39, 0xb9, 0x27, 0x37, 0x16,
	0xd2, 0x81, 0x6e, 0x02, 0xa8, 0xaf, 0x09, 0x6d, 0x8e, 0xb2, 0x36, 0x1f, 0xae, 0x66, 0xf4, 0x61,
	0x51, 0x7d, 0x1f, 0xa4, 0xb2, 0x0d, 0x1d, 0xc2, 0x4e, 0x39, 0x9f, 0xb2, 0x6e, 0xef, 0x9b, 0xdb,
	0xd3, 0x16, 0x2c, 0x25, 0x92, 0x56, 0x0f, 0x23, 0xf5, 0x61, 0x31, 0x91, 0x7c, 0xf4, 0x20, 0xce,
	0x5b, 0x58, 0x39, 0xae, 0xf2, 0x2c, 0x63, 0x55, 0x74, 0x10, 0x7e, 0x4d, 0x7e, 0x95, 0xe9, 0x9d,
	0x67, 0xb5, 0x03, 0xcf, 0x77, 0xa5, 0xa4, 0xe9, 0x79, 0x22, 0xc7, 0x67, 0x67, 0xc7, 0x3c, 0xe1,
	0x05, 0xfb, 0x22, 0xcf, 0x59, 0xa5, 0xf1, 0x9d, 0xd4, 0x0f, 0xd0, 0xd2, 0xb7, 0xe1, 0x64, 0x80,
	0x9e, 0x39, 0xce, 0xdc, 0xdb, 0xd4, 0xc5, 0x37, 0x0b, 0xe6, 0xea, 0xec, 0xbd, 0x82, 0x17, 0x25,
	0x93, 0xee, 0xbb, 0x6c, 0x5f, 0x6a, 0xf5, 0x34, 0xcf, 0x58, 0xa7, 0x4d, 0xfd, 0x2e, 0x6f, 0xfd,
	0x0b, 0x00, 0x00, 0xff, 0xff, 0x18, 0x7d, 0xdd, 0x22, 0x0a, 0x06, 0x00, 0x00,
}
