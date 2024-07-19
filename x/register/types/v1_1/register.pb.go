// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: stratos/register/v1_1/register.proto

package v1_1

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type Description struct {
	Moniker         string `protobuf:"bytes,1,opt,name=moniker,proto3" json:"moniker" yaml:"moniker"`
	Identity        string `protobuf:"bytes,2,opt,name=identity,proto3" json:"identity" yaml:"identity",omitempty`
	Website         string `protobuf:"bytes,3,opt,name=website,proto3" json:"website" yaml:"website",omitempty`
	SecurityContact string `protobuf:"bytes,4,opt,name=security_contact,json=securityContact,proto3" json:"security_contact" yaml:"security_contact",omitempty`
	Details         string `protobuf:"bytes,5,opt,name=details,proto3" json:"details" yaml:"details",omitempty`
}

func (m *Description) Reset()         { *m = Description{} }
func (m *Description) String() string { return proto.CompactTextString(m) }
func (*Description) ProtoMessage()    {}
func (*Description) Descriptor() ([]byte, []int) {
	return fileDescriptor_9b740b7b32b37484, []int{0}
}
func (m *Description) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Description) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Description.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Description) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Description.Merge(m, src)
}
func (m *Description) XXX_Size() int {
	return m.Size()
}
func (m *Description) XXX_DiscardUnknown() {
	xxx_messageInfo_Description.DiscardUnknown(m)
}

var xxx_messageInfo_Description proto.InternalMessageInfo

func (m *Description) GetMoniker() string {
	if m != nil {
		return m.Moniker
	}
	return ""
}

func (m *Description) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

func (m *Description) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *Description) GetSecurityContact() string {
	if m != nil {
		return m.SecurityContact
	}
	return ""
}

func (m *Description) GetDetails() string {
	if m != nil {
		return m.Details
	}
	return ""
}

func init() {
	proto.RegisterType((*Description)(nil), "stratos.register.v1_1.Description")
}

func init() {
	proto.RegisterFile("stratos/register/v1_1/register.proto", fileDescriptor_9b740b7b32b37484)
}

var fileDescriptor_9b740b7b32b37484 = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0xcf, 0x4e, 0xdb, 0x40,
	0x10, 0xc6, 0xe3, 0x34, 0x6d, 0x5a, 0x57, 0xea, 0x1f, 0xab, 0x95, 0x9c, 0x48, 0xf5, 0xb6, 0x56,
	0xa5, 0xb6, 0x52, 0xc9, 0x2a, 0xe2, 0x80, 0x92, 0x1b, 0x01, 0x5e, 0x20, 0x17, 0x24, 0x2e, 0x91,
	0xed, 0x2c, 0xce, 0x8a, 0x78, 0xd7, 0xf2, 0x4e, 0x00, 0xbf, 0x05, 0x47, 0x8e, 0x1c, 0x79, 0x14,
	0x8e, 0x39, 0x72, 0xb2, 0x90, 0x73, 0xcb, 0xd1, 0x2f, 0x00, 0xf2, 0xee, 0xda, 0x09, 0x70, 0xb1,
	0x66, 0x7e, 0xdf, 0x37, 0x9f, 0x47, 0xda, 0x31, 0x7f, 0x0b, 0x48, 0x3c, 0xe0, 0x02, 0x27, 0x24,
	0xa4, 0x02, 0x48, 0x82, 0xcf, 0xfb, 0x93, 0x7e, 0xdd, 0xf5, 0xe2, 0x84, 0x03, 0xb7, 0xbe, 0x6b,
	0x57, 0xaf, 0xe6, 0xa5, 0xab, 0xfb, 0x2d, 0xe4, 0x21, 0x97, 0x0e, 0x5c, 0x56, 0xca, 0xdc, 0xed,
	0x84, 0x9c, 0x87, 0x73, 0x82, 0x65, 0xe7, 0x2f, 0x4e, 0xb1, 0xc7, 0xd2, 0x4a, 0x0a, 0xb8, 0x88,
	0xb8, 0x98, 0xa8, 0x19, 0xd5, 0x68, 0xe9, 0xab, 0x17, 0x51, 0xc6, 0xb1, 0xfc, 0x2a, 0xe4, 0x3e,
	0x36, 0xcd, 0x8f, 0x87, 0x44, 0x04, 0x09, 0x8d, 0x81, 0x72, 0x66, 0xed, 0x99, 0xed, 0x88, 0x33,
	0x7a, 0x46, 0x12, 0xdb, 0xf8, 0x69, 0xfc, 0xfd, 0x30, 0xfa, 0xb1, 0xce, 0x50, 0x85, 0x8a, 0x0c,
	0x7d, 0x4a, 0xbd, 0x68, 0x3e, 0x74, 0x35, 0x70, 0xc7, 0x95, 0x64, 0x1d, 0x99, 0xef, 0xe9, 0x94,
	0x30, 0xa0, 0x90, 0xda, 0x4d, 0x39, 0xf9, 0x6f, 0x9d, 0xa1, 0x9a, 0x15, 0x19, 0xea, 0xa8, 0xd1,
	0x8a, 0xb8, 0xff, 0x79, 0x44, 0x81, 0x44, 0x31, 0xa4, 0xe3, 0xda, 0x66, 0xed, 0x9b, 0xed, 0x0b,
	0xe2, 0x0b, 0x0a, 0xc4, 0x7e, 0x23, 0x53, 0xfe, 0x94, 0xff, 0xd7, 0xa8, 0xc8, 0x90, 0xad, 0x42,
	0x34, 0xd8, 0xce, 0xa8, 0x4c, 0xd6, 0xd4, 0xfc, 0x22, 0x48, 0xb0, 0x48, 0x28, 0xa4, 0x93, 0x80,
	0x33, 0xf0, 0x02, 0xb0, 0x5b, 0x32, 0x6b, 0xb0, 0xce, 0xd0, 0x2b, 0xad, 0xc8, 0xd0, 0x2f, 0x15,
	0xfa, 0x52, 0xd9, 0x4e, 0xff, 0x5c, 0x89, 0x07, 0x4a, 0x2b, 0x17, 0x9d, 0x12, 0xf0, 0xe8, 0x5c,
	0xd8, 0x6f, 0x37, 0x8b, 0x6a, 0xb4, 0x59, 0x54, 0x83, 0x67, 0x8b, 0x6a, 0x36, 0x6c, 0x5d, 0xdf,
	0x20, 0x63, 0x74, 0x7c, 0x9b, 0x3b, 0xc6, 0x5d, 0xee, 0x18, 0xcb, 0xdc, 0x31, 0x1e, 0x72, 0xc7,
	0xb8, 0x5a, 0x39, 0x8d, 0xe5, 0xca, 0x69, 0xdc, 0xaf, 0x9c, 0xc6, 0xc9, 0x20, 0xa4, 0x30, 0x5b,
	0xf8, 0xbd, 0x80, 0x47, 0x58, 0x1f, 0x08, 0x23, 0x50, 0x95, 0x3b, 0xc1, 0xcc, 0xa3, 0x0c, 0x5f,
	0x6e, 0x2e, 0x0b, 0xd2, 0x98, 0x08, 0x79, 0x5f, 0xfe, 0x3b, 0xf9, 0xc2, 0xbb, 0x4f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xd4, 0xbe, 0x94, 0x11, 0x7f, 0x02, 0x00, 0x00,
}

