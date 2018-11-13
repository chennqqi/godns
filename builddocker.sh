#!/bin/bash
go build -o godns -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -v 
sudo docker build -t "sort/godns:$(cat VERSION)" -f Dockerfile.local .

