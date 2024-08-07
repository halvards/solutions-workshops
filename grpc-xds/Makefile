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

# kind-create-multi-cluster creates two kind Kubernetes cluster, each with multiple nodes.
# It also installs prerequisites such as cert-manager in both clusters.
.PHONY: kind-create-multi-cluster
kind-create-multi-cluster: check-kind check-skaffold check-yq
	./hack/kind-multi-cluster.sh

# kind-create creates a kind Kubernetes cluster with multiple nodes.
# It also installs prerequisites such as cert-manager.
.PHONY: kind-create
kind-create: kind-create-cluster-only cert-manager

# kind-create-cluster-only creates a kind Kubernetes cluster with multiple nodes.
.PHONY: kind-create-cluster-only
kind-create-cluster-only: check-kind check-kubectl
	kind create cluster --config=hack/kind-cluster-config.yaml
	kubectl config set-context --current --namespace=xds

# kind-delete-multi-cluster deletes the two kind Kubernetes clusters.
.PHONY: kind-delete-multi-cluster
kind-delete-multi-cluster: check-kind
	kind delete cluster --name=grpc-xds
	kind delete cluster --name=grpc-xds-2

# kind-delete deletes the kind Kubernetes cluster.
.PHONY: kind-delete
kind-delete: check-kind
	kind delete cluster --name=grpc-xds

# run-go-multi-cluster builds and deploys control-plane-go and greeter-go to one cluster,
# and greeter-go to the other cluster.
.PHONY: run-go-multi-cluster
run-go-multi-cluster: check-skaffold check-go
	./hack/deploy-multi-cluster.sh go

# run-go-federation builds and deploys the xDS control plane and greeter-intermediary to one cluster,
# and the xDS control plane and greeter-leaf to the other cluster.
.PHONY: run-go-federation
run-go-federation: check-skaffold check-go
	./hack/deploy-federation.sh go

# run-go builds and deploys control-plane-go and greeter-go.
.PHONY: run-go
run-go: check-skaffold check-go
	skaffold run --module=go --port-forward=user

# run-go-td builds and deploys greeter-go for use with Traffic Director/Cloud Service Mesh
# on a GKE cluster.
.PHONY: run-go-td
run-go-td: check-skaffold check-go
	skaffold run --module=go-td --port-forward=user --profile=traffic-director

# run-go-td-off-gcp builds and deploys greeter-go for use with Traffic Director/Cloud Service Mesh
# on a non-GKE Kubernetes cluster.
.PHONY: run-go-td-off-gcp
run-go-td-off-gcp: check-skaffold check-go
	skaffold run --module=go-td --port-forward=user --profile=traffic-director-off-gcp

# run-java-multi-cluster builds and deploys control-plane-java and greeter-java to one cluster,
# and greeter-java to the other cluster.
.PHONY: run-java-multi-cluster
run-java-multi-cluster: check-skaffold check-javac
	./hack/deploy-multi-cluster.sh java

# run-java-federation builds and deploys the xDS control plane and greeter-intermediary to one cluster,
# and the xDS control plane and greeter-leaf to the other cluster.
.PHONY: run-java-federation
run-java-federation: check-skaffold check-javac
	./hack/deploy-federation.sh java

# run-java builds and deploys control-plane-java and greeter-java.
.PHONY: run-java
run-java: check-skaffold check-javac
	skaffold run --module=java --port-forward=user

# run-java-td builds and deploys greeter-java for use with Traffic Director/Cloud Service Mesh
# on a GKE cluster.
.PHONY: run-java-td
run-java-td: check-skaffold check-javac
	skaffold run --module=java-td --port-forward=user --profile=traffic-director

# run-java-td-off-gcp builds and deploys greeter-java for use with
# Traffic Director/Cloud Service Mesh on a non-GKE Kubernetes cluster.
.PHONY: run-java-td-off-gcp
run-java-td-off-gcp: check-skaffold check-javac
	skaffold run --module=java-td --port-forward=user --profile=traffic-director-off-gcp

