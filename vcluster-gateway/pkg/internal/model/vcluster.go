package model

import (
	"encoding/base64"
	"time"

	"gorm.io/gorm"
)

type VCluster struct {
	gorm.Model

	// 时间戳字段
	CreateTime *time.Time `gorm:"column:create_time" json:"createTime"`
	DeleteTime *time.Time `gorm:"column:delete_time" json:"deleteTime"`
	StartTime  *time.Time `gorm:"column:started_time" json:"startedTime"`

	// 标识和管理字段
	VClusterId string `gorm:"column:vcluster_id;size:12;uniqueIndex" json:"id"` // 唯一标识
	InstanceId string `gorm:"column:instance_id;size:36;uniqueIndex" json:"instanceId"`
	ManageBy   string `gorm:"column:manage_by;size:16;default:'raw'" json:"manageBy"`
	UserName   string `gorm:"column:user_name;size:60" json:"-" swaggerignore:"true"`
	TenantId   string `gorm:"column:tenant_id;size:36" json:"-" swaggerignore:"true"`

	// 状态字段
	Status       string `gorm:"column:status" json:"status"`              // vCluster 的内部状态
	ServerStatus string `gorm:"column:server_status" json:"serverStatus"` // vCluster 的业务状态
	IsDeleted    int    `gorm:"column:is_deleted;size:1" json:"-" swaggerignore:"true"`
	Reason       string `gorm:"column:reason;type:text" json:"reason"`

	// 描述字段
	VClusterName    string `gorm:"column:vcluster_name;size:30;uniqueIndex" json:"name"`
	Namespace       string `gorm:"column:namespace;size:24" json:"namespace"`
	RootClusterName string `gorm:"column:root_cluster_name;size:30" json:"context"`
	Comment         string `gorm:"column:comment;size:100" json:"comment"`
	InstanceSpec    string `gorm:"column:instance_spec;type:text" json:"instanceSpec"`
}

func (vc *VCluster) DecodeInstanceSpec() {
	data, _ := base64.StdEncoding.DecodeString(vc.InstanceSpec)
	vc.InstanceSpec = string(data)
}

// Workflow Status:
//  1. Pending
//  2. Running
//  3. Succeeded
//  4. Failed
//  5. Paused
//
// Workflow StepStatus:
//  1. Pending
//  2. Running
//  3. Completed
//  4. Failed
type Workflow struct {
	gorm.Model
	WorkflowID   string `gorm:"column:workflow_id;size:36;uniqueIndex;not null;default:''" json:"workflow_id"`
	Name         string `gorm:"column:name;size:36;default:''" json:"name"`
	VClusterID   string `gorm:"column:vcluster_id;size:36;not null;default:''" json:"vcluster_id"`
	InstanceId   string `gorm:"column:instance_id;size:64;not null;default:''" json:"instance_id"`
	VClusterName string `gorm:"column:vcluster_name;size:64;not null;default:''" json:"vcluster_name"`
	UserName     string `gorm:"column:user_name;size:60;default:''" json:"user_name"`
	TenantId     string `gorm:"column:tenant_id;size:36;default:''" json:"tenant_id"`
	Status       string `gorm:"column:status;size:20;default:''" json:"status"`
	StepStatus   string `gorm:"column:step_status;size:20;default:''" json:"step_status"`     // running, completed, failed
	WorkflowStep string `gorm:"column:workflow_step;size:10;default:''" json:"workflow_step"` // e.g., "1/3"
	Cleanup      bool   `gorm:"column:cleanup;default:false" json:"cleanup"`
}

type VStorage struct {
	gorm.Model
	VClusterID       string `gorm:"column:vcluster_id;index;size:12;not null;default:''"`
	VStorageType     string `gorm:"column:vstorage_type;index;size:60;not null;default:''"`
	VStorageCapacity int    `gorm:"column:vstorage_capacity;size:60;not null;default:0"`
	IsDeleted        int    `gorm:"column:is_deleted;size:1;default:0"`
	Name             string `gorm:"column:name;index;size:60;default:''"`
}

type VGpu struct {
	gorm.Model
	ClusterID    string `gorm:"index;size:12;not null;default:''"`
	GpuType      string `gorm:"size:50;not null;default:''"`
	ResourceName string `gorm:"size:50;not null;default:''"`
}

// The TableName function is a custom table name method implemented in GORM.
func (u VCluster) TableName() string {
	return "vcluster"
}

// The TableName function is a custom table name method implemented in GORM.
func (Workflow) TableName() string {
	return "workflow"
}

// The TableName function is a custom table name method implemented in GORM.
func (u VStorage) TableName() string {
	return "vstorage"
}

// The TableName function is a custom table name method implemented in GORM.
func (u VGpu) TableName() string {
	return "vgpu"
}
