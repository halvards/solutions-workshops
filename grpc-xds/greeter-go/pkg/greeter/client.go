// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package greeter

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	orcav3 "github.com/cncf/xds/go/xds/data/orca/v3"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	xdscredentials "google.golang.org/grpc/credentials/xds"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	"github.com/googlecloudplatform/solutions-workshops/grpc-xds/greeter-go/pkg/interceptors"
	"github.com/googlecloudplatform/solutions-workshops/grpc-xds/greeter-go/pkg/logging"
	helloworldpb "github.com/googlecloudplatform/solutions-workshops/grpc-xds/greeter-go/third_party/google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	grpcClientDialTimeout      = 10 * time.Second
	grpcClientKeepaliveTime    = 30 * time.Second
	grpcClientKeepaliveTimeout = 5 * time.Second
	grpcClientIdleTimeout      = math.MaxInt64 // good idea?
	TrailerMetadataKey         = "endpoint-load-metrics-bin"
)

var errMultipleORCALoadReportsInMetadata = errors.New("multiple ORCA load reports found in provided metadata")

type Client struct {
	logger  logr.Logger
	nextHop string
	client  helloworldpb.GreeterClient
}

func NewClient(ctx context.Context, nextHop string, useXDSCredentials bool) (*Client, error) {
	logger := logging.FromContext(ctx)
	dialOpts, err := dialOptions(logger, useXDSCredentials)
	if err != nil {
		return nil, fmt.Errorf("could not configure greeter client connection dial options: %w", err)
	}
	clientConn, err := grpc.NewClient(nextHop, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("could not create a virtual connection to target=%s: %w", nextHop, err)
	}
	addClientConnectionCloseBehavior(ctx, logger, clientConn)
	return &Client{
		client:  helloworldpb.NewGreeterClient(clientConn),
		logger:  logger,
		nextHop: nextHop,
	}, nil
}

// ToLoadReport extracts a binary ORCA load report from the trailer metadata,
// using the metadata key "endpoint-load-metrics-bin". Returns a nil report
// reference and nil error if there was no load report in the metadata map.
// This function was adapted from
// https://github.com/grpc/grpc-go/blob/48b6b11b388f4b26e03fa722dd6ed60ae27a0fa7/orca/internal/internal.go#L49-L71
// Copyright 2022 gRPC authors. Licensed under the Apache License, Version 2.0.
func ToLoadReport(md metadata.MD) (*orcav3.OrcaLoadReport, error) {
	loadMetrics := md.Get(TrailerMetadataKey)
	if len(loadMetrics) == 0 {
		return nil, nil
	}
	if len(loadMetrics) > 1 {
		return nil, errMultipleORCALoadReportsInMetadata
	}
	loadReport := &orcav3.OrcaLoadReport{}
	if err := proto.Unmarshal([]byte(loadMetrics[0]), loadReport); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ORCA load report found in metadata: %w", err)
	}
	return loadReport, nil
}

// SayHello implements helloworld.Greeter/SayHello.
func (c *Client) SayHello(requestCtx context.Context, name string) (string, error) {
	var header, trailer metadata.MD
	resp, err := c.client.SayHello(requestCtx, &helloworldpb.HelloRequest{Name: name}, grpc.WaitForReady(true), grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return "", fmt.Errorf("could not greet name=%s at target=%s: %w", name, c.nextHop, err)
	}
	_ = header
	loadReport, err := ToLoadReport(trailer)
	if err != nil {
		return "", fmt.Errorf("could not get load report: %w", err)
	}
	return resp.GetMessage() + " " + loadReport.String(), nil
}

// dialOptions sets parameters for client connection establishment.
func dialOptions(logger logr.Logger, useXDSCredentials bool) ([]grpc.DialOption, error) {
	clientCredentials := insecure.NewCredentials()
	if useXDSCredentials {
		logger.V(1).Info("Using xDS client-side credentials, with insecure as fallback")
		var err error
		if clientCredentials, err = xdscredentials.NewClientCredentials(xdscredentials.ClientOptions{FallbackCreds: insecure.NewCredentials()}); err != nil {
			return nil, fmt.Errorf("could not create client-side transport credentials for xDS: %w", err)
		}
	}
	return []grpc.DialOption{
		grpc.WithChainStreamInterceptor(interceptors.StreamClientLogging(logger)),
		grpc.WithChainUnaryInterceptor(interceptors.UnaryClientLogging(logger)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcClientKeepaliveTime,
			Timeout:             grpcClientKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithIdleTimeout(time.Duration(grpcClientIdleTimeout)),
		grpc.WithTransportCredentials(clientCredentials),
	}, nil
}

func addClientConnectionCloseBehavior(ctx context.Context, logger logr.Logger, clientConn *grpc.ClientConn) {
	go func(cc *grpc.ClientConn) {
		<-ctx.Done()
		logger.Info("Closing the greeter client connection")
		err := cc.Close()
		if err != nil {
			logger.Error(err, "Error when closing the greeter client connection")
		}
	}(clientConn)
}
