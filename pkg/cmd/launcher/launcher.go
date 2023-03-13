package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rocketblend/rocketblend-launcher/pkg/cmd/launcher/config"
)

var Name = "rocketblend-launcher"

func Launch() error {
	config, err := config.Load(Name)
	if err != nil {
		return fmt.Errorf("error loading config: %s", err)
	}

	launch := config.GetString("previous")
	if len((os.Args)) > 1 {
		launch = os.Args[1]
	}

	path := filepath.Dir(launch)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path '%s' does not exist", path)
	}

	cmd := exec.Command("rocketblend", "install", "-d", path)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error installing packages: %s", err)
	}

	cmd = exec.Command("rocketblend", "start", "-d", path)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting file: %s", err)
	}

	config.Set("previous", launch)

	return nil
}
