/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an existing Odin/OLS installation",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		content, err := os.ReadFile(OVM_CFG)
		if err != nil {
			fmt.Printf("Failed to read base path from config: %s\n", err)
			return
		}
		basePath := string(content)

		odinPath := fmt.Sprintf("%s/odin", basePath)
		olsPath := fmt.Sprintf("%s/odin-lsp", basePath)

		err = os.RemoveAll(odinPath)
		if err != nil {
			fmt.Printf("Failed to remove Odin directory: %s\n", err)
		} else {
			fmt.Println("Odin directory removed")
		}

		lspInstalled := checkExists(olsPath)
		err = os.RemoveAll(olsPath)
		if err != nil {
			fmt.Printf("Failed to remove OLS directory: %s\n", err)
		} else if lspInstalled {
			fmt.Println("OLS directory removed")
		}

		err = os.RemoveAll(OVM_CFG)
		if err != nil {
			fmt.Printf("Failed to remove OVM config file: %s\n", err)
		} else {
			fmt.Println("OVM config removed")
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
