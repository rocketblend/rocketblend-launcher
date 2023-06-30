package launcher

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/hashicorp/go-version"
	"github.com/rocketblend/rocketblend-launcher/pkg/cmd/launcher/config"
)

type launcher struct {
	projectPath        string
	rocketBlendVersion string
}

const (
	Name               = "rocketblend-launcher"
	Alias              = "rocketblend"
	RocketBlendVersion = "0.8.0"
)

func New(projectPath string) *launcher {
	return &launcher{
		projectPath:        projectPath,
		rocketBlendVersion: RocketBlendVersion,
	}
}

func (l *launcher) Launch() error {
	if err := l.WithSpinner("Checking if Rocketblend is available... ", l.checkAvailablity); err != nil {
		fmt.Println("Rocketblend is not found. Please ensure it is installed and available in PATH.")
		fmt.Println("You can download Rocketblend from: https://github.com/rocketblend/rocketblend/releases/latest")
		return err
	}

	if err := l.WithSpinner("Checking if Project is valid...", l.checkProject); err != nil {
		fmt.Println("Invalid project path. Please ensure the project path is valid.")
		return err
	}

	if err := l.WithSpinner("Starting Rocketblend...", l.startProject); err != nil {
		fmt.Println("Failed to start Rocketblend.")
		return err
	}

	if err := l.WithSpinner("Updating config...", l.updateConfig); err != nil {
		fmt.Println("Failed to update config.")
		return err
	}

	fmt.Println("Rocketblend started successfully!")

	return nil
}

func (l *launcher) updateConfig() error {
	config, err := config.Load(Name)
	if err != nil {
		return err
	}

	config.Set("previous", l.projectPath)
	if err := config.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func (l *launcher) checkProject() error {
	config, err := config.Load(Name)
	if err != nil {
		return err
	}

	projectPath := l.projectPath
	if projectPath == "" {
		projectPath = config.GetString("previous")
	}

	if err := l.isValidProjectPath(filepath.Dir(projectPath)); err != nil {
		return err
	}

	l.projectPath = projectPath

	return nil
}

func (l *launcher) startProject() error {
	cmd := exec.Command(Alias, "run", "-d", filepath.Dir(l.projectPath))
	err := cmd.Start()
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)

	go func() {
		errChan <- cmd.Wait()
	}()

	go func() {
		err := <-errChan
		if err != nil {
			fmt.Printf("ERROR: %v", err)
		}
	}()

	return nil
}

func (l *launcher) checkAvailablity() error {
	cmd := exec.Command(Alias, "--version")
	out, err := cmd.Output()
	if err != nil {
		return errors.New("rocketblend not found")
	}

	versionStr := strings.TrimSpace(string(out))
	if versionStr == "dev" {
		return nil // Bypass version check for "dev" version
	}

	v, err := version.NewVersion(versionStr)
	if err != nil {
		return fmt.Errorf("failed to parse Rocketblend version: %w", err)
	}

	constraints, err := version.NewConstraint(">= " + l.rocketBlendVersion)
	if err != nil {
		return fmt.Errorf("failed to parse version constraint: %w", err)
	}

	if !constraints.Check(v) {
		return fmt.Errorf("rocketblend version is too old, minimum version required is %s", l.rocketBlendVersion)
	}

	return nil
}

func (l *launcher) isValidProjectPath(path string) error {
	// Check if path exists and is a directory
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		return fmt.Errorf("the provided path either does not exist or is not a directory: %v", err)
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
		return fmt.Errorf("an error occurred while walking the directory: %v", err)
	}

	if !blendFileFound {
		return errors.New("no .blend file found in the provided path")
	}

	if !rocketFileFound {
		return errors.New("no rocketfile.yaml file found in the provided path")
	}

	return nil
}

func (l *launcher) WithSpinner(prefix string, f func() error) error {
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.HideCursor = true
	spinner.Prefix = prefix
	spinner.Start()
	err := f()
	spinner.Stop()
	if err != nil {
		return err
	}

	return nil
}
