package handler

import (
	"fmt"
	"os"

	"github.com/bitly/go-simplejson"
	gcpctx "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/context"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/response"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/framework"
)

var _ framework.ProfilesHandlerInterface = &ProfileHandler{}

type ProfileHandler struct {
	controller framework.Interface
	repo       pkg.VirtualServerRepo
}

func NewProfileHandler(controller framework.Interface, repoImpl pkg.VirtualServerRepo) *ProfileHandler {
	return &ProfileHandler{
		controller: controller,
		repo:       repoImpl,
	}
}

const (
	EMPTY_STORAGE            int64 = 0
	DEFAULT_MIN_STORAGE      int64 = 20
	DEFAULT_MAX_STORAGE      int64 = 100
	DEFAULT_DEFAULT_STORGAGE int64 = 50
)

// ListProductProfiles godoc
//
//	@Summary		List Product Profiles
//	@Description	List Product Profiles
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse{data=v1.Profiles}
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/product-profiles [get]
func (w *ProfileHandler) ListProductProfiles(c *gcpctx.GCPContext) {
	// username := c.GetUesrName()

	// // logger: log with username and sessionuuid
	// logger := c.Logger()

	// logger.Infof("username: %s", username)
	productsFilePath := w.controller.ComponentConfig().GetProductsFilePath()
	if productsFilePath != "" {
		imagesBytes, err := os.ReadFile(productsFilePath)
		if err != nil {
			response.Response(c, response.ErrJsonUnmarshal, nil)
			return
		}
		profiles, err := simplejson.NewJson(imagesBytes)
		if err != nil {
			// logger.Errorf("Failed to parse JSON data:", err)
			response.Response(c, response.ErrJsonUnmarshal, nil)
			return
		}
		response.Response(c, response.SuccessGCPResponse, &profiles)
		return
	}
	profilesStr := `{
        "product_categories": [
          {
            "code": "CP4",
            "value": "NVIDIA P4",
            "seq": "1"
          },
          {
            "code": "CA100",
            "value": "NVIDIA A100",
            "seq": "2"
          },
          {
            "code": "CH800",
            "value": "NVIDIA H800 SXM",
            "seq": "3"
          },
          {
            "code": "CA800-pcie",
            "value": "NVIDIA A800 PCIe",
            "seq": "4"
          },
          {
            "code": "CL40s",
            "value": "NVIDIA L40s",
            "seq": "5"
          }
        ],
        "products": [
          {
            "code": "DCX-CP4-1",
            "name": "产品1",
            "category": "CP4",
            "configs": [
              {
                "configKey": "gpu_type",
                "configValue": "Nvidia Tesla p4"
              },
              {
                "configKey": "gpu_quantity",
                "configValue": 1
              },
              {
                "configKey": "gpu_mem",
                "configValue": "8G"
              },
              {
                "configKey": "mem",
                "configValue": "32G"
              },
              {
                "configKey": "cpu",
                "configValue": 8
              }
            ]
          },
          {
            "code": "DCX-CA100-1",
            "name": "产品2",
            "category": "CA100",
            "configs": [
              {
                "configKey": "gpu_type",
                "configValue": "Nvidia A100"
              },
              {
                "configKey": "gpu_quantity",
                "configValue": 1
              },
              {
                "configKey": "gpu_mem",
                "configValue": "40G"
              },
              {
                "configKey": "mem",
                "configValue": "32G"
              },
              {
                "configKey": "cpu",
                "configValue": 8
              }
            ]
          },
          {
            "code": "DCX-CA800-1",
            "name": "产品3",
            "category": "CA800",
            "configs": [
              {
                "configKey": "gpu_type",
                "configValue": "Nvidia A800 SXM"
              },
              {
                "configKey": "gpu_quantity",
                "configValue": 1
              },
              {
                "configKey": "gpu_mem",
                "configValue": "80G"
              },
              {
                "configKey": "mem",
                "configValue": "32G"
              },
              {
                "configKey": "cpu",
                "configValue": 8
              }
            ]
          },
          {
            "code": "DCX-CH800-1",
            "name": "产品4",
            "category": "CH800",
            "configs": [
              {
                "configKey": "gpu_type",
                "configValue": "Nvidia H800 SXM"
              },
              {
                "configKey": "gpu_quantity",
                "configValue": 1
              },
              {
                "configKey": "gpu_mem",
                "configValue": "80G"
              },
              {
                "configKey": "mem",
                "configValue": "32G"
              },
              {
                "configKey": "cpu",
                "configValue": 8
              }
            ]
          },
          {
            "code": "DCX-CL40s-1",
            "name": "产品4",
            "category": "CL40s",
            "configs": [
              {
                "configKey": "gpu_type",
                "configValue": "Nvidia CL40s SXM"
              },
              {
                "configKey": "gpu_quantity",
                "configValue": 1
              },
              {
                "configKey": "gpu_mem",
                "configValue": "80G"
              },
              {
                "configKey": "mem",
                "configValue": "32G"
              },
              {
                "configKey": "cpu",
                "configValue": 8
              }
            ]
          }
        ]
      }`

	// 将字符串解析为simplejson对象
	profiles, err := simplejson.NewJson([]byte(profilesStr))
	if err != nil {
		response.Response(c, response.ErrJsonUnmarshal, nil)
		return
	}

	response.Response(c, response.SuccessGCPResponse, &profiles)
}

