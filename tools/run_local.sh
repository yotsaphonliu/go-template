#!/bin/bash

trap killgroup SIGINT

killgroup(){
  echo killing...
  kill 0
}

rm -rf go-template

go build -o go-template src/main.go
TZ=Asia/Bangkok ./go-template serve-http-api --config ./cfg/config.yaml

# go build -o go-template src/main.go
# TZ=Asia/Bangkok ./go-template migrate-db --config ./cfg/config.yaml &
# TZ=Asia/Bangkok ./go-template serve-http-api --config ./cfg/config.yaml

wait