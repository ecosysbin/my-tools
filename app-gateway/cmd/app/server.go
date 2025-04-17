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

package app

import (
	"context"

	"github.com/endverse/go-kit/signals"
	"github.com/endverse/go-kit/tmpl"
	"github.com/spf13/cobra"
	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"

	"gitlab.datacanvas.com/aidc/app-gateway/cmd/app/config"
	"gitlab.datacanvas.com/aidc/app-gateway/cmd/app/options"
	"gitlab.datacanvas.com/aidc/app-gateway/pkg/controller"
)

func NewAppGatewayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "app-gateway",
	}

	cmd.AddCommand(newAppGatewayCommand())

	log.AddFlags(cmd.PersistentFlags(), cmd.Name())

	return cmd
}

func newAppGatewayCommand() *cobra.Command {
	opts, err := options.NewOptions()
	if err != nil {
		log.Fatalf("unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run GCP App Gateway Server.",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.InitLogger()
			defer logger.Sync()

			// set up signals, so we handle the first shutdown signal gracefully
			stopCh := signals.SetupSignalHandler()

			if err := runCommand(opts, stopCh); err != nil {
				log.Fatal(err)
			}
		},
	}

	fs := cmd.Flags()
	namedFlagSets := opts.Flags()
	for _, f := range namedFlagSets.FlagSets {
		fs.AddFlagSet(f)
	}

	tmpl.SetHelpAndUsageFunc(cmd, namedFlagSets)

	return cmd
}

func runCommand(opts *options.Options, stopCh <-chan struct{}) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	c, err := opts.Config()
	if err != nil {
		return err
	}

	// Get the completed config
	cc := c.Complete()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// listen for interrupts or the Linux SIGTERM signal and cancel
	// our context, which the leader election code will observe and
	// step down
	go func() {
		<-stopCh
		log.Info("Received termination, signaling shutdown")
		cancel()
	}()

	return Run(ctx, cc, stopCh)
}

func Run(ctx context.Context, cc config.CompletedConfig, stopCh <-chan struct{}) error {
	// Start watch config file change.
	if cc.ComponentConfig.EnableWatch {
		go cc.ComponentConfig.Watch(ctx)
	}

	ctrl, err := controller.New(
		controller.WithComponentConfig(cc.ComponentConfig),
	)
	if err != nil {
		return err
	}
	if err := ctrl.Run(ctx.Done()); err != nil {
		return err
	}

	log.Info("Stop Server With Graceful.")

	return nil
}