// ListImageProfiles godoc
//
//	@Summary		List Image profiles
//	@Description	List Image profiles
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse{data=v1.Profiles}
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/image-profiles [get]
func (w *ProfileHandler) ListImageProfiles(c *gcpctx.GCPContext) {
	// username := c.GetUesrName()

	// logger: log with username and sessionuuid
	// logger := c.Logger()
	// logger.Infof("username: %s", username)
	imagesFilePath := w.controller.ComponentConfig().GetImagesFilePath()
	if imagesFilePath != "" {
		imagesBytes, err := os.ReadFile(imagesFilePath)
		if err != nil {
			response.Response(c, response.ErrJsonUnmarshal, nil)
			return
		}
		profiles, err := simplejson.NewJson(imagesBytes)
		if err != nil {
			// logger.Errorf("Failed to parse JSON data:", err)
			response.Response(c, response.ErrJsonUnmarshal, nil)
			return
		}
		response.Response(c, response.SuccessGCPResponse, &profiles)
		return
	}

	profilesStr := `{
        "vmimage_categories": [
          {
            "code": "ubuntu",
            "value": "Ubuntu",
            "seq": "1"
          },
          {
            "code": "centos",
            "value": "CentOs",
            "seq": "2"
          },
          {
            "code": "redHat",
            "value": "Red Hat",
            "seq": "3"
          }
        ],
        "vmimages": [
          {
            "imageName": "aidc-ubuntu-testimage",
            "category": "ubuntu",
            "config": {
              "gpudriver": "",
              "cuda": ""
            }
          },
          {
            "imageName": "aidc-ubuntu-20.04-server-cloudimg-amd64-wfdsgsfb",
            "category": "ubuntu",
            "config": {
              "gpudriver": "470.239.06",
              "cuda": "11.4"
            }
          },
          {
            "imageName": "aidc-ubuntu-22.04-server-cloudimg-amd64-undmkhnd",
            "category": "ubuntu",
            "config": {
              "gpudriver": "545.29.06",
              "cuda": "12.3"
            }
          },
          {
            "imageName": "aidc-centos-7-generic-cloud-crgsgsdh",
            "category": "centos",
            "config": {
              "gpudriver": "515.65.01",
              "cuda": "11.7"
            }
          },
          {
            "imageName": "aidc-centos-8-generic-cloud-cgkdmtib",
            "category": "centos",
            "config": {
              "gpudriver": "510.47.03",
              "cuda": "11.6"
            }
          }
        ]
      }`

	// 将字符串解析为simplejson对象
	profiles, err := simplejson.NewJson([]byte(profilesStr))
	if err != nil {
		// logger.Errorf("Failed to parse JSON data:", err)
		response.Response(c, response.ErrJsonUnmarshal, nil)
		return
	}
	response.Response(c, response.SuccessGCPResponse, &profiles)
}

