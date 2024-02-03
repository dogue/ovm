package cli

import (
	"fmt"

	"github.com/charmbracelet/log"
)

func (o *OVM) ListVersions(remote bool) error {
	if remote {
		versions, err := GetGitHubReleases("odin-lang", "odin")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Odin versions available for download:")
		for _, v := range versions {
			fmt.Printf("%s\n", v.TagName)
		}
	} else {
		fmt.Println("Odin versions installed locally (*active):")
		for _, v := range o.Config.InstalledVersions {
			if v == o.Config.ActiveVersion {
				fmt.Print("*")
			}
			fmt.Println(v)
		}
	}

	return nil
}
