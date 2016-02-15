// Code generated by protoc-gen-go.
// source: ghs.proto
// DO NOT EDIT!

/*
Package ghs is a generated protocol buffer package.

It is generated from these files:
	ghs.proto

It has these top-level messages:
	GHSMessage
*/
package ghs

import proto "github.com/golang/protobuf/proto"
import graphalg "github.com/tcolgate/vonq/graphalg"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type GHSMessage_Type int32

const (
	GHSMessage_CONNECT    GHSMessage_Type = 0
	GHSMessage_INITIATE   GHSMessage_Type = 1
	GHSMessage_TEST       GHSMessage_Type = 2
	GHSMessage_ACCEPT     GHSMessage_Type = 3
	GHSMessage_REJECT     GHSMessage_Type = 4
	GHSMessage_REPORT     GHSMessage_Type = 5
	GHSMessage_CHANGEROOT GHSMessage_Type = 6
	GHSMessage_HALT       GHSMessage_Type = 7
)

var GHSMessage_Type_name = map[int32]string{
	0: "CONNECT",
	1: "INITIATE",
	2: "TEST",
	3: "ACCEPT",
	4: "REJECT",
	5: "REPORT",
	6: "CHANGEROOT",
	7: "HALT",
}
var GHSMessage_Type_value = map[string]int32{
	"CONNECT":    0,
	"INITIATE":   1,
	"TEST":       2,
	"ACCEPT":     3,
	"REJECT":     4,
	"REPORT":     5,
	"CHANGEROOT": 6,
	"HALT":       7,
}

func (x GHSMessage_Type) String() string {
	return proto.EnumName(GHSMessage_Type_name, int32(x))
}

type GHSMessage_Initiate_NodeState int32

const (
	GHSMessage_Initiate_sleeping GHSMessage_Initiate_NodeState = 0
	GHSMessage_Initiate_find     GHSMessage_Initiate_NodeState = 1
	GHSMessage_Initiate_found    GHSMessage_Initiate_NodeState = 2
)

var GHSMessage_Initiate_NodeState_name = map[int32]string{
	0: "sleeping",
	1: "find",
	2: "found",
}
var GHSMessage_Initiate_NodeState_value = map[string]int32{
	"sleeping": 0,
	"find":     1,
	"found":    2,
}

func (x GHSMessage_Initiate_NodeState) String() string {
	return proto.EnumName(GHSMessage_Initiate_NodeState_name, int32(x))
}

type GHSMessage struct {
	Type       GHSMessage_Type        `protobuf:"varint,1,opt,name=type,enum=ghs.GHSMessage_Type" json:"type,omitempty"`
	Connect    *GHSMessage_Connect    `protobuf:"bytes,2,opt,name=connect" json:"connect,omitempty"`
	Initiate   *GHSMessage_Initiate   `protobuf:"bytes,3,opt,name=initiate" json:"initiate,omitempty"`
	Test       *GHSMessage_Test       `protobuf:"bytes,4,opt,name=test" json:"test,omitempty"`
	Accept     *GHSMessage_Accept     `protobuf:"bytes,5,opt,name=accept" json:"accept,omitempty"`
	Reject     *GHSMessage_Reject     `protobuf:"bytes,6,opt,name=reject" json:"reject,omitempty"`
	Report     *GHSMessage_Report     `protobuf:"bytes,7,opt,name=report" json:"report,omitempty"`
	Changeroot *GHSMessage_ChangeRoot `protobuf:"bytes,8,opt,name=changeroot" json:"changeroot,omitempty"`
	Halt       *GHSMessage_Halt       `protobuf:"bytes,9,opt,name=halt" json:"halt,omitempty"`
}

func (m *GHSMessage) Reset()         { *m = GHSMessage{} }
func (m *GHSMessage) String() string { return proto.CompactTextString(m) }
func (*GHSMessage) ProtoMessage()    {}

func (m *GHSMessage) GetConnect() *GHSMessage_Connect {
	if m != nil {
		return m.Connect
	}
	return nil
}

func (m *GHSMessage) GetInitiate() *GHSMessage_Initiate {
	if m != nil {
		return m.Initiate
	}
	return nil
}

func (m *GHSMessage) GetTest() *GHSMessage_Test {
	if m != nil {
		return m.Test
	}
	return nil
}

func (m *GHSMessage) GetAccept() *GHSMessage_Accept {
	if m != nil {
		return m.Accept
	}
	return nil
}

