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

package version

import (
	"fmt"
	"runtime"

	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/apis/response"

	gcpctx "gitlab.datacanvas.com/aidc/gcpctl/gokit/gin/context"

	"github.com/endverse/go-kit/tmpl"
	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
)

var (
	// Version shows the version of kubevirt-gateway.
	Version = "Not provided."
	// GitCommit shows the git commit id of  kubevirt-gateway.
	GitCommit = "Not provided."
	// BuildAt shows the built time of the binary.
	BuildAt = "Not provided."

	apiVersion = "v1"
)

func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print Versions of GCP kubevirt-gateway Gataway.",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}

	tmpl.SetHelpAndUsageFunc(cmd, cliflag.NamedFlagSets{})

	return cmd
}

var versionFormat = `
****************************************
* GCP kubevirt-gateway Gataway:
*
*     Version:            %s
*     GOOS:               %s
*     GOARCH:             %s
*     Git Commit:         %s
*     Build Time:         %s
*     API Version:        %s
****************************************
`

func printVersion() {
	fmt.Printf(versionFormat, Version, runtime.GOOS, runtime.GOARCH, GitCommit, BuildAt, apiVersion)
}

type versionInformation struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildAt   string `json:"buildAt"`
}

func VersionInformation(c *gcpctx.GCPContext) {
	response.Response(c, response.SuccessGCPResponse, &versionInformation{
		Version:   Version,
		GitCommit: GitCommit,
		BuildAt:   BuildAt,
	})
}
