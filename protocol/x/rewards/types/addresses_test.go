package types_test

import (
	"github.com/nemo-network/v4-chain/protocol/x/rewards/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTreasuryModuleAddress(t *testing.T) {
	require.Equal(t, "nemo16wrau2x4tsg033xfrrdpae6kxfn9kyueprnegt", types.TreasuryModuleAddress.String())
}
