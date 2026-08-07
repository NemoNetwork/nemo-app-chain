package cli

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/nemo-network/v4-chain/protocol/dtypes"
	satypes "github.com/nemo-network/v4-chain/protocol/x/subaccounts/types"
	"github.com/nemo-network/v4-chain/protocol/x/vault/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdDepositToMegavault())
	cmd.AddCommand(CmdWithdrawFromMegavault())

	return cmd
}

func CmdDepositToMegavault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-to-megavault [depositor_owner] [depositor_number] [quantums]",
		Short: "Broadcast message DepositToMegavault",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Parse depositor number.
			depositorNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			// Parse quantums.
			quantums, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Create MsgDepositToMegavault.
			msg := &types.MsgDepositToMegavault{
				SubaccountId: &satypes.SubaccountId{
					Owner:  args[0],
					Number: depositorNumber,
				},
				QuoteQuantums: dtypes.NewIntFromUint64(quantums),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdWithdrawFromMegavault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-from-megavault [withdrawer_owner] [withdrawer_number] [shares] [min_quote_quantums]",
		Short: "Broadcast message WithdrawFromMegavault",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Parse withdrawer owner and number
			withdrawerOwner := args[0]
			withdrawerNumber, err := cast.ToUint32E(args[1])
			if err != nil {
				return err
			}

			// Parse shares.
			shares, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			// Parse min quote quantums.
			minQuoteQuantums, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			// Create MsgWithdrawFromMegavault.
			msg := &types.MsgWithdrawFromMegavault{
				SubaccountId: satypes.SubaccountId{
					Owner:  withdrawerOwner,
					Number: withdrawerNumber,
				},
				Shares:           types.NumShares{NumShares: dtypes.NewIntFromUint64(shares)},
				MinQuoteQuantums: dtypes.NewIntFromUint64(minQuoteQuantums),
			}

			// Broadcast or generate the transaction.
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	// Add the necessary flags.
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
