package keeper

import (
	"github.com/nemo-network/v4-chain/protocol/x/vest/types"
)

var _ types.QueryServer = Keeper{}
