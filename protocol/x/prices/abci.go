package prices

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/nemo-network/v4-chain/protocol/x/prices/types"
)

// PreBlocker executes all ABCI PreBlock logic respective to the prices module.
func PreBlocker(
	ctx sdk.Context,
	keeper types.PricesKeeper,
) {
	keeper.InitializeCurrencyPairIdCache(ctx)
}
