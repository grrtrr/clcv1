/*
 * List the snapshots of a server.
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

	snapshots, err := client.GetSnapshots(flag.Arg(0), *acctAlias)
	if err != nil {
		exit.Fatalf("Failed to list the snapshotss of server %q: %s", flag.Arg(0), err)
	}

	fmt.Printf("Snapshots of %s:\n", flag.Arg(0))
	if len(snapshots) == 0 {
		println("No snapshots.")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoFormatHeaders(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(true)

		table.SetHeader([]string{ "Name", "Description", "Date Created" })
		for _, s := range snapshots {
			table.Append([]string{ s.Name, s.Description, fmt.Sprint(s.DateCreated) })
		}
		table.Render()
	}
}
