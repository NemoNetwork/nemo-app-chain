package types_test

import (
	"github.com/nemo-network/v4-chain/protocol/x/delaymsg/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModuleAddress(t *testing.T) {
	require.Equal(t, "nemo1mkkvp26dngu6n8rmalaxyp3gwkjuzztqkzp33f", types.ModuleAddress.String())
}
