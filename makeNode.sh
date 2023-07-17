#!/bin/bash
export IP=$(hostname -I)
curl -X POST -d "{\"ip\": \"$IP\"}" 194.163.166.108:9999/node/make
