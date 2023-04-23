package launcher

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rocketblend/rocketblend-launcher/pkg/cmd/launcher/config"
)

const (
	Name  = "rocketblend-launcher"
	Alias = "rocketblend"
)

func Launch() error {
	log.Println("Checking if rocketblend is available...")

	if !isRocketBlendAvailable() {
		return fmt.Errorf("rocketblend is not available")
	}

	log.Println("Rocketblend is available!")

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

	log.Println("Starting project...")
	cmd := exec.Command(Alias, "start", "-d", path)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed to start project: %v, output: %s", err, output)
		return fmt.Errorf("failed to start project: %v", err)
	}

	log.Println("Project started successfully!")
	config.Set("previous", launch)
	log.Println("Updating last launched...")

	err = config.WriteConfig()
	if err != nil {
		return fmt.Errorf("error writing config: %s", err)
	}

	log.Println("Updated successfully!")
	return nil
}

func isRocketBlendAvailable() bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", Alias+".exe")
	} else {
		cmd = exec.Command("which", Alias)
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