# run-go-tls-cert-manager builds and deploys control-plane-go and greeter-go,
# with workload TLS certificates from CA managed by cert-manager,
# even if the current kubecontext points to a GKE cluster.
.PHONY: run-go-tls-cert-manager
run-go-tls-cert-manager: check-skaffold check-go
	skaffold run --module=go --port-forward=user --profile=tls-gke-workload-certs,-tls-cert-manager

# run-java-tls-cert-manager builds and deploys control-plane-java and greeter-java,
# with workload TLS certificates from CA managed by cert-manager,
# even if the current kubecontext points to a GKE cluster.
.PHONY: run-java-tls-cert-manager
run-java-tls-cert-manager: check-skaffold check-javac
	skaffold run --module=java --port-forward=user --profile=tls-gke-workload-certs,-tls-cert-manager

# debug-go builds and deploys control-plane-go and greeter-go, with remote debugging enabled.
.PHONY: debug-go
debug-go: check-skaffold check-go
	skaffold debug --module=go --port-forward=user

# debug-java builds and deploys control-plane-java and greeter-java, with remote debugging enabled.
.PHONY: debug-java
debug-java: check-skaffold check-javac
	skaffold debug --module=java --port-forward=user

# dev-go sets up file watching with automatic image rebuilds and redeploys for control-plane-go and greeter-go.
.PHONY: dev-go
dev-go: check-skaffold check-go
	skaffold dev --module=go

# dev-java sets up file watching with automatic image rebuilds and redeploys for control-plane-java and greeter-java.
.PHONY: dev-java
dev-java: check-skaffold check-javac
	skaffold dev --module=java

# host-certs does the following:
# 1. creates a temporary pod with a workload TLS certificate
# 2. copies the certificate to the local developer workstation
# 3. deletes the temporary pod and associated k8s resources
.PHONY: host-certs
host-certs: check-skaffold
	skaffold run --filename=k8s/host-certs/skaffold.yaml
	skaffold delete --filename=k8s/host-certs/skaffold.yaml
	@echo "*************************************************"
	@echo " Certificate Validity :"
	@openssl x509 -text -noout -in certs/certificates.pem | grep After
	@echo "*************************************************"

# copy-bastion-certs copies the workload TLS certificates from the bastion pod to the local developer workstation.
.PHONY: copy-bastion-certs
copy-bastion-certs: check-kubectl
	mkdir -p certs
	kubectl exec deployment/bastion --namespace=xds --container=app -- tar -chf - /var/run/secrets/workload-spiffe-credentials/ | tar -C certs -xf - --strip-components=4 --exclude='..*'

# request sends a plaintext request to the greeter-intermediary server using the DNS resolver.
# Port forwarding must be running. Verify that the port in the command below matches the local forwarded port.
.PHONY: request
request: check-grpcurl
	grpcurl -max-time 10 -proto ./greeter-java/third_party/src/main/proto/helloworld/helloworld.proto -plaintext -d '{"name": "World"}' dns:///localhost:50055 helloworld.Greeter/SayHello

# request-leaf sends a plaintext request to the greeter-leaf server using the DNS resolver.
# Port forwarding must be running. Verify that the port in the command below matches the local forwarded port.
.PHONY: request-leaf
request-leaf: check-grpcurl
	grpcurl -proto ./greeter-java/third_party/src/main/proto/helloworld/helloworld.proto -plaintext -d '{"name": "World"}' dns:///localhost:50057 helloworld.Greeter/SayHello

# request-mtls sends an mTLS request to the greeter-intermediary server using the DNS resolver.
# Port forwarding must be running. Verify that the port in the command below matches the local forwarded port.
.PHONY: request-mtls
request-mtls: check-grpcurl certs/ca_certificates.pem
	grpcurl -proto ./greeter-java/third_party/src/main/proto/helloworld/helloworld.proto -authority greeter-intermediary -cacert ./certs/ca_certificates.pem -cert ./certs/certificates.pem -key ./certs/private_key.pem -d '{"name": "World"}' dns:///localhost:50055 helloworld.Greeter/SayHello