// ListStorageProfiles godoc
//
//	@Summary		List Storage profiles
//	@Description	List Storage profiles
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse{data=v1.Profiles}
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/storage-profiles [get]
func (w *ProfileHandler) ListStorageProfiles(c *gcpctx.GCPContext) {
	// username := c.GetUesrName()

	// logger: log with username and sessionuuid
	// logger := c.Logger()

	minDisk := DEFAULT_MIN_STORAGE
	maxDisk := DEFAULT_MAX_STORAGE
	defaultDisk := DEFAULT_DEFAULT_STORGAGE
	if w.controller.ComponentConfig().GetStorageMin() != EMPTY_STORAGE {
		minDisk = w.controller.ComponentConfig().GetStorageMin()
	}

	if w.controller.ComponentConfig().GetStorageMax() != EMPTY_STORAGE {
		maxDisk = w.controller.ComponentConfig().GetStorageMax()
	}

	if w.controller.ComponentConfig().GetStorageDefault() != EMPTY_STORAGE {
		defaultDisk = w.controller.ComponentConfig().GetStorageDefault()
	}

	minDiskStr := fmt.Sprintf("%d", minDisk)
	maxDiskStr := fmt.Sprintf("%d", maxDisk)
	defaultDiskStr := fmt.Sprintf("%d", defaultDisk)

	// logger.Infof("username: %s", username)
	profilesStr := `{
	      "storage": {
	          "storage_system": {
	              "min": ` + minDiskStr + `,
	              "max": ` + maxDiskStr + `,
	              "default": ` + defaultDiskStr + `
	          },
	          "storage_data": {
	              "min": ` + minDiskStr + `,
	              "max": ` + maxDiskStr + `,
	              "default": ` + defaultDiskStr + `
	          }
	      }
	    }`

	// 将字符串解析为simplejson对象
	profiles, err := simplejson.NewJson([]byte(profilesStr))
	if err != nil {
		// logger.Errorf("Failed to parse JSON data:", err)
		response.Response(c, response.ErrJsonUnmarshal, nil)
		return
	}

	response.Response(c, response.SuccessGCPResponse, &profiles)
}

// Available godoc
//
//	@Summary		Get Products available amount
//	@Description	get products available amount
//	@Tags			Profiles
//	@Accept			json
//	@Produce		json
//	@Param			X-Access-Token	header		string	true	"用户 JWT token"
//	@Success		200				{object}	response.BaseResponse{data=v1.Profiles}
//	@Failure		400				{object}	response.BaseResponse
//	@Failure		401				{object}	response.BaseResponse
//	@Failure		404				{object}	response.BaseResponse
//	@Failure		500				{object}	response.BaseResponse
//	@Router			/api/kvm/v1/product/available [get]
func (w *ProfileHandler) Available(c *gcpctx.GCPContext) {
	userName := c.GetUesrName()
	logger := c.Logger()

	logger.Infof("username: %s", userName)
	// 获取config.yaml里的products（后续产品的配置，以及产品和GPU节点label的对应关系建议放到configmap里）
	products := w.controller.ComponentConfig().GetProducts()
	available := make(map[string]int)

	// 查询全量虚拟机列表
	vms, err := w.repo.List("")
	if err != nil {
		logger.Infof("list vm from repo err, %v", err)
	}
	// 查询各产品的使用数
	usedProducts := map[string]int{}
	for _, vm := range vms {
		// 已删除的虚拟机产品会被释放
		usedProducts[vm.Product] = usedProducts[vm.Product] + 1
	}
	for key, product := range products {
		product := product.(map[string]interface{})
		amount := product["amount"].(int)
		quantity := usedProducts[key]
		// 从数据库查询该product的使用数量返回，测试阶段mock quantity=0
		available[key] = amount - quantity
	}

	response.Response(c, response.SuccessGCPResponse, available)

}
