package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/nemo-network/v4-chain/protocol/lib"
	"github.com/nemo-network/v4-chain/protocol/x/vault/types"
)

// UpdateMegavaultParams updates the megavault-level parameters.
//
// Note: fork-local message with no upstream dydxprotocol equivalent.
func (k msgServer) UpdateMegavaultParams(
	goCtx context.Context,
	msg *types.MsgUpdateMegavaultParams,
) (*types.MsgUpdateMegavaultParamsResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	operator := k.GetOperatorParams(ctx).Operator

	// Check if authority is valid (must be a module authority or operator).
	if !k.HasAuthority(msg.Authority) && msg.Authority != operator {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidAuthority,
			"invalid authority %s",
			msg.Authority,
		)
	}

	if err := k.Keeper.SetMegavaultParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateMegavaultParamsResponse{}, nil
}
