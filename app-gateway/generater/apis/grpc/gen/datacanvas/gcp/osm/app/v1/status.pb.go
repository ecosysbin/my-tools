// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: datacanvas/gcp/osm/app/v1/status.proto

package appv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VClusterStatus int32

const (
	VClusterStatus_VCLUSTER_STATUS_UNSPECIFIED VClusterStatus = 0
	VClusterStatus_VCLUSTER_STATUS_OK          VClusterStatus = 200
	// EN: Access Deny
	VClusterStatus_VCLUSTER_STATUS_ACCESS_DENY VClusterStatus = 142001
	// EN: Jwt Token Verify Error
	VClusterStatus_VCLUSTER_STATUS_JWT_TOKEN_VERIFY_ERROR VClusterStatus = 142002
	// EN: Create Cluster Error
	VClusterStatus_VCLUSTER_STATUS_CREATE_CLUSTER_ERROR VClusterStatus = 142003
	// EN: Select cluster error
	VClusterStatus_VCLUSTER_STATUS_SELECT_CLUSTER_ERROR VClusterStatus = 142004
	// EN: Delete cluster error
	VClusterStatus_VCLUSTER_STATUS_DELETE_CLUSTER_ERROR VClusterStatus = 142005
	// EN: Request token error
	VClusterStatus_VCLUSTER_STATUS_REQUEST_TOKEN_ERROR VClusterStatus = 142006
	// EN: request body Verify error
	VClusterStatus_VCLUSTER_STATUS_REQUEST_BODY_VERIFY_ERROR VClusterStatus = 142007
	// EN: Get cluster event error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_EVENT_ERROR VClusterStatus = 142008
	// EN: Get cluster namespace error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_NAMESPACE_ERROR VClusterStatus = 142009
	// EN: Get cluster deploy error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_DEPLOY_ERROR VClusterStatus = 142010
	// EN: Get cluster statefulset error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_STATEFULSET_ERROR VClusterStatus = 142011
	// EN: Get cluster ingress error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_INGRESS_ERROR VClusterStatus = 142012
	// EN: Get cluster pod error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_POD_ERROR VClusterStatus = 142013
	// EN: Get cluster secret error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_SECRET_ERROR VClusterStatus = 142014
	// EN: Get cluster service error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_SERVICE_ERROR VClusterStatus = 142015
	// EN: Get cluster configmap error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_CONFIGMAP_ERROR VClusterStatus = 142016
	// EN: Get cluster resource list error
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_RESOURCE_LIST_ERROR VClusterStatus = 142017
	// EN: Get cluster resource no exist
	VClusterStatus_VCLUSTER_STATUS_GET_CLUSTER_RESOURCE_NO_EXIST VClusterStatus = 142018
	// EN: Cluster structure assemble error
	VClusterStatus_VCLUSTER_STATUS_CLUSTER_STRUCTURE_ASSEMBLE_ERROR VClusterStatus = 142019
	// EN: Create Cluster Workflow Error
	VClusterStatus_VCLUSTER_STATUS_CREATE_CLUSTER_WORKFLOW_ERROR VClusterStatus = 142020
	// EN: Pause VCluster Error
	VClusterStatus_VCLUSTER_STATUS_PAUSE_VCLUSTER_ERROR VClusterStatus = 142021
	// EN: Resume VCluster Error
	VClusterStatus_VCLUSTER_STATUS_RESUME_VCLUSTER_ERROR VClusterStatus = 142022
	// EN: Update VCluster Service Error
	VClusterStatus_VCLUSTER_STATUS_UPDATE_SERVICE_ERROR VClusterStatus = 142023
	// EN: Create Cluster Error, cluster already exist
	VClusterStatus_VCLUSTER_STATUS_CLUSTER_ALREADY_EXIST_ERROR VClusterStatus = 142024
)

