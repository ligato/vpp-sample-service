#!/usr/bin/env bash
LIGATO_BGP_AGENT="/vendor/github.com/ligato/bgp-agent"
source .$LIGATO_BGP_AGENT/scripts/testOuput.sh
exitCode=0

VENDOR_BGP_AGENT_DOCKER=$LIGATO_BGP_AGENT"/docker"
VENDOR_BGP_AGENT_RR_DOCKER=$VENDOR_BGP_AGENT_DOCKER"/gobgp_route_reflector"

#####Setup Docker Network ##################################################
.$VENDOR_BGP_AGENT_RR_DOCKER/create-ligato-network-for-docker.sh
#####Download Dockers ##################################################
.$VENDOR_BGP_AGENT_DOCKER/gobgp_for_rr/pull-docker.sh
./docker/dev_bgp_to_vpp/pull-docker.sh

#####Run VPP Dockers ##################################################
.$VENDOR_BGP_AGENT_RR_DOCKER/start-routereflector.sh gobgp-client-in-docker
./docker/dev_bgp_to_vpp/scripts/start-bgp-to-vpp.sh


#####Add Path##################################################
.$VENDOR_BGP_AGENT_RR_DOCKER/addPath.sh &
#####Start Vpp##################################################
./docker/dev_bgp_to_vpp/scripts/start-vpp.sh &
#
##Validate l3plugin advertized path
expected=("SendStaticRouteToVPP &{65001 101.0.0.0/24 101.0.10.1}
")
#
./docker/dev_bgp_to_vpp/scripts/run_vpp_l3_bgp_routes_example.sh > log &
sleep 50
testOutput "$(less log)" "${expected}"
#
###check FIB
./docker/dev_bgp_to_vpp/scripts/check-vpp-fib.sh > fib &
sleep 5
expected=("101.0.10.1/32
")

testOutput "$(less fib)" "${expected}"

#####Close Dockers##################################################
.$VENDOR_BGP_AGENT_RR_DOCKER/stop-routereflector.sh
./docker/dev_bgp_to_vpp/scripts/stop-bgp-to-vpp.sh
.$VENDOR_BGP_AGENT_RR_DOCKER/remove-ligato-network-for-docker.sh
##########################################################################
exit ${exitCode}