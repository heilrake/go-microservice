#!/bin/bash
# Keeps postgres port-forwards alive for pgAdmin
# trip-postgres    → localhost:5433
# driver-postgres  → localhost:5434
# user-postgres    → localhost:5435
# payment-postgres → localhost:5436

forward() {
  local svc=$1
  local port=$2
  while true; do
    echo "[db-forward] Starting $svc on localhost:$port"
    kubectl port-forward service/$svc $port:5432
    echo "[db-forward] $svc disconnected, restarting in 2s..."
    sleep 2
  done
}

forward trip-postgres 5433 &
forward driver-postgres 5434 &
forward user-postgres 5435 &
forward payment-postgres 5436 &

echo "Port-forwards running. Press Ctrl+C to stop all."
wait
