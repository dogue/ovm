/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing Odin/OLS installation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		content, err := os.ReadFile(OVM_CFG)
		if err != nil {
			fmt.Printf("Failed to read base path from OVM config: %s\n", err)
			return
		}
		basePath := string(content)

		odinPath := path.Join(basePath, "odin")
		olsPath := path.Join(basePath, "odin-lsp")

		updateOdin(odinPath)
		updateOLS(olsPath)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

}

func updateOdin(path string) {
	fmt.Println("Pulling latest changes for Odin")
	if !pullRepo(path) {
		return
	}

	buildOdin(path)
}

func updateOLS(path string) {
	fmt.Println("Pulling latest changes for OLS")
	if !pullRepo(path) {
		return
	}

	buildOLS(path)
}

func pullRepo(path string) bool {
	repo, err := git.PlainOpen(path)
	if err != nil {
		fmt.Printf("Failed to open repository: %s\n", err)
		return false
	}

	wt, err := repo.Worktree()
	if err != nil {
		fmt.Printf("Failed to get working tree: %s\n", err)
		return false
	}

	err = wt.Pull(&git.PullOptions{Progress: os.Stdout})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("Already up to date")
			return false
		} else {
			fmt.Printf("Failed to pull latest changes: %s\n", err)
			return false
		}
	}

	return true
}
