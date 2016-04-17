// Code generated by protoc-gen-go.
// source: graph.proto
// DO NOT EDIT!

/*
Package graph is a generated protocol buffer package.

It is generated from these files:
	graph.proto

It has these top-level messages:
	GraphID
	AlgorithmID
*/
package graph

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

// GraphID is an identified for a specific arrangemnt of nodes and edges.
// We can use this to select the specific topology we wish to pass messages
// to.
type GraphID struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *GraphID) Reset()                    { *m = GraphID{} }
func (m *GraphID) String() string            { return proto.CompactTextString(m) }
func (*GraphID) ProtoMessage()               {}
func (*GraphID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// AlgorithmID identifies a specific instance of a running algorithm. It
// should be possible to run multiple instances of any algorithms
// over a given graph at any one time
type AlgorithmID struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *AlgorithmID) Reset()                    { *m = AlgorithmID{} }
func (m *AlgorithmID) String() string            { return proto.CompactTextString(m) }
func (*AlgorithmID) ProtoMessage()               {}
func (*AlgorithmID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*GraphID)(nil), "graph.GraphID")
	proto.RegisterType((*AlgorithmID)(nil), "graph.AlgorithmID")
}

var fileDescriptor0 = []byte{
	// 84 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0x2f, 0x4a, 0x2c,
	0xc8, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x24, 0xb9, 0xd8, 0xdd,
	0x41, 0x0c, 0x4f, 0x17, 0x21, 0x3e, 0x2e, 0xa6, 0xcc, 0x14, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce,
	0x20, 0x20, 0x4b, 0x49, 0x96, 0x8b, 0xdb, 0x31, 0x27, 0x3d, 0xbf, 0x28, 0xb3, 0x24, 0x23, 0x17,
	0x53, 0x3a, 0x89, 0x0d, 0x6c, 0x8e, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xf3, 0x63, 0x75, 0x94,
	0x56, 0x00, 0x00, 0x00,
}
