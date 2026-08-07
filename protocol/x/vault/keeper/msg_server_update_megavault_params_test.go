package keeper_test

import (
	"testing"

	"github.com/nemo-network/v4-chain/protocol/dtypes"
	"github.com/nemo-network/v4-chain/protocol/lib"

	testapp "github.com/nemo-network/v4-chain/protocol/testutil/app"
	"github.com/nemo-network/v4-chain/protocol/testutil/constants"

	"github.com/nemo-network/v4-chain/protocol/x/vault/keeper"
	"github.com/nemo-network/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateMegavaultParams(t *testing.T) {
	operator := constants.AliceAccAddress.String()

	tests := map[string]struct {
		// Msg.
		msg *types.MsgUpdateMegavaultParams
		// Expected error.
		expectedErr string
	}{
		"Success - Gov authority sets a cap": {
			msg: &types.MsgUpdateMegavaultParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.MegavaultParams{
					DepositCapQuoteQuantums: dtypes.NewInt(1_000_000_000),
				},
			},
		},
		"Success - Operator sets a cap": {
			msg: &types.MsgUpdateMegavaultParams{
				Authority: operator,
				Params: types.MegavaultParams{
					DepositCapQuoteQuantums: dtypes.NewInt(7_654_321),
				},
			},
		},
		"Success - Cap of zero means uncapped": {
			msg: &types.MsgUpdateMegavaultParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.MegavaultParams{
					DepositCapQuoteQuantums: dtypes.NewInt(0),
				},
			},
		},
		"Failure - Neither authority nor operator": {
			msg: &types.MsgUpdateMegavaultParams{
				Authority: constants.BobAccAddress.String(),
				Params: types.MegavaultParams{
					DepositCapQuoteQuantums: dtypes.NewInt(1_000),
				},
			},
			expectedErr: types.ErrInvalidAuthority.Error(),
		},
		"Failure - Empty authority": {
			msg: &types.MsgUpdateMegavaultParams{
				Authority: "",
				Params: types.MegavaultParams{
					DepositCapQuoteQuantums: dtypes.NewInt(1_000),
				},
			},
			expectedErr: types.ErrInvalidAuthority.Error(),
		},
		"Failure - Negative deposit cap": {
			msg: &types.MsgUpdateMegavaultParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.MegavaultParams{
					DepositCapQuoteQuantums: dtypes.NewInt(-1),
				},
			},
			expectedErr: types.ErrNegativeDepositCap.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper
			ms := keeper.NewMsgServerImpl(k)

			// Appoint an operator that is not a module authority, so that the
			// operator-gating half of the check is actually exercised.
			require.NoError(t, k.SetOperatorParams(ctx, types.OperatorParams{
				Operator: operator,
			}))
			paramsBefore := k.GetMegavaultParams(ctx)

			_, err := ms.UpdateMegavaultParams(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				require.Equal(t, paramsBefore, k.GetMegavaultParams(ctx))
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.Params, k.GetMegavaultParams(ctx))
			}
		})
	}
}

func TestMegavaultParamsDefaultsToUncapped(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	require.Equal(t, types.DefaultMegavaultParams(), k.GetMegavaultParams(ctx))
	require.Equal(t, 0, k.GetMegavaultParams(ctx).DepositCapQuoteQuantums.Sign())
}
