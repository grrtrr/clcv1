/*
 * Convert template into a server
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
	var acctAlias = flag.String("a",    "", "Account alias to use")
	var hwGrpUUID = flag.String("u",    "", "UUID of the Hardware Group to place the converted server in")
	var password  = flag.String("pass", "", "New administrator/root password for the converted server")
	var network   = flag.String("net",  "", "Name of the network to use for the converted server")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <server-name>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 || *password == "" || *hwGrpUUID == "" || *network == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	reqId, err := client.ConvertTemplateToServer(flag.Arg(0), *password, *hwGrpUUID, *network, *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to generate a server from %s: %s", flag.Arg(0), err)
	}

	fmt.Println("Request ID for converting template:", reqId)
}
