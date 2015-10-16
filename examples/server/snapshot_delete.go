/*
 * Delete a named snapshot for a specified server
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
	var acctAlias = flag.String("a", "", "Account alias to use")
	var snapName  = flag.String("s", "", "The name of the Snapshot to delete")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <server-name>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 || *snapName == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	err = client.DeleteSnapshot(*snapName, flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to delete snapshot on %s: %s", flag.Arg(0), err)
	}

	fmt.Printf("Successfully deleted snapshot %q from %s.\n", *snapName, flag.Arg(0))
}
