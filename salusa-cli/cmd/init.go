/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new salusa repo",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := exec.Command("git", "clone", "https://github.com/abibby/salusa-template.git", ".").Run()
		if err != nil {
			return err
		}

		err = os.RemoveAll("./.git")
		if err != nil {
			return err
		}

		err = exec.Command("git", "init").Run()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
