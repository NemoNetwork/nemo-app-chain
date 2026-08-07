package keeper

import (
	"context"
	"math/big"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nemo-network/v4-chain/protocol/dtypes"
	"github.com/nemo-network/v4-chain/protocol/lib"
	"github.com/nemo-network/v4-chain/protocol/x/vault/types"
)

// MegavaultFeeState returns the megavault fee accounting state along with the
// current NAV per share, so that a client can tell how far above or below the
// high-water mark the megavault currently is.
//
// Note: fork-local; see x/vault/spec/adr-001-megavault-fees.md.
func (k Keeper) MegavaultFeeState(
	goCtx context.Context,
	req *types.QueryMegavaultFeeStateRequest,
) (*types.QueryMegavaultFeeStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	navPerShare, _, _, exists, err := k.GetNavPerShare(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		// NAV per share is undefined with no shares or non-positive equity.
		// Report zero rather than failing the query.
		navPerShare = big.NewInt(0)
	}

	return &types.QueryMegavaultFeeStateResponse{
		FeeState:    k.GetFeeState(ctx),
		NavPerShare: dtypes.NewIntFromBigInt(navPerShare),
	}, nil
}
