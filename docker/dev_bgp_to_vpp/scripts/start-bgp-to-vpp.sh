#!/bin/bash

VM_GOROOT="/usr/local/go"
VM_GOPATH="/root/go"
VPP_ENDPOINT_IP_ADDRESS="172.18.0.1"
RELATIVE_PATH_TO_CODE_BASE="src/github.com"

#run Vpp with Vpp-agent
sudo docker run -d --name vpp-endpoint --net ligato-bgp-network --ip $VPP_ENDPOINT_IP_ADDRESS --privileged -p 5001:5002 -v /$GOPATH/$RELATIVE_PATH_TO_CODE_BASE:$VM_GOPATH/$RELATIVE_PATH_TO_CODE_BASE -w /root/go --rm ligato/dev-bgp-to-vpp:v1.5