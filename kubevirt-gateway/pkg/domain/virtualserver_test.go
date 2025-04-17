package domain

import (
	"fmt"
	"testing"

	v1 "gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/kubevirt_gateway/v1"
)

func TestStartVirtualserverWorkFlow(t *testing.T) {
	virtualServerManager := VirtualServerManager{}
	virtualserver := v1.VirtualServer{
		Name: "testVirtualServer",
	}

	flow1 := VirtulServerFlow{
		Work: testFun1,
	}

	flow2 := VirtulServerFlow{
		Work:     testFun2,
		RollBack: testFun2Rollback,
	}

	flow3 := VirtulServerFlow{
		Work:     testFun3Err,
		RollBack: testFun3Rollback,
	}

	virtulServerCreateWorkFlow := VirtulServerWorkFlow{
		Metadata: virtualserver,
		Works:    []VirtulServerFlow{flow1, flow2, flow3},
	}
	virtualServerManager.StartVirtualserverWorkFlow(virtulServerCreateWorkFlow)
}

var testFun1 = func(virtualServer v1.VirtualServer) error {
	fmt.Printf("=====================> test Fun1  \n")
	return nil
}

var testFun2 = func(virtualServer v1.VirtualServer) error {
	fmt.Printf("=====================> test Fun2 \n")
	return nil
}

var testFun2Rollback = func(virtualServer v1.VirtualServer) error {
	fmt.Printf("=====================> test Fun2 rollBack \n")
	return nil
}

var testFun3Err = func(virtualServer v1.VirtualServer) error {
	fmt.Printf("=====================> test Fun3 err \n")
	return fmt.Errorf("test fun2 err")
}

var testFun3Rollback = func(virtualServer v1.VirtualServer) error {
	fmt.Printf("=====================> test Fun3 rollBack \n")
	return nil
}
