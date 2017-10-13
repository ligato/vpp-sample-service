#!/bin/bash

docker exec -it bgp-to-vpp bash -c "vpp unix { interactive } plugins { plugin dpdk_plugin.so { disable } }"