// Enum value maps for VClusterStatus.
var (
	VClusterStatus_name = map[int32]string{
		0:      "VCLUSTER_STATUS_UNSPECIFIED",
		200:    "VCLUSTER_STATUS_OK",
		142001: "VCLUSTER_STATUS_ACCESS_DENY",
		142002: "VCLUSTER_STATUS_JWT_TOKEN_VERIFY_ERROR",
		142003: "VCLUSTER_STATUS_CREATE_CLUSTER_ERROR",
		142004: "VCLUSTER_STATUS_SELECT_CLUSTER_ERROR",
		142005: "VCLUSTER_STATUS_DELETE_CLUSTER_ERROR",
		142006: "VCLUSTER_STATUS_REQUEST_TOKEN_ERROR",
		142007: "VCLUSTER_STATUS_REQUEST_BODY_VERIFY_ERROR",
		142008: "VCLUSTER_STATUS_GET_CLUSTER_EVENT_ERROR",
		142009: "VCLUSTER_STATUS_GET_CLUSTER_NAMESPACE_ERROR",
		142010: "VCLUSTER_STATUS_GET_CLUSTER_DEPLOY_ERROR",
		142011: "VCLUSTER_STATUS_GET_CLUSTER_STATEFULSET_ERROR",
		142012: "VCLUSTER_STATUS_GET_CLUSTER_INGRESS_ERROR",
		142013: "VCLUSTER_STATUS_GET_CLUSTER_POD_ERROR",
		142014: "VCLUSTER_STATUS_GET_CLUSTER_SECRET_ERROR",
		142015: "VCLUSTER_STATUS_GET_CLUSTER_SERVICE_ERROR",
		142016: "VCLUSTER_STATUS_GET_CLUSTER_CONFIGMAP_ERROR",
		142017: "VCLUSTER_STATUS_GET_CLUSTER_RESOURCE_LIST_ERROR",
		142018: "VCLUSTER_STATUS_GET_CLUSTER_RESOURCE_NO_EXIST",
		142019: "VCLUSTER_STATUS_CLUSTER_STRUCTURE_ASSEMBLE_ERROR",
		142020: "VCLUSTER_STATUS_CREATE_CLUSTER_WORKFLOW_ERROR",
		142021: "VCLUSTER_STATUS_PAUSE_VCLUSTER_ERROR",
		142022: "VCLUSTER_STATUS_RESUME_VCLUSTER_ERROR",
		142023: "VCLUSTER_STATUS_UPDATE_SERVICE_ERROR",
		142024: "VCLUSTER_STATUS_CLUSTER_ALREADY_EXIST_ERROR",
	}
	VClusterStatus_value = map[string]int32{
		"VCLUSTER_STATUS_UNSPECIFIED":                      0,
		"VCLUSTER_STATUS_OK":                               200,
		"VCLUSTER_STATUS_ACCESS_DENY":                      142001,
		"VCLUSTER_STATUS_JWT_TOKEN_VERIFY_ERROR":           142002,
		"VCLUSTER_STATUS_CREATE_CLUSTER_ERROR":             142003,
		"VCLUSTER_STATUS_SELECT_CLUSTER_ERROR":             142004,
		"VCLUSTER_STATUS_DELETE_CLUSTER_ERROR":             142005,
		"VCLUSTER_STATUS_REQUEST_TOKEN_ERROR":              142006,
		"VCLUSTER_STATUS_REQUEST_BODY_VERIFY_ERROR":        142007,
		"VCLUSTER_STATUS_GET_CLUSTER_EVENT_ERROR":          142008,
		"VCLUSTER_STATUS_GET_CLUSTER_NAMESPACE_ERROR":      142009,
		"VCLUSTER_STATUS_GET_CLUSTER_DEPLOY_ERROR":         142010,
		"VCLUSTER_STATUS_GET_CLUSTER_STATEFULSET_ERROR":    142011,
		"VCLUSTER_STATUS_GET_CLUSTER_INGRESS_ERROR":        142012,
		"VCLUSTER_STATUS_GET_CLUSTER_POD_ERROR":            142013,
		"VCLUSTER_STATUS_GET_CLUSTER_SECRET_ERROR":         142014,
		"VCLUSTER_STATUS_GET_CLUSTER_SERVICE_ERROR":        142015,
		"VCLUSTER_STATUS_GET_CLUSTER_CONFIGMAP_ERROR":      142016,
		"VCLUSTER_STATUS_GET_CLUSTER_RESOURCE_LIST_ERROR":  142017,
		"VCLUSTER_STATUS_GET_CLUSTER_RESOURCE_NO_EXIST":    142018,
		"VCLUSTER_STATUS_CLUSTER_STRUCTURE_ASSEMBLE_ERROR": 142019,
		"VCLUSTER_STATUS_CREATE_CLUSTER_WORKFLOW_ERROR":    142020,
		"VCLUSTER_STATUS_PAUSE_VCLUSTER_ERROR":             142021,
		"VCLUSTER_STATUS_RESUME_VCLUSTER_ERROR":            142022,
		"VCLUSTER_STATUS_UPDATE_SERVICE_ERROR":             142023,
		"VCLUSTER_STATUS_CLUSTER_ALREADY_EXIST_ERROR":      142024,
	}
)

