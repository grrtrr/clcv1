/*
 * Change the Administrator credentials of a given server
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
	var acctAlias = flag.String("a",   "", "Account alias to use")
	var oldPasswd = flag.String("old", "", "The existing password (for authentication)")
	var newPasswd = flag.String("new", "", "The new password to apply")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <server-name>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 || *oldPasswd == "" || *newPasswd == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	err = client.ServerChangePassword(flag.Arg(0), *acctAlias, *oldPasswd, *newPasswd)
	if err != nil {
		exit.Fatalf("Failed to change the password on %q: %s", flag.Arg(0), err)
	}

	fmt.Printf("Successfully changed the password on %s to %s\n", flag.Arg(0), *newPasswd)
}
