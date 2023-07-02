#!/bin/bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o k8swt

sudo docker buildx build --output type=docker -t k8swt:latest .