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
	"os"
	"time"

	appconfig "gitlab.datacanvas.com/aidc/kubevirt-gateway/cmd/app/config"

	"gitlab.datacanvas.com/aidc/kubevirt-gateway/cmd/app/options"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/pkg/controller"
	"gitlab.datacanvas.com/aidc/kubevirt-gateway/version"

	"gitlab.datacanvas.com/aidc/gcpctl/gokit/log"

	"github.com/endverse/go-kit/signals"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func NewKubevirtGatewayCommand() *cobra.Command {
	cmd := newKubevirtGatewayCommand()
	cmd.AddCommand(version.NewVersionCommand())

	log.AddFlags(cmd.PersistentFlags(), cmd.Name())

	return cmd
}

func newKubevirtGatewayCommand() *cobra.Command {
	opts, err := options.NewOptions()
	if err != nil {
		log.Fatalf("unable to initialize command options: %v", err)
	}

	cmd := &cobra.Command{
		Use:   "kubevirt-gateway",
		Short: "Run GCP KubevirtGateway Gateway Server.",
		Run: func(cmd *cobra.Command, args []string) {
			logger := log.InitLogger()
			defer logger.Sync()
			// set up signals so we handle the first shutdown signal gracefully
			stopCh := signals.SetupSignalHandler()

			if err := runCommand(opts, stopCh); err != nil {
				log.Fatal(err)
			}
		},
	}
	opts.Flags(cmd.PersistentFlags())
	return cmd
}

func runCommand(opts *options.Options, stopCh <-chan struct{}) error {
	config, err := opts.Config()
	if err != nil {
		log.Info("get config err, %v", err)
		return err
	}

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

	return Run(ctx, config, stopCh)
}

func Run(ctx context.Context, config *appconfig.Config, stopCh <-chan struct{}) error {
	// Start watch config file change.
	if config.ComponentConfig.EnableWatch {
		go config.ComponentConfig.Watch(ctx)
	}

	ctrl := controller.New(config, stopCh)

	name, err := os.Hostname()
	if err != nil {
		panic("get hostname failed")
	}
	// 抢占主实例
	rl, err := resourcelock.NewFromKubeconfig(
		resourcelock.EndpointsLeasesResourceLock,
		config.KubeConfig.Leaderelection.Namespace,
		config.KubeConfig.Leaderelection.Name,
		resourcelock.ResourceLockConfig{
			Identity: name + "_" + string(uuid.NewUUID()),
		},
		config.KubeConfig.KubeRestConfig,
		15*time.Second)

	if err != nil {
		return err
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: rl,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				if err := ctrl.Run(ctx.Done()); err != nil {
					panic(err)
				}
			},
			OnStoppedLeading: func() {
				// we can do cleanup here
				log.Warn("leaderelection lost")
				os.Exit(1)
			},
		},
	})

	log.Info("Stop Server With Graceful.")

	return nil
}