# request-leaf-mtls sends an mTLS request to the greeter-leaf server using the DNS resolver.
# Port forwarding must be running. Verify that the port in the command below matches the local forwarded port.
.PHONY: request-leaf-mtls
request-leaf-mtls: check-grpcurl certs/ca_certificates.pem
	grpcurl -proto ./greeter-java/third_party/src/main/proto/helloworld/helloworld.proto -authority greeter-leaf -cacert ./certs/ca_certificates.pem -cert ./certs/certificates.pem -key ./certs/private_key.pem -d '{"name": "World"}' dns:///localhost:50057 helloworld.Greeter/SayHello

# request-mtls-noverify sends an mTLS request to the greeter-intermediary server using the DNS resolver,
# with certificate validation turned off. Use this target with GKE workload TLS certificates,
# as these certificates do not have common names or subject alternate names that grpcurl can verify.
# Port forwarding must be running. Verify that the port in the command below matches the local forwarded port.
.PHONY: request-mtls-noverify
request-mtls-noverify: check-grpcurl certs/ca_certificates.pem
	grpcurl -proto ./greeter-java/third_party/src/main/proto/helloworld/helloworld.proto -insecure -cacert ./certs/ca_certificates.pem -cert ./certs/certificates.pem -key ./certs/private_key.pem -d '{"name": "World"}' dns:///localhost:50055 helloworld.Greeter/SayHello

# request-leaf-mtls-noverify sends an mTLS request to the greeter-leaf server using the DNS resolver,
# with certificate validation turned off. Use this target with GKE workload TLS certificates,
# as these certificates do not have common names or subject alternate names that grpcurl can verify.
# Port forwarding must be running. Verify that the port in the command below matches the local forwarded port.
.PHONY: request-leaf-mtls-noverify
request-leaf-mtls-noverify: check-grpcurl certs/ca_certificates.pem
	grpcurl -proto ./greeter-java/third_party/src/main/proto/helloworld/helloworld.proto -insecure -cacert ./certs/ca_certificates.pem -cert ./certs/certificates.pem -key ./certs/private_key.pem -d '{"name": "World"}' dns:///localhost:50057 helloworld.Greeter/SayHello

# tail-control-plane tails the control-plane logs
.PHONY: tail-control-plane
tail-control-plane: check-kubectl
	kubectl logs --all-containers --follow --namespace=xds deployment/control-plane

# tail-greeter-intermediary tails the greeter-intermediary logs
.PHONY: tail-greeter-intermediary
tail-greeter-intermediary: check-kubectl
	kubectl logs --all-containers --follow --namespace=xds deployment/greeter-intermediary

# tail-greeter-intermediary tails the greeter-leaf logs
.PHONY: tail-greeter-leaf
tail-greeter-leaf: check-kubectl
	kubectl logs --all-containers --follow --namespace=xds deployment/greeter-leaf

# xds-status-intermediary shows the current xDS resources ACKed by greeter-intermediary
.PHONY: xds-status-intermediary
xds-status-intermediary: check-grpcdebug
	grpcdebug localhost:50055 xds status

# xds-status-leaf shows the current xDS resources ACKed by greeter-leaf
.PHONY: xds-status-leaf
xds-status-leaf: check-grpcdebug
	grpcdebug localhost:50057 xds status

# xds-config-intermediary dumps the current xDS configuration of greeter-intermediary
.PHONY: xds-config-intermediary
xds-config-intermediary: check-grpcdebug check-yq
	grpcdebug localhost:50055 xds config | yq --input-format=json --prettyPrint

# xds-config-intermediary-lds dumps the current LDS configuration of greeter-intermediary
.PHONY: xds-config-intermediary-lds
xds-config-intermediary-lds: check-grpcdebug check-yq
	grpcdebug localhost:50055 xds config --type=LDS | yq --input-format=json --prettyPrint

# xds-config-leaf dumps the current xDS configuration of greeter-leaf
.PHONY: xds-config-leaf
xds-config-leaf: check-grpcdebug check-yq
	grpcdebug localhost:50057 xds config | yq --input-format=json --prettyPrint

