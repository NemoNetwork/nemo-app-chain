#!/bin/bash
set -e

COMPOSE_FILE="docker-compose-local-deployment.yml"
BRANCH="feat/update-indexer"

echo "[1/4] Pulling latest code from $BRANCH..."
git fetch origin
git checkout "$BRANCH"
git pull origin "$BRANCH"

echo "[2/4] Building postgres-package image..."
docker compose -f "$COMPOSE_FILE" build postgres-package

echo "[3/4] Running database migrations..."
docker compose -f "$COMPOSE_FILE" run --rm postgres-package

echo "[4/4] Rebuilding and restarting comlink..."
docker compose -f "$COMPOSE_FILE" build comlink
docker compose -f "$COMPOSE_FILE" up -d --no-deps comlink

echo ""
echo "Deployment complete."
echo "Verify migration:"
echo "  docker compose -f $COMPOSE_FILE exec postgres psql -U nemo_dev -d nemo_dev -c '\d user_complaints'"
echo "Verify endpoint:"
echo "  curl -s -o /dev/null -w '%{http_code}' -X POST http://localhost:3002/v4/userComplaints -H 'Content-Type: application/json' -d '{\"message\": \"test\"}'"
