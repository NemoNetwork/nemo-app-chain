#!/bin/bash
set -eo pipefail

# For legacy reasons, our internal dev environment runs `/nemo-network/start.sh` as the entrypoint for
# cosmovisor images. This file serves a stub we can provide in those images.

cosmovisor "$@"
