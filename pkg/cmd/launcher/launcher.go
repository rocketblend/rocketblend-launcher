package launcher

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rocketblend/rocketblend-launcher/pkg/cmd/launcher/config"
)

var Name = "rocketblend-launcher"

func Launch() error {
	fmt.Println("Checking if rocketblend is available...")

	if !isRocketBlendAvailable() {
		return fmt.Errorf("rocketblend is not available")
	}

	fmt.Println("Rocketblend is available!")

	config, err := config.Load(Name)
	if err != nil {
		return fmt.Errorf("error loading config: %s", err)
	}

	launch := config.GetString("previous")
	if len(os.Args) > 1 {
		launch = os.Args[1]
	}

	path := filepath.Dir(launch)
	if !isValidRocketPath(path) {
		return fmt.Errorf("invalid rocket path: %s", path)
	}

	fmt.Println("Starting project...")

	cmd := exec.Command("rocketblend", "start", "-d", path)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error starting project: %s", err)
	}

	fmt.Println("Project started successfully!")

	config.Set("previous", launch)

	fmt.Println("Updating last launched...")

	err = config.WriteConfig()
	if err != nil {
		return fmt.Errorf("error writing config: %s", err)
	}

	fmt.Println("Updated successfully!")

	return nil
}

func isRocketBlendAvailable() bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "rocketblend.exe")
	} else {
		cmd = exec.Command("which", "rocketblend")
	}
	err := cmd.Run()
	return err == nil
}

func isValidRocketPath(path string) bool {
	// Check if path exists and is a directory
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		return false
	}

	// Check if path contains a .blend file and a rocketfile.yaml file
	var blendFileFound, rocketFileFound bool
	err = filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		switch ext := filepath.Ext(path); ext {
		case ".blend":
			blendFileFound = true
		case ".yaml", ".yml":
			if strings.TrimSuffix(d.Name(), ext) == "rocketfile" {
				rocketFileFound = true
			}
		}
		return nil
	})

	if err != nil {
		return false
	}

	return blendFileFound && rocketFileFound
}
