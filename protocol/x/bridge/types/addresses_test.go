package types_test

import (
	"github.com/nemo-network/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModuleAddress(t *testing.T) {
	require.Equal(t, "nemo1zlefkpe3g0vvm9a4h0jf9000lmqutlh9sw4c2x", types.ModuleAddress.String())
}
