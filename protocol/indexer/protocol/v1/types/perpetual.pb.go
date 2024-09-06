// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: nemo_network/indexer/protocol/v1/perpetual.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Market type of perpetual.
// Defined in perpetual.
type PerpetualMarketType int32

const (
	// Unspecified market type.
	PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED PerpetualMarketType = 0
	// Market type for cross margin perpetual markets.
	PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS PerpetualMarketType = 1
	// Market type for isolated margin perpetual markets.
	PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED PerpetualMarketType = 2
)

var PerpetualMarketType_name = map[int32]string{
	0: "PERPETUAL_MARKET_TYPE_UNSPECIFIED",
	1: "PERPETUAL_MARKET_TYPE_CROSS",
	2: "PERPETUAL_MARKET_TYPE_ISOLATED",
}

var PerpetualMarketType_value = map[string]int32{
	"PERPETUAL_MARKET_TYPE_UNSPECIFIED": 0,
	"PERPETUAL_MARKET_TYPE_CROSS":       1,
	"PERPETUAL_MARKET_TYPE_ISOLATED":    2,
}

func (x PerpetualMarketType) String() string {
	return proto.EnumName(PerpetualMarketType_name, int32(x))
}

func (PerpetualMarketType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bba195cde109f518, []int{0}
}

func init() {
	proto.RegisterEnum("nemo_network.indexer.protocol.v1.PerpetualMarketType", PerpetualMarketType_name, PerpetualMarketType_value)
}

func init() {
	proto.RegisterFile("nemo_network/indexer/protocol/v1/perpetual.proto", fileDescriptor_bba195cde109f518)
}

var fileDescriptor_bba195cde109f518 = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0xc8, 0x4b, 0xcd, 0xcd,
	0x8f, 0xcf, 0x4b, 0x2d, 0x29, 0xcf, 0x2f, 0xca, 0xd6, 0xcf, 0xcc, 0x4b, 0x49, 0xad, 0x48, 0x2d,
	0xd2, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x2f, 0x33, 0xd4, 0x2f, 0x48, 0x2d,
	0x2a, 0x48, 0x2d, 0x29, 0x4d, 0xcc, 0xd1, 0x03, 0x8b, 0x0a, 0x29, 0x20, 0xeb, 0xd0, 0x83, 0xea,
	0xd0, 0x83, 0xe9, 0xd0, 0x2b, 0x33, 0xd4, 0x6a, 0x64, 0xe4, 0x12, 0x0e, 0x80, 0xe9, 0xf2, 0x4d,
	0x2c, 0xca, 0x4e, 0x2d, 0x09, 0xa9, 0x2c, 0x48, 0x15, 0x52, 0xe5, 0x52, 0x0c, 0x70, 0x0d, 0x0a,
	0x70, 0x0d, 0x09, 0x75, 0xf4, 0x89, 0xf7, 0x75, 0x0c, 0xf2, 0x76, 0x0d, 0x89, 0x0f, 0x89, 0x0c,
	0x70, 0x8d, 0x0f, 0xf5, 0x0b, 0x0e, 0x70, 0x75, 0xf6, 0x74, 0xf3, 0x74, 0x75, 0x11, 0x60, 0x10,
	0x92, 0xe7, 0x92, 0xc6, 0xae, 0xcc, 0x39, 0xc8, 0x3f, 0x38, 0x58, 0x80, 0x51, 0x48, 0x89, 0x4b,
	0x0e, 0xbb, 0x02, 0xcf, 0x60, 0x7f, 0x1f, 0xc7, 0x10, 0x57, 0x17, 0x01, 0x26, 0xa7, 0xd8, 0x13,
	0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86,
	0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0x72, 0x4e, 0xcf, 0x2c, 0xc9, 0x28, 0x4d,
	0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x07, 0x79, 0x45, 0x17, 0xe6, 0xf9, 0x32, 0x13, 0xdd, 0xe4, 0x8c,
	0xc4, 0xcc, 0x3c, 0x84, 0xef, 0xb1, 0x05, 0x47, 0x49, 0x65, 0x41, 0x6a, 0x71, 0x12, 0x1b, 0x58,
	0xc8, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x63, 0x73, 0xfa, 0xf6, 0x3f, 0x01, 0x00, 0x00,
}
