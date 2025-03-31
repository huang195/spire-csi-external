module github.com/huang195/spire-csi

go 1.23.0

toolchain go1.23.2

require (
	github.com/container-storage-interface/spec v1.10.0
	github.com/go-logr/logr v1.4.2
	github.com/go-logr/zapr v1.3.0
	github.com/spiffe/spiffe-csi v0.2.6
	github.com/spiffe/spire v1.10.4
	go.uber.org/zap v1.27.0
	golang.org/x/sys v0.25.0
	google.golang.org/grpc v1.65.0
	k8s.io/apimachinery v0.31.1
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240730163845-b1a4ccb954bf // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
