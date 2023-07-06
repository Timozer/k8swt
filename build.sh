#!/bin/bash

if [ $# -eq 0 ]; then
    echo 'please input the docker version'
    exit 1
fi

CGO_ENABLED=0 go build -ldflags="-s -w" -o k8swt

sudo docker buildx build --output type=docker -t 878592748/k8swt:latest .

sudo docker tag 878592748/k8swt:latest 878592748/k8swt:$1

sudo docker push 878592748/k8swt:latest
sudo docker push 878592748/k8swt:$1