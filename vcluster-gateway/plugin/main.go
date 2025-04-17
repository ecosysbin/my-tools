package main

import (
	"github.com/loft-sh/vcluster-pod-hooks/hooks"
	"github.com/loft-sh/vcluster-sdk/plugin"
)

func main() {
	ctx := plugin.MustInit()
	plugin.MustRegister(hooks.NewPodHook(ctx))
	// plugin.MustRegister(hooks.NewServiceHook())
	// plugin.MustRegister(hooks.NewSecretHook())
	plugin.MustRegister(hooks.NewScHook(ctx))
	plugin.MustRegister(hooks.NewPvcHook(ctx))
	plugin.MustStart()
}
