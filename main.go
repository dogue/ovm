package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"ovm/cli"
	"ovm/cli/meta"

	// "strings"

	"github.com/charmbracelet/log"
	flag "github.com/spf13/pflag"
	"github.com/tristanisham/clr"
	// "github.com/tristanisham/clr"
)

//go:embed help.txt
var helpText string

func main() {
	args := os.Args[1:]
	if _, ok := os.LookupEnv("OVM_DEBUG"); ok {
		log.SetLevel(log.DebugLevel)
	}

	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	installFlagSet := flag.NewFlagSet("install", flag.ExitOnError)
	installLsp := flag.BoolP("lsp", "l", false, "Specify if OLS should be installed with Odin")
	installFlagSet.AddFlag(flag.ShorthandLookup("l"))

	lsFlagSet := flag.NewFlagSet("ls", flag.ExitOnError)
	lsRemote := flag.BoolP("remote", "r", false, "List Odin versions available for download")
	lsFlagSet.AddFlag(flag.ShorthandLookup("r"))

	verboseMode := flag.BoolP("verbose", "v", false, "Show extra output during operations")
	flag.Parse()

	ovm := cli.Initialize(*verboseMode)
	args = flag.Args()

	for i, arg := range args {
		switch arg {

		case "install", "i":
			installFlagSet.Parse(args[i+1:])

			// default to "latest" if no version given
			var requestedVersion string
			if len(args) > i+1 {
				requestedVersion = installFlagSet.Arg(0)
			} else {
				requestedVersion = "latest"
			}

			targetVersion := cli.ValidateTargetVersion(requestedVersion)

			if ovm.Verbose {
				var outVer string
				if ovm.Config.UseColor {
					outVer = clr.Green(targetVersion.Tag)
				} else {
					outVer = targetVersion.Tag
				}

				fmt.Printf("Installing version %s...\n", outVer)
			}

			if err := ovm.Install(targetVersion, *installLsp); err != nil {
				log.Fatal(err)
			}
			return

		case "use", "switch":
			if len(args) > i+1 {
				version := args[i+1]
				if err := ovm.Use(version); err != nil {
					log.Fatal(err)
				}
			}
			return

		case "ls", "list":
			lsFlagSet.Parse(args[i+1:])
			err := ovm.ListVersions(*lsRemote)
			if err != nil {
				log.Warn(err)
			}

		case "remove", "rm":
			if len(args) > i+1 {
				err := ovm.Uninstall(args[i+1])
				if err != nil {
					log.Warn(err)
				}
			}

		case "upgrade", "u":
			if err := ovm.Upgrade(); err != nil {
				log.Fatal(err)
			}

		case "colors":
			var prompt string
			if ovm.Config.UseColor {
				prompt = fmt.Sprintf("Colors are currently %s. Would you like to disable them? [y/n]", ovm.Colored("enabled", "green"))
			} else {
				prompt = "Colors are currently disabled. Would you like to enable them? [y/n]"
			}

			fmt.Println(prompt)
			if cli.GetConfirmation() {
				ovm.Config.ToggleColor()
			}

		case "version":
			fmt.Printf("OVM %s\n", meta.VERSION)
			return

		case "help":
			printHelp()
			return

		default:
			log.Fatalf("invalid argument %q. Have a look at `ovm help` for usage.\n", arg)
		}
	}

}

func printHelp() {
	helpTemplate, err := template.New("help").Parse(helpText)
	if err != nil {
		fmt.Printf("Rendering error (%q). Version: %s\n", err, meta.VERSION)
		fmt.Println(helpText)
		return
	}

	if err := helpTemplate.Execute(os.Stdout, map[string]string{"Version": meta.VERSION}); err != nil {
		fmt.Printf("Rendering error (%q). Version: %s\n", err, meta.VERSION)
		fmt.Println(helpText)
		return
	}
}
