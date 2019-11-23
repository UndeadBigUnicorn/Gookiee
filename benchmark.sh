#!/bin/bash

set -e

pl=$1
if [ "$pl" == "" ]; then
    pl="1"
fi

echo ""
echo "--- BENCH GOOKIE PIPELINE $pl START ---"
echo ""

cd $(dirname "${BASH_SOURCE[0]}")

function cleanup {
    echo "--- BENCH GOOKIE PIPELINE $pl DONE ---"
    kill -9 $(jobs -rp)
    wait $(jobs -rp) 2>/dev/null
}
trap cleanup EXIT

mkdir -p bin

function gobench {
    echo "--- $1 ---"
    if [ "$3" != "" ]; then
        go build -o $2 $3
    fi
    GOMAXPROCS=2 $2 --port $4 &
    sleep 1
    echo "*** $5 connections, 10 seconds, 6 byte packets"
    nl=$'\r\n'
    tcpkali --workers 1 -c "$5" -T 10s -m "PING{$nl}" 127.0.0.1:"$4"
    echo "--- DONE ---"
    echo ""
}

gobench "GOOKIE network 50 connections" bin/gookiee-network ./main.go 8964 50