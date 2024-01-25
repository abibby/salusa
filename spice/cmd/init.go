/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path"

	"github.com/abibby/salusa/static"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <package>",
	Short: "Initialize a new salusa repo",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := copyDir(static.Content, "template", ".", args[0])
		if err != nil {
			return err
		}

		err = exec.Command("git", "init").Run()
		if err != nil {
			return err
		}

		err = exec.Command("go", "mod", "init", args[0]).Run()
		if err != nil {
			return err
		}

		// TODO: set git origin

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func copyDir(root fs.FS, src, dist, pkgPath string) error {
	files, err := fs.ReadDir(root, src)
	if err != nil {
		return err
	}
	for _, f := range files {
		srcPath := path.Join(src, f.Name())
		distPath := path.Join(dist, f.Name())
		if f.IsDir() {
			err = os.MkdirAll(distPath, 0755)
			if err != nil {
				return err
			}
			err = copyDir(root, srcPath, distPath, pkgPath)
			if err != nil {
				return err
			}
		} else {
			b, err := fs.ReadFile(root, srcPath)
			if err != nil {
				return err
			}
			b = bytes.ReplaceAll(b, []byte("github.com/abibby/salusa/static/template"), []byte(pkgPath))

			err = os.WriteFile(distPath, b, 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
