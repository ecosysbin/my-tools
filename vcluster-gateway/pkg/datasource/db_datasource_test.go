package datasource

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/dig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"

	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/internal/model"
)

func setupTestContainer(ctx context.Context) (string, func(), error) {
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0.36",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "testdb",
		},
		WaitingFor: wait.ForLog("ready for connections").WithPollInterval(1 * time.Second),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, err
	}

	host, err := mysqlC.Host(ctx)
	if err != nil {
		return "", nil, err
	}

	port, err := mysqlC.MappedPort(ctx, "3306")
	if err != nil {
		return "", nil, err
	}

	dsn := fmt.Sprintf("root:root@tcp(%s:%s)/testdb?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())
	tearDown := func() {
		if err := mysqlC.Terminate(ctx); err != nil {
			log.Fatalf("Error terminating container: %s", err)
		}
	}

	return dsn, tearDown, nil
}

var terminateMysqlContainer func() = func() {}

func initGormDB() (*gorm.DB, error) {
	ctx := context.Background()
	dsn, tearDown, err := setupTestContainer(ctx)
	if err != nil {
		log.Fatalf("Could not set up test container: %s", err)
	}

	terminateMysqlContainer = tearDown

	// wait mysql init
	time.Sleep(5 * time.Second)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

var diContainer = dig.New()

func init() {
	diContainer.Provide(initGormDB)
	diContainer.Provide(NewDBDataSource)
}

func TestMain(m *testing.M) {
	// Run the tests
	code := m.Run()

	if terminateMysqlContainer != nil {
		terminateMysqlContainer()
	}

	// Exit with the correct code
	os.Exit(code)
}

func TestDBDataSourceFindVlusters(t *testing.T) {
	modelVclusters := []model.VCluster{
		{
			Model: gorm.Model{
				ID: 1,
			},
			UserName:        "test_user",
			TenantId:        "tenant_id_1",
			VClusterId:      "vclusterid1",
			VClusterName:    "vcluster_name_1",
			RootClusterName: "root_cluster_name_1",
			IsDeleted:       0,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			UserName:        "test_user",
			TenantId:        "tenant_id_2",
			VClusterId:      "vclusterid2",
			VClusterName:    "vcluster_name_2",
			RootClusterName: "root_cluster_name_1",
			IsDeleted:       0,
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			UserName:        "test_user",
			TenantId:        "tenant_id_1",
			VClusterId:      "vclusterid3",
			VClusterName:    "vcluster_name_3",
			RootClusterName: "root_cluster_name_2",
			IsDeleted:       0,
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			UserName:        "test_user",
			TenantId:        "tenant_id_1",
			VClusterId:      "vclusterid4",
			VClusterName:    "vcluster_name_4",
			RootClusterName: "root_cluster_name_1",
			IsDeleted:       1,
		},
	}

	// insert test datas
	diContainer.Invoke(func(db *gorm.DB) {
		db.AutoMigrate(&model.VCluster{})
		for i, vc := range modelVclusters {
			if err := db.Model(&model.VCluster{}).Create(&vc).Error; err != nil {
				t.Fatal(err)
			}

			modelVclusters[i] = vc
		}
	})

	diContainer.Invoke(func(dbDataSource VClusterDBDataSource) {
		testCases := []struct {
			description     string
			tenantType      string
			tenantId        string
			rootClusterName string
			isDeleted       int32
			permission      []string
			exceptRows      []model.VCluster
		}{
			{
				description:     "tc_1",
				tenantType:      "TENANT_TYPE_PLATFORM",
				tenantId:        "tenant_id_1",
				rootClusterName: "root_cluster_name_1",
				isDeleted:       0,
				permission:      []string{"*"},
				exceptRows:      []model.VCluster{modelVclusters[0], modelVclusters[1]},
			},
			{
				description:     "tc_2",
				tenantType:      "",
				tenantId:        "tenant_id_2",
				rootClusterName: "root_cluster_name_1",
				isDeleted:       0,
				permission:      []string{"*"},
				exceptRows:      []model.VCluster{modelVclusters[1]},
			},
			{
				description:     "tc_3",
				tenantType:      "TENANT_TYPE_PLATFORM",
				tenantId:        "tenant_id_1",
				rootClusterName: "root_cluster_name_1",
				isDeleted:       1,
				permission:      []string{"*"},
				exceptRows:      []model.VCluster{modelVclusters[3]},
			},
			{
				description:     "tc_4",
				tenantType:      "",
				tenantId:        "tenant_id_1",
				rootClusterName: "root_cluster_name_1",
				isDeleted:       0,
				permission:      []string{"vclusterid1"},
				exceptRows:      []model.VCluster{modelVclusters[0]},
			},
		}

		for _, tc := range testCases {
			ret, err := dbDataSource.FindVlusters(tc.tenantType, tc.tenantId, tc.rootClusterName, tc.isDeleted, tc.permission)
			assert.Nil(t, err)
			assert.Equal(t, tc.exceptRows, ret, tc.description)
		}
	})
}

func TestDBDataSourceFindVClusterGPUs(t *testing.T) {
	modelVGPUs := []model.VGpu{
		{
			ClusterID:    "cluster1",
			GpuType:      "H100",
			ResourceName: "NVIDIA-H100",
		},
		{
			ClusterID:    "cluster2",
			GpuType:      "H200",
			ResourceName: "NVIDIA-H200",
		},
		{
			ClusterID:    "cluster1",
			GpuType:      "H300",
			ResourceName: "NVIDIA-H300",
		},
	}

	// insert test datas
	diContainer.Invoke(func(db *gorm.DB) {
		db.AutoMigrate(&model.VGpu{})
		for i, gpu := range modelVGPUs {
			if err := db.Model(&model.VGpu{}).Create(&gpu).Error; err != nil {
				t.Fatal(err)
			}

			modelVGPUs[i] = gpu
		}
	})

	diContainer.Invoke(func(dbDataSource VClusterDBDataSource) {
		testCases := []struct {
			description string
			vclusterId  string
			exceptRows  []model.VGpu
		}{
			{
				description: "tc_1",
				vclusterId:  "cluster1",
				exceptRows:  []model.VGpu{modelVGPUs[0], modelVGPUs[2]},
			},
			{
				description: "tc_2",
				vclusterId:  "cluster2",
				exceptRows:  []model.VGpu{modelVGPUs[1]},
			},
		}

		for _, tc := range testCases {
			ret, err := dbDataSource.FindVClusterGPUs(tc.vclusterId)
			assert.Nil(t, err)
			assert.Equal(t, tc.exceptRows, ret, tc.description)
		}
	})
}

func TestDBDataSourceFindVClusterStorages(t *testing.T) {
	modelVStorages := []model.VStorage{
		{
			VClusterID:       "cluster1",
			VStorageType:     "hdd",
			VStorageCapacity: 20,
			IsDeleted:        0,
			Name:             "test_hdd",
		},
		{
			VClusterID:       "cluster1",
			VStorageType:     "ssd",
			VStorageCapacity: 100,
			IsDeleted:        0,
			Name:             "test_ssd",
		},
	}

	// insert test datas
	diContainer.Invoke(func(db *gorm.DB) {
		db.AutoMigrate(&model.VStorage{})
		for i, sc := range modelVStorages {
			if err := db.Model(&model.VStorage{}).Create(&sc).Error; err != nil {
				t.Fatal(err)
			}

			modelVStorages[i] = sc
		}
	})

	diContainer.Invoke(func(dbDataSource VClusterDBDataSource) {
		testCases := []struct {
			description string
			vclusterId  string
			exceptRows  []model.VStorage
		}{
			{
				description: "tc_1",
				vclusterId:  "cluster1",
				exceptRows:  modelVStorages,
			},
		}

		for _, tc := range testCases {
			ret, err := dbDataSource.FindVClusterStorages(tc.vclusterId)
			assert.Nil(t, err)
			assert.Equal(t, tc.exceptRows, ret, tc.description)
		}
	})
}
