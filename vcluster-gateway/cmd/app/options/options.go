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
	"github.com/endverse/go-kit/aggregate"
	cliflag "k8s.io/component-base/cli/flag"

	appconfig "vcluster-gateway/cmd/app/config"
	configv1 "vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
)

type Options struct {
	VclusterGateway *VclusterGatewayOptions
}

func NewOptions() (*Options, error) {
	cfg, err := newDefaultComponentConfig()
	if err != nil {
		return nil, err
	}

	o := &Options{
		VclusterGateway: &VclusterGatewayOptions{VclusterGatewayConfiguration: cfg},
	}

	return o, nil
}

func newDefaultComponentConfig() (*configv1.VclusterGatewayConfiguration, error) {
	cfg := configv1.VclusterGatewayConfiguration{}
	if err := configv1.Default(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (o *Options) Flags() cliflag.NamedFlagSets {
	fss := cliflag.NamedFlagSets{}
	o.VclusterGateway.AddFlags(fss.FlagSet("VclusterGateway"))

	return fss
}

func (o *Options) ApplyTo(c *appconfig.Config) error {
	if err := o.VclusterGateway.ApplyTo(c.ComponentConfig); err != nil {
		return err
	}

	return nil
}

func (o *Options) Validate() error {
	var errs []error

	errs = append(errs, o.VclusterGateway.Validate()...)

	return aggregate.NewAggregate(errs)
}

func (o *Options) Config() (*appconfig.Config, error) {
	c := &appconfig.Config{
		ComponentConfig: o.VclusterGateway.VclusterGatewayConfiguration,
	}
	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}
