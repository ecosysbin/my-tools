package tests

import "gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/repo"

type FakeDb struct{}

func (fakedb *FakeDb) ListAll(userName string) ([]repo.VirtualServer, error) {
	return nil, nil
}

func (fakedb *FakeDb) List(userName string) ([]repo.VirtualServer, error) {
	return nil, nil
}

func (fakedb *FakeDb) ListDeletedVms(userName string) ([]repo.VirtualServer, error) {
	return nil, nil
}

func (fakedb *FakeDb) Store(vm repo.VirtualServer) error {
	return nil
}

func (fakedb *FakeDb) Update(vm *repo.VirtualServer) error {
	return nil
}

func (fakedb *FakeDb) GetVmById(userName string, id string) (repo.VirtualServer, error) {
	return repo.VirtualServer{}, nil
}
