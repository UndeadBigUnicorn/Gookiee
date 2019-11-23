#!/bin/bash

set -e

echo ""
echo "--- BENCH ECHO START ---"
echo ""

cd $(dirname "${BASH_SOURCE[0]}")
function cleanup {
    echo "--- BENCH ECHO DONE ---"
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
    GOMAXPROCS=4 $2 --port $4 &
    sleep 1
    echo "*** 50 connections, 10 seconds, 6 byte packets"
    nl=$'\r\n'
    tcpkali --workers 1 -c 50 -T 10s -m "PING{$nl}" 127.0.0.1:$4
    echo "--- DONE ---"
    echo ""
}

gobench "EVIO" bin/echo ./main.go 5002