package keeper

import (
	"github.com/nemo-network/v4-chain/protocol/x/delaymsg/types"
)

var _ types.QueryServer = Keeper{}
