package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/nemo-network/v4-chain/protocol/lib"
	"github.com/nemo-network/v4-chain/protocol/x/bridge/types"
)

// UpdateSafetyParams updates the SafetyParams in state.
func (k msgServer) UpdateSafetyParams(
	goCtx context.Context,
	msg *types.MsgUpdateSafetyParams,
) (*types.MsgUpdateSafetyParamsResponse, error) {
	if !k.Keeper.HasAuthority(msg.GetAuthority()) {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidAuthority,
			"message authority %s is not valid for sending update safety params messages",
			msg.GetAuthority(),
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if err := k.Keeper.UpdateSafetyParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateSafetyParamsResponse{}, nil
}