func (this *Description) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Description)
	if !ok {
		that2, ok := that.(Description)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Moniker != that1.Moniker {
		return false
	}
	if this.Identity != that1.Identity {
		return false
	}
	if this.Website != that1.Website {
		return false
	}
	if this.SecurityContact != that1.SecurityContact {
		return false
	}
	if this.Details != that1.Details {
		return false
	}
	return true
}
func (m *Description) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Description) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Description) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Details) > 0 {
		i -= len(m.Details)
		copy(dAtA[i:], m.Details)
		i = encodeVarintRegister(dAtA, i, uint64(len(m.Details)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.SecurityContact) > 0 {
		i -= len(m.SecurityContact)
		copy(dAtA[i:], m.SecurityContact)
		i = encodeVarintRegister(dAtA, i, uint64(len(m.SecurityContact)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Website) > 0 {
		i -= len(m.Website)
		copy(dAtA[i:], m.Website)
		i = encodeVarintRegister(dAtA, i, uint64(len(m.Website)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Identity) > 0 {
		i -= len(m.Identity)
		copy(dAtA[i:], m.Identity)
		i = encodeVarintRegister(dAtA, i, uint64(len(m.Identity)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Moniker) > 0 {
		i -= len(m.Moniker)
		copy(dAtA[i:], m.Moniker)
		i = encodeVarintRegister(dAtA, i, uint64(len(m.Moniker)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRegister(dAtA []byte, offset int, v uint64) int {
	offset -= sovRegister(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Description) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Moniker)
	if l > 0 {
		n += 1 + l + sovRegister(uint64(l))
	}
	l = len(m.Identity)
	if l > 0 {
		n += 1 + l + sovRegister(uint64(l))
	}
	l = len(m.Website)
	if l > 0 {
		n += 1 + l + sovRegister(uint64(l))
	}
	l = len(m.SecurityContact)
	if l > 0 {
		n += 1 + l + sovRegister(uint64(l))
	}
	l = len(m.Details)
	if l > 0 {
		n += 1 + l + sovRegister(uint64(l))
	}
	return n
}

func sovRegister(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRegister(x uint64) (n int) {
	return sovRegister(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Description) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRegister
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
			return fmt.Errorf("proto: Description: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Description: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Moniker", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegister
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
				return ErrInvalidLengthRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRegister
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Moniker = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Identity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegister
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
				return ErrInvalidLengthRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRegister
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Identity = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Website", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegister
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
				return ErrInvalidLengthRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRegister
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Website = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SecurityContact", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegister
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
				return ErrInvalidLengthRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRegister
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SecurityContact = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Details", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRegister
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
				return ErrInvalidLengthRegister
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRegister
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Details = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRegister(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRegister
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
func skipRegister(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRegister
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
					return 0, ErrIntOverflowRegister
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
					return 0, ErrIntOverflowRegister
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
				return 0, ErrInvalidLengthRegister
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRegister
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRegister
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRegister        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRegister          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRegister = fmt.Errorf("proto: unexpected end of group")
)
