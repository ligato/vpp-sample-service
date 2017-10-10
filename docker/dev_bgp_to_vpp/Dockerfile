FROM ligato/dev-vpp-agent:v1.0.5

MAINTAINER Filip Gschwandtner <filip.gschwandtner@pantheon.tech>

EXPOSE 179

RUN mkdir -p /root/go/src/github.com/ligato/vpp-sample-service/examples/vpp_l3_bgp_routes

# adding go binary to path
ENV PATH /root/go/src/github.com/ligato/vpp-sample-service/examples/vpp_l3_bgp_routes:$PATH

ENV PATH /usr/local/go/bin/go:$PATH

WORKDIR /root/

# disables the start of vpp and vpp-agent at the start of docker instance
RUN sed -i.bak -e '6,$d' /etc/supervisord.conf

# run supervisor as the default executable
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]