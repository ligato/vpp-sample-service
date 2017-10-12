## VPP L3 BGP

Implements the `Vpp sample service` plugin that allows plugin
to render learned IP-based routes to l3 plugin. Plugin retrieves the information from a Watcher plugin and
publish it to l3 plugin, acting as a translator. Plugin provides [core agent](https://github.com/ligato/cn-infra/tree/master/core) functionality. 

To acquire IP-based routes been advertized to l3 plugin we must do 2 things:

1. Create plugin implementing [Watcher interface](https://github.com/ligato/bgp-agent/tree/master/bgp/bgp_api.go)
2. Create [vppl3bgp plugin](https://github.com/ligato/vpp-sample-service/tree/master/plugins/vppl3bgp/plugin.go) injecting `watcher plugin` into constructor `vppl3bgp.New(...)`, i.e.:

```
...
	goBgpPlugin := gobgp.New(gobgp.Deps{
		PluginInfraDeps: deps,
		SessionConfig:   goBgpConfig})

	// Create BGP-to-L3 plugin that is plugin of the Vpp Agent
	bgptol3 := vppl3bgp.New(vppl3bgp.Deps{
		PluginInfraDeps: deps,
		Watcher:         goBgpPlugin,
	})
...
```

For further usage please look into our [example](https://github.com/ligato/vpp-sample-service/tree/master/examples/vpp_l3_bgp_routes/).