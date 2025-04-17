package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/endverse/go-kit/signals"

	. "github.com/smartystreets/goconvey/convey"

	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg"
	configv1 "gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/config/kubevirt_gateway/v1"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/kube"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/response"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/handler"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/domain"
	"k8s.io/client-go/kubernetes"
	"kubevirt.io/client-go/kubecli"
)

func TestGetStorageProfiles(t *testing.T) {
	Convey("Given mock config storage profile", t, func() {
		// given mock data
		patches := MockDb()
		defer patches.Reset()
		patches1 := MockKubernetes()
		defer patches1.Reset()
		patches2 := MockKubevirt()
		defer patches2.Reset()
		patches3 := MockVirtualServerHandler()
		defer patches3.Reset()
		var con *controller.Controller
		patches5 := ApplyMethod(con, "SetGlobalMiddleware", func(_ *controller.Controller) {
		})
		defer patches5.Reset()
		appConfig := MockConfig()
		appConfig.ComponentConfig.Storage = configv1.Storage{
			Min:     1,
			Max:     100,
			Default: 10,
		}
		stopCh := signals.SetupSignalHandler()
		contro := controller.New(appConfig, stopCh)

		Convey("When Get storage-profiles", func() {
			// when
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/kvm/v1/storage-profiles", nil)
			contro.Handler.ServeHTTP(w, req)
			Convey("then status shold be ok, storage data shold equal to mock data", func() {
				// then
				So(w.Code, ShouldEqual, http.StatusOK)
				body, err := io.ReadAll(w.Body)
				if err != nil {
					t.Errorf("read body err, %v", err)
				}
				var resp = response.BaseResponse{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Errorf("unmarhsal body err, %v", err)
				}
				So(resp.Status, ShouldEqual, response.SuccessGCPResponse.Code())
				var dataMap = resp.Data.(map[string]interface{})
				var system_storage = dataMap["storage"].(map[string]interface{})["storage_system"]
				var data_storage = dataMap["storage"].(map[string]interface{})["storage_data"]
				system_min := system_storage.(map[string]interface{})["min"].(float64)
				system_max := system_storage.(map[string]interface{})["max"].(float64)
				system_default := system_storage.(map[string]interface{})["default"].(float64)
				data_min := data_storage.(map[string]interface{})["min"].(float64)
				data_max := data_storage.(map[string]interface{})["max"].(float64)
				data_default := data_storage.(map[string]interface{})["default"].(float64)
				So(system_min, ShouldEqual, 1)
				So(system_max, ShouldEqual, 100)
				So(system_default, ShouldEqual, 10)
				So(data_min, ShouldEqual, 1)
				So(data_max, ShouldEqual, 100)
				So(data_default, ShouldEqual, 10)
			})
		})
	})
}

func TestGetStorageProfilesWithDefault(t *testing.T) {
	Convey("get storage profile with default", t, func() {
		// given mock data
		patches := MockDb()
		defer patches.Reset()
		patches1 := MockKubernetes()
		defer patches1.Reset()
		patches2 := MockKubevirt()
		defer patches2.Reset()
		patches3 := MockVirtualServerHandler()
		defer patches3.Reset()
		var con *controller.Controller
		patches5 := ApplyMethod(con, "SetGlobalMiddleware", func(_ *controller.Controller) {
		})
		defer patches5.Reset()
		appConfig := MockConfig()
		// appConfig.ComponentConfig.Storage = configv1.Storage{
		// 	Min:     1,
		// 	Max:     100,
		// 	Default: 10,
		// }
		stopCh := signals.SetupSignalHandler()
		contro := controller.New(appConfig, stopCh)

		Convey("When Get storage-profiles", func() {
			// when
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/kvm/v1/storage-profiles", nil)
			contro.Handler.ServeHTTP(w, req)
			Convey("then status shold be ok, storage data shold equal default storage data", func() {
				// then
				So(w.Code, ShouldEqual, http.StatusOK)
				body, err := io.ReadAll(w.Body)
				if err != nil {
					t.Errorf("read body err, %v", err)
				}
				var resp = response.BaseResponse{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Errorf("unmarhsal body err, %v", err)
				}
				So(resp.Status, ShouldEqual, response.SuccessGCPResponse.Code())
				var dataMap = resp.Data.(map[string]interface{})
				var system_storage = dataMap["storage"].(map[string]interface{})["storage_system"]
				var data_storage = dataMap["storage"].(map[string]interface{})["storage_data"]
				system_min := system_storage.(map[string]interface{})["min"].(float64)
				system_max := system_storage.(map[string]interface{})["max"].(float64)
				system_default := system_storage.(map[string]interface{})["default"].(float64)
				data_min := data_storage.(map[string]interface{})["min"].(float64)
				data_max := data_storage.(map[string]interface{})["max"].(float64)
				data_default := data_storage.(map[string]interface{})["default"].(float64)
				So(system_min, ShouldEqual, handler.DEFAULT_MIN_STORAGE)
				So(system_max, ShouldEqual, handler.DEFAULT_MAX_STORAGE)
				So(system_default, ShouldEqual, handler.DEFAULT_DEFAULT_STORGAGE)
				So(data_min, ShouldEqual, handler.DEFAULT_MIN_STORAGE)
				So(data_max, ShouldEqual, handler.DEFAULT_MAX_STORAGE)
				So(data_default, ShouldEqual, handler.DEFAULT_DEFAULT_STORGAGE)
			})
		})
	})
}

