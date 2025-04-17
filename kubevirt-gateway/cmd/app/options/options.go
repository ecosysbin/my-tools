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
	appconfig "gitlab.datacanvas.com/aidc/kubevirt-gateway/cmd/app/config"
	configv1 "gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/config/kubevirt_gateway/v1"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/kube"

	"github.com/endverse/go-kit/aggregate"
	"github.com/spf13/pflag"
)

type Options struct {
	KubevirtGateway *KubevirtGatewayOptions
}

func NewOptions() (*Options, error) {
	o := &Options{
		KubevirtGateway: &KubevirtGatewayOptions{KubevirtGatewayConfiguration: &configv1.KubevirtGatewayConfiguration{}},
	}

	return o, nil
}

func (o *Options) Flags(fs *pflag.FlagSet) {
	o.KubevirtGateway.AddFlags(fs)
	// return fss
}

func (o *Options) ApplyTo(c *appconfig.Config) error {
	if err := o.KubevirtGateway.ApplyTo(c.ComponentConfig); err != nil {
		return err
	}

	return nil
}

func (o *Options) Validate() error {
	var errs []error

	errs = append(errs, o.KubevirtGateway.Validate()...)

	return aggregate.NewAggregate(errs)
}

func (o *Options) Config() (*appconfig.Config, error) {
	c := &appconfig.Config{
		ComponentConfig: o.KubevirtGateway.KubevirtGatewayConfiguration,
	}
	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}
	kubeConfig, err := kube.KubeConfig(c.ComponentConfig.GetHttpKubeApiserver(), c.ComponentConfig.GetHttpKubeConfig())
	if err != nil {
		return nil, err
	}
	c.KubeConfig = kubeConfig
	return c, nil
}
