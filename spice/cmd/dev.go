/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmtrek/air/runner"
	"github.com/spf13/cobra"
)

// devCmd represents the run command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		cfg, err := runner.InitConfig("")
		if err != nil {
			return err
		}

		r, err := runner.NewEngineWithConfig(cfg, false)
		if err != nil {
			return err
		}
		go func() {
			<-sigs
			r.Stop()
		}()

		defer func() {
			if e := recover(); e != nil {
				log.Fatalf("PANIC: %+v", e)
			}
		}()

		r.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}
