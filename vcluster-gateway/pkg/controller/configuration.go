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
	configv1 "vcluster-gateway/pkg/apis/config/vcluster_gateway/v1"
)

type controllerConfiguration struct {
	componentConfig *configv1.VclusterGatewayConfiguration
}

type Option func(*controllerConfiguration)

func WithComponentConfig(config *configv1.VclusterGatewayConfiguration) Option {
	return func(cc *controllerConfiguration) {
		cc.componentConfig = config
	}
}
