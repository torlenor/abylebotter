package giveawayplugin

import (
	"fmt"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
)

type config struct {
	Mods     []string
	OnlyMods bool
}

func parseConfig(c botconfig.PluginConfig) (config, error) {
	if c.Type != "giveaway" {
		return config{}, fmt.Errorf("Not a Giveaway Plugin config")
	}

	var mods []string
	var onlyMods bool

	var ok bool
	if mods, ok = c.Config["mods"].([]string); !ok {
	}

	if onlyMods, ok = c.Config["onlymods"].(bool); !ok {
	}

	if onlyMods && len(mods) == 0 {
		return config{}, fmt.Errorf("Cannot have a Giveaway Plugin configuration with OnlyMods = true but no Mods defined")
	}

	cfg := config{
		Mods:     mods,
		OnlyMods: onlyMods,
	}

	return cfg, nil
}
