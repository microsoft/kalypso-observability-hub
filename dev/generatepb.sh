#!/bin/sh
# Note: This is expected to be ran from the root of this repo, e.g. '$ ./dev/generatepb.sh'
# This requires that the following dependencies are already installed and properly configured on your build machine
# protoc: https://grpc.io/docs/protoc-installation/
# protoc-gen-go v1.28.1 with: go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
# protoc-gen-go-grpc v1.2.0 with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
# This also requires that your install directory for these plugins are accessible via your $PATH - typically '~/go/bin'
protoc --proto_path=storage/api/grpc/proto --go_out=storage/api/grpc/proto --go-grpc_out=storage/api/grpc/proto --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative storage.proto