# xds-config-leaf-lds-rds dumps the current LDS and RDS configuration of greeter-leaf
.PHONY: xds-config-leaf-lds-rds
xds-config-leaf-lds-rds: check-grpcdebug check-yq
	grpcdebug localhost:50057 xds config --type=LDS,RDS | yq --input-format=json --prettyPrint

# envoy-config dumps the current configuration of the Envoy proxy
.PHONY: envoy-config
envoy-config: check-yq
	curl -sL localhost:19000/config_dump?include_eds | yq --input-format=json --prettyPrint

# envoy-config-xds dumps the current LDS, RDS, CDS, and EDS configuration of the Envoy proxy
.PHONY: envoy-config-xds
envoy-config-xds: check-jq check-yq
	curl -sL localhost:19000/config_dump?include_eds | jq '.configs[] | select(."@type" == "type.googleapis.com/envoy.admin.v3.ListenersConfigDump" or ."@type" == "type.googleapis.com/envoy.admin.v3.RoutesConfigDump" or ."@type" == "type.googleapis.com/envoy.admin.v3.ClustersConfigDump" or ."@type" == "type.googleapis.com/envoy.admin.v3.EndpointsConfigDump")' | yq --input-format=json --prettyPrint

# envoy-config-lds dumps the current LDS configuration of the Envoy proxy
.PHONY: envoy-config-lds
envoy-config-lds: check-yq
	curl -sL localhost:19000/config_dump?include_eds | yq --input-format=json --prettyPrint '.configs[] | select(.@type == "type.googleapis.com/envoy.admin.v3.ListenersConfigDump")'

# cert-manager deploys cert-manager and its custom resource definitions (CRDs) to the Kubernetes cluster,
# and creates a root CA with a cluster issuer.
.PHONY: cert-manager
cert-manager: check-skaffold
	skaffold run --filename=k8s/cert-manager/skaffold.yaml --module=cert-manager
	skaffold run --filename=k8s/cert-manager/skaffold.yaml --module=root-ca

# cert-manager-cas deploys cert-manager and its custom resource definitions (CRDs) to a GKE
# cluster, as well as the cert-manager issuer controller for CA Service from Google Cloud.
# It also creates a cluster issuer that references a certificate authority from CA Service.
# This target depends on a Kubernetes manifest that is created by following the Cloud Service Mesh
# setup instructions.
.PHONY: cert-manager-cas
cert-manager-cas: check-skaffold k8s/cert-manager/ca-service-cluster-issuer/ca-service-cluster-issuer.yaml
	skaffold run --filename=k8s/cert-manager/skaffold.yaml --module=cert-manager
	skaffold run --filename=k8s/cert-manager/skaffold.yaml --module=ca-service-issuer-controller
	skaffold run --filename=k8s/cert-manager/skaffold.yaml --module=ca-service-cluster-issuer

# run-bastion deploys a bastion Pod with troubleshooting tools for use with the sample xDS control
# plane implementations in this repository.
.PHONY: run-bastion
run-bastion: check-skaffold
	skaffold run --filename=k8s/troubleshoot/skaffold.yaml

# run-bastion-td deploys a bastion Pod with troubleshooting tools for use with
# Traffic Director/Cloud Service Mesh on a GKE cluster.
.PHONY: run-bastion-td
run-bastion-td: check-skaffold
	skaffold run --filename=k8s/troubleshoot/skaffold.yaml --profile=traffic-director

# run-bastion-td-off-gcp deploys a bastion Pod with troubleshooting tools for use with
# Traffic Director/Cloud Service Mesh on a non-GKE Kubernetes cluster.
.PHONY: run-bastion-td-off-gcp
run-bastion-td-off-gcp: check-skaffold
	skaffold run --filename=k8s/troubleshoot/skaffold.yaml --profile=traffic-director-off-gcp

# troubleshoot executes a shell in the bastion Pod, in raw terminal mode.
.PHONY: troubleshoot
troubleshoot: check-kubectl
	kubectl exec --container=app --namespace=xds --stdin --tty deployment/bastion -- /bin/sh

