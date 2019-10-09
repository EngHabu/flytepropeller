package task

import (
	"context"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flytestdlib/logger"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/lyft/flytepropeller/pkg/controller/nodes/task/config"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/task/k8s"
)

func WranglePluginsAndGenerateFinalList(ctx context.Context, cfg *config.TaskPluginConfig, pr PluginRegistryIface) ([]core.PluginEntry, error) {
	allPluginsEnabled := false
	enabledPlugins := sets.NewString()
	if cfg != nil {
		enabledPlugins = cfg.GetEnabledPluginsSet()
	}
	if enabledPlugins.Len() == 0 {
		allPluginsEnabled = true
	}

	var finalizedPlugins []core.PluginEntry
	logger.Infof(ctx, "Enabled plugins: %v", enabledPlugins.List())
	logger.Infof(ctx, "Loading core Plugins, plugin configuration [all plugins enabled: %v]", allPluginsEnabled)
	for _, cpe := range pr.GetCorePlugins() {
		if !allPluginsEnabled && !enabledPlugins.Has(cpe.ID) {
			logger.Infof(ctx, "Plugin [%s] is DISABLED.", cpe.ID)
		} else {
			logger.Infof(ctx, "Plugin [%s] ENABLED", cpe.ID)
			finalizedPlugins = append(finalizedPlugins, cpe)
		}
	}

	k8sPlugins := pr.GetK8sPlugins()
	for i := range k8sPlugins {
		kpe := k8sPlugins[i]
		if !allPluginsEnabled && !enabledPlugins.Has(kpe.ID) {
			logger.Infof(ctx, "K8s Plugin [%s] is DISABLED.", kpe.ID)
		} else {
			logger.Infof(ctx, "K8s Plugin [%s] is ENABLED.", kpe.ID)
			finalizedPlugins = append(finalizedPlugins, core.PluginEntry{
				ID:                  kpe.ID,
				RegisteredTaskTypes: kpe.RegisteredTaskTypes,
				LoadPlugin: func(ctx context.Context, iCtx core.SetupContext) (plugin core.Plugin, e error) {
					return k8s.NewPluginManager(ctx, iCtx, kpe)
				},
				IsDefault: kpe.IsDefault,
			})
		}
	}
	return finalizedPlugins, nil
}
