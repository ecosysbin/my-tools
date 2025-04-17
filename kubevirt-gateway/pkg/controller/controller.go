//
// Copyright 2023 The Zetyun.GCP Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package controller

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"
	appconfig "gitlab.datacanvas.com/aidc/kubevirt-gateway/cmd/app/config"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/kube"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/casdoor"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/handler"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/middleware"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/domain"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"

	"kubevirt.io/client-go/kubecli"
)

var _ framework.Interface = &Controller{}

type Controller struct {
	cfg *appconfig.Config

	// controllers
	*casdoor.CasdoorController
	*middleware.MiddlewareController

	// handlers
	*handler.ProfileHandler
	*handler.VirtualServerHandler

	// httpclient
	httpclient *resty.Client
	// handler is the root node of the routers.
	Handler *gin.Engine
}

func New(c *appconfig.Config, stopCh <-chan struct{}) *Controller {
	controller := Controller{
		cfg: c,
		httpclient: resty.New().
			SetTransport(&http.Transport{
				DisableKeepAlives: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}).
			SetHeader("Content-Type", "application/json").
			SetHeader("Connection", "close").
			SetRetryCount(3).
			SetTimeout(3 * time.Second),
		Handler: gin.New(),
	}
	KubeApiserver := c.ComponentConfig.GetHttpKubeApiserver()
	KubeConfig := c.ComponentConfig.GetHttpKubeConfig() // 配置为空时，走服务账号
	// 初始化kubevirt客户端
	kubvirtClient, err := kubecli.GetKubevirtClientFromFlags(KubeApiserver, KubeConfig)
	if err != nil {
		log.Infof("get kubevirt client err. %v", err)
		return nil
	}
	// 初始化kubernetes客户端
	kubernetes, err := kube.CreateClients(&c.KubeConfig)
	if err != nil {
		log.Infof("get kubernetes client err. %v", err)
		return nil
	}
	// 初始化repo客户端
	repoImpl, err := pkg.NewVirtualServerRepo(c.ComponentConfig.GetMysqlUrl())
	if err != nil {
		log.Infof("new virtualserver repo err. %v", err)
		return nil
	}
	// Init controllers
	controller.CasdoorController = casdoor.New(c.ComponentConfig)
	controller.MiddlewareController = middleware.New(&controller)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	virtualserverManager, err := domain.NewVirtualServerManager(ctx, repoImpl, kubvirtClient, kubernetes)
	if err != nil {
		log.Infof("new virtualServerManager err, %v", err)
		return nil
	}
	// 初始化profileHandler
	controller.ProfileHandler = handler.NewProfileHandler(&controller, repoImpl)
	// 初始化virtualserverHandler
	controller.VirtualServerHandler, err = handler.NewVirtualServerHandler(&controller, virtualserverManager)
	if err != nil {
		log.Infof("new virtualServerHandler err, %v", err)
		return nil
	}
	controller.registerHandlers()

	return &controller
}

// ****** Controller Funcs ******
func (c *Controller) HttpClient() *resty.Client {
	return c.httpclient
}

func (c *Controller) ComponentConfig() framework.ComponentConfigInterface {
	return c.cfg.ComponentConfig
}

// ****** Start Controller ******
func (c *Controller) Run(stopCh <-chan struct{}) error {
	c.Serve(stopCh)
	return nil
}

func (c *Controller) Serve(stopCh <-chan struct{}) {
	c.startServerWithGracefulShutdown(stopCh)
}

func (c *Controller) startServerWithGracefulShutdown(stopCh <-chan struct{}) {
	srv := &http.Server{
		Addr:    ":" + c.ComponentConfig().GetServerPort(),
		Handler: c.Handler,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Infof("start listen server port %s", c.ComponentConfig().GetServerPort())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof(" ****** GCP KubevirtGateway Gateway Server is Start at [:%s] ****** ", c.ComponentConfig().GetServerPort())

	<-stopCh
	log.Info("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("Server exiting")
}
