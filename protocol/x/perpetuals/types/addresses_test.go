package types_test

import (
	"testing"

	"github.com/nemo-network/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestInsuranceFundModuleAddress(t *testing.T) {
	require.Equal(t, "nemo1c7ptc87hkd54e3r7zjy92q29xkq7t79wc4h5e2", types.InsuranceFundModuleAddress.String())
}
