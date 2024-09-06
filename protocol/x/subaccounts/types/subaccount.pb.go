// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: nemo_network/subaccounts/subaccount.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// SubaccountId defines a unique identifier for a Subaccount.
type SubaccountId struct {
	// The address of the wallet that owns this subaccount.
	Owner string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	// The unique number of this subaccount for the owner.
	// Currently limited to 128*1000 subaccounts per owner.
	Number uint32 `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
}

func (m *SubaccountId) Reset()         { *m = SubaccountId{} }
func (m *SubaccountId) String() string { return proto.CompactTextString(m) }
func (*SubaccountId) ProtoMessage()    {}
func (*SubaccountId) Descriptor() ([]byte, []int) {
	return fileDescriptor_506647a44082e5f3, []int{0}
}
func (m *SubaccountId) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubaccountId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubaccountId.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubaccountId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubaccountId.Merge(m, src)
}
func (m *SubaccountId) XXX_Size() int {
	return m.Size()
}
func (m *SubaccountId) XXX_DiscardUnknown() {
	xxx_messageInfo_SubaccountId.DiscardUnknown(m)
}

var xxx_messageInfo_SubaccountId proto.InternalMessageInfo

func (m *SubaccountId) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *SubaccountId) GetNumber() uint32 {
	if m != nil {
		return m.Number
	}
	return 0
}

// Subaccount defines a single sub-account for a given address.
// Subaccounts are uniquely indexed by a subaccountNumber/owner pair.
type Subaccount struct {
	// The Id of the Subaccount
	Id *SubaccountId `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// All `AssetPosition`s associated with this subaccount.
	// Always sorted ascending by `asset_id`.
	AssetPositions []*AssetPosition `protobuf:"bytes,2,rep,name=asset_positions,json=assetPositions,proto3" json:"asset_positions,omitempty"`
	// All `PerpetualPosition`s associated with this subaccount.
	// Always sorted ascending by `perpetual_id.
	PerpetualPositions []*PerpetualPosition `protobuf:"bytes,3,rep,name=perpetual_positions,json=perpetualPositions,proto3" json:"perpetual_positions,omitempty"`
	// Set by the owner. If true, then margin trades can be made in this
	// subaccount.
	MarginEnabled bool `protobuf:"varint,4,opt,name=margin_enabled,json=marginEnabled,proto3" json:"margin_enabled,omitempty"`
}

func (m *Subaccount) Reset()         { *m = Subaccount{} }
func (m *Subaccount) String() string { return proto.CompactTextString(m) }
func (*Subaccount) ProtoMessage()    {}
func (*Subaccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_506647a44082e5f3, []int{1}
}
func (m *Subaccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Subaccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Subaccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Subaccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subaccount.Merge(m, src)
}
func (m *Subaccount) XXX_Size() int {
	return m.Size()
}
func (m *Subaccount) XXX_DiscardUnknown() {
	xxx_messageInfo_Subaccount.DiscardUnknown(m)
}

var xxx_messageInfo_Subaccount proto.InternalMessageInfo

func (m *Subaccount) GetId() *SubaccountId {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Subaccount) GetAssetPositions() []*AssetPosition {
	if m != nil {
		return m.AssetPositions
	}
	return nil
}

func (m *Subaccount) GetPerpetualPositions() []*PerpetualPosition {
	if m != nil {
		return m.PerpetualPositions
	}
	return nil
}

func (m *Subaccount) GetMarginEnabled() bool {
	if m != nil {
		return m.MarginEnabled
	}
	return false
}

func init() {
	proto.RegisterType((*SubaccountId)(nil), "nemo_network.subaccounts.SubaccountId")
	proto.RegisterType((*Subaccount)(nil), "nemo_network.subaccounts.Subaccount")
}

func init() {
	proto.RegisterFile("nemo_network/subaccounts/subaccount.proto", fileDescriptor_506647a44082e5f3)
}

