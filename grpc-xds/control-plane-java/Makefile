# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Set global Skaffold options
export SKAFFOLD_BUILD_CONCURRENCY:=0
export SKAFFOLD_CLEANUP:=false
export SKAFFOLD_DETECT_MINIKUBE:=false
export SKAFFOLD_INTERACTIVE:=false
export SKAFFOLD_SKIP_TESTS:=true
export SKAFFOLD_UPDATE_CHECK:=false

.PHONY: all
all: clean proto format lint test

# run builds and deploys the xDS control plane and sets up port forwarding.
.PHONY: run
run:
	skaffold run --port-forward=user

# run-no-port-forward builds and deploys the xDS control plane.
.PHONY: run-no-port-forward
run-no-port-forward:
	skaffold run --port-forward=off

# debug builds and deploys the xDS control plane, and sets up remote debugging.
.PHONY: debug
debug:
	skaffold debug --port-forward=user

# dev builds and deploys the xDS control plane, and sets up file watching for automatic image
# rebuilds and redeploys on source code changes.
.PHONY: dev
dev:
	skaffold dev

# tail the logs of the xDS control plane.
.PHONY: tail
tail:
	kubectl logs --all-containers --follow --namespace=xds deployment/control-plane

# delete the xds control plane Deployment and Service resources from the k8s cluster.
.PHONY: delete
delete:
	kubectl delete --ignore-not-found --namespace=xds deployment control-plane
	kubectl delete --ignore-not-found --namespace=xds service control-plane
	kubectl delete --ignore-not-found --namespace=xds configmaps --selector="app.kubernetes.io/part-of"=grpc-xds

.PHONY: build
build:
	../gradlew assembleDist

.PHONY: clean
clean:
	../gradlew clean

.PHONY: depupdates
depupdates:
	../gradlew dependencyUpdates

.PHONY: format
format:
	../gradlew googleJavaFormat

.PHONY: image
image:
	skaffold build

.PHONY: imagejib
imageko:
	../gradlew jibDockerBuild

.PHONY: lint
lint: format
	../gradlew verifyGoogleJavaFormat checkstyleMain spotbugsMain

.PHONY: manifests
manifests:
	skaffold render --digest-source=local --loud=false

.PHONY: proto
proto:
	../gradlew generateProto

.PHONY: test
test:
	../gradlew test
