// Code generated by protoc-gen-go.
// source: ghs.proto
// DO NOT EDIT!

/*
Package ghs is a generated protocol buffer package.

It is generated from these files:
	ghs.proto

It has these top-level messages:
	Message
*/
package ghs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import graphalg "github.com/tcolgate/radia/graphalg"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type Message_InitiateMsg_NodeState int32

const (
	Message_InitiateMsg_sleeping Message_InitiateMsg_NodeState = 0
	Message_InitiateMsg_find     Message_InitiateMsg_NodeState = 1
	Message_InitiateMsg_found    Message_InitiateMsg_NodeState = 2
)

var Message_InitiateMsg_NodeState_name = map[int32]string{
	0: "sleeping",
	1: "find",
	2: "found",
}
var Message_InitiateMsg_NodeState_value = map[string]int32{
	"sleeping": 0,
	"find":     1,
	"found":    2,
}

func (x Message_InitiateMsg_NodeState) String() string {
	return proto.EnumName(Message_InitiateMsg_NodeState_name, int32(x))
}
func (Message_InitiateMsg_NodeState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 1, 0}
}

type Message struct {
	// Types that are valid to be assigned to U:
	//	*Message_Connect
	//	*Message_Initiate
	//	*Message_Test
	//	*Message_Accept
	//	*Message_Reject
	//	*Message_Report
	//	*Message_ChangeRoot
	//	*Message_Halt
	U isMessage_U `protobuf_oneof:"U"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type isMessage_U interface {
	isMessage_U()
}

type Message_Connect struct {
	Connect *Message_ConnectMsg `protobuf:"bytes,2,opt,name=connect,oneof"`
}
type Message_Initiate struct {
	Initiate *Message_InitiateMsg `protobuf:"bytes,3,opt,name=initiate,oneof"`
}
type Message_Test struct {
	Test *Message_TestMsg `protobuf:"bytes,4,opt,name=test,oneof"`
}
type Message_Accept struct {
	Accept *Message_AcceptMsg `protobuf:"bytes,5,opt,name=accept,oneof"`
}
type Message_Reject struct {
	Reject *Message_RejectMsg `protobuf:"bytes,6,opt,name=reject,oneof"`
}
type Message_Report struct {
	Report *Message_ReportMsg `protobuf:"bytes,7,opt,name=report,oneof"`
}
type Message_ChangeRoot struct {
	ChangeRoot *Message_ChangeRootMsg `protobuf:"bytes,8,opt,name=changeRoot,oneof"`
}
type Message_Halt struct {
	Halt *Message_HaltMsg `protobuf:"bytes,9,opt,name=halt,oneof"`
}

func (*Message_Connect) isMessage_U()    {}
func (*Message_Initiate) isMessage_U()   {}
func (*Message_Test) isMessage_U()       {}
func (*Message_Accept) isMessage_U()     {}
func (*Message_Reject) isMessage_U()     {}
func (*Message_Report) isMessage_U()     {}
func (*Message_ChangeRoot) isMessage_U() {}
func (*Message_Halt) isMessage_U()       {}

func (m *Message) GetU() isMessage_U {
	if m != nil {
		return m.U
	}
	return nil
}

func (m *Message) GetConnect() *Message_ConnectMsg {
	if x, ok := m.GetU().(*Message_Connect); ok {
		return x.Connect
	}
	return nil
}

func (m *Message) GetInitiate() *Message_InitiateMsg {
	if x, ok := m.GetU().(*Message_Initiate); ok {
		return x.Initiate
	}
	return nil
}

func (m *Message) GetTest() *Message_TestMsg {
	if x, ok := m.GetU().(*Message_Test); ok {
		return x.Test
	}
	return nil
}

func (m *Message) GetAccept() *Message_AcceptMsg {
	if x, ok := m.GetU().(*Message_Accept); ok {
		return x.Accept
	}
	return nil
}

func (m *Message) GetReject() *Message_RejectMsg {
	if x, ok := m.GetU().(*Message_Reject); ok {
		return x.Reject
	}
	return nil
}

func (m *Message) GetReport() *Message_ReportMsg {
	if x, ok := m.GetU().(*Message_Report); ok {
		return x.Report
	}
	return nil
}

func (m *Message) GetChangeRoot() *Message_ChangeRootMsg {
	if x, ok := m.GetU().(*Message_ChangeRoot); ok {
		return x.ChangeRoot
	}
	return nil
}

func (m *Message) GetHalt() *Message_HaltMsg {
	if x, ok := m.GetU().(*Message_Halt); ok {
		return x.Halt
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Message) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Message_OneofMarshaler, _Message_OneofUnmarshaler, _Message_OneofSizer, []interface{}{
		(*Message_Connect)(nil),
		(*Message_Initiate)(nil),
		(*Message_Test)(nil),
		(*Message_Accept)(nil),
		(*Message_Reject)(nil),
		(*Message_Report)(nil),
		(*Message_ChangeRoot)(nil),
		(*Message_Halt)(nil),
	}
}

func _Message_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Message)
	// U
	switch x := m.U.(type) {
	case *Message_Connect:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Connect); err != nil {
			return err
		}
	case *Message_Initiate:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Initiate); err != nil {
			return err
		}
	case *Message_Test:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Test); err != nil {
			return err
		}
	case *Message_Accept:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Accept); err != nil {
			return err
		}
	case *Message_Reject:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Reject); err != nil {
			return err
		}
	case *Message_Report:
		b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Report); err != nil {
			return err
		}
	case *Message_ChangeRoot:
		b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChangeRoot); err != nil {
			return err
		}
	case *Message_Halt:
		b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Halt); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Message.U has unexpected type %T", x)
	}
	return nil
}

func _Message_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Message)
	switch tag {
	case 2: // U.connect
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message_ConnectMsg)
		err := b.DecodeMessage(msg)
		m.U = &Message_Connect{msg}
		return true, err
	case 3: // U.initiate
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message_InitiateMsg)
		err := b.DecodeMessage(msg)
		m.U = &Message_Initiate{msg}
		return true, err
	case 4: // U.test
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message_TestMsg)
		err := b.DecodeMessage(msg)
		m.U = &Message_Test{msg}
		return true, err
	case 5: // U.accept
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message_AcceptMsg)
		err := b.DecodeMessage(msg)
		m.U = &Message_Accept{msg}
		return true, err
	case 6: // U.reject
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message_RejectMsg)
		err := b.DecodeMessage(msg)
		m.U = &Message_Reject{msg}
		return true, err
	case 7: // U.report
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message_ReportMsg)
		err := b.DecodeMessage(msg)
		m.U = &Message_Report{msg}
		return true, err
	case 8: // U.changeRoot
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message_ChangeRootMsg)
		err := b.DecodeMessage(msg)
		m.U = &Message_ChangeRoot{msg}
		return true, err
	case 9: // U.halt
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Message_HaltMsg)
		err := b.DecodeMessage(msg)
		m.U = &Message_Halt{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Message_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Message)
	// U
	switch x := m.U.(type) {
	case *Message_Connect:
		s := proto.Size(x.Connect)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_Initiate:
		s := proto.Size(x.Initiate)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_Test:
		s := proto.Size(x.Test)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_Accept:
		s := proto.Size(x.Accept)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_Reject:
		s := proto.Size(x.Reject)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_Report:
		s := proto.Size(x.Report)
		n += proto.SizeVarint(7<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_ChangeRoot:
		s := proto.Size(x.ChangeRoot)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Message_Halt:
		s := proto.Size(x.Halt)
		n += proto.SizeVarint(9<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Message_ConnectMsg struct {
	Level uint32 `protobuf:"varint,1,opt,name=Level,json=level" json:"Level,omitempty"`
}

func (m *Message_ConnectMsg) Reset()                    { *m = Message_ConnectMsg{} }
func (m *Message_ConnectMsg) String() string            { return proto.CompactTextString(m) }
func (*Message_ConnectMsg) ProtoMessage()               {}
func (*Message_ConnectMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Message_InitiateMsg struct {
	Level     uint32                        `protobuf:"varint,1,opt,name=level" json:"level,omitempty"`
	Fragment  *graphalg.Weight              `protobuf:"bytes,2,opt,name=fragment" json:"fragment,omitempty"`
	NodeState Message_InitiateMsg_NodeState `protobuf:"varint,3,opt,name=nodeState,enum=ghs.Message_InitiateMsg_NodeState" json:"nodeState,omitempty"`
}

func (m *Message_InitiateMsg) Reset()                    { *m = Message_InitiateMsg{} }
func (m *Message_InitiateMsg) String() string            { return proto.CompactTextString(m) }
func (*Message_InitiateMsg) ProtoMessage()               {}
func (*Message_InitiateMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *Message_InitiateMsg) GetFragment() *graphalg.Weight {
	if m != nil {
		return m.Fragment
	}
	return nil
}

type Message_TestMsg struct {
	Level    uint32           `protobuf:"varint,1,opt,name=level" json:"level,omitempty"`
	Fragment *graphalg.Weight `protobuf:"bytes,2,opt,name=fragment" json:"fragment,omitempty"`
}

func (m *Message_TestMsg) Reset()                    { *m = Message_TestMsg{} }
func (m *Message_TestMsg) String() string            { return proto.CompactTextString(m) }
func (*Message_TestMsg) ProtoMessage()               {}
func (*Message_TestMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 2} }

func (m *Message_TestMsg) GetFragment() *graphalg.Weight {
	if m != nil {
		return m.Fragment
	}
	return nil
}

type Message_AcceptMsg struct {
}

func (m *Message_AcceptMsg) Reset()                    { *m = Message_AcceptMsg{} }
func (m *Message_AcceptMsg) String() string            { return proto.CompactTextString(m) }
func (*Message_AcceptMsg) ProtoMessage()               {}
func (*Message_AcceptMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 3} }

type Message_RejectMsg struct {
}

func (m *Message_RejectMsg) Reset()                    { *m = Message_RejectMsg{} }
func (m *Message_RejectMsg) String() string            { return proto.CompactTextString(m) }
func (*Message_RejectMsg) ProtoMessage()               {}
func (*Message_RejectMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 4} }

type Message_ReportMsg struct {
	Weight *graphalg.Weight `protobuf:"bytes,1,opt,name=weight" json:"weight,omitempty"`
}

func (m *Message_ReportMsg) Reset()                    { *m = Message_ReportMsg{} }
func (m *Message_ReportMsg) String() string            { return proto.CompactTextString(m) }
func (*Message_ReportMsg) ProtoMessage()               {}
func (*Message_ReportMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 5} }

func (m *Message_ReportMsg) GetWeight() *graphalg.Weight {
	if m != nil {
		return m.Weight
	}
	return nil
}

type Message_ChangeRootMsg struct {
}

func (m *Message_ChangeRootMsg) Reset()                    { *m = Message_ChangeRootMsg{} }
func (m *Message_ChangeRootMsg) String() string            { return proto.CompactTextString(m) }
func (*Message_ChangeRootMsg) ProtoMessage()               {}
func (*Message_ChangeRootMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 6} }

type Message_HaltMsg struct {
}

func (m *Message_HaltMsg) Reset()                    { *m = Message_HaltMsg{} }
func (m *Message_HaltMsg) String() string            { return proto.CompactTextString(m) }
func (*Message_HaltMsg) ProtoMessage()               {}
func (*Message_HaltMsg) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 7} }

func init() {
	proto.RegisterType((*Message)(nil), "ghs.Message")
	proto.RegisterType((*Message_ConnectMsg)(nil), "ghs.Message.ConnectMsg")
	proto.RegisterType((*Message_InitiateMsg)(nil), "ghs.Message.InitiateMsg")
	proto.RegisterType((*Message_TestMsg)(nil), "ghs.Message.TestMsg")
	proto.RegisterType((*Message_AcceptMsg)(nil), "ghs.Message.AcceptMsg")
	proto.RegisterType((*Message_RejectMsg)(nil), "ghs.Message.RejectMsg")
	proto.RegisterType((*Message_ReportMsg)(nil), "ghs.Message.ReportMsg")
	proto.RegisterType((*Message_ChangeRootMsg)(nil), "ghs.Message.ChangeRootMsg")
	proto.RegisterType((*Message_HaltMsg)(nil), "ghs.Message.HaltMsg")
	proto.RegisterEnum("ghs.Message_InitiateMsg_NodeState", Message_InitiateMsg_NodeState_name, Message_InitiateMsg_NodeState_value)
}

var fileDescriptor0 = []byte{
	// 429 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x93, 0xcb, 0x8e, 0xd3, 0x30,
	0x14, 0x86, 0x27, 0xd3, 0xe6, 0x76, 0xca, 0x40, 0x65, 0x8d, 0xc0, 0xf2, 0x0a, 0x65, 0x35, 0x42,
	0x28, 0x1d, 0xcd, 0x08, 0x56, 0x2c, 0xb8, 0x6c, 0x06, 0x89, 0xb2, 0x08, 0x20, 0xd6, 0x99, 0xd4,
	0x75, 0x82, 0x52, 0x3b, 0x24, 0x6e, 0x79, 0x40, 0xd6, 0xbc, 0x13, 0xf1, 0x25, 0x6e, 0x23, 0xc2,
	0x6e, 0x76, 0x3e, 0xfd, 0xbf, 0x2f, 0xd1, 0xf1, 0x9f, 0x42, 0xcc, 0xca, 0x2e, 0x6d, 0x5a, 0x21,
	0x05, 0x9a, 0xf5, 0x47, 0x72, 0xcd, 0x2a, 0x59, 0xee, 0xef, 0xd3, 0x42, 0xec, 0x56, 0xb2, 0x10,
	0x35, 0xcb, 0x25, 0x5d, 0x1d, 0x04, 0xff, 0xb9, 0x62, 0x6d, 0xde, 0x94, 0x79, 0xcd, 0xdc, 0xc1,
	0x68, 0xc9, 0x9f, 0x00, 0xc2, 0x35, 0xed, 0xba, 0x9c, 0x51, 0x74, 0x0b, 0x61, 0x21, 0x38, 0xa7,
	0x85, 0xc4, 0xe7, 0xcf, 0xbd, 0xab, 0xc5, 0xcd, 0xb3, 0x54, 0x3d, 0xdf, 0xc6, 0xe9, 0x07, 0x93,
	0xad, 0x3b, 0x76, 0x77, 0x96, 0x0d, 0x24, 0x7a, 0x0d, 0x51, 0xc5, 0x2b, 0x59, 0xf5, 0x6f, 0xc2,
	0x33, 0x6d, 0xe1, 0x91, 0xf5, 0xd1, 0x86, 0x46, 0x73, 0x2c, 0x7a, 0x01, 0x73, 0x49, 0x3b, 0x89,
	0xe7, 0xda, 0xb9, 0x1c, 0x39, 0x5f, 0xfb, 0xc0, 0xf0, 0x9a, 0x41, 0xd7, 0x10, 0xe4, 0x45, 0x41,
	0x1b, 0x89, 0x7d, 0x4d, 0x3f, 0x1d, 0xd1, 0xef, 0x74, 0x64, 0x78, 0xcb, 0x29, 0xa3, 0xa5, 0x3f,
	0xd4, 0x26, 0xc1, 0x84, 0x91, 0xe9, 0xc8, 0x1a, 0x86, 0x33, 0x46, 0x23, 0x5a, 0x89, 0xc3, 0x49,
	0x43, 0x45, 0xce, 0x50, 0x03, 0x7a, 0x03, 0x50, 0x94, 0x39, 0x67, 0x34, 0x13, 0x42, 0xe2, 0x48,
	0x5b, 0x64, 0x7c, 0x63, 0x2e, 0x36, 0xe6, 0x09, 0xaf, 0xf6, 0xef, 0x6b, 0x90, 0x38, 0x9e, 0xd8,
	0xff, 0xae, 0x0f, 0xec, 0xfe, 0x8a, 0x21, 0x09, 0xc0, 0xf1, 0xf2, 0xd1, 0x25, 0xf8, 0x9f, 0xe8,
	0x81, 0xd6, 0xd8, 0xeb, 0xd5, 0x8b, 0xcc, 0xaf, 0xd5, 0x40, 0x7e, 0x7b, 0xb0, 0x38, 0xb9, 0x6b,
	0x45, 0xd5, 0xff, 0x50, 0xe8, 0x25, 0x44, 0xdb, 0x36, 0x67, 0x3b, 0xca, 0x87, 0x8e, 0x97, 0xa9,
	0xfb, 0x22, 0xbe, 0xd3, 0x8a, 0x95, 0x32, 0x73, 0x04, 0x7a, 0x0b, 0x31, 0x17, 0x1b, 0xfa, 0x45,
	0x0e, 0xe5, 0x3e, 0xbe, 0x49, 0xfe, 0x57, 0x6e, 0xfa, 0x79, 0x20, 0xb3, 0xa3, 0x94, 0xa4, 0x10,
	0xbb, 0xdf, 0xd1, 0x23, 0x88, 0xba, 0x9a, 0xd2, 0xa6, 0xe2, 0x6c, 0x79, 0x86, 0x22, 0x98, 0x6f,
	0x2b, 0xbe, 0x59, 0x7a, 0x28, 0x06, 0x7f, 0x2b, 0xf6, 0xfd, 0xf1, 0x9c, 0xac, 0x21, 0xb4, 0xe5,
	0x3f, 0xc4, 0x02, 0x64, 0x01, 0xb1, 0xfb, 0x3a, 0xd4, 0xe0, 0x8a, 0x27, 0xaf, 0xd4, 0x60, 0x3b,
	0x45, 0x57, 0x10, 0xfc, 0xd2, 0xaa, 0x7e, 0xd7, 0xd4, 0x23, 0x6d, 0x4e, 0x9e, 0xc0, 0xc5, 0xa8,
	0x54, 0x12, 0x43, 0x68, 0xdb, 0x7a, 0x3f, 0x03, 0xef, 0xdb, 0x7d, 0xa0, 0xff, 0x56, 0xb7, 0x7f,
	0x03, 0x00, 0x00, 0xff, 0xff, 0xea, 0x0e, 0xbd, 0xba, 0x9a, 0x03, 0x00, 0x00,
}
