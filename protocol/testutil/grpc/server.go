package grpc

import pricetypes "github.com/nemo-network/v4-chain/protocol/x/prices/types"

type QueryServer interface {
	pricetypes.QueryServer
}
