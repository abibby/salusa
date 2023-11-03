/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// devCmd represents the run command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("go", append([]string{"run", "-tags", "dev", "main.go"}, args...)...)
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr

		err := c.Run()
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
