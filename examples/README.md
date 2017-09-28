#TODO: change links to vpp-agent github pull request, test this procedure, resolve other todos
# Examples

Examples package provide examples of BGP-to-L3 plugin usage.

## [End-to-End example](end_to_end.example.go) 
This is a runtime test that show the capability to retrieve and transport the basic BGP information (prefix/next hop) from the BGP Route Reflector node to the VPP node.
### Architecture & Data flow
![arch](../doc/img/endtoendBGPExample.png "End-to-end BGP Example")

The example integrates together all needed part to achieve the VPP configuration based on the BGP information. It integrates the [BGPtoL3 plugin](../bgptol3plugin/README.md) with the [VPP-Agent](https://github.com/ligato/vpp-agent) and the [BGP-Agent](../../agent/README.md). The VPP agent is configured to connect to the VPP and the GoBGP plugin to connect to the Route Reflector node.

The data flow starts in the Route Reflector node, where we manually insert the new path into the [Route Reflector node](../../agent/route-reflector-gobgp-docker/README.md). Remote GoBGP plugin will be advertized about this new path via standard BGP protocol. Then this information will be translated to the API format that is exposed to all API clients of the BGP-Agent. This translation is needed due to possible future plugins that can retrieve the BGP information using other frameworks(exaBGP, Quagga,...) and therefore in different formats. 

The BGPtoL3 plugin is registered as watcher client in the BGP-Agent and will be notified by any new update. The BGPtoL3 plugin will translate this information into the VPP configuration structures provided by the L3 plugin inside the VPP-Agent. From there, the VPP-Agent will handle the rest of data flow down to the VPP.
### Infrastructure setup
To be able to run this example you must setup the infrastructure first. We will use docker containers. This is the docker containers architecture for this example:
![arch](../doc/img/dockerArchitecture.png "docker container architecture of end-to-end BGP Example") 

It consist of the [VPP-endpoint](../docker/README.md) docker container and the [Route reflector](../../agent/route-reflector-gobgp-docker/README.md) docker container. In infrastructure setup we will only create images for these containers. 

At this point, it is expected that you have:
 * Installed the docker CE (docker engine version >=1.10). If you haven't please install it(for [ubuntu](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/) or [other](https://docs.docker.com/engine/installation/)). 
 * Properly filled the [GOPATH environment variable](https://github.com/golang/go/wiki/Setting-GOPATH). Installing go is not necessary. 
 * Downloaded the [Ligato VPP-Agent](https://github.com/ligato/vpp-agent). In case of installed go, just run:
```
go get github.com/ligato/vpp-agent
```
Change the path to the VPP-endpoint docker folder
```
cd $GOPATH/src/github.com/ligato/vpp-agent/docker/bgpexample
```
Build the VPP-endpoint image
```
./build-vpp-agent-image.sh
``` 
Change path to the Route reflector docker folder
```
cd $GOPATH/src/github.com/ligato/vpp-agent/vendor/github.com/ligato/bgp-agent/route-reflector-gobgp-docker
```
Build the Route reflector image
```
./build-image-routereflector.sh
``` 
Now you should see something like this:

```
REPOSITORY                 TAG                 IMAGE ID            CREATED             SIZE
routereflector             latest              a6d47c8559da        11 seconds ago      982MB
ligato-bgp/dev-vpp-agent   v1.0.4              1ea7a0617fe7        2 minutes ago       4.52GB
```
Process of building of the images has downloaded also other images that served as base images in the creation process. You can delete these base images if you want.

To be able to have static ip addresses for running docker images, we need to create separate network that can be used by docker.
```
./create-ligato-network-for-docker.sh
```    
### Example run    
We will need 4 linux terminals. To differentiate commands in terminal we will use different [prompt string](http://www.linuxnix.com/linuxunix-shell-ps1-prompt-explained-in-detail) for each terminal:

`[vpp-endpoint-vpp-console]$` 

* The terminal for the VPP console inside the VPP-endpoint docker container. We can interact with the VPP here, but also see its new log entries.

`[vpp-endpoint-example-run]$`

* The terminal for Running go example inside the VPP-endpoint docker container. 

`[rr-bgp-server]$` 

* The terminal for BGP server inside the Route reflector docker container that acts as a Route reflector. We can see logs of the BGP server here.

`[rr-manual-info-addition]$` 

* The terminal for adding the prefix/nexthop information directly to the BGP server(acting like Route reflector) in the Route reflector docker container. 

Lets run the example:

<b>1. Start the docker containers, VPP and BGP server.</b>

Change the directory so we can use the helper scripts in the VPP-endpoint docker folder 
```
[vpp-endpoint-vpp-console]$ cd $GOPATH/src/github.com/ligato/vpp-agent/docker/bgpexample
``` 
Start the VPP-endpoint docker container
```
[vpp-endpoint-vpp-console]$ ./start-vpp-endpoint.sh
```
Start the VPP inside the VPP-endpoint docker container
```
[vpp-endpoint-vpp-console]$ ./start-vpp.sh
```
Switch to the ```[rr-bgp-server]``` terminal and change the directory so we can use the helper scripts in the route reflector docker folder
```
[rr-bgp-server]$ cd $GOPATH/src/github.com/ligato/vpp-agent/vendor/github.com/ligato/bgp-agent/route-reflector-gobgp-docker
```
Start the route reflector docker container
```
[rr-bgp-server]$ ./start-routereflector-for-client-in-docker.sh
```
<b>2. Run the go code example</b> 

Switch to the ```[vpp-endpoint-example-run]$``` and change the directory so we can use the helper scripts in the VPP-endpoint docker folder 
```
[vpp-endpoint-example-run]$ cd $GOPATH/src/github.com/ligato/vpp-agent/docker/bgpexample
``` 
Run the go example (it will be build automatically before running)
```
[vpp-endpoint-example-run]$ ./build-and-start-vpp-agent-with-bgp-plugin.sh
```
This example has limited duration to show the graceful stop of components. The current duration is set to 4 minutes, but you can change it by modifying the ```exampleDuration``` constant in the [end_to_end_example.go](end_to_end_example.go) file.
 
(Note: The run of the example will take short time to compile and when running it will initially take some time, ~20 seconds, to receive the first BGP information due to the session initialization)

<b>2. Add new route information to the Route reflector</b>
Switch to the ```[rr-manual-info-addition]``` terminal and change the directory so we can use the helper scripts in the route reflector docker folder
```
[rr-manual-info-addition]$ cd $GOPATH/src/github.com/ligato/vpp-agent/vendor/github.com/ligato/bgp-agent/route-reflector-gobgp-docker
```
Connect to the bash console inside the Route reflector docker container
```
[rr-manual-info-addition]$ ./connect-to-routereflector.sh
```
Add new prefix(`101.0.0.0/24`)/nexthop(`192.168.1.1`) information to the Route reflector
```
[rr-manual-info-addition]$ gobgp global rib add -a ipv4 101.0.0.0/24 nexthop 192.168.1.1
``` 
(Note: `192.168.1.1` is the IP address of the VPPs virtual interface)

<b>3. Verify that the prefix/nexthop information did arrive in the VPP</b>
Switch back to the `[vpp-endpoint-vpp-console]` and list in the VPP console the appropriate configuration (IP Fib configuration)
```
[vpp-endpoint-vpp-console]$ vpp> show ip fib
```
You should see something like this:
```
TODO add example of vpp success output
```
Additionally, you can check also the go runtime logs for adding prefix/nexthop in the `[vpp-endpoint-example-run]` terminal.