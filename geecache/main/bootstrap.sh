#!/bin/bash
trap "kill 0" EXIT

go run main/main.go -port=8001 &
go run main/main.go -port=8002 &
go run main/main.go -port=8003 -api=1 &

sleep 4
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait
