// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: stratos/register/v1/genesis.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types1 "github.com/cosmos/cosmos-sdk/x/staking/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the register module's genesis state.
type GenesisState struct {
	Params              *Params                                `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty" yaml:"params"`
	ResourceNodes       *ResourceNodes                         `protobuf:"bytes,2,opt,name=resourceNodes,proto3" json:"resourceNodes,omitempty" yaml:"resource_nodes"`
	IndexingNodes       *IndexingNodes                         `protobuf:"bytes,3,opt,name=indexingNodes,proto3" json:"indexingNodes,omitempty" yaml:"indexing_nodes"`
	InitialUozPrice     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=initialUozPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"initialUozPrice" yaml:"initial_uoz_price"`
	TotalUnissuedPrepay github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=totalUnissuedPrepay,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"totalUnissuedPrepay" yaml:"total_unissued_prepay"`
	Slashing            []*Slashing                            `protobuf:"bytes,6,rep,name=slashing,proto3" json:"slashing,omitempty" yaml:"slashing_info"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_5bdab54ebea9e48e, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() *Params {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *GenesisState) GetResourceNodes() *ResourceNodes {
	if m != nil {
		return m.ResourceNodes
	}
	return nil
}

func (m *GenesisState) GetIndexingNodes() *IndexingNodes {
	if m != nil {
		return m.IndexingNodes
	}
	return nil
}

func (m *GenesisState) GetSlashing() []*Slashing {
	if m != nil {
		return m.Slashing
	}
	return nil
}

type GenesisIndexingNode struct {
	NetworkAddr  string                                 `protobuf:"bytes,1,opt,name=networkAddr,proto3" json:"networkAddr,omitempty" yaml:"network_address"`
	PubKey       *types.Any                             `protobuf:"bytes,2,opt,name=pubKey,proto3" json:"pubKey,omitempty" yaml:"pubkey"`
	Suspend      bool                                   `protobuf:"varint,3,opt,name=suspend,proto3" json:"suspend,omitempty" yaml:"suspend"`
	Status       types1.BondStatus                      `protobuf:"varint,4,opt,name=status,proto3,enum=cosmos.staking.v1beta1.BondStatus" json:"status,omitempty" yaml:"status"`
	Token        github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=token,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"token" yaml:"token"`
	OwnerAddress string                                 `protobuf:"bytes,6,opt,name=ownerAddress,proto3" json:"ownerAddress,omitempty" yaml:"owner_address"`
	Description  *Description                           `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty" yaml:"description",omitempty`
}

func (m *GenesisIndexingNode) Reset()         { *m = GenesisIndexingNode{} }
func (m *GenesisIndexingNode) String() string { return proto.CompactTextString(m) }
func (*GenesisIndexingNode) ProtoMessage()    {}
func (*GenesisIndexingNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_5bdab54ebea9e48e, []int{1}
}
func (m *GenesisIndexingNode) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisIndexingNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisIndexingNode.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisIndexingNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisIndexingNode.Merge(m, src)
}
func (m *GenesisIndexingNode) XXX_Size() int {
	return m.Size()
}
func (m *GenesisIndexingNode) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisIndexingNode.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisIndexingNode proto.InternalMessageInfo

func (m *GenesisIndexingNode) GetNetworkAddr() string {
	if m != nil {
		return m.NetworkAddr
	}
	return ""
}

func (m *GenesisIndexingNode) GetPubKey() *types.Any {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *GenesisIndexingNode) GetSuspend() bool {
	if m != nil {
		return m.Suspend
	}
	return false
}

func (m *GenesisIndexingNode) GetStatus() types1.BondStatus {
	if m != nil {
		return m.Status
	}
	return types1.Unspecified
}

func (m *GenesisIndexingNode) GetOwnerAddress() string {
	if m != nil {
		return m.OwnerAddress
	}
	return ""
}

func (m *GenesisIndexingNode) GetDescription() *Description {
	if m != nil {
		return m.Description
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "stratos.register.v1.GenesisState")
	proto.RegisterType((*GenesisIndexingNode)(nil), "stratos.register.v1.GenesisIndexingNode")
}

func init() { proto.RegisterFile("stratos/register/v1/genesis.proto", fileDescriptor_5bdab54ebea9e48e) }

var fileDescriptor_5bdab54ebea9e48e = []byte{
	// 725 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0x4f, 0x4f, 0xdb, 0x3c,
	0x1c, 0xc7, 0xdb, 0x07, 0x28, 0x90, 0x02, 0x8f, 0x9e, 0xd0, 0x67, 0x0a, 0x6c, 0x6b, 0x3a, 0x6b,
	0x9a, 0x98, 0x04, 0x89, 0xca, 0x76, 0x9a, 0xd0, 0x24, 0x22, 0xb4, 0x89, 0x4d, 0x43, 0x55, 0x3a,
	0x34, 0x69, 0x97, 0xc8, 0x4d, 0x4c, 0xb0, 0xda, 0xda, 0x91, 0xed, 0x00, 0xe1, 0xb4, 0x97, 0xb0,
	0xf7, 0xb0, 0xb7, 0xb0, 0x17, 0x81, 0x76, 0xe2, 0x38, 0xed, 0x10, 0x4d, 0x70, 0xdf, 0xa1, 0xaf,
	0x60, 0xaa, 0xed, 0xd0, 0x30, 0x55, 0x93, 0x38, 0xd5, 0xee, 0xef, 0xeb, 0xcf, 0xd7, 0xfe, 0xfd,
	0x89, 0xf1, 0x88, 0x0b, 0x06, 0x05, 0xe5, 0x2e, 0x43, 0x31, 0xe6, 0x02, 0x31, 0xf7, 0xa4, 0xed,
	0xc6, 0x88, 0x20, 0x8e, 0xb9, 0x93, 0x30, 0x2a, 0xa8, 0xb9, 0xaa, 0x25, 0x4e, 0x21, 0x71, 0x4e,
	0xda, 0xeb, 0x6b, 0x31, 0xa5, 0xf1, 0x00, 0xb9, 0x52, 0xd2, 0x4b, 0x8f, 0x5c, 0x48, 0x32, 0xa5,
	0x5f, 0x6f, 0xc4, 0x34, 0xa6, 0x72, 0xe9, 0x8e, 0x57, 0xfa, 0xdf, 0xb5, 0x90, 0xf2, 0x21, 0xe5,
	0x81, 0x0a, 0xa8, 0x8d, 0x0e, 0x3d, 0x56, 0x3b, 0x97, 0x0b, 0xd8, 0xc7, 0x24, 0x76, 0x4f, 0xda,
	0x3d, 0x24, 0x60, 0xbb, 0xd8, 0x6b, 0x15, 0x98, 0x76, 0xd3, 0x9b, 0x2b, 0x49, 0x0d, 0xf8, 0x35,
	0x6b, 0x2c, 0xbd, 0x56, 0x97, 0xef, 0x0a, 0x28, 0x90, 0xf9, 0xca, 0xa8, 0x25, 0x90, 0xc1, 0x21,
	0xb7, 0xaa, 0xad, 0xea, 0x46, 0x7d, 0xfb, 0xbe, 0x33, 0xe5, 0x31, 0x4e, 0x47, 0x4a, 0xbc, 0xff,
	0x46, 0xb9, 0xbd, 0x9c, 0xc1, 0xe1, 0xe0, 0x05, 0x50, 0x87, 0x80, 0xaf, 0x4f, 0x9b, 0xa1, 0xb1,
	0xcc, 0x10, 0xa7, 0x29, 0x0b, 0xd1, 0x01, 0x8d, 0x10, 0xb7, 0xfe, 0x91, 0x38, 0x30, 0x15, 0xe7,
	0x97, 0x95, 0xde, 0xda, 0x28, 0xb7, 0xff, 0x57, 0xd4, 0x02, 0x11, 0x90, 0x71, 0x04, 0xf8, 0xb7,
	0x99, 0x63, 0x13, 0x4c, 0x22, 0x74, 0x86, 0x49, 0xac, 0x4c, 0x66, 0xfe, 0x62, 0xb2, 0x5f, 0x56,
	0x96, 0x4d, 0x0a, 0xc4, 0x8d, 0xc9, 0x2d, 0xa6, 0x29, 0x8c, 0x7f, 0x31, 0xc1, 0x02, 0xc3, 0xc1,
	0x21, 0x3d, 0xef, 0x30, 0x1c, 0x22, 0x6b, 0xb6, 0x55, 0xdd, 0x58, 0xf4, 0xde, 0x5c, 0xe4, 0x76,
	0xe5, 0x47, 0x6e, 0x3f, 0x89, 0xb1, 0x38, 0x4e, 0x7b, 0x4e, 0x48, 0x87, 0xba, 0x4c, 0xfa, 0x67,
	0x8b, 0x47, 0x7d, 0x57, 0x64, 0x09, 0xe2, 0xce, 0x1e, 0x0a, 0x47, 0xb9, 0x6d, 0x15, 0x86, 0x12,
	0x17, 0xa4, 0xf4, 0x3c, 0x48, 0xc6, 0x40, 0xe0, 0xff, 0x69, 0x61, 0x7e, 0xaa, 0x1a, 0xab, 0x82,
	0x0a, 0x38, 0x38, 0x24, 0x98, 0xf3, 0x14, 0x45, 0x1d, 0x86, 0x12, 0x98, 0x59, 0x73, 0xd2, 0xfa,
	0xe0, 0x0e, 0xd6, 0xfb, 0x44, 0x8c, 0x72, 0xfb, 0x81, 0xb2, 0x96, 0xc8, 0x20, 0xd5, 0xcc, 0x20,
	0x91, 0x50, 0xe0, 0x4f, 0xb3, 0x32, 0xbb, 0xc6, 0x02, 0x1f, 0x40, 0x7e, 0x8c, 0x49, 0x6c, 0xd5,
	0x5a, 0x33, 0x1b, 0xf5, 0xed, 0x87, 0x53, 0x13, 0xdb, 0xd5, 0x22, 0xcf, 0x1a, 0xe5, 0x76, 0x43,
	0xf9, 0x14, 0x07, 0x03, 0x4c, 0x8e, 0x28, 0xf0, 0x6f, 0x40, 0xe0, 0xcb, 0xac, 0xb1, 0xaa, 0x1b,
	0xae, 0x5c, 0x10, 0x73, 0xc7, 0xa8, 0x13, 0x24, 0x4e, 0x29, 0xeb, 0xef, 0x46, 0x11, 0x93, 0xcd,
	0xb7, 0xe8, 0xad, 0x8f, 0x72, 0xfb, 0x9e, 0x02, 0xea, 0x60, 0x00, 0xa3, 0x88, 0x21, 0xce, 0x81,
	0x5f, 0x96, 0x9b, 0x1f, 0x8c, 0x5a, 0x92, 0xf6, 0xde, 0xa2, 0x4c, 0xb7, 0x59, 0xc3, 0x51, 0xd3,
	0xe6, 0x14, 0xd3, 0xe6, 0xec, 0x92, 0xcc, 0x7b, 0x5a, 0x6a, 0xd7, 0xb4, 0xd7, 0x47, 0x19, 0xf8,
	0xf6, 0x75, 0xab, 0xa1, 0x27, 0x2b, 0x64, 0x59, 0x22, 0xa8, 0xd3, 0x91, 0x18, 0x5f, 0xe3, 0xcc,
	0x4d, 0x63, 0x9e, 0xa7, 0x3c, 0x41, 0x24, 0x92, 0xbd, 0xb5, 0xe0, 0x99, 0xa3, 0xdc, 0x5e, 0xd1,
	0x6f, 0x54, 0x01, 0xe0, 0x17, 0x12, 0xf3, 0x9d, 0x51, 0xe3, 0x02, 0x8a, 0x94, 0xcb, 0x0e, 0x59,
	0xd9, 0x06, 0x8e, 0x86, 0x17, 0x83, 0xa9, 0x07, 0xd5, 0xf1, 0x28, 0x89, 0xba, 0x52, 0x59, 0x9e,
	0x21, 0x75, 0x16, 0xf8, 0x1a, 0x62, 0xbe, 0x37, 0xe6, 0x04, 0xed, 0x23, 0xa2, 0x8b, 0xfe, 0xf2,
	0xce, 0x45, 0x5f, 0x2a, 0x8a, 0xde, 0x47, 0x04, 0xf8, 0x0a, 0x66, 0xee, 0x18, 0x4b, 0xf4, 0x94,
	0x20, 0xb6, 0xab, 0x32, 0x69, 0xd5, 0x24, 0xbc, 0x54, 0x3b, 0x19, 0x9d, 0x24, 0xfa, 0x96, 0xda,
	0x8c, 0x8c, 0x7a, 0x84, 0x78, 0xc8, 0x70, 0x22, 0x30, 0x25, 0xd6, 0xbc, 0x4c, 0x77, 0x6b, 0x6a,
	0x5f, 0xec, 0x4d, 0x74, 0x5e, 0x6b, 0xd2, 0x82, 0xa5, 0xe3, 0x60, 0x93, 0x0e, 0xb1, 0x40, 0xc3,
	0x44, 0x64, 0x7e, 0x19, 0xeb, 0x1d, 0x5c, 0x5c, 0x35, 0xab, 0x97, 0x57, 0xcd, 0xea, 0xcf, 0xab,
	0x66, 0xf5, 0xf3, 0x75, 0xb3, 0x72, 0x79, 0xdd, 0xac, 0x7c, 0xbf, 0x6e, 0x56, 0x3e, 0x3e, 0x2f,
	0x3d, 0x5e, 0x9b, 0x12, 0x24, 0x8a, 0xe5, 0x56, 0x78, 0x0c, 0x31, 0x71, 0xcf, 0x26, 0x9f, 0x3c,
	0x99, 0x8e, 0x5e, 0x4d, 0xf6, 0xc1, 0xb3, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x0b, 0xd8, 0x26,
	0x1f, 0xbd, 0x05, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Slashing) > 0 {
		for iNdEx := len(m.Slashing) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Slashing[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	{
		size := m.TotalUnissuedPrepay.Size()
		i -= size
		if _, err := m.TotalUnissuedPrepay.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.InitialUozPrice.Size()
		i -= size
		if _, err := m.InitialUozPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.IndexingNodes != nil {
		{
			size, err := m.IndexingNodes.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.ResourceNodes != nil {
		{
			size, err := m.ResourceNodes.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Params != nil {
		{
			size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *GenesisIndexingNode) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisIndexingNode) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisIndexingNode) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Description != nil {
		{
			size, err := m.Description.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if len(m.OwnerAddress) > 0 {
		i -= len(m.OwnerAddress)
		copy(dAtA[i:], m.OwnerAddress)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.OwnerAddress)))
		i--
		dAtA[i] = 0x32
	}
	{
		size := m.Token.Size()
		i -= size
		if _, err := m.Token.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.Status != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Status))
		i--
		dAtA[i] = 0x20
	}
	if m.Suspend {
		i--
		if m.Suspend {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.PubKey != nil {
		{
			size, err := m.PubKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.NetworkAddr) > 0 {
		i -= len(m.NetworkAddr)
		copy(dAtA[i:], m.NetworkAddr)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.NetworkAddr)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Params != nil {
		l = m.Params.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.ResourceNodes != nil {
		l = m.ResourceNodes.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.IndexingNodes != nil {
		l = m.IndexingNodes.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = m.InitialUozPrice.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.TotalUnissuedPrepay.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Slashing) > 0 {
		for _, e := range m.Slashing {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *GenesisIndexingNode) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.NetworkAddr)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.PubKey != nil {
		l = m.PubKey.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Suspend {
		n += 2
	}
	if m.Status != 0 {
		n += 1 + sovGenesis(uint64(m.Status))
	}
	l = m.Token.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = len(m.OwnerAddress)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Description != nil {
		l = m.Description.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Params == nil {
				m.Params = &Params{}
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ResourceNodes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ResourceNodes == nil {
				m.ResourceNodes = &ResourceNodes{}
			}
			if err := m.ResourceNodes.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IndexingNodes", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.IndexingNodes == nil {
				m.IndexingNodes = &IndexingNodes{}
			}
			if err := m.IndexingNodes.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialUozPrice", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InitialUozPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalUnissuedPrepay", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalUnissuedPrepay.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Slashing", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Slashing = append(m.Slashing, &Slashing{})
			if err := m.Slashing[len(m.Slashing)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GenesisIndexingNode) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GenesisIndexingNode: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisIndexingNode: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetworkAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PubKey == nil {
				m.PubKey = &types.Any{}
			}
			if err := m.PubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Suspend", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Suspend = bool(v != 0)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= types1.BondStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Token.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OwnerAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OwnerAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Description == nil {
				m.Description = &Description{}
			}
			if err := m.Description.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
