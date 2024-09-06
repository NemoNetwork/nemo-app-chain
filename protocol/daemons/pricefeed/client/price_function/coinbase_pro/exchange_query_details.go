package coinbase_pro

import (
	"github.com/nemo-network/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/nemo-network/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	CoinbaseProDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_COINBASE_PRO,
		Url:           "https://api.pro.coinbase.com/products/$/ticker",
		PriceFunction: CoinbaseProPriceFunction,
	}
)
