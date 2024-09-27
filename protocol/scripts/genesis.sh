#!/usr/bin/env bash
set -eux

go run $SCRIPT_DIR/genesis.go --use-core=$USE_CORE_MARKETS --use-raydium=$USE_RAYDIUM_MARKETS \
    --use-uniswapv3-base=$USE_UNISWAPV3_BASE_MARKETS --use-coingecko=$USE_COINGECKO_MARKETS \
    --use-polymarket=$USE_POLYMARKET_MARKETS --use-coinmarketcap=$USE_COINMARKETCAP_MARKETS \
    --use-osmosis=$USE_OSMOSIS_MARKETS --temp-file=markets.json
MARKETS=$(cat markets.json)

echo "MARKETS content: $MARKETS"

NUM_MARKETS=$(echo "$MARKETS" | jq '.markets | length + 1')

jq '.consensus["params"]["abci"]["vote_extensions_enable_height"] = "2"' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
NUM_MARKETS=$NUM_MARKETS; jq --arg num "$NUM_MARKETS" '.app_state["oracle"]["next_id"] = $num' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["marketmap"]["market_map"] = ($markets | fromjson)' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["oracle"]["currency_pair_genesis"] += [$markets | fromjson | .markets | values | .[].ticker.currency_pair | {"currency_pair": {"Base": .Base, "Quote": .Quote}, "currency_pair_price": null, "nonce": 0} ]' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
MARKETS=$MARKETS; jq --arg markets "$MARKETS" '.app_state["oracle"]["currency_pair_genesis"] |= (to_entries | map(.value += {id: (.key + 1)} | .value))' "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"


# Check if MARKET_MAP_AUTHORITY environment variable exists and is not empty
if [ -n "${MARKET_MAP_AUTHORITY+x}" ] && [ -n "$MARKET_MAP_AUTHORITY" ]; then
    MARKET_MAP_AUTHORITY=$(printenv MARKET_MAP_AUTHORITY)
    jq --arg authority "$MARKET_MAP_AUTHORITY" \
    '.app_state["marketmap"]["params"]["market_authorities"] += [$authority]' \
    "$GENESIS" > "$GENESIS_TMP" && mv "$GENESIS_TMP" "$GENESIS"
fi

rm markets.json
