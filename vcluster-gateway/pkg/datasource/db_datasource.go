package datasource

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"

	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/model"
)

var _ VClusterDBDataSource = &DBDataSource{}

func wrapDBError(err error) error {
	return errors.Wrap(err, "DB Error")
}

type DBDataSource struct {
	db *gorm.DB
}

func NewDBDataSource(db *gorm.DB) VClusterDBDataSource {
	dbs := &DBDataSource{
		db: db,
	}

	// 自动同步表结构
	//db.Set("gorm:table_options", "ENGINE=InnoDB")
	//db.Set("gorm:charset", "utf8mb4")
	//
	//err := db.AutoMigrate(&model.VCluster{}, &model.VStorage{}, &model.VGpu{}, &model.Workflow{})
	//if err != nil {
	//	panic(err)
	//}

	return dbs
}

// VCluster DataSource Functions
func (dbs *DBDataSource) FindVlusters(tenantType, tenantId string, rootClusterName string, isDeleted int32, permission []string) ([]model.VCluster, error) {
	var vclusters []model.VCluster
	var tx *gorm.DB

	switch tenantType {
	case "3", "TENANT_TYPE_PLATFORM":
		tx = dbs.db.Where("is_deleted = ?", isDeleted)
	default:
		tx = dbs.db.Where("tenant_id = ? and is_deleted = ?", tenantId, isDeleted)
	}

	tx = tx.Order("created_at DESC")

	if err := tx.Find(&vclusters).Error; err != nil {
		return nil, wrapDBError(err)
	}

	return vclusters, nil
}

func (dbs *DBDataSource) FindVClusterGPUs(vclusterId string) ([]model.VGpu, error) {
	var gpus []model.VGpu
	if err := dbs.db.Model(&model.VGpu{}).Where("cluster_id = ?", vclusterId).Scan(&gpus).Error; err != nil {
		return nil, wrapDBError(err)
	}

	return gpus, nil
}

func (dbs *DBDataSource) FindVClusterStorages(vclusterId string) ([]model.VStorage, error) {
	var storages []model.VStorage
	if err := dbs.db.Model(&model.VStorage{}).Where("vcluster_id = ?", vclusterId).Scan(&storages).Error; err != nil {
		return nil, wrapDBError(err)
	}

	return storages, nil
}

func (dbs *DBDataSource) FindVStorageByVClusterIdAndNameType(vclusterId, name, storageType string) ([]*model.VStorage, error) {
	var storages []*model.VStorage

	if err := dbs.db.Model(&model.VStorage{}).Where("vcluster_id = ? and name = ? and vstorage_type = ?",
		vclusterId, name, storageType).
		Find(&storages).Error; err != nil {
		return nil, wrapDBError(err)
	}

	return storages, nil
}

func (dbs *DBDataSource) GetVClusterById(vClusterId string) (*model.VCluster, error) {
	var vCluster model.VCluster

	err := dbs.db.Model(&model.VCluster{}).
		Where("vcluster_id = ?", vClusterId).
		First(&vCluster).Error
	if err != nil {
		return nil, err
	}

	return &vCluster, nil
}

func (dbs *DBDataSource) GetVClusterByInstanceId(instanceId string) (*model.VCluster, error) {
	var vCluster model.VCluster

	err := dbs.db.Model(&model.VCluster{}).
		Where("instance_id = ?", instanceId).
		First(&vCluster).Error
	if err != nil {
		return nil, err
	}

	return &vCluster, nil
}

func (dbs *DBDataSource) GetVStorageByVClusterId(vclusterId string) (*model.VStorage, error) {
	var vStorage model.VStorage

	err := dbs.db.Model(&model.VStorage{}).
		Where("vcluster_id = ?", vclusterId).
		First(&vStorage).Error
	if err != nil {
		return nil, err
	}

	return &vStorage, nil
}

func (dbs *DBDataSource) CheckVClusterNameExistByTenantId(tenantId string, name string) bool {
	var count int64

	err := dbs.db.Model(&model.VCluster{}).Unscoped().
		Where("tenant_id = ? and vcluster_name = ?", tenantId, name).
		Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		return true
	}

	return count > 0
}

func (dbs *DBDataSource) CheckVClusterExistByTenantIdAndVClusterId(vClusterId string, tenantId string) bool {
	var count int64
	err := dbs.db.Model(&model.VCluster{}).
		Where("vcluster_id = ? and tenant_id = ?", vClusterId, tenantId).
		Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("CheckVClusterExistByTenantIdAndVClusterId error: %v", err)
		}
		return false
	}

	return count > 0
}

func (dbs *DBDataSource) CheckVClusterExistById(vClusterId string) bool {
	var count int64

	err := dbs.db.Model(&model.VCluster{}).
		Where("vcluster_id = ?", vClusterId).
		Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		return true
	}

	return count > 0
}

func (dbs *DBDataSource) CheckVStorageExistByVClusterId(vclusterId string) bool {
	var count int64

	err := dbs.db.Model(&model.VStorage{}).
		Where("vcluster_id = ?", vclusterId).
		Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("CheckVStorageExistByVClusterId error: %v", err)
		}
		return false
	}

	return count > 0
}

