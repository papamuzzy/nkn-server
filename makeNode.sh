#!/bin/bash
IP=$(hostname -I)
curl -X POST -d "{\"ip\": \"$IP\"}" http://localhost:9999/node/make
