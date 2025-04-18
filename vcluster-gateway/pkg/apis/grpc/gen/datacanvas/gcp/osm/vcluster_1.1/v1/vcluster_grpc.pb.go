// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: datacanvas/gcp/osm/vcluster_1.1/v1/vcluster.proto

package vclusterv1

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
	VClusterGatewayService_CheckHealth_FullMethodName                = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/CheckHealth"
	VClusterGatewayService_VersionInformation_FullMethodName         = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/VersionInformation"
	VClusterGatewayService_CreateVCluster_FullMethodName             = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/CreateVCluster"
	VClusterGatewayService_UpdateVCluster_FullMethodName             = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/UpdateVCluster"
	VClusterGatewayService_DeleteVCluster_FullMethodName             = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/DeleteVCluster"
	VClusterGatewayService_GetKubeConfig_FullMethodName              = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/GetKubeConfig"
	VClusterGatewayService_GetKubeConfigBase64_FullMethodName        = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/GetKubeConfigBase64"
	VClusterGatewayService_PauseVCluster_FullMethodName              = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/PauseVCluster"
	VClusterGatewayService_QueryOperateStatus_FullMethodName         = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/QueryOperateStatus"
	VClusterGatewayService_ResumeVCluster_FullMethodName             = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/ResumeVCluster"
	VClusterGatewayService_GetVClusterStatus_FullMethodName          = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/GetVClusterStatus"
	VClusterGatewayService_GetVClusterResourceDetails_FullMethodName = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/GetVClusterResourceDetails"
	VClusterGatewayService_GetVClusterContainerID_FullMethodName     = "/datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService/GetVClusterContainerID"
)

// VClusterGatewayServiceClient is the client API for VClusterGatewayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VClusterGatewayServiceClient interface {
	// 检查服务健康状态
	CheckHealth(ctx context.Context, in *CheckHealthRequest, opts ...grpc.CallOption) (*CheckHealthResponse, error)
	// 版本信息
	VersionInformation(ctx context.Context, in *VersionInformationRequest, opts ...grpc.CallOption) (*VersionInformationResponse, error)
	// 创建 vCluster
	CreateVCluster(ctx context.Context, in *CreateVClusterRequest, opts ...grpc.CallOption) (*CreateVClusterResponse, error)
	// 更新 vCluster
	UpdateVCluster(ctx context.Context, in *UpdateVClusterRequest, opts ...grpc.CallOption) (*UpdateVClusterResponse, error)
	// 删除 vCluster
	DeleteVCluster(ctx context.Context, in *DeleteVClusterRequest, opts ...grpc.CallOption) (*DeleteVClusterResponse, error)
	// 获取 vCluster KubeConfig
	GetKubeConfig(ctx context.Context, in *GetKubeConfigRequest, opts ...grpc.CallOption) (*GetKubeConfigResponse, error)
	// 获取 vCluster KubeConfig, 返回值 base64 编码
	GetKubeConfigBase64(ctx context.Context, in *GetKubeConfigRequest, opts ...grpc.CallOption) (*GetKubeConfigBase64Response, error)
	// 暂停 vCluster
	PauseVCluster(ctx context.Context, in *PauseVClusterRequest, opts ...grpc.CallOption) (*PauseVClusterResponse, error)
	// 查询操作 vCluster 集群操作的状态
	// 例如 type=create, 返回的 status 为 success 即为成功
	QueryOperateStatus(ctx context.Context, in *QueryOperateStatusRequest, opts ...grpc.CallOption) (*QueryOperateStatusResponse, error)
	// 恢复 vCluster
	ResumeVCluster(ctx context.Context, in *ResumeVClusterRequest, opts ...grpc.CallOption) (*ResumeVClusterResponse, error)
	// 获取 vCluster 集群状态
	GetVClusterStatus(ctx context.Context, in *GetVClusterStatusRequest, opts ...grpc.CallOption) (*GetVClusterStatusResponse, error)
	// 获取 vCluster 集群资源配额
	GetVClusterResourceDetails(ctx context.Context, in *GetVClusterResourceDetailsRequest, opts ...grpc.CallOption) (*GetVClusterResourceDetailsResponse, error)
	// 获取指定 vCluster 的容器 ID
	GetVClusterContainerID(ctx context.Context, in *GetVClusterContainerIDRequest, opts ...grpc.CallOption) (*GetVClusterContainerIDResponse, error)
}

type vClusterGatewayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVClusterGatewayServiceClient(cc grpc.ClientConnInterface) VClusterGatewayServiceClient {
	return &vClusterGatewayServiceClient{cc}
}

func (c *vClusterGatewayServiceClient) CheckHealth(ctx context.Context, in *CheckHealthRequest, opts ...grpc.CallOption) (*CheckHealthResponse, error) {
	out := new(CheckHealthResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_CheckHealth_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) VersionInformation(ctx context.Context, in *VersionInformationRequest, opts ...grpc.CallOption) (*VersionInformationResponse, error) {
	out := new(VersionInformationResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_VersionInformation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) CreateVCluster(ctx context.Context, in *CreateVClusterRequest, opts ...grpc.CallOption) (*CreateVClusterResponse, error) {
	out := new(CreateVClusterResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_CreateVCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) UpdateVCluster(ctx context.Context, in *UpdateVClusterRequest, opts ...grpc.CallOption) (*UpdateVClusterResponse, error) {
	out := new(UpdateVClusterResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_UpdateVCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) DeleteVCluster(ctx context.Context, in *DeleteVClusterRequest, opts ...grpc.CallOption) (*DeleteVClusterResponse, error) {
	out := new(DeleteVClusterResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_DeleteVCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) GetKubeConfig(ctx context.Context, in *GetKubeConfigRequest, opts ...grpc.CallOption) (*GetKubeConfigResponse, error) {
	out := new(GetKubeConfigResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_GetKubeConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) GetKubeConfigBase64(ctx context.Context, in *GetKubeConfigRequest, opts ...grpc.CallOption) (*GetKubeConfigBase64Response, error) {
	out := new(GetKubeConfigBase64Response)
	err := c.cc.Invoke(ctx, VClusterGatewayService_GetKubeConfigBase64_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) PauseVCluster(ctx context.Context, in *PauseVClusterRequest, opts ...grpc.CallOption) (*PauseVClusterResponse, error) {
	out := new(PauseVClusterResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_PauseVCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) QueryOperateStatus(ctx context.Context, in *QueryOperateStatusRequest, opts ...grpc.CallOption) (*QueryOperateStatusResponse, error) {
	out := new(QueryOperateStatusResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_QueryOperateStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) ResumeVCluster(ctx context.Context, in *ResumeVClusterRequest, opts ...grpc.CallOption) (*ResumeVClusterResponse, error) {
	out := new(ResumeVClusterResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_ResumeVCluster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) GetVClusterStatus(ctx context.Context, in *GetVClusterStatusRequest, opts ...grpc.CallOption) (*GetVClusterStatusResponse, error) {
	out := new(GetVClusterStatusResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_GetVClusterStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) GetVClusterResourceDetails(ctx context.Context, in *GetVClusterResourceDetailsRequest, opts ...grpc.CallOption) (*GetVClusterResourceDetailsResponse, error) {
	out := new(GetVClusterResourceDetailsResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_GetVClusterResourceDetails_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vClusterGatewayServiceClient) GetVClusterContainerID(ctx context.Context, in *GetVClusterContainerIDRequest, opts ...grpc.CallOption) (*GetVClusterContainerIDResponse, error) {
	out := new(GetVClusterContainerIDResponse)
	err := c.cc.Invoke(ctx, VClusterGatewayService_GetVClusterContainerID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VClusterGatewayServiceServer is the server API for VClusterGatewayService service.
// All implementations must embed UnimplementedVClusterGatewayServiceServer
// for forward compatibility
type VClusterGatewayServiceServer interface {
	// 检查服务健康状态
	CheckHealth(context.Context, *CheckHealthRequest) (*CheckHealthResponse, error)
	// 版本信息
	VersionInformation(context.Context, *VersionInformationRequest) (*VersionInformationResponse, error)
	// 创建 vCluster
	CreateVCluster(context.Context, *CreateVClusterRequest) (*CreateVClusterResponse, error)
	// 更新 vCluster
	UpdateVCluster(context.Context, *UpdateVClusterRequest) (*UpdateVClusterResponse, error)
	// 删除 vCluster
	DeleteVCluster(context.Context, *DeleteVClusterRequest) (*DeleteVClusterResponse, error)
	// 获取 vCluster KubeConfig
	GetKubeConfig(context.Context, *GetKubeConfigRequest) (*GetKubeConfigResponse, error)
	// 获取 vCluster KubeConfig, 返回值 base64 编码
	GetKubeConfigBase64(context.Context, *GetKubeConfigRequest) (*GetKubeConfigBase64Response, error)
	// 暂停 vCluster
	PauseVCluster(context.Context, *PauseVClusterRequest) (*PauseVClusterResponse, error)
	// 查询操作 vCluster 集群操作的状态
	// 例如 type=create, 返回的 status 为 success 即为成功
	QueryOperateStatus(context.Context, *QueryOperateStatusRequest) (*QueryOperateStatusResponse, error)
	// 恢复 vCluster
	ResumeVCluster(context.Context, *ResumeVClusterRequest) (*ResumeVClusterResponse, error)
	// 获取 vCluster 集群状态
	GetVClusterStatus(context.Context, *GetVClusterStatusRequest) (*GetVClusterStatusResponse, error)
	// 获取 vCluster 集群资源配额
	GetVClusterResourceDetails(context.Context, *GetVClusterResourceDetailsRequest) (*GetVClusterResourceDetailsResponse, error)
	// 获取指定 vCluster 的容器 ID
	GetVClusterContainerID(context.Context, *GetVClusterContainerIDRequest) (*GetVClusterContainerIDResponse, error)
	mustEmbedUnimplementedVClusterGatewayServiceServer()
}

// UnimplementedVClusterGatewayServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVClusterGatewayServiceServer struct {
}

func (UnimplementedVClusterGatewayServiceServer) CheckHealth(context.Context, *CheckHealthRequest) (*CheckHealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckHealth not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) VersionInformation(context.Context, *VersionInformationRequest) (*VersionInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VersionInformation not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) CreateVCluster(context.Context, *CreateVClusterRequest) (*CreateVClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVCluster not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) UpdateVCluster(context.Context, *UpdateVClusterRequest) (*UpdateVClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVCluster not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) DeleteVCluster(context.Context, *DeleteVClusterRequest) (*DeleteVClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVCluster not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) GetKubeConfig(context.Context, *GetKubeConfigRequest) (*GetKubeConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKubeConfig not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) GetKubeConfigBase64(context.Context, *GetKubeConfigRequest) (*GetKubeConfigBase64Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKubeConfigBase64 not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) PauseVCluster(context.Context, *PauseVClusterRequest) (*PauseVClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PauseVCluster not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) QueryOperateStatus(context.Context, *QueryOperateStatusRequest) (*QueryOperateStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryOperateStatus not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) ResumeVCluster(context.Context, *ResumeVClusterRequest) (*ResumeVClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResumeVCluster not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) GetVClusterStatus(context.Context, *GetVClusterStatusRequest) (*GetVClusterStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVClusterStatus not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) GetVClusterResourceDetails(context.Context, *GetVClusterResourceDetailsRequest) (*GetVClusterResourceDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVClusterResourceDetails not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) GetVClusterContainerID(context.Context, *GetVClusterContainerIDRequest) (*GetVClusterContainerIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVClusterContainerID not implemented")
}
func (UnimplementedVClusterGatewayServiceServer) mustEmbedUnimplementedVClusterGatewayServiceServer() {
}

// UnsafeVClusterGatewayServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VClusterGatewayServiceServer will
// result in compilation errors.
type UnsafeVClusterGatewayServiceServer interface {
	mustEmbedUnimplementedVClusterGatewayServiceServer()
}

func RegisterVClusterGatewayServiceServer(s grpc.ServiceRegistrar, srv VClusterGatewayServiceServer) {
	s.RegisterService(&VClusterGatewayService_ServiceDesc, srv)
}

func _VClusterGatewayService_CheckHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckHealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).CheckHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_CheckHealth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).CheckHealth(ctx, req.(*CheckHealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_VersionInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VersionInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).VersionInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_VersionInformation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).VersionInformation(ctx, req.(*VersionInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_CreateVCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateVClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).CreateVCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_CreateVCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).CreateVCluster(ctx, req.(*CreateVClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_UpdateVCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).UpdateVCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_UpdateVCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).UpdateVCluster(ctx, req.(*UpdateVClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_DeleteVCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).DeleteVCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_DeleteVCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).DeleteVCluster(ctx, req.(*DeleteVClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_GetKubeConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKubeConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).GetKubeConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_GetKubeConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).GetKubeConfig(ctx, req.(*GetKubeConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_GetKubeConfigBase64_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKubeConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).GetKubeConfigBase64(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_GetKubeConfigBase64_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).GetKubeConfigBase64(ctx, req.(*GetKubeConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_PauseVCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PauseVClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).PauseVCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_PauseVCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).PauseVCluster(ctx, req.(*PauseVClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_QueryOperateStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOperateStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).QueryOperateStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_QueryOperateStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).QueryOperateStatus(ctx, req.(*QueryOperateStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_ResumeVCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResumeVClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).ResumeVCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_ResumeVCluster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).ResumeVCluster(ctx, req.(*ResumeVClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_GetVClusterStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVClusterStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).GetVClusterStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_GetVClusterStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).GetVClusterStatus(ctx, req.(*GetVClusterStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_GetVClusterResourceDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVClusterResourceDetailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).GetVClusterResourceDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_GetVClusterResourceDetails_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).GetVClusterResourceDetails(ctx, req.(*GetVClusterResourceDetailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VClusterGatewayService_GetVClusterContainerID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVClusterContainerIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VClusterGatewayServiceServer).GetVClusterContainerID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: VClusterGatewayService_GetVClusterContainerID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VClusterGatewayServiceServer).GetVClusterContainerID(ctx, req.(*GetVClusterContainerIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VClusterGatewayService_ServiceDesc is the grpc.ServiceDesc for VClusterGatewayService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VClusterGatewayService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "datacanvas.gcp.osm.vcluster.v1.VClusterGatewayService",
	HandlerType: (*VClusterGatewayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckHealth",
			Handler:    _VClusterGatewayService_CheckHealth_Handler,
		},
		{
			MethodName: "VersionInformation",
			Handler:    _VClusterGatewayService_VersionInformation_Handler,
		},
		{
			MethodName: "CreateVCluster",
			Handler:    _VClusterGatewayService_CreateVCluster_Handler,
		},
		{
			MethodName: "UpdateVCluster",
			Handler:    _VClusterGatewayService_UpdateVCluster_Handler,
		},
		{
			MethodName: "DeleteVCluster",
			Handler:    _VClusterGatewayService_DeleteVCluster_Handler,
		},
		{
			MethodName: "GetKubeConfig",
			Handler:    _VClusterGatewayService_GetKubeConfig_Handler,
		},
		{
			MethodName: "GetKubeConfigBase64",
			Handler:    _VClusterGatewayService_GetKubeConfigBase64_Handler,
		},
		{
			MethodName: "PauseVCluster",
			Handler:    _VClusterGatewayService_PauseVCluster_Handler,
		},
		{
			MethodName: "QueryOperateStatus",
			Handler:    _VClusterGatewayService_QueryOperateStatus_Handler,
		},
		{
			MethodName: "ResumeVCluster",
			Handler:    _VClusterGatewayService_ResumeVCluster_Handler,
		},
		{
			MethodName: "GetVClusterStatus",
			Handler:    _VClusterGatewayService_GetVClusterStatus_Handler,
		},
		{
			MethodName: "GetVClusterResourceDetails",
			Handler:    _VClusterGatewayService_GetVClusterResourceDetails_Handler,
		},
		{
			MethodName: "GetVClusterContainerID",
			Handler:    _VClusterGatewayService_GetVClusterContainerID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "datacanvas/gcp/osm/vcluster_1.1/v1/vcluster.proto",
}
