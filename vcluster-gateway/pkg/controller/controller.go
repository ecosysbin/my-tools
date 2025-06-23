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
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	v1 "vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	"vcluster-gateway/pkg/controller/framework"
	"vcluster-gateway/pkg/controller/handler"
	"vcluster-gateway/pkg/controller/services"
	"vcluster-gateway/pkg/controller/vcluster"
	"vcluster-gateway/pkg/dicontainer"
)

var _ framework.Interface = &Controller{}

type Controller struct {
	cfg *controllerConfiguration

	vClusterController *vcluster.VClusterController
	diContainer        *dicontainer.DIContainer

	// servers
	vclusterServer *services.VClusterServer

	// httpclient
	httpclient *resty.Client
	// handler is the root node of the routers.
	Router *gin.Engine

	vClusterHandler *handler.VClusterHandler
}

func New(opts ...Option) *Controller {
	c := controllerConfiguration{}
	for _, opt := range opts {
		opt(&c)
	}

	controller := Controller{
		cfg: &c,
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
		// new gin restHandler
		Router: gin.Default(),
	}

	controller.diContainer = dicontainer.NewDIContainer(&controller)

	controller.vClusterController = vcluster.NewVClusterController(&controller)

	controller.vclusterServer = services.NewVClusterServer(&controller)

	// new
	controller.vClusterHandler = handler.NewVClusterHandler(&controller)
	controller.registerHandlers()

	return &controller
}

// ****** Controller Funcs ******
func (c *Controller) HttpClient() *resty.Client {
	return c.httpclient
}

func (c *Controller) ComponentConfig() *v1.VclusterGatewayConfiguration {
	return c.cfg.componentConfig
}

func (c *Controller) VClusterController() framework.VClusterControllerInterface {
	return c.vClusterController
}

func (c *Controller) DIContainer() framework.DIContainerInterface {
	return c.diContainer
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
	addr := fmt.Sprintf(":%s", c.ComponentConfig().GetServerPort())
	log.Infof("Serving gRPC-Gateway on http://0.0.0.0" + addr)

	go func() {
		if err := c.Router.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-stopCh
	log.Info("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling

	log.Info("Server exiting")
}
