package config

import "fork_go_im/pkg/config"

func init() {
	config.Add("cache", config.StrMap{
		"redis": map[string]interface{}{
			"addr":     config.Env("REDIS_HOST", "localhost"),
			"port":     config.Env("REDIS_PORT", "6379"),
			"password": config.Env("REDIS_PASSWORD", ""),
			"db":       config.Env("REDIS_DB", 0),
		},
	})
}
