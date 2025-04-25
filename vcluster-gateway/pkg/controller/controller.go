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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	v1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
	vclusterv1 "gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/apis/grpc/gen/datacanvas/gcp/osm/vcluster_1.1/v1/vclusterv1connect"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/controller/services"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/controller/vcluster"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/datasource"
	"gitlab.datacanvas.com/aidc/vcluster-gateway/pkg/dicontainer"
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
	handler *gin.Engine
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
		handler: gin.New(),
	}

	controller.diContainer = dicontainer.NewDIContainer(&controller)

	controller.vClusterController = vcluster.NewVClusterController(&controller)

	controller.vclusterServer = services.NewVClusterServer(&controller)

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
	var kubernetesDataSource datasource.VClusterKubernetesDataSource

	err := c.DIContainer().Invoke(func(ds datasource.VClusterKubernetesDataSource) {
		kubernetesDataSource = ds
	})
	if err != nil {
		panic(err)
	}

	// go c.VClusterController().SyncWorkflows()
	// go c.VClusterController().ReProcessVCluster()

	defer kubernetesDataSource.Shutdown()

	c.Serve(stopCh)
	return nil
}

func (c *Controller) Serve(stopCh <-chan struct{}) {
	c.startServerWithGracefulShutdown(stopCh)
}

func (c *Controller) startServerWithGracefulShutdown(stopCh <-chan struct{}) {
	mux := http.NewServeMux()
	svcPath, handler := vclusterv1connect.NewVClusterGatewayServiceHandler(c.vclusterServer)
	mux.Handle(svcPath, handler)

	srv := &http.Server{
		Addr:    ":" + "8088",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof(" ****** GCP VclusterGateway Gateway Server is Start at [:%s] ****** ", "8088")

	conn, err := grpc.Dial(
		"localhost:8088",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			md := metadata.MD{}

			log.Infof("http method: %v", r.Method)
			log.Infof("http Path: %v", r.URL.Path)

			return md
		}),
	)

	err = vclusterv1.RegisterVClusterGatewayServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", c.ComponentConfig().GetServerPort()),
		Handler: gwmux,
	}

	log.Infof("Serving gRPC-Gateway on http://0.0.0.0" + fmt.Sprintf(":%s", c.ComponentConfig().GetServerPort()))

	go func() {
		if err := gwServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-stopCh
	log.Info("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := gwServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Info("Server exiting")
}