func TestGetImagesProfileDefault(t *testing.T) {
	Convey("get images profile", t, func() {
		// given mock data
		patches := MockDb()
		defer patches.Reset()
		patches1 := MockKubernetes()
		defer patches1.Reset()
		patches2 := MockKubevirt()
		defer patches2.Reset()
		patches3 := MockVirtualServerHandler()
		defer patches3.Reset()
		var con *controller.Controller
		patches5 := ApplyMethod(con, "SetGlobalMiddleware", func(_ *controller.Controller) {
		})
		defer patches5.Reset()
		appConfig := MockConfig()
		stopCh := signals.SetupSignalHandler()
		contro := controller.New(appConfig, stopCh)

		Convey("When Get images profiles with mock data", func() {
			// when
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/kvm/v1/image-profiles", nil)
			contro.Handler.ServeHTTP(w, req)
			Convey("then status shold be ok, storage data shold equal default storage data", func() {
				// then
				So(w.Code, ShouldEqual, http.StatusOK)
				body, err := io.ReadAll(w.Body)
				if err != nil {
					t.Errorf("read body err, %v", err)
				}
				var resp = response.BaseResponse{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Errorf("unmarhsal body err, %v", err)
				}
				So(resp.Status, ShouldEqual, response.SuccessGCPResponse.Code())
				var dataMap = resp.Data.(map[string]interface{})
				var vmimages = dataMap["vmimages"].([]interface{})
				So(len(vmimages), ShouldEqual, 5)
			})
		})
	})
}

func TestGetImagesProfileWithMockImagesPath(t *testing.T) {
	Convey("get images profile", t, func() {
		// given mock data
		patches := MockDb()
		defer patches.Reset()
		patches1 := MockKubernetes()
		defer patches1.Reset()
		patches2 := MockKubevirt()
		defer patches2.Reset()
		patches3 := MockVirtualServerHandler()
		defer patches3.Reset()
		var con *controller.Controller
		patches5 := ApplyMethod(con, "SetGlobalMiddleware", func(_ *controller.Controller) {
		})
		defer patches5.Reset()
		appConfig := MockConfig()
		// mock imagesFilePath
		appConfig.ComponentConfig.Images.ImageFilePath = "./images.json"
		stopCh := signals.SetupSignalHandler()
		contro := controller.New(appConfig, stopCh)

		Convey("When Get images profiles with mock data", func() {
			// when
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/kvm/v1/image-profiles", nil)
			contro.Handler.ServeHTTP(w, req)
			Convey("then status shold be ok, storage data shold equal default storage data", func() {
				// then
				So(w.Code, ShouldEqual, http.StatusOK)
				body, err := io.ReadAll(w.Body)
				if err != nil {
					t.Errorf("read body err, %v", err)
				}
				var resp = response.BaseResponse{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Errorf("unmarhsal body err, %v", err)
				}
				So(resp.Status, ShouldEqual, response.SuccessGCPResponse.Code())
				var dataMap = resp.Data.(map[string]interface{})
				var vmimages = dataMap["vmimages"].([]interface{})
				So(len(vmimages), ShouldEqual, 4)
			})
		})
	})
}

