package keeper

import (
	"github.com/nemo-network/v4-chain/protocol/x/listing/types"
)

var _ types.QueryServer = Keeper{}
