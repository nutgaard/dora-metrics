package main

import (
	"nutgaard/dora-metrics/cmd/dora-metrics-api"
	"nutgaard/dora-metrics/internal/config"
)

func main() {
	dora_metrics_api.RunApp(config.ReadConfig())
}
