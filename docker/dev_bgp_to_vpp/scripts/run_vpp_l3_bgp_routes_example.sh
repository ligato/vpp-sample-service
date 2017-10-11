#!/usr/bin/env bash

VM_GOROOT="/usr/local/go"

docker exec -t bgp-to-vpp sh -c "cd src/github.com/ligato/vpp-sample-service/examples/vpp_l3_bgp_routes;"$VM_GOROOT"/bin/go run vpp_l3_bgp_routes.go"