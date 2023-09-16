/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path"

	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/spice/pkg"
	"github.com/abibby/salusa/spice/util"
	"github.com/spf13/cobra"
)

// makeMigrationCmd represents the makeMigration command
var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration",
	Short: "",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := util.Name(args)

		root, _, err := util.PackageRoot(".")
		if err != nil {
			return err
		}

		packageName := "migrations"
		migrationsDir := path.Join(root, packageName)

		err = os.MkdirAll(migrationsDir, 0755)
		if err != nil {
			return err
		}
		run := "schema.Run(func(ctx context.Context, tx builder.QueryExecer) error {\n" +
			"return nil\n" +
			"})"
		src, err := migrate.SrcFile(name, packageName, pkg.Raw(run), pkg.Raw(run))
		if err != nil {
			return err
		}

		return os.WriteFile(path.Join(migrationsDir, name+".go"), []byte(src), 0644)
	},
}

func init() {
	rootCmd.AddCommand(makeMigrationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// makeMigrationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// makeMigrationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
