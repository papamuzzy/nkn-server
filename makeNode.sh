#!/bin/bash
IP=$(hostname -I)
curl -X POST -d "{\"ip\": \"$IP\"}"5.180.183.19:9999/node/make
