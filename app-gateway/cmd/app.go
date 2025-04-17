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

package main

import (
	"os"

	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"

	"gitlab.datacanvas.com/aidc/app-gateway/cmd/app"
)

//	@title			GCP VclusterGateway Gateway API
//	@version		1.0
//	@description	This is an authorization authentication proxy service.

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host	localhost:8083
//	@BasePath

// @securityDefinitions.basic	BasicAuth
func main() {
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)

	command := app.NewAppGatewayCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
