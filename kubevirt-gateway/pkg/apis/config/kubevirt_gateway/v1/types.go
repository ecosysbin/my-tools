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

package v1

import (
	"context"
	"fmt"
	"os"

	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller/framework"

	"github.com/fsnotify/fsnotify"
	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"
	"gopkg.in/yaml.v3"
)

var _ framework.ComponentConfigInterface = &KubevirtGatewayConfiguration{}

type KubevirtGatewayConfiguration struct {
	ConfigFilePath     string `json:"configFilePath" yaml:"configFilePath"`
	PlatformConfigPath string `json:"platformConfigPath" yaml:"platformConfigPath"`
	EnableWatch        bool   `json:"enableWatch" yaml:"enableWatch"`

	Server   Server   `json:"server" yaml:"server"`
	Http     Http     `json:"http" yaml:"http"`
	Storage  Storage  `json:"storage" yaml:"storage"`
	Images   Images   `json:"images" yaml:"images"`
	Products Products `json:"products" yaml:"products"`
	Mysql    Mysql    `json:"mysql" yaml:"mysql"`
	Logger   Logger   `json:"logger" yaml:"logger"`
	Casdoor  Casdoor  `json:"casdoor" yaml:"casdoor"`
	Platform Platform `json:"platform" yaml:"platform"`
}

type Logger struct {
	LogLevel string `json:"products" yaml:"loglevel"`
}

type Server struct {
	IP           string                 `json:"-" yaml:"-" mapstructure:"-"`
	Port         string                 `json:"port" yaml:"port"`
	TokenKey     string                 `json:"tokenKey" yaml:"tokenKey"`
	RespCacheKey string                 `json:"respCacheKey" yaml:"respCacheKey"`
	Products     map[string]interface{} `json:"products" yaml:"products"`
}

type Http struct {
	KubeApiserver string `json:"kubeApiserver" yaml:"kubeApiserver"`
	KubeConfig    string `json:"kubeConfig" yaml:"kubeConfig"`
	VncServer     string `json:"vncServer" yaml:"vncServer"`
}

type Storage struct {
	Min     int64 `json:"min" yaml:"min"`
	Max     int64 `json:"max" yaml:"max"`
	Default int64 `json:"default" yaml:"default"`
}

type Images struct {
	ImageFilePath string `json:"imageFilePath" yaml:"imageFilePath"`
}

type Products struct {
	ProductFilePath string `json:"productFilePath" yaml:"productFilePath"`
}

type Mysql struct {
	PoolMax int    `json:"poolMax" yaml:"poolMax"`
	Url     string `json:"url" yaml:"url"`
}

type Casdoor struct {
	Endpoint         string `json:"endpoint" yaml:"endpoint"`
	ClientId         string `json:"clientId" yaml:"clientId"`
	ClientSecret     string `json:"clientSecret" yaml:"clientSecret"`
	OrganizationName string `json:"organizationName" yaml:"organizationName"`
	ApplicationName  string `json:"applicationName" yaml:"applicationName"`
	Certificate      string `json:"certificate" yaml:"certificate"`
}

type Platform struct {
	Domain string `json:"domain" yaml:"domain"`
}

func (c *KubevirtGatewayConfiguration) GetServerIP() string {
	return c.Server.IP
}

func (c *KubevirtGatewayConfiguration) GetServerPort() string {
	return c.Server.Port
}

func (c *KubevirtGatewayConfiguration) GetTokenKey() string {
	return c.Server.TokenKey
}

func (c *KubevirtGatewayConfiguration) GetRespCacheKey() string {
	return c.Server.RespCacheKey
}

func (c *KubevirtGatewayConfiguration) GetHttpKubeApiserver() string {
	return c.Http.KubeApiserver
}

func (c *KubevirtGatewayConfiguration) GetHttpKubeConfig() string {
	return c.Http.KubeConfig
}

