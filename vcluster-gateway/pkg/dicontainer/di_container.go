package dicontainer

import (
	"go.uber.org/dig"
	"gorm.io/gorm"

	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/datasource"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/repository"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/usecase"
)

var globalDIContainer *DIContainer

var GlobalDIContainer = globalDIContainer

type DIContainer struct {
	controller framework.Interface
	container  *dig.Container
}

func NewDIContainer(controller framework.Interface) *DIContainer {
	di := &DIContainer{
		controller: controller,
		container:  dig.New(),
	}

	di.build()

	return di
}

func (di *DIContainer) build() {
	di.container.Provide(func() *gorm.DB {
		return di.controller.ComponentConfig().AllCluster.DB
	})

	di.container.Provide(datasource.NewDBDataSource)

	di.container.Provide(datasource.NewKubernetesDataSource)

	di.container.Provide(repository.NewVClusterRepository)

	di.container.Provide(usecase.NewVClusterUseCase)
}

func (di *DIContainer) Invoke(function interface{}, opts ...dig.InvokeOption) error {
	return di.container.Invoke(function, opts...)
}
