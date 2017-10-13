## Development Docker Image

This image can be used to get started with the vpp-sample-service Go code. It 
contains:

- The development environment with all libs & dependencies required 
  to build both the VPP itself and the VPP sample service
- A pre-built vpp ready to be used
- A link with local environment Vpp sample service example

This folder also contains some helper scripts. These scripts are meant to run without `sudo` command, therefore the environment for docker must be [altered accordingly](https://docs.docker.com/engine/installation/linux/linux-postinstall/#manage-docker-as-a-non-root-user).
### Getting an Image from Dockerhub
For a quick start with the Development image, you can use pre-built 
Development docker images based on [Dockerhub](https://hub.docker.com/r/ligato/dev-vpp-agent/).
Images has been modified, removing default start of VPP.
The pre-built Development docker images are available from [Dockerhub](https://hub.docker.com/r/ligato/dev-bgp-to-vpp/),
or you can just type:
```
docker pull ligato/dev-bgp-to-vpp
```
Then you can start the downloaded Development image using provided script.
```
./docker/dev_bgp_to_vpp/scripts/start-bgp-to-vpp.sh
```

Start VPP.
```
./docker/dev_bgp_to_vpp/scripts/start-vpp.sh
```

To open another terminal into the image:
```
./docker/dev_bgp_to_vpp/scripts/connect-to-bgp-to-vpp.sh
```

Remove docker.
```
./docker/dev_bgp_to_vpp/scripts/stop-bgp-to-vpp.sh
```

### Building Locally
To build the docker image on your local machine,  type:
```
./build-image.sh
```

#### Verifying a Created or Downloaded Image
You can verify the newly built or downloaded image as follows:

```
docker images
``` 

You should see something like this:

```
REPOSITORY                            TAG                 IMAGE ID            CREATED             SIZE
ligato/dev-bgp-to-vpp                 v1.5                b53924b8255a        5 minutes ago          4.78GB
```
Get the details of the newly built or downloaded image:

```
docker image inspect dev-bgp-to-vpp
docker image history dev-bgp-to-vpp