func (dbs *DBDataSource) CheckVGpuExistByVClusterId(vclusterId string) bool {
	var count int64

	err := dbs.db.Model(&model.VGpu{}).
		Where("cluster_id = ?", vclusterId).
		Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("CheckVGpuExistByVClusterId error: %v", err)
		}
		return false
	}

	return count > 0
}

func (dbs *DBDataSource) CheckWorkflowExistById(workflowId string) bool {
	var count int64

	err := dbs.db.Model(&model.Workflow{}).
		Where("workflow_id = ?", workflowId).
		Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("CheckWorkflowExistById error: %v", err)
		}
		return false
	}

	return count > 0
}

// DeleteVClusterById 根据 vClusterId 删除 VCluster
func (dbs *DBDataSource) DeleteVClusterById(id string) error {
	err := dbs.db.Unscoped().Where("vcluster_id = ?", id).Delete(&model.VCluster{}).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteVStorageById 根据 vClusterId 删除 VStorage
func (dbs *DBDataSource) DeleteVStorageById(id string) error {
	err := dbs.db.Unscoped().Where("vcluster_id = ?", id).Delete(&model.VStorage{}).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteVGpuById 根据 clusterId 删除 VGpu
func (dbs *DBDataSource) DeleteVGpuById(id string) error {
	err := dbs.db.Unscoped().Where("cluster_id = ?", id).Delete(&model.VGpu{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DBDataSource) UpdateVCluster(oldModel *model.VCluster, needUpdate *model.VCluster) error {
	result := dbs.db.Model(oldModel).Updates(needUpdate)

	if result.Error != nil || result.RowsAffected == 0 {
		return errors.Wrapf(result.Error, "update vcluster failed, vcluster_id: %s", oldModel.VClusterId)
	}

	return nil
}

func (dbs *DBDataSource) UpdateVClusterSingle(vcluster *model.VCluster) error {
	result := dbs.db.Model(vcluster).Updates(vcluster)

	if result.Error != nil || result.RowsAffected == 0 {
		return errors.Wrapf(result.Error, "update vcluster failed, vcluster_id: %s", vcluster.VClusterId)
	}

	return nil
}

func (dbs *DBDataSource) UpdateVStorage(oldModel *model.VStorage, needUpdate *model.VStorage) error {
	result := dbs.db.Model(oldModel).Updates(needUpdate)

	if result.Error != nil || result.RowsAffected == 0 {
		return errors.Wrapf(result.Error, "update vStorage failed, vcluster_id: %s", oldModel.VClusterID)
	}

	return nil
}

func (dbs *DBDataSource) CreateVCluster(vc *model.VCluster) error {
	err := dbs.db.Create(vc).Error
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DBDataSource) CreateVStorages(storages []*model.VStorage) error {
	err := dbs.db.Create(storages).Error
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DBDataSource) CreateVGpus(gpus []*model.VGpu) error {
	err := dbs.db.Create(gpus).Error
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DBDataSource) CreateWorkflow(workflow *model.Workflow) error {
	err := dbs.db.Create(workflow).Error
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DBDataSource) ListWorkflows(status string) ([]model.Workflow, error) {
	workflowList := []model.Workflow{}

	db := dbs.db.Model(&model.Workflow{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Find(&workflowList).Error; err != nil {
		return nil, err
	}

	return workflowList, nil
}

func (dbs *DBDataSource) UpdateWorkflow(wf *model.Workflow) error {
	return dbs.db.Model(&model.Workflow{}).Where("workflow_id = ?", wf.WorkflowID).Updates(wf).Error
}

func (dbs *DBDataSource) UpdateWorkflowByInstanceId(ctx context.Context, instanceId string, wf *model.Workflow) error {
	err := dbs.db.Model(&model.Workflow{}).Where("instance_id = ?", instanceId).Updates(wf).Error
	if err != nil {
		log.Errorf("UpdateWorkflowByInstanceId error: %v", err)
		return err
	}
	return nil
}

func (dbs *DBDataSource) CheckVClusterNameExist(vclusterName string) bool {
	var count int64

	err := dbs.db.Model(&model.VCluster{}).
		Where("vcluster_name = ?", vclusterName).
		Count(&count).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("CheckVClusterNameExist error: %v", err)
		}
		return false
	}

	return count > 0
}

func (dbs *DBDataSource) CheckVClusterNameExistAndDeleted(vclusterName string, isDeleted int) bool {
	var count int64

	err := dbs.db.Model(&model.VCluster{}).
		Where("vcluster_name = ? AND is_deleted = ?", vclusterName, isDeleted).
		Count(&count).Error
	if err != nil {
		log.Errorf("CheckVClusterNameExistAndDeleted error: %v", err)
		return false
	}
	return count > 0
}

func (dbs *DBDataSource) CheckInstanceIdExist(id string) bool {
	var count int64

	err := dbs.db.Model(&model.VCluster{}).
		Where("instance_id = ?", id).
		Count(&count).Error
	if err != nil {
		log.Errorf("CheckInstanceIdExist error: %v", err)
		return false
	}

	return count > 0
}
