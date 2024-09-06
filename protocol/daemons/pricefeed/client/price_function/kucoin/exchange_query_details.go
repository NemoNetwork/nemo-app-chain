package kucoin

import (
	"github.com/nemo-network/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/nemo-network/v4-chain/protocol/daemons/pricefeed/client/types"
)

var (
	KucoinDetails = types.ExchangeQueryDetails{
		Exchange:      exchange_common.EXCHANGE_ID_KUCOIN,
		Url:           "https://api.kucoin.com/api/v1/market/allTickers",
		PriceFunction: KucoinPriceFunction,
		IsMultiMarket: true,
	}
)
