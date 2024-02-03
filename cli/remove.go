package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func (o *OVM) Uninstall(version string) error {
	targetPath := filepath.Join(o.baseDir, version)

	if _, err := os.Stat(targetPath); err == nil {
		if err := os.RemoveAll(targetPath); err != nil {
			return err
		}

		if err := o.Config.RemoveInstalledVersion(version); err != nil {
			return err
		}

		fmt.Printf("✔ Uninstalled %s.\nRun `ovm ls` to view installed versions.\n", version)
		return nil
	}

	fmt.Printf("Version %s doesn't appear to be installed.\n", o.Colored(version, "red"))
	return o.ListVersions(false)
}
