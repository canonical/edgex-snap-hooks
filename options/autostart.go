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
		autostart := options.Apps[app].Autostart
		if autostart == nil {
			// no autostart option for the app
			continue
		}

		bValue, err := parseAutostart(*autostart)
		if err != nil {
			return nil, fmt.Errorf("error parsing autostart for %s: %v", app, err)
		}
		if bValue != nil {
			appAutostart[app] = *&bValue
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
		bValue, err := parseAutostart(autostart)
		if err != nil {
			return nil, fmt.Errorf("error parsing autostart for %s: %v", app, err)
		}
		if bValue != nil {
			appAutostart[app] = *&bValue
		}
	}

	return appAutostart, nil
}

func parseAutostart(value string) (*bool, error) {
	value = strings.ToLower(value)
	switch value {
	case "":
		return nil, nil
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

// ProcessAutoStart will start and enable the listed app(s)
// based on the value of autostart snap option
func ProcessAutoStart(apps ...string) error {

	globalAppAutostart, err := processGlobalAutostartOptions(apps)
	if err != nil {
		return fmt.Errorf("error processing global autostart option: %s", err)
	}

	appAutostart, err := processGlobalAutostartOptions(apps)
	if err != nil {
		return fmt.Errorf("error processing global autostart option: %s", err)
	}

	for _, app := range apps {
		autostart := globalAppAutostart[app]
		if a, found := appAutostart[app]; found {
			autostart = a
		}
		log.Debugf("app %s: autostart=%t", app, appAutostart[app])
		if autostart != nil && *autostart {
			err = snapctl.Start(env.SnapName + "." + app).Enable().Run()
			if err != nil {
				return fmt.Errorf("error starting service: %s", err)
			}
		}
	}

	return nil
}
