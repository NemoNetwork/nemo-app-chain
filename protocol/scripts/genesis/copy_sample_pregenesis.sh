#!/bin/bash

./scripts/genesis/prod_pregenesis.sh nemod
cp /tmp/prod-chain/.dydxprotocol/config/sorted_genesis.json ./scripts/genesis/sample_pregenesis.json
