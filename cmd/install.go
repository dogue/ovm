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

var OVM_CFG = path.Join(xdg.ConfigHome, "ovm")

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
		odinPath := fmt.Sprintf("%s/odin", basePath)
		olsPath := fmt.Sprintf("%s/odin-lsp", basePath)

		err := os.WriteFile(OVM_CFG, []byte(basePath), 0777)
		if err != nil {
			fmt.Printf("Failed to write base path to file at '%s': %s", OVM_CFG, err)
			return
		}

		if checkExists(odinPath) {
			if !force {
				fmt.Printf("Odin directory '%s' exists. Use -f/--force to overwrite\n", odinPath)
				return
			}

			os.RemoveAll(odinPath)
		}

		fmt.Println("Cloning Odin repository")
		cloneRepo(ODIN_URL, odinPath, odinVersion)

		buildOdin(odinPath)

		if installLsp {
			olsVersion, _ := cmd.Flags().GetString("ols-version")

			if checkExists(olsPath) {
				if !force {
					fmt.Printf("OLS directory '%s' exists. Use -f/--force to overwrite\n", olsPath)
					return
				}

				os.RemoveAll(olsPath)
			}

			fmt.Println("Cloning OLS respository")
			cloneRepo(OLS_URL, olsPath, olsVersion)

			buildOLS(olsPath)
		}

		printPathHelp(installLsp, odinPath, olsPath)
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

func cloneRepo(url string, path string, version string) {
	cloneOpts := &git.CloneOptions{
		URL:           url,
		Depth:         1,
		ReferenceName: plumbing.Master,
		Progress:      os.Stdout,
	}

	// regular clone with history if not using latest commit
	if version != "HEAD" {
		cloneOpts.Depth = 0
	}

	repo, err := git.PlainClone(path, false, cloneOpts)
	if err != nil {
		fmt.Printf("Failed to clone repository: %s\n", err)
		cleanUp(path)
	}

	// no need to checkout if using latest, bail
	if version == "HEAD" {
		return
	}

	hash := plumbing.NewHash(version)
	commit, err := repo.CommitObject(hash)
	if err != nil {
		fmt.Printf("Failed to find specified commit: %s\n", err)
		cleanUp(path)
	}

	wt, err := repo.Worktree()
	if err != nil {
		fmt.Printf("Failed to retrieve worktree: %s\n", err)
		cleanUp(path)
	}

	checkOutOpts := &git.CheckoutOptions{
		Hash: commit.Hash,
	}

	err = wt.Checkout(checkOutOpts)
	if err != nil {
		fmt.Printf("Failed to checkout specified commit: %s\n", err)
		cleanUp(path)
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

func printPathHelp(lspInstalled bool, odinPath string, olsPath string) {
	// odinPath, _ = filepath.Abs(odinPath)

	fmt.Println("Add the following to your shell config (.bash_profile, .zshrc, etc):")
	fmt.Printf("export PATH=$PATH:%s", odinPath)

	if lspInstalled {
		// olsPath, _ = filepath.Abs(olsPath)
		fmt.Printf(":%s", olsPath)
	}

	fmt.Println()
}
