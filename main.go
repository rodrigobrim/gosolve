package main

import (
	"config"
	. "deploy"
	"os"
)

func main() {

	call := os.Args[1]
	env := os.Args[2]

	cfg := config.GetConfig(env)

	// Todo
	// Set cfg.LaunchTemplateInput to build.Build(cfg.MyBuildConfig)

	switch call {
	case "apply":
		Apply(cfg)
	case "destroy":
		Destroy(cfg)
	}
}
