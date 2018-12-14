// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transaction.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Transaction struct {
	ClientID int32 `protobuf:"varint,1,opt,name=ClientID,proto3" json:"ClientID,omitempty"`
	TType    int32 `protobuf:"varint,2,opt,name=TType,proto3" json:"TType,omitempty"`
	// For TType specifically, the value is such:
	// Frontend -> Backend
	// [0, 1, 2, 3, 4, 5] represent [CREATE, READ, UPDATE, DELETE, READALL, PING] calls
	// Backend -> Frontend
	// [1, 2] represent [VALID, INVALID] operations (i.e. if an operation could be performed)
	// Use cases:
	// Used to notify if Review exists within database
	// Used to respond to PING call
	Reviews              []*Transaction_RBody `protobuf:"bytes,3,rep,name=Reviews,proto3" json:"Reviews,omitempty"`
	Valid                int32                `protobuf:"varint,4,opt,name=Valid,proto3" json:"Valid,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0}
}

func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (m *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(m, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetClientID() int32 {
	if m != nil {
		return m.ClientID
	}
	return 0
}

func (m *Transaction) GetTType() int32 {
	if m != nil {
		return m.TType
	}
	return 0
}

func (m *Transaction) GetReviews() []*Transaction_RBody {
	if m != nil {
		return m.Reviews
	}
	return nil
}

func (m *Transaction) GetValid() int32 {
	if m != nil {
		return m.Valid
	}
	return 0
}

// Defines review struct data
type Transaction_RBody struct {
	RID                  int32    `protobuf:"varint,1,opt,name=RID,proto3" json:"RID,omitempty"`
	Album                string   `protobuf:"bytes,2,opt,name=Album,proto3" json:"Album,omitempty"`
	Artist               string   `protobuf:"bytes,3,opt,name=Artist,proto3" json:"Artist,omitempty"`
	Rating               int32    `protobuf:"varint,4,opt,name=Rating,proto3" json:"Rating,omitempty"`
	Body                 string   `protobuf:"bytes,5,opt,name=Body,proto3" json:"Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transaction_RBody) Reset()         { *m = Transaction_RBody{} }
func (m *Transaction_RBody) String() string { return proto.CompactTextString(m) }
func (*Transaction_RBody) ProtoMessage()    {}
func (*Transaction_RBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0, 0}
}

func (m *Transaction_RBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction_RBody.Unmarshal(m, b)
}
func (m *Transaction_RBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction_RBody.Marshal(b, m, deterministic)
}
func (m *Transaction_RBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction_RBody.Merge(m, src)
}
func (m *Transaction_RBody) XXX_Size() int {
	return xxx_messageInfo_Transaction_RBody.Size(m)
}
func (m *Transaction_RBody) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction_RBody.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction_RBody proto.InternalMessageInfo

func (m *Transaction_RBody) GetRID() int32 {
	if m != nil {
		return m.RID
	}
	return 0
}

func (m *Transaction_RBody) GetAlbum() string {
	if m != nil {
		return m.Album
	}
	return ""
}

func (m *Transaction_RBody) GetArtist() string {
	if m != nil {
		return m.Artist
	}
	return ""
}

func (m *Transaction_RBody) GetRating() int32 {
	if m != nil {
		return m.Rating
	}
	return 0
}

func (m *Transaction_RBody) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func init() {
	proto.RegisterType((*Transaction)(nil), "pb.Transaction")
	proto.RegisterType((*Transaction_RBody)(nil), "pb.Transaction.RBody")
}

func init() { proto.RegisterFile("transaction.proto", fileDescriptor_2cc4e03d2c28c490) }

var fileDescriptor_2cc4e03d2c28c490 = []byte{
	// 208 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0xe9, 0x66, 0xb3, 0xea, 0xec, 0x45, 0x07, 0x95, 0xd0, 0x53, 0xf1, 0xd4, 0x53, 0x04,
	0x7d, 0x82, 0xaa, 0x17, 0xaf, 0xa1, 0x78, 0x4f, 0x6c, 0x90, 0x40, 0x4d, 0x42, 0x13, 0x95, 0x3e,
	0xbb, 0x17, 0xc9, 0xb4, 0x56, 0x6f, 0xf3, 0xcd, 0xfc, 0xf9, 0x7e, 0x08, 0x5c, 0xe4, 0x49, 0xfb,
	0xa4, 0x5f, 0xb3, 0x0b, 0x5e, 0xc6, 0x29, 0xe4, 0x80, 0xbb, 0x68, 0x6e, 0xbe, 0x2b, 0x38, 0xf6,
	0x7f, 0x17, 0xac, 0xe1, 0xf4, 0x71, 0x74, 0xd6, 0xe7, 0xe7, 0x27, 0x51, 0x35, 0x55, 0xcb, 0xd5,
	0xc6, 0x78, 0x09, 0xbc, 0xef, 0xe7, 0x68, 0xc5, 0x8e, 0x0e, 0x0b, 0xe0, 0x2d, 0x9c, 0x28, 0xfb,
	0xe9, 0xec, 0x57, 0x12, 0xac, 0x61, 0xed, 0xf1, 0xee, 0x4a, 0x46, 0x23, 0xff, 0x39, 0xa5, 0x7a,
	0x08, 0xc3, 0xac, 0x7e, 0x53, 0x45, 0xf3, 0xa2, 0x47, 0x37, 0x88, 0xfd, 0xa2, 0x21, 0xa8, 0x13,
	0x70, 0xca, 0xe1, 0x39, 0x30, 0xb5, 0x95, 0x97, 0xb1, 0x3c, 0xe8, 0x46, 0xf3, 0xf1, 0x4e, 0xbd,
	0x67, 0x6a, 0x01, 0xbc, 0x86, 0x43, 0x37, 0x65, 0x97, 0xb2, 0x60, 0xb4, 0x5e, 0xa9, 0xec, 0x95,
	0xce, 0xce, 0xbf, 0xad, 0xfe, 0x95, 0x10, 0x61, 0x5f, 0xfc, 0x82, 0x53, 0x9a, 0x66, 0x73, 0xa0,
	0x8f, 0xb8, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0x0c, 0x27, 0x31, 0x35, 0x1d, 0x01, 0x00, 0x00,
}
