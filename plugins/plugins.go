package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/plumbie/plumbie/config"
	"github.com/plumbie/plumbie/sdk"

	log "github.com/sirupsen/logrus"
)

var Apps []sdk.Application

func Load() error {
	var files []string
	err := filepath.Walk(config.Plugins.Path, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("plugins: Invalid path")
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	Apps = []sdk.Application{}

	for _, file := range files {
		plug, err := plugin.Open(file)
		if err != nil {
			return err
		}
		log.Debugf("plugins: Plugin loaded: %s", file)

		symApplication, err := plug.Lookup("Application")
		if err != nil {
			return err
		}

		application, ok := symApplication.(sdk.Application)
		if !ok {
			log.Errorf("plugins: Application %s is not of type sdk.Application", symApplication)
		}

		log.Debugf("plugins: Application symbol loaded: %s", application)
		Apps = append(Apps, application)
	}
	return nil
}