# delete-control-plane deletes the Deployments, Services, and ConfigMaps for control-plane-[go|java].
.PHONY: delete-control-plane
delete-control-plane: check-kubectl
	kubectl delete --ignore-not-found --namespace=xds deployment control-plane
	kubectl delete --ignore-not-found --namespace=xds service control-plane
	kubectl delete --ignore-not-found --namespace=xds configmaps --selector='app.kubernetes.io/part-of=grpc-xds'

# delete-greeter deletes the Deployments and Services for greeter-[intermediary|leaf].
.PHONY: delete-greeter
delete-greeter: check-kubectl
	kubectl delete --ignore-not-found --namespace=xds deployment greeter-intermediary greeter-leaf
	kubectl delete --ignore-not-found --namespace=xds service greeter-intermediary greeter-leaf

# delete deletes the Services and Deployments for control-plane and greeter-[intermediary|leaf],
# but leaves the bastion Pod running.
.PHONY: delete
delete: delete-control-plane delete-greeter

# delete-apps deletes the Deployments, Service, and ConfigMaps for control-plane and
# greeter-[intermediary|leaf] in both Kubernetes clusters, but leaves any bastion Pods running.
.PHONY: delete-apps
delete-apps: check-kubectl
	./hack/delete-apps.sh

# clean removes running instances of control-plane-[go|java], greeter-[go|java], Envoy front proxy, and bastion.
.PHONY: clean
clean: check-skaffold
	skaffold delete --module=go
	skaffold delete --module=java
	skaffold delete --filename=k8s/envoy/skaffold.yaml
	skaffold delete --filename=k8s/troubleshoot/skaffold.yaml

.PHONY: check-go
check-go:
	@if ! which go > /dev/null; then \
		echo "error: go is not installed: https://go.dev/doc/install" >&2; \
		exit 1; \
	fi

.PHONY: check-grpcdebug
check-grpcdebug:
	@if ! which grpcdebug > /dev/null; then \
		echo "error: grpcdebug is not installed: https://github.com/grpc-ecosystem/grpcdebug#installation" >&2; \
		exit 1; \
	fi

.PHONY: check-grpcurl
check-grpcurl:
	@if ! which grpcurl > /dev/null; then \
		echo "error: grpcurl is not installed: https://github.com/fullstorydev/grpcurl#installation" >&2; \
		exit 1; \
	fi

.PHONY: check-javac
check-javac:
	@if ! which javac > /dev/null; then \
		echo "error: javac is not installed: https://adoptium.net/" >&2; \
		exit 1; \
	fi

.PHONY: check-jq
check-jq:
	@if ! which jq > /dev/null; then \
		echo "error: jq is not installed: https://jqlang.github.io/jq/download/" >&2; \
		exit 1; \
	fi

.PHONY: check-kind
check-kind:
	@if ! which kind > /dev/null; then \
		echo "error: kind is not installed: https://kind.sigs.k8s.io/#installation-and-usage" >&2; \
		exit 1; \
	fi

.PHONY: check-kubectl
check-kubectl:
	@if ! which kubectl > /dev/null; then \
		echo "error: kubectl is not installed: https://kubernetes.io/docs/tasks/tools/#kubectl" >&2; \
		exit 1; \
	fi

.PHONY: check-kustomize
check-kustomize:
	@if ! which kustomize > /dev/null; then \
		echo "error: kustomize is not installed: https://kubectl.docs.kubernetes.io/installation/kustomize/" >&2; \
		exit 1; \
	fi

.PHONY: check-skaffold
check-skaffold: check-kubectl check-kustomize
	@if ! which skaffold > /dev/null; then \
		echo "error: skaffold is not installed: https://skaffold.dev/docs/install/" >&2; \
		exit 1; \
	fi

.PHONY: check-yq
check-yq:
	@if ! which yq > /dev/null; then \
		echo "error: yq is not installed: https://mikefarah.gitbook.io/yq/#install" >&2; \
		exit 1; \
	fi
