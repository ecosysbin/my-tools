package repo

import (
	"strings"
	"time"

	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

const (
	ERR_RECORD_NOT_FOUND = "record not found"
	ADD_END_TIME         = " 23:59:59"
)

// // VirtualServerRepo -.
type AppRepoMysqlImpl struct {
	*gorm.DB
}

func NewAppMysqlImpl(datasource string) (*AppRepoMysqlImpl, error) {
	db, err := NewMysqlClient(datasource)
	if err != nil {
		return nil, err
	}
	// 数据库表结构映射，创建或更新表结构
	// db.AutoMigrate(&AppRecord{})
	// if err := db.AutoMigrate(&AppRecord{}, &AppConfig{}); err != nil {
	// 	return nil, err
	// }
	return &AppRepoMysqlImpl{db}, nil
}

type AppConfig struct {
	AppType    string `primaryKey;gorm:"size:255"`
	Domain     string `gorm:"size:255"`
	Version    string `gorm:"size:255"`
	KubeConfig string `gorm:"text"`
}

type AppRecord struct {
	Id            string `primaryKey;gorm:"size:255"`
	InstanceId    string `gorm:"size:255"`
	Name          string `gorm:"size:255"`
	TenantId      string `gorm:"size:255"`
	CreateUser    string `gorm:"size:255"`
	AppId         string `gorm:"size:255"`
	Url           string `gorm:"size:255"`
	MonitorUrl    string `gorm:"size:255"`
	Desc          string `gorm:"size:255"`
	ManageBy      string `gorm:"size:255"`
	Status        string `gorm:"size:255"`
	Conditions    string `gorm:"text"`
	Deleted       int32
	StartedTime   *time.Time
	CreateTime    *time.Time
	DeleteTime    *time.Time
	Reason        string `gorm:"text"`
	OriginMessage string `gorm:"text"`
	Message       string `gorm:"text"`
}

func (repo *AppRepoMysqlImpl) GetConfig(appType string) (AppConfig, error) {
	var config AppConfig
	result := repo.DB.Where("app_type = ?", appType).First(&config)
	if result.Error != nil {
		return AppConfig{}, errors.Errorf("get appType %s err, %v", appType, result.Error)
	}
	return config, nil
}

func (repo *AppRepoMysqlImpl) AddConfig(config AppConfig) error {
	// 执行创建或者更新操作（主键已存在）
	result := repo.DB.Create(&config)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate") {
			//处理违反唯一键约束错误
			return repo.DB.Where("app_type = ?", config.AppType).Updates(&config).Error
		}
		return errors.Errorf("add appType %s err, %v", config.AppType, result.Error)
	}
	return nil
}

func (repo *AppRepoMysqlImpl) DeleteConfig(appType string) error {
	return repo.DB.Where("app_type = ?", appType).Delete(&AppConfig{}).Error
}

func (repo *AppRepoMysqlImpl) ListAppConfig() ([]AppConfig, error) {
	var configs []AppConfig
	result := repo.DB.Find(&configs)

	if result.Error != nil {
		return nil, result.Error
	}
	return configs, nil
}

// 根据appId查询
func (repo *AppRepoMysqlImpl) GetByAppId(appId string) (AppRecord, error) {
	var app AppRecord
	// 后续考虑vc的id和appId统一
	result := repo.DB.Where("id = ?", appId).First(&app)
	if result.Error != nil {
		return AppRecord{}, errors.Errorf("get app %s err, %v", appId, result.Error)
	}
	return app, nil
}

// 更新app
func (repo *AppRepoMysqlImpl) Update(app AppRecord) error {
	result := repo.DB.Updates(&app)
	if result.Error != nil {
		return errors.Errorf("update app %s err, %v", app.Name, result.Error)
	}
	return nil
}

