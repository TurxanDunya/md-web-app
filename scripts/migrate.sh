#!/bin/bash
cd "$(dirname "$0")/.."
set -a
source .env
set +a

command=$1
name=$2

case $command in
  up)
    migrate -path migrations -database "$DATABASE_URL" up
    ;;
  down)
    count=${name:-1}
    read -p "Rolling back $count migration(s). Continue? [y/N] " confirm
    if [ "$confirm" = "y" ]; then
      migrate -path migrations -database "$DATABASE_URL" down "$count"
    fi
    ;;
  create)
    migrate create -ext sql -dir migrations -seq "$name"
    ;;
  force)
    migrate -path migrations -database "$DATABASE_URL" force "$name"
    ;;
esac