func (m *GHSMessage) GetReject() *GHSMessage_Reject {
	if m != nil {
		return m.Reject
	}
	return nil
}

func (m *GHSMessage) GetReport() *GHSMessage_Report {
	if m != nil {
		return m.Report
	}
	return nil
}

func (m *GHSMessage) GetChangeroot() *GHSMessage_ChangeRoot {
	if m != nil {
		return m.Changeroot
	}
	return nil
}

func (m *GHSMessage) GetHalt() *GHSMessage_Halt {
	if m != nil {
		return m.Halt
	}
	return nil
}

type GHSMessage_Connect struct {
	Level uint32 `protobuf:"varint,1,opt" json:"Level,omitempty"`
}

func (m *GHSMessage_Connect) Reset()         { *m = GHSMessage_Connect{} }
func (m *GHSMessage_Connect) String() string { return proto.CompactTextString(m) }
func (*GHSMessage_Connect) ProtoMessage()    {}

type GHSMessage_Initiate struct {
	Level     uint32                        `protobuf:"varint,1,opt,name=level" json:"level,omitempty"`
	Fragment  *graphalg.Weight              `protobuf:"bytes,2,opt,name=fragment" json:"fragment,omitempty"`
	NodeState GHSMessage_Initiate_NodeState `protobuf:"varint,3,opt,name=nodeState,enum=ghs.GHSMessage_Initiate_NodeState" json:"nodeState,omitempty"`
}

func (m *GHSMessage_Initiate) Reset()         { *m = GHSMessage_Initiate{} }
func (m *GHSMessage_Initiate) String() string { return proto.CompactTextString(m) }
func (*GHSMessage_Initiate) ProtoMessage()    {}

func (m *GHSMessage_Initiate) GetFragment() *graphalg.Weight {
	if m != nil {
		return m.Fragment
	}
	return nil
}

type GHSMessage_Test struct {
	Level    uint32           `protobuf:"varint,1,opt,name=level" json:"level,omitempty"`
	Fragment *graphalg.Weight `protobuf:"bytes,2,opt,name=fragment" json:"fragment,omitempty"`
}

func (m *GHSMessage_Test) Reset()         { *m = GHSMessage_Test{} }
func (m *GHSMessage_Test) String() string { return proto.CompactTextString(m) }
func (*GHSMessage_Test) ProtoMessage()    {}

func (m *GHSMessage_Test) GetFragment() *graphalg.Weight {
	if m != nil {
		return m.Fragment
	}
	return nil
}

type GHSMessage_Accept struct {
}

func (m *GHSMessage_Accept) Reset()         { *m = GHSMessage_Accept{} }
func (m *GHSMessage_Accept) String() string { return proto.CompactTextString(m) }
func (*GHSMessage_Accept) ProtoMessage()    {}

type GHSMessage_Reject struct {
}

func (m *GHSMessage_Reject) Reset()         { *m = GHSMessage_Reject{} }
func (m *GHSMessage_Reject) String() string { return proto.CompactTextString(m) }
func (*GHSMessage_Reject) ProtoMessage()    {}

type GHSMessage_Report struct {
	Weight *graphalg.Weight `protobuf:"bytes,1,opt,name=weight" json:"weight,omitempty"`
}

func (m *GHSMessage_Report) Reset()         { *m = GHSMessage_Report{} }
func (m *GHSMessage_Report) String() string { return proto.CompactTextString(m) }
func (*GHSMessage_Report) ProtoMessage()    {}

func (m *GHSMessage_Report) GetWeight() *graphalg.Weight {
	if m != nil {
		return m.Weight
	}
	return nil
}

type GHSMessage_ChangeRoot struct {
}

func (m *GHSMessage_ChangeRoot) Reset()         { *m = GHSMessage_ChangeRoot{} }
func (m *GHSMessage_ChangeRoot) String() string { return proto.CompactTextString(m) }
func (*GHSMessage_ChangeRoot) ProtoMessage()    {}

type GHSMessage_Halt struct {
}

func (m *GHSMessage_Halt) Reset()         { *m = GHSMessage_Halt{} }
func (m *GHSMessage_Halt) String() string { return proto.CompactTextString(m) }
func (*GHSMessage_Halt) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("ghs.GHSMessage_Type", GHSMessage_Type_name, GHSMessage_Type_value)
	proto.RegisterEnum("ghs.GHSMessage_Initiate_NodeState", GHSMessage_Initiate_NodeState_name, GHSMessage_Initiate_NodeState_value)
}