var fileDescriptor_506647a44082e5f3 = []byte{
	// 367 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x86, 0x9b, 0xf4, 0xde, 0x72, 0xef, 0xd4, 0x56, 0x18, 0x45, 0x62, 0x17, 0x21, 0x14, 0xd4,
	0x88, 0x24, 0xc1, 0x2a, 0xee, 0x5c, 0xb4, 0xe0, 0xc2, 0x5d, 0x49, 0x41, 0x41, 0x84, 0x30, 0x49,
	0x86, 0x76, 0xb0, 0x99, 0x09, 0x33, 0x13, 0xab, 0x6f, 0xe1, 0xde, 0xd7, 0xf0, 0x21, 0x5c, 0x16,
	0x57, 0x2e, 0xa5, 0x7d, 0x11, 0x69, 0x92, 0xd6, 0xb4, 0x92, 0x5d, 0xce, 0xc9, 0xf7, 0x7f, 0x9c,
	0x39, 0x07, 0x1c, 0x53, 0x1c, 0x31, 0x8f, 0x62, 0x39, 0x61, 0xfc, 0xc1, 0x11, 0x89, 0x8f, 0x82,
	0x80, 0x25, 0x54, 0x8a, 0xc2, 0xb7, 0x1d, 0x73, 0x26, 0x19, 0xd4, 0x8a, 0xa8, 0x5d, 0x40, 0x5b,
	0xfb, 0x01, 0x13, 0x11, 0x13, 0x5e, 0xca, 0x39, 0x59, 0x91, 0x85, 0x5a, 0x56, 0xa9, 0x1f, 0x09,
	0x81, 0xa5, 0x17, 0x33, 0x41, 0x24, 0x61, 0x34, 0xc7, 0x4f, 0x4b, 0xf1, 0x18, 0xf3, 0x18, 0xcb,
	0x04, 0x8d, 0x37, 0x22, 0xed, 0x1b, 0xb0, 0x35, 0x58, 0x71, 0xd7, 0x21, 0xb4, 0xc1, 0x5f, 0x36,
	0xa1, 0x98, 0x6b, 0x8a, 0xa1, 0x98, 0xff, 0x7b, 0xda, 0xc7, 0x9b, 0xb5, 0x9b, 0x8f, 0xd4, 0x0d,
	0x43, 0x8e, 0x85, 0x18, 0x48, 0x4e, 0xe8, 0xd0, 0xcd, 0x30, 0xb8, 0x07, 0x6a, 0x34, 0x89, 0x7c,
	0xcc, 0x35, 0xd5, 0x50, 0xcc, 0x86, 0x9b, 0x57, 0xed, 0x57, 0x15, 0x80, 0x1f, 0x31, 0xbc, 0x00,
	0x2a, 0x09, 0x53, 0x67, 0xbd, 0x73, 0x68, 0x97, 0xad, 0xc2, 0x2e, 0x8e, 0xe2, 0xaa, 0x24, 0x84,
	0x7d, 0xb0, 0xbd, 0xfe, 0x52, 0xa1, 0xa9, 0x46, 0xd5, 0xac, 0x77, 0x8e, 0xca, 0x25, 0xdd, 0x45,
	0xa0, 0x9f, 0xf3, 0x6e, 0x13, 0x15, 0x4b, 0x01, 0xef, 0xc1, 0xce, 0xef, 0x65, 0x08, 0xad, 0x9a,
	0x5a, 0x4f, 0xca, 0xad, 0xfd, 0x65, 0x68, 0x65, 0x86, 0xf1, 0x66, 0x4b, 0xc0, 0x03, 0xd0, 0x8c,
	0x10, 0x1f, 0x12, 0xea, 0x61, 0x8a, 0xfc, 0x31, 0x0e, 0xb5, 0x3f, 0x86, 0x62, 0xfe, 0x73, 0x1b,
	0x59, 0xf7, 0x2a, 0x6b, 0xf6, 0x6e, 0xdf, 0x67, 0xba, 0x32, 0x9d, 0xe9, 0xca, 0xd7, 0x4c, 0x57,
	0x5e, 0xe6, 0x7a, 0x65, 0x3a, 0xd7, 0x2b, 0x9f, 0x73, 0xbd, 0x72, 0x77, 0x39, 0x24, 0x72, 0x94,
	0xf8, 0x76, 0xc0, 0x22, 0x67, 0x31, 0x8b, 0xb5, 0xbc, 0xe6, 0xe3, 0xb9, 0x15, 0x8c, 0x10, 0xa1,
	0x4e, 0x7a, 0xb7, 0x80, 0x8d, 0x9d, 0xa7, 0xb5, 0x0b, 0xcb, 0xe7, 0x18, 0x0b, 0xbf, 0x96, 0xfe,
	0x3d, 0xfb, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xa8, 0xa0, 0x58, 0x50, 0x99, 0x02, 0x00, 0x00,
}

func (m *SubaccountId) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubaccountId) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubaccountId) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Number != 0 {
		i = encodeVarintSubaccount(dAtA, i, uint64(m.Number))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintSubaccount(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Subaccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Subaccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Subaccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MarginEnabled {
		i--
		if m.MarginEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.PerpetualPositions) > 0 {
		for iNdEx := len(m.PerpetualPositions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PerpetualPositions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSubaccount(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.AssetPositions) > 0 {
		for iNdEx := len(m.AssetPositions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AssetPositions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSubaccount(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Id != nil {
		{
			size, err := m.Id.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSubaccount(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSubaccount(dAtA []byte, offset int, v uint64) int {
	offset -= sovSubaccount(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SubaccountId) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovSubaccount(uint64(l))
	}
	if m.Number != 0 {
		n += 1 + sovSubaccount(uint64(m.Number))
	}
	return n
}

func (m *Subaccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != nil {
		l = m.Id.Size()
		n += 1 + l + sovSubaccount(uint64(l))
	}
	if len(m.AssetPositions) > 0 {
		for _, e := range m.AssetPositions {
			l = e.Size()
			n += 1 + l + sovSubaccount(uint64(l))
		}
	}
	if len(m.PerpetualPositions) > 0 {
		for _, e := range m.PerpetualPositions {
			l = e.Size()
			n += 1 + l + sovSubaccount(uint64(l))
		}
	}
	if m.MarginEnabled {
		n += 2
	}
	return n
}

func sovSubaccount(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSubaccount(x uint64) (n int) {
	return sovSubaccount(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SubaccountId) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSubaccount
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
			return fmt.Errorf("proto: SubaccountId: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubaccountId: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Number", wireType)
			}
			m.Number = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Number |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSubaccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSubaccount
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
func (m *Subaccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSubaccount
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
			return fmt.Errorf("proto: Subaccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Subaccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Id == nil {
				m.Id = &SubaccountId{}
			}
			if err := m.Id.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetPositions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AssetPositions = append(m.AssetPositions, &AssetPosition{})
			if err := m.AssetPositions[len(m.AssetPositions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PerpetualPositions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PerpetualPositions = append(m.PerpetualPositions, &PerpetualPosition{})
			if err := m.PerpetualPositions[len(m.PerpetualPositions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MarginEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
			m.MarginEnabled = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipSubaccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSubaccount
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
func skipSubaccount(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSubaccount
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
					return 0, ErrIntOverflowSubaccount
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
					return 0, ErrIntOverflowSubaccount
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
				return 0, ErrInvalidLengthSubaccount
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSubaccount
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSubaccount
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSubaccount        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSubaccount          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSubaccount = fmt.Errorf("proto: unexpected end of group")
)
