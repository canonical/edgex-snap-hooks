package options

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/canonical/edgex-snap-hooks/v2/env"
	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
)

func processAppAutostartOptions(apps []string) (map[string]*bool, error) {
	// get the apps' json structure
	jsonString, err := snapctl.Get("apps").Document().Run()
	if err != nil {
		return nil, fmt.Errorf("error reading 'apps' option: %s", err)
	}
	var options snapOptions
	err = json.Unmarshal([]byte(jsonString), &options)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling 'apps' option: %s", err)
	}

	appAutostart := make(map[string]*bool)
	for _, app := range apps {
		// get the configuration specified for each app
		// autostart := options.Apps[app].Autostart
		// if autostart == nil {
		// 	// no autostart option for the app
		// 	continue
		// }

		appAutostart[app] = options.Apps[app].Autostart

		if appAutostart[app] != nil {
			log.Debugf("%s: autostart=%t (app setting)", app, *appAutostart[app])
		}
	}

	return appAutostart, nil
}

func processGlobalAutostartOptions(apps []string) (map[string]*bool, error) {
	autostart, err := snapctl.Get("autostart").Run()
	if err != nil {
		return nil, fmt.Errorf("error reading 'autostart' option: %s", err)
	}

	appAutostart := make(map[string]*bool)
	for _, app := range apps {
		appAutostart[app], err = parseAutostart(autostart)
		if err != nil {
			return nil, fmt.Errorf("error parsing autostart for %s: %v", app, err)
		}
		if appAutostart[app] != nil {
			log.Debugf("%s: autostart=%t (global setting)", app, *appAutostart[app])
		}
	}

	return appAutostart, nil
}

func parseAutostart(value string) (*bool, error) {
	value = strings.ToLower(value)
	switch value {
	case "":
		return nil, nil
	// need to accept yes/no for EdgeX 2 backward compatibility
	case "true", "yes":
		b := true
		return &b, nil
	case "false", "no":
		b := false
		return &b, nil
	default:
		return nil, fmt.Errorf("invalid value for 'autostart': '%s'", value)
	}
}

// ProcessAutostart will start and enable the listed app(s)
// based on the value of autostart snap option
func ProcessAutostart(apps ...string) error {
	if len(apps) == 0 {
		return fmt.Errorf("empty apps list")
	}

	log.Infof("Processing autostart for: %v", apps)

	globalAppAutostart, err := processGlobalAutostartOptions(apps)
	if err != nil {
		return fmt.Errorf("error processing global autostart option: %s", err)
	}

	appAutostart, err := processAppAutostartOptions(apps)
	if err != nil {
		return fmt.Errorf("error processing global autostart option: %s", err)
	}

	for _, app := range apps {
		autostart := globalAppAutostart[app]
		// app setting takes precedence
		if appAutostart[app] != nil {
			autostart = appAutostart[app]
		}

		if autostart != nil {
			if *autostart {
				log.Infof("%s: will start and enable", app)
				err = snapctl.Start(env.SnapName + "." + app).Enable().Run()
				if err != nil {
					return fmt.Errorf("error starting service: %s", err)
				}
			} else {
				log.Infof("%s: will stop and disable", app)
				err = snapctl.Stop(env.SnapName + "." + app).Disable().Run()
				if err != nil {
					return fmt.Errorf("error stopping service: %s", err)
				}
			}
		}
	}

	return nil
}
