/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"

	"github.com/abibby/salusa/spice/util"
	"github.com/spf13/cobra"
)

var srcMain = `package main

import (
	"errors"
	"log"
	"os"

	"github.com/abibby/salusa/database/migrate"
	%#v
	%#v
)

func main() {
	m := %s.Use()

	src, err := m.GenerateMigration(%#v, %#v, &%s.%s{})
	if errors.Is(err, migrate.ErrNoChanges) {
		return
	} else if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(%#v, []byte(src), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
`

var srcMigrations = `package %s

import (
	"github.com/abibby/salusa/database/migrate"
)

var migrations = migrate.New()

func Use() *migrate.Migrations {
	return migrations
}
`

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate:migration",
	Short: "Run from go generate",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := util.PkgInfo(".")
		if err != nil {
			return err
		}

		packageName := "migrations"
		migrationsDir := path.Join(info.PackageRoot, packageName)

		err = os.MkdirAll(migrationsDir, 0755)
		if err != nil {
			return err
		}
		err = os.WriteFile(path.Join(migrationsDir, "migrations.go"), []byte(fmt.Sprintf(srcMigrations, packageName)), 0644)
		if err != nil {
			return err
		}

		modelPackage := os.Getenv("GOPACKAGE")
		modelFile := os.Getenv("GOFILE")
		modelLineStr := os.Getenv("GOLINE")
		modelLine, err := strconv.Atoi(modelLineStr)
		if err != nil {
			return err
		}

		b, err := os.ReadFile(modelFile)
		if err != nil {
			return err
		}

		line := bytes.Split(b, []byte("\n"))[modelLine]

		matches := regexp.MustCompile(`type ([A-Z]\w+) struct`).FindSubmatch(line)
		if len(matches) < 2 {
			return fmt.Errorf("could not find model struct")
		}
		name := util.Name([]string{string(matches[1])})
		migrationFile := path.Join(migrationsDir, name+".go")
		migrationsImport := path.Join(info.RootPackage, packageName)
		src := []byte(fmt.Sprintf(srcMain, migrationsImport, info.CurrentPackage, packageName, name, packageName, modelPackage, matches[1], migrationFile))
		// fmt.Printf("%s\n", src)

		tmp := os.TempDir()
		outFile := path.Join(tmp, fmt.Sprintf("bob-generate-main-%s.go", name))

		err = os.WriteFile(outFile, src, 0644)
		if err != nil {
			return err
		}
		defer os.Remove(outFile)
		err = run("go", "run", outFile)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	// cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
