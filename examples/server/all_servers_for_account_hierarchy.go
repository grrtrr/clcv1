/*
 * Deep list of all Servers, grouped by account alias.
 */
package main

import (
	"github.com/grrtrr/clcv1/utils"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"flag"
	"log"
_	"fmt"
	"os"
)

func main() {
	var acctAlias = flag.String("a", "", "Account alias of the account that owns the servers")
	var location  = flag.String("l", "", "The data center location")
	var simple    = flag.Bool("simple", false, "Use simple (debugging) output format")

	flag.Parse()

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	servers, err := client.GetAllServersForAccountHierarchy(*acctAlias, *location)
	if err != nil {
		exit.Fatalf("Failed to list all servers: %s", err)
	}

	// FIXME: simple representation only, since currently not able to test
	//       (StatusCode 2, null result)
	*simple = true
	if len(servers) == 0 {
		println("Empty result.")
	} else if *simple {
		for _, s := range servers {
			utils.PrintStruct(s)
		}
	}
}