func (c *KubevirtGatewayConfiguration) GetHttpVncServer() string {
	return c.Http.VncServer
}

func (c *KubevirtGatewayConfiguration) GetLoggerLogLevel() string {
	return c.Logger.LogLevel
}

func (c *KubevirtGatewayConfiguration) GetMysqlPoolMax() int {
	return c.Mysql.PoolMax
}

func (c *KubevirtGatewayConfiguration) GetMysqlUrl() string {
	return c.Mysql.Url
}

func (c *KubevirtGatewayConfiguration) GetCasdoorEndpoint() string {
	return c.Casdoor.Endpoint
}

func (c *KubevirtGatewayConfiguration) GetCasdoorClientId() string {
	return c.Casdoor.ClientId
}

func (c *KubevirtGatewayConfiguration) GetCasdoorClientSecret() string {
	return c.Casdoor.ClientSecret
}

func (c *KubevirtGatewayConfiguration) GetCasdoorOrganizationName() string {
	return c.Casdoor.OrganizationName
}

func (c *KubevirtGatewayConfiguration) GetCasdoorApplicationName() string {
	return c.Casdoor.ApplicationName
}

func (c *KubevirtGatewayConfiguration) GetCasdoorCertificate() string {
	return c.Casdoor.Certificate
}

func (c *KubevirtGatewayConfiguration) GetProducts() map[string]interface{} {
	return c.Server.Products
}

func (c *KubevirtGatewayConfiguration) GetPlatformDomain() string {
	return c.Platform.Domain
}

func (c *KubevirtGatewayConfiguration) GetStorageMin() int64 {
	return c.Storage.Min
}

func (c *KubevirtGatewayConfiguration) GetStorageMax() int64 {
	return c.Storage.Max
}

func (c *KubevirtGatewayConfiguration) GetStorageDefault() int64 {
	return c.Storage.Default
}

func (c *KubevirtGatewayConfiguration) GetImagesFilePath() string {
	return c.Images.ImageFilePath
}

func (c *KubevirtGatewayConfiguration) GetProductsFilePath() string {
	return c.Products.ProductFilePath
}

func (c *KubevirtGatewayConfiguration) ReadConfFromFile() error {
	file, err := os.ReadFile(c.ConfigFilePath)
	if err != nil {
		log.Infof("read config file %s err, %v", c.ConfigFilePath, err)
		return err
	}

	if err := yaml.Unmarshal(file, c); err != nil {
		return fmt.Errorf("unmarshl config file err, %v", err)
	}
	// 读取平台配置
	platConfigfile, err := os.ReadFile(c.PlatformConfigPath)
	if err != nil {
		log.Infof("read platConfigfile file %s err, %v", c.PlatformConfigPath, err)
		return err
	}
	platform := Platform{}
	if err := yaml.Unmarshal(platConfigfile, &platform); err != nil {
		return fmt.Errorf("unmarshl config file err, %v", err)
	}
	c.Platform = platform
	return nil
}

func (c *KubevirtGatewayConfiguration) Watch(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Warnf("Create new config file watcher failed with error: %v", err)
		return
	}

	defer watcher.Close()

	if err := watcher.Add(c.ConfigFilePath); err != nil {
		log.Warnf("Watching config file %s with error: %v", c.ConfigFilePath, err)
		return
	}

	if err := watcher.Add(c.PlatformConfigPath); err != nil {
		log.Warnf("Watching config file %s with error: %v", c.ConfigFilePath, err)
		return
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Infof("Config file %s is changed, refresh the configuration.", c.ConfigFilePath)

				if err := c.ReadConfFromFile(); err != nil {
					log.Warnf("Refresh config file with error: %v", err)
				}

				log.Infof("refresh conf: %#v", c)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Warnf("Watch config file %s with error: %v", c.ConfigFilePath, err)

		case <-ctx.Done():
			log.Info("Get shutdown signal.")
			return
		}
	}
}
