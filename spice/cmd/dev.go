/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/cosmtrek/air/runner"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// devCmd represents the run command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		r, err := runner.NewEngineWithConfig(&runner.Config{}, false)
		if err != nil {
			log.Fatal(err)
			return
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
		// cfg, err := util.LoadConfig("./")
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// watcher, changes := devWatch(cfg.Root)
		// defer watcher.Close()

		// fmt.Println("Watching")

		// for {
		// 	c := exec.Command("go", append([]string{"run", "-tags", "dev", "main.go"}, args...)...)
		// 	c.Stdout = os.Stdout
		// 	// c.Stdin = os.Stdin
		// 	c.Stderr = os.Stderr
		// 	err := c.Start()
		// 	if err != nil {
		// 		log.Print(err)
		// 	}

		// 	<-changes

		// 	c.Process.Kill()

		// 	fmt.Println("Change detected restarting server")

		// 	err = c.Wait()
		// 	if err != nil {
		// 		log.Print(err)
		// 	}
		// }
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
}

func devWatch(root string) (*fsnotify.Watcher, chan struct{}) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// defer .Close()

	changes := make(chan struct{})

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op.Has(fsnotify.Chmod) {
					continue
				}
				changes <- struct{}{}

				if event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
					dir, err := isDir(event.Name)
					if err != nil {
						log.Printf("failed to check dir status %s: %v", event.Name, err)
						continue
					}
					if !dir {
						continue
					}
					if event.Has(fsnotify.Create) {
						watcher.Add(event.Name)
					}
					if event.Has(fsnotify.Remove) {
						watcher.Remove(event.Name)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = filepath.Walk(root, func(p string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		return watcher.Add(p)
	})
	if err != nil {
		log.Fatal(err)
	}

	return watcher, changes
}

func isDir(p string) (bool, error) {
	fileInfo, err := os.Stat(p)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
