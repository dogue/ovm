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

	ovm_path := filepath.Join(home, ".ovm")
	if _, err := os.Stat(ovm_path); errors.Is(err, fs.ErrNotExist) {
		if verbose {
			fmt.Printf("OVM directory not found at `%s`, creating it now\n", ovm_path)
		}

		if err := os.MkdirAll(filepath.Join(ovm_path, "self"), 0775); err != nil {
			log.Fatal(err)
		}
	}

	ovm := &OVM{
		baseDir: ovm_path,
		Verbose: verbose,
	}
	ovm.Config.basePath = filepath.Join(ovm_path, "config.toml")

	if err := ovm.loadConfig(); err != nil {
		if errors.Is(err, ErrNoConfig) {
			if ovm.Verbose {
				fmt.Println("Config file not found. Creating default.")
			}

			ovm.Config = Config{
				UseColor: true,
			}

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
