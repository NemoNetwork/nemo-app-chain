// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: nemo_network/indexer/protocol/v1/subaccount.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_nemo_network_v4_chain_protocol_dtypes "github.com/nemo-network/v4-chain/protocol/dtypes"
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

// IndexerSubaccountId defines a unique identifier for a Subaccount.
type IndexerSubaccountId struct {
	// The address of the wallet that owns this subaccount.
	Owner string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	// < 128 Since 128 should be enough to start and it fits within
	// 1 Byte (1 Bit needed to indicate that the first byte is the last).
	Number uint32 `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
}

func (m *IndexerSubaccountId) Reset()         { *m = IndexerSubaccountId{} }
func (m *IndexerSubaccountId) String() string { return proto.CompactTextString(m) }
func (*IndexerSubaccountId) ProtoMessage()    {}
func (*IndexerSubaccountId) Descriptor() ([]byte, []int) {
	return fileDescriptor_787802078ec9eda3, []int{0}
}
func (m *IndexerSubaccountId) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IndexerSubaccountId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IndexerSubaccountId.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IndexerSubaccountId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexerSubaccountId.Merge(m, src)
}
func (m *IndexerSubaccountId) XXX_Size() int {
	return m.Size()
}
func (m *IndexerSubaccountId) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexerSubaccountId.DiscardUnknown(m)
}

var xxx_messageInfo_IndexerSubaccountId proto.InternalMessageInfo

func (m *IndexerSubaccountId) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *IndexerSubaccountId) GetNumber() uint32 {
	if m != nil {
		return m.Number
	}
	return 0
}

// IndexerPerpetualPosition are an account’s positions of a `Perpetual`.
// Therefore they hold any information needed to trade perpetuals.
type IndexerPerpetualPosition struct {
	// The `Id` of the `Perpetual`.
	PerpetualId uint32 `protobuf:"varint,1,opt,name=perpetual_id,json=perpetualId,proto3" json:"perpetual_id,omitempty"`
	// The size of the position in base quantums.
	Quantums github_com_nemo_network_v4_chain_protocol_dtypes.SerializableInt `protobuf:"bytes,2,opt,name=quantums,proto3,customtype=github.com/nemo-network/v4-chain/protocol/dtypes.SerializableInt" json:"quantums"`
	// The funding_index of the `Perpetual` the last time this position was
	// settled.
	FundingIndex github_com_nemo_network_v4_chain_protocol_dtypes.SerializableInt `protobuf:"bytes,3,opt,name=funding_index,json=fundingIndex,proto3,customtype=github.com/nemo-network/v4-chain/protocol/dtypes.SerializableInt" json:"funding_index"`
	// Amount of funding payment (in quote quantums).
	// Note: 1. this field is not cumulative.
	// 2. a positive value means funding payment was paid out and
	// a negative value means funding payment was received.
	FundingPayment github_com_nemo_network_v4_chain_protocol_dtypes.SerializableInt `protobuf:"bytes,4,opt,name=funding_payment,json=fundingPayment,proto3,customtype=github.com/nemo-network/v4-chain/protocol/dtypes.SerializableInt" json:"funding_payment"`
}

func (m *IndexerPerpetualPosition) Reset()         { *m = IndexerPerpetualPosition{} }
func (m *IndexerPerpetualPosition) String() string { return proto.CompactTextString(m) }
func (*IndexerPerpetualPosition) ProtoMessage()    {}
func (*IndexerPerpetualPosition) Descriptor() ([]byte, []int) {
	return fileDescriptor_787802078ec9eda3, []int{1}
}
func (m *IndexerPerpetualPosition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IndexerPerpetualPosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IndexerPerpetualPosition.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IndexerPerpetualPosition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexerPerpetualPosition.Merge(m, src)
}
func (m *IndexerPerpetualPosition) XXX_Size() int {
	return m.Size()
}
func (m *IndexerPerpetualPosition) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexerPerpetualPosition.DiscardUnknown(m)
}

var xxx_messageInfo_IndexerPerpetualPosition proto.InternalMessageInfo

func (m *IndexerPerpetualPosition) GetPerpetualId() uint32 {
	if m != nil {
		return m.PerpetualId
	}
	return 0
}

// IndexerAssetPosition define an account’s positions of an `Asset`.
// Therefore they hold any information needed to trade on Spot and Margin.
type IndexerAssetPosition struct {
	// The `Id` of the `Asset`.
	AssetId uint32 `protobuf:"varint,1,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
	// The absolute size of the position in base quantums.
	Quantums github_com_nemo_network_v4_chain_protocol_dtypes.SerializableInt `protobuf:"bytes,2,opt,name=quantums,proto3,customtype=github.com/nemo-network/v4-chain/protocol/dtypes.SerializableInt" json:"quantums"`
	// The `Index` (either `LongIndex` or `ShortIndex`) of the `Asset` the last
	// time this position was settled
	// TODO(DEC-582): pending margin trading being added.
	Index uint64 `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

func (m *IndexerAssetPosition) Reset()         { *m = IndexerAssetPosition{} }
func (m *IndexerAssetPosition) String() string { return proto.CompactTextString(m) }
func (*IndexerAssetPosition) ProtoMessage()    {}
func (*IndexerAssetPosition) Descriptor() ([]byte, []int) {
	return fileDescriptor_787802078ec9eda3, []int{2}
}
func (m *IndexerAssetPosition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IndexerAssetPosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IndexerAssetPosition.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IndexerAssetPosition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexerAssetPosition.Merge(m, src)
}
func (m *IndexerAssetPosition) XXX_Size() int {
	return m.Size()
}
func (m *IndexerAssetPosition) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexerAssetPosition.DiscardUnknown(m)
}

var xxx_messageInfo_IndexerAssetPosition proto.InternalMessageInfo

func (m *IndexerAssetPosition) GetAssetId() uint32 {
	if m != nil {
		return m.AssetId
	}
	return 0
}

func (m *IndexerAssetPosition) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func init() {
	proto.RegisterType((*IndexerSubaccountId)(nil), "nemo_network.indexer.protocol.v1.IndexerSubaccountId")
	proto.RegisterType((*IndexerPerpetualPosition)(nil), "nemo_network.indexer.protocol.v1.IndexerPerpetualPosition")
	proto.RegisterType((*IndexerAssetPosition)(nil), "nemo_network.indexer.protocol.v1.IndexerAssetPosition")
}

func init() {
	proto.RegisterFile("nemo_network/indexer/protocol/v1/subaccount.proto", fileDescriptor_787802078ec9eda3)
}

var fileDescriptor_787802078ec9eda3 = []byte{
	// 417 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x93, 0xbd, 0xae, 0xd3, 0x30,
	0x18, 0x86, 0x13, 0xce, 0x0f, 0x07, 0xd3, 0x82, 0x14, 0x22, 0x94, 0x73, 0x86, 0x9c, 0xd2, 0xe9,
	0x2c, 0x4d, 0x54, 0xc1, 0x05, 0x70, 0xca, 0x42, 0xb6, 0x2a, 0xdd, 0x90, 0xaa, 0xc8, 0x89, 0x4d,
	0x6a, 0x91, 0xd8, 0xa9, 0x7f, 0x5a, 0xca, 0xce, 0xce, 0x8d, 0xb0, 0x71, 0x11, 0x1d, 0x2b, 0x26,
	0xc4, 0x50, 0xa1, 0xf6, 0x46, 0x50, 0xed, 0x24, 0x65, 0x60, 0x60, 0xa8, 0xd8, 0xf2, 0x3d, 0xf2,
	0xf7, 0xbe, 0xd1, 0x63, 0x19, 0x0c, 0x29, 0x2e, 0x59, 0x42, 0xb1, 0x5c, 0x32, 0xfe, 0x21, 0x24,
	0x14, 0xe1, 0x8f, 0x98, 0x87, 0x15, 0x67, 0x92, 0x65, 0xac, 0x08, 0x17, 0xc3, 0x50, 0xa8, 0x14,
	0x66, 0x19, 0x53, 0x54, 0x06, 0x1a, 0x3b, 0xbd, 0x3f, 0x57, 0x82, 0x7a, 0x25, 0x68, 0x56, 0x82,
	0xc5, 0xf0, 0xe6, 0x3a, 0x63, 0xa2, 0x64, 0x22, 0xd1, 0x2c, 0x34, 0x83, 0x39, 0x70, 0xe3, 0xe6,
	0x2c, 0x67, 0x86, 0x1f, 0xbe, 0x0c, 0xed, 0x4f, 0xc1, 0xb3, 0xc8, 0xe4, 0x4c, 0xda, 0xb6, 0x08,
	0x39, 0x01, 0xb8, 0x60, 0x4b, 0x8a, 0xb9, 0x67, 0xf7, 0xec, 0xbb, 0x47, 0x23, 0xef, 0xfb, 0xb7,
	0x81, 0x5b, 0xa7, 0xdd, 0x23, 0xc4, 0xb1, 0x10, 0x13, 0xc9, 0x09, 0xcd, 0x63, 0x73, 0xcc, 0x79,
	0x0e, 0x2e, 0xa9, 0x2a, 0x53, 0xcc, 0xbd, 0x07, 0x3d, 0xfb, 0xae, 0x1b, 0xd7, 0x53, 0xff, 0xf3,
	0x19, 0xf0, 0xea, 0xfc, 0x31, 0xe6, 0x15, 0x96, 0x0a, 0x16, 0x63, 0x26, 0x88, 0x24, 0x8c, 0x3a,
	0x2f, 0x40, 0xa7, 0x6a, 0x60, 0x42, 0x90, 0xee, 0xea, 0xc6, 0x8f, 0x5b, 0x16, 0x21, 0x07, 0x81,
	0xab, 0xb9, 0x82, 0x54, 0xaa, 0x52, 0xe8, 0xe4, 0xce, 0xe8, 0xed, 0x7a, 0x7b, 0x6b, 0xfd, 0xdc,
	0xde, 0xbe, 0xce, 0x89, 0x9c, 0xa9, 0x34, 0xc8, 0x58, 0x19, 0x1e, 0xb4, 0x0c, 0x1a, 0x93, 0x8b,
	0x57, 0x83, 0x6c, 0x06, 0x09, 0x3d, 0xaa, 0x44, 0x72, 0x55, 0x61, 0x11, 0x4c, 0x30, 0x27, 0xb0,
	0x20, 0x9f, 0x60, 0x5a, 0xe0, 0x88, 0xca, 0xb8, 0x4d, 0x76, 0x4a, 0xd0, 0x7d, 0xaf, 0x28, 0x22,
	0x34, 0x4f, 0xb4, 0x54, 0xef, 0xec, 0xc4, 0x55, 0x9d, 0x3a, 0x5e, 0xab, 0x70, 0xe6, 0xe0, 0x69,
	0x53, 0x57, 0xc1, 0x55, 0x89, 0xa9, 0xf4, 0xce, 0x4f, 0x5c, 0xf8, 0xa4, 0x2e, 0x18, 0x9b, 0xfc,
	0xfe, 0x57, 0x1b, 0xb8, 0xf5, 0x3d, 0xdc, 0x0b, 0x81, 0x65, 0x7b, 0x07, 0xd7, 0xe0, 0x0a, 0x1e,
	0xc0, 0xd1, 0xff, 0x43, 0x3d, 0xff, 0x37, 0xf7, 0x2e, 0xb8, 0x38, 0x3a, 0x3f, 0x8f, 0xcd, 0x30,
	0x9a, 0xae, 0x77, 0xbe, 0xbd, 0xd9, 0xf9, 0xf6, 0xaf, 0x9d, 0x6f, 0x7f, 0xd9, 0xfb, 0xd6, 0x66,
	0xef, 0x5b, 0x3f, 0xf6, 0xbe, 0xf5, 0xee, 0xcd, 0xbf, 0x77, 0xff, 0xed, 0x4d, 0xe9, 0xdf, 0x49,
	0x2f, 0x35, 0x7a, 0xf9, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x26, 0x35, 0xed, 0x7a, 0x84, 0x03, 0x00,
	0x00,
}

func (m *IndexerSubaccountId) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IndexerSubaccountId) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerSubaccountId) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *IndexerPerpetualPosition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IndexerPerpetualPosition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerPerpetualPosition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.FundingPayment.Size()
		i -= size
		if _, err := m.FundingPayment.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSubaccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.FundingIndex.Size()
		i -= size
		if _, err := m.FundingIndex.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSubaccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Quantums.Size()
		i -= size
		if _, err := m.Quantums.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSubaccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.PerpetualId != 0 {
		i = encodeVarintSubaccount(dAtA, i, uint64(m.PerpetualId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *IndexerAssetPosition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IndexerAssetPosition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerAssetPosition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Index != 0 {
		i = encodeVarintSubaccount(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.Quantums.Size()
		i -= size
		if _, err := m.Quantums.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSubaccount(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.AssetId != 0 {
		i = encodeVarintSubaccount(dAtA, i, uint64(m.AssetId))
		i--
		dAtA[i] = 0x8
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
func (m *IndexerSubaccountId) Size() (n int) {
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

func (m *IndexerPerpetualPosition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PerpetualId != 0 {
		n += 1 + sovSubaccount(uint64(m.PerpetualId))
	}
	l = m.Quantums.Size()
	n += 1 + l + sovSubaccount(uint64(l))
	l = m.FundingIndex.Size()
	n += 1 + l + sovSubaccount(uint64(l))
	l = m.FundingPayment.Size()
	n += 1 + l + sovSubaccount(uint64(l))
	return n
}

func (m *IndexerAssetPosition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AssetId != 0 {
		n += 1 + sovSubaccount(uint64(m.AssetId))
	}
	l = m.Quantums.Size()
	n += 1 + l + sovSubaccount(uint64(l))
	if m.Index != 0 {
		n += 1 + sovSubaccount(uint64(m.Index))
	}
	return n
}

func sovSubaccount(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSubaccount(x uint64) (n int) {
	return sovSubaccount(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *IndexerSubaccountId) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: IndexerSubaccountId: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IndexerSubaccountId: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *IndexerPerpetualPosition) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: IndexerPerpetualPosition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IndexerPerpetualPosition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PerpetualId", wireType)
			}
			m.PerpetualId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PerpetualId |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quantums", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Quantums.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FundingIndex", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.FundingIndex.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FundingPayment", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.FundingPayment.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
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
func (m *IndexerAssetPosition) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: IndexerAssetPosition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IndexerAssetPosition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetId", wireType)
			}
			m.AssetId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AssetId |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Quantums", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Quantums.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint64(b&0x7F) << shift
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
