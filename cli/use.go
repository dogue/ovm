package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (o *OVM) Use(version string) error {
	targetPath := filepath.Join(o.baseDir, version)
	var err error

	if _, err = os.Stat(targetPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("It looks like %s isn't installed. Would you like to install it? [y/n]\n", version)
		if GetConfirmation() {
			targetVersion := ValidateTargetVersion(version)
			err = o.Install(targetVersion, false)
		} else {
			return fmt.Errorf("Version %s is not installed", version)
		}
	}

	if err != nil {
		return err
	}
	return o.setBin(version)
}

func (o *OVM) setBin(version string) error {
	targetPath := filepath.Join(o.baseDir, version, "odin")
	o.createSymlink(targetPath, "bin")

	o.linkCollections(version)

	o.Config.ActiveVersion = version
	if err := o.Config.save(); err != nil {
		return err
	}

	fmt.Printf("Active version set to %s\n", o.Colored(version, "green"))

	return nil
}

func GetConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	answer := strings.TrimSpace(strings.ToLower(text))
	return answer == "y" || answer == "ye" || answer == "yes"

}
