package v1

import (
	"fmt"

	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/repo"
)

func TransformOSMVirtualServerReqToOSMVirtualServer(osmvirtualServerReq OSMVirtualServerReq) VirtualServer {
	return VirtualServer{
		Name:        osmvirtualServerReq.Name,
		Desc:        osmvirtualServerReq.Desc,
		Image:       osmvirtualServerReq.Image,
		Storage:     osmvirtualServerReq.Storage,
		CloudInit:   osmvirtualServerReq.CloudInit,
		ProductCode: osmvirtualServerReq.ProductInfo.ProductCode,
	}
}

func TransformBSMVirtualServerReqToVirtualServer(bsmvirtualServerReq BSMVirtualServerReq) VirtualServer {
	return VirtualServer{
		Name:        bsmvirtualServerReq.Name,
		Desc:        bsmvirtualServerReq.Desc,
		Image:       bsmvirtualServerReq.Image,
		Storage:     bsmvirtualServerReq.Storage,
		CloudInit:   bsmvirtualServerReq.CloudInit,
		ProductCode: bsmvirtualServerReq.ProductInfo.ProductCode,
		// bsm场景，账单是和实例绑定，为了所有产品实例id格式统一且全局唯一，则由账单服务统一生成实例id
		State: State{
			InstanceId: bsmvirtualServerReq.InstanceId,
		},
	}
}

// 当前在只有OSM的场景，产品管理在子产品，产品详情和数量直接配置在启动配置文件里。类型强转进行模型转换，配置错误则启动报错。后续配置方式变更，修改这边即可。
func TransformLocalProductsToOSMProduct(products map[string]interface{}) (OSMProduct, error) {
	osmProduct := map[string]LocalProduct{}
	for productCode, v := range products {
		product, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("transform product %s failed", productCode)
		}
		prodctNum, ok := product["amount"].(int)
		if !ok {
			return nil, fmt.Errorf("transform product %s amount failed", productCode)
		}
		gpuNum, ok := product["gpu_num"].(int)
		if !ok {
			return nil, fmt.Errorf("transform product %s gpu_num failed", productCode)
		}
		gpuk8sresource, ok := product["gpu_k8s_resource"].(string)
		if !ok {
			return nil, fmt.Errorf("transform product %s gpu_k8s_resource failed", productCode)
		}
		cpuRequest, ok := product["cpu"].(string)
		if !ok {
			return nil, fmt.Errorf("transform product %s cpu failed", productCode)
		}
		MemoryRequest, ok := product["mem"].(string)
		if !ok {
			return nil, fmt.Errorf("transform product %s mem failed", productCode)
		}
		osmProduct[productCode] = LocalProduct{
			Num: prodctNum,
			ProductResource: &ProductResource{
				GpuK8sResource: gpuk8sresource,
				GpuNum:         gpuNum,
				Mem:            MemoryRequest,
				Cpu:            cpuRequest,
			},
		}
	}
	return osmProduct, nil
}

func TransformVirtualServerToRepo(vm VirtualServer) repo.VirtualServer {
	virtualServerRepo := repo.VirtualServer{
		Id:          vm.State.InstanceId,
		Name:        vm.Name,
		Image:       vm.Image,
		Desc:        vm.Desc,
		Product:     vm.ProductCode,
		CreateUser:  vm.CreateUser,
		SshPort:     vm.State.SshPort,
		Status:      vm.State.Status,
		Deleted:     vm.Deleted,
		CreateTime:  vm.CreateTime,
		StartedTime: vm.StartedTime,
		DeleteTime:  vm.DeleteTime,
	}
	if vm.State.FailedStatus != nil {
		virtualServerRepo.Reason = vm.State.FailedStatus.Reason
		virtualServerRepo.Message = vm.State.FailedStatus.Message
	}
	return virtualServerRepo
}

func TransformRepoToVirtualServer(repo repo.VirtualServer) VirtualServer {
	virtualServer := VirtualServer{
		Name:        repo.Name,
		Desc:        repo.Desc,
		Image:       repo.Image,
		ProductCode: repo.Product,
		CreateUser:  repo.CreateUser,
		State: State{
			InstanceId: repo.Id,
			Status:     repo.Status,
			Vnc:        vnc(repo.Id),
			SshPort:    repo.SshPort,
		},
		CreateTime:  repo.CreateTime,
		StartedTime: repo.StartedTime,
		DeleteTime:  repo.DeleteTime,
		Deleted:     repo.Deleted,
	}
	if repo.Reason != "" {
		virtualServer.State.FailedStatus = &FailedStatus{
			Reason:  repo.Reason,
			Message: repo.Message,
		}
	}
	return virtualServer
}

func TransformRepoToVirtualServers(vmRepos []repo.VirtualServer) []VirtualServer {
	vertualServers := []VirtualServer{}
	for _, repo := range vmRepos {
		vertualServers = append(vertualServers, TransformRepoToVirtualServer(repo))
	}
	return vertualServers
}
