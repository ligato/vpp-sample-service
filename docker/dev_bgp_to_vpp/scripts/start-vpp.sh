#!/bin/bash

sudo docker exec -d vpp-endpoint bash -c "vpp unix { nodaemon cli-listen 0.0.0.0:5002 cli-no-pager } plugins { plugin dpdk_plugin.so { disable } }"