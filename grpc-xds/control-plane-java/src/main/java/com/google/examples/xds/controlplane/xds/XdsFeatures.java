// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package com.google.examples.xds.controlplane.xds;

import java.util.Map;
import java.util.Objects;
import org.jetbrains.annotations.NotNull;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Contains flags to enable and disable xDS features. */
public record XdsFeatures(
    boolean serverListenerUsesRds,
    boolean enableControlPlaneTls,
    boolean requireControlPlaneClientCerts,
    boolean enableDataPlaneTls,
    boolean requireDataPlaneClientCerts,
    boolean enableFederation) {

  private static final Logger LOG = LoggerFactory.getLogger(XdsFeatures.class);

  /** Canonical constructor. */
  public XdsFeatures {
    if (!enableControlPlaneTls && requireControlPlaneClientCerts) {
      throw new IllegalArgumentException(
          "xDS feature flags: enableControlPlaneTls=true is required when"
              + " requireControlPlaneClientCerts=true");
    }
    if (!enableDataPlaneTls && requireDataPlaneClientCerts) {
      throw new IllegalArgumentException(
          "xDS feature flags: enableDataPlaneTls=true is required when"
              + " requireDataPlaneClientCerts=true");
    }
    if (serverListenerUsesRds) {
      LOG.warn(
          "xDS clients implemented using Go must use gRPC-Go v1.61.0 or later for dynamic RouteConfiguration via RDS for server Listeners, see https://github.com/grpc/grpc-go/issues/6788");
    }
  }

  /**
   * Constructor used after parsing a YAML file (or similar) to a Map.
   *
   * <p>The map keys must match the instance variable names in this class.
   *
   * <p>Missing or null flags default to false.
   *
   * @param features xDS feature toggles
   */
  public XdsFeatures(@NotNull Map<String, Boolean> features) {
    this(
        Objects.requireNonNullElse(features.get("serverListenerUsesRds"), Boolean.FALSE),
        Objects.requireNonNullElse(features.get("enableControlPlaneTls"), Boolean.FALSE),
        Objects.requireNonNullElse(features.get("requireControlPlaneClientCerts"), Boolean.FALSE),
        Objects.requireNonNullElse(features.get("enableDataPlaneTls"), Boolean.FALSE),
        Objects.requireNonNullElse(features.get("requireDataPlaneClientCerts"), Boolean.FALSE),
        Objects.requireNonNullElse(features.get("enableFederation"), Boolean.FALSE));
  }
}
