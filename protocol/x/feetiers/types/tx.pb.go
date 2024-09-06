// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: nemo_network/feetiers/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type.
type MsgUpdatePerpetualFeeParams struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// Defines the parameters to update. All parameters must be supplied.
	Params PerpetualFeeParams `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgUpdatePerpetualFeeParams) Reset()         { *m = MsgUpdatePerpetualFeeParams{} }
func (m *MsgUpdatePerpetualFeeParams) String() string { return proto.CompactTextString(m) }
func (*MsgUpdatePerpetualFeeParams) ProtoMessage()    {}
func (*MsgUpdatePerpetualFeeParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_a289fc262eca568a, []int{0}
}
func (m *MsgUpdatePerpetualFeeParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdatePerpetualFeeParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdatePerpetualFeeParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdatePerpetualFeeParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdatePerpetualFeeParams.Merge(m, src)
}
func (m *MsgUpdatePerpetualFeeParams) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdatePerpetualFeeParams) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdatePerpetualFeeParams.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdatePerpetualFeeParams proto.InternalMessageInfo

func (m *MsgUpdatePerpetualFeeParams) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgUpdatePerpetualFeeParams) GetParams() PerpetualFeeParams {
	if m != nil {
		return m.Params
	}
	return PerpetualFeeParams{}
}

// MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
// response type.
type MsgUpdatePerpetualFeeParamsResponse struct {
}

func (m *MsgUpdatePerpetualFeeParamsResponse) Reset()         { *m = MsgUpdatePerpetualFeeParamsResponse{} }
func (m *MsgUpdatePerpetualFeeParamsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdatePerpetualFeeParamsResponse) ProtoMessage()    {}
func (*MsgUpdatePerpetualFeeParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a289fc262eca568a, []int{1}
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdatePerpetualFeeParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdatePerpetualFeeParamsResponse.Merge(m, src)
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdatePerpetualFeeParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdatePerpetualFeeParamsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgUpdatePerpetualFeeParams)(nil), "nemo_network.feetiers.MsgUpdatePerpetualFeeParams")
	proto.RegisterType((*MsgUpdatePerpetualFeeParamsResponse)(nil), "nemo_network.feetiers.MsgUpdatePerpetualFeeParamsResponse")
}

func init() { proto.RegisterFile("nemo_network/feetiers/tx.proto", fileDescriptor_a289fc262eca568a) }

var fileDescriptor_a289fc262eca568a = []byte{
	// 342 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcb, 0x4b, 0xcd, 0xcd,
	0x8f, 0xcf, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0x4f, 0x4b, 0x4d, 0x2d, 0xc9, 0x4c, 0x2d,
	0x2a, 0xd6, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x45, 0x96, 0xd7, 0x83,
	0xc9, 0x4b, 0x49, 0x26, 0xe7, 0x17, 0xe7, 0xe6, 0x17, 0xc7, 0x83, 0x15, 0xe9, 0x43, 0x38, 0x10,
	0x1d, 0x52, 0xe2, 0x10, 0x9e, 0x7e, 0x6e, 0x71, 0xba, 0x7e, 0x99, 0x21, 0x88, 0x82, 0x4a, 0x28,
	0x61, 0xb7, 0xaa, 0x20, 0xb1, 0x28, 0x31, 0x17, 0xa6, 0x59, 0x24, 0x3d, 0x3f, 0x3d, 0x1f, 0x62,
	0x28, 0x88, 0x05, 0x11, 0x55, 0x5a, 0xc7, 0xc8, 0x25, 0xed, 0x5b, 0x9c, 0x1e, 0x5a, 0x90, 0x92,
	0x58, 0x92, 0x1a, 0x90, 0x5a, 0x54, 0x90, 0x5a, 0x52, 0x9a, 0x98, 0xe3, 0x96, 0x9a, 0x1a, 0x00,
	0xd6, 0x2b, 0x64, 0xc6, 0xc5, 0x99, 0x58, 0x5a, 0x92, 0x91, 0x5f, 0x94, 0x59, 0x52, 0x29, 0xc1,
	0xa8, 0xc0, 0xa8, 0xc1, 0xe9, 0x24, 0x71, 0x69, 0x8b, 0xae, 0x08, 0xd4, 0x5d, 0x8e, 0x29, 0x29,
	0x45, 0xa9, 0xc5, 0xc5, 0xc1, 0x25, 0x45, 0x99, 0x79, 0xe9, 0x41, 0x08, 0xa5, 0x42, 0xee, 0x5c,
	0x6c, 0x10, 0xdb, 0x25, 0x98, 0x14, 0x18, 0x35, 0xb8, 0x8d, 0x34, 0xf5, 0xb0, 0xfa, 0x56, 0x0f,
	0xd3, 0x4a, 0x27, 0x96, 0x13, 0xf7, 0xe4, 0x19, 0x82, 0xa0, 0xda, 0xad, 0xf8, 0x9a, 0x9e, 0x6f,
	0xd0, 0x42, 0x18, 0xac, 0xa4, 0xca, 0xa5, 0x8c, 0xc7, 0xbd, 0x41, 0xa9, 0xc5, 0x05, 0xf9, 0x79,
	0xc5, 0xa9, 0x46, 0x93, 0x18, 0xb9, 0x98, 0x7d, 0x8b, 0xd3, 0x85, 0xba, 0x18, 0xb9, 0x24, 0x70,
	0x7a, 0xce, 0x08, 0x87, 0xa3, 0xf0, 0x58, 0x20, 0x65, 0x45, 0xba, 0x1e, 0x98, 0xa3, 0x9c, 0x42,
	0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5,
	0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0xca, 0x2a, 0x3d, 0xb3, 0x24, 0xa3,
	0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x1f, 0x64, 0xbe, 0x2e, 0x2c, 0x2e, 0xcb, 0x4c, 0x74, 0x93,
	0x33, 0x12, 0x33, 0xf3, 0xf4, 0xc1, 0xd1, 0x95, 0x9c, 0x9f, 0xa3, 0x5f, 0x81, 0x94, 0x94, 0x2a,
	0x0b, 0x52, 0x8b, 0x93, 0xd8, 0xc0, 0x52, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x47, 0x16,
	0x6e, 0xb5, 0x70, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// UpdatePerpetualFeeParams updates the PerpetualFeeParams in state.
	UpdatePerpetualFeeParams(ctx context.Context, in *MsgUpdatePerpetualFeeParams, opts ...grpc.CallOption) (*MsgUpdatePerpetualFeeParamsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpdatePerpetualFeeParams(ctx context.Context, in *MsgUpdatePerpetualFeeParams, opts ...grpc.CallOption) (*MsgUpdatePerpetualFeeParamsResponse, error) {
	out := new(MsgUpdatePerpetualFeeParamsResponse)
	err := c.cc.Invoke(ctx, "/nemo_network.feetiers.Msg/UpdatePerpetualFeeParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// UpdatePerpetualFeeParams updates the PerpetualFeeParams in state.
	UpdatePerpetualFeeParams(context.Context, *MsgUpdatePerpetualFeeParams) (*MsgUpdatePerpetualFeeParamsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) UpdatePerpetualFeeParams(ctx context.Context, req *MsgUpdatePerpetualFeeParams) (*MsgUpdatePerpetualFeeParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePerpetualFeeParams not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_UpdatePerpetualFeeParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdatePerpetualFeeParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdatePerpetualFeeParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nemo_network.feetiers.Msg/UpdatePerpetualFeeParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdatePerpetualFeeParams(ctx, req.(*MsgUpdatePerpetualFeeParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "nemo_network.feetiers.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdatePerpetualFeeParams",
			Handler:    _Msg_UpdatePerpetualFeeParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nemo_network/feetiers/tx.proto",
}

func (m *MsgUpdatePerpetualFeeParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdatePerpetualFeeParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdatePerpetualFeeParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdatePerpetualFeeParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdatePerpetualFeeParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdatePerpetualFeeParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgUpdatePerpetualFeeParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Params.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgUpdatePerpetualFeeParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgUpdatePerpetualFeeParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgUpdatePerpetualFeeParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdatePerpetualFeeParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgUpdatePerpetualFeeParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgUpdatePerpetualFeeParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdatePerpetualFeeParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
