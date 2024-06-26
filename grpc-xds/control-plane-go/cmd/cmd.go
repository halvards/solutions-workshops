// Copyright 2023 Google LLC
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

package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/googlecloudplatform/solutions-workshops/grpc-xds/control-plane-go/pkg/auth"
	"github.com/googlecloudplatform/solutions-workshops/grpc-xds/control-plane-go/pkg/config"
	"github.com/googlecloudplatform/solutions-workshops/grpc-xds/control-plane-go/pkg/informers"
	"github.com/googlecloudplatform/solutions-workshops/grpc-xds/control-plane-go/pkg/logging"
	"github.com/googlecloudplatform/solutions-workshops/grpc-xds/control-plane-go/pkg/server"
	"github.com/googlecloudplatform/solutions-workshops/grpc-xds/control-plane-go/pkg/signals"
)

func Run(ctx context.Context, flagset *flag.FlagSet, args []string) error {
	ctx = signals.SetupSignalHandler(ctx)
	logging.InitFlags(flagset)
	informers.InitFlags(flagset)
	if err := flagset.Parse(args); err != nil {
		return fmt.Errorf("could not parse command line flags args=%+v: %w", args, err)
	}
	logger := logging.NewLogger()
	logging.SetGRPCLogger(logger)
	ctx = logging.NewContext(ctx, logger)
	auth.RegisterAll(ctx, logger)
	servingPort, err := config.ServingPort()
	if err != nil {
		return fmt.Errorf("could not configure management server listening port: %w", err)
	}
	healthPort, err := config.HealthPort()
	if err != nil {
		return fmt.Errorf("could not configure management server health checking port: %w", err)
	}
	kubecontexts, err := config.Kubecontexts(logger)
	if err != nil {
		return fmt.Errorf("could not initialize informer configuration: %w", err)
	}
	xdsFeatures, err := config.XDSFeatures(logger)
	if err != nil {
		return fmt.Errorf("could not initialize xDS feature flags: %w", err)
	}
	authority, err := config.AuthorityName(logger)
	if err != nil {
		return fmt.Errorf("could not determine control plane authority name: %w", err)
	}
	if xdsFeatures.EnableFederation {
		logger.V(2).Info("Enabling xDS federation", "authority", authority)
	}
	return server.Run(ctx, servingPort, healthPort, kubecontexts, xdsFeatures, authority)
}
