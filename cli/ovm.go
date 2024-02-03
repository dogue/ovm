package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

type OVM struct {
	baseDir string
	Verbose bool
	Config  Config
}

func Initialize(verbose bool) *OVM {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	ovmPath := filepath.Join(home, ".ovm")
	if _, err := os.Stat(ovmPath); errors.Is(err, fs.ErrNotExist) {
		if verbose {
			fmt.Printf("OVM directory not found at `%s`, creating it now\n", ovmPath)
		}

		if err := os.MkdirAll(filepath.Join(ovmPath, "self"), 0775); err != nil {
			log.Fatal(err)
		}
	}

	ovm := &OVM{
		baseDir: ovmPath,
		Verbose: verbose,
	}
	ovm.Config.basePath = filepath.Join(ovmPath, "config.toml")

	if err := ovm.loadConfig(); err != nil {
		if errors.Is(err, ErrNoConfig) {
			if ovm.Verbose {
				fmt.Println("Config file not found. Creating default.")
			}

			ovm.Config.UseColor = true

			if err := ovm.Config.save(); err != nil {
				log.Warn("Failed to create config.toml file", err)
			}
		}
	}

	return ovm
}

func (o *OVM) loadConfig() error {
	set_path := o.Config.basePath
	if _, err := os.Stat(set_path); errors.Is(err, os.ErrNotExist) {
		return ErrNoConfig
	}

	data, err := os.ReadFile(set_path)
	if err != nil {
		return err
	}

	return toml.Unmarshal(data, &o.Config)
}
