#!/bin/bash
set -e

# Always run from the indexer root regardless of where the script is called from
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

COMPOSE_FILE="docker-compose-local-deployment.yml"
BRANCH="feat/update-indexer"

echo "[1/4] Pulling latest code from $BRANCH..."
# Stash any local changes so checkout never aborts
git stash --include-untracked --quiet
git fetch origin
git checkout "$BRANCH"
git pull origin "$BRANCH"

echo "[2/4] Building postgres-package image..."
docker compose -f "$COMPOSE_FILE" build postgres-package

echo "[3/4] Running database migrations..."
# postgres-package depends_on postgres (service_healthy), so it waits automatically
docker compose -f "$COMPOSE_FILE" run --rm postgres-package

echo "[4/4] Rebuilding and restarting comlink..."
docker compose -f "$COMPOSE_FILE" build comlink
# --no-deps avoids touching postgres/redis/kafka; migration already ran in step 3
docker compose -f "$COMPOSE_FILE" up -d --no-deps comlink

echo ""
echo "Deployment complete."
echo ""
echo "Verify migration:"
echo "  docker compose -f $COMPOSE_FILE exec postgres psql -U nemo_dev -d nemo_dev -c '\d user_complaints'"
echo ""
echo "Verify endpoint:"
echo "  curl -s -o /dev/null -w '%{http_code}' -X POST http://localhost:3002/v4/userComplaints -H 'Content-Type: application/json' -d '{\"message\": \"test\"}'"
