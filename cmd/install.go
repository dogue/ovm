/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	ODIN_URL = "https://github.com/odin-lang/Odin"
	OLS_URL  = "https://github.com/DanielGavin/ols"
)

type Tool struct {
	name        string
	url         string
	buildScript string
	path        string
	commitHash  string
}

var ovmConfig = path.Join(xdg.ConfigHome, "ovm")

var Odin = Tool{
	name:        "Odin",
	url:         "https://github.com/odin-lang/Odin",
	buildScript: "build_odin.sh",
}

var Ols = Tool{
	name:        "OLS",
	url:         "https://github.com/DanielGavin/ols",
	buildScript: "build.sh",
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Odin/OLS",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		basePath, _ := cmd.Flags().GetString("path")
		odinVersion, _ := cmd.Flags().GetString("odin-version")
		force, _ := cmd.Flags().GetBool("force")
		installLsp, _ := cmd.Flags().GetBool("lsp")

		basePath, _ = filepath.Abs(basePath)
		Odin.path = filepath.Join(basePath, "odin")
		Ols.path = filepath.Join(basePath, "odin-lsp")

		err := os.WriteFile(ovmConfig, []byte(basePath), 0777)
		if err != nil {
			fmt.Printf("Failed to write base path to OVM config: %s", err)
			return
		}

		if checkExists(Odin.path) {
			if !force {
				fmt.Printf("Odin directory '%s' exists. Use -f/--force to overwrite\n", Odin.path)
				return
			}

			os.RemoveAll(Odin.path)
		}

		if odinVersion != "HEAD" {
			Odin.commitHash = odinVersion
		}

		cloneRepo(Odin)
		buildTool(Odin)

		if installLsp {
			olsVersion, _ := cmd.Flags().GetString("ols-version")

			if checkExists(Ols.path) {
				if !force {
					fmt.Printf("OLS directory '%s' exists. Use -f/--force to overwrite\n", Ols.path)
					return
				}

				os.RemoveAll(Ols.path)
			}

			if olsVersion != "HEAD" {
				Ols.commitHash = olsVersion
			}

			cloneRepo(Ols)
			buildTool(Ols)
		}

		printPathHelp(installLsp)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringP("path", "p", xdg.DataHome, "Install path prefix for Odin and OLS")
	installCmd.Flags().StringP("odin-version", "c", plumbing.HEAD.String(), "Specific commit for Odin to install")
	installCmd.Flags().BoolP("lsp", "l", false, "Install the Odin Language Server alongside Odin")
	installCmd.Flags().StringP("ols-version", "s", plumbing.HEAD.String(), "Specific commit for OLS to install")
	installCmd.Flags().BoolP("force", "f", false, "Force overwriting an existing installation")
}

func cleanUp(path string) {
	fmt.Printf("Cleaning up in %s\n", path)

	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Failed to clean up install directory: %s\n", err)
	}

	os.Exit(1)
}

func cloneRepo(tool Tool) {
	cloneOpts := &git.CloneOptions{
		URL:           tool.url,
		ReferenceName: plumbing.Master,
		Progress:      os.Stdout,
	}

	fmt.Printf("Cloning %s repository\n", tool.name)
	defer fmt.Println("Done!")
	repo, err := git.PlainClone(tool.path, false, cloneOpts)
	if err != nil {
		fmt.Printf("Failed to clone repository: %s\n", err)
		cleanUp(tool.path)
	}

	// no need to checkout if using latest, bail
	if tool.commitHash == "" {
		return
	}

	hash := plumbing.NewHash(tool.commitHash)
	commit, err := repo.CommitObject(hash)
	if err != nil {
		fmt.Printf("Failed to find specified commit: %s\n", err)
		cleanUp(tool.path)
	}

	wt, err := repo.Worktree()
	if err != nil {
		fmt.Printf("Failed to retrieve worktree: %s\n", err)
		cleanUp(tool.path)
	}

	checkOutOpts := &git.CheckoutOptions{
		Hash: commit.Hash,
	}

	err = wt.Checkout(checkOutOpts)
	if err != nil {
		fmt.Printf("Failed to checkout specified commit: %s\n", err)
		cleanUp(tool.path)
	}
}

func checkExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func buildOdin(path string) {
	cwd, _ := os.Getwd()
	os.Chdir(path)
	defer os.Chdir(cwd)

	build := exec.Command("./build_odin.sh")

	fmt.Println("Building Odin")

	err := build.Run()
	if err != nil {
		fmt.Printf("Failed to build Odin: %s\n", err)
		return
	}

	fmt.Println("Done!")
}

func buildOLS(path string) {
	cwd, _ := os.Getwd()
	os.Chdir(path)
	defer os.Chdir(cwd)
	build := exec.Command("./build.sh")

	fmt.Println("Building OLS")

	err := build.Run()
	if err != nil {
		fmt.Printf("Failed to build OLS: %s\n", err)
		return
	}

	fmt.Println("Done!")
}

func buildTool(tool Tool) {
	cwd, _ := os.Getwd()
	os.Chdir(tool.path)
	defer os.Chdir(cwd)

	cmdStr := filepath.Join(".", tool.buildScript)
	build := exec.Command(cmdStr)

	fmt.Printf("Building %s\n", tool.name)

	err := build.Run()
	if err != nil {
		fmt.Printf("Failed to build %s: %s\n", tool.name, err)
		return
	}

	fmt.Println("Done!")
}

func printPathHelp(lspInstalled bool) {
	fmt.Println("Add the following to your shell config (.bash_profile, .zshrc, etc):")
	fmt.Printf("export PATH=$PATH:%s", Odin.path)

	if lspInstalled {
		fmt.Printf(":%s", Ols.path)
	}

	fmt.Println()
}
