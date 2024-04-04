package main

import (
	"website-telemetry-demo/cmd/app"
	"website-telemetry-demo/configs"
)

func main() {
	staticCfg, err := configs.NewStaticConfig()
	if err != nil {
		panic(err)
		return
	}

	app.StartApp(staticCfg)
}
