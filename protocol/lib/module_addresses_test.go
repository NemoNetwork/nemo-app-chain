package lib_test

import (
	"github.com/nemo-network/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGovModuleAddress(t *testing.T) {
	require.Equal(t, "nemo10d07y265gmmuvt4z0w9aw880jnsr700j3m62vw", lib.GovModuleAddress.String())
}
