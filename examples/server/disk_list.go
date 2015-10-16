/*
 * List the disks of a given server
 */
package main

import (
	"github.com/olekukonko/tablewriter"
	"github.com/grrtrr/clcv1"
	"github.com/grrtrr/exit"
	"path"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {
	var acctAlias = flag.String("a",    "",    "Account alias to use")
	var diskNames = flag.Bool("q",      false, "Enable to list disk mount points / drive letters.")

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

	hasSnapshot, disks, err := client.ListDisks(flag.Arg(0), *acctAlias, *diskNames)
	if err != nil {
		exit.Fatalf("Failed to list disks of server %q: %s", flag.Arg(0), err)
	}

	fmt.Printf("Server %s (", flag.Arg(0))
	if hasSnapshot {
		fmt.Printf("has")
	} else {
		fmt.Printf("no")
	}
	fmt.Printf(" snapshot) disks:\n")


	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(true)

	table.SetHeader([]string{ "Name", "Size/GB", "SCSI Bus ID", "SCSI Dev ID" })
	for _, d := range disks {
		table.Append([]string{ d.Name, fmt.Sprint(d.SizeGB), d.ScsiBusID, d.ScsiDeviceID })
	}
	table.Render()
}
