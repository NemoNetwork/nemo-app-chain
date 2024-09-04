FROM nemo-network-base

COPY ./testing/containertest/nemo-networkd_preupgrade* /bin/
COPY ./testing/containertest/containertest.sh /nemo-network/containertest.sh
COPY ./testing/containertest/preupgrade_entrypoint.sh /nemo-network/preupgrade_entrypoint.sh
COPY ./testing/containertest/preupgrade_genesis.json /nemo-network/preupgrade_genesis.json
COPY ./testing/genesis.sh /nemo-network/genesis.sh
COPY ./daemons/pricefeed/client/constants/testdata /nemo-network/exchange_config
COPY ./testing/delaymsg_config /nemo-network/delaymsg_config
COPY ./testing/version/. /nemo-network/

RUN go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0

RUN /nemo-network/containertest.sh

COPY ./testing/containertest/config/pricefeed_exchange_config.toml /nemo-network/chain/.alice/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /nemo-network/chain/.bob/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /nemo-network/chain/.carl/config/pricefeed_exchange_config.toml
COPY ./testing/containertest/config/pricefeed_exchange_config.toml /nemo-network/chain/.dave/config/pricefeed_exchange_config.toml

ENTRYPOINT ["nemod"]
