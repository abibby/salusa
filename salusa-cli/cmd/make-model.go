/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// makeModelCmd represents the makeModel command
var makeModelCmd = &cobra.Command{
	Use:   "make:model",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(makeModelCmd)
}
