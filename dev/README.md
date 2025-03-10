# Kalypso Observability Hub - Developer Tooling

Developer tooling to assist and accelerate development efforts on the Kalypso Observability Hub

## Contents

### Generate Protobuf

Generate Protobuf files for the Storage API [generatepb.sh](./generatepb.sh)

To generate the storage.pb.go and storage_grpc.pb.go files in the `storage/api/grpc/proto` directory, execute the following from the root of this repo:

```sh
 ./dev/generatepb.sh
 ```

This script requires that the following dependencies are already installed and properly configured on your build machine

* [protoc](https://grpc.io/docs/protoc-installation/)
* protoc-gen-go v1.28.1 via `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1`
* protoc-gen-go-grpc v1.2.0 via `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0`

> Note: For the above plugins to work, it requires that your go plugin directory is acceissible via your `$PATH` variable - typically `~/go/bin`
