/*
 * Print the Administrator credentials of one server
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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <server-name>\n", path.Base(os.Args[0]))
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

	credentials, err := client.GetServerCredentials(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to obtain the credentials of server %q: %s", flag.Arg(0), err)
	}

	fmt.Printf("Credentials for %s:\n", flag.Arg(0))
	fmt.Printf("User:     %s\n", credentials.Username)
	fmt.Printf("Password: %s\n", credentials.Password)
}
