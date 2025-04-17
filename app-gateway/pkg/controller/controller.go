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
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/authz"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/casdoor"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/connectrpc/interceptors"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/connectrpc/interceptors/authentication"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller/app"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller/framework"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/repo"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/services"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	appv1 "gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1"
	"gitlab.datacanvas.com/aidc/app-gateway/generater/apis/grpc/gen/datacanvas/gcp/osm/app/v1/appv1connect"

	v1 "gitlab.datacanvas.com/aidc/app-gateway/pkg/apis/config/app_gateway/v1"
)

var _ framework.Interface = &Controller{}

type Controller struct {
	cfg *controllerConfiguration

	// controllers
	*casdoor.CasdoorController
	*authz.GcpAuthorization
	// *middleware.MiddlewareController
	appController *app.AppController
	interceptors.Interceptors

	// servers
	appServer *services.AppServer
	appRepo   framework.AppRepoInterface
	// httpclient
	httpclient *resty.Client
	// handler is the root node of the routers.
	handler *gin.Engine
}

func New(opts ...Option) (*Controller, error) {
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

	// Init controllers
	casdoorController, err := casdoor.New(&casdoor.CasdoorControllerConfiguration{Endpoint: c.componentConfig.Casdoor.Endpoint})
	if err != nil {
		panic(err)
	}

	controller.CasdoorController = casdoorController
	authenticationInterceptor, err := authentication.NewGCPAuthenticationInterceptor(&authentication.GCPAuthenticationInterceptorConfiguration{
		XTokenHeader:      c.componentConfig.GetTokenKey(),
		CasdoorController: casdoorController,
	})
	if err != nil {
		panic(err)
	}

	controller.Interceptors = interceptors.NewInterceptor()
	controller.Interceptors.SetAuthenticationInterceptor(authenticationInterceptor)

	controller.appController = app.NewAppController(&controller)
	repo, err := repo.NewAppMysqlImpl(controller.ComponentConfig().GetMysqlUrl())
	if err != nil {
		return nil, fmt.Errorf("new app repo failed %v", err)
	}
	controller.appRepo = repo

	controller.appServer = services.NewAppServer(&controller)
	if err := controller.appServer.SyncAppStatus(); err != nil {
		log.Infof("sync app status failed: %v", err)
	}
	return &controller, nil
}

// ****** Controller Funcs ******
func (c *Controller) HttpClient() *resty.Client {
	return c.httpclient
}

func (c *Controller) AppController() framework.AppControllerInterface {
	return c.appController
}

func (c *Controller) AppRepo() framework.AppRepoInterface {
	return c.appRepo
}

func (c *Controller) ComponentConfig() *v1.VclusterGatewayConfiguration {
	return c.cfg.componentConfig
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

	interceptorsList := c.Interceptors.InterceptorsList()

	// 检查interceptorsList的长度，确保至少有三个元素
	if len(interceptorsList) >= 3 {
		// 删除第三个元素
		interceptorsList = append(interceptorsList[:2], interceptorsList[3:]...)
	}

	interceptors := connect.WithInterceptors(interceptorsList...)

	mux := http.NewServeMux()
	svcPath, handler := appv1connect.NewAppServiceHandler(c.appServer, interceptors)
	mux.Handle(svcPath, handler)
	// newhandler := Cors(mux)
	srv := &http.Server{
		Addr:    ":" + "8088",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
		// Handler: h2c.NewHandler(newhandler, &http2.Server{}),
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof(" ****** App Gateway Server is Start at [:%s] ****** ", "8088")

	conn, err := grpc.Dial(
		"localhost:8088",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	gwmux := runtime.NewServeMux(
		c.Interceptors.SetTokenHeader(),
		c.Interceptors.SetGWApiKeyHeader(),
		c.Interceptors.SetHttpConfigs(),
		c.Interceptors.HandleGRPCError(),
		runtime.WithMetadata(func(ctx context.Context, r *http.Request) metadata.MD {
			md := metadata.MD{}
			md.Set(c.cfg.componentConfig.GetTokenKey(), r.Header.Get(c.cfg.componentConfig.GetTokenKey()))
			if !strings.Contains(r.URL.Path, "health") {
				log.Infof("http method: %v", r.Method)
				log.Infof("http Path: %v", r.URL.Path)
			}
			return md
		}),
		// runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
		// 	w.Header().Set("Access-Control-Allow-Origin", "https://bsm.ops01.zetyun.cn")
		// 	w.Header().Set("Access-Control-Allow-Credentials", "true")
		// 	return nil
		// }),
	)
	// RegisterVClusterGatewayServiceHandler
	// log.Info("mysql url: %v", c.ComponentConfig().GetMysqlUrl())
	err = appv1.RegisterAppServiceHandler(context.Background(), gwmux, conn)
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
