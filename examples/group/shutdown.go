/*
 * Shut down the Hardware Group along with all child groups and servers.
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
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <HW Group UUID>\n", path.Base(os.Args[0]))
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

	reqId, err := client.ShutdownHardwareGroup(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to shut down hardware group: %s", err)
	}

	fmt.Println("Request ID for shutting down hardware group:", reqId)
}