func TestGetProductProfileDefault(t *testing.T) {
	Convey("get images profile", t, func() {
		// given mock data
		patches := MockDb()
		defer patches.Reset()
		patches1 := MockKubernetes()
		defer patches1.Reset()
		patches2 := MockKubevirt()
		defer patches2.Reset()
		patches3 := MockVirtualServerHandler()
		defer patches3.Reset()
		var con *controller.Controller
		patches5 := ApplyMethod(con, "SetGlobalMiddleware", func(_ *controller.Controller) {
		})
		defer patches5.Reset()
		appConfig := MockConfig()
		stopCh := signals.SetupSignalHandler()
		contro := controller.New(appConfig, stopCh)

		Convey("When Get Product profiles with default data", func() {
			// when
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/kvm/v1/product-profiles", nil)
			contro.Handler.ServeHTTP(w, req)
			Convey("then status shold be ok, storage data shold equal default storage data", func() {
				// then
				So(w.Code, ShouldEqual, http.StatusOK)
				body, err := io.ReadAll(w.Body)
				if err != nil {
					t.Errorf("read body err, %v", err)
				}
				var resp = response.BaseResponse{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Errorf("unmarhsal body err, %v", err)
				}
				So(resp.Status, ShouldEqual, response.SuccessGCPResponse.Code())
				var dataMap = resp.Data.(map[string]interface{})
				var products = dataMap["products"].([]interface{})
				So(len(products), ShouldEqual, 4)
			})
		})
	})
}

func TestGetProductProfileWithMockProductFilePath(t *testing.T) {
	Convey("get images profile", t, func() {
		// given mock data
		patches := MockDb()
		defer patches.Reset()
		patches1 := MockKubernetes()
		defer patches1.Reset()
		patches2 := MockKubevirt()
		defer patches2.Reset()
		patches3 := MockVirtualServerHandler()
		defer patches3.Reset()
		var con *controller.Controller
		patches5 := ApplyMethod(con, "SetGlobalMiddleware", func(_ *controller.Controller) {
		})
		defer patches5.Reset()
		appConfig := MockConfig()
		appConfig.ComponentConfig.Products.ProductFilePath = "./products.json"
		stopCh := signals.SetupSignalHandler()
		contro := controller.New(appConfig, stopCh)

		Convey("When Get Product profiles with default data", func() {
			// when
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/kvm/v1/product-profiles", nil)
			contro.Handler.ServeHTTP(w, req)
			Convey("then status shold be ok, storage data shold equal default storage data", func() {
				// then
				So(w.Code, ShouldEqual, http.StatusOK)
				body, err := io.ReadAll(w.Body)
				if err != nil {
					t.Errorf("read body err, %v", err)
				}
				var resp = response.BaseResponse{}
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Errorf("unmarhsal body err, %v", err)
				}
				So(resp.Status, ShouldEqual, response.SuccessGCPResponse.Code())
				var dataMap = resp.Data.(map[string]interface{})
				var products = dataMap["products"].([]interface{})
				So(len(products), ShouldEqual, 5)
			})
		})
	})
}

func MockDb() *Patches {
	patches := ApplyFunc(pkg.NewVirtualServerRepo, func(_ string) (pkg.VirtualServerRepo, error) {
		return &FakeDb{}, nil
	})
	return patches
}

func MockKubernetes() *Patches {
	patches := ApplyFunc(kube.CreateClients, func(_ *kube.KubeConfiguration) (kubernetes.Interface, error) {
		return nil, nil
	})
	return patches
}

func MockKubevirt() *Patches {
	patches := ApplyFunc(kubecli.GetKubevirtClientFromFlags, func(_ string, _ string) (kubecli.KubevirtClient, error) {
		return nil, nil
	})
	return patches
}

func MockVirtualServerHandler() *Patches {
	patches := ApplyFunc(handler.NewVirtualServerHandler, func(_ framework.Interface, _ domain.VirtualServerManager) (*handler.VirtualServerHandler, error) {
		return nil, nil
	})
	return patches
}

// 1. 只是简单对某个函数的依赖，可以简单对函数mock
// 2. 对某个客户端的依赖，除非想要对客户端的每个方法mock，不然建议：1. 对客户端接口依赖 2. 创建fake客户端，使用fake客户端对真正的客户端进行mock
