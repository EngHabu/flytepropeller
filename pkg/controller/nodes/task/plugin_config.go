package task

import (
	"context"
	"strings"

	"github.com/lyft/flytepropeller/pkg/controller/nodes/task/backoff"

	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flytestdlib/logger"

	"github.com/lyft/flytepropeller/pkg/controller/nodes/task/config"
	"github.com/lyft/flytepropeller/pkg/controller/nodes/task/k8s"
)

func WranglePluginsAndGenerateFinalList(ctx context.Context, cfg *config.TaskPluginConfig, pr PluginRegistryIface) ([]core.PluginEntry, error) {
	allPluginsEnabled := false
	enabledPlugins := make(map[string]config.PluginConfig)
	if cfg != nil {
		enabledPlugins = cfg.GetEnabledPlugins()
	}
	if len(enabledPlugins) == 0 {
		allPluginsEnabled = true
	}

	var finalizedPlugins []core.PluginEntry
	logger.Infof(ctx, "Enabled plugins: %+v", enabledPlugins)
	logger.Infof(ctx, "Loading core Plugins, plugin configuration [all plugins enabled: %v]", allPluginsEnabled)
	for _, cpe := range pr.GetCorePlugins() {
		id := strings.ToLower(cpe.ID)
		pluginCfg, pluginEnabled := enabledPlugins[id]
		if !allPluginsEnabled && !pluginEnabled {
			logger.Infof(ctx, "Plugin [%s] is DISABLED (not found in enabled plugins list).", id)
		} else {
			logger.Infof(ctx, "Plugin [%s] ENABLED", id)
			cpe.DefaultForTaskTypes = pluginCfg.DefaultForTaskTypes
			finalizedPlugins = append(finalizedPlugins, cpe)
		}
	}

	// Create a single backOffManager for all the plugins
	backOffController := backoff.NewController(ctx)

	// Create a single resource monitor object for all plugins to use
	monitorIndex := k8s.NewResourceMonitorIndex()

	k8sPlugins := pr.GetK8sPlugins()
	for i := range k8sPlugins {
		kpe := k8sPlugins[i]
		id := strings.ToLower(kpe.ID)
		pluginConfig, pluginEnabled := enabledPlugins[id]
		if !allPluginsEnabled && !pluginEnabled {
			logger.Infof(ctx, "K8s Plugin [%s] is DISABLED (not found in enabled plugins list).", id)
		} else {
			logger.Infof(ctx, "K8s Plugin [%s] is ENABLED.", id)
			plugin := core.PluginEntry{
				ID:                  id,
				RegisteredTaskTypes: kpe.RegisteredTaskTypes,
				LoadPlugin: func(ctx context.Context, iCtx core.SetupContext) (plugin core.Plugin, e error) {
					return k8s.NewPluginManagerWithBackOff(ctx, iCtx, kpe, backOffController, monitorIndex)
				},
				IsDefault:           kpe.IsDefault,
				DefaultForTaskTypes: pluginConfig.DefaultForTaskTypes,
			}
			finalizedPlugins = append(finalizedPlugins, plugin)
		}
	}
	return finalizedPlugins, nil
}
