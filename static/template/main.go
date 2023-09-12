package main

import (
	"log"

	"github.com/abibby/salusa/static/template/app"
)

func main() {
	err := app.Kernel.Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Kernel.Run()
	if err != nil {
		log.Fatal(err)
	}
}
