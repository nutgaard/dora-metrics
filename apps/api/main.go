package main

import (
	ConfigLoader "nutgaard/dora-metrics/config"
)

func main() {
	runApp(ConfigLoader.ReadConfig())
}
