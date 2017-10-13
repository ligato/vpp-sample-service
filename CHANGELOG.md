# Release v1.0.0 (2017-10-13)

## Initial release

Release contains [basic example of usage](https://github.com/ligato/vpp-sample-service/tree/master/examples/vpp_l3_bgp_routes) 
of this plugin implementation.

Implementation contains basic [VPP L3 BGP plugin](https://github.com/ligato/vpp-sample-service/tree/master/plugins/vppl3bgp) 
that serves as an IP-based route provider to l3 plugin. 
It uses [BGP Agent plugin](https://github.com/ligato/bgp-agent) library as 
source of BGP information and it provides reachables IPv4 routes.

## Known Issues
Not known issue at the moment.

## Known Limitiations
The information flow is unidirectional, and flow from watcher plugin to l3 plugin.