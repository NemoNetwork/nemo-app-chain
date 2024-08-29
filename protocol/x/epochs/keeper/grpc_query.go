package keeper

import (
	"github.com/nemo-network/v4-chain/protocol/x/epochs/types"
)

var _ types.QueryServer = Keeper{}
