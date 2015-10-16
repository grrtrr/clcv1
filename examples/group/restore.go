/*
 * Restore an archived Hardware Group.
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
	var acctAlias  = flag.String("a",      "", "Account alias to use")
	var parentUuid = flag.String("parent", "", "Parent group UUID to restore the group into")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <HW Group UUID>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 || *parentUuid == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	reqId, err := client.RestoreHardwareGroup(flag.Arg(0), *parentUuid, *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to restore Hardware Group %s: %s", flag.Arg(0), err)
	}

	fmt.Println("Request ID for group restoration:", reqId)
}
