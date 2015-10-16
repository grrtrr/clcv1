/*
 * Convert server into a template
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
	var acctAlias  = flag.String("a",     "",      "Account alias to use")
	var password   = flag.String("pass",  "",      "Administrator/root password for the server to convert")
	var templAlias = flag.String("alias", "TEMPL", "The alias for the Template to create")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <server-name>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 || *password == "" || *templAlias == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	reqId, err := client.ConvertServerToTemplate(flag.Arg(0), *password, *templAlias, *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to convert %s into a template: %s", flag.Arg(0), err)
	}

	fmt.Println("Request ID for converting server:", reqId)
}
