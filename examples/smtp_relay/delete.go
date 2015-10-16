/*
 * Delete an existing SMTP relay alias
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
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <Relay Alias>\n", path.Base(os.Args[0]))
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

	if err = client.RemoveAlias(flag.Arg(0)); err != nil {
		exit.Fatalf("Failed to delete Relay Alias %s: %s", flag.Arg(0), err)
	}
	fmt.Printf("Successfully deleted SMTP Relay Alias %q\n", flag.Arg(0))
}
