package main

import (
	"flag"
	"fmt"
)

func main() {
	var mode *string = flag.String("mode", "server", "start mode whether is batch or server. ")
	flag.Parse()

	if *mode == "batch" {
		fmt.Println("batch mode")
	} else if *mode == "server" {
		fmt.Println("server mode")
	} else {
		panic(`
	   Usage:
	      -mode [batch or server]
	   `)
	}
}
