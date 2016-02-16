// Package tracer implements a service, and a set of implementations,
// that recieve and display updates and log message from the nodes involved
// in a graphalg
//go:generate protoc -I $GOPATH/src:. --go_out=plugins=grpc:. tracer.proto
package tracer
