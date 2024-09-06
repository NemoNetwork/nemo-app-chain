package constants

import "github.com/nemo-network/v4-chain/protocol/testutil/encoding"

var (
	TestEncodingCfg = encoding.GetTestEncodingCfg()
	TestTxBuilder   = TestEncodingCfg.TxConfig.NewTxBuilder()
)
