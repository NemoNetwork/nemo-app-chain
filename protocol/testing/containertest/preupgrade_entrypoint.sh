#!/bin/bash
set -eo pipefail

# This entrypoint configures the chain to start in a preupgrade state before forwarding args to cosmovisor.
# The following need to be done:
#   - Set the preupgrade binary (downloaded in containertest.sh) to be the genesis binary
#   - Set the binary at the current commit to be the upgrade binary
#   - Set the genesis to be the preupgrade genesis. This is generated by genesis generation script at the
#     preupgrade commit. It is copied directly from a container.
#   - Set the voting period of the preupgrade genesis to 15s.

if [[ -z "${UPGRADE_TO_VERSION}" ]]; then
    echo >&2 "UPGRADE_TO_VERSION must be set"
    exit 1
fi

if [[ -z "${DAEMON_NAME}" ]]; then
    echo >&2 "DAEMON_NAME must be set"
    exit 1
fi

if [[ -z "${DAEMON_HOME}" ]]; then
    echo >&2 "DAEMON_HOME must be set"
    exit 1
fi

rm "$DAEMON_HOME/cosmovisor/genesis/bin/nemod"
ln -s /bin/nemod_preupgrade "$DAEMON_HOME/cosmovisor/genesis/bin/nemod"
mkdir -p "$DAEMON_HOME/cosmovisor/upgrades/$UPGRADE_TO_VERSION/bin/"
ln -s /bin/nemod "$DAEMON_HOME/cosmovisor/upgrades/$UPGRADE_TO_VERSION/bin/nemod"

rm "$DAEMON_HOME/config/genesis.json"
cp "$HOME/preupgrade_genesis.json" "$DAEMON_HOME/config/genesis.json"
dasel put -t string -f "$DAEMON_HOME/config/genesis.json" '.app_state.gov.params.voting_period' -v '15s'

cosmovisor run "$@"