// 查询全部的App列表
func (repo *AppRepoMysqlImpl) ListAll(tenantId string) ([]AppRecord, error) {
	var apps []AppRecord
	var result *gorm.DB
	if tenantId == "" {
		result = repo.DB.Where("deleted = ?", 0).Find(&apps)
	} else {
		result = repo.DB.Where("deleted = ?", 0).Where("tenant_id = ?", tenantId).Order("create_time desc").Find(&apps)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return apps, nil
}

type ListOptions struct {
	TenantId   string
	Page       int
	PageSize   int
	Id         string
	Name       string
	InstanceId string
	Deleted    int32
	ManageBy   string
	CreateBy   string
	Status     AppStatus
	TenantIds  TenantIds
	CreateTime *PeriodTime
	DeleteTime *PeriodTime
}

type PeriodTime struct {
	StartTime string
	EndTime   string
}

type TenantIds []string

type AppStatus []string

func (repo *AppRepoMysqlImpl) ListPageAll(options ListOptions) ([]AppRecord, int64, error) {
	offset := (options.Page - 1) * options.PageSize
	var totalCount int64
	var apps []AppRecord

	var sql = repo.DB.Where("deleted = ?", options.Deleted)
	if options.TenantId != "" {
		sql = sql.Where("tenant_id = ?", options.TenantId)
		// 分类在TenantId不为空时才生效，admin时查询所有类型
	}

	if options.ManageBy == "vcluster" {
		sql = sql.Where("manage_by = ?", options.ManageBy)
	} else {
		sql = sql.Where("manage_by != ?", "vcluster") // 当前vcluster类型特殊处理
	}

	if options.Id != "" {
		sql = sql.Where("id like ?", "%"+options.Id+"%")
	}
	if options.Name != "" {
		sql = sql.Where("name like ?", "%"+options.Name+"%")
	}
	if options.InstanceId != "" {
		sql = sql.Where("instance_id like ?", "%"+options.InstanceId+"%")
	}
	if len(options.Status) > 0 {
		querySql := "status in ("
		for i := 0; i < len(options.Status); i++ {
			querySql += "'" + options.Status[i] + "'"
			if i < len(options.Status)-1 {
				querySql += ","
			}
		}
		querySql += ")"
		sql = sql.Where(querySql)
	}
	if len(options.TenantIds) > 0 {
		querySql := "tenant_id in ("
		for i := 0; i < len(options.TenantIds); i++ {
			querySql += "'" + options.TenantIds[i] + "'"
			if i < len(options.TenantIds)-1 {
				querySql += ","
			}
		}
		querySql += ")"
		sql = sql.Where(querySql)
	}
	// create time period
	if options.CreateTime.StartTime != "" {
		sql = sql.Where("create_time >= ?", options.CreateTime.StartTime)
	}
	if options.CreateTime.EndTime != "" {
		options.CreateTime.EndTime = options.CreateTime.EndTime + ADD_END_TIME
		sql = sql.Where("create_time <= ?", options.CreateTime.EndTime)
	}
	// delete time period
	if options.DeleteTime.StartTime != "" {
		sql = sql.Where("delete_time >= ?", options.DeleteTime.StartTime)
	}
	if options.DeleteTime.EndTime != "" {
		options.DeleteTime.EndTime = options.DeleteTime.EndTime + ADD_END_TIME
		sql = sql.Where("delete_time <= ?", options.DeleteTime.EndTime)
	}
	if options.CreateBy != "" {
		sql = sql.Where("create_user = ?", options.CreateBy)
	}
	// 计数
	countResult := sql.Model(&AppRecord{}).Count(&totalCount)
	if countResult.Error != nil {
		return nil, -1, countResult.Error
	}
	// 查询列表
	var listResult *gorm.DB
	sql = sql.Order("create_time desc")
	if options.Page == 0 {
		listResult = sql.Find(&apps)
	} else {
		listResult = sql.Offset(offset).Limit(options.PageSize).Find(&apps)
	}
	if listResult.Error != nil {
		return nil, -1, listResult.Error
	}

	return apps, totalCount, nil
}

func (repo *AppRepoMysqlImpl) Store(app AppRecord) error {
	return repo.DB.Create(&app).Error
}
