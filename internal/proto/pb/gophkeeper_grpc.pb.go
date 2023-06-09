// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.4
// source: internal/proto/gophkeeper.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GophKeeperService_Register_FullMethodName       = "/pb.GophKeeperService/Register"
	GophKeeperService_Authenticate_FullMethodName   = "/pb.GophKeeperService/Authenticate"
	GophKeeperService_Authorize_FullMethodName      = "/pb.GophKeeperService/Authorize"
	GophKeeperService_AddData_FullMethodName        = "/pb.GophKeeperService/AddData"
	GophKeeperService_RetrieveData_FullMethodName   = "/pb.GophKeeperService/RetrieveData"
	GophKeeperService_GetInformation_FullMethodName = "/pb.GophKeeperService/GetInformation"
)

// GophKeeperServiceClient is the client API for GophKeeperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GophKeeperServiceClient interface {
	Register(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*RegistrationResponse, error)
	Authenticate(ctx context.Context, in *AuthenticationRequest, opts ...grpc.CallOption) (*AuthenticationResponse, error)
	Authorize(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationResponse, error)
	AddData(ctx context.Context, in *AddDataRequest, opts ...grpc.CallOption) (*AddDataResponse, error)
	RetrieveData(ctx context.Context, in *RetrieveDataRequest, opts ...grpc.CallOption) (*RetrieveDataResponse, error)
	GetInformation(ctx context.Context, in *GetInformationRequest, opts ...grpc.CallOption) (*GetInformationResponse, error)
}

type gophKeeperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGophKeeperServiceClient(cc grpc.ClientConnInterface) GophKeeperServiceClient {
	return &gophKeeperServiceClient{cc}
}

func (c *gophKeeperServiceClient) Register(ctx context.Context, in *RegistrationRequest, opts ...grpc.CallOption) (*RegistrationResponse, error) {
	out := new(RegistrationResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) Authenticate(ctx context.Context, in *AuthenticationRequest, opts ...grpc.CallOption) (*AuthenticationResponse, error) {
	out := new(AuthenticationResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_Authenticate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) Authorize(ctx context.Context, in *AuthorizationRequest, opts ...grpc.CallOption) (*AuthorizationResponse, error) {
	out := new(AuthorizationResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_Authorize_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) AddData(ctx context.Context, in *AddDataRequest, opts ...grpc.CallOption) (*AddDataResponse, error) {
	out := new(AddDataResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_AddData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) RetrieveData(ctx context.Context, in *RetrieveDataRequest, opts ...grpc.CallOption) (*RetrieveDataResponse, error) {
	out := new(RetrieveDataResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_RetrieveData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gophKeeperServiceClient) GetInformation(ctx context.Context, in *GetInformationRequest, opts ...grpc.CallOption) (*GetInformationResponse, error) {
	out := new(GetInformationResponse)
	err := c.cc.Invoke(ctx, GophKeeperService_GetInformation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GophKeeperServiceServer is the server API for GophKeeperService service.
// All implementations must embed UnimplementedGophKeeperServiceServer
// for forward compatibility
type GophKeeperServiceServer interface {
	Register(context.Context, *RegistrationRequest) (*RegistrationResponse, error)
	Authenticate(context.Context, *AuthenticationRequest) (*AuthenticationResponse, error)
	Authorize(context.Context, *AuthorizationRequest) (*AuthorizationResponse, error)
	AddData(context.Context, *AddDataRequest) (*AddDataResponse, error)
	RetrieveData(context.Context, *RetrieveDataRequest) (*RetrieveDataResponse, error)
	GetInformation(context.Context, *GetInformationRequest) (*GetInformationResponse, error)
	mustEmbedUnimplementedGophKeeperServiceServer()
}

// UnimplementedGophKeeperServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGophKeeperServiceServer struct {
}

func (UnimplementedGophKeeperServiceServer) Register(context.Context, *RegistrationRequest) (*RegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedGophKeeperServiceServer) Authenticate(context.Context, *AuthenticationRequest) (*AuthenticationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedGophKeeperServiceServer) Authorize(context.Context, *AuthorizationRequest) (*AuthorizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (UnimplementedGophKeeperServiceServer) AddData(context.Context, *AddDataRequest) (*AddDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddData not implemented")
}
func (UnimplementedGophKeeperServiceServer) RetrieveData(context.Context, *RetrieveDataRequest) (*RetrieveDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveData not implemented")
}
func (UnimplementedGophKeeperServiceServer) GetInformation(context.Context, *GetInformationRequest) (*GetInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetInformation not implemented")
}
func (UnimplementedGophKeeperServiceServer) mustEmbedUnimplementedGophKeeperServiceServer() {}

// UnsafeGophKeeperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GophKeeperServiceServer will
// result in compilation errors.
type UnsafeGophKeeperServiceServer interface {
	mustEmbedUnimplementedGophKeeperServiceServer()
}

func RegisterGophKeeperServiceServer(s grpc.ServiceRegistrar, srv GophKeeperServiceServer) {
	s.RegisterService(&GophKeeperService_ServiceDesc, srv)
}

func _GophKeeperService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegistrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).Register(ctx, req.(*RegistrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_Authenticate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).Authenticate(ctx, req.(*AuthenticationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_Authorize_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).Authorize(ctx, req.(*AuthorizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_AddData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).AddData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_AddData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).AddData(ctx, req.(*AddDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_RetrieveData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).RetrieveData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_RetrieveData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).RetrieveData(ctx, req.(*RetrieveDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GophKeeperService_GetInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GophKeeperServiceServer).GetInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GophKeeperService_GetInformation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GophKeeperServiceServer).GetInformation(ctx, req.(*GetInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GophKeeperService_ServiceDesc is the grpc.ServiceDesc for GophKeeperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GophKeeperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.GophKeeperService",
	HandlerType: (*GophKeeperServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _GophKeeperService_Register_Handler,
		},
		{
			MethodName: "Authenticate",
			Handler:    _GophKeeperService_Authenticate_Handler,
		},
		{
			MethodName: "Authorize",
			Handler:    _GophKeeperService_Authorize_Handler,
		},
		{
			MethodName: "AddData",
			Handler:    _GophKeeperService_AddData_Handler,
		},
		{
			MethodName: "RetrieveData",
			Handler:    _GophKeeperService_RetrieveData_Handler,
		},
		{
			MethodName: "GetInformation",
			Handler:    _GophKeeperService_GetInformation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/gophkeeper.proto",
}