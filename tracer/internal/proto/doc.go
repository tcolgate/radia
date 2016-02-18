// Package proto describes the RPC interface to tracer, it should not
// be used directly, but via the tracer package
package proto

//go:generate protoc -I $GOPATH/src:. --go_out=plugins=grpc:. tracer.proto
