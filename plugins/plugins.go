package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/plumbie/plumbie/config"

	log "github.com/sirupsen/logrus"
)

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

	for _, file := range files {
		plug, err := plugin.Open(file)
		if err != nil {
			return err
		}
		log.Debugf("plugins: Plugin loaded: %s", file)

		application, err := plug.Lookup("Application")
		if err != nil {
			return err
		}
		log.Debugf("plugins: Application symbol loaded: %s", application)
	}
	return nil
}
