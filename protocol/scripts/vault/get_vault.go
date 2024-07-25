package main

import (
	"fmt"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	appconfig "github.com/dydxprotocol/v4-chain/protocol/app/config"
	vaultcli "github.com/dydxprotocol/v4-chain/protocol/x/vault/client/cli"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

/*
main derives a vault (in essence a subaccount) given its type and number.
This output can be used in a transaction to deposit into a vault.

Usage:

	go run scripts/vault/get_vault.go -type <vault_type> -number <vault_number>
*/
func main() {
	// ------------ FLAGS ------------
	// Get flags.
	var vaultTypeStr string = "clob"

	// Convert vault type string to VaultType.
	vaultType, err := vaultcli.GetVaultTypeFromString(vaultTypeStr)
	if err != nil {
		log.Fatal(err)
	}

	// Convert vault number string to uint32.
	vaultRange := 1000

	for i := 0; i < vaultRange; i++ {
		vaultNumber := uint32(i)
		// ------------ LOGIC ------------
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount(appconfig.Bech32PrefixAccAddr, appconfig.Bech32PrefixAccPub)

		vaultId := vaulttypes.VaultId{
			Type:   vaultType,
			Number: vaultNumber,
		}
		if err != nil {
			log.Fatal(err)
		}

		vault := vaultId.ToSubaccountId()

		// ------------ OUTPUT ------------
		fmt.Println("\"" + vault.Owner + "\"" + ",")
	}
}
