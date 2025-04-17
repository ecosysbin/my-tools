package repo

import (
	"fmt"
	"time"

	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"
	"gorm.io/gorm"
)

const (
	ERR_RECORD_NOT_FOUND = "record not found"
)

// // VirtualServerRepo -.
type VirtualServerMysqlImpl struct {
	*gorm.DB
}

func NewVirtualServerMysqlImpl(datasource string) (*VirtualServerMysqlImpl, error) {
	db, err := NewMysqlClient(datasource)
	if err != nil {
		return nil, err
	}
	// 数据库表结构映射，创建或更新表结构
	db.AutoMigrate(&VirtualServer{})
	if err := db.AutoMigrate(&VirtualServer{}); err != nil {
		log.Infof("register db err,%v", err)
		return nil, err
	}
	return &VirtualServerMysqlImpl{db}, nil
}

// // // New -.
// func NewVirtualServerRepo(db *gorm.DB) pkg.VirtualServerRepo {
// 	// 数据库表结构映射，创建或更新表结构
// 	db.AutoMigrate(&VirtualServer{})
// 	if err := db.AutoMigrate(&VirtualServer{}); err != nil {
// 		log.Infof("register db err,%v", err)
// 		return nil
// 	}
// 	return &VirtualServerRepoImpl{db}
// }

type VirtualServer struct {
	Id          string `gorm:"size:255"`
	Name        string `gorm:"size:255"`
	Image       string `gorm:"size:255"`
	CreateUser  string `gorm:"size:255"`
	Product     string `gorm:"size:255"`
	Desc        string `gorm:"size:255"`
	SshPort     int32
	Status      string `gorm:"size:255"`
	Deleted     int32
	StartedTime *time.Time `gorm:"column:start_time"`
	CreateTime  *time.Time `gorm:"column:create_time"`
	DeleteTime  *time.Time `gorm:"column:delete_time"`
	Reason      string     `gorm:"size:255"`
	Message     string     `gorm:"text"`
}

func (repo *VirtualServerMysqlImpl) GetVmById(userName string, id string) (VirtualServer, error) {
	var vm VirtualServer
	var result *gorm.DB
	if userName == "" {
		result = repo.DB.Where("id = ?", id).First(&vm)
	} else {
		result = repo.DB.Where("id = ?", id).Where("create_user = ?", userName).First(&vm)
	}

	if result.Error != nil {
		return VirtualServer{}, fmt.Errorf("get vm %s err, %v", vm.Name, result.Error)
	}
	return vm, nil
}

func (repo *VirtualServerMysqlImpl) Update(vm *VirtualServer) error {
	result := repo.DB.Updates(vm)
	if result.Error != nil {
		return fmt.Errorf("update vm %s err, %v", vm.Name, result.Error)
	}
	return nil
}

// 查询全部的虚拟机列表
func (repo *VirtualServerMysqlImpl) ListAll(userName string) ([]VirtualServer, error) {
	var vms []VirtualServer
	var result *gorm.DB
	if userName == "" {
		result = repo.DB.Find(&vms)
	} else {
		result = repo.DB.Where("create_user = ?", userName).Find(&vms)

	}
	if result.Error != nil {
		return nil, result.Error
	}
	return vms, nil
}

// 查询vm时默认都不会查询已删除的虚拟机
func (repo *VirtualServerMysqlImpl) List(userName string) ([]VirtualServer, error) {
	var vms []VirtualServer
	var result *gorm.DB
	if userName == "" {
		result = repo.DB.Where("deleted = ?", 0).Find(&vms)
	} else {
		result = repo.DB.Where("deleted = ?", 0).Where("create_user = ?", userName).Find(&vms)

	}
	if result.Error != nil {
		return nil, result.Error
	}
	return vms, nil
}

// 查询已删除的虚拟机
func (repo *VirtualServerMysqlImpl) ListDeletedVms(userName string) ([]VirtualServer, error) {
	var vms []VirtualServer
	var result *gorm.DB
	if userName == "" {
		result = repo.DB.Where("deleted = ?", 1).Find(&vms)
	} else {
		result = repo.DB.Where("deleted = ?", 1).Where("create_user = ?", userName).Find(&vms)

	}
	if result.Error != nil {
		return nil, result.Error
	}
	return vms, nil
}

func (repo *VirtualServerMysqlImpl) Store(vm VirtualServer) error {
	return repo.DB.Create(&vm).Error
}
