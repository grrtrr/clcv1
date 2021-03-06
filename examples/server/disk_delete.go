/*
 * Delete a disk from a server
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
	var busId     = flag.String("bus", "", "The SCSI bus ID of the disk")
	var devId     = flag.String("dev", "", "The SCSI device ID of the disk")
	var force     = flag.Bool("force", false, "Override safety checks required for primary/OS disks")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]  <server-name>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 || *busId == "" || *devId == "" {
		flag.Usage()
		os.Exit(1)
	}

	client, err := clcv1.NewClient(log.New(os.Stdout, "", log.LstdFlags | log.Ltime))
	if err != nil {
		exit.Fatal(err.Error())
	} else if err := client.Logon("", ""); err != nil {
		exit.Fatalf("Login failed: %s", err)
	}

	reqId, err := client.DeleteDisk(flag.Arg(0), *acctAlias, *busId, *devId, *force)
	if err != nil {
		exit.Fatalf("Failed to delete disk on %s: %s", flag.Arg(0), err)
	}

	fmt.Println("Request ID for server disk deletion:", reqId)
}
