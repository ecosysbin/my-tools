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

package options

import (
	"fmt"

	"github.com/spf13/pflag"
	configv1 "gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/config/kubevirt_gateway/v1"
)

type KubevirtGatewayOptions struct {
	*configv1.KubevirtGatewayConfiguration
}

func (o *KubevirtGatewayOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.ConfigFilePath, "config", o.ConfigFilePath, "Config file path.")
	fs.StringVar(&o.PlatformConfigPath, "gcp-config", o.PlatformConfigPath, "gcp-Config file path.")
	fs.BoolVar(&o.EnableWatch, "enable-watch", o.EnableWatch, "If true, will auto watch config file and refresh configuration.")

	fs.StringVar(&o.Server.IP, "ip", o.Server.IP, "Http server listen address.")
	fs.StringVar(&o.Server.Port, "port", o.Server.Port, "Http server listen port.")
	fs.StringVar(&o.Server.TokenKey, "token-key", o.Server.TokenKey, "Http server token key.")
	fs.StringVar(&o.Server.RespCacheKey, "resp-cache-key", o.Server.RespCacheKey, "Http server resp cache key.")

	// NOTE:
	// If you have input a custom config,
	// then you need to add flags at this point.
	//
	// fs.StringVar(&o.Example.WebServer, "example-webserver", o.Example.WebServer, "Example web server address.")

	fs.StringVar(&o.Casdoor.Endpoint, "casdoor-endpoint", o.Casdoor.Endpoint, "Casdoor server endpoint.")
	fs.StringVar(&o.Casdoor.ClientId, "casdoor-client-id", o.Casdoor.ClientId, "Casdoor server client id.")
	fs.StringVar(&o.Casdoor.ClientSecret, "casdoor-client-secret", o.Casdoor.ClientSecret, "Casdoor server client secret.")
	fs.StringVar(&o.Casdoor.OrganizationName, "casdoor-organization-name", o.Casdoor.OrganizationName, "Casdoor server organization name.")
	fs.StringVar(&o.Casdoor.ApplicationName, "casdoor-application-name", o.Casdoor.ApplicationName, "Casdoor server application name.")
	fs.StringVar(&o.Casdoor.Certificate, "casdoor-certificate", o.Casdoor.Certificate, "Casdoor server certificate.")
}

const (
	DEFAULT_CONFIG_PATH          = "/etc/config.yaml"
	DEFAULT_PLATFORM_CONFIG_PATH = "/etc/gcp.yaml"
)

func (o *KubevirtGatewayOptions) ApplyTo(cfg *configv1.KubevirtGatewayConfiguration) error {
	// logger := log.InitLogger()
	if o.ConfigFilePath == "" {
		o.ConfigFilePath = DEFAULT_CONFIG_PATH
	}
	if o.PlatformConfigPath == "" {
		o.PlatformConfigPath = DEFAULT_PLATFORM_CONFIG_PATH
	}

	cfg.ConfigFilePath = o.ConfigFilePath
	if err := cfg.ReadConfFromFile(); err != nil {
		return err
	}
	if o.EnableWatch {
		cfg.EnableWatch = o.EnableWatch
		// If watch config file changes is enabled,
		// the command line arguments will be overwritten.
		return nil
	}

	if o.Server.TokenKey != "" {
		cfg.Server.TokenKey = o.Server.TokenKey
	}
	if o.Server.RespCacheKey != "" {
		cfg.Server.RespCacheKey = o.Server.RespCacheKey
	}

	// NOTE:
	// If you have input a custom config,
	// then you need to convert value into 'cfg' at this point.
	//
	// if o.Example.WebServer != "" {
	// 	   cfg.Example.WebServer = o.Example.WebServer
	// }

	if o.Casdoor.Endpoint != "" {
		cfg.Casdoor.Endpoint = o.Casdoor.Endpoint
	}
	if o.Casdoor.ClientId != "" {
		cfg.Casdoor.ClientId = o.Casdoor.ClientId
	}
	if o.Casdoor.ClientSecret != "" {
		cfg.Casdoor.ClientSecret = o.Casdoor.ClientSecret
	}
	if o.Casdoor.OrganizationName != "" {
		cfg.Casdoor.OrganizationName = o.Casdoor.OrganizationName
	}
	if o.Casdoor.ApplicationName != "" {
		cfg.Casdoor.ApplicationName = o.Casdoor.ApplicationName
	}
	if o.Casdoor.Certificate != "" {
		cfg.Casdoor.Certificate = o.Casdoor.Certificate
	}
	return nil
}

func (o *KubevirtGatewayOptions) Validate() []error {
	if o == nil {
		return nil
	}

	errs := []error{}

	if o.ConfigFilePath == "" {
		if o.Server.TokenKey == "" {
			errs = append(errs, fmt.Errorf("server TokenKey is must"))
		}

		if o.Server.RespCacheKey == "" {
			errs = append(errs, fmt.Errorf("server RespCacheKey is must"))
		}

		// NOTE:
		// If you have input a custom config,
		// then you need to perform validation at this point.
		//
		// if o.Example.WebServer == "" {
		//	  errs = append(errs, fmt.Errorf("example WebServer is must"))
		// }

		if o.Casdoor.Endpoint == "" {
			errs = append(errs, fmt.Errorf("casdoor Endpoint is must"))
		}

		if o.Casdoor.ClientId == "" {
			errs = append(errs, fmt.Errorf("casdoor ClientId is must"))
		}

		if o.Casdoor.ClientSecret == "" {
			errs = append(errs, fmt.Errorf("casdoor ClientSecret is must"))
		}

		if o.Casdoor.OrganizationName == "" {
			errs = append(errs, fmt.Errorf("casdoor OrganizationName is must"))
		}

		if o.Casdoor.ApplicationName == "" {
			errs = append(errs, fmt.Errorf("casdoor ApplicationName is must"))
		}

		if o.Casdoor.Certificate == "" {
			errs = append(errs, fmt.Errorf("casdoor Certificate is must"))
		}
	}

	return errs
}
