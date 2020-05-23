// Package plugins provides a source for cchat services as Go plugins. This
// package looks in UserConfigDir()/cchat/plugins/ by default.
//
// Usage
//
// The package can easily be used by just dash importing it:
//
//    _ "github.com/diamondburned/cchat/services/plugins"
//
package plugins

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"

	"github.com/diamondburned/cchat/services"
	"github.com/pkg/errors"
)

var pluginPath string

// SetPluginPath sets the plugin path before loading plugins. This only works
// until LoadPlugins is called.
func SetPluginPath(path string) {
	pluginPath = path
}

// TryConfigPath returns a path to $cfgDir/suffix. cfgDir is from
// os.UserConfigDir.
func TryConfigPath(suffix ...string) (string, error) {
	d, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(append([]string{d}, suffix...)...), nil
}

func init() {
	services.RegisterSource(loadPlugins)
}

func loadPlugins() (errs []error) {
	if pluginPath == "" {
		p, err := TryConfigPath("cchat", "plugins")
		if err != nil {
			errs = []error{errors.Wrap(err, "Failed to get config path")}
			return
		}
		pluginPath = p
	}

	d, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		// If the directory does not exist, then make one and exit.
		if os.IsNotExist(err) {
			if err := os.MkdirAll(pluginPath, 0755); err != nil {
				errs = []error{errors.Wrap(err, "Failed to make plugins dir")}
			}
			return
		}

		errs = []error{errors.Wrap(err, "Failed to read plugin path")}
		return
	}

	for _, f := range d {
		// We only need the plugin to call its init() function.
		_, err := plugin.Open(filepath.Join(pluginPath, f.Name()))
		if err != nil {
			errs = append(errs, errors.Wrap(err, "Failed to open plugin"))
			continue
		}
	}

	return
}
