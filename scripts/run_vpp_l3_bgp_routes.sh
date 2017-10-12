#!/usr/bin/env bash
LIGATO_BGP_AGENT="/vendor/github.com/ligato/bgp-agent"
source .$LIGATO_BGP_AGENT/scripts/testOuput.sh
exitCode=0

VENDOR_BGP_AGENT_RR_DOCKER=$LIGATO_BGP_AGENT"/docker/gobgp_route_reflector"
VENDOR_BGP_AGENT_RR_DOCKER_SCRIPTS=$VENDOR_BGP_AGENT_RR_DOCKER"/usage_scripts"

echo "## setup the infrastructure"
#####Setup Docker Network ##################################################
echo "### setup the docker network"
.$VENDOR_BGP_AGENT_RR_DOCKER_SCRIPTS/create-ligato-network-for-docker.sh
echo "done"
#####Download Dockers ##################################################
echo "### pull the Route reflector docker image"
.$VENDOR_BGP_AGENT_RR_DOCKER/pull-docker.sh
echo "done"
echo "### pull the bgp-to-vpp docker image"
./docker/dev_bgp_to_vpp/pull-docker.sh
echo "done"

echo ""
echo "## running the examples"
echo "### running the vpp_l3_bgp_routes example"
#####Run VPP Dockers ##################################################
echo "#### starting the Route reflector docker container"
.$VENDOR_BGP_AGENT_RR_DOCKER_SCRIPTS/start-routereflector.sh gobgp-client-in-docker
echo "done"
echo "#### starting the bgp-to-vpp docker container"
./docker/dev_bgp_to_vpp/scripts/start-bgp-to-vpp.sh
echo "done"


#####Add Path##################################################
echo "#### advertizing the path to the Route reflector docker container (background run)"
.$VENDOR_BGP_AGENT_RR_DOCKER_SCRIPTS/addPath.sh &
echo "launched in the background"
#####Start Vpp##################################################
echo "#### starting of the VPP inside the bgp-to-vpp docker container (background run)"
./docker/dev_bgp_to_vpp/scripts/start-vpp.sh &
echo "launched in the background"

#####Start Go Example###########################################
echo "#### running the Go example (bgptol3 plugin, goBGP plugin, pluginInterface)"
./docker/dev_bgp_to_vpp/scripts/run_vpp_l3_bgp_routes_example.sh > log &
sleep 50
echo "done"

##Validate l3plugin advertized path
echo "#### validating Go example output"
expected=("SendStaticRouteToVPP &{65001 101.0.0.0/24 101.0.10.1}
")
testOutput "$(less log)" "${expected}"
echo "done"

###check FIB
echo "#### checking the VPP's FIB table for the entry corresponding to the route added into the Route reflector"
./docker/dev_bgp_to_vpp/scripts/check-vpp-fib.sh > fib &
sleep 5
expected=("101.0.10.1/32
")
testOutput "$(less fib)" "${expected}"
echo "done"

echo ""
echo "## cleanup"
#####Close Dockers##################################################
echo "### stop and remove the Route reflector docker container"
.$VENDOR_BGP_AGENT_RR_DOCKER_SCRIPTS/stop-routereflector.sh
echo "done"
echo "### stop and remove the bgp-to-vpp docker container"
./docker/dev_bgp_to_vpp/scripts/stop-bgp-to-vpp.sh
echo "done"
echo "### remove the docker network"
.$VENDOR_BGP_AGENT_RR_DOCKER_SCRIPTS/remove-ligato-network-for-docker.sh
echo "done"
##########################################################################
exit ${exitCode}