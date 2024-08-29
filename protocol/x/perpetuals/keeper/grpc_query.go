package keeper

import (
	"github.com/nemo-network/v4-chain/protocol/x/perpetuals/types"
)

var _ types.QueryServer = Keeper{}
