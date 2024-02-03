package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
	"github.com/tristanisham/clr"
)

type Config struct {
	basePath          string
	UseColor          bool
	ActiveVersion     string
	InstalledVersions []string
}

func (c *Config) save() error {
	serialized, err := toml.Marshal(&c)
	if err != nil {
		return fmt.Errorf("failed to serialize configuration %v", err)
	}

	if err := os.WriteFile(c.basePath, serialized, 0755); err != nil {
		return fmt.Errorf("failed to write config.toml file %v", err)
	}

	return nil
}

func (c *Config) ToggleColor() {
	if c.UseColor {
		c.DisableColor()
	} else {
		c.EnableColor()
	}
}

func (c *Config) DisableColor() {
	c.UseColor = false
	if err := c.save(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Output color disabled")
}

func (c *Config) EnableColor() {
	c.UseColor = true
	if err := c.save(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Output color %s\n", clr.Green("enabled"))
}

func (c *Config) AddInstalledVersion(version string) error {
	for _, v := range c.InstalledVersions {
		if v == version {
			return nil
		}
	}

	c.InstalledVersions = append(c.InstalledVersions, version)
	return c.save()
}

func (c *Config) RemoveInstalledVersion(version string) error {
	for i, v := range c.InstalledVersions {
		if v == version {
			c.InstalledVersions[i] = c.InstalledVersions[len(c.InstalledVersions)-1]
			c.InstalledVersions = c.InstalledVersions[:len(c.InstalledVersions)-1]
		}
	}
	return c.save()
}
