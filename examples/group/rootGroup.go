/*
 * Prints the root Hardware Group for a given location
 */
package main

import (
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"path"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Location>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	rootGroup, err := client.GetRootGroup(flag.Arg(0))
	if err != nil {
		exit.Fatalf("Failed to list root group: %s", err)
	}

	fmt.Println("Location:   ", flag.Arg(0))
	fmt.Println("Root Group: ", rootGroup.Name)
	fmt.Println("UUID:       ", rootGroup.UUID)
}
