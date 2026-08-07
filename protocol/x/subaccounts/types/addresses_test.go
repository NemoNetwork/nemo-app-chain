package types_test

import (
	"github.com/nemo-network/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModuleAddress(t *testing.T) {
	require.Equal(t, "nemo1v88c3xv9xyv3eetdx0tvcmq7ung3dywpkux9zs", types.ModuleAddress.String())
}
