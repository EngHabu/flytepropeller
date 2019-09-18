package launchplan

import (
	ctrlConfig "github.com/lyft/flytepropeller/pkg/controller/config"
)

// TODO: Consider moving this to FlyteIdl/Admin client config

//go:generate pflags AdminConfig --default-var defaultAdminConfig

var (
	defaultAdminConfig = &AdminConfig{
		TPS:          5,
		Burst:        10,
		MaxCacheSize: 10000,
	}

	adminConfigSection = ctrlConfig.ConfigSection.MustRegisterSection("admin-launcher", defaultAdminConfig)
)

type AdminConfig struct {
	// TPS indicates the maximum transactions per second to flyte admin from this client.
	// If it's zero, the created client will use DefaultTPS: 5
	TPS int64 `json:"tps" pflag:",The maximum number of transactions per second to flyte admin from this client."`

	// Maximum burst for throttle.
	// If it's zero, the created client will use DefaultBurst: 10.
	Burst int `json:"burst" pflag:",Maximum burst for throttle"`

	MaxCacheSize int `json:"cacheSize" pflag:",Maximum cache in terms of number of items stored."`
}

func GetAdminConfig() *AdminConfig {
	return adminConfigSection.GetConfig().(*AdminConfig)
}
