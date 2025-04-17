package pkg

import (
	"fmt"

	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/repo"
)

// // New -.
func NewVirtualServerRepo(datasource string) (VirtualServerRepo, error) {
	fmt.Println("newVirtualServerRepo")
	return repo.NewVirtualServerMysqlImpl(datasource)
}

type (
	// Translation -.
	VirtualServerRepo interface {
		ListAll(userName string) ([]repo.VirtualServer, error)
		List(userName string) ([]repo.VirtualServer, error)
		ListDeletedVms(userName string) ([]repo.VirtualServer, error)
		Store(vm repo.VirtualServer) error
		Update(vm *repo.VirtualServer) error
		GetVmById(userName string, id string) (repo.VirtualServer, error)
	}
)