func (x VClusterStatus) Enum() *VClusterStatus {
	p := new(VClusterStatus)
	*p = x
	return p
}

func (x VClusterStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VClusterStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_datacanvas_gcp_osm_app_v1_status_proto_enumTypes[0].Descriptor()
}

func (VClusterStatus) Type() protoreflect.EnumType {
	return &file_datacanvas_gcp_osm_app_v1_status_proto_enumTypes[0]
}

func (x VClusterStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VClusterStatus.Descriptor instead.
func (VClusterStatus) EnumDescriptor() ([]byte, []int) {
	return file_datacanvas_gcp_osm_app_v1_status_proto_rawDescGZIP(), []int{0}
}

var File_datacanvas_gcp_osm_app_v1_status_proto protoreflect.FileDescriptor

var file_datacanvas_gcp_osm_app_v1_status_proto_rawDesc = []byte{
	0x0a, 0x26, 0x64, 0x61, 0x74, 0x61, 0x63, 0x61, 0x6e, 0x76, 0x61, 0x73, 0x2f, 0x67, 0x63, 0x70,
	0x2f, 0x6f, 0x73, 0x6d, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x64, 0x61, 0x74, 0x61, 0x63, 0x61,
	0x6e, 0x76, 0x61, 0x73, 0x2e, 0x67, 0x63, 0x70, 0x2e, 0x6f, 0x73, 0x6d, 0x2e, 0x61, 0x70, 0x70,
	0x2e, 0x76, 0x31, 0x2a, 0xc5, 0x09, 0x0a, 0x0e, 0x56, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x1b, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54,
	0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x12, 0x56, 0x43, 0x4c, 0x55, 0x53,
	0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x4f, 0x4b, 0x10, 0xc8, 0x01,
	0x12, 0x21, 0x0a, 0x1b, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x5f, 0x44, 0x45, 0x4e, 0x59, 0x10,
	0xb1, 0xd5, 0x08, 0x12, 0x2c, 0x0a, 0x26, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x4a, 0x57, 0x54, 0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e,
	0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x59, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xb2, 0xd5,
	0x08, 0x12, 0x2a, 0x0a, 0x24, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x4c, 0x55, 0x53,
	0x54, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xb3, 0xd5, 0x08, 0x12, 0x2a, 0x0a,
	0x24, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x5f, 0x53, 0x45, 0x4c, 0x45, 0x43, 0x54, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xb4, 0xd5, 0x08, 0x12, 0x2a, 0x0a, 0x24, 0x56, 0x43, 0x4c,
	0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x45, 0x4c,
	0x45, 0x54, 0x45, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f,
	0x52, 0x10, 0xb5, 0xd5, 0x08, 0x12, 0x29, 0x0a, 0x23, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45,
	0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54,
	0x5f, 0x54, 0x4f, 0x4b, 0x45, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xb6, 0xd5, 0x08,
	0x12, 0x2f, 0x0a, 0x29, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f, 0x42, 0x4f, 0x44, 0x59,
	0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x59, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xb7, 0xd5,
	0x08, 0x12, 0x2d, 0x0a, 0x27, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52,
	0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xb8, 0xd5, 0x08,
	0x12, 0x31, 0x0a, 0x2b, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f,
	0x4e, 0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10,
	0xb9, 0xd5, 0x08, 0x12, 0x2e, 0x0a, 0x28, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54,
	0x45, 0x52, 0x5f, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10,
	0xba, 0xd5, 0x08, 0x12, 0x33, 0x0a, 0x2d, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54,
	0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x46, 0x55, 0x4c, 0x53, 0x45, 0x54, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0xbb, 0xd5, 0x08, 0x12, 0x2f, 0x0a, 0x29, 0x56, 0x43, 0x4c, 0x55,
	0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f,
	0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x47, 0x52, 0x45, 0x53, 0x53, 0x5f,
	0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xbc, 0xd5, 0x08, 0x12, 0x2b, 0x0a, 0x25, 0x56, 0x43, 0x4c,
	0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54,
	0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x50, 0x4f, 0x44, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0xbd, 0xd5, 0x08, 0x12, 0x2e, 0x0a, 0x28, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54,
	0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4c,
	0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x45, 0x43, 0x52, 0x45, 0x54, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0xbe, 0xd5, 0x08, 0x12, 0x2f, 0x0a, 0x29, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54,
	0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4c,
	0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x10, 0xbf, 0xd5, 0x08, 0x12, 0x31, 0x0a, 0x2b, 0x56, 0x43, 0x4c, 0x55, 0x53,
	0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x43,
	0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x4d, 0x41, 0x50,
	0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xc0, 0xd5, 0x08, 0x12, 0x35, 0x0a, 0x2f, 0x56, 0x43,
	0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45,
	0x54, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52,
	0x43, 0x45, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xc1, 0xd5,
	0x08, 0x12, 0x33, 0x0a, 0x2d, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x47, 0x45, 0x54, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52,
	0x5f, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x4e, 0x4f, 0x5f, 0x45, 0x58, 0x49,
	0x53, 0x54, 0x10, 0xc2, 0xd5, 0x08, 0x12, 0x36, 0x0a, 0x30, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54,
	0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45,
	0x52, 0x5f, 0x53, 0x54, 0x52, 0x55, 0x43, 0x54, 0x55, 0x52, 0x45, 0x5f, 0x41, 0x53, 0x53, 0x45,
	0x4d, 0x42, 0x4c, 0x45, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xc3, 0xd5, 0x08, 0x12, 0x33,
	0x0a, 0x2d, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52,
	0x5f, 0x57, 0x4f, 0x52, 0x4b, 0x46, 0x4c, 0x4f, 0x57, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10,
	0xc4, 0xd5, 0x08, 0x12, 0x2a, 0x0a, 0x24, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x50, 0x41, 0x55, 0x53, 0x45, 0x5f, 0x56, 0x43, 0x4c,
	0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xc5, 0xd5, 0x08, 0x12,
	0x2b, 0x0a, 0x25, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4d, 0x45, 0x5f, 0x56, 0x43, 0x4c, 0x55, 0x53, 0x54,
	0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xc6, 0xd5, 0x08, 0x12, 0x2a, 0x0a, 0x24,
	0x56, 0x43, 0x4c, 0x55, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f,
	0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x53, 0x45, 0x52, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0xc7, 0xd5, 0x08, 0x12, 0x31, 0x0a, 0x2b, 0x56, 0x43, 0x4c, 0x55,
	0x53, 0x54, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x43, 0x4c, 0x55, 0x53,
	0x54, 0x45, 0x52, 0x5f, 0x41, 0x4c, 0x52, 0x45, 0x41, 0x44, 0x59, 0x5f, 0x45, 0x58, 0x49, 0x53,
	0x54, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0xc8, 0xd5, 0x08, 0x42, 0x95, 0x02, 0x0a, 0x1d,
	0x63, 0x6f, 0x6d, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x63, 0x61, 0x6e, 0x76, 0x61, 0x73, 0x2e, 0x67,
	0x63, 0x70, 0x2e, 0x6f, 0x73, 0x6d, 0x2e, 0x61, 0x70, 0x70, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5e, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2e, 0x64, 0x61, 0x74, 0x61, 0x63, 0x61, 0x6e, 0x76, 0x61, 0x73, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x64, 0x63, 0x2f, 0x61, 0x70, 0x70, 0x2d, 0x67, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x72, 0x2f, 0x61,
	0x70, 0x69, 0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x64, 0x61, 0x74,
	0x61, 0x63, 0x61, 0x6e, 0x76, 0x61, 0x73, 0x2f, 0x67, 0x63, 0x70, 0x2f, 0x6f, 0x73, 0x6d, 0x2f,
	0x61, 0x70, 0x70, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x70, 0x70, 0x76, 0x31, 0xa2, 0x02, 0x04, 0x44,
	0x47, 0x4f, 0x41, 0xaa, 0x02, 0x19, 0x44, 0x61, 0x74, 0x61, 0x63, 0x61, 0x6e, 0x76, 0x61, 0x73,
	0x2e, 0x47, 0x63, 0x70, 0x2e, 0x4f, 0x73, 0x6d, 0x2e, 0x41, 0x70, 0x70, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x19, 0x44, 0x61, 0x74, 0x61, 0x63, 0x61, 0x6e, 0x76, 0x61, 0x73, 0x5c, 0x47, 0x63, 0x70,
	0x5c, 0x4f, 0x73, 0x6d, 0x5c, 0x41, 0x70, 0x70, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x25, 0x44, 0x61,
	0x74, 0x61, 0x63, 0x61, 0x6e, 0x76, 0x61, 0x73, 0x5c, 0x47, 0x63, 0x70, 0x5c, 0x4f, 0x73, 0x6d,
	0x5c, 0x41, 0x70, 0x70, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x1d, 0x44, 0x61, 0x74, 0x61, 0x63, 0x61, 0x6e, 0x76, 0x61, 0x73,
	0x3a, 0x3a, 0x47, 0x63, 0x70, 0x3a, 0x3a, 0x4f, 0x73, 0x6d, 0x3a, 0x3a, 0x41, 0x70, 0x70, 0x3a,
	0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_datacanvas_gcp_osm_app_v1_status_proto_rawDescOnce sync.Once
	file_datacanvas_gcp_osm_app_v1_status_proto_rawDescData = file_datacanvas_gcp_osm_app_v1_status_proto_rawDesc
)

func file_datacanvas_gcp_osm_app_v1_status_proto_rawDescGZIP() []byte {
	file_datacanvas_gcp_osm_app_v1_status_proto_rawDescOnce.Do(func() {
		file_datacanvas_gcp_osm_app_v1_status_proto_rawDescData = protoimpl.X.CompressGZIP(file_datacanvas_gcp_osm_app_v1_status_proto_rawDescData)
	})
	return file_datacanvas_gcp_osm_app_v1_status_proto_rawDescData
}

var file_datacanvas_gcp_osm_app_v1_status_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_datacanvas_gcp_osm_app_v1_status_proto_goTypes = []interface{}{
	(VClusterStatus)(0), // 0: datacanvas.gcp.osm.app.v1.VClusterStatus
}
var file_datacanvas_gcp_osm_app_v1_status_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_datacanvas_gcp_osm_app_v1_status_proto_init() }
func file_datacanvas_gcp_osm_app_v1_status_proto_init() {
	if File_datacanvas_gcp_osm_app_v1_status_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_datacanvas_gcp_osm_app_v1_status_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_datacanvas_gcp_osm_app_v1_status_proto_goTypes,
		DependencyIndexes: file_datacanvas_gcp_osm_app_v1_status_proto_depIdxs,
		EnumInfos:         file_datacanvas_gcp_osm_app_v1_status_proto_enumTypes,
	}.Build()
	File_datacanvas_gcp_osm_app_v1_status_proto = out.File
	file_datacanvas_gcp_osm_app_v1_status_proto_rawDesc = nil
	file_datacanvas_gcp_osm_app_v1_status_proto_goTypes = nil
	file_datacanvas_gcp_osm_app_v1_status_proto_depIdxs = nil
